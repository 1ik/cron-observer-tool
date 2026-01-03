package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string   `json:"error" example:"Invalid request"`
	Details []string `json:"details,omitempty" example:"project_id is required"`
}
