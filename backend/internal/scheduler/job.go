package scheduler

import (
	"log"

	"github.com/yourusername/cron-observer/backend/internal/models"
)

// TaskJob represents a cron job for a task
type TaskJob struct {
	Task *models.Task
}

// Run executes the task job (currently just logs)
func (j *TaskJob) Run() {
	log.Printf("[CRON] Task triggered: %s (UUID: %s)", j.Task.Name, j.Task.UUID)
	// Future: Execute HTTP trigger based on j.Task.TriggerConfig
}
