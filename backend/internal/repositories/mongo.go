package repositories

import (
	"context"

	"github.com/yourusername/cron-observer/backend/internal/database"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *MongoRepository) GetTasksByProjectID(ctx context.Context, projectID primitive.ObjectID) ([]*models.Task, error) {
	collection := r.db.Collection(database.CollectionTasks)
	cursor, err := collection.Find(ctx, bson.M{"project_id": projectID})
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

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}
