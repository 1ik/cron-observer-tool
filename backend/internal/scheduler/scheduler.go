package scheduler

import (
	"context"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
)

// Scheduler manages cron jobs for tasks
type Scheduler struct {
	cron     *cron.Cron
	jobs     map[string]cron.EntryID // taskUUID -> entryID
	mu       sync.RWMutex
	eventBus *events.EventBus
}

// New creates a new Scheduler instance
func New(eventBus *events.EventBus) *Scheduler {
	c := cron.New(cron.WithSeconds()) // Using WithSeconds for more precise scheduling

	return &Scheduler{
		cron:     c,
		jobs:     make(map[string]cron.EntryID),
		eventBus: eventBus,
	}
}

// Start starts the scheduler and begins listening for events
func (s *Scheduler) Start(ctx context.Context) {
	// Start the cron engine
	s.cron.Start()
	log.Println("Scheduler started")

	// Subscribe to task events
	taskCreatedCh := s.eventBus.Subscribe(events.TaskCreated)
	taskUpdatedCh := s.eventBus.Subscribe(events.TaskUpdated)
	taskDeletedCh := s.eventBus.Subscribe(events.TaskDeleted)

	// Start event listener goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Scheduler context cancelled, stopping event listener")
				return
			case event, ok := <-taskCreatedCh:
				if !ok {
					log.Println("TaskCreated channel closed")
					return
				}
				s.handleTaskCreated(event)
			case event, ok := <-taskUpdatedCh:
				if !ok {
					log.Println("TaskUpdated channel closed")
					return
				}
				s.handleTaskUpdated(event)
			case event, ok := <-taskDeletedCh:
				if !ok {
					log.Println("TaskDeleted channel closed")
					return
				}
				s.handleTaskDeleted(event)
			}
		}
	}()
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped")
}

// LoadAllActiveTasks loads all active tasks from the repository and registers them
func (s *Scheduler) LoadAllActiveTasks(ctx context.Context, repo repositories.Repository) error {
	if repo == nil {
		return nil
	}

	tasks, err := repo.GetAllActiveTasks(ctx)
	if err != nil {
		return err
	}

	log.Printf("Loading %d active tasks into scheduler", len(tasks))

	for _, task := range tasks {
		if err := s.registerTask(task); err != nil {
			log.Printf("Failed to register task %s: %v", task.UUID, err)
			continue
		}
	}

	return nil
}

// registerTask registers a task as a cron job
func (s *Scheduler) registerTask(task *models.Task) error {
	// Only register tasks with cron expressions and ACTIVE status
	if task.ScheduleConfig.CronExpression == "" {
		return nil
	}

	if task.Status != models.TaskStatusActive {
		return nil
	}

	job := &TaskJob{Task: task}
	entryID, err := s.cron.AddJob(task.ScheduleConfig.CronExpression, job)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.jobs[task.UUID] = entryID
	s.mu.Unlock()

	log.Printf("Registered cron job for task %s (UUID: %s) with expression: %s", task.Name, task.UUID, task.ScheduleConfig.CronExpression)
	return nil
}

// unregisterTask removes a task's cron job
func (s *Scheduler) unregisterTask(taskUUID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, exists := s.jobs[taskUUID]
	if !exists {
		return
	}

	s.cron.Remove(entryID)
	delete(s.jobs, taskUUID)
	log.Printf("Unregistered cron job for task UUID: %s", taskUUID)
}

// handleTaskCreated handles TaskCreated events
func (s *Scheduler) handleTaskCreated(event events.Event) {
	payload, ok := event.Payload.(events.TaskPayload)
	if !ok {
		log.Printf("Invalid payload for TaskCreated event")
		return
	}

	if err := s.registerTask(payload.Task); err != nil {
		log.Printf("Failed to register task from event: %v", err)
	}
}

// handleTaskUpdated handles TaskUpdated events
func (s *Scheduler) handleTaskUpdated(event events.Event) {
	payload, ok := event.Payload.(events.TaskPayload)
	if !ok {
		log.Printf("Invalid payload for TaskUpdated event")
		return
	}

	// Remove old job if exists
	s.unregisterTask(payload.Task.UUID)

	// Register new job (will check if task is ACTIVE and has cron expression)
	if err := s.registerTask(payload.Task); err != nil {
		log.Printf("Failed to register updated task: %v", err)
	}
}

// handleTaskDeleted handles TaskDeleted events
func (s *Scheduler) handleTaskDeleted(event events.Event) {
	payload, ok := event.Payload.(events.TaskDeletedPayload)
	if !ok {
		log.Printf("Invalid payload for TaskDeleted event")
		return
	}

	s.unregisterTask(payload.TaskUUID)
}
