package deletequeue

import "time"

// DeleteTaskMessage is the message contract for enqueueing a task deletion job.
// It is serialized to JSON when publishing to the message broker.
type DeleteTaskMessage struct {
	TaskUUID    string    `json:"task_uuid"`
	ProjectID   string    `json:"project_id"`
	RequestedAt time.Time `json:"requested_at"`
	RequestID   string    `json:"request_id,omitempty"`
}
