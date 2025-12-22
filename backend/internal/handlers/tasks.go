package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"github.com/yourusername/cron-observer/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	repo repositories.Repository
}

func NewTaskHandler(repo repositories.Repository) *TaskHandler {
	return &TaskHandler{
		repo: repo,
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

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasksByProject(c *gin.Context) {
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

	// Get tasks by project ID
	tasks, err := h.repo.GetTasksByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tasks",
		})
		return
	}

	// Return empty slice if nil (following GetAllProjects pattern)
	if tasks == nil {
		tasks = []*models.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}
