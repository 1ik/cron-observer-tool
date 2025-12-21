package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a scheduled task entity
type Task struct {
	ID                 primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UUID               string                 `json:"uuid" bson:"uuid"`
	ProjectID          primitive.ObjectID     `json:"project_id" bson:"project_id"`
	Name               string                 `json:"name" bson:"name"`
	Description        string                 `json:"description,omitempty" bson:"description,omitempty"`
	ScheduleType       ScheduleType           `json:"schedule_type" bson:"schedule_type"`
	Status             TaskStatus             `json:"status" bson:"status"`
	ScheduleConfig     ScheduleConfig         `json:"schedule_config" bson:"schedule_config"`
	Metadata           map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	NotificationConfig *NotificationConfig    `json:"notification_config,omitempty" bson:"notification_config,omitempty"`
	CreatedAt          time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" bson:"updated_at"`
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
type ScheduleConfig struct {
	// For RECURRING tasks
	CronExpression string     `json:"cron_expression,omitempty" bson:"cron_expression,omitempty"`
	Timezone       string     `json:"timezone" bson:"timezone"`
	TimeRange      *TimeRange `json:"time_range,omitempty" bson:"time_range,omitempty"`
	DaysOfWeek     []int      `json:"days_of_week,omitempty" bson:"days_of_week,omitempty"`
	Exclusions     []int      `json:"exclusions,omitempty" bson:"exclusions,omitempty"`

	// For ONEOFF tasks
	ExecuteAt *time.Time `json:"execute_at,omitempty" bson:"execute_at,omitempty"`
}

// TimeRange defines a time range for task execution
type TimeRange struct {
	Start string `json:"start" bson:"start"` // Format: "HH:MM"
	End   string `json:"end" bson:"end"`     // Format: "HH:MM"
}

// NotificationConfig holds notification settings for a task
type NotificationConfig struct {
	OnSuccess bool                  `json:"on_success" bson:"on_success"`
	OnFailure bool                  `json:"on_failure" bson:"on_failure"`
	Channels  []NotificationChannel `json:"channels" bson:"channels"`
}

// NotificationChannel defines a notification channel
type NotificationChannel struct {
	Type     NotificationChannelType `json:"type" bson:"type"`
	Endpoint string                  `json:"endpoint" bson:"endpoint"`
	Template string                  `json:"template,omitempty" bson:"template,omitempty"`
}

// NotificationChannelType defines notification channel types
type NotificationChannelType string

const (
	NotificationChannelEmail   NotificationChannelType = "EMAIL"
	NotificationChannelWebhook NotificationChannelType = "WEBHOOK"
	NotificationChannelSlack   NotificationChannelType = "SLACK"
	NotificationChannelDiscord NotificationChannelType = "DISCORD"
)
