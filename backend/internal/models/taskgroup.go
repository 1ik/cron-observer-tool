package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskGroup represents a group of tasks that can be controlled together
// @Description TaskGroup represents a group of tasks that can be controlled together
type TaskGroup struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"507f1f77bcf86cd799439011"`
	UUID        string             `json:"uuid" bson:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProjectID   primitive.ObjectID `json:"project_id" bson:"project_id" example:"507f1f77bcf86cd799439011"`
	Name        string             `json:"name" bson:"name" example:"Morning Tasks"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" example:"Tasks that run in the morning"`
	Status      TaskGroupStatus    `json:"status" bson:"status" enums:"ACTIVE,DISABLED" example:"ACTIVE"`
	State       TaskGroupState     `json:"state" bson:"state" enums:"RUNNING,NOT_RUNNING" example:"NOT_RUNNING"`    // System-controlled: based on time window
	StartTime   string             `json:"start_time,omitempty" bson:"start_time,omitempty" example:"09:00"`        // Format: "HH:MM"
	EndTime     string             `json:"end_time,omitempty" bson:"end_time,omitempty" example:"17:00"`            // Format: "HH:MM"
	Timezone    string             `json:"timezone,omitempty" bson:"timezone,omitempty" example:"America/New_York"` // IANA timezone (e.g., "America/New_York")
	CreatedAt   time.Time          `json:"created_at" bson:"created_at" example:"2025-01-15T10:00:00Z"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at" example:"2025-01-15T10:00:00Z"`
}

// TaskGroupStatus defines the status of a task group
type TaskGroupStatus string

const (
	TaskGroupStatusActive   TaskGroupStatus = "ACTIVE"
	TaskGroupStatusDisabled TaskGroupStatus = "DISABLED"
)

// TaskGroupState defines the runtime state of a task group (system-controlled)
type TaskGroupState string

const (
	TaskGroupStateRunning    TaskGroupState = "RUNNING"
	TaskGroupStateNotRunning TaskGroupState = "NOT_RUNNING"
)

// CreateTaskGroupRequest represents the request DTO for creating a task group
type CreateTaskGroupRequest struct {
	ProjectID   string          `json:"project_id" binding:"required,objectid"`
	Name        string          `json:"name" binding:"required,min=1,max=255"`
	Description string          `json:"description,omitempty" binding:"omitempty,max=1000"`
	Status      TaskGroupStatus `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE DISABLED"`
	StartTime   string          `json:"start_time,omitempty" binding:"omitempty,time_format"` // Format: "HH:MM"
	EndTime     string          `json:"end_time,omitempty" binding:"omitempty,time_format"`   // Format: "HH:MM"
	Timezone    string          `json:"timezone,omitempty" binding:"omitempty,timezone"`
}

// UpdateTaskGroupRequest represents the request DTO for updating a task group
type UpdateTaskGroupRequest struct {
	Name        string          `json:"name" binding:"required,min=1,max=255"`
	Description string          `json:"description,omitempty" binding:"omitempty,max=1000"`
	Status      TaskGroupStatus `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE DISABLED"`
	StartTime   string          `json:"start_time,omitempty" binding:"omitempty,time_format"` // Format: "HH:MM"
	EndTime     string          `json:"end_time,omitempty" binding:"omitempty,time_format"`   // Format: "HH:MM"
	Timezone    string          `json:"timezone,omitempty" binding:"omitempty,timezone"`
}
