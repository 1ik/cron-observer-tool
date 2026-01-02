package handlers

import (
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
	repo     repositories.Repository
	eventBus *events.EventBus
}

func NewTaskHandler(repo repositories.Repository, eventBus *events.EventBus) *TaskHandler {
	return &TaskHandler{
		repo:     repo,
		eventBus: eventBus,
	}
}

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

	// Convert request DTO to Task model
	task := &models.Task{
		ProjectID:    projectID,
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

	// Convert TriggerConfig from request DTO to model
	task.TriggerConfig = models.TriggerConfig{
		Type: req.TriggerConfig.Type,
		HTTP: &models.HTTPTriggerConfig{
			URL:     req.TriggerConfig.HTTP.URL,
			Method:  req.TriggerConfig.HTTP.Method,
			Headers: req.TriggerConfig.HTTP.Headers,
			Body:    req.TriggerConfig.HTTP.Body,
			Timeout: req.TriggerConfig.HTTP.Timeout,
		},
	}

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

	// Convert TriggerConfig from request DTO to model
	task.TriggerConfig = models.TriggerConfig{
		Type: req.TriggerConfig.Type,
		HTTP: &models.HTTPTriggerConfig{
			URL:     req.TriggerConfig.HTTP.URL,
			Method:  req.TriggerConfig.HTTP.Method,
			Headers: req.TriggerConfig.HTTP.Headers,
			Body:    req.TriggerConfig.HTTP.Body,
			Timeout: req.TriggerConfig.HTTP.Timeout,
		},
	}

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
