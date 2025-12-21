# Module 2: Data Models

## Overview

This module defines all data models, entities, relationships, enums, and validation rules for the Cron Observer system.

## Entity Relationship Diagram

```
┌──────────┐         ┌─────────────┐
│   Task   │────────▶│  Execution  │
└──────────┘   1:N   └─────────────┘
     │                    │
     │                    │
     │                    ▼
     │              ┌──────────────┐
     │              │ ExecutionLog │
     │              └──────────────┘
     │
     │
     ▼
┌──────────┐
│ Exclusion│ (Many-to-Many with Task)
└──────────┘

┌──────────┐
│   User   │ (Future - 1:N with Task)
└──────────┘
```

## Core Entities

### 1. Task Entity

**Purpose**: Represents a scheduled task that external systems will execute.

**Table Name**: `tasks`

**Fields**:

| Field Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Internal unique identifier |
| `uuid` | VARCHAR(36) | UNIQUE, NOT NULL, INDEX | Public UUID for external systems |
| `name` | VARCHAR(255) | NOT NULL | Human-readable task name |
| `description` | TEXT | NULL | Optional task description |
| `schedule_type` | ENUM | NOT NULL | 'RECURRING' or 'ONEOFF' |
| `status` | ENUM | NOT NULL, DEFAULT 'ACTIVE' | 'ACTIVE', 'PAUSED', 'DISABLED' |
| `schedule_config` | JSON | NOT NULL | Schedule configuration (see below) |
| `metadata` | JSON | NULL | Task-specific metadata for external systems |
| `notification_config` | JSON | NULL | Notification configuration (see below) |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Creation timestamp |
| `updated_at` | TIMESTAMP | NOT NULL, DEFAULT NOW(), ON UPDATE NOW() | Last update timestamp |
| `created_by` | BIGINT | NULL, FOREIGN KEY | User ID (future - nullable for MVP) |

**Indexes**:
- Primary: `id`
- Unique: `uuid`
- Index: `status`
- Index: `schedule_type`
- Index: `created_at`

**Schedule Config Structure** (JSON):

For RECURRING tasks:
```json
{
  "cron_expression": "0 10 * * 0-4",
  "timezone": "America/New_York",
  "time_range": {
    "start": "10:00",
    "end": "14:30"
  },
  "days_of_week": [0, 1, 2, 3, 4],
  "exclusions": [1, 2, 3]
}
```

For ONEOFF tasks:
```json
{
  "execute_at": "2025-12-25T10:00:00Z",
  "timezone": "America/New_York"
}
```

**Metadata Structure** (JSON - Flexible):
```json
{
  "endpoint_url": "https://api.example.com/endpoint",
  "api_key": "secret-key",
  "custom_config": {}
}
```

**Notification Config Structure** (JSON):
```json
{
  "on_success": true,
  "on_failure": true,
  "channels": [
    {
      "type": "EMAIL",
      "endpoint": "admin@example.com",
      "template": "custom-template"
    },
    {
      "type": "WEBHOOK",
      "endpoint": "https://webhook.example.com",
      "template": null
    }
  ]
}
```

### 2. Execution Entity

**Purpose**: Tracks individual task execution attempts.

**Table Name**: `executions`

**Fields**:

| Field Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Internal unique identifier |
| `execution_uuid` | VARCHAR(36) | UNIQUE, NOT NULL, INDEX | Public UUID for SDK calls |
| `task_id` | BIGINT | NOT NULL, FOREIGN KEY | Reference to task |
| `task_uuid` | VARCHAR(36) | NOT NULL, INDEX | Task UUID (for SDK lookups) |
| `scheduled_at` | TIMESTAMP | NOT NULL | When execution was scheduled |
| `started_at` | TIMESTAMP | NULL | When execution started (status → RUNNING) |
| `completed_at` | TIMESTAMP | NULL | When execution completed (status → FINISHED/FAILED) |
| `status` | ENUM | NOT NULL, DEFAULT 'PENDING' | 'PENDING', 'RUNNING', 'FINISHED', 'FAILED', 'CANCELLED' |
| `trigger_type` | ENUM | NOT NULL | 'SCHEDULED' or 'MANUAL' |
| `duration_ms` | BIGINT | NULL | Calculated duration in milliseconds |
| `result_data` | JSON | NULL | Execution results/metadata |
| `error_message` | TEXT | NULL | Error details if failed |
| `updated_by` | VARCHAR(255) | NULL | External system identifier |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Creation timestamp |
| `updated_at` | TIMESTAMP | NOT NULL, DEFAULT NOW(), ON UPDATE NOW() | Last update timestamp |

**Indexes**:
- Primary: `id`
- Unique: `execution_uuid`
- Index: `task_id`
- Index: `task_uuid`
- Index: `status`
- Index: `scheduled_at`
- Composite: `(task_id, scheduled_at)`
- Composite: `(status, scheduled_at)`

**Status Transitions**:
- Initial: `PENDING`
- Valid transitions:
  - `PENDING` → `RUNNING`
  - `PENDING` → `CANCELLED`
  - `RUNNING` → `FINISHED`
  - `RUNNING` → `FAILED`
  - `RUNNING` → `CANCELLED`
- Terminal states: `FINISHED`, `FAILED`, `CANCELLED`

### 3. ExecutionLog Entity

**Purpose**: Stores append-only log entries for executions.

**Table Name**: `execution_logs`

**Fields**:

| Field Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Internal unique identifier |
| `execution_id` | BIGINT | NOT NULL, FOREIGN KEY | Reference to execution |
| `execution_uuid` | VARCHAR(36) | NOT NULL, INDEX | Execution UUID (for SDK lookups) |
| `timestamp` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Log entry timestamp |
| `level` | ENUM | NOT NULL | 'INFO', 'WARN', 'ERROR', 'DEBUG' |
| `message` | TEXT | NOT NULL | Log message text |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Creation timestamp |

**Indexes**:
- Primary: `id`
- Index: `execution_id`
- Index: `execution_uuid`
- Index: `timestamp`
- Composite: `(execution_id, timestamp)`

**Note**: Logs are append-only. No updates or deletes allowed.

### 4. Exclusion Entity

**Purpose**: Defines date exclusions (holidays, maintenance windows) for recurring tasks.

**Table Name**: `exclusions`

**Fields**:

| Field Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Internal unique identifier |
| `name` | VARCHAR(255) | NOT NULL | Description (e.g., "New Year's Day") |
| `type` | ENUM | NOT NULL | 'SINGLE_DATE', 'DATE_RANGE', 'RECURRING_YEARLY' |
| `start_date` | DATE | NOT NULL | Start date |
| `end_date` | DATE | NULL | End date (for DATE_RANGE) |
| `recurring_pattern` | JSON | NULL | For RECURRING_YEARLY: `{"month": 1, "day": 1}` |
| `applies_to_all_tasks` | BOOLEAN | NOT NULL, DEFAULT false | Global vs task-specific |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Creation timestamp |
| `updated_at` | TIMESTAMP | NOT NULL, DEFAULT NOW(), ON UPDATE NOW() | Last update timestamp |

**Indexes**:
- Primary: `id`
- Index: `type`
- Index: `start_date`
- Index: `applies_to_all_tasks`

### 5. TaskExclusion Junction Table

**Purpose**: Many-to-many relationship between tasks and exclusions.

**Table Name**: `task_exclusions`

**Fields**:

| Field Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `task_id` | BIGINT | NOT NULL, FOREIGN KEY | Reference to task |
| `exclusion_id` | BIGINT | NOT NULL, FOREIGN KEY | Reference to exclusion |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Creation timestamp |

**Indexes**:
- Primary: `(task_id, exclusion_id)`
- Index: `task_id`
- Index: `exclusion_id`

### 6. User Entity (Future - MVP Optional)

**Purpose**: User accounts for multi-user support.

**Table Name**: `users`

**Fields**:

| Field Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Internal unique identifier |
| `email` | VARCHAR(255) | UNIQUE, NOT NULL | User email |
| `username` | VARCHAR(100) | UNIQUE, NULL | Optional username |
| `password_hash` | VARCHAR(255) | NULL | Hashed password (future) |
| `role` | ENUM | NOT NULL, DEFAULT 'USER' | 'ADMIN' or 'USER' |
| `api_key` | VARCHAR(255) | UNIQUE, NULL, INDEX | API key for programmatic access |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Account creation timestamp |
| `updated_at` | TIMESTAMP | NOT NULL, DEFAULT NOW(), ON UPDATE NOW() | Last update timestamp |

**Indexes**:
- Primary: `id`
- Unique: `email`
- Unique: `username`
- Unique: `api_key`

### 7. NotificationLog Entity (Future)

**Purpose**: Tracks sent notifications.

**Table Name**: `notification_logs`

**Fields**:

| Field Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `id` | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Internal unique identifier |
| `execution_id` | BIGINT | NOT NULL, FOREIGN KEY | Reference to execution |
| `channel_type` | ENUM | NOT NULL | 'EMAIL', 'WEBHOOK', 'SLACK', 'DISCORD' |
| `sent_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | When notification was sent |
| `status` | ENUM | NOT NULL | 'SENT' or 'FAILED' |
| `recipient` | VARCHAR(255) | NOT NULL | Email/URL where sent |
| `error_message` | TEXT | NULL | Error details if failed |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT NOW() | Creation timestamp |

**Indexes**:
- Primary: `id`
- Index: `execution_id`
- Index: `sent_at`
- Index: `status`

## Enums

### ScheduleType
- `RECURRING`: Task repeats based on cron expression
- `ONEOFF`: Task executes once at specified datetime

### TaskStatus
- `ACTIVE`: Task is active and will be scheduled
- `PAUSED`: Task is paused (schedules continue but executions skipped)
- `DISABLED`: Task is disabled (no scheduling)

### ExecutionStatus
- `PENDING`: Execution created, waiting for external system
- `RUNNING`: External system is executing
- `FINISHED`: Execution completed successfully
- `FAILED`: Execution failed
- `CANCELLED`: Execution was cancelled

### TriggerType
- `SCHEDULED`: Execution triggered by schedule
- `MANUAL`: Execution triggered manually by user

### LogLevel
- `INFO`: Informational message
- `WARN`: Warning message
- `ERROR`: Error message
- `DEBUG`: Debug message

### ExclusionType
- `SINGLE_DATE`: Single date exclusion
- `DATE_RANGE`: Date range exclusion
- `RECURRING_YEARLY`: Recurring yearly exclusion (e.g., holidays)

### NotificationChannelType
- `EMAIL`: Email notification
- `WEBHOOK`: Webhook callback
- `SLACK`: Slack integration (future)
- `DISCORD`: Discord integration (future)

### NotificationStatus
- `SENT`: Notification sent successfully
- `FAILED`: Notification failed to send

## Validation Rules

### Task Validation
- `name`: Required, max 255 characters
- `uuid`: Auto-generated, must be unique
- `schedule_type`: Must be RECURRING or ONEOFF
- `schedule_config`: Must be valid JSON matching structure
  - RECURRING: Must have `cron_expression` and `timezone`
  - ONEOFF: Must have `execute_at` and `timezone`
- `status`: Must be ACTIVE, PAUSED, or DISABLED

### Execution Validation
- `execution_uuid`: Auto-generated, must be unique
- `task_id`: Must reference existing task
- `status`: Must follow valid status transitions
- `trigger_type`: Must be SCHEDULED or MANUAL
- `scheduled_at`: Required, must be valid timestamp

### ExecutionLog Validation
- `execution_id`: Must reference existing execution
- `level`: Must be INFO, WARN, ERROR, or DEBUG
- `message`: Required, cannot be empty
- `timestamp`: Required, must be valid timestamp

### Exclusion Validation
- `name`: Required, max 255 characters
- `type`: Must be SINGLE_DATE, DATE_RANGE, or RECURRING_YEARLY
- `start_date`: Required
- `end_date`: Required if type is DATE_RANGE
- `recurring_pattern`: Required if type is RECURRING_YEARLY, must have `month` and `day`

## Relationships

1. **Task → Execution**: One-to-Many
   - One task can have many executions
   - Cascade: Delete executions when task is deleted

2. **Execution → ExecutionLog**: One-to-Many
   - One execution can have many log entries
   - Cascade: Delete logs when execution is deleted

3. **Task → Exclusion**: Many-to-Many
   - Tasks can have multiple exclusions
   - Exclusions can apply to multiple tasks
   - Junction table: `task_exclusions`

4. **User → Task**: One-to-Many (Future)
   - One user can have many tasks
   - Cascade: Set to NULL when user is deleted (soft delete)

5. **Execution → NotificationLog**: One-to-Many (Future)
   - One execution can have multiple notification attempts

## Data Constraints

### Business Rules
1. Task UUID must be unique across all tasks
2. Execution UUID must be unique across all executions
3. Execution status transitions must be valid
4. Logs are append-only (no updates/deletes)
5. Exclusions with `applies_to_all_tasks=true` apply to all tasks
6. Scheduled executions must have valid future datetime
7. Execution duration is calculated: `completed_at - started_at`

### Database Constraints
- Foreign keys with appropriate cascade rules
- Unique constraints on UUIDs
- Check constraints on enum values
- Not null constraints on required fields
- Default values for status fields

## Next Steps

After completing this module:
1. Review and approve data models
2. Proceed to Module 3: Database Setup & Migrations
3. Create migration scripts based on these models

