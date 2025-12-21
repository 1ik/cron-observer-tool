# Module 1: Project Structure & Basic API Server

## Overview

**Goal**: Get a basic API server running with a simple "Hello World" endpoint.

This module focuses on:
- Setting up the Go project structure
- Installing Gin framework
- Creating a minimal working API server
- Testing the server with a Hello World endpoint

**Success Criteria**: 
- ✅ Server starts successfully
- ✅ `GET /api/v1/hello` returns `{"message": "Hello World"}`
- ✅ Server runs on port 8080 (or configured port)

## Project Structure (Minimal for Module 1)

### Directory Layout

For Module 1, we only need a minimal structure:

```
cron-observer/
├── README.md
├── MASTER_PLAN.md
├── docs/
│   └── MODULE_01_PROJECT_STRUCTURE.md (this file)
└── backend/
    ├── cmd/
    │   └── server/
    │       └── main.go          # Application entry point
    ├── go.mod                   # Go module file
    └── go.sum                   # Go dependencies (auto-generated)
```

**Note**: More complex structure (handlers, services, repositories) will be added in later modules. For now, we keep it simple.

## Technology Stack

### Backend: Go (Golang) + Gin Framework

**Selected**: 
- **Go 1.21+**: Programming language
- **Gin** (`github.com/gin-gonic/gin`): HTTP web framework

**Why Gin?**
- Fast and lightweight
- Perfect for REST APIs
- Great documentation
- Easy to get started

**Note**: Database, models, services, and other components will be added in later modules.

## Architecture (Module 1 - Simple)

For Module 1, we have a very simple architecture:

```
┌─────────────────────────────────────┐
│         main.go                      │  ← Entry point
│         - Gin router setup          │
│         - Hello World endpoint      │
└─────────────────────────────────────┘
```

**Note**: More complex architecture (handlers, services, repositories) will be added in later modules.

## Step-by-Step Setup

### Prerequisites
- **Go 1.21+**: [Download](https://go.dev/dl/)
- **Git**: Version control
- **Terminal/Command Line**: To run commands

### Step 1: Create Project Directory

```bash
mkdir -p cron-observer/backend/cmd/server
cd cron-observer
```

### Step 2: Initialize Go Module

```bash
cd backend
go mod init github.com/yourusername/cron-observer/backend
```

### Step 3: Install Gin Framework

```bash
go get github.com/gin-gonic/gin
```

### Step 4: Create Main Server File

Create `cmd/server/main.go`:

```go
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create Gin router
	router := gin.Default()

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Hello World endpoint
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello World",
			})
		})
	}

	// Start server on port 8080
	router.Run(":8080")
}
```

### Step 5: Run the Server

```bash
go run cmd/server/main.go
```

You should see:
```
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] GET    /api/v1/hello         --> main.main.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is not safe. We recommend you to set the value according to https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
[GIN-debug] Listening and serving HTTP on :8080
```

### Step 6: Test the API

Open a new terminal and test the endpoint:

```bash
# Using curl
curl http://localhost:8080/api/v1/hello

# Expected response:
# {"message":"Hello World"}
```

Or open in browser: `http://localhost:8080/api/v1/hello`

You should see:
```json
{"message":"Hello World"}
```

## Success Checklist

- [ ] Go module initialized
- [ ] Gin framework installed
- [ ] `main.go` file created
- [ ] Server starts without errors
- [ ] `GET /api/v1/hello` returns `{"message":"Hello World"}`
- [ ] Server runs on port 8080

## Dependencies (Module 1)

For Module 1, we only need:

```go
github.com/gin-gonic/gin  // HTTP web framework
```

### `go.mod` (After Setup)

Your `go.mod` should look like:

```go
module github.com/yourusername/cron-observer/backend

go 1.21

require github.com/gin-gonic/gin v1.9.1
```

**Note**: More dependencies will be added in later modules (database, cron parser, etc.)

## Understanding the Code

### What Each Part Does

```go
router := gin.Default()
```
- Creates a Gin router with default middleware (logging and recovery)

```go
api := router.Group("/api/v1")
```
- Creates a route group for API versioning
- All routes under this group will have `/api/v1` prefix

```go
api.GET("/hello", func(c *gin.Context) { ... })
```
- Defines a GET endpoint at `/api/v1/hello`
- The function is the handler that processes the request

```go
c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
```
- Sends a JSON response with status 200 (OK)
- `gin.H` is a shortcut for `map[string]interface{}`

```go
router.Run(":8080")
```
- Starts the HTTP server on port 8080

## Troubleshooting

### Common Issues

**Issue**: `go: cannot find module`
- **Solution**: Make sure you're in the `backend` directory and run `go mod init` first

**Issue**: Port 8080 already in use
- **Solution**: Change the port in `main.go`: `router.Run(":8081")`

**Issue**: Import errors
- **Solution**: Run `go mod tidy` to download dependencies

**Issue**: Server not starting
- **Solution**: Check that Go is installed: `go version`

## Next Steps

After completing this module:
1. ✅ Basic API server running
2. ✅ Hello World endpoint working
3. ✅ Server responds on port 8080

**Proceed to Module 2: Data Models** - where we'll define the data structures for tasks, executions, etc.

## Notes

- This is a minimal setup to get started
- More structure (handlers, services, repositories) will be added in later modules
- Database will be added in Module 3
- Keep it simple for now - just get the server running!

