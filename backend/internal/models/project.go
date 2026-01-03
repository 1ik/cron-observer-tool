package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project represents a project entity that contains tasks
// @Description Project represents a project entity that contains tasks
type Project struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"507f1f77bcf86cd799439011"`
	UUID        string             `json:"uuid" bson:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string             `json:"name" bson:"name" example:"My Project"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" example:"Project description"`
	APIKey      string             `json:"api_key" bson:"api_key" example:"sk_live_abc123..."`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at" example:"2025-01-15T10:00:00Z"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at" example:"2025-01-15T10:00:00Z"`
}

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "ACTIVE"
	ProjectStatusInactive ProjectStatus = "INACTIVE"
)
