package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
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
