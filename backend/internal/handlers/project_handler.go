package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/cron-observer/backend/internal/middleware"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"github.com/yourusername/cron-observer/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectHandler struct {
	repo          repositories.Repository
	superAdminMap map[string]bool
}

func NewProjectHandler(repo repositories.Repository, superAdmins []string) *ProjectHandler {
	// Create a map for O(1) lookup
	superAdminMap := make(map[string]bool)
	for _, admin := range superAdmins {
		normalizedAdmin := strings.ToLower(strings.TrimSpace(admin))
		if normalizedAdmin != "" {
			superAdminMap[normalizedAdmin] = true
		}
	}

	return &ProjectHandler{
		repo:          repo,
		superAdminMap: superAdminMap,
	}
}

// isSuperAdmin checks if the given email is a super admin
func (h *ProjectHandler) isSuperAdmin(email string) bool {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	return h.superAdminMap[normalizedEmail]
}

// GetAllProjects retrieves all projects
// @Summary      Get all projects
// @Description  Retrieve a list of all projects. Super admins get all projects, regular users get only projects they are members of.
// @Tags         projects
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Project
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects [get]
func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	// Get authenticated user from context
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var projects []*models.Project
	var err error

	// Check if user is a super admin
	if h.isSuperAdmin(user.Email) {
		// Super admin - return all projects
		log.Printf("Super admin %s requesting all projects", user.Email)
		projects, err = h.repo.GetAllProjects(c.Request.Context())
	} else {
		// Regular user - return only projects they are members of
		log.Printf("User %s requesting their projects", user.Email)
		projects, err = h.repo.GetUserProjects(c.Request.Context(), user.Email)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch projects",
		})
		return
	}

	if projects == nil {
		projects = []*models.Project{}
	}

	// Ensure project_users is always initialized (not nil) for JSON serialization
	for _, project := range projects {
		if project.ProjectUsers == nil {
			project.ProjectUsers = []models.ProjectUser{}
		}
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
	// Get authenticated user from context
	user, exists := middleware.GetUserFromContext(c)
	if exists {
		log.Printf("User %s (%s) is creating a project", user.Name, user.Email)
	}

	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
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

	// Check project name is unique (case-insensitive)
	existing, getErr := h.repo.GetProjectByName(c.Request.Context(), strings.TrimSpace(req.Name))
	if getErr == nil && existing != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "A project with this name already exists",
		})
		return
	}

	// Create project model from request
	now := time.Now()
	name := strings.TrimSpace(req.Name)
	project := &models.Project{
		ID:                primitive.NewObjectID(),
		Name:              name,
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

// UpdateProject updates an existing project
// @Summary      Update a project
// @Description  Update an existing project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project_id path string true "Project ID"
// @Param        project body models.UpdateProjectRequest true "Project update request"
// @Success      200  {object}  models.Project
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /projects/{project_id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": []string{err.Error()},
		})
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

	// Convert project_id to ObjectID
	projectID, err := primitive.ObjectIDFromHex(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project_id format in path",
		})
		return
	}

	// Get existing project to preserve UUID, APIKey, and timestamps
	existingProject, err := h.repo.GetProjectByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Project not found",
		})
		return
	}

	// Update only provided fields
	now := time.Now()
	updatedProject := &models.Project{
		ID:                existingProject.ID,
		UUID:              existingProject.UUID,   // UUID cannot be changed
		APIKey:            existingProject.APIKey, // API key cannot be changed
		Name:              existingProject.Name,
		Description:       existingProject.Description,
		ExecutionEndpoint: existingProject.ExecutionEndpoint,
		AlertEmails:       existingProject.AlertEmails,
		ProjectUsers:      existingProject.ProjectUsers, // Preserve existing users
		CreatedAt:         existingProject.CreatedAt,    // Preserve original creation time
		UpdatedAt:         now,
	}

	// Update fields if provided in request
	if req.Name != "" {
		newName := strings.TrimSpace(req.Name)
		// Check name is unique (case-insensitive), excluding current project
		existingByName, getErr := h.repo.GetProjectByName(c.Request.Context(), newName)
		if getErr == nil && existingByName != nil && existingByName.ID != projectID {
			c.JSON(http.StatusConflict, gin.H{
				"error": "A project with this name already exists",
			})
			return
		}
		updatedProject.Name = newName
	}
	if req.Description != "" {
		updatedProject.Description = req.Description
	} else if req.Description == "" && c.GetHeader("Content-Type") == "application/json" {
		// Allow clearing description by sending empty string
		updatedProject.Description = ""
	}
	if req.ExecutionEndpoint != "" {
		updatedProject.ExecutionEndpoint = req.ExecutionEndpoint
	} else if req.ExecutionEndpoint == "" && c.GetHeader("Content-Type") == "application/json" {
		// Allow clearing execution endpoint by sending empty string
		updatedProject.ExecutionEndpoint = ""
	}
	if req.AlertEmails != "" {
		updatedProject.AlertEmails = req.AlertEmails
	} else if req.AlertEmails == "" && c.GetHeader("Content-Type") == "application/json" {
		// Allow clearing alert emails by sending empty string
		updatedProject.AlertEmails = ""
	}
	if req.ProjectUsers != nil {
		updatedProject.ProjectUsers = req.ProjectUsers
		log.Printf("Updating project_users: %d users", len(req.ProjectUsers))
		for i, user := range req.ProjectUsers {
			log.Printf("  User %d: email=%s, role=%s", i+1, user.Email, user.Role)
		}
	} else {
		log.Printf("ProjectUsers not provided in request, preserving existing: %d users", len(updatedProject.ProjectUsers))
	}

	// Update the project
	log.Printf("Updating project with ProjectUsers: %v", updatedProject.ProjectUsers)
	err = h.repo.UpdateProject(c.Request.Context(), projectID, updatedProject)
	if err != nil {
		log.Printf("Failed to update project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update project",
		})
		return
	}

	log.Printf("Project updated successfully: ID=%s, UUID=%s, Name=%s", updatedProject.ID.Hex(), updatedProject.UUID, updatedProject.Name)
	c.JSON(http.StatusOK, updatedProject)
}
