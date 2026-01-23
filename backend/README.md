# Cron Observer Backend

## Quick Start

### 1. Start MongoDB

```bash
# From project root
docker-compose up -d
```

This will start MongoDB on `localhost:27017` with data persisted in a Docker volume.

### 2. Run Migrations

```bash
# From backend directory
go run cmd/migrate/main.go create-collections
```

This will:
- Create the `projects` collection with indexes
- Create the `tasks` collection with indexes
- Create the `task_groups` collection with indexes

### 3. Start the Server

```bash
# From backend directory
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### 4. Test the API

```bash
# Health check with database status
curl http://localhost:8080/api/v1/health
```

See [API Endpoints](#api-endpoints) section below for all available endpoints.

## Environment Variables

Create a `.env` file in the backend directory (optional):

```bash
MONGODB_URI=mongodb://localhost:27017
DB_NAME=cronobserver
```

Default values are used if not specified.

## Project Structure

```
backend/
├── api-docs/            # Generated OpenAPI specification
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── cmd/
│   ├── server/          # Main API server
│   │   └── main.go
│   └── migrate/         # Migration CLI
│       └── main.go
├── internal/
│   ├── database/        # Database connection & collections
│   │   ├── mongo.go
│   │   └── collections.go
│   ├── events/          # Event bus for event-driven architecture
│   │   ├── bus.go
│   │   └── event.go
│   ├── handlers/        # HTTP request handlers
│   │   ├── project_handler.go
│   │   ├── tasks.go
│   │   └── taskgroup_handler.go
│   ├── models/          # Data models
│   │   ├── error.go
│   │   ├── project.go
│   │   ├── task.go
│   │   └── taskgroup.go
│   ├── repositories/    # Repository pattern implementation
│   │   ├── mongo.go
│   │   └── repository.go
│   ├── scheduler/       # Scheduler engine
│   │   ├── group_job.go
│   │   ├── job.go
│   │   └── scheduler.go
│   ├── utils/           # Utility functions
│   │   ├── api_key.go
│   │   └── validation.go
│   └── validators/      # Custom validators
│       └── custom.go
├── go.mod
└── go.sum
```

## Database Schema

### Collections

#### Projects
- `uuid` (string, unique) - Public identifier
- `name` (string) - Project name
- `description` (string) - Optional description
- `api_key` (string, unique) - API key for authentication
- `created_at` (timestamp)
- `updated_at` (timestamp)

**Indexes**: uuid, api_key, created_at

#### Tasks
- `uuid` (string, unique) - Public identifier
- `project_id` (ObjectID) - Reference to project
- `task_group_id` (ObjectID, optional) - Reference to task group
- `name` (string) - Task name
- `description` (string) - Optional description
- `schedule_type` (enum) - RECURRING or ONEOFF
- `status` (enum) - ACTIVE, PAUSED, or DISABLED
- `schedule_config` (object) - Schedule configuration
  - `cron_expression` (string, optional) - Cron expression
  - `timezone` (string) - IANA timezone
  - `time_range` (object, optional) - Time range with frequency
  - `days_of_week` (array, optional) - Days of week (0-6)
  - `exclusions` (array, optional) - Excluded days
- `trigger_config` (object) - Trigger configuration (HTTP)
- `metadata` (object, optional) - Custom metadata
- `created_at` (timestamp)
- `updated_at` (timestamp)

**Indexes**: uuid, project_id, task_group_id, status, schedule_type, created_at, project_status (compound), project_created (compound)

#### Task Groups
- `uuid` (string, unique) - Public identifier
- `project_id` (ObjectID) - Reference to project
- `name` (string) - Task group name
- `description` (string, optional) - Optional description
- `status` (enum) - ACTIVE, PAUSED, or DISABLED
- `start_time` (string, optional) - Start time (HH:MM format)
- `end_time` (string, optional) - End time (HH:MM format)
- `timezone` (string, optional) - IANA timezone for time windows
- `created_at` (timestamp)
- `updated_at` (timestamp)

**Indexes**: uuid, project_id, status, created_at

## Development Commands

```bash
# Build server binary
go build -o server cmd/server/main.go

# Build migration CLI
go build -o migrate cmd/migrate/main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## Docker Commands

```bash
# Start MongoDB
docker-compose up -d

# Stop MongoDB
docker-compose down

# Stop MongoDB and remove data
docker-compose down -v

# View MongoDB logs
docker-compose logs -f mongodb

# Access MongoDB shell
docker exec -it cronobserver-mongodb mongosh
```

## Migration Commands

```bash
# Create collections and indexes
# This creates: projects, tasks, and task_groups collections with all indexes
go run cmd/migrate/main.go create-collections

# View available commands
go run cmd/migrate/main.go --help
```

## API Endpoints

All endpoints are under `/api/v1` base path.

### Projects

- `GET /projects` - Get all projects
- `POST /projects` - Create a new project

### Tasks

- `POST /projects/{project_id}/tasks` - Create a new task
- `PUT /projects/{project_id}/tasks/{task_uuid}` - Update a task
- `DELETE /projects/{project_id}/tasks/{task_uuid}` - Delete a task

### Task Groups

- `POST /projects/{project_id}/task-groups` - Create a new task group
- `GET /projects/{project_id}/task-groups/{group_uuid}` - Get a task group
- `PUT /projects/{project_id}/task-groups/{group_uuid}` - Update a task group
- `DELETE /projects/{project_id}/task-groups/{group_uuid}` - Delete a task group
- `POST /projects/{project_id}/task-groups/{group_uuid}/start` - Start all tasks in a group
- `POST /projects/{project_id}/task-groups/{group_uuid}/stop` - Stop all tasks in a group
- `GET /projects/{project_id}/task-groups/{group_uuid}/tasks` - Get all tasks in a group

### Health Check

- `GET /health` - Health check with database status

## OpenAPI Specification

The API is documented using OpenAPI v3 specification. The specification is auto-generated from code annotations using the `swag` tool.

### Generate OpenAPI Specification

#### Quick Start (Recommended)

Use the unified script to generate both backend docs and frontend client:

```bash
# From project root
./scripts/generate-openapi.sh
```

This will:
1. Generate backend Swagger/OpenAPI documentation from Go code annotations
2. Convert Swagger 2.0 to OpenAPI 3.0 format
3. Generate TypeScript API client for the frontend

#### Manual Generation

**Backend only:**

```bash
# Install swag tool (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate specification
export PATH=$PATH:$(go env GOPATH)/bin
swag init -g ./cmd/server/main.go -o api-docs
```

**Frontend only:**

```bash
# From UI directory
cd ../UI
pnpm gen:api
```

This generates:
- `api-docs/swagger.json` - OpenAPI specification in JSON format (Swagger 2.0)
- `api-docs/swagger.yaml` - OpenAPI specification in YAML format (Swagger 2.0)
- `api-docs/openapi.json` - OpenAPI 3.0 specification (for frontend client)
- `api-docs/docs.go` - Generated Go code (can be ignored)
- `UI/packages/lib/src/api-client.ts` - TypeScript API client

The specification documents all endpoints, request/response schemas, and error responses.

### Viewing the Specification

You can use the generated `swagger.json` file with:
- Postman (import OpenAPI spec)
- Insomnia (import OpenAPI spec)
- Swagger UI (if configured)
- API documentation tools

## Next Steps

- Implement SDK/API endpoints for external systems to report execution status
- Implement execution tracking system
- Add comprehensive testing
- Implement frontend UI
- Add API key authentication middleware

