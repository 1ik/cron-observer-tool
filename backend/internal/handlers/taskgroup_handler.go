package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"github.com/yourusername/cron-observer/backend/internal/scheduler"
	"github.com/yourusername/cron-observer/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskGroupHandler struct {
	repo      repositories.Repository
	eventBus  *events.EventBus
	scheduler *scheduler.Scheduler
}

func NewTaskGroupHandler(repo repositories.Repository, eventBus *events.EventBus, sched *scheduler.Scheduler) *TaskGroupHandler {
	return &TaskGroupHandler{
		repo:      repo,
		eventBus:  eventBus,
		scheduler: sched,
	}
}

// GetTaskGroupsByProject retrieves all task groups for a project
// @Summary      Get task groups by project
// @Description  Retrieve all task groups belonging to a project
// @Tags         task-groups
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Success      200  {array}   models.TaskGroup
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/task-groups [get]
func (h *TaskGroupHandler) GetTaskGroupsByProject(c *gin.Context) {
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

	// Get all task groups for this project
	taskGroups, err := h.repo.GetTaskGroupsByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get task groups for project",
		})
		return
	}

	if taskGroups == nil {
		taskGroups = []*models.TaskGroup{}
	}

	c.JSON(http.StatusOK, taskGroups)
}

// CreateTaskGroup creates a new task group
// @Summary      Create a new task group
// @Description  Create a new task group in a project
// @Tags         task-groups
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_group body models.CreateTaskGroupRequest true "Task group creation request"
// @Success      201  {object}  models.TaskGroup
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/task-groups [post]
func (h *TaskGroupHandler) CreateTaskGroup(c *gin.Context) {
	var req models.CreateTaskGroupRequest

	// Bind JSON and validate
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
		status = models.TaskGroupStatusActive
	}

	// Set default timezone if not provided
	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	// Calculate initial state based on time window
	state := models.TaskGroupStateNotRunning
	if req.StartTime != "" && req.EndTime != "" {
		// Check if current time is within the window
		tempTaskGroup := &models.TaskGroup{
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			Timezone:  timezone,
		}
		if h.scheduler.IsWithinGroupWindow(c.Request.Context(), tempTaskGroup) {
			state = models.TaskGroupStateRunning
		}
	}

	// Convert request DTO to TaskGroup model
	taskGroup := &models.TaskGroup{
		ProjectID:   projectID,
		UUID:        uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
		State:       state, // Set calculated state
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Timezone:    timezone,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create the task group
	err = h.repo.CreateTaskGroup(c.Request.Context(), projectIDParam, taskGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task group",
		})
		return
	}

	// Publish TaskGroupCreated event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskGroupCreated,
		Payload: events.TaskGroupPayload{TaskGroup: taskGroup},
	})

	c.JSON(http.StatusCreated, taskGroup)
}

// GetTaskGroup retrieves a task group by UUID
// @Summary      Get a task group
// @Description  Retrieve a task group by its UUID
// @Tags         task-groups
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        group_uuid path string true "Task Group UUID"
// @Success      200  {object}  models.TaskGroup
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/task-groups/{group_uuid} [get]
func (h *TaskGroupHandler) GetTaskGroup(c *gin.Context) {
	taskGroupUUID := c.Param("group_uuid")
	if taskGroupUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group_uuid is required in path",
		})
		return
	}

	taskGroup, err := h.repo.GetTaskGroupByUUID(c.Request.Context(), taskGroupUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task group not found",
		})
		return
	}

	c.JSON(http.StatusOK, taskGroup)
}

// UpdateTaskGroup updates an existing task group
// @Summary      Update a task group
// @Description  Update an existing task group
// @Tags         task-groups
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        group_uuid path string true "Task Group UUID"
// @Param        task_group body models.UpdateTaskGroupRequest true "Task group update request"
// @Success      200  {object}  models.TaskGroup
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/task-groups/{group_uuid} [put]
func (h *TaskGroupHandler) UpdateTaskGroup(c *gin.Context) {
	var req models.UpdateTaskGroupRequest

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Get project_id and group_uuid from path parameters
	projectIDParam := c.Param("project_id")
	taskGroupUUIDParam := c.Param("group_uuid")

	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}

	if taskGroupUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group_uuid is required in path",
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

	// Get existing task group to preserve ID, UUID, ProjectID and timestamps
	existingTaskGroup, err := h.repo.GetTaskGroupByUUID(c.Request.Context(), taskGroupUUIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task group not found",
		})
		return
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = existingTaskGroup.Status
	}

	// Set default timezone if not provided
	timezone := req.Timezone
	if timezone == "" {
		timezone = existingTaskGroup.Timezone
		if timezone == "" {
			timezone = "UTC"
		}
	}

	// Calculate state based on time window
	state := existingTaskGroup.State // Preserve existing state by default

	// If status is being changed to ACTIVE, recalculate state based on current time window
	if req.Status == models.TaskGroupStatusActive && existingTaskGroup.Status != models.TaskGroupStatusActive {
		// Status changed to ACTIVE, recalculate state based on current window
		if req.StartTime != "" && req.EndTime != "" {
			tempTaskGroup := &models.TaskGroup{
				StartTime: req.StartTime,
				EndTime:   req.EndTime,
				Timezone:  timezone,
			}
			if h.scheduler.IsWithinGroupWindow(c.Request.Context(), tempTaskGroup) {
				state = models.TaskGroupStateRunning
			} else {
				state = models.TaskGroupStateNotRunning
			}
		}
	} else if req.StartTime != "" && req.EndTime != "" {
		// Check if start_time or end_time changed
		if req.StartTime != existingTaskGroup.StartTime || req.EndTime != existingTaskGroup.EndTime || req.Timezone != existingTaskGroup.Timezone {
			// Calculate new state based on updated window
			tempTaskGroup := &models.TaskGroup{
				StartTime: req.StartTime,
				EndTime:   req.EndTime,
				Timezone:  timezone,
			}
			if h.scheduler.IsWithinGroupWindow(c.Request.Context(), tempTaskGroup) {
				state = models.TaskGroupStateRunning
			} else {
				state = models.TaskGroupStateNotRunning
			}
		}
	} else if req.StartTime == "" || req.EndTime == "" {
		// Window removed, set to NOT_RUNNING
		state = models.TaskGroupStateNotRunning
	}

	// Update task group fields
	taskGroup := &models.TaskGroup{
		ID:          existingTaskGroup.ID,
		UUID:        existingTaskGroup.UUID, // UUID cannot be changed
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
		State:       state, // Set calculated state
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Timezone:    timezone,
		CreatedAt:   existingTaskGroup.CreatedAt, // Preserve original creation time
		UpdatedAt:   time.Now(),
	}

	// Update the task group
	err = h.repo.UpdateTaskGroup(c.Request.Context(), taskGroupUUIDParam, taskGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task group",
		})
		return
	}

	// Update state separately to ensure it's persisted
	if err := h.repo.UpdateTaskGroupState(c.Request.Context(), taskGroupUUIDParam, state); err != nil {
		log.Printf("Failed to update task group state: %v", err)
	}

	// Update ALL tasks' states in this group BEFORE returning the response
	// This ensures the frontend gets the updated task states when it refetches
	tasks, err := h.repo.GetTasksByGroupID(c.Request.Context(), taskGroup.ID)
	if err != nil {
		log.Printf("Failed to get tasks for group %s: %v", taskGroup.UUID, err)
	} else {
		// Determine task state based on group state
		var taskState models.TaskState
		if state == models.TaskGroupStateRunning {
			taskState = models.TaskStateRunning
		} else {
			taskState = models.TaskStateNotRunning
		}

		// Update all tasks' states
		for _, task := range tasks {
			if err := h.repo.UpdateTaskState(c.Request.Context(), task.UUID, taskState); err != nil {
				log.Printf("Failed to update task %s state to %s: %v", task.UUID, taskState, err)
			}
		}
		log.Printf("[GROUP] Updated %d tasks' state to %s for group %s", len(tasks), taskState, taskGroup.UUID)
	}

	// Publish TaskGroupUpdated event (for scheduler to register/unregister cron jobs)
	h.eventBus.Publish(events.Event{
		Type:    events.TaskGroupUpdated,
		Payload: events.TaskGroupPayload{TaskGroup: taskGroup},
	})

	c.JSON(http.StatusOK, taskGroup)
}

// DeleteTaskGroup deletes a task group
// @Summary      Delete a task group
// @Description  Delete an existing task group
// @Tags         task-groups
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        group_uuid path string true "Task Group UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/task-groups/{group_uuid} [delete]
func (h *TaskGroupHandler) DeleteTaskGroup(c *gin.Context) {
	taskGroupUUIDParam := c.Param("group_uuid")

	if taskGroupUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group_uuid is required in path",
		})
		return
	}

	// Delete the task group
	err := h.repo.DeleteTaskGroup(c.Request.Context(), taskGroupUUIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task group",
		})
		return
	}

	// Publish TaskGroupDeleted event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskGroupDeleted,
		Payload: events.TaskGroupDeletedPayload{TaskGroupUUID: taskGroupUUIDParam},
	})

	c.Status(http.StatusNoContent)
}

// StartGroup starts all tasks in a task group
// @Summary      Start a task group
// @Description  Manually start all tasks in a task group
// @Tags         task-groups
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        group_uuid path string true "Task Group UUID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/task-groups/{group_uuid}/start [post]
func (h *TaskGroupHandler) StartGroup(c *gin.Context) {
	taskGroupUUIDParam := c.Param("group_uuid")

	if taskGroupUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group_uuid is required in path",
		})
		return
	}

	err := h.scheduler.StartGroup(c.Request.Context(), taskGroupUUIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to start group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Group started successfully",
	})
}

// StopGroup stops all tasks in a task group
// @Summary      Stop a task group
// @Description  Manually stop all tasks in a task group
// @Tags         task-groups
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        group_uuid path string true "Task Group UUID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/task-groups/{group_uuid}/stop [post]
func (h *TaskGroupHandler) StopGroup(c *gin.Context) {
	taskGroupUUIDParam := c.Param("group_uuid")

	if taskGroupUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group_uuid is required in path",
		})
		return
	}

	err := h.scheduler.StopGroup(c.Request.Context(), taskGroupUUIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to stop group",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Group stopped successfully",
	})
}

// GetTasksByGroup retrieves all tasks in a task group
// @Summary      Get tasks in a group
// @Description  Retrieve all tasks belonging to a task group
// @Tags         task-groups
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        group_uuid path string true "Task Group UUID"
// @Success      200  {array}   models.Task
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/task-groups/{group_uuid}/tasks [get]
func (h *TaskGroupHandler) GetTasksByGroup(c *gin.Context) {
	taskGroupUUIDParam := c.Param("group_uuid")

	if taskGroupUUIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group_uuid is required in path",
		})
		return
	}

	// Get task group to get its ID
	taskGroup, err := h.repo.GetTaskGroupByUUID(c.Request.Context(), taskGroupUUIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task group not found",
		})
		return
	}

	// Get all tasks in this group
	tasks, err := h.repo.GetTasksByGroupID(c.Request.Context(), taskGroup.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get tasks for group",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
