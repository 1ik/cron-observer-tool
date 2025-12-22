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
