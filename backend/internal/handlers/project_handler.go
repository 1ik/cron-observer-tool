package handlers

import (
	"log"
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

// GetAllProjects retrieves all projects
// @Summary      Get all projects
// @Description  Retrieve a list of all projects
// @Tags         projects
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Project
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects [get]
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

// CreateProject creates a new project
// @Summary      Create a new project
// @Description  Create a new project with auto-generated UUID and API key
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project body models.CreateProjectRequest true "Project creation request"
// @Success      201  {object}  models.Project
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": []string{err.Error()},
		})
		return
	}

	log.Printf("Parsed request: Name=%s, Description=%s, ExecutionEndpoint=%s", req.Name, req.Description, req.ExecutionEndpoint)

	// Validate that name is not empty (additional check)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project name is required",
		})
		return
	}

	// Create project model from request
	now := time.Now()
	project := &models.Project{
		Name:              req.Name,
		Description:       req.Description,
		ExecutionEndpoint: req.ExecutionEndpoint,
		UUID:              uuid.New().String(),
		APIKey:            utils.GenerateAPIKey(),
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Log the project being created for debugging
	log.Printf("Creating project: Name=%s, Description=%s, ExecutionEndpoint=%s", project.Name, project.Description, project.ExecutionEndpoint)

	// create the project
	err := h.repo.CreateProject(c.Request.Context(), project)
	if err != nil {
		log.Printf("Failed to create project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create project",
		})
		return
	}

	log.Printf("Project created successfully: ID=%s, UUID=%s, Name=%s", project.ID.Hex(), project.UUID, project.Name)
	c.JSON(http.StatusCreated, project)
}
