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

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}
