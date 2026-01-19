package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
)

type ExecutionHandler struct {
	repo repositories.Repository
}

func NewExecutionHandler(repo repositories.Repository) *ExecutionHandler {
	return &ExecutionHandler{
		repo: repo,
	}
}

// GetExecutionsByTaskUUID retrieves executions for a specific task
// @Summary      Get executions for a task
// @Description  Retrieve all executions for a specific task filtered by date
// @Tags         executions
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_uuid path string true "Task UUID"
// @Param        date query string true "Filter by date (YYYY-MM-DD format). Returns executions for that date only"
// @Success      200  {array}   models.Execution
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/tasks/{task_uuid}/executions [get]
func (h *ExecutionHandler) GetExecutionsByTaskUUID(c *gin.Context) {
	projectID := c.Param("project_id")
	taskUUID := c.Param("task_uuid")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}
	if taskUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task_uuid is required in path",
		})
		return
	}

	// Date parameter is required
	dateParam := c.Query("date")
	if dateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date parameter is required (YYYY-MM-DD format)",
		})
		return
	}

	// Parse date in YYYY-MM-DD format (will be parsed as UTC)
	parsedDate, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	// MongoDB stores times in UTC, so we need to create date range in UTC
	// Set startDate to beginning of day in UTC
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	startDate := &startOfDay
	// Set endDate to end of day in UTC
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, time.UTC)
	endDate := &endOfDay

	executions, err := h.repo.GetExecutionsByTaskUUID(c.Request.Context(), taskUUID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get executions for task %s: %v", taskUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get executions",
		})
		return
	}

	// Ensure we always return an empty array instead of null
	if executions == nil {
		executions = []*models.Execution{}
	}

	c.JSON(http.StatusOK, executions)
}

// AppendLogToExecution appends a log entry to an execution
// @Summary      Append log to execution
// @Description  Append a log entry to an execution by execution UUID
// @Tags         executions
// @Accept       json
// @Produce      json
// @Param        execution_uuid path string true "Execution UUID"
// @Param        log body object true "Log entry" example({"message": "Processing started", "level": "info"})
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /executions/{execution_uuid}/logs [post]
func (h *ExecutionHandler) AppendLogToExecution(c *gin.Context) {
	executionUUID := c.Param("execution_uuid")
	if executionUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "execution_uuid is required in path",
		})
		return
	}

	var logRequest struct {
		Message string `json:"message" binding:"required"`
		Level   string `json:"level" binding:"required"`
	}

	if err := c.ShouldBindJSON(&logRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": []string{err.Error()},
		})
		return
	}

	// Validate log level
	validLevels := map[string]bool{"info": true, "warn": true, "error": true}
	if !validLevels[logRequest.Level] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid log level. Must be one of: info, warn, error",
		})
		return
	}

	logEntry := models.LogEntry{
		Message:   logRequest.Message,
		Level:     logRequest.Level,
		Timestamp: time.Now(),
	}

	if err := h.repo.AppendLogToExecution(c.Request.Context(), executionUUID, logEntry); err != nil {
		log.Printf("Failed to append log to execution %s: %v", executionUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to append log",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Log appended successfully",
	})
}

// UpdateExecutionStatus updates the status of an execution
// @Summary      Update execution status
// @Description  Update the status of an execution (SUCCESS, FAILED, RUNNING)
// @Tags         executions
// @Accept       json
// @Produce      json
// @Param        execution_uuid path string true "Execution UUID"
// @Param        status body object true "Status update" example({"status": "SUCCESS"}) or example({"status": "FAILED", "error": "Error message"})
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /executions/{execution_uuid}/status [patch]
func (h *ExecutionHandler) UpdateExecutionStatus(c *gin.Context) {
	executionUUID := c.Param("execution_uuid")
	if executionUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "execution_uuid is required in path",
		})
		return
	}

	var statusRequest struct {
		Status string `json:"status" binding:"required"`
		Error  string `json:"error,omitempty"`
	}

	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": []string{err.Error()},
		})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"PENDING": true,
		"RUNNING": true,
		"SUCCESS": true,
		"FAILED":  true,
	}
	if !validStatuses[statusRequest.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status. Must be one of: PENDING, RUNNING, SUCCESS, FAILED",
		})
		return
	}

	var errorMsg *string
	if statusRequest.Error != "" {
		errorMsg = &statusRequest.Error
	}

	if err := h.repo.UpdateExecutionStatus(
		c.Request.Context(),
		executionUUID,
		models.ExecutionStatus(statusRequest.Status),
		errorMsg,
	); err != nil {
		log.Printf("Failed to update execution status for %s: %v", executionUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update execution status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Execution status updated successfully",
		"status":  statusRequest.Status,
	})
}
