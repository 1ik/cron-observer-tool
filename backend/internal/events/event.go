package events

import "github.com/yourusername/cron-observer/backend/internal/models"

// EventType defines the type of event
type EventType string

const (
	TaskCreated       EventType = "task.created"
	TaskUpdated       EventType = "task.updated"
	TaskDeleted       EventType = "task.deleted" // Published after a task is hard-deleted (e.g. by delete worker); scheduler unregisters it.
	TaskGroupCreated  EventType = "taskgroup.created"
	TaskGroupUpdated  EventType = "taskgroup.updated"
	TaskGroupDeleted  EventType = "taskgroup.deleted"
	ExecutionFailed   EventType = "execution.failed"
	ExecutionTimedOut EventType = "execution.timed_out"
)

// Event represents an event in the system
type Event struct {
	Type    EventType
	Payload interface{}
}

// TaskPayload contains the task data for created/updated events
type TaskPayload struct {
	Task *models.Task
}

// TaskDeletedPayload contains the task UUID for TaskDeleted events. Used when publishing after a hard delete.
type TaskDeletedPayload struct {
	TaskUUID string
}

// TaskGroupPayload contains the task group data for created/updated events
type TaskGroupPayload struct {
	TaskGroup *models.TaskGroup
}

// TaskGroupDeletedPayload contains the task group UUID for deleted events
type TaskGroupDeletedPayload struct {
	TaskGroupUUID string
}

// ExecutionFailedPayload contains execution and task data for failed execution events
type ExecutionFailedPayload struct {
	Execution *models.Execution
	Task      *models.Task
}

// ExecutionTimedOutPayload contains execution UUID and timeout information
type ExecutionTimedOutPayload struct {
	ExecutionUUID  string
	TaskUUID       string
	TimeoutSeconds int
}
