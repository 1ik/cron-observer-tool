package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cron-observer/backend/internal/deletequeue"
	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"
)

// mockScheduler implements the scheduler interface needed by TaskHandler
type mockScheduler struct {
	unregisterTaskCalled bool
	taskUUID             string
}

func (m *mockScheduler) RegisterTask(ctx context.Context, task *models.Task) error {
	return nil
}

func (m *mockScheduler) UnregisterTask(taskUUID string) {
	m.unregisterTaskCalled = true
	m.taskUUID = taskUUID
}

func (m *mockScheduler) IsWithinGroupWindow(ctx context.Context, taskGroup *models.TaskGroup) bool {
	return false
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestTaskHandler_DeleteTask_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	projectID := primitive.NewObjectID()
	taskUUID := "test-task-uuid"
	taskName := "test-task"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: projectID,
		Name:      taskName,
		Status:    models.TaskStatusActive,
	}

	repo := mocks.NewMockRepository(ctrl)
	eventBus := events.NewEventBus(100)
	defer eventBus.Close()
	scheduler := &mockScheduler{}
	deletePublisher := mocks.NewMockDeleteJobPublisher(ctrl)

	handler := NewTaskHandler(repo, eventBus, scheduler, []string{}, deletePublisher)

	// Expectations
	// Handler calls GetTaskByUUID once to fetch task
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	// Handler publishes to RabbitMQ
	deletePublisher.EXPECT().
		PublishDeleteTask(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, msg deletequeue.DeleteTaskMessage) error {
			if msg.TaskUUID != taskUUID {
				t.Errorf("Expected TaskUUID %s, got %s", taskUUID, msg.TaskUUID)
			}
			if msg.ProjectID != projectID.Hex() {
				t.Errorf("Expected ProjectID %s, got %s", projectID.Hex(), msg.ProjectID)
			}
			return nil
		}).
		Times(1)

	// Setup router
	router := setupRouter()
	router.DELETE("/api/v1/projects/:project_id/tasks/:task_uuid", handler.DeleteTask)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.Hex()+"/tasks/"+taskUUID, nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, got %d", http.StatusAccepted, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "PENDING_DELETE" {
		t.Errorf("Expected status 'PENDING_DELETE', got '%v'", response["status"])
	}

	if response["task_uuid"] != taskUUID {
		t.Errorf("Expected task_uuid '%s', got '%v'", taskUUID, response["task_uuid"])
	}

	// Scheduler should NOT be called in handler (happens in worker)
	if scheduler.unregisterTaskCalled {
		t.Error("UnregisterTask should not be called in handler (happens in worker)")
	}
}

func TestTaskHandler_DeleteTask_TaskAlreadyDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	projectID := primitive.NewObjectID()
	taskUUID := "test-task-uuid"

	repo := mocks.NewMockRepository(ctrl)
	eventBus := events.NewEventBus(100)
	defer eventBus.Close()
	scheduler := &mockScheduler{}
	deletePublisher := mocks.NewMockDeleteJobPublisher(ctrl)

	handler := NewTaskHandler(repo, eventBus, scheduler, []string{}, deletePublisher)

	// Expectations - task already deleted (idempotent)
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(nil, mongo.ErrNoDocuments).
		Times(1)

	// Setup router
	router := setupRouter()
	router.DELETE("/api/v1/projects/:project_id/tasks/:task_uuid", handler.DeleteTask)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.Hex()+"/tasks/"+taskUUID, nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "ALREADY_DELETED" {
		t.Errorf("Expected status 'ALREADY_DELETED', got '%v'", response["status"])
	}

	// Publisher should NOT be called for already deleted tasks (no expectation needed since it returns early)
}

func TestTaskHandler_DeleteTask_MissingProjectID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	eventBus := events.NewEventBus(100)
	defer eventBus.Close()
	scheduler := &mockScheduler{}
	deletePublisher := mocks.NewMockDeleteJobPublisher(ctrl)

	handler := NewTaskHandler(repo, eventBus, scheduler, []string{}, deletePublisher)

	// Setup router
	router := setupRouter()
	router.DELETE("/api/v1/projects/:project_id/tasks/:task_uuid", handler.DeleteTask)

	// Create request with empty project_id
	req, _ := http.NewRequest("DELETE", "/api/v1/projects//tasks/test-uuid", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTaskHandler_DeleteTask_MissingTaskUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := primitive.NewObjectID()

	repo := mocks.NewMockRepository(ctrl)
	eventBus := events.NewEventBus(100)
	defer eventBus.Close()
	scheduler := &mockScheduler{}
	deletePublisher := mocks.NewMockDeleteJobPublisher(ctrl)

	handler := NewTaskHandler(repo, eventBus, scheduler, []string{}, deletePublisher)

	// Test by calling the handler directly with empty task_uuid param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.Hex()+"/tasks/", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "project_id", Value: projectID.Hex()},
		{Key: "task_uuid", Value: ""}, // Empty task_uuid
	}

	handler.DeleteTask(c)

	// Verify
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTaskHandler_DeleteTask_GetTaskByUUIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	projectID := primitive.NewObjectID()
	taskUUID := "test-task-uuid"
	dbErr := errors.New("database connection failed")

	repo := mocks.NewMockRepository(ctrl)
	eventBus := events.NewEventBus(100)
	defer eventBus.Close()
	scheduler := &mockScheduler{}
	deletePublisher := mocks.NewMockDeleteJobPublisher(ctrl)

	handler := NewTaskHandler(repo, eventBus, scheduler, []string{}, deletePublisher)

	// Expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(nil, dbErr).
		Times(1)

	// Setup router
	router := setupRouter()
	router.DELETE("/api/v1/projects/:project_id/tasks/:task_uuid", handler.DeleteTask)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.Hex()+"/tasks/"+taskUUID, nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == nil {
		t.Error("Expected error message in response")
	}
}

func TestTaskHandler_DeleteTask_PublishFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	projectID := primitive.NewObjectID()
	taskUUID := "test-task-uuid"
	taskName := "test-task"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: projectID,
		Name:      taskName,
		Status:    models.TaskStatusActive,
	}
	publishErr := errors.New("failed to publish to RabbitMQ")

	repo := mocks.NewMockRepository(ctrl)
	eventBus := events.NewEventBus(100)
	defer eventBus.Close()
	scheduler := &mockScheduler{}
	deletePublisher := mocks.NewMockDeleteJobPublisher(ctrl)

	handler := NewTaskHandler(repo, eventBus, scheduler, []string{}, deletePublisher)

	// Expectations
	// Handler calls GetTaskByUUID once to fetch task
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	// Publisher fails to publish
	deletePublisher.EXPECT().
		PublishDeleteTask(gomock.Any(), gomock.Any()).
		Return(publishErr).
		Times(1)

	// Setup router
	router := setupRouter()
	router.DELETE("/api/v1/projects/:project_id/tasks/:task_uuid", handler.DeleteTask)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.Hex()+"/tasks/"+taskUUID, nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == nil {
		t.Error("Expected error message in response")
	}
}

func TestTaskHandler_DeleteTask_NoPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	projectID := primitive.NewObjectID()
	taskUUID := "test-task-uuid"
	taskName := "test-task"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: projectID,
		Name:      taskName,
		Status:    models.TaskStatusActive,
	}

	repo := mocks.NewMockRepository(ctrl)
	eventBus := events.NewEventBus(100)
	defer eventBus.Close()
	scheduler := &mockScheduler{}

	// Handler with nil publisher (RabbitMQ not configured)
	handler := NewTaskHandler(repo, eventBus, scheduler, []string{}, nil)

	// Expectations
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	// Setup router
	router := setupRouter()
	router.DELETE("/api/v1/projects/:project_id/tasks/:task_uuid", handler.DeleteTask)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.Hex()+"/tasks/"+taskUUID, nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == nil {
		t.Error("Expected error message in response")
	}
}

// Note: Event publishing is now handled asynchronously by the worker,
// so we can't test it in the handler tests. Event tests should be in worker tests.

func TestTaskHandler_DeleteTask_NilScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	projectID := primitive.NewObjectID()
	taskUUID := "test-task-uuid"
	taskName := "test-task"
	task := &models.Task{
		ID:        primitive.NewObjectID(),
		UUID:      taskUUID,
		ProjectID: projectID,
		Name:      taskName,
		Status:    models.TaskStatusActive,
	}

	repo := mocks.NewMockRepository(ctrl)
	eventBus := events.NewEventBus(100)
	defer eventBus.Close()
	deletePublisher := mocks.NewMockDeleteJobPublisher(ctrl)

	// Create handler with nil scheduler (scheduler is optional)
	handler := NewTaskHandler(repo, eventBus, nil, []string{}, deletePublisher)

	// Expectations
	// Handler calls GetTaskByUUID once to fetch task
	repo.EXPECT().
		GetTaskByUUID(gomock.Any(), taskUUID).
		Return(task, nil).
		Times(1)

	// Handler publishes to RabbitMQ
	deletePublisher.EXPECT().
		PublishDeleteTask(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	// Setup router
	router := setupRouter()
	router.DELETE("/api/v1/projects/:project_id/tasks/:task_uuid", handler.DeleteTask)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.Hex()+"/tasks/"+taskUUID, nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, got %d", http.StatusAccepted, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "PENDING_DELETE" {
		t.Errorf("Expected status 'PENDING_DELETE', got '%v'", response["status"])
	}
}
