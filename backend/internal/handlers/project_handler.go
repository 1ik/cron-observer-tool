package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"github.com/yourusername/cron-observer/backend/internal/utils"
)

type ProjectHandler struct {
	repo repositories.Repository
}

func NewProjectHandler(repo repositories.Repository) *ProjectHandler {
	return &ProjectHandler{
		repo: repo,
	}
}

func (h *ProjectHandler) GetAllProjects(c *gin.Context) {

	projects, err := h.repo.GetAllProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch projects",
		})
		return
	}

	if projects == nil {
		projects = []*models.Project{}
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Set timestamps
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	// Generate UUID and API key
	project.UUID = uuid.New().String()
	project.APIKey = utils.GenerateAPIKey()

	// create the project
	err := h.repo.CreateProject(c.Request.Context(), &project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create project",
		})
		return
	}

	c.JSON(http.StatusCreated, project)
}
