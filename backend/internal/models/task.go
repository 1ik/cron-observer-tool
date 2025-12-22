package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a scheduled task entity
type Task struct {
	ID             primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UUID           string                 `json:"uuid" bson:"uuid"`
	ProjectID      primitive.ObjectID     `json:"project_id" bson:"project_id"`
	Name           string                 `json:"name" bson:"name"`
	Description    string                 `json:"description,omitempty" bson:"description,omitempty"`
	ScheduleType   ScheduleType           `json:"schedule_type" bson:"schedule_type"`
	Status         TaskStatus             `json:"status" bson:"status"`
	ScheduleConfig ScheduleConfig         `json:"schedule_config" bson:"schedule_config"`
	TriggerConfig  TriggerConfig          `json:"trigger_config" bson:"trigger_config"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
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
	CronExpression string     `json:"cron_expression,omitempty" bson:"cron_expression,omitempty"` // If provided, TimeRange and DaysOfWeek are ignored
	Timezone       string     `json:"timezone" bson:"timezone"`
	TimeRange      *TimeRange `json:"time_range,omitempty" bson:"time_range,omitempty"`     // Used only if CronExpression is not provided
	DaysOfWeek     []int      `json:"days_of_week,omitempty" bson:"days_of_week,omitempty"` // Used only if CronExpression is not provided
	Exclusions     []int      `json:"exclusions,omitempty" bson:"exclusions,omitempty"`
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
	Value int           `json:"value" bson:"value"` // Numeric value (e.g., 15)
	Unit  FrequencyUnit `json:"unit" bson:"unit"`   // Unit: "s" (seconds), "m" (minutes), "h" (hours)
}

// TimeRange defines a time range for task execution with frequency
type TimeRange struct {
	Start     string     `json:"start" bson:"start"`         // Format: "HH:MM"
	End       string     `json:"end" bson:"end"`             // Format: "HH:MM"
	Frequency *Frequency `json:"frequency" bson:"frequency"` // Frequency with value and unit (e.g., {value: 15, unit: "m"})
}

// CreateTaskRequest represents the request DTO for creating a task
type CreateTaskRequest struct {
	ProjectID      string                 `json:"project_id" binding:"required,objectid"`
	Name           string                 `json:"name" binding:"required,min=1,max=255"`
	Description    string                 `json:"description,omitempty" binding:"omitempty,max=1000"`
	ScheduleType   ScheduleType           `json:"schedule_type" binding:"required,oneof=RECURRING ONEOFF"`
	Status         TaskStatus             `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE PAUSED DISABLED"`
	ScheduleConfig CreateScheduleConfig   `json:"schedule_config" binding:"required"`
	TriggerConfig  CreateTriggerConfig    `json:"trigger_config" binding:"required"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// CreateScheduleConfig represents the schedule configuration in the request
type CreateScheduleConfig struct {
	CronExpression string           `json:"cron_expression,omitempty" binding:"omitempty,cron"`
	Timezone       string           `json:"timezone" binding:"required,timezone"`
	TimeRange      *CreateTimeRange `json:"time_range,omitempty" binding:"omitempty"`
	DaysOfWeek     []int            `json:"days_of_week,omitempty" binding:"omitempty,dive,min=0,max=6"`
	Exclusions     []int            `json:"exclusions,omitempty" binding:"omitempty,dive,min=0,max=6"`
}

// CreateTimeRange represents the time range in the request
type CreateTimeRange struct {
	Start     string           `json:"start" binding:"required,time_format"`
	End       string           `json:"end" binding:"required,time_format"`
	Frequency *CreateFrequency `json:"frequency" binding:"required"`
}

// CreateFrequency represents the frequency in the request
type CreateFrequency struct {
	Value int           `json:"value" binding:"required,min=1"`
	Unit  FrequencyUnit `json:"unit" binding:"required,oneof=s m h"`
}

// TriggerType defines the type of trigger
type TriggerType string

const (
	TriggerTypeHTTP TriggerType = "HTTP"
)

// HTTPTriggerConfig holds the HTTP trigger configuration
type HTTPTriggerConfig struct {
	URL     string            `json:"url" bson:"url"`
	Method  string            `json:"method" bson:"method"`
	Headers map[string]string `json:"headers,omitempty" bson:"headers,omitempty"`
	Body    interface{}       `json:"body,omitempty" bson:"body,omitempty"`
	Timeout int               `json:"timeout,omitempty" bson:"timeout,omitempty"`
}

// TriggerConfig holds the trigger configuration for a task
type TriggerConfig struct {
	Type TriggerType        `json:"type" bson:"type"`
	HTTP *HTTPTriggerConfig `json:"http" bson:"http"`
}

// CreateHTTPTriggerConfig represents the HTTP trigger configuration in the request
type CreateHTTPTriggerConfig struct {
	URL     string            `json:"url" binding:"required,url"`
	Method  string            `json:"method" binding:"required,http_method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    interface{}       `json:"body,omitempty"`
	Timeout int               `json:"timeout,omitempty" binding:"omitempty,min=1,max=300"`
}

// CreateTriggerConfig represents the trigger configuration in the request
type CreateTriggerConfig struct {
	Type TriggerType             `json:"type" binding:"required,oneof=HTTP"`
	HTTP CreateHTTPTriggerConfig `json:"http" binding:"required"`
}
