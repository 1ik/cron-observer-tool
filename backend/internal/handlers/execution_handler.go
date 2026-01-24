package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExecutionHandler struct {
	repo     repositories.Repository
	eventBus *events.EventBus
}

func NewExecutionHandler(repo repositories.Repository, eventBus *events.EventBus) *ExecutionHandler {
	return &ExecutionHandler{
		repo:     repo,
		eventBus: eventBus,
	}
}

// GetExecutionsByTaskUUID retrieves executions for a specific task
// @Summary      Get executions for a task
// @Description  Retrieve paginated executions for a specific task filtered by date
// @Tags         executions
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        task_uuid path string true "Task UUID"
// @Param        date query string true "Filter by date (YYYY-MM-DD format). Returns executions for that date only"
// @Param        page query int false "Page number (default: 1)"
// @Param        page_size query int false "Page size (default: 100)"
// @Success      200  {object}  models.PaginatedExecutionsResponse
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

	// Parse pagination parameters with defaults
	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 100
	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeParam); err == nil && parsedPageSize > 0 {
			// Limit max page size to prevent abuse
			if parsedPageSize > 100 {
				pageSize = 100
			} else {
				pageSize = parsedPageSize
			}
		}
	}

	// MongoDB stores times in UTC, so we need to create date range in UTC
	// Set startDate to beginning of day in UTC
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	startDate := &startOfDay
	// Set endDate to end of day in UTC
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, time.UTC)
	endDate := &endOfDay

	executions, totalCount, err := h.repo.GetExecutionsByTaskUUIDPaginated(c.Request.Context(), taskUUID, startDate, endDate, page, pageSize)
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

	// Calculate total pages
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	response := models.PaginatedExecutionsResponse{
		Data:       executions,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
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
			"error":   "Invalid request body",
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

	// Emit ExecutionFailed event if status is FAILED
	if models.ExecutionStatus(statusRequest.Status) == models.ExecutionStatusFailed {
		// Fetch execution and task for event payload
		execution, err := h.repo.GetExecutionByUUID(c.Request.Context(), executionUUID)
		if err == nil && execution != nil {
			task, err := h.repo.GetTaskByUUID(c.Request.Context(), execution.TaskUUID)
			if err == nil && task != nil {
				h.eventBus.Publish(events.Event{
					Type: events.ExecutionFailed,
					Payload: events.ExecutionFailedPayload{
						Execution: execution,
						Task:      task,
					},
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Execution status updated successfully",
		"status":  statusRequest.Status,
	})
}

// GetFailedExecutionsStats retrieves failure statistics for a project
// @Summary      Get failure statistics for a project
// @Description  Retrieve failed executions grouped by date for the last N days
// @Tags         executions
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        days query int false "Number of days to look back (default: 7)"
// @Success      200  {object}  models.FailedExecutionsStatsResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id}/executions/failed-stats [get]
func (h *ExecutionHandler) GetFailedExecutionsStats(c *gin.Context) {
	projectIDParam := c.Param("project_id")
	if projectIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required in path",
		})
		return
	}

	// Parse project ID
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project_id format",
		})
		return
	}

	// Parse optional days parameter (default: 7, max: 30)
	days := 7
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 {
			if parsedDays > 30 {
				days = 30
			} else {
				days = parsedDays
			}
		}
	}

	// Get failure stats
	stats, total, err := h.repo.GetFailureStatsByProject(c.Request.Context(), projectID, days)
	if err != nil {
		log.Printf("Failed to get failure stats for project %s: %v", projectIDParam, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get failure statistics",
		})
		return
	}

	// Convert pointers to values
	statsValues := make([]models.FailedExecutionStats, len(stats))
	for i, stat := range stats {
		statsValues[i] = *stat
	}

	response := models.FailedExecutionsStatsResponse{
		Stats: statsValues,
		Total: total,
	}

	c.JSON(http.StatusOK, response)
}
