package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string   `json:"error" example:"Invalid request"`
	Details []string `json:"details,omitempty" example:"project_id is required"`
}

// DeleteTaskResponse represents the response for async task deletion
type DeleteTaskResponse struct {
	Status   string `json:"status" example:"PENDING_DELETE" enums:"PENDING_DELETE,ALREADY_DELETED"`
	TaskUUID string `json:"task_uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message  string `json:"message" example:"Task deletion has been scheduled"`
}
