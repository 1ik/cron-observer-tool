# NestJS Example Client for Cron Observer

Simple NestJS application that exposes an execution endpoint for cron-observer integration.

## Features

- ✅ Exposes a single API endpoint (`POST /api/execute`) that receives task execution notifications
- ✅ Logs all execution requests from the cron server

## Prerequisites

- Node.js 18+ and npm/yarn/pnpm

## Installation

```bash
cd example-client-nestjs
npm install
```

## Running

### Development Mode

```bash
npm run start:dev
```

### Production Mode

```bash
npm run build
npm run start:prod
```

## Setup Instructions

1. **Start the NestJS server:**
   ```bash
   npm run start:dev
   ```
   The server will run on `http://localhost:5202` (or the port specified in `PORT` env variable).

2. **Create a project in cron-observer UI:**
   - Go to the cron-observer UI
   - Create a new project
   - Set the **Execution Endpoint** to: `http://localhost:5202/api/execute`
   - Save the project

3. **Create a task in cron-observer UI:**
   - In your project, create a new task
   - Set the schedule (e.g., every 10 seconds: `*/10 * * * * *`)
   - Set status to `ACTIVE`
   - Save the task

4. **Watch the logs:**
   - When the task executes, you'll see execution requests logged in the console

## API Endpoint

### POST /api/execute

Receives execution notifications from cron-observer.

**Request Body:**
```json
{
  "task_name": "string",
  "execution_id": "string"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Execution received",
  "received_at": "2025-01-05T15:30:00.000Z"
}
```

## Example Output

```
🚀 NestJS Example Client running on http://localhost:5202
📡 Execution endpoint: http://localhost:5202/api/execute

[When a task executes, you'll see:]
═══════════════════════════════════════════════════
📨 Received execution request from cron-observer
═══════════════════════════════════════════════════
Task Name: My Task
Execution ID: 507f1f77bcf86cd799439011
Full Body: {
  "task_name": "My Task",
  "execution_id": "507f1f77bcf86cd799439011"
}
Timestamp: 2025-01-05T15:30:00.000Z
═══════════════════════════════════════════════════
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port for NestJS server | `5202` |

## Next Steps

- Modify the execution endpoint to perform actual work
- Add database integration
- Implement retry logic
- Add authentication/authorization

