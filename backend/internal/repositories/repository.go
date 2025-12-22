package repositories

import (
	"context"

	"github.com/yourusername/cron-observer/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository defines project-related repository operations
type Repository interface {
	GetAllProjects(ctx context.Context) ([]*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) error

	// tasks
	CreateTask(ctx context.Context, projectID string, task *models.Task) error
	GetTasksByProjectID(ctx context.Context, projectID primitive.ObjectID) ([]*models.Task, error)
}
