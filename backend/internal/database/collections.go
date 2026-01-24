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
	CollectionProjects              = "projects"
	CollectionTasks                 = "tasks"
	CollectionTaskGroups            = "task_groups"
	CollectionExecutions            = "executions"
	CollectionExecutionFailureStats = "execution_failure_stats"
	CollectionTaskFailureStats      = "task_failure_stats"
)

// GetProjectsCollection returns the projects collection
func (d *Database) GetProjectsCollection() *mongo.Collection {
	return d.DB.Collection(CollectionProjects)
}

// GetTasksCollection returns the tasks collection
func (d *Database) GetTasksCollection() *mongo.Collection {
	return d.DB.Collection(CollectionTasks)
}

// GetTaskGroupsCollection returns the task_groups collection
func (d *Database) GetTaskGroupsCollection() *mongo.Collection {
	return d.DB.Collection(CollectionTaskGroups)
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

	// Create indexes for task_groups collection
	if err := d.createTaskGroupIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create task group indexes: %w", err)
	}

	// Create indexes for execution_failure_stats collection
	if err := d.createExecutionFailureStatsIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create execution failure stats indexes: %w", err)
	}

	// Create indexes for task_failure_stats collection
	if err := d.createTaskFailureStatsIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create task failure stats indexes: %w", err)
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
		{
			Keys:    bson.D{{Key: "task_group_id", Value: 1}},
			Options: options.Index().SetName("idx_task_group_id"),
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

// createTaskGroupIndexes creates indexes for the task_groups collection
func (d *Database) createTaskGroupIndexes(ctx context.Context) error {
	collection := d.GetTaskGroupsCollection()
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

// createExecutionFailureStatsIndexes creates indexes for the execution_failure_stats collection
func (d *Database) createExecutionFailureStatsIndexes(ctx context.Context) error {
	collection := d.DB.Collection(CollectionExecutionFailureStats)
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "date", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("idx_project_date"),
		},
		{
			Keys:    bson.D{{Key: "project_id", Value: 1}},
			Options: options.Index().SetName("idx_project_id"),
		},
		{
			Keys:    bson.D{{Key: "date", Value: -1}},
			Options: options.Index().SetName("idx_date"),
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

// createTaskFailureStatsIndexes creates indexes for the task_failure_stats collection
func (d *Database) createTaskFailureStatsIndexes(ctx context.Context) error {
	collection := d.DB.Collection(CollectionTaskFailureStats)
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "date", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("idx_project_date"),
		},
		{
			Keys:    bson.D{{Key: "project_id", Value: 1}},
			Options: options.Index().SetName("idx_project_id"),
		},
		{
			Keys:    bson.D{{Key: "date", Value: -1}},
			Options: options.Index().SetName("idx_date"),
		},
		{
			Keys:    bson.D{{Key: "calculated_at", Value: -1}},
			Options: options.Index().SetName("idx_calculated_at"),
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
