package events

import "github.com/yourusername/cron-observer/backend/internal/models"

// EventType defines the type of event
type EventType string

const (
	TaskCreated EventType = "task.created"
	TaskUpdated EventType = "task.updated"
	TaskDeleted EventType = "task.deleted"
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
