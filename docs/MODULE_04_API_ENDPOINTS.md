# Module 4: Core API Endpoints

## Overview

This module defines all REST API endpoints for task management, execution history, and manual operations. These endpoints are used by the frontend UI and for programmatic access.

## API Design Principles

- **RESTful**: Follow REST conventions
- **Versioned**: All endpoints under `/api/v1/`
- **JSON**: Request/response in JSON format
- **Error Handling**: Consistent error response format
- **Authentication**: API key or token-based (future)
- **Pagination**: For list endpoints
- **Filtering**: Query parameters for filtering

## Base URL

```
http://localhost:8080/api/v1
```

## Common Response Formats

### Success Response
```json
{
  "data": { ... },
  "message": "Success message (optional)"
}
```

### Error Response
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": { ... }
  }
}
```

### Paginated Response
```json
{
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

## HTTP Status Codes

- `200 OK`: Successful GET, PUT, PATCH
- `201 Created`: Successful POST (resource created)
- `204 No Content`: Successful DELETE
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., duplicate UUID)
- `422 Unprocessable Entity`: Validation errors
- `500 Internal Server Error`: Server error

## Task Management Endpoints

### 1. Create Task

**Endpoint**: `POST /api/v1/tasks`

**Request Body**:
```json
{
  "name": "DSE News scrape",
  "description": "Scrape news from DSE website",
  "schedule_type": "RECURRING",
  "schedule_config": {
    "cron_expression": "0 10 * * 0-4",
    "timezone": "America/New_York",
    "time_range": {
      "start": "10:00",
      "end": "14:30"
    },
    "days_of_week": [0, 1, 2, 3, 4],
    "exclusions": []
  },
  "metadata": {
    "endpoint_url": "https://api.example.com/scrape",
    "api_key": "secret-key"
  },
  "notification_config": {
    "on_success": false,
    "on_failure": true,
    "channels": []
  }
}
```

**Response** (201 Created):
```json
{
  "data": {
    "id": 1,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "name": "DSE News scrape",
    "description": "Scrape news from DSE website",
    "schedule_type": "RECURRING",
    "status": "ACTIVE",
    "schedule_config": { ... },
    "metadata": { ... },
    "notification_config": { ... },
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

**Validation Rules**:
- `name`: Required, max 255 characters
- `schedule_type`: Required, must be "RECURRING" or "ONEOFF"
- `schedule_config`: Required, must match structure
- UUID is auto-generated

### 2. Get All Tasks

**Endpoint**: `GET /api/v1/tasks`

**Query Parameters**:
- `status`: Filter by status (ACTIVE, PAUSED, DISABLED)
- `schedule_type`: Filter by type (RECURRING, ONEOFF)
- `search`: Search by name or description
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 20, max: 100)
- `sort_by`: Sort field (name, created_at, updated_at)
- `sort_order`: Sort order (asc, desc)

**Example**: `GET /api/v1/tasks?status=ACTIVE&page=1&page_size=20`

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "name": "DSE News scrape",
      "schedule_type": "RECURRING",
      "status": "ACTIVE",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 50,
    "total_pages": 3
  }
}
```

### 3. Get Task by UUID

**Endpoint**: `GET /api/v1/tasks/{task_uuid}`

**Response** (200 OK):
```json
{
  "data": {
    "id": 1,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "name": "DSE News scrape",
    "description": "Scrape news from DSE website",
    "schedule_type": "RECURRING",
    "status": "ACTIVE",
    "schedule_config": { ... },
    "metadata": { ... },
    "notification_config": { ... },
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z",
    "statistics": {
      "total_executions": 150,
      "success_rate": 0.95,
      "avg_duration_ms": 3000,
      "last_execution": "2025-01-15T10:01:10Z",
      "next_execution": "2025-01-16T10:00:00Z"
    }
  }
}
```

**Error** (404 Not Found):
```json
{
  "error": {
    "code": "TASK_NOT_FOUND",
    "message": "Task with UUID '550e8400-e29b-41d4-a716-446655440000' not found"
  }
}
```

### 4. Update Task

**Endpoint**: `PUT /api/v1/tasks/{task_uuid}`

**Request Body**: Same as Create Task (all fields required)

**Response** (200 OK): Same as Get Task by UUID

**Notes**:
- UUID cannot be changed
- Updating schedule_config may affect future executions
- Status can be updated (ACTIVE, PAUSED, DISABLED)

### 5. Partial Update Task

**Endpoint**: `PATCH /api/v1/tasks/{task_uuid}`

**Request Body**: Only fields to update
```json
{
  "status": "PAUSED",
  "metadata": {
    "api_key": "new-key"
  }
}
```

**Response** (200 OK): Updated task object

### 6. Delete Task

**Endpoint**: `DELETE /api/v1/tasks/{task_uuid}`

**Response** (204 No Content)

**Notes**:
- Cascades to delete all executions and logs
- Cannot delete if there are RUNNING executions (return 409 Conflict)

### 7. Pause Task

**Endpoint**: `POST /api/v1/tasks/{task_uuid}/pause`

**Response** (200 OK):
```json
{
  "data": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "status": "PAUSED",
    "message": "Task paused successfully"
  }
}
```

### 8. Resume Task

**Endpoint**: `POST /api/v1/tasks/{task_uuid}/resume`

**Response** (200 OK):
```json
{
  "data": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "status": "ACTIVE",
    "message": "Task resumed successfully"
  }
}
```

### 9. Manual Trigger Task

**Endpoint**: `POST /api/v1/tasks/{task_uuid}/trigger`

**Response** (201 Created):
```json
{
  "data": {
    "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
    "task_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "status": "PENDING",
    "trigger_type": "MANUAL",
    "scheduled_at": "2025-01-15T12:00:00Z",
    "message": "Execution created successfully"
  }
}
```

**Notes**:
- Creates execution record with PENDING status
- External system can pick up and execute
- Bypasses schedule (immediate execution)

## Execution History Endpoints

### 10. Get Executions for Task

**Endpoint**: `GET /api/v1/tasks/{task_uuid}/executions`

**Query Parameters**:
- `status`: Filter by status (PENDING, RUNNING, FINISHED, FAILED, CANCELLED)
- `trigger_type`: Filter by type (SCHEDULED, MANUAL)
- `start_date`: Filter from date (ISO 8601)
- `end_date`: Filter to date (ISO 8601)
- `page`: Page number
- `page_size`: Items per page
- `sort_by`: scheduled_at, created_at
- `sort_order`: asc, desc

**Example**: `GET /api/v1/tasks/{uuid}/executions?status=FINISHED&start_date=2025-01-01&page=1`

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
      "task_uuid": "550e8400-e29b-41d4-a716-446655440000",
      "scheduled_at": "2025-01-15T10:00:00Z",
      "started_at": "2025-01-15T10:00:10Z",
      "completed_at": "2025-01-15T10:00:13Z",
      "status": "FINISHED",
      "trigger_type": "SCHEDULED",
      "duration_ms": 3000,
      "result_data": null,
      "error_message": null
    }
  ],
  "pagination": { ... }
}
```

### 11. Get Execution by UUID

**Endpoint**: `GET /api/v1/executions/{execution_uuid}`

**Response** (200 OK):
```json
{
  "data": {
    "id": 1,
    "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
    "task_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "task": {
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "name": "DSE News scrape"
    },
    "scheduled_at": "2025-01-15T10:00:00Z",
    "started_at": "2025-01-15T10:00:10Z",
    "completed_at": "2025-01-15T10:00:13Z",
    "status": "FINISHED",
    "trigger_type": "SCHEDULED",
    "duration_ms": 3000,
    "result_data": {
      "records_processed": 150,
      "status": "success"
    },
    "error_message": null,
    "updated_by": "external-system-1",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:13Z",
    "logs_count": 5
  }
}
```

### 12. Get Execution Logs

**Endpoint**: `GET /api/v1/executions/{execution_uuid}/logs`

**Query Parameters**:
- `level`: Filter by level (INFO, WARN, ERROR, DEBUG)
- `start_timestamp`: Filter from timestamp
- `end_timestamp`: Filter to timestamp
- `limit`: Max number of logs (default: 1000, max: 10000)

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
      "timestamp": "2025-01-15T10:00:10Z",
      "level": "INFO",
      "message": "Task started"
    },
    {
      "id": 2,
      "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
      "timestamp": "2025-01-15T10:00:10Z",
      "level": "INFO",
      "message": "Fetch started for DSE News"
    }
  ],
  "total": 5
}
```

### 13. Cancel Execution

**Endpoint**: `POST /api/v1/executions/{execution_uuid}/cancel`

**Response** (200 OK):
```json
{
  "data": {
    "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
    "status": "CANCELLED",
    "message": "Execution cancelled successfully"
  }
}
```

**Notes**:
- Only PENDING or RUNNING executions can be cancelled
- Returns 400 if execution is already finished

## Date-Based Navigation Endpoints

### 14. Get Execution Dates

**Endpoint**: `GET /api/v1/executions/dates`

**Query Parameters**:
- `start_date`: Start date (ISO 8601, default: 30 days ago)
- `end_date`: End date (ISO 8601, default: today)
- `task_uuid`: Filter by task (optional)

**Response** (200 OK):
```json
{
  "data": [
    {
      "date": "2025-01-15",
      "execution_count": 5,
      "status_summary": {
        "finished": 4,
        "failed": 1,
        "pending": 0,
        "running": 0
      }
    },
    {
      "date": "2025-01-14",
      "execution_count": 3,
      "status_summary": {
        "finished": 3,
        "failed": 0,
        "pending": 0,
        "running": 0
      }
    }
  ]
}
```

### 15. Get Executions by Date

**Endpoint**: `GET /api/v1/executions/by-date/{date}`

**Path Parameter**: `date` (ISO 8601 date: YYYY-MM-DD)

**Query Parameters**:
- `task_uuid`: Filter by task (optional)
- `status`: Filter by status (optional)

**Response** (200 OK):
```json
{
  "data": [
    {
      "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
      "task": {
        "uuid": "550e8400-e29b-41d4-a716-446655440000",
        "name": "DSE News scrape"
      },
      "scheduled_at": "2025-01-15T10:00:00Z",
      "status": "FINISHED",
      "duration_ms": 3000,
      "summary": "Task executed successfully"
    }
  ],
  "date": "2025-01-15",
  "total": 5
}
```

## Exclusion Management Endpoints

### 16. Create Exclusion

**Endpoint**: `POST /api/v1/exclusions`

**Request Body**:
```json
{
  "name": "New Year's Day",
  "type": "RECURRING_YEARLY",
  "start_date": "2025-01-01",
  "recurring_pattern": {
    "month": 1,
    "day": 1
  },
  "applies_to_all_tasks": true
}
```

**Response** (201 Created): Exclusion object

### 17. Get All Exclusions

**Endpoint**: `GET /api/v1/exclusions`

**Response** (200 OK): List of exclusions

### 18. Delete Exclusion

**Endpoint**: `DELETE /api/v1/exclusions/{exclusion_id}`

**Response** (204 No Content)

## Error Codes

| Code | Description |
|------|-------------|
| `TASK_NOT_FOUND` | Task with given UUID not found |
| `EXECUTION_NOT_FOUND` | Execution with given UUID not found |
| `INVALID_SCHEDULE_CONFIG` | Schedule configuration is invalid |
| `INVALID_STATUS_TRANSITION` | Status transition is not allowed |
| `TASK_HAS_RUNNING_EXECUTIONS` | Cannot delete task with running executions |
| `VALIDATION_ERROR` | Request validation failed |
| `INTERNAL_ERROR` | Internal server error |

## Authentication (Future)

Currently, authentication is deferred. For MVP, use:
- API key in header: `X-API-Key: your-api-key`
- Or simple token: `Authorization: Bearer token`

Future implementation will support:
- JWT tokens
- OAuth2
- API key management

## Rate Limiting (Future)

- 100 requests per minute per API key
- 1000 requests per hour per API key
- Burst: 20 requests per second

## Next Steps

After completing this module:
1. Implement endpoint handlers
2. Add request validation
3. Add error handling
4. Add authentication middleware (basic for MVP)
5. Proceed to Module 5: Scheduler Engine

