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
	GetAllActiveTasks(ctx context.Context) ([]*models.Task, error)
	GetTaskByUUID(ctx context.Context, taskUUID string) (*models.Task, error)
	UpdateTask(ctx context.Context, taskUUID string, task *models.Task) error
	DeleteTask(ctx context.Context, taskUUID string) error

	// task groups
	CreateTaskGroup(ctx context.Context, projectID string, taskGroup *models.TaskGroup) error
	GetTaskGroupByUUID(ctx context.Context, taskGroupUUID string) (*models.TaskGroup, error)
	GetTaskGroupByID(ctx context.Context, taskGroupID primitive.ObjectID) (*models.TaskGroup, error)
	UpdateTaskGroup(ctx context.Context, taskGroupUUID string, taskGroup *models.TaskGroup) error
	DeleteTaskGroup(ctx context.Context, taskGroupUUID string) error
	GetTasksByGroupID(ctx context.Context, taskGroupID primitive.ObjectID) ([]*models.Task, error)
	GetActiveTaskGroupsWithWindows(ctx context.Context) ([]*models.TaskGroup, error)
}
