package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
)

// ProjectContextKey is the key for storing project info in gin context
const ProjectContextKey = "project"

// APIKeyMiddleware validates API key authentication for SDK endpoints
// It validates that the API key matches the project that owns the execution
func APIKeyMiddleware(repo repositories.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract API key from Authorization header (raw format, no prefix)
		apiKey := c.GetHeader("Authorization")
		if apiKey == "" {
			log.Printf("[API_KEY] Missing Authorization header for %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Get execution UUID from path parameter
		executionUUID := c.Param("execution_uuid")
		if executionUUID == "" {
			log.Printf("[API_KEY] Missing execution_uuid parameter for %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "execution_uuid is required",
			})
			c.Abort()
			return
		}

		// Get execution by UUID
		execution, err := repo.GetExecutionByUUID(c.Request.Context(), executionUUID)
		if err != nil {
			log.Printf("[API_KEY] Execution not found: %s, error: %v", executionUUID, err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Execution not found",
			})
			c.Abort()
			return
		}

		// Get task by execution's TaskUUID
		task, err := repo.GetTaskByUUID(c.Request.Context(), execution.TaskUUID)
		if err != nil {
			log.Printf("[API_KEY] Task not found for execution %s: %v", executionUUID, err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			c.Abort()
			return
		}

		// Get project by task's ProjectID
		project, err := repo.GetProjectByID(c.Request.Context(), task.ProjectID)
		if err != nil {
			log.Printf("[API_KEY] Project not found for task %s: %v", execution.TaskUUID, err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			c.Abort()
			return
		}

		// Match API key from Authorization header with project's API key
		if project.APIKey != apiKey {
			log.Printf("[API_KEY] API key mismatch for execution %s (project: %s)", executionUUID, project.ID.Hex())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		// Store project info in context for handlers to access
		c.Set(ProjectContextKey, project)

		// Continue to next handler
		c.Next()
	}
}

// GetProjectFromContext extracts project info from gin context
func GetProjectFromContext(c *gin.Context) (*models.Project, bool) {
	project, exists := c.Get(ProjectContextKey)
	if !exists {
		return nil, false
	}

	projectInfo, ok := project.(*models.Project)
	if !ok {
		return nil, false
	}

	return projectInfo, true
}
