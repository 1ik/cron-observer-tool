package handlers

import (
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

	// Convert request DTO to TaskGroup model
	taskGroup := &models.TaskGroup{
		ProjectID:   projectID,
		UUID:        uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
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

	// Update task group fields
	taskGroup := &models.TaskGroup{
		ID:          existingTaskGroup.ID,
		UUID:        existingTaskGroup.UUID, // UUID cannot be changed
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
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

	// Publish TaskGroupUpdated event
	h.eventBus.Publish(events.Event{
		Type:    events.TaskGroupUpdated,
		Payload: events.TaskGroupPayload{TaskGroup: taskGroup},
	})

	c.JSON(http.StatusOK, taskGroup)
}

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
