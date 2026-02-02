package deleteworker

import (
	"context"
	"errors"
	"log"

	"github.com/yourusername/cron-observer/backend/internal/deletequeue"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskUnregisterer is the minimal scheduler interface needed for the delete worker.
type TaskUnregisterer interface {
	UnregisterTask(taskUUID string)
}

// EventPublisher is the minimal event bus interface needed for the delete worker.
type EventPublisher interface {
	Publish(event events.Event)
}

// Worker processes delete job messages: stops cron, hard-deletes the task, publishes TaskDeleted.
type Worker struct {
	repo         repositories.Repository
	scheduler    TaskUnregisterer // optional; nil-safe
	eventPublisher EventPublisher
}

// NewWorker creates a delete worker with the given dependencies.
func NewWorker(repo repositories.Repository, scheduler TaskUnregisterer, eventPublisher EventPublisher) *Worker {
	return &Worker{
		repo:          repo,
		scheduler:     scheduler,
		eventPublisher: eventPublisher,
	}
}

// ProcessDeleteTask performs the delete workflow for one message. Idempotent and retryable.
// Returns nil to ack the message; non-nil to trigger broker retry/DLQ.
func (w *Worker) ProcessDeleteTask(ctx context.Context, msg deletequeue.DeleteTaskMessage) error {
	// Step 1: Fetch task from repository
	task, err := w.repo.GetTaskByUUID(ctx, msg.TaskUUID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("[Worker] Task already deleted (idempotent success): TaskUUID=%s", msg.TaskUUID)
			return nil
		}
		log.Printf("[Worker] ERROR: Failed to fetch task: TaskUUID=%s, error=%v", msg.TaskUUID, err)
		return err
	}

	// Start delete process
	log.Printf("[Worker] Starting task delete process: TaskUUID=%s, TaskName=%s", 
		task.UUID, task.Name)

	// Step 2: Stop cron scheduler
	log.Printf("[Worker] Unregistering task from scheduler: TaskUUID=%s, TaskName=%s", 
		task.UUID, task.Name)
	if w.scheduler != nil {
		w.scheduler.UnregisterTask(task.UUID)
		log.Printf("[Worker] Task unregistered from scheduler: TaskUUID=%s", task.UUID)
	} else {
		log.Printf("[Worker] WARNING: Scheduler is nil, skipping UnregisterTask: TaskUUID=%s", task.UUID)
	}

	// Step 3: Hard delete from MongoDB
	log.Printf("[Worker] Deleting task from database: TaskUUID=%s, TaskName=%s", 
		task.UUID, task.Name)
	if err := w.repo.DeleteTask(ctx, task.UUID); err != nil {
		log.Printf("[Worker] ERROR: Failed to delete task from database: TaskUUID=%s, TaskName=%s, error=%v", 
			task.UUID, task.Name, err)
		
		// Mark as DELETE_FAILED for observability
		if updateErr := w.repo.UpdateTaskStatus(ctx, task.UUID, models.TaskStatusDeleteFailed); updateErr != nil {
			log.Printf("[Worker] WARNING: Failed to update status to DELETE_FAILED: TaskUUID=%s, error=%v", 
				task.UUID, updateErr)
		} else {
			log.Printf("[Worker] Task marked as DELETE_FAILED: TaskUUID=%s, TaskName=%s", 
				task.UUID, task.Name)
		}
		
		return err
	}

	log.Printf("[Worker] Task successfully deleted from database: TaskUUID=%s, TaskName=%s", 
		task.UUID, task.Name)

	// Step 4: Publish TaskDeleted event
	if w.eventPublisher != nil {
		event := events.Event{
			Type: events.TaskDeleted,
			Payload: events.TaskDeletedPayload{
				TaskUUID: task.UUID,
			},
		}
		w.eventPublisher.Publish(event)
		log.Printf("[Worker] TaskDeleted event published: TaskUUID=%s, TaskName=%s", 
			task.UUID, task.Name)
	}

	log.Printf("[Worker] Task delete process completed successfully: TaskUUID=%s, TaskName=%s", 
		task.UUID, task.Name)
	return nil
}
