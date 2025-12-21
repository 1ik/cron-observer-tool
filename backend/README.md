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

# Hello World endpoint
curl http://localhost:8080/api/v1/hello
```

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
├── cmd/
│   ├── server/          # Main API server
│   │   └── main.go
│   └── migrate/         # Migration CLI
│       └── main.go
├── internal/
│   ├── database/        # Database connection & collections
│   │   ├── mongo.go
│   │   └── collections.go
│   └── models/          # Data models
│       ├── project.go
│       └── task.go
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
- `name` (string) - Task name
- `description` (string) - Optional description
- `schedule_type` (enum) - RECURRING or ONEOFF
- `status` (enum) - ACTIVE, PAUSED, or DISABLED
- `schedule_config` (object) - Schedule configuration
- `metadata` (object) - Custom metadata
- `notification_config` (object) - Notification settings
- `created_at` (timestamp)
- `updated_at` (timestamp)

**Indexes**: uuid, project_id, status, schedule_type, created_at, project_status (compound), project_created (compound)

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
go run cmd/migrate/main.go create-collections

# View available commands
go run cmd/migrate/main.go --help
```

## Next Steps

- Add API handlers for Project CRUD operations
- Add API handlers for Task CRUD operations
- Implement API key authentication middleware
- Add repository pattern for database operations

