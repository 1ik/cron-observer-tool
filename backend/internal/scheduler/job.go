package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskJob represents a cron job for a task
type TaskJob struct {
	Task *models.Task
	Repo repositories.Repository
}

// Run executes the task job
func (j *TaskJob) Run() {
	ctx := context.Background()
	log.Printf("[CRON] Task triggered: %s (UUID: %s)", j.Task.Name, j.Task.UUID)

	// Get the project to retrieve execution_endpoint
	project, err := j.Repo.GetProjectByID(ctx, j.Task.ProjectID)
	if err != nil {
		log.Printf("[CRON] Failed to get project for task %s: %v", j.Task.UUID, err)
		return
	}

	// Check if execution_endpoint is set
	if project.ExecutionEndpoint == "" {
		log.Printf("[CRON] No execution_endpoint set for project %s, skipping execution", project.UUID)
		return
	}

	// Create execution record
	executionUUID := uuid.New().String()
	executionID := primitive.NewObjectID()
	now := time.Now()

	execution := &models.Execution{
		ID:        executionID,
		UUID:      executionUUID,
		TaskID:    j.Task.ID,
		TaskUUID:  j.Task.UUID,
		Status:    models.ExecutionStatusPending,
		StartedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save execution record
	if err := j.Repo.CreateExecution(ctx, execution); err != nil {
		log.Printf("[CRON] Failed to create execution record for task %s: %v", j.Task.UUID, err)
		return
	}

	// Prepare request body with task name and execution ID
	requestBody := map[string]interface{}{
		"task_name":    j.Task.Name,
		"execution_id": executionUUID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("[CRON] Failed to marshal request body for task %s: %v", j.Task.UUID, err)
		return
	}

	// Send POST request to execution_endpoint
	req, err := http.NewRequest("POST", project.ExecutionEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("[CRON] Failed to create HTTP request for task %s: %v", j.Task.UUID, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[CRON] Failed to send POST request for task %s: %v", j.Task.UUID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("[CRON] Successfully executed task %s (execution: %s)", j.Task.UUID, executionUUID)
	} else {
		log.Printf("[CRON] Execution endpoint returned non-2xx status for task %s: %d", j.Task.UUID, resp.StatusCode)
	}
}
