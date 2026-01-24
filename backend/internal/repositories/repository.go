package repositories

import (
	"context"
	"time"

	"github.com/yourusername/cron-observer/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository defines project-related repository operations
type Repository interface {
	GetAllProjects(ctx context.Context) ([]*models.Project, error)
	GetProjectByID(ctx context.Context, projectID primitive.ObjectID) (*models.Project, error)
	GetUserProjects(ctx context.Context, email string) ([]*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) error
	UpdateProject(ctx context.Context, projectID primitive.ObjectID, project *models.Project) error

	// tasks
	CreateTask(ctx context.Context, projectID string, task *models.Task) error
	GetAllActiveTasks(ctx context.Context) ([]*models.Task, error)
	GetTasksByProjectID(ctx context.Context, projectID primitive.ObjectID) ([]*models.Task, error)
	GetTaskByUUID(ctx context.Context, taskUUID string) (*models.Task, error)
	UpdateTask(ctx context.Context, taskUUID string, task *models.Task) error
	DeleteTask(ctx context.Context, taskUUID string) error

	// task groups
	CreateTaskGroup(ctx context.Context, projectID string, taskGroup *models.TaskGroup) error
	GetTaskGroupsByProjectID(ctx context.Context, projectID primitive.ObjectID) ([]*models.TaskGroup, error)
	GetTaskGroupByUUID(ctx context.Context, taskGroupUUID string) (*models.TaskGroup, error)
	GetTaskGroupByID(ctx context.Context, taskGroupID primitive.ObjectID) (*models.TaskGroup, error)
	UpdateTaskGroup(ctx context.Context, taskGroupUUID string, taskGroup *models.TaskGroup) error
	UpdateTaskGroupStatus(ctx context.Context, taskGroupUUID string, status models.TaskGroupStatus) error
	UpdateTaskGroupState(ctx context.Context, taskGroupUUID string, state models.TaskGroupState) error
	DeleteTaskGroup(ctx context.Context, taskGroupUUID string) error
	GetTasksByGroupID(ctx context.Context, taskGroupID primitive.ObjectID) ([]*models.Task, error)
	GetActiveTaskGroupsWithWindows(ctx context.Context) ([]*models.TaskGroup, error)
	UpdateTaskStatus(ctx context.Context, taskUUID string, status models.TaskStatus) error
	UpdateTaskState(ctx context.Context, taskUUID string, state models.TaskState) error

	// executions
	CreateExecution(ctx context.Context, execution *models.Execution) error
	GetExecutionsByTaskUUID(ctx context.Context, taskUUID string, startDate, endDate *time.Time) ([]*models.Execution, error)
	GetExecutionsByTaskUUIDPaginated(ctx context.Context, taskUUID string, startDate, endDate *time.Time, page, pageSize int) ([]*models.Execution, int64, error)
	AppendLogToExecution(ctx context.Context, executionUUID string, logEntry models.LogEntry) error
	UpdateExecutionStatus(ctx context.Context, executionUUID string, status models.ExecutionStatus, errorMessage *string) error
	GetExecutionByUUID(ctx context.Context, executionUUID string) (*models.Execution, error)

	// failure statistics
	IncrementFailureStat(ctx context.Context, projectID primitive.ObjectID, date string) error
	GetFailureStatsByProject(ctx context.Context, projectID primitive.ObjectID, days int) ([]*models.FailedExecutionStats, int, error)
	
	// execution statistics
	GetExecutionStatsByProject(ctx context.Context, projectID primitive.ObjectID, days int) ([]*models.ExecutionStats, error)
}
