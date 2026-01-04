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

func (r *MongoRepository) GetTasksByProjectID(ctx context.Context, projectID primitive.ObjectID) ([]*models.Task, error) {
	collection := r.db.Collection(database.CollectionTasks)

	filter := bson.M{"project_id": projectID}

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

// TaskGroup repository methods

func (r *MongoRepository) CreateTaskGroup(ctx context.Context, projectID string, taskGroup *models.TaskGroup) error {
	collection := r.db.Collection(database.CollectionTaskGroups)
	_, err := collection.InsertOne(ctx, taskGroup)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) GetTaskGroupsByProjectID(ctx context.Context, projectID primitive.ObjectID) ([]*models.TaskGroup, error) {
	collection := r.db.Collection(database.CollectionTaskGroups)

	filter := bson.M{"project_id": projectID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var taskGroups []*models.TaskGroup
	err = cursor.All(ctx, &taskGroups)
	if err != nil {
		return nil, err
	}
	return taskGroups, nil
}

func (r *MongoRepository) GetTaskGroupByUUID(ctx context.Context, taskGroupUUID string) (*models.TaskGroup, error) {
	collection := r.db.Collection(database.CollectionTaskGroups)

	var taskGroup models.TaskGroup
	err := collection.FindOne(ctx, bson.M{"uuid": taskGroupUUID}).Decode(&taskGroup)
	if err != nil {
		return nil, err
	}
	return &taskGroup, nil
}

func (r *MongoRepository) GetTaskGroupByID(ctx context.Context, taskGroupID primitive.ObjectID) (*models.TaskGroup, error) {
	collection := r.db.Collection(database.CollectionTaskGroups)

	var taskGroup models.TaskGroup
	err := collection.FindOne(ctx, bson.M{"_id": taskGroupID}).Decode(&taskGroup)
	if err != nil {
		return nil, err
	}
	return &taskGroup, nil
}

func (r *MongoRepository) UpdateTaskGroup(ctx context.Context, taskGroupUUID string, taskGroup *models.TaskGroup) error {
	collection := r.db.Collection(database.CollectionTaskGroups)

	filter := bson.M{"uuid": taskGroupUUID}
	update := bson.M{"$set": taskGroup}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) DeleteTaskGroup(ctx context.Context, taskGroupUUID string) error {
	collection := r.db.Collection(database.CollectionTaskGroups)

	_, err := collection.DeleteOne(ctx, bson.M{"uuid": taskGroupUUID})
	return err
}

func (r *MongoRepository) GetTasksByGroupID(ctx context.Context, taskGroupID primitive.ObjectID) ([]*models.Task, error) {
	collection := r.db.Collection(database.CollectionTasks)

	filter := bson.M{"task_group_id": taskGroupID}

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

func (r *MongoRepository) GetActiveTaskGroupsWithWindows(ctx context.Context) ([]*models.TaskGroup, error) {
	collection := r.db.Collection(database.CollectionTaskGroups)

	// Filter for active groups with start and end times
	filter := bson.M{
		"status":     models.TaskGroupStatusActive,
		"start_time": bson.M{"$ne": ""},
		"end_time":   bson.M{"$ne": ""},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var taskGroups []*models.TaskGroup
	err = cursor.All(ctx, &taskGroups)
	if err != nil {
		return nil, err
	}
	return taskGroups, nil
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}
