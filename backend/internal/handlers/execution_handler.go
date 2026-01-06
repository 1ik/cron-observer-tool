package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, executions)
}

