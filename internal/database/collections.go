package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Collection names
	CollectionProjects = "projects"
	CollectionTasks    = "tasks"
)

// GetProjectsCollection returns the projects collection
func (d *Database) GetProjectsCollection() *mongo.Collection {
	return d.DB.Collection(CollectionProjects)
}

// GetTasksCollection returns the tasks collection
func (d *Database) GetTasksCollection() *mongo.Collection {
	return d.DB.Collection(CollectionTasks)
}

// CreateIndexes creates all necessary indexes for collections
func (d *Database) CreateIndexes(ctx context.Context) error {
	// Create indexes for projects collection
	if err := d.createProjectIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create project indexes: %w", err)
	}

	// Create indexes for tasks collection
	if err := d.createTaskIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create task indexes: %w", err)
	}

	return nil
}

// createProjectIndexes creates indexes for the projects collection
func (d *Database) createProjectIndexes(ctx context.Context) error {
	collection := d.GetProjectsCollection()
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_uuid"),
		},
		{
			Keys:    bson.D{{Key: "api_key", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_api_key"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// createTaskIndexes creates indexes for the tasks collection
func (d *Database) createTaskIndexes(ctx context.Context) error {
	collection := d.GetTasksCollection()
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_uuid"),
		},
		{
			Keys:    bson.D{{Key: "project_id", Value: 1}},
			Options: options.Index().SetName("idx_project_id"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "schedule_type", Value: 1}},
			Options: options.Index().SetName("idx_schedule_type"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().SetName("idx_project_status"),
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("idx_project_created"),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
