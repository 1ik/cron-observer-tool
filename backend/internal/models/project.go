package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project represents a project entity that contains tasks
// @Description Project represents a project entity that contains tasks
type Project struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"507f1f77bcf86cd799439011"`
	UUID              string             `json:"uuid" bson:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name              string             `json:"name" bson:"name" example:"My Project"`
	Description       string             `json:"description,omitempty" bson:"description,omitempty" example:"Project description"`
	APIKey            string             `json:"api_key" bson:"api_key" example:"sk_live_abc123..."`
	ExecutionEndpoint string             `json:"execution_endpoint" bson:"execution_endpoint" binding:"omitempty,url" example:"https://api.example.com/execute"`
	AlertEmails       string             `json:"alert_emails,omitempty" bson:"alert_emails,omitempty" example:"admin@example.com,ops@example.com"`
	ProjectUsers      []ProjectUser      `json:"project_users" bson:"project_users,omitempty"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at" example:"2025-01-15T10:00:00Z"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at" example:"2025-01-15T10:00:00Z"`
}

// CreateProjectRequest represents the request DTO for creating a project
type CreateProjectRequest struct {
	Name              string `json:"name" binding:"required,min=1,max=255"`
	Description       string `json:"description,omitempty" binding:"omitempty,max=1000"`
	ExecutionEndpoint string `json:"execution_endpoint,omitempty" binding:"omitempty,url"`
}

// UpdateProjectRequest represents the request DTO for updating a project
type UpdateProjectRequest struct {
	Name              string        `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description       string        `json:"description,omitempty" binding:"omitempty,max=1000"`
	ExecutionEndpoint string        `json:"execution_endpoint,omitempty" binding:"omitempty,url"`
	AlertEmails       string        `json:"alert_emails,omitempty" binding:"omitempty"`
	ProjectUsers      []ProjectUser `json:"project_users,omitempty" binding:"omitempty,dive"`
}

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "ACTIVE"
	ProjectStatusInactive ProjectStatus = "INACTIVE"
)

// ProjectUserRole represents the role of a user in a project
type ProjectUserRole string

const (
	ProjectUserRoleAdmin    ProjectUserRole = "admin"
	ProjectUserRoleReadonly ProjectUserRole = "readonly"
)

// ProjectUser represents a user associated with a project
// @Description ProjectUser represents a user associated with a project
type ProjectUser struct {
	Email string          `json:"email" bson:"email" binding:"required,email" example:"user@example.com"`
	Role  ProjectUserRole `json:"role" bson:"role" binding:"required,oneof=admin readonly" example:"admin"`
}
