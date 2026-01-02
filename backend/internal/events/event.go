package events

import "github.com/yourusername/cron-observer/backend/internal/models"

// EventType defines the type of event
type EventType string

const (
	TaskCreated      EventType = "task.created"
	TaskUpdated      EventType = "task.updated"
	TaskDeleted      EventType = "task.deleted"
	TaskGroupCreated EventType = "taskgroup.created"
	TaskGroupUpdated EventType = "taskgroup.updated"
	TaskGroupDeleted EventType = "taskgroup.deleted"
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

// TaskDeletedPayload contains the task UUID for deleted events
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
