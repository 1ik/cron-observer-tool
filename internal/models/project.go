package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project represents a project entity that contains tasks
type Project struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	APIKey      string             `json:"api_key" bson:"api_key"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "ACTIVE"
	ProjectStatusInactive ProjectStatus = "INACTIVE"
)
