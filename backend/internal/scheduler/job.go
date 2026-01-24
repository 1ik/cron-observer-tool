package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// ExecuteTask creates an execution record and sends it to the execution endpoint.
// Returns the execution UUID and any error encountered during execution creation.
// The actual HTTP request to the execution endpoint is sent asynchronously.
func ExecuteTask(ctx context.Context, task *models.Task, repo repositories.Repository, logPrefix string) (string, error) {
	// Get the project to retrieve execution_endpoint
	project, err := repo.GetProjectByID(ctx, task.ProjectID)
	if err != nil {
		log.Printf("[%s] Failed to get project for task %s: %v", logPrefix, task.UUID, err)
		return "", err
	}

	// Check if execution_endpoint is set
	if project.ExecutionEndpoint == "" {
		log.Printf("[%s] No execution_endpoint set for project %s, skipping execution", logPrefix, project.UUID)
		return "", fmt.Errorf("no execution_endpoint set for project")
	}

	// Create execution record
	executionUUID := uuid.New().String()
	executionID := primitive.NewObjectID()
	now := time.Now()

	execution := &models.Execution{
		ID:        executionID,
		UUID:      executionUUID,
		TaskID:    task.ID,
		TaskUUID:  task.UUID,
		Status:    models.ExecutionStatusPending,
		StartedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save execution record
	if err := repo.CreateExecution(ctx, execution); err != nil {
		log.Printf("[%s] Failed to create execution record for task %s: %v", logPrefix, task.UUID, err)
		return "", err
	}

	// Send execution to the execution endpoint asynchronously (don't wait for response)
	go func() {
		// Prepare request body with task name and execution ID
		requestBody := map[string]interface{}{
			"task_name":    task.Name,
			"execution_id": executionUUID,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			log.Printf("[%s] Failed to marshal request body for task %s: %v", logPrefix, task.UUID, err)
			return
		}

		// Send POST request to execution_endpoint
		req, err := http.NewRequest("POST", project.ExecutionEndpoint, bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Printf("[%s] Failed to create HTTP request for task %s: %v", logPrefix, task.UUID, err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[%s] Failed to send POST request for task %s: %v", logPrefix, task.UUID, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("[%s] Successfully executed task %s (execution: %s)", logPrefix, task.UUID, executionUUID)
		} else {
			log.Printf("[%s] Execution endpoint returned non-2xx status for task %s: %d", logPrefix, task.UUID, resp.StatusCode)
		}
	}()

	return executionUUID, nil
}

// Run executes the task job
func (j *TaskJob) Run() {
	ctx := context.Background()
	log.Printf("[CRON] Task triggered: %s (UUID: %s)", j.Task.Name, j.Task.UUID)

	_, err := ExecuteTask(ctx, j.Task, j.Repo, "CRON")
	if err != nil {
		// Error already logged in ExecuteTask
		return
	}
}
