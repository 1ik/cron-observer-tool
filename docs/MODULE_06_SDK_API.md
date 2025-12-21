# Module 6: SDK/API for External Systems

## Overview

This module implements the SDK/API endpoints that allow external systems to discover pending executions, update execution status, and append logs. This is the interface between Cron Observer and external task executors.

## SDK Architecture

### Communication Flow

```
External System                    Cron Observer
     │                                  │
     │  1. Discover Pending            │
     ├─────────────────────────────────>│
     │                                  │
     │  2. Get Execution Details        │
     ├─────────────────────────────────>│
     │                                  │
     │  3. Update Status: RUNNING       │
     ├─────────────────────────────────>│
     │                                  │
     │  4. Append Logs (during exec)  │
     ├─────────────────────────────────>│
     │                                  │
     │  5. Update Status: FINISHED     │
     ├─────────────────────────────────>│
     │                                  │
```

## Authentication

### API Key Authentication

For MVP, use simple API key authentication:

**Header**: `X-API-Key: your-api-key-here`

**Future**: JWT tokens, OAuth2

### API Key Management

- API keys stored in database (users table)
- Each API key associated with a user (future)
- Rate limiting per API key
- Key rotation support (future)

## SDK Endpoints

### Base URL

```
http://localhost:8080/api/v1/sdk
```

## 1. Discover Pending Executions

### Endpoint

**GET** `/api/v1/sdk/tasks/{task_uuid}/executions/pending`

**Purpose**: Allow external systems to find pending executions for a specific task.

**Headers**:
```
X-API-Key: your-api-key
```

**Path Parameters**:
- `task_uuid`: UUID of the task

**Query Parameters**:
- `limit`: Maximum number of executions to return (default: 10, max: 100)

**Response** (200 OK):
```json
{
  "data": [
    {
      "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
      "task_uuid": "550e8400-e29b-41d4-a716-446655440000",
      "scheduled_at": "2025-01-15T10:00:00Z",
      "status": "PENDING",
      "trigger_type": "SCHEDULED",
      "task_metadata": {
        "endpoint_url": "https://api.example.com/scrape",
        "api_key": "secret-key"
      }
    }
  ],
  "count": 1
}
```

**Response** (404 Not Found):
```json
{
  "error": {
    "code": "TASK_NOT_FOUND",
    "message": "Task with UUID '550e8400-e29b-41d4-a716-446655440000' not found"
  }
}
```

**Implementation Notes**:
- Only returns executions with status PENDING
- Ordered by scheduled_at (oldest first)
- Includes task metadata for execution context
- External system should claim execution by updating status to RUNNING

### Alternative: Get All Pending Executions

**GET** `/api/v1/sdk/executions/pending`

**Purpose**: Get all pending executions across all tasks (for systems managing multiple tasks).

**Query Parameters**:
- `limit`: Maximum number of executions (default: 50, max: 200)
- `task_uuid`: Optional filter by task

**Response**: Same format as above, but may include multiple tasks.

## 2. Get Execution Details

### Endpoint

**GET** `/api/v1/sdk/executions/{execution_uuid}`

**Purpose**: Get full execution details including task metadata.

**Response** (200 OK):
```json
{
  "data": {
    "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
    "task_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "task": {
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "name": "DSE News scrape",
      "metadata": {
        "endpoint_url": "https://api.example.com/scrape",
        "api_key": "secret-key"
      }
    },
    "scheduled_at": "2025-01-15T10:00:00Z",
    "status": "PENDING",
    "trigger_type": "SCHEDULED"
  }
}
```

## 3. Update Execution Status

### Endpoint

**PUT** `/api/v1/sdk/executions/{execution_uuid}/status`

**Purpose**: Update execution status (PENDING → RUNNING → FINISHED/FAILED).

**Request Body**:
```json
{
  "status": "RUNNING",
  "result_data": null
}
```

**Status Values**:
- `RUNNING`: Execution has started
- `FINISHED`: Execution completed successfully
- `FAILED`: Execution failed

**Response** (200 OK):
```json
{
  "data": {
    "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
    "status": "RUNNING",
    "started_at": "2025-01-15T10:00:10Z",
    "message": "Status updated successfully"
  }
}
```

**Status Transition Rules**:

| From Status | To Status | Allowed | Notes |
|-------------|-----------|---------|-------|
| PENDING | RUNNING | ✅ | Sets started_at |
| PENDING | CANCELLED | ✅ | Manual cancellation |
| RUNNING | FINISHED | ✅ | Sets completed_at, calculates duration |
| RUNNING | FAILED | ✅ | Sets completed_at, requires error_message |
| RUNNING | CANCELLED | ✅ | Manual cancellation |
| FINISHED | - | ❌ | Terminal state |
| FAILED | - | ❌ | Terminal state |
| CANCELLED | - | ❌ | Terminal state |

**Validation**:
- Invalid status transitions return 400 Bad Request
- FINISHED status requires result_data (optional but recommended)
- FAILED status should include error_message in separate call or result_data

**Implementation**:

```java
@PutMapping("/executions/{execution_uuid}/status")
public ResponseEntity<ExecutionStatusResponse> updateStatus(
        @PathVariable String execution_uuid,
        @RequestBody StatusUpdateRequest request,
        @RequestHeader("X-API-Key") String apiKey) {
    
    // Authenticate
    validateApiKey(apiKey);
    
    // Get execution
    Execution execution = executionRepository.findByExecutionUuid(execution_uuid)
        .orElseThrow(() -> new ExecutionNotFoundException(execution_uuid));
    
    // Validate status transition
    ExecutionStatus newStatus = ExecutionStatus.valueOf(request.getStatus());
    if (!isValidTransition(execution.getStatus(), newStatus)) {
        throw new InvalidStatusTransitionException(
            execution.getStatus(), 
            newStatus
        );
    }
    
    // Update status
    execution.setStatus(newStatus);
    
    // Set timestamps
    if (newStatus == ExecutionStatus.RUNNING && execution.getStartedAt() == null) {
        execution.setStartedAt(Instant.now());
    }
    
    if (newStatus == ExecutionStatus.FINISHED || 
        newStatus == ExecutionStatus.FAILED) {
        execution.setCompletedAt(Instant.now());
        
        // Calculate duration
        if (execution.getStartedAt() != null) {
            long duration = Duration.between(
                execution.getStartedAt(), 
                execution.getCompletedAt()
            ).toMillis();
            execution.setDurationMs(duration);
        }
    }
    
    // Set result data
    if (request.getResultData() != null) {
        execution.setResultData(request.getResultData());
    }
    
    // Track who updated
    execution.setUpdatedBy(getSystemIdentifier(apiKey));
    execution.setUpdatedAt(Instant.now());
    
    execution = executionRepository.save(execution);
    
    return ResponseEntity.ok(new ExecutionStatusResponse(execution));
}
```

## 4. Append Logs

### Endpoint

**POST** `/api/v1/sdk/executions/{execution_uuid}/logs`

**Purpose**: Append log entries to an execution (append-only operation).

**Request Body**:
```json
{
  "logs": [
    {
      "timestamp": "2025-01-15T10:00:10Z",
      "level": "INFO",
      "message": "Task started"
    },
    {
      "timestamp": "2025-01-15T10:00:10Z",
      "level": "INFO",
      "message": "Fetch started for DSE News"
    },
    {
      "timestamp": "2025-01-15T10:00:12Z",
      "level": "WARN",
      "message": "Rate limit approaching"
    }
  ]
}
```

**Log Entry Fields**:
- `timestamp`: ISO 8601 timestamp (optional, defaults to server time)
- `level`: INFO, WARN, ERROR, DEBUG
- `message`: Log message text (required)

**Response** (200 OK):
```json
{
  "data": {
    "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
    "logs_added": 3,
    "total_logs": 5,
    "message": "Logs appended successfully"
  }
}
```

**Validation**:
- Execution must exist
- Execution status must be PENDING or RUNNING (cannot append to finished executions)
- Logs array cannot be empty
- Each log entry must have valid level and non-empty message
- Timestamp is optional (defaults to current time)

**Implementation**:

```java
@PostMapping("/executions/{execution_uuid}/logs")
public ResponseEntity<LogAppendResponse> appendLogs(
        @PathVariable String execution_uuid,
        @RequestBody LogAppendRequest request,
        @RequestHeader("X-API-Key") String apiKey) {
    
    // Authenticate
    validateApiKey(apiKey);
    
    // Get execution
    Execution execution = executionRepository.findByExecutionUuid(execution_uuid)
        .orElseThrow(() -> new ExecutionNotFoundException(execution_uuid));
    
    // Validate execution status
    if (execution.getStatus() == ExecutionStatus.FINISHED ||
        execution.getStatus() == ExecutionStatus.FAILED ||
        execution.getStatus() == ExecutionStatus.CANCELLED) {
        throw new ExecutionNotWritableException(
            "Cannot append logs to finished execution"
        );
    }
    
    // Create log entries
    List<ExecutionLog> logs = new ArrayList<>();
    Instant now = Instant.now();
    
    for (LogEntryRequest logRequest : request.getLogs()) {
        ExecutionLog log = new ExecutionLog();
        log.setExecutionId(execution.getId());
        log.setExecutionUuid(execution_uuid);
        log.setLevel(LogLevel.valueOf(logRequest.getLevel()));
        log.setMessage(logRequest.getMessage());
        
        // Use provided timestamp or current time
        if (logRequest.getTimestamp() != null) {
            log.setTimestamp(Instant.parse(logRequest.getTimestamp()));
        } else {
            log.setTimestamp(now);
        }
        
        logs.add(log);
    }
    
    // Save logs (batch insert for performance)
    executionLogRepository.saveAll(logs);
    
    // Get total log count
    long totalLogs = executionLogRepository.countByExecutionId(execution.getId());
    
    return ResponseEntity.ok(new LogAppendResponse(
        execution_uuid,
        logs.size(),
        totalLogs
    ));
}
```

## 5. Update Execution with Error

### Endpoint

**PUT** `/api/v1/sdk/executions/{execution_uuid}/error`

**Purpose**: Mark execution as failed with error message.

**Request Body**:
```json
{
  "error_message": "Connection timeout after 30 seconds",
  "result_data": {
    "attempts": 3,
    "last_error": "Connection timeout"
  }
}
```

**Response** (200 OK):
```json
{
  "data": {
    "execution_uuid": "660e8400-e29b-41d4-a716-446655440001",
    "status": "FAILED",
    "error_message": "Connection timeout after 30 seconds",
    "completed_at": "2025-01-15T10:00:40Z",
    "message": "Execution marked as failed"
  }
}
```

**Implementation**:
- Updates status to FAILED
- Sets error_message
- Sets completed_at
- Calculates duration
- Can be called instead of status update with FAILED

## Error Handling

### Error Codes

| Code | Description |
|------|-------------|
| `EXECUTION_NOT_FOUND` | Execution with UUID not found |
| `INVALID_STATUS_TRANSITION` | Status transition not allowed |
| `EXECUTION_NOT_WRITABLE` | Cannot modify finished execution |
| `INVALID_LOG_LEVEL` | Invalid log level |
| `INVALID_API_KEY` | API key is invalid or missing |
| `RATE_LIMIT_EXCEEDED` | Too many requests |

### Error Response Format

```json
{
  "error": {
    "code": "INVALID_STATUS_TRANSITION",
    "message": "Cannot transition from FINISHED to RUNNING",
    "details": {
      "current_status": "FINISHED",
      "requested_status": "RUNNING",
      "allowed_transitions": ["N/A - terminal state"]
    }
  }
}
```

## Rate Limiting

### Rate Limits

- **Per API Key**: 100 requests per minute
- **Per Execution**: 10 status updates per minute
- **Log Appending**: 50 log append requests per minute per execution

### Rate Limit Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642248000
```

### Rate Limit Response

**429 Too Many Requests**:
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Try again in 5 seconds.",
    "retry_after": 5
  }
}
```

## SDK Client Examples

### Python SDK Example

```python
import requests

class CronObserverSDK:
    def __init__(self, base_url, api_key):
        self.base_url = base_url
        self.headers = {"X-API-Key": api_key}
    
    def get_pending_executions(self, task_uuid, limit=10):
        url = f"{self.base_url}/api/v1/sdk/tasks/{task_uuid}/executions/pending"
        response = requests.get(url, headers=self.headers, params={"limit": limit})
        response.raise_for_status()
        return response.json()["data"]
    
    def update_status(self, execution_uuid, status, result_data=None):
        url = f"{self.base_url}/api/v1/sdk/executions/{execution_uuid}/status"
        data = {"status": status}
        if result_data:
            data["result_data"] = result_data
        response = requests.put(url, headers=self.headers, json=data)
        response.raise_for_status()
        return response.json()["data"]
    
    def append_logs(self, execution_uuid, logs):
        url = f"{self.base_url}/api/v1/sdk/executions/{execution_uuid}/logs"
        response = requests.post(url, headers=self.headers, json={"logs": logs})
        response.raise_for_status()
        return response.json()["data"]

# Usage
sdk = CronObserverSDK("http://localhost:8080", "your-api-key")

# Get pending executions
executions = sdk.get_pending_executions("550e8400-e29b-41d4-a716-446655440000")

for execution in executions:
    exec_uuid = execution["execution_uuid"]
    
    # Start execution
    sdk.update_status(exec_uuid, "RUNNING")
    
    # Append logs
    sdk.append_logs(exec_uuid, [
        {"level": "INFO", "message": "Task started"},
        {"level": "INFO", "message": "Fetching data..."}
    ])
    
    # Do actual work
    # ... execute task ...
    
    # Finish execution
    sdk.update_status(exec_uuid, "FINISHED", {
        "records_processed": 150
    })
```

### Node.js SDK Example

```javascript
class CronObserverSDK {
    constructor(baseUrl, apiKey) {
        this.baseUrl = baseUrl;
        this.headers = { "X-API-Key": apiKey };
    }
    
    async getPendingExecutions(taskUuid, limit = 10) {
        const response = await fetch(
            `${this.baseUrl}/api/v1/sdk/tasks/${taskUuid}/executions/pending?limit=${limit}`,
            { headers: this.headers }
        );
        const data = await response.json();
        return data.data;
    }
    
    async updateStatus(executionUuid, status, resultData = null) {
        const response = await fetch(
            `${this.baseUrl}/api/v1/sdk/executions/${executionUuid}/status`,
            {
                method: "PUT",
                headers: { ...this.headers, "Content-Type": "application/json" },
                body: JSON.stringify({ status, result_data: resultData })
            }
        );
        const data = await response.json();
        return data.data;
    }
    
    async appendLogs(executionUuid, logs) {
        const response = await fetch(
            `${this.baseUrl}/api/v1/sdk/executions/${executionUuid}/logs`,
            {
                method: "POST",
                headers: { ...this.headers, "Content-Type": "application/json" },
                body: JSON.stringify({ logs })
            }
        );
        const data = await response.json();
        return data.data;
    }
}

// Usage
const sdk = new CronObserverSDK("http://localhost:8080", "your-api-key");

async function executeTask() {
    const executions = await sdk.getPendingExecutions("550e8400-e29b-41d4-a716-446655440000");
    
    for (const execution of executions) {
        await sdk.updateStatus(execution.execution_uuid, "RUNNING");
        
        await sdk.appendLogs(execution.execution_uuid, [
            { level: "INFO", message: "Task started" }
        ]);
        
        // Do work...
        
        await sdk.updateStatus(execution.execution_uuid, "FINISHED", {
            records_processed: 150
        });
    }
}
```

## Best Practices

### For External Systems

1. **Polling Interval**: Poll for pending executions every 30-60 seconds
2. **Status Updates**: Update to RUNNING immediately when starting work
3. **Log Frequency**: Append logs periodically (not every line)
4. **Error Handling**: Always update status to FAILED on errors
5. **Idempotency**: Handle duplicate status updates gracefully
6. **Timeout Handling**: Set reasonable timeouts for API calls

### For Cron Observer

1. **Validation**: Strictly validate all status transitions
2. **Atomicity**: Use database transactions for status updates
3. **Idempotency**: Allow duplicate status updates (idempotent)
4. **Log Limits**: Consider max log entries per execution (e.g., 10,000)
5. **Monitoring**: Track SDK API usage and errors

## Next Steps

After completing this module:
1. Implement SDK endpoints
2. Add authentication middleware
3. Add rate limiting
4. Add status transition validation
5. Create SDK client libraries (optional)
6. Proceed to Module 7: Execution Tracking

