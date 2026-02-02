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
	task, err := w.repo.GetTaskByUUID(ctx, msg.TaskUUID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Already deleted – safe to ack.
			return nil
		}
		return err
	}

	// 1) Stop cron (idempotent: safe to call many times).
	if w.scheduler != nil {
		w.scheduler.UnregisterTask(task.UUID)
	}

	// 2) Hard delete from Mongo.
	if err := w.repo.DeleteTask(ctx, task.UUID); err != nil {
		_ = w.repo.UpdateTaskStatus(ctx, task.UUID, models.TaskStatusDeleteFailed)
		log.Printf("[deleteworker] DeleteTask failed for %s: %v", task.UUID, err)
		return err
	}

	// 3) Publish TaskDeleted event (scheduler may unregister again; no-op).
	w.eventPublisher.Publish(events.Event{
		Type: events.TaskDeleted,
		Payload: events.TaskDeletedPayload{
			TaskUUID: task.UUID,
		},
	})

	return nil
}
