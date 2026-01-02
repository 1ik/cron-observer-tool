package repositories

import (
	"context"

	"github.com/yourusername/cron-observer/backend/internal/database"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	db *mongo.Database
}

func (r *MongoRepository) GetAllProjects(ctx context.Context) ([]*models.Project, error) {
	collection := r.db.Collection(database.CollectionProjects)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*models.Project
	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *MongoRepository) CreateProject(ctx context.Context, project *models.Project) error {
	collection := r.db.Collection(database.CollectionProjects)
	_, err := collection.InsertOne(ctx, project)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) CreateTask(ctx context.Context, projectID string, task *models.Task) error {
	collection := r.db.Collection(database.CollectionTasks)
	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) GetAllActiveTasks(ctx context.Context) ([]*models.Task, error) {
	collection := r.db.Collection(database.CollectionTasks)

	// Filter for active tasks with cron expressions
	filter := bson.M{
		"status":                          models.TaskStatusActive,
		"schedule_config.cron_expression": bson.M{"$ne": ""},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *MongoRepository) GetTaskByUUID(ctx context.Context, taskUUID string) (*models.Task, error) {
	collection := r.db.Collection(database.CollectionTasks)

	var task models.Task
	err := collection.FindOne(ctx, bson.M{"uuid": taskUUID}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *MongoRepository) UpdateTask(ctx context.Context, taskUUID string, task *models.Task) error {
	collection := r.db.Collection(database.CollectionTasks)

	filter := bson.M{"uuid": taskUUID}
	update := bson.M{"$set": task}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) DeleteTask(ctx context.Context, taskUUID string) error {
	collection := r.db.Collection(database.CollectionTasks)

	_, err := collection.DeleteOne(ctx, bson.M{"uuid": taskUUID})
	return err
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}
