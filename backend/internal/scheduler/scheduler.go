package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
)

// Scheduler manages cron jobs for tasks
type Scheduler struct {
	cron      *cron.Cron
	jobs      map[string]cron.EntryID            // taskUUID -> entryID
	groupJobs map[string]map[string]cron.EntryID // groupUUID -> {"start": entryID, "end": entryID}
	mu        sync.RWMutex
	eventBus  *events.EventBus
	repo      repositories.Repository
}

// New creates a new Scheduler instance
func New(eventBus *events.EventBus, repo repositories.Repository) *Scheduler {
	c := cron.New(cron.WithSeconds()) // Using WithSeconds for more precise scheduling

	return &Scheduler{
		cron:      c,
		jobs:      make(map[string]cron.EntryID),
		groupJobs: make(map[string]map[string]cron.EntryID),
		eventBus:  eventBus,
		repo:      repo,
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

	// Subscribe to task group events
	taskGroupCreatedCh := s.eventBus.Subscribe(events.TaskGroupCreated)
	taskGroupUpdatedCh := s.eventBus.Subscribe(events.TaskGroupUpdated)
	taskGroupDeletedCh := s.eventBus.Subscribe(events.TaskGroupDeleted)

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
			case event, ok := <-taskGroupCreatedCh:
				if !ok {
					log.Println("TaskGroupCreated channel closed")
					return
				}
				s.handleTaskGroupCreated(event)
			case event, ok := <-taskGroupUpdatedCh:
				if !ok {
					log.Println("TaskGroupUpdated channel closed")
					return
				}
				s.handleTaskGroupUpdated(event)
			case event, ok := <-taskGroupDeletedCh:
				if !ok {
					log.Println("TaskGroupDeleted channel closed")
					return
				}
				s.handleTaskGroupDeleted(event)
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
func (s *Scheduler) LoadAllActiveTasks(ctx context.Context) error {
	// Load active task groups with windows
	taskGroups, err := s.repo.GetActiveTaskGroupsWithWindows(ctx)
	if err != nil {
		log.Printf("Failed to load active task groups: %v", err)
	} else {
		log.Printf("Loading %d active task groups with windows", len(taskGroups))
		for _, group := range taskGroups {
			if err := s.registerGroupWindowJobs(group); err != nil {
				log.Printf("Failed to register window jobs for group %s: %v", group.UUID, err)
			}
		}
	}

	tasks, err := s.repo.GetAllActiveTasks(ctx)
	if err != nil {
		return err
	}

	log.Printf("Loading %d active tasks into scheduler", len(tasks))

	for _, task := range tasks {
		if err := s.registerTask(ctx, task); err != nil {
			log.Printf("Failed to register task %s: %v", task.UUID, err)
			continue
		}
	}

	return nil
}

// registerTask registers a task as a cron job
func (s *Scheduler) registerTask(ctx context.Context, task *models.Task) error {
	// Only register tasks with cron expressions
	if task.ScheduleConfig.CronExpression == "" {
		return nil
	}

	// If task belongs to a group, check group status and window
	if task.TaskGroupID != nil {
		taskGroup, err := s.repo.GetTaskGroupByID(ctx, *task.TaskGroupID)
		if err != nil {
			log.Printf("Failed to get task group for task %s: %v", task.UUID, err)
			return nil // Don't register if group lookup fails
		}

		// Only register if group is ACTIVE and current time is within window
		if taskGroup.Status != models.TaskGroupStatusActive {
			return nil // Group is not active
		}

		// Check if current time is within group window
		if !s.isWithinGroupWindow(ctx, taskGroup) {
			return nil // Not within window
		}
	} else {
		// For tasks without groups, check individual status
		if task.Status != models.TaskStatusActive {
			return nil
		}
	}

	job := &TaskJob{Task: task, Repo: s.repo}
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

	ctx := context.Background()
	if err := s.registerTask(ctx, payload.Task); err != nil {
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
	ctx := context.Background()
	if err := s.registerTask(ctx, payload.Task); err != nil {
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

// handleTaskGroupCreated handles TaskGroupCreated events
func (s *Scheduler) handleTaskGroupCreated(event events.Event) {
	payload, ok := event.Payload.(events.TaskGroupPayload)
	if !ok {
		log.Printf("Invalid payload for TaskGroupCreated event")
		return
	}

	// Only register window jobs if group has start and end times
	if payload.TaskGroup.StartTime != "" && payload.TaskGroup.EndTime != "" {
		if err := s.registerGroupWindowJobs(payload.TaskGroup); err != nil {
			log.Printf("Failed to register group window jobs: %v", err)
		}
	}
}

// handleTaskGroupUpdated handles TaskGroupUpdated events
func (s *Scheduler) handleTaskGroupUpdated(event events.Event) {
	payload, ok := event.Payload.(events.TaskGroupPayload)
	if !ok {
		log.Printf("Invalid payload for TaskGroupUpdated event")
		return
	}

	// Remove old jobs
	s.unregisterGroupWindowJobs(payload.TaskGroup.UUID)

	// Register new jobs if group has windows
	if payload.TaskGroup.StartTime != "" && payload.TaskGroup.EndTime != "" {
		if err := s.registerGroupWindowJobs(payload.TaskGroup); err != nil {
			log.Printf("Failed to register updated group window jobs: %v", err)
		}
	}
}

// handleTaskGroupDeleted handles TaskGroupDeleted events
func (s *Scheduler) handleTaskGroupDeleted(event events.Event) {
	payload, ok := event.Payload.(events.TaskGroupDeletedPayload)
	if !ok {
		log.Printf("Invalid payload for TaskGroupDeleted event")
		return
	}

	s.unregisterGroupWindowJobs(payload.TaskGroupUUID)
}

// registerGroupWindowJobs registers cron jobs for group start and end times
// Creates two daily cron jobs: one at start time (registers all tasks) and one at end time (unregisters all tasks)
func (s *Scheduler) registerGroupWindowJobs(taskGroup *models.TaskGroup) error {
	if taskGroup.StartTime == "" || taskGroup.EndTime == "" {
		return nil // No window defined
	}

	// Convert start time to cron expression
	startCron, err := timeToCronExpression(taskGroup.StartTime, taskGroup.Timezone)
	if err != nil {
		return err
	}

	// Convert end time to cron expression
	endCron, err := timeToCronExpression(taskGroup.EndTime, taskGroup.Timezone)
	if err != nil {
		return err
	}

	// Create start job
	startJob := &GroupStartJob{
		TaskGroupID: taskGroup.ID,
		Scheduler:   s,
		Repo:        s.repo,
	}
	startEntryID, err := s.cron.AddJob(startCron, startJob)
	if err != nil {
		return err
	}

	// Create end job
	endJob := &GroupEndJob{
		TaskGroupID: taskGroup.ID,
		Scheduler:   s,
		Repo:        s.repo,
	}
	endEntryID, err := s.cron.AddJob(endCron, endJob)
	if err != nil {
		// Remove start job if end job fails
		s.cron.Remove(startEntryID)
		return err
	}

	// Store both entry IDs
	s.mu.Lock()
	if s.groupJobs[taskGroup.UUID] == nil {
		s.groupJobs[taskGroup.UUID] = make(map[string]cron.EntryID)
	}
	s.groupJobs[taskGroup.UUID]["start"] = startEntryID
	s.groupJobs[taskGroup.UUID]["end"] = endEntryID
	s.mu.Unlock()

	log.Printf("Registered window jobs for group %s: start=%s, end=%s", taskGroup.UUID, startCron, endCron)
	return nil
}

// unregisterGroupWindowJobs removes cron jobs for a group's window
func (s *Scheduler) unregisterGroupWindowJobs(groupUUID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobs, exists := s.groupJobs[groupUUID]
	if !exists {
		return
	}

	if startID, ok := jobs["start"]; ok {
		s.cron.Remove(startID)
	}
	if endID, ok := jobs["end"]; ok {
		s.cron.Remove(endID)
	}

	delete(s.groupJobs, groupUUID)
	log.Printf("Unregistered window jobs for group UUID: %s", groupUUID)
}

// isWithinGroupWindow checks if current time is within the group's time window
func (s *Scheduler) isWithinGroupWindow(ctx context.Context, taskGroup *models.TaskGroup) bool {
	if taskGroup.StartTime == "" || taskGroup.EndTime == "" {
		return true // No window defined, always within
	}

	// Parse times and check current time
	now := time.Now()

	// Load location for timezone
	loc, err := time.LoadLocation(taskGroup.Timezone)
	if err != nil {
		log.Printf("Invalid timezone %s for group %s: %v", taskGroup.Timezone, taskGroup.UUID, err)
		return false
	}

	// Parse start and end times
	startTime, err := parseTimeInLocation(taskGroup.StartTime, loc, now)
	if err != nil {
		log.Printf("Failed to parse start time %s: %v", taskGroup.StartTime, err)
		return false
	}

	endTime, err := parseTimeInLocation(taskGroup.EndTime, loc, now)
	if err != nil {
		log.Printf("Failed to parse end time %s: %v", taskGroup.EndTime, err)
		return false
	}

	// Check if current time is within window
	nowInLoc := now.In(loc)
	currentTime := time.Date(nowInLoc.Year(), nowInLoc.Month(), nowInLoc.Day(), nowInLoc.Hour(), nowInLoc.Minute(), 0, 0, loc)

	return (currentTime.Equal(startTime) || currentTime.After(startTime)) && currentTime.Before(endTime)
}

// timeToCronExpression converts HH:MM time to daily cron expression
// Assumes time is in the given timezone, converts to UTC for cron
func timeToCronExpression(timeStr, timezone string) (string, error) {
	// Parse time (HH:MM format)
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	// Parse the time string
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return "", err
	}

	// Create a time for today in the group's timezone
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, loc)

	// Convert to UTC
	utcTime := today.UTC()

	// Create cron expression: second minute hour day month weekday
	// Format: "second minute hour * * *"
	return fmt.Sprintf("%d %d %d * * *", utcTime.Second(), utcTime.Minute(), utcTime.Hour()), nil
}

// parseTimeInLocation parses HH:MM time string in the given location for today
func parseTimeInLocation(timeStr string, loc *time.Location, reference time.Time) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	refInLoc := reference.In(loc)
	return time.Date(refInLoc.Year(), refInLoc.Month(), refInLoc.Day(), t.Hour(), t.Minute(), 0, 0, loc), nil
}

// StartGroup manually registers all tasks in a group
func (s *Scheduler) StartGroup(ctx context.Context, groupUUID string) error {
	taskGroup, err := s.repo.GetTaskGroupByUUID(ctx, groupUUID)
	if err != nil {
		return err
	}

	if taskGroup.Status != models.TaskGroupStatusActive {
		return fmt.Errorf("task group is not ACTIVE")
	}

	tasks, err := s.repo.GetTasksByGroupID(ctx, taskGroup.ID)
	if err != nil {
		return err
	}

	log.Printf("[GROUP] Manually starting group %s, registering %d tasks", groupUUID, len(tasks))

	for _, task := range tasks {
		if err := s.registerTask(ctx, task); err != nil {
			log.Printf("[GROUP] Failed to register task %s: %v", task.UUID, err)
			continue
		}
	}

	return nil
}

// StopGroup manually unregisters all tasks in a group
func (s *Scheduler) StopGroup(ctx context.Context, groupUUID string) error {
	taskGroup, err := s.repo.GetTaskGroupByUUID(ctx, groupUUID)
	if err != nil {
		return err
	}

	tasks, err := s.repo.GetTasksByGroupID(ctx, taskGroup.ID)
	if err != nil {
		return err
	}

	log.Printf("[GROUP] Manually stopping group %s, unregistering %d tasks", groupUUID, len(tasks))

	for _, task := range tasks {
		s.unregisterTask(task.UUID)
	}

	return nil
}
