package repositories

import (
	"context"
	"time"

	"github.com/yourusername/cron-observer/backend/internal/database"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *MongoRepository) GetProjectByID(ctx context.Context, projectID primitive.ObjectID) (*models.Project, error) {
	collection := r.db.Collection(database.CollectionProjects)

	var project models.Project
	err := collection.FindOne(ctx, bson.M{"_id": projectID}).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *MongoRepository) CreateProject(ctx context.Context, project *models.Project) error {
	collection := r.db.Collection(database.CollectionProjects)
	_, err := collection.InsertOne(ctx, project)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) UpdateProject(ctx context.Context, projectID primitive.ObjectID, project *models.Project) error {
	collection := r.db.Collection(database.CollectionProjects)

	update := bson.M{
		"$set": bson.M{
			"name":               project.Name,
			"description":        project.Description,
			"execution_endpoint": project.ExecutionEndpoint,
			"alert_emails":       project.AlertEmails,
			"updated_at":         project.UpdatedAt,
		},
	}

	// Always include project_users in the update (even if empty array)
	// This ensures the field exists in MongoDB
	// If nil, initialize as empty array
	projectUsers := project.ProjectUsers
	if projectUsers == nil {
		projectUsers = []models.ProjectUser{}
	}
	update["$set"].(bson.M)["project_users"] = projectUsers

	_, err := collection.UpdateOne(ctx, bson.M{"_id": projectID}, update)
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

func (r *MongoRepository) UpdateTaskStatus(ctx context.Context, taskUUID string, status models.TaskStatus) error {
	collection := r.db.Collection(database.CollectionTasks)

	filter := bson.M{"uuid": taskUUID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) UpdateTaskState(ctx context.Context, taskUUID string, state models.TaskState) error {
	collection := r.db.Collection(database.CollectionTasks)

	filter := bson.M{"uuid": taskUUID}
	update := bson.M{
		"$set": bson.M{
			"state":      state,
			"updated_at": time.Now(),
		},
	}

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

func (r *MongoRepository) UpdateTaskGroupStatus(ctx context.Context, taskGroupUUID string, status models.TaskGroupStatus) error {
	collection := r.db.Collection(database.CollectionTaskGroups)

	filter := bson.M{"uuid": taskGroupUUID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) UpdateTaskGroupState(ctx context.Context, taskGroupUUID string, state models.TaskGroupState) error {
	collection := r.db.Collection(database.CollectionTaskGroups)

	filter := bson.M{"uuid": taskGroupUUID}
	update := bson.M{
		"$set": bson.M{
			"state":      state,
			"updated_at": time.Now(),
		},
	}

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

func (r *MongoRepository) CreateExecution(ctx context.Context, execution *models.Execution) error {
	collection := r.db.Collection(database.CollectionExecutions)
	_, err := collection.InsertOne(ctx, execution)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) GetExecutionsByTaskUUID(ctx context.Context, taskUUID string, startDate, endDate *time.Time) ([]*models.Execution, error) {
	collection := r.db.Collection(database.CollectionExecutions)

	filter := bson.M{"task_uuid": taskUUID}

	// Add date filtering if provided
	if startDate != nil || endDate != nil {
		dateFilter := bson.M{}
		if startDate != nil {
			// Ensure startDate is in UTC for MongoDB comparison
			startUTC := startDate.UTC()
			dateFilter["$gte"] = startUTC
		}
		if endDate != nil {
			// Ensure endDate is in UTC for MongoDB comparison
			endUTC := endDate.UTC()
			dateFilter["$lte"] = endUTC
		}
		filter["started_at"] = dateFilter
	}

	opts := options.Find().SetSort(bson.M{"started_at": -1}) // Most recent first
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var executions []*models.Execution
	err = cursor.All(ctx, &executions)
	if err != nil {
		return nil, err
	}

	// Ensure we always return an empty slice instead of nil
	if executions == nil {
		executions = []*models.Execution{}
	}

	return executions, nil
}

func (r *MongoRepository) AppendLogToExecution(ctx context.Context, executionUUID string, logEntry models.LogEntry) error {
	collection := r.db.Collection(database.CollectionExecutions)

	filter := bson.M{"uuid": executionUUID}
	update := bson.M{
		"$push": bson.M{
			"logs": logEntry,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) UpdateExecutionStatus(ctx context.Context, executionUUID string, status models.ExecutionStatus, errorMessage *string) error {
	collection := r.db.Collection(database.CollectionExecutions)

	filter := bson.M{"uuid": executionUUID}
	now := time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": now,
		},
	}

	// Set ended_at if status is SUCCESS or FAILED
	if status == models.ExecutionStatusSuccess || status == models.ExecutionStatusFailed {
		update["$set"].(bson.M)["ended_at"] = now
	}

	// Set error message if provided
	if errorMessage != nil {
		update["$set"].(bson.M)["error"] = *errorMessage
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}
