package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cron-observer/backend/internal/middleware"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectAuthGuard checks if the current user has admin access to a project
// Returns true if:
//   - User is a super admin, OR
//   - User is in project's project_users with role 'admin'
//
// Returns false otherwise
func ProjectAuthGuard(c *gin.Context, repo repositories.Repository, projectID primitive.ObjectID, superAdminMap map[string]bool) bool {
	// Get authenticated user from context
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		log.Printf("[AUTH GUARD] User not authenticated")
		return false
	}

	userEmail := strings.ToLower(strings.TrimSpace(user.Email))
	if userEmail == "" {
		log.Printf("[AUTH GUARD] User email is empty")
		return false
	}

	// Check if user is a super admin
	if superAdminMap[userEmail] {
		log.Printf("[AUTH GUARD] User %s is a super admin, access granted", userEmail)
		return true
	}

	// Get project to check project_users
	project, err := repo.GetProjectByID(c.Request.Context(), projectID)
	if err != nil {
		log.Printf("[AUTH GUARD] Failed to get project %s: %v", projectID.Hex(), err)
		return false
	}

	// Check if user is in project_users with role 'admin'
	for _, projectUser := range project.ProjectUsers {
		projectUserEmail := strings.ToLower(strings.TrimSpace(projectUser.Email))
		if projectUserEmail == userEmail && projectUser.Role == models.ProjectUserRoleAdmin {
			log.Printf("[AUTH GUARD] User %s is admin in project %s, access granted", userEmail, projectID.Hex())
			return true
		}
	}

	log.Printf("[AUTH GUARD] User %s does not have admin access to project %s", userEmail, projectID.Hex())
	return false
}

// RequireProjectAdmin is a middleware-like function that checks authorization and returns error if not authorized
func RequireProjectAdmin(c *gin.Context, repo repositories.Repository, projectID primitive.ObjectID, superAdminMap map[string]bool) bool {
	if !ProjectAuthGuard(c, repo, projectID, superAdminMap) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to perform this action. Admin role or super admin access required.",
		})
		c.Abort()
		return false
	}
	return true
}
