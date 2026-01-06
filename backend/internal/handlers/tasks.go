package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"github.com/yourusername/cron-observer/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	repo      repositories.Repository
	eventBus  *events.EventBus
	scheduler interface {
		RegisterTask(ctx context.Context, task *models.Task) error
		UnregisterTask(taskUUID string)
	}
}

func NewTaskHandler(repo repositories.Repository, eventBus *events.EventBus, scheduler interface {
	RegisterTask(ctx context.Context, task *models.Task) error
	UnregisterTask(taskUUID string)
}) *TaskHandler {
	return &TaskHandler{
		repo:      repo,
		eventBus:  eventBus,
		scheduler: scheduler, // Can be nil if scheduler is not needed
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

	// Set default status if not provided
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

	// Convert request DTO to Task model
	task := &models.Task{
		ProjectID:    projectID,
		TaskGroupID:  taskGroupID,
		UUID:         uuid.New().String(),
		Name:         req.Name,
		Description:  req.Description,
		ScheduleType: req.ScheduleType,
		Status:       status,
		ScheduleConfig: models.ScheduleConfig{
			CronExpression: req.ScheduleConfig.CronExpression,
			Timezone:       req.ScheduleConfig.Timezone,
			DaysOfWeek:     req.ScheduleConfig.DaysOfWeek,
			Exclusions:     req.ScheduleConfig.Exclusions,
		},
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

	// Get existing task to preserve UUID and timestamps
	existingTask, err := h.repo.GetTaskByUUID(c.Request.Context(), taskUUIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = models.TaskStatusActive
	}

	// Update task fields
	task := &models.Task{
		ID:           existingTask.ID,
		UUID:         existingTask.UUID, // UUID cannot be changed
		ProjectID:    projectID,
		Name:         req.Name,
		Description:  req.Description,
		ScheduleType: req.ScheduleType,
		Status:       status,
		ScheduleConfig: models.ScheduleConfig{
			CronExpression: req.ScheduleConfig.CronExpression,
			Timezone:       req.ScheduleConfig.Timezone,
			DaysOfWeek:     req.ScheduleConfig.DaysOfWeek,
			Exclusions:     req.ScheduleConfig.Exclusions,
		},
		Metadata:  req.Metadata,
		CreatedAt: existingTask.CreatedAt, // Preserve original creation time
		UpdatedAt: time.Now(),
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

	// Publish TaskUpdated event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskUpdated,
		Payload: events.TaskPayload{Task: task},
	})

	c.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task
// @Summary      Delete a task
// @Description  Delete an existing scheduled task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_uuid path string true "Task UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/tasks/{task_uuid} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// Get task_uuid from path parameter
	taskUUIDParam := c.Param("task_uuid")

	if taskUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task_uuid is required in path",
		})
		return
	}

	// Delete the task
	err := h.repo.DeleteTask(c.Request.Context(), taskUUIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task",
		})
		return
	}

	// Publish TaskDeleted event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskDeleted,
		Payload: events.TaskDeletedPayload{TaskUUID: taskUUIDParam},
	})

	c.Status(http.StatusNoContent)
}

// UpdateTaskStatus updates a task's status (pause/play)
// @Summary      Update task status
// @Description  Update a task's status (ACTIVE or PAUSED) and update scheduler accordingly
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_uuid path string true "Task UUID"
// @Param        request body object true "Status update request" example({"status": "PAUSED"})
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

	// Parse request body
	var req struct {
		Status models.TaskStatus `json:"status" binding:"required,oneof=ACTIVE PAUSED"`
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

	// Update task status
	updatedTask := *existingTask
	updatedTask.Status = req.Status
	updatedTask.UpdatedAt = time.Now()

	// Update in database
	err = h.repo.UpdateTask(c.Request.Context(), taskUUIDParam, &updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task status",
		})
		return
	}

	// Update scheduler if available
	if h.scheduler != nil {
		if req.Status == models.TaskStatusActive {
			// Register task in scheduler
			if err := h.scheduler.RegisterTask(c.Request.Context(), &updatedTask); err != nil {
				log.Printf("Failed to register task %s in scheduler: %v", taskUUIDParam, err)
				// Don't fail the request, just log the error
			}
		} else if req.Status == models.TaskStatusPaused {
			// Unregister task from scheduler
			h.scheduler.UnregisterTask(taskUUIDParam)
		}
	}

	// Publish TaskUpdated event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskUpdated,
		Payload: events.TaskPayload{Task: &updatedTask},
	})

	c.JSON(http.StatusOK, &updatedTask)
}
