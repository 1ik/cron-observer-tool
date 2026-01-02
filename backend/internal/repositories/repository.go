package repositories

import (
	"context"

	"github.com/yourusername/cron-observer/backend/internal/models"
)

// Repository defines project-related repository operations
type Repository interface {
	GetAllProjects(ctx context.Context) ([]*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) error

	// tasks
	CreateTask(ctx context.Context, projectID string, task *models.Task) error
	GetAllActiveTasks(ctx context.Context) ([]*models.Task, error)
	GetTaskByUUID(ctx context.Context, taskUUID string) (*models.Task, error)
	UpdateTask(ctx context.Context, taskUUID string, task *models.Task) error
	DeleteTask(ctx context.Context, taskUUID string) error
}
