package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a scheduled task entity
// @Description Task represents a scheduled task entity
type Task struct {
	ID             primitive.ObjectID     `json:"id" bson:"_id,omitempty" example:"507f1f77bcf86cd799439011"`
	UUID           string                 `json:"uuid" bson:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProjectID      primitive.ObjectID     `json:"project_id" bson:"project_id" example:"507f1f77bcf86cd799439011"`
	TaskGroupID    *primitive.ObjectID    `json:"task_group_id,omitempty" bson:"task_group_id,omitempty" example:"507f1f77bcf86cd799439011"` // Optional reference to task group
	Name           string                 `json:"name" bson:"name" example:"Daily Backup"`
	Description    string                 `json:"description,omitempty" bson:"description,omitempty" example:"Backup database daily"`
	ScheduleType   ScheduleType           `json:"schedule_type" bson:"schedule_type" enums:"RECURRING,ONEOFF" example:"RECURRING"`
	Status         TaskStatus             `json:"status" bson:"status" enums:"ACTIVE,PAUSED,DISABLED" example:"ACTIVE"`
	ScheduleConfig ScheduleConfig         `json:"schedule_config" bson:"schedule_config"`
	TriggerConfig  TriggerConfig          `json:"trigger_config,omitempty" bson:"trigger_config,omitempty"` // Deprecated: Tasks now use project's execution_endpoint
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2025-01-15T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" example:"2025-01-15T10:00:00Z"`
}

// ScheduleType defines the type of schedule
type ScheduleType string

const (
	ScheduleTypeRecurring ScheduleType = "RECURRING"
	ScheduleTypeOneOff    ScheduleType = "ONEOFF"
)

// TaskStatus defines the status of a task
type TaskStatus string

const (
	TaskStatusActive   TaskStatus = "ACTIVE"
	TaskStatusPaused   TaskStatus = "PAUSED"
	TaskStatusDisabled TaskStatus = "DISABLED"
)

// ScheduleConfig holds the schedule configuration for a task
// Behavior:
//   - If CronExpression is provided: TimeRange and DaysOfWeek are ignored, schedule follows cron expression only
//   - If CronExpression is not provided: TimeRange and DaysOfWeek are used to determine execution schedule
type ScheduleConfig struct {
	CronExpression string     `json:"cron_expression,omitempty" bson:"cron_expression,omitempty" binding:"omitempty,cron"` // If provided, TimeRange and DaysOfWeek are ignored
	Timezone       string     `json:"timezone" bson:"timezone" binding:"required,timezone"`
	TimeRange      *TimeRange `json:"time_range,omitempty" bson:"time_range,omitempty" binding:"omitempty"`                      // Used only if CronExpression is not provided
	DaysOfWeek     []int      `json:"days_of_week,omitempty" bson:"days_of_week,omitempty" binding:"omitempty,dive,min=0,max=6"` // Used only if CronExpression is not provided
	Exclusions     []int      `json:"exclusions,omitempty" bson:"exclusions,omitempty" binding:"omitempty,dive,min=0,max=6"`
}

// FrequencyUnit defines the unit for frequency
type FrequencyUnit string

const (
	FrequencyUnitSecond FrequencyUnit = "s"
	FrequencyUnitMinute FrequencyUnit = "m"
	FrequencyUnitHour   FrequencyUnit = "h"
)

// Frequency defines how often a task should run within a time range
type Frequency struct {
	Value int           `json:"value" bson:"value" binding:"required,min=1"`     // Numeric value (e.g., 15)
	Unit  FrequencyUnit `json:"unit" bson:"unit" binding:"required,oneof=s m h"` // Unit: "s" (seconds), "m" (minutes), "h" (hours)
}

// TimeRange defines a time range for task execution with frequency
type TimeRange struct {
	Start     string     `json:"start" bson:"start" binding:"required,time_format"` // Format: "HH:MM"
	End       string     `json:"end" bson:"end" binding:"required,time_format"`     // Format: "HH:MM"
	Frequency *Frequency `json:"frequency" bson:"frequency" binding:"required"`     // Frequency with value and unit (e.g., {value: 15, unit: "m"})
}

// CreateTaskRequest represents the request DTO for creating a task
type CreateTaskRequest struct {
	ProjectID      string                 `json:"project_id" binding:"required,objectid"`
	TaskGroupID    string                 `json:"task_group_id,omitempty" binding:"omitempty,objectid"` // Optional task group ID
	Name           string                 `json:"name" binding:"required,min=1,max=255"`
	Description    string                 `json:"description,omitempty" binding:"omitempty,max=1000"`
	ScheduleType   ScheduleType           `json:"schedule_type" binding:"required,oneof=RECURRING ONEOFF"`
	Status         TaskStatus             `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE PAUSED DISABLED"`
	ScheduleConfig ScheduleConfig         `json:"schedule_config" binding:"required"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateTaskRequest represents the request DTO for full task update (PUT)
// Same structure as CreateTaskRequest but without ProjectID (comes from path parameter)
type UpdateTaskRequest struct {
	TaskGroupID    string                 `json:"task_group_id,omitempty" binding:"omitempty,objectid"` // Optional task group ID
	Name           string                 `json:"name" binding:"required,min=1,max=255"`
	Description    string                 `json:"description,omitempty" binding:"omitempty,max=1000"`
	ScheduleType   ScheduleType           `json:"schedule_type" binding:"required,oneof=RECURRING ONEOFF"`
	Status         TaskStatus             `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE PAUSED DISABLED"`
	ScheduleConfig ScheduleConfig         `json:"schedule_config" binding:"required"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// TriggerType defines the type of trigger
type TriggerType string

const (
	TriggerTypeHTTP TriggerType = "HTTP"
)

// HTTPTriggerConfig holds the HTTP trigger configuration
type HTTPTriggerConfig struct {
	URL     string            `json:"url" bson:"url" binding:"required,url"`
	Method  string            `json:"method" bson:"method" binding:"required,http_method"`
	Headers map[string]string `json:"headers,omitempty" bson:"headers,omitempty"`
	Body    interface{}       `json:"body,omitempty" bson:"body,omitempty"`
	Timeout int               `json:"timeout,omitempty" bson:"timeout,omitempty" binding:"omitempty,min=1,max=300"`
}

// TriggerConfig holds the trigger configuration for a task
type TriggerConfig struct {
	Type TriggerType        `json:"type,omitempty" bson:"type,omitempty" binding:"omitempty,oneof=HTTP"`
	HTTP *HTTPTriggerConfig `json:"http,omitempty" bson:"http,omitempty" binding:"omitempty"`
}
