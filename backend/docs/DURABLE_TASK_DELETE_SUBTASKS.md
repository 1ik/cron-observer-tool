# Durable Task Delete - Subtasks Breakdown

This document breaks down the durable task delete implementation into small, sequential tasks that can be executed and reviewed one by one.

---

## Task 1: Extend TaskStatus Enum with Internal Delete States [Done]

**Goal**: Add `PENDING_DELETE` and `DELETE_FAILED` status values to the TaskStatus enum.

**Files to Review/Modify**:
- `backend/internal/models/task.go` - Find the `TaskStatus` type definition

**What to do**:
1. Locate the `TaskStatus` type (likely a string type or const enum)
2. Add two new constants:
   - `TaskStatusPendingDelete = "PENDING_DELETE"`
   - `TaskStatusDeleteFailed = "DELETE_FAILED"`
3. Ensure these are documented as internal-only states

**What to verify**:
- [ ] Both new statuses are defined
- [ ] They follow the existing naming convention
- [ ] They are clearly marked as internal-only (not for external API use)

**Dependencies**: None (first task)

---

## Task 2: Review and Update Status Validation Logic [Done]

**Goal**: Ensure external APIs reject `PENDING_DELETE` and `DELETE_FAILED` statuses from clients.

**Files to Review**:
- `backend/internal/handlers/tasks.go` - Look for `CreateTask`, `UpdateTask`, `UpdateTaskStatus` handlers
- `backend/internal/validators/custom.go` - Check for status validation logic
- `backend/internal/utils/validation.go` - Check for validation utilities

**What to do**:
1. Find where task status is validated in API handlers
2. Ensure validation only allows `ACTIVE` and `DISABLED` from external clients
3. Document that `PENDING_DELETE` and `DELETE_FAILED` are backend-only

**What to verify**:
- [ ] External APIs reject `PENDING_DELETE` and `DELETE_FAILED` if sent by clients
- [ ] Validation logic is clear and well-documented
- [ ] Scheduler code already ignores non-`ACTIVE` tasks (verify this assumption)

**Dependencies**: Task 1

---

## Task 3: Create Delete Queue Message Structure

**Goal**: Define the message structure for delete job messages.

**Files to Create/Modify**:
- `backend/internal/deletequeue/message.go` (new file)

**What to do**:
1. Create the new file `backend/internal/deletequeue/message.go`
2. Define `DeleteTaskMessage` struct with fields:
   - `TaskUUID string`
   - `ProjectID string`
   - `RequestedAt time.Time`
   - `RequestID string` (optional, omitempty)
3. Add JSON tags for serialization
4. Add any necessary validation or helper methods

**What to verify**:
- [ ] Struct is properly defined with JSON tags
- [ ] Fields match the specification in the main doc
- [ ] File follows Go naming conventions

**Dependencies**: None

---

## Task 4: Create Delete Queue Interface Abstractions

**Goal**: Define broker-agnostic interfaces for publishing and consuming delete jobs.

**Files to Create/Modify**:
- `backend/internal/deletequeue/queue.go` (new file)

**What to do**:
1. Create the new file `backend/internal/deletequeue/queue.go`
2. Define `DeleteJobPublisher` interface with method:
   - `PublishDeleteTask(ctx context.Context, msg DeleteTaskMessage) error`
3. Define `DeleteJobConsumer` interface with method:
   - `Start(ctx context.Context, handler func(context.Context, DeleteTaskMessage) error) error`
4. Add documentation explaining these are broker-agnostic abstractions

**What to verify**:
- [ ] Interfaces are properly defined
- [ ] Method signatures match the specification
- [ ] Documentation explains the abstraction purpose

**Dependencies**: Task 3

---

## Task 5: Review Repository Methods for Delete Operations

**Goal**: Ensure repository has methods needed for delete workflow.

**Files to Review**:
- `backend/internal/repositories/repository.go` - Check interface definition
- `backend/internal/repositories/mongo.go` - Check implementation
- `backend/internal/database/collections.go` - Check database operations

**What to do**:
1. Verify `GetTaskByUUID(ctx context.Context, uuid string) (*Task, error)` exists
2. Verify `UpdateTaskStatus(ctx context.Context, taskUUID string, status TaskStatus) error` exists
3. Verify `DeleteTask(ctx context.Context, taskUUID string) error` exists (hard delete)
4. If any are missing, add them to both interface and implementation

**What to verify**:
- [ ] All three methods exist in repository interface
- [ ] All three methods are implemented in MongoDB repository
- [ ] `DeleteTask` performs hard delete (not soft delete)
- [ ] Methods handle `mongo.ErrNoDocuments` appropriately

**Dependencies**: Task 1

---

## Task 6: Update TaskHandler.DeleteTask API Handler

**Goal**: Modify the delete endpoint to schedule deletion instead of executing it immediately.

**Files to Modify**:
- `backend/internal/handlers/tasks.go` - Find `DeleteTask` handler method

**What to do**:
1. Locate the `DeleteTask` handler method
2. Update it to:
   - Parse `project_id` and `task_uuid` from request
   - Call `repo.GetTaskByUUID` - if not found, return `202 Accepted` with `{"status": "ALREADY_DELETED"}`
   - Call `repo.UpdateTaskStatus(ctx, taskUUID, TaskStatusPendingDelete)`
   - Publish `DeleteTaskMessage` via `DeleteJobPublisher` (will be wired later)
   - Return `202 Accepted` with JSON body:
     ```json
     {
       "status": "PENDING_DELETE",
       "task_uuid": "<uuid>",
       "message": "Task deletion has been scheduled"
     }
     ```
3. Handle broker publish failure (start with rollback approach - set status back and return 500)

**What to verify**:
- [ ] Handler returns `202 Accepted` instead of `204 No Content`
- [ ] Status is set to `PENDING_DELETE` before publishing
- [ ] Idempotent behavior when task already deleted
- [ ] Error handling for broker publish failures

**Dependencies**: Tasks 1, 3, 4, 5

---

## Task 7: Review Scheduler UnregisterTask Method

**Goal**: Verify scheduler has `UnregisterTask` method and understand its behavior.

**Files to Review**:
- `backend/internal/scheduler/scheduler.go` - Check scheduler interface
- `backend/internal/scheduler/job.go` - Check job management
- `backend/internal/scheduler/group_job.go` - Check group job handling

**What to do**:
1. Find `UnregisterTask(uuid string)` method in scheduler
2. Verify it's idempotent (safe to call multiple times)
3. Understand how it stops cron jobs
4. Document its behavior if not already documented

**What to verify**:
- [ ] `UnregisterTask` method exists
- [ ] Method is idempotent (no errors if task already unregistered)
- [ ] Method properly stops cron execution
- [ ] Method signature matches what worker will need

**Dependencies**: None (can be done in parallel with other tasks)

---

## Task 8: Review EventBus for TaskDeleted Event

**Goal**: Verify EventBus can publish `TaskDeleted` events.

**Files to Review**:
- `backend/internal/events/event.go` - Check event type definitions
- `backend/internal/events/bus.go` - Check EventBus interface

**What to do**:
1. Check if `TaskDeleted` event type exists
2. Check if `TaskDeletedPayload` struct exists with `TaskUUID` field
3. Verify `EventBus.Publish` method signature
4. If missing, add the event type and payload struct

**What to verify**:
- [ ] `TaskDeleted` event type is defined
- [ ] `TaskDeletedPayload` struct exists with `TaskUUID string`
- [ ] EventBus can publish events
- [ ] Scheduler subscribes to `TaskDeleted` events (verify this)

**Dependencies**: None (can be done in parallel)

---

## Task 9: Implement Delete Worker Core Logic

**Goal**: Create the worker that processes delete job messages.

**Files to Create**:
- `backend/internal/deleteworker/worker.go` (new file)

**What to do**:
1. Create new file `backend/internal/deleteworker/worker.go`
2. Define `Worker` struct with dependencies:
   - `repo repositories.Repository`
   - `scheduler` (interface with `UnregisterTask` method)
   - `eventBus *events.EventBus`
   - `consumer DeleteJobConsumer` (will be wired later)
3. Implement `processDeleteTask(ctx context.Context, msg DeleteTaskMessage) error`:
   - Get task via `repo.GetTaskByUUID`
   - If `mongo.ErrNoDocuments`, return nil (already deleted, idempotent)
   - Call `scheduler.UnregisterTask(task.UUID)` (idempotent)
   - Call `repo.DeleteTask(ctx, task.UUID)` for hard delete
   - If delete fails, optionally set status to `DELETE_FAILED` and return error
   - Publish `TaskDeleted` event via `eventBus.Publish`
   - Return nil to ack the message
4. Add proper error handling and logging

**What to verify**:
- [ ] Worker struct has all required dependencies
- [ ] `processDeleteTask` is idempotent
- [ ] Error handling properly distinguishes retryable vs non-retryable errors
- [ ] Hard delete is performed (not soft delete)
- [ ] Event is published after successful delete

**Dependencies**: Tasks 1, 3, 4, 5, 7, 8

---

## Task 10: Implement Delete Job Consumer Start Method

**Goal**: Implement the consumer that receives messages and calls the worker.

**Files to Create/Modify**:
- `backend/internal/deletequeue/consumer.go` (new file, or add to existing queue.go)

**What to do**:
1. Create consumer implementation (will depend on chosen broker - RabbitMQ, SQS, Redis, etc.)
2. Implement `Start` method that:
   - Subscribes to delete job queue
   - For each message, calls `processDeleteTask`
   - Only acks message if `processDeleteTask` returns nil
   - Lets broker handle retries/DLQ for errors
3. Add proper context handling and graceful shutdown

**What to verify**:
- [ ] Consumer properly subscribes to queue
- [ ] Messages are only acked on success
- [ ] Errors trigger broker retry mechanism
- [ ] Context cancellation is handled gracefully

**Dependencies**: Tasks 4, 9

**Note**: This task requires choosing a message broker (RabbitMQ, SQS, Redis Queue, etc.)

---

## Task 11: Implement Delete Job Publisher

**Goal**: Implement the publisher that sends delete job messages to the broker.

**Files to Create/Modify**:
- `backend/internal/deletequeue/publisher.go` (new file, or add to existing queue.go)

**What to do**:
1. Create publisher implementation for chosen broker
2. Implement `PublishDeleteTask` method that:
   - Serializes `DeleteTaskMessage` to JSON
   - Publishes to delete job queue
   - Returns error if publish fails
3. Add proper error handling and logging

**What to verify**:
- [ ] Publisher properly serializes messages
- [ ] Messages are published to correct queue
- [ ] Error handling is appropriate
- [ ] Logging helps with debugging

**Dependencies**: Tasks 3, 4

**Note**: This task requires choosing a message broker (RabbitMQ, SQS, Redis Queue, etc.)

---

## Task 12: Wire DeleteJobPublisher into TaskHandler

**Goal**: Connect the publisher to the delete handler.

**Files to Modify**:
- `backend/internal/handlers/tasks.go` - Update `TaskHandler` struct and `DeleteTask` method
- `backend/cmd/server/main.go` - Wire dependencies

**What to do**:
1. Add `deletePublisher DeleteJobPublisher` field to `TaskHandler` struct
2. Update `TaskHandler` constructor to accept publisher
3. Update `DeleteTask` handler to use `deletePublisher.PublishDeleteTask`
4. In `main.go`, create publisher instance and pass to handler

**What to verify**:
- [ ] Publisher is properly injected into handler
- [ ] Handler uses publisher to enqueue delete jobs
- [ ] Dependency injection is clean and testable

**Dependencies**: Tasks 6, 11

---

## Task 13: Implement Delete Reconciler

**Goal**: Create periodic reconciler to handle stuck delete tasks.

**Files to Create**:
- `backend/internal/reconciler/delete_reconciler.go` (new file)

**What to do**:
1. Create new file `backend/internal/reconciler/delete_reconciler.go`
2. Define `DeleteReconciler` struct with:
   - `repo repositories.Repository`
   - `publisher DeleteJobPublisher`
   - `ticker *time.Ticker` (e.g., 5 minutes)
   - `threshold time.Duration` (e.g., 10 minutes for updated_at)
3. Implement `Start(ctx context.Context)` method that:
   - Runs on ticker interval
   - Queries tasks with `status IN (PENDING_DELETE, DELETE_FAILED)`
   - For each task, check if `updated_at` is older than threshold
   - Re-publish `DeleteTaskMessage` for eligible tasks
   - Log actions for observability
4. Implement `Stop()` method for graceful shutdown

**What to verify**:
- [ ] Reconciler queries correct statuses
- [ ] Threshold prevents constant re-queues
- [ ] Re-publishing works correctly
- [ ] Graceful shutdown is handled

**Dependencies**: Tasks 3, 4, 5, 11

---

## Task 14: Wire Delete Worker and Consumer in Main

**Goal**: Start the delete worker and consumer in the main application.

**Files to Modify**:
- `backend/cmd/server/main.go` - Add worker startup

**What to do**:
1. Create delete worker instance with all dependencies
2. Create delete consumer instance
3. Start consumer in a goroutine with proper context handling
4. Ensure graceful shutdown stops consumer
5. Add configuration for broker connection (queue name, connection string, etc.)

**What to verify**:
- [ ] Worker is created with all dependencies
- [ ] Consumer starts and processes messages
- [ ] Graceful shutdown works correctly
- [ ] Configuration is externalized (env vars, config file)

**Dependencies**: Tasks 9, 10, 12

---

## Task 15: Wire Delete Reconciler in Main

**Goal**: Start the delete reconciler in the main application.

**Files to Modify**:
- `backend/cmd/server/main.go` - Add reconciler startup

**What to do**:
1. Create delete reconciler instance
2. Start reconciler in a goroutine with proper context handling
3. Ensure graceful shutdown stops reconciler
4. Add configuration for reconciler interval and threshold

**What to verify**:
- [ ] Reconciler starts correctly
- [ ] Runs on configured interval
- [ ] Graceful shutdown works correctly
- [ ] Configuration is externalized

**Dependencies**: Task 13, 14

---

## Task 16: Add Configuration for Message Broker

**Goal**: Add configuration options for message broker connection and settings.

**Files to Review/Modify**:
- `backend/internal/config/config.go` - Add broker config struct
- `backend/internal/config/loader.go` - Load broker config from env/file
- `.env.example` - Add example broker configuration

**What to do**:
1. Add broker configuration struct (connection string, queue name, etc.)
2. Load configuration from environment variables or config file
3. Add example values to `.env.example`
4. Document configuration options

**What to verify**:
- [ ] Configuration struct includes all necessary fields
- [ ] Configuration loads from environment
- [ ] Example values are provided
- [ ] Documentation is clear

**Dependencies**: None (can be done early)

---

## Task 17: Update API Response Types (if needed)

**Goal**: Ensure API response types include new delete response structure.

**Files to Review**:
- `backend/internal/models/error.go` - Check response models
- `backend/api-docs/openapi.json` - Update OpenAPI spec if needed

**What to do**:
1. Check if response types need updating for `202 Accepted` with JSON body
2. Update OpenAPI/Swagger documentation if needed
3. Ensure response structure matches specification

**What to verify**:
- [ ] Response types support new delete response format
- [ ] OpenAPI spec is updated
- [ ] Response structure matches specification

**Dependencies**: Task 6

---

## Task 18: Add Unit Tests for Delete Worker

**Goal**: Write unit tests for delete worker logic.

**Files to Create**:
- `backend/internal/deleteworker/worker_test.go` (new file)

**What to do**:
1. Create test file
2. Test `processDeleteTask` with various scenarios:
   - Task already deleted (idempotent success)
   - Successful delete flow
   - Delete failure (sets DELETE_FAILED)
   - Scheduler unregister failure
   - Event publish failure
3. Use mocks for dependencies (repo, scheduler, eventBus)

**What to verify**:
- [ ] All scenarios are covered
- [ ] Idempotency is tested
- [ ] Error cases are tested
- [ ] Mocks are properly set up

**Dependencies**: Task 9

---

## Task 19: Add Integration Tests for Delete Flow

**Goal**: Write integration tests for the full delete workflow.

**Files to Create**:
- `backend/internal/handlers/tasks_test.go` - Add delete handler tests
- Or create separate integration test file

**What to do**:
1. Test `DeleteTask` handler:
   - Returns 202 with PENDING_DELETE status
   - Handles already-deleted task (idempotent)
   - Handles broker publish failure
2. Test worker processes delete job:
   - Stops cron
   - Deletes from database
   - Publishes event
3. Test full flow end-to-end (if possible)

**What to verify**:
- [ ] Handler tests pass
- [ ] Worker tests pass
- [ ] Integration tests cover happy path and error cases

**Dependencies**: Tasks 6, 9, 12

---

## Task 20: Review Frontend Delete Implementation

**Goal**: Update frontend to handle async delete contract.

**Files to Review**:
- Frontend task deletion code (likely in UI components)
- API client code (likely in `UI/packages/lib/src/api.ts`)

**What to do**:
1. Review current delete implementation
2. Update to handle `202 Accepted` response
3. Show toast notification "Task deletion scheduled"
4. Decide on Option A (optimistic removal) or Option B (show PENDING_DELETE badge)
5. Handle `ALREADY_DELETED` status gracefully
6. Ensure re-delete is safe (idempotent)

**What to verify**:
- [ ] Frontend handles 202 response correctly
- [ ] User feedback is clear
- [ ] Re-delete works correctly
- [ ] UI updates appropriately

**Dependencies**: Task 6

---

## Summary Checklist

Use this checklist to track overall progress:

- [ ] **Phase 1: Foundation** (Tasks 1-5)
  - [ ] Task 1: Extend TaskStatus enum
  - [ ] Task 2: Update status validation
  - [ ] Task 3: Create message structure
  - [ ] Task 4: Create queue interfaces
  - [ ] Task 5: Review repository methods

- [ ] **Phase 2: Core Implementation** (Tasks 6-11)
  - [ ] Task 6: Update DeleteTask handler
  - [ ] Task 7: Review scheduler UnregisterTask
  - [ ] Task 8: Review EventBus
  - [ ] Task 9: Implement delete worker
  - [ ] Task 10: Implement consumer
  - [ ] Task 11: Implement publisher

- [ ] **Phase 3: Integration** (Tasks 12-15)
  - [ ] Task 12: Wire publisher to handler
  - [ ] Task 13: Implement reconciler
  - [ ] Task 14: Wire worker in main
  - [ ] Task 15: Wire reconciler in main

- [ ] **Phase 4: Configuration & Testing** (Tasks 16-19)
  - [ ] Task 16: Add broker configuration
  - [ ] Task 17: Update API response types
  - [ ] Task 18: Add unit tests
  - [ ] Task 19: Add integration tests

- [ ] **Phase 5: Frontend** (Task 20)
  - [ ] Task 20: Update frontend delete implementation

---

## Notes

- **Message Broker Choice**: Tasks 10 and 11 require choosing a message broker (RabbitMQ, AWS SQS, Redis Queue, etc.). This decision should be made early.
- **Testing Strategy**: Consider testing each component in isolation before integration testing.
- **Error Handling**: Pay special attention to error handling and idempotency throughout.
- **Observability**: Add logging and metrics at each step for debugging and monitoring.
