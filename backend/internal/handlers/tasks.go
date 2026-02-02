package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/cron-observer/backend/internal/deletequeue"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/middleware"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"github.com/yourusername/cron-observer/backend/internal/scheduler"
	"github.com/yourusername/cron-observer/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskHandler struct {
	repo      repositories.Repository
	eventBus  *events.EventBus
	scheduler interface {
		RegisterTask(ctx context.Context, task *models.Task) error
		UnregisterTask(taskUUID string)
		IsWithinGroupWindow(ctx context.Context, taskGroup *models.TaskGroup) bool
	}
	superAdminMap   map[string]bool
	deletePublisher deletequeue.DeleteJobPublisher // optional until wired in main
}

func NewTaskHandler(repo repositories.Repository, eventBus *events.EventBus, scheduler interface {
	RegisterTask(ctx context.Context, task *models.Task) error
	UnregisterTask(taskUUID string)
	IsWithinGroupWindow(ctx context.Context, taskGroup *models.TaskGroup) bool
}, superAdmins []string, deletePublisher deletequeue.DeleteJobPublisher) *TaskHandler {

	// Create a map for O(1) lookup
	superAdminMap := make(map[string]bool)
	for _, admin := range superAdmins {
		normalizedAdmin := strings.ToLower(strings.TrimSpace(admin))
		if normalizedAdmin != "" {
			superAdminMap[normalizedAdmin] = true
		}
	}

	return &TaskHandler{
		repo:            repo,
		eventBus:        eventBus,
		scheduler:       scheduler, // Can be nil if scheduler is not needed
		superAdminMap:   superAdminMap,
		deletePublisher: deletePublisher, // optional until wired in main
	}
}

// GetTasksByProject retrieves all tasks for a project
// @Summary      Get tasks by project
// @Description  Retrieve all tasks belonging to a project
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Success      200  {array}   models.Task
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/tasks [get]
func (h *TaskHandler) GetTasksByProject(c *gin.Context) {
	projectIDParam := c.Param("project_id")
	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}

	// Convert project_id to ObjectID
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project_id format in path",
		})
		return
	}

	// Get all tasks for this project
	tasks, err := h.repo.GetTasksByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get tasks for project",
		})
		return
	}

	if tasks == nil {
		tasks = []*models.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}

// CreateTask creates a new task
// @Summary      Create a new task
// @Description  Create a new scheduled task in a project
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task body models.CreateTaskRequest true "Task creation request"
// @Success      201  {object}  models.Task
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	// Get authenticated user from context
	user, exists := middleware.GetUserFromContext(c)
	if exists {
		log.Printf("User %s (%s) is creating a task", user.Name, user.Email)
	}

	var req models.CreateTaskRequest

	// Bind JSON and validate using Gin's binding
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Get project_id from path parameter
	projectIDParam := c.Param("project_id")
	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}

	// Convert project_id path parameter to ObjectID
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project_id format in path",
		})
		return
	}

	// Also validate project_id from request body matches path parameter
	reqProjectID, err := primitive.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project_id format in request body",
		})
		return
	}

	// Ensure project_id in body matches path parameter
	if reqProjectID != projectID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id in path and body must match",
		})
		return
	}

	// Set default status if not provided. Binding restricts client input to ACTIVE/DISABLED only (PENDING_DELETE/DELETE_FAILED are backend-only).
	status := req.Status
	if status == "" {
		status = models.TaskStatusActive
	}

	// Convert TaskGroupID if provided
	var taskGroupID *primitive.ObjectID
	if req.TaskGroupID != "" {
		groupID, err := primitive.ObjectIDFromHex(req.TaskGroupID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid task_group_id format",
			})
			return
		}
		taskGroupID = &groupID
	}

	// Calculate initial state based on task group window (if task belongs to a group)
	state := models.TaskStateNotRunning
	if taskGroupID != nil {
		// Get the task group to check its window
		taskGroup, err := h.repo.GetTaskGroupByID(c.Request.Context(), *taskGroupID)
		if err == nil && taskGroup != nil && taskGroup.StartTime != "" && taskGroup.EndTime != "" {
			// Check if current time is within the group's window
			// Note: We can't use scheduler here as it's an interface, so we'll calculate state after creation
			// For now, default to NOT_RUNNING - it will be updated by scheduler when group window starts
			state = models.TaskStateNotRunning
		}
	}

	// Convert request DTO to Task model
	task := &models.Task{
		ProjectID:    projectID,
		TaskGroupID:  taskGroupID,
		UUID:         uuid.New().String(),
		Name:         req.Name,
		Description:  req.Description,
		ScheduleType: req.ScheduleType,
		Status:       status,
		State:        state, // Set initial state
		ScheduleConfig: models.ScheduleConfig{
			CronExpression: req.ScheduleConfig.CronExpression,
			Timezone:       req.ScheduleConfig.Timezone,
			DaysOfWeek:     req.ScheduleConfig.DaysOfWeek,
			Exclusions:     req.ScheduleConfig.Exclusions,
		},
		TimeoutSeconds: req.TimeoutSeconds,
		Metadata:       req.Metadata,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Convert TimeRange if provided
	if req.ScheduleConfig.TimeRange != nil {
		task.ScheduleConfig.TimeRange = &models.TimeRange{
			Start: req.ScheduleConfig.TimeRange.Start,
			End:   req.ScheduleConfig.TimeRange.End,
			Frequency: &models.Frequency{
				Value: req.ScheduleConfig.TimeRange.Frequency.Value,
				Unit:  req.ScheduleConfig.TimeRange.Frequency.Unit,
			},
		}
	}

	// TriggerConfig is no longer required - tasks use project's execution_endpoint
	// Leave TriggerConfig empty/zero value for new tasks

	// Create the task
	err = h.repo.CreateTask(c.Request.Context(), projectIDParam, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task",
		})
		return
	}

	// Publish TaskCreated event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskCreated,
		Payload: events.TaskPayload{Task: task},
	})

	c.JSON(http.StatusCreated, task)
}

// UpdateTask updates an existing task
// @Summary      Update a task
// @Description  Update an existing scheduled task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_uuid path string true "Task UUID"
// @Param        task body models.UpdateTaskRequest true "Task update request"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/tasks/{task_uuid} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req models.UpdateTaskRequest

	// Bind JSON and validate using Gin's binding
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Get project_id and task_uuid from path parameters
	projectIDParam := c.Param("project_id")
	taskUUIDParam := c.Param("task_uuid")

	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}

	if taskUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task_uuid is required in path",
		})
		return
	}

	// Convert project_id to ObjectID
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project_id format in path",
		})
		return
	}

	// Check authorization: user must be admin in project or super admin
	if !RequireProjectAdmin(c, h.repo, projectID, h.superAdminMap) {
		return
	}

	// Get existing task to preserve UUID and timestamps
	existingTask, err := h.repo.GetTaskByUUID(c.Request.Context(), taskUUIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Set default status if not provided. Binding restricts client input to ACTIVE/DISABLED only (PENDING_DELETE/DELETE_FAILED are backend-only).
	status := req.Status
	if status == "" {
		status = existingTask.Status
	}

	// Handle TaskGroupID - preserve existing if not provided in request
	var taskGroupID *primitive.ObjectID
	if req.TaskGroupID != "" {
		groupID, err := primitive.ObjectIDFromHex(req.TaskGroupID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid task_group_id format",
			})
			return
		}
		taskGroupID = &groupID
	} else {
		// Preserve existing TaskGroupID if not provided
		taskGroupID = existingTask.TaskGroupID
	}

	// Determine task state
	// If status is being set to DISABLED, set state to NOT_RUNNING
	// If status is being set to ACTIVE and task belongs to an ACTIVE group within window, set state to RUNNING
	// Otherwise, preserve existing state (it's system-controlled)
	state := existingTask.State
	if status == models.TaskStatusDisabled {
		state = models.TaskStateNotRunning
	} else if status == models.TaskStatusActive && existingTask.Status != models.TaskStatusActive {
		// Status changed to ACTIVE - check if task belongs to an active group within window
		if taskGroupID != nil && h.scheduler != nil {
			taskGroup, err := h.repo.GetTaskGroupByID(c.Request.Context(), *taskGroupID)
			if err == nil && taskGroup != nil {
				if taskGroup.Status == models.TaskGroupStatusActive {
					// Check if group is within time window
					if h.scheduler.IsWithinGroupWindow(c.Request.Context(), taskGroup) {
						state = models.TaskStateRunning
					}
				}
			}
		}
	}

	// Update task fields
	task := &models.Task{
		ID:           existingTask.ID,
		UUID:         existingTask.UUID, // UUID cannot be changed
		ProjectID:    projectID,
		TaskGroupID:  taskGroupID,
		Name:         req.Name,
		Description:  req.Description,
		ScheduleType: req.ScheduleType,
		Status:       status,
		State:        state,
		ScheduleConfig: models.ScheduleConfig{
			CronExpression: req.ScheduleConfig.CronExpression,
			Timezone:       req.ScheduleConfig.Timezone,
			DaysOfWeek:     req.ScheduleConfig.DaysOfWeek,
			Exclusions:     req.ScheduleConfig.Exclusions,
		},
		TimeoutSeconds: req.TimeoutSeconds,
		Metadata:       req.Metadata,
		CreatedAt:      existingTask.CreatedAt, // Preserve original creation time
		UpdatedAt:      time.Now(),
	}

	// Convert TimeRange if provided
	if req.ScheduleConfig.TimeRange != nil {
		task.ScheduleConfig.TimeRange = &models.TimeRange{
			Start: req.ScheduleConfig.TimeRange.Start,
			End:   req.ScheduleConfig.TimeRange.End,
			Frequency: &models.Frequency{
				Value: req.ScheduleConfig.TimeRange.Frequency.Value,
				Unit:  req.ScheduleConfig.TimeRange.Frequency.Unit,
			},
		}
	}

	// TriggerConfig is no longer required - tasks use project's execution_endpoint
	// Preserve existing TriggerConfig if it exists, otherwise leave empty
	task.TriggerConfig = existingTask.TriggerConfig

	// Update the task
	err = h.repo.UpdateTask(c.Request.Context(), taskUUIDParam, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	// If status changed to DISABLED, update state and unregister cron job immediately
	if status == models.TaskStatusDisabled && existingTask.Status != models.TaskStatusDisabled {
		// Update state to NOT_RUNNING
		if err := h.repo.UpdateTaskState(c.Request.Context(), taskUUIDParam, models.TaskStateNotRunning); err != nil {
			log.Printf("Failed to update task %s state to NOT_RUNNING: %v", taskUUIDParam, err)
		}

		// Unregister task from scheduler immediately
		if h.scheduler != nil {
			h.scheduler.UnregisterTask(taskUUIDParam)
			log.Printf("Unregistered cron job for task %s (status set to DISABLED)", taskUUIDParam)
		}
	}

	// If status changed to ACTIVE, check if we need to register cron job and update state
	if status == models.TaskStatusActive && existingTask.Status != models.TaskStatusActive {
		// Check if task belongs to an active group within window
		if task.TaskGroupID != nil && h.scheduler != nil {
			taskGroup, err := h.repo.GetTaskGroupByID(c.Request.Context(), *task.TaskGroupID)
			if err == nil && taskGroup != nil {
				if taskGroup.Status == models.TaskGroupStatusActive {
					// Check if group is within time window
					if h.scheduler.IsWithinGroupWindow(c.Request.Context(), taskGroup) {
						// Update state to RUNNING
						if err := h.repo.UpdateTaskState(c.Request.Context(), taskUUIDParam, models.TaskStateRunning); err != nil {
							log.Printf("Failed to update task %s state to RUNNING: %v", taskUUIDParam, err)
						}
						// Register cron job
						if err := h.scheduler.RegisterTask(c.Request.Context(), task); err != nil {
							log.Printf("Failed to register task %s: %v", taskUUIDParam, err)
						} else {
							log.Printf("Registered cron job for task %s (status set to ACTIVE, group within window)", taskUUIDParam)
						}
					}
				}
			}
		} else if task.TaskGroupID == nil && h.scheduler != nil {
			// Task doesn't belong to a group - register directly if it has a cron expression
			if task.ScheduleConfig.CronExpression != "" {
				if err := h.scheduler.RegisterTask(c.Request.Context(), task); err != nil {
					log.Printf("Failed to register task %s in scheduler: %v", taskUUIDParam, err)
				} else {
					log.Printf("Registered cron job for task %s (status set to ACTIVE, no group)", taskUUIDParam)
				}
			}
		}
	}

	// Publish TaskUpdated event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskUpdated,
		Payload: events.TaskPayload{Task: task},
	})

	c.JSON(http.StatusOK, task)
}

// DeleteTask schedules a task for deletion (async). Returns 202 Accepted with PENDING_DELETE or ALREADY_DELETED.
// @Summary      Delete a task (async)
// @Description  Schedule a task for deletion. Deletion is performed asynchronously; use 202 response and status in body.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_uuid path string true "Task UUID"
// @Success      202  {object}  models.DeleteTaskResponse "Task deletion scheduled"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/tasks/{task_uuid} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	projectIDParam := c.Param("project_id")
	taskUUIDParam := c.Param("task_uuid")

	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}
	if taskUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task_uuid is required in path",
		})
		return
	}

	ctx := c.Request.Context()

	// Idempotent: if task already gone, treat as success
	task, err := h.repo.GetTaskByUUID(ctx, taskUUIDParam)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusAccepted, gin.H{
				"status":    "ALREADY_DELETED",
				"task_uuid": taskUUIDParam,
				"message":   "Task already deleted or not found",
			})
			return
		}
		log.Printf("DeleteTask GetTaskByUUID error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load task",
		})
		return
	}

	if h.deletePublisher == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Delete queue not configured",
		})
		return
	}

	// Mark as PENDING_DELETE before enqueueing
	if err := h.repo.UpdateTaskStatus(ctx, taskUUIDParam, models.TaskStatusPendingDelete); err != nil {
		log.Printf("DeleteTask UpdateTaskStatus error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to schedule deletion",
		})
		return
	}

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    task.UUID,
		ProjectID:   projectIDParam,
		RequestedAt: time.Now(),
	}

	if err := h.deletePublisher.PublishDeleteTask(ctx, msg); err != nil {
		log.Printf("DeleteTask PublishDeleteTask error: %v", err)
		// Rollback status to previous
		_ = h.repo.UpdateTaskStatus(ctx, taskUUIDParam, task.Status)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to enqueue delete job",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":    "PENDING_DELETE",
		"task_uuid": taskUUIDParam,
		"message":   "Task deletion has been scheduled",
	})
}

// UpdateTaskStatus updates a task's status
// @Summary      Update task status
// @Description  Update a task's status (ACTIVE or DISABLED) and update scheduler accordingly
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_uuid path string true "Task UUID"
// @Param        request body object true "Status update request" example({"status": "DISABLED"})
// @Success      200  {object}  models.Task
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/tasks/{task_uuid}/status [patch]
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	// Get project_id and task_uuid from path parameters
	projectIDParam := c.Param("project_id")
	taskUUIDParam := c.Param("task_uuid")

	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}

	if taskUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task_uuid is required in path",
		})
		return
	}

	// Convert project_id to ObjectID
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project_id format in path",
		})
		return
	}

	// Parse request body. Only ACTIVE and DISABLED are accepted; PENDING_DELETE/DELETE_FAILED are backend-only.
	var req struct {
		Status models.TaskStatus `json:"status" binding:"required,oneof=ACTIVE DISABLED"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Get existing task
	existingTask, err := h.repo.GetTaskByUUID(c.Request.Context(), taskUUIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Verify project_id matches
	if existingTask.ProjectID != projectID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task does not belong to this project",
		})
		return
	}

	// Don't update if status is already the same
	if existingTask.Status == req.Status {
		c.JSON(http.StatusOK, existingTask)
		return
	}

	// Determine task state based on status change
	// If status is being set to DISABLED, set state to NOT_RUNNING
	// If status is being set to ACTIVE and task belongs to an ACTIVE group within window, set state to RUNNING
	// Otherwise, preserve existing state (it's system-controlled)
	state := existingTask.State
	if req.Status == models.TaskStatusDisabled {
		state = models.TaskStateNotRunning
	} else if req.Status == models.TaskStatusActive && existingTask.Status != models.TaskStatusActive {
		// Status changed to ACTIVE - check if task belongs to an active group within window
		if existingTask.TaskGroupID != nil && h.scheduler != nil {
			taskGroup, err := h.repo.GetTaskGroupByID(c.Request.Context(), *existingTask.TaskGroupID)
			if err == nil && taskGroup != nil {
				if taskGroup.Status == models.TaskGroupStatusActive {
					// Check if group is within time window
					if h.scheduler.IsWithinGroupWindow(c.Request.Context(), taskGroup) {
						state = models.TaskStateRunning
					}
				}
			}
		}
	}

	// Update task status
	updatedTask := *existingTask
	updatedTask.Status = req.Status
	updatedTask.State = state
	updatedTask.UpdatedAt = time.Now()

	// Update in database
	err = h.repo.UpdateTask(c.Request.Context(), taskUUIDParam, &updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task status",
		})
		return
	}

	// Handle immediate actions based on status change
	if req.Status == models.TaskStatusDisabled && existingTask.Status != models.TaskStatusDisabled {
		// Status changed to DISABLED - update state and unregister cron job immediately
		if err := h.repo.UpdateTaskState(c.Request.Context(), taskUUIDParam, models.TaskStateNotRunning); err != nil {
			log.Printf("Failed to update task %s state to NOT_RUNNING: %v", taskUUIDParam, err)
		}

		// Unregister task from scheduler immediately
		if h.scheduler != nil {
			h.scheduler.UnregisterTask(taskUUIDParam)
			log.Printf("Unregistered cron job for task %s (status set to DISABLED)", taskUUIDParam)
		}
	} else if req.Status == models.TaskStatusActive && existingTask.Status != models.TaskStatusActive {
		// Status changed to ACTIVE - check if we need to register cron job and update state
		if updatedTask.TaskGroupID != nil && h.scheduler != nil {
			taskGroup, err := h.repo.GetTaskGroupByID(c.Request.Context(), *updatedTask.TaskGroupID)
			if err == nil && taskGroup != nil {
				if taskGroup.Status == models.TaskGroupStatusActive {
					// Check if group is within time window
					if h.scheduler.IsWithinGroupWindow(c.Request.Context(), taskGroup) {
						// Update state to RUNNING
						if err := h.repo.UpdateTaskState(c.Request.Context(), taskUUIDParam, models.TaskStateRunning); err != nil {
							log.Printf("Failed to update task %s state to RUNNING: %v", taskUUIDParam, err)
						}
						// Register cron job
						if err := h.scheduler.RegisterTask(c.Request.Context(), &updatedTask); err != nil {
							log.Printf("Failed to register task %s: %v", taskUUIDParam, err)
						} else {
							log.Printf("Registered cron job for task %s (status set to ACTIVE, group within window)", taskUUIDParam)
						}
					}
				}
			}
		} else {
			// Task doesn't belong to a group - register directly if it has a cron expression
			if updatedTask.ScheduleConfig.CronExpression != "" && h.scheduler != nil {
				if err := h.scheduler.RegisterTask(c.Request.Context(), &updatedTask); err != nil {
					log.Printf("Failed to register task %s in scheduler: %v", taskUUIDParam, err)
				}
			}
		}
	}

	// Publish TaskUpdated event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskUpdated,
		Payload: events.TaskPayload{Task: &updatedTask},
	})

	c.JSON(http.StatusOK, &updatedTask)
}

// TriggerTask manually triggers a task execution
// @Summary      Trigger task manually
// @Description  Manually trigger a task execution outside of cron schedule. Creates an execution record and sends it to the project's execution endpoint.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_uuid path string true "Task UUID"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/tasks/{task_uuid}/trigger [post]
func (h *TaskHandler) TriggerTask(c *gin.Context) {
	// Get project_id and task_uuid from path parameters
	projectIDParam := c.Param("project_id")
	taskUUIDParam := c.Param("task_uuid")

	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}

	if taskUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task_uuid is required in path",
		})
		return
	}

	// Get the task
	task, err := h.repo.GetTaskByUUID(c.Request.Context(), taskUUIDParam)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get task",
			})
		}
		return
	}

	// Use the shared ExecuteTask function from scheduler package
	executionUUID, err := scheduler.ExecuteTask(c.Request.Context(), task, h.repo, h.eventBus, "TRIGGER")
	if err != nil {
		if err.Error() == "no execution_endpoint set for project" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No execution_endpoint set for this project",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create execution record",
		})
		return
	}

	// Return immediately with the execution UUID
	now := time.Now()
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"execution_uuid": executionUUID,
			"task_uuid":      task.UUID,
			"status":         "PENDING",
			"trigger_type":   "MANUAL",
			"scheduled_at":   now.Format(time.RFC3339),
			"message":        "Execution created successfully",
		},
	})
}
