package scheduler

import (
	"context"
	"log"

	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GroupStartJob represents a cron job that registers all tasks in a group
type GroupStartJob struct {
	TaskGroupID primitive.ObjectID
	Scheduler   *Scheduler
	Repo        repositories.Repository
}

// Run executes the group start job - registers all tasks in the group
func (j *GroupStartJob) Run() {
	ctx := context.Background()

	// Get the task group
	taskGroup, err := j.Repo.GetTaskGroupByID(ctx, j.TaskGroupID)
	if err != nil {
		log.Printf("[GROUP] Failed to get task group %s: %v", j.TaskGroupID.Hex(), err)
		return
	}

	// Only register if group is ACTIVE
	if taskGroup.Status != models.TaskGroupStatusActive {
		log.Printf("[GROUP] Task group %s is not ACTIVE, skipping registration", taskGroup.UUID)
		return
	}

	// Update group state to RUNNING (status remains ACTIVE)
	if err := j.Repo.UpdateTaskGroupState(ctx, taskGroup.UUID, models.TaskGroupStateRunning); err != nil {
		log.Printf("[GROUP] Failed to update group %s state to RUNNING: %v", taskGroup.UUID, err)
	}

	// Get all tasks in this group
	tasks, err := j.Repo.GetTasksByGroupID(ctx, j.TaskGroupID)
	if err != nil {
		log.Printf("[GROUP] Failed to get tasks for group %s: %v", taskGroup.UUID, err)
		return
	}

	log.Printf("[GROUP] Registering %d tasks for group %s (start time: %s)", len(tasks), taskGroup.UUID, taskGroup.StartTime)

	// Register each task and update state to RUNNING (status remains ACTIVE)
	for _, task := range tasks {
		if err := j.Scheduler.registerTask(ctx, task); err != nil {
			log.Printf("[GROUP] Failed to register task %s in group %s: %v", task.UUID, taskGroup.UUID, err)
			continue
		}
		// Update task state to RUNNING (status remains ACTIVE)
		if err := j.Repo.UpdateTaskState(ctx, task.UUID, models.TaskStateRunning); err != nil {
			log.Printf("[GROUP] Failed to update task %s state to RUNNING: %v", task.UUID, err)
		}
	}
}

// GroupEndJob represents a cron job that unregisters all tasks in a group
type GroupEndJob struct {
	TaskGroupID primitive.ObjectID
	Scheduler   *Scheduler
	Repo        repositories.Repository
}

// Run executes the group end job - unregisters all tasks in the group
func (j *GroupEndJob) Run() {
	ctx := context.Background()

	// Get the task group
	taskGroup, err := j.Repo.GetTaskGroupByID(ctx, j.TaskGroupID)
	if err != nil {
		log.Printf("[GROUP] Failed to get task group %s: %v", j.TaskGroupID.Hex(), err)
		return
	}

	// Update group state to NOT_RUNNING (status remains ACTIVE)
	if err := j.Repo.UpdateTaskGroupState(ctx, taskGroup.UUID, models.TaskGroupStateNotRunning); err != nil {
		log.Printf("[GROUP] Failed to update group %s state to NOT_RUNNING: %v", taskGroup.UUID, err)
	}

	// Get all tasks in this group
	tasks, err := j.Repo.GetTasksByGroupID(ctx, j.TaskGroupID)
	if err != nil {
		log.Printf("[GROUP] Failed to get tasks for group %s: %v", taskGroup.UUID, err)
		return
	}

	log.Printf("[GROUP] Unregistering %d tasks for group %s (end time: %s)", len(tasks), taskGroup.UUID, taskGroup.EndTime)

	// Unregister each task and update state to NOT_RUNNING (status remains ACTIVE)
	for _, task := range tasks {
		j.Scheduler.unregisterTask(task.UUID)
		// Update task state to NOT_RUNNING (status remains ACTIVE)
		if err := j.Repo.UpdateTaskState(ctx, task.UUID, models.TaskStateNotRunning); err != nil {
			log.Printf("[GROUP] Failed to update task %s state to NOT_RUNNING: %v", task.UUID, err)
		}
	}
}
