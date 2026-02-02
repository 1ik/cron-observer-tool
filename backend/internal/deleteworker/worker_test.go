package deleteworker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/cron-observer/backend/internal/deletequeue"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"
)

func TestWorker_ProcessDeleteTask_TaskAlreadyDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    "test-uuid",
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	// Expectations - task already deleted (idempotent success)
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), "test-uuid").
		Return(nil, mongo.ErrNoDocuments).
		Times(1)

	// Execute
	err := worker.ProcessDeleteTask(context.Background(), msg)

	// Verify
	if err != nil {
		t.Errorf("Expected nil error (idempotent success), got: %v", err)
	}
}

func TestWorker_ProcessDeleteTask_SuccessfulDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	taskUUID := "test-uuid"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: primitive.NewObjectID(),
		Status:    models.TaskStatusPendingDelete,
	}

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    taskUUID,
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	// Expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	scheduler.EXPECT().
		UnregisterTask(taskUUID).
		Times(1)

	repo.EXPECT().
		DeleteTask(gomock.Any(), taskUUID).
		Return(nil).
		Times(1)

	// Verify event is published with correct payload
	eventPublisher.EXPECT().
		Publish(gomock.Any()).
		Do(func(event events.Event) {
			if event.Type != events.TaskDeleted {
				t.Errorf("Expected event type TaskDeleted, got %v", event.Type)
			}
			payload, ok := event.Payload.(events.TaskDeletedPayload)
			if !ok {
				t.Fatalf("Expected TaskDeletedPayload, got %T", event.Payload)
			}
			if payload.TaskUUID != taskUUID {
				t.Errorf("Expected TaskUUID %s, got %s", taskUUID, payload.TaskUUID)
			}
		}).
		Times(1)

	// Execute
	err := worker.ProcessDeleteTask(context.Background(), msg)

	// Verify
	if err != nil {
		t.Errorf("Expected nil error (success), got: %v", err)
	}
}

func TestWorker_ProcessDeleteTask_DeleteFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	taskUUID := "test-uuid"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: primitive.NewObjectID(),
		Status:    models.TaskStatusPendingDelete,
	}

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    taskUUID,
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	deleteErr := errors.New("database connection failed")

	// Expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	scheduler.EXPECT().
		UnregisterTask(taskUUID).
		Times(1)

	repo.EXPECT().
		DeleteTask(gomock.Any(), taskUUID).
		Return(deleteErr).
		Times(1)

	// UpdateTaskStatus should be called to mark task as DELETE_FAILED
	repo.EXPECT().
		UpdateTaskStatus(gomock.Any(), taskUUID, models.TaskStatusDeleteFailed).
		Return(nil).
		Times(1)

	// Event should NOT be published on failure
	eventPublisher.EXPECT().
		Publish(gomock.Any()).
		Times(0)

	// Execute
	err := worker.ProcessDeleteTask(context.Background(), msg)

	// Verify
	if err == nil {
		t.Error("Expected error when DeleteTask fails, got nil")
	}
	if err.Error() != "database connection failed" {
		t.Errorf("Expected error 'database connection failed', got: %v", err)
	}
}

func TestWorker_ProcessDeleteTask_DeleteFailure_UpdateStatusFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	taskUUID := "test-uuid"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: primitive.NewObjectID(),
		Status:    models.TaskStatusPendingDelete,
	}

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    taskUUID,
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	deleteErr := errors.New("database connection failed")
	updateErr := errors.New("failed to update status")

	// Expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	scheduler.EXPECT().
		UnregisterTask(taskUUID).
		Times(1)

	repo.EXPECT().
		DeleteTask(gomock.Any(), taskUUID).
		Return(deleteErr).
		Times(1)

	// UpdateTaskStatus fails (error is ignored in worker, but we verify it's called)
	repo.EXPECT().
		UpdateTaskStatus(gomock.Any(), taskUUID, models.TaskStatusDeleteFailed).
		Return(updateErr).
		Times(1)

	// Execute
	err := worker.ProcessDeleteTask(context.Background(), msg)

	// Verify - should still return the DeleteTask error, not UpdateTaskStatus error
	if err == nil {
		t.Error("Expected error when DeleteTask fails, got nil")
	}
	if err.Error() != "database connection failed" {
		t.Errorf("Expected error 'database connection failed', got: %v", err)
	}
}

func TestWorker_ProcessDeleteTask_GetTaskByUUIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    "test-uuid",
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	dbErr := errors.New("database error")

	// Expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), "test-uuid").
		Return(nil, dbErr).
		Times(1)

	// Scheduler and event publisher should NOT be called
	scheduler.EXPECT().
		UnregisterTask(gomock.Any()).
		Times(0)

	eventPublisher.EXPECT().
		Publish(gomock.Any()).
		Times(0)

	// Execute
	err := worker.ProcessDeleteTask(context.Background(), msg)

	// Verify
	if err == nil {
		t.Error("Expected error when GetTaskByUUID fails (non-ErrNoDocuments), got nil")
	}
	if err.Error() != "database error" {
		t.Errorf("Expected error 'database error', got: %v", err)
	}
}

func TestWorker_ProcessDeleteTask_NilScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	taskUUID := "test-uuid"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: primitive.NewObjectID(),
		Status:    models.TaskStatusPendingDelete,
	}

	repo := mocks.NewMockRepository(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, nil, eventPublisher) // nil scheduler

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    taskUUID,
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	// Expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	repo.EXPECT().
		DeleteTask(gomock.Any(), taskUUID).
		Return(nil).
		Times(1)

	eventPublisher.EXPECT().
		Publish(gomock.Any()).
		Do(func(event events.Event) {
			if event.Type != events.TaskDeleted {
				t.Errorf("Expected event type TaskDeleted, got %v", event.Type)
			}
			payload, ok := event.Payload.(events.TaskDeletedPayload)
			if !ok {
				t.Fatalf("Expected TaskDeletedPayload, got %T", event.Payload)
			}
			if payload.TaskUUID != taskUUID {
				t.Errorf("Expected TaskUUID %s, got %s", taskUUID, payload.TaskUUID)
			}
		}).
		Times(1)

	// Execute
	err := worker.ProcessDeleteTask(context.Background(), msg)

	// Verify
	if err != nil {
		t.Errorf("Expected nil error (success), got: %v", err)
	}
}

func TestWorker_ProcessDeleteTask_NilEventPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	taskUUID := "test-uuid"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: primitive.NewObjectID(),
		Status:    models.TaskStatusPendingDelete,
	}

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)

	// This should panic or handle gracefully - let's test what happens
	// Actually, looking at the code, eventPublisher is required, so this is a test
	// to ensure the code handles nil gracefully (it will panic, which is expected)
	// But let's test with a real nil to see behavior
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when eventPublisher is nil, but no panic occurred")
		}
	}()

	worker := NewWorker(repo, scheduler, nil) // nil eventPublisher

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    taskUUID,
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil)

	scheduler.EXPECT().
		UnregisterTask(taskUUID)

	repo.EXPECT().
		DeleteTask(gomock.Any(), taskUUID).
		Return(nil)

	// This should panic when trying to publish
	_ = worker.ProcessDeleteTask(context.Background(), msg)
}

func TestWorker_ProcessDeleteTask_Idempotency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	taskUUID := "test-uuid"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: primitive.NewObjectID(),
		Status:    models.TaskStatusPendingDelete,
	}

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    taskUUID,
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	// First call expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	scheduler.EXPECT().
		UnregisterTask(taskUUID).
		Times(1)

	repo.EXPECT().
		DeleteTask(gomock.Any(), taskUUID).
		Return(nil).
		Times(1)

	eventPublisher.EXPECT().
		Publish(gomock.Any()).
		Times(1)

	// Execute first time
	err1 := worker.ProcessDeleteTask(context.Background(), msg)
	if err1 != nil {
		t.Fatalf("First delete failed: %v", err1)
	}

	// Second call expectations (task already deleted - idempotent)
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(nil, mongo.ErrNoDocuments).
		Times(1)

	// Execute second time (task already deleted)
	err2 := worker.ProcessDeleteTask(context.Background(), msg)

	// Verify idempotency
	if err2 != nil {
		t.Errorf("Expected nil error on second call (idempotent), got: %v", err2)
	}
}

func TestWorker_ProcessDeleteTask_ContextCancellation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    "test-uuid",
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Expectations - GetTaskByUUID should respect context cancellation
	// Note: The actual behavior depends on repository implementation,
	// but we can test that context is passed through
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), "test-uuid").
		Do(func(ctx context.Context, uuid string) {
			if ctx.Err() == nil {
				t.Error("Expected cancelled context")
			}
		}).
		Return(nil, context.Canceled).
		Times(1)

	// Execute
	err := worker.ProcessDeleteTask(ctx, msg)

	// Verify - should return context error
	if err == nil {
		t.Error("Expected error when context is cancelled, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}

func TestWorker_ProcessDeleteTask_EventPayloadValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	taskUUID := "test-uuid-123"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: primitive.NewObjectID(),
		Status:    models.TaskStatusPendingDelete,
	}

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	msg := deletequeue.DeleteTaskMessage{
		TaskUUID:    taskUUID,
		ProjectID:   "project-123",
		RequestedAt: time.Now(),
	}

	// Expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	scheduler.EXPECT().
		UnregisterTask(taskUUID).
		Times(1)

	repo.EXPECT().
		DeleteTask(gomock.Any(), taskUUID).
		Return(nil).
		Times(1)

	// Verify event payload structure
	eventPublisher.EXPECT().
		Publish(gomock.Any()).
		Do(func(event events.Event) {
			// Verify event type
			if event.Type != events.TaskDeleted {
				t.Errorf("Expected event type %s, got %s", events.TaskDeleted, event.Type)
			}

			// Verify payload type
			payload, ok := event.Payload.(events.TaskDeletedPayload)
			if !ok {
				t.Fatalf("Expected TaskDeletedPayload, got %T", event.Payload)
			}

			// Verify TaskUUID matches
			if payload.TaskUUID != taskUUID {
				t.Errorf("Expected TaskUUID %s, got %s", taskUUID, payload.TaskUUID)
			}
		}).
		Times(1)

	// Execute
	err := worker.ProcessDeleteTask(context.Background(), msg)

	// Verify
	if err != nil {
		t.Errorf("Expected nil error (success), got: %v", err)
	}
}

func TestNewWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	scheduler := mocks.NewMockTaskUnregisterer(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, scheduler, eventPublisher)

	if worker == nil {
		t.Fatal("Expected non-nil worker, got nil")
	}

	if worker.repo != repo {
		t.Error("Repository not set correctly")
	}

	if worker.scheduler != scheduler {
		t.Error("Scheduler not set correctly")
	}

	if worker.eventPublisher != eventPublisher {
		t.Error("EventPublisher not set correctly")
	}
}

func TestNewWorker_WithNilScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	eventPublisher := mocks.NewMockEventPublisher(ctrl)

	worker := NewWorker(repo, nil, eventPublisher)

	if worker == nil {
		t.Fatal("Expected non-nil worker, got nil")
	}

	if worker.scheduler != nil {
		t.Error("Expected nil scheduler, got non-nil")
	}
}
