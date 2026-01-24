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

func (r *MongoRepository) GetUserProjects(ctx context.Context, email string) ([]*models.Project, error) {
	collection := r.db.Collection(database.CollectionProjects)

	// Find projects where the user's email exists in the project_users array
	filter := bson.M{
		"project_users.email": email,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*models.Project
	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, err
	}

	// Ensure project_users is always initialized (not nil) for JSON serialization
	for _, project := range projects {
		if project.ProjectUsers == nil {
			project.ProjectUsers = []models.ProjectUser{}
		}
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

func (r *MongoRepository) GetExecutionsByTaskUUIDPaginated(ctx context.Context, taskUUID string, startDate, endDate *time.Time, page, pageSize int) ([]*models.Execution, int64, error) {
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

	// Get total count
	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip value
	skip := (page - 1) * pageSize

	// Set up pagination options
	opts := options.Find().
		SetSort(bson.M{"started_at": -1}). // Most recent first
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var executions []*models.Execution
	err = cursor.All(ctx, &executions)
	if err != nil {
		return nil, 0, err
	}

	// Ensure we always return an empty slice instead of nil
	if executions == nil {
		executions = []*models.Execution{}
	}

	return executions, totalCount, nil
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

func (r *MongoRepository) GetExecutionByUUID(ctx context.Context, executionUUID string) (*models.Execution, error) {
	collection := r.db.Collection(database.CollectionExecutions)

	var execution models.Execution
	err := collection.FindOne(ctx, bson.M{"uuid": executionUUID}).Decode(&execution)
	if err != nil {
		return nil, err
	}

	return &execution, nil
}

func (r *MongoRepository) IncrementFailureStat(ctx context.Context, projectID primitive.ObjectID, date string) error {
	collection := r.db.Collection(database.CollectionExecutionFailureStats)

	filter := bson.M{
		"project_id": projectID,
		"date":        date,
	}

	update := bson.M{
		"$inc": bson.M{"count": 1},
		"$set": bson.M{"updated_at": time.Now()},
		"$setOnInsert": bson.M{
			"project_id": projectID,
			"date":       date,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoRepository) GetFailureStatsByProject(ctx context.Context, projectID primitive.ObjectID, days int) ([]*models.FailedExecutionStats, int, error) {
	collection := r.db.Collection(database.CollectionExecutionFailureStats)

	// Calculate date range (last N days)
	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, -days)
	startDateStr := startDate.Format("2006-01-02")

	// Build filter
	filter := bson.M{
		"project_id": projectID,
		"date": bson.M{
			"$gte": startDateStr,
		},
	}

	// Find all stats for the project in the date range
	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "date", Value: -1}}))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var stats []*models.ExecutionFailureStat
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, 0, err
	}

	// Convert to response format
	result := make([]*models.FailedExecutionStats, 0, len(stats))
	total := 0
	for _, stat := range stats {
		result = append(result, &models.FailedExecutionStats{
			Date:  stat.Date,
			Count: stat.Count,
		})
		total += stat.Count
	}

	return result, total, nil
}

func (r *MongoRepository) GetExecutionStatsByProject(ctx context.Context, projectID primitive.ObjectID, days int) ([]*models.ExecutionStats, error) {
	collection := r.db.Collection(database.CollectionExecutions)

	// Calculate date range (last N days)
	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, -days)
	startOfDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)

	// Get all tasks for this project
	tasksCollection := r.db.Collection(database.CollectionTasks)
	taskCursor, err := tasksCollection.Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return nil, err
	}
	defer taskCursor.Close(ctx)

	var tasks []models.Task
	if err := taskCursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []*models.ExecutionStats{}, nil
	}

	// Get task IDs
	taskIDs := make([]primitive.ObjectID, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}

	// Aggregate executions by date and status
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"task_id": bson.M{"$in": taskIDs},
				"started_at": bson.M{
					"$gte": startOfDay,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"date": bson.M{
						"$dateToString": bson.M{
							"format": "%Y-%m-%d",
							"date":   "$started_at",
						},
					},
					"status": "$status",
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id.date",
				"statuses": bson.M{
					"$push": bson.M{
						"status": "$_id.status",
						"count":  "$count",
					},
				},
			},
		},
		{
			"$sort": bson.M{"_id": -1},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type aggregationResult struct {
		ID       string `bson:"_id"`
		Statuses []struct {
			Status string `bson:"status"`
			Count  int    `bson:"count"`
		} `bson:"statuses"`
	}

	var results []aggregationResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Convert to ExecutionStats format
	stats := make([]*models.ExecutionStats, 0, len(results))
	for _, result := range results {
		stat := &models.ExecutionStats{
			Date:     result.ID,
			Failures: 0,
			Success:  0,
			Total:    0,
		}

		for _, statusCount := range result.Statuses {
			count := statusCount.Count
			stat.Total += count

			switch models.ExecutionStatus(statusCount.Status) {
			case models.ExecutionStatusFailed:
				stat.Failures += count
			case models.ExecutionStatusSuccess:
				stat.Success += count
			}
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}
