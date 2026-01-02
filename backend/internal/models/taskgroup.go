package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskGroup represents a group of tasks that can be controlled together
type TaskGroup struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid"`
	ProjectID   primitive.ObjectID `json:"project_id" bson:"project_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Status      TaskGroupStatus    `json:"status" bson:"status"`
	StartTime   string             `json:"start_time,omitempty" bson:"start_time,omitempty"` // Format: "HH:MM"
	EndTime     string             `json:"end_time,omitempty" bson:"end_time,omitempty"`     // Format: "HH:MM"
	Timezone    string             `json:"timezone,omitempty" bson:"timezone,omitempty"`     // IANA timezone (e.g., "America/New_York")
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// TaskGroupStatus defines the status of a task group
type TaskGroupStatus string

const (
	TaskGroupStatusActive   TaskGroupStatus = "ACTIVE"
	TaskGroupStatusPaused   TaskGroupStatus = "PAUSED"
	TaskGroupStatusDisabled TaskGroupStatus = "DISABLED"
)

// CreateTaskGroupRequest represents the request DTO for creating a task group
type CreateTaskGroupRequest struct {
	ProjectID   string          `json:"project_id" binding:"required,objectid"`
	Name        string          `json:"name" binding:"required,min=1,max=255"`
	Description string          `json:"description,omitempty" binding:"omitempty,max=1000"`
	Status      TaskGroupStatus `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE PAUSED DISABLED"`
	StartTime   string          `json:"start_time,omitempty" binding:"omitempty,time_format"` // Format: "HH:MM"
	EndTime     string          `json:"end_time,omitempty" binding:"omitempty,time_format"`   // Format: "HH:MM"
	Timezone    string          `json:"timezone,omitempty" binding:"omitempty,timezone"`
}

// UpdateTaskGroupRequest represents the request DTO for updating a task group
type UpdateTaskGroupRequest struct {
	Name        string          `json:"name" binding:"required,min=1,max=255"`
	Description string          `json:"description,omitempty" binding:"omitempty,max=1000"`
	Status      TaskGroupStatus `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE PAUSED DISABLED"`
	StartTime   string          `json:"start_time,omitempty" binding:"omitempty,time_format"` // Format: "HH:MM"
	EndTime     string          `json:"end_time,omitempty" binding:"omitempty,time_format"`   // Format: "HH:MM"
	Timezone    string          `json:"timezone,omitempty" binding:"omitempty,timezone"`
}
