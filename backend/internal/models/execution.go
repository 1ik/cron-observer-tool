package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LogEntry represents a single log entry for an execution
type LogEntry struct {
	Message   string    `json:"message" bson:"message"`
	Level     string    `json:"level" bson:"level"` // info, warn, error
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// Execution represents a task execution record
// @Description Execution represents a task execution record
type Execution struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"507f1f77bcf86cd799439011"`
	UUID      string             `json:"uuid" bson:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	TaskID    primitive.ObjectID `json:"task_id" bson:"task_id" example:"507f1f77bcf86cd799439011"`
	TaskUUID  string             `json:"task_uuid" bson:"task_uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status    ExecutionStatus    `json:"status" bson:"status" enums:"PENDING,RUNNING,SUCCESS,FAILED" example:"PENDING"`
	StartedAt time.Time          `json:"started_at" bson:"started_at" example:"2025-01-15T10:00:00Z"`
	EndedAt   *time.Time         `json:"ended_at,omitempty" bson:"ended_at,omitempty" example:"2025-01-15T10:00:05Z"`
	Error     string             `json:"error,omitempty" bson:"error,omitempty" example:"Connection timeout"`
	Logs      []LogEntry         `json:"logs,omitempty" bson:"logs,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at" example:"2025-01-15T10:00:00Z"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" example:"2025-01-15T10:00:00Z"`
}

// ExecutionStatus defines the status of an execution
type ExecutionStatus string

const (
	ExecutionStatusPending ExecutionStatus = "PENDING"
	ExecutionStatusRunning ExecutionStatus = "RUNNING"
	ExecutionStatusSuccess ExecutionStatus = "SUCCESS"
	ExecutionStatusFailed  ExecutionStatus = "FAILED"
)

// PaginatedExecutionsResponse represents a paginated response for executions
type PaginatedExecutionsResponse struct {
	Data       []*Execution `json:"data"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalCount int64        `json:"total_count"`
	TotalPages int          `json:"total_pages"`
}
