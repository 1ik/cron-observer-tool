package reconciler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yourusername/cron-observer/backend/internal/deletequeue"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
)

// DeleteReconciler periodically re-enqueues stuck PENDING_DELETE and DELETE_FAILED tasks.
type DeleteReconciler struct {
	repo      repositories.Repository
	publisher deletequeue.DeleteJobPublisher
	ticker    *time.Ticker
	interval  time.Duration
	threshold time.Duration
	mu        sync.RWMutex
	running   bool
	stopCh    chan struct{}
}

// NewDeleteReconciler creates a new delete reconciler.
// interval: how often to run (e.g., 5 minutes)
// threshold: only re-enqueue tasks older than this (e.g., 10 minutes)
func NewDeleteReconciler(repo repositories.Repository, publisher deletequeue.DeleteJobPublisher, interval, threshold time.Duration) *DeleteReconciler {
	return &DeleteReconciler{
		repo:      repo,
		publisher: publisher,
		ticker:    time.NewTicker(interval),
		interval:  interval,
		threshold: threshold,
		stopCh:    make(chan struct{}),
	}
}

// Start begins the reconciler loop. Runs until ctx is cancelled or Stop() is called.
func (r *DeleteReconciler) Start(ctx context.Context) error {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return ErrReconcilerAlreadyRunning
	}
	r.running = true
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		r.running = false
		r.ticker.Stop()
		r.mu.Unlock()
	}()

	log.Printf("[reconciler] Delete reconciler started (interval=%v, threshold=%v)", r.interval, r.threshold)

	// Run immediately on start
	r.reconcile(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[reconciler] Delete reconciler context cancelled, stopping")
			return ctx.Err()
		case <-r.stopCh:
			log.Printf("[reconciler] Delete reconciler stopped")
			return nil
		case <-r.ticker.C:
			r.reconcile(ctx)
		}
	}
}

// Stop stops the reconciler gracefully.
func (r *DeleteReconciler) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.running {
		close(r.stopCh)
	}
}

// reconcile queries stuck tasks and re-enqueues them.
func (r *DeleteReconciler) reconcile(ctx context.Context) {
	// Query tasks with PENDING_DELETE or DELETE_FAILED status
	statuses := []models.TaskStatus{
		models.TaskStatusPendingDelete,
		models.TaskStatusDeleteFailed,
	}

	tasks, err := r.repo.GetTasksByStatus(ctx, statuses)
	if err != nil {
		log.Printf("[reconciler] Failed to query stuck delete tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		return // No stuck tasks
	}

	now := time.Now()
	reEnqueuedCount := 0

	for _, task := range tasks {
		// Only re-enqueue if updated_at is older than threshold
		age := now.Sub(task.UpdatedAt)
		if age < r.threshold {
			continue // Task is too recent, skip
		}

		// Re-publish delete job
		msg := deletequeue.DeleteTaskMessage{
			TaskUUID:    task.UUID,
			ProjectID:   task.ProjectID.Hex(),
			RequestedAt: time.Now(),
		}

		if err := r.publisher.PublishDeleteTask(ctx, msg); err != nil {
			log.Printf("[reconciler] Failed to re-enqueue delete job for task %s: %v", task.UUID, err)
			continue
		}

		reEnqueuedCount++
		log.Printf("[reconciler] Re-enqueued delete job for task %s (status=%s, age=%v)", task.UUID, task.Status, age)
	}

	if reEnqueuedCount > 0 {
		log.Printf("[reconciler] Re-enqueued %d stuck delete task(s)", reEnqueuedCount)
	}
}

// Errors
var (
	ErrReconcilerAlreadyRunning = &ReconcilerError{Message: "reconciler is already running"}
)

// ReconcilerError represents a reconciler error.
type ReconcilerError struct {
	Message string
}

func (e *ReconcilerError) Error() string {
	return e.Message
}
