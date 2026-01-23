#!/bin/bash

echo "🛑 Stopping all services..."

# Kill processes by port
echo "Killing processes on ports 8080, 5202, 5201..."
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
lsof -ti:5202 | xargs kill -9 2>/dev/null || true
lsof -ti:5201 | xargs kill -9 2>/dev/null || true

# Kill any remaining Go processes related to the backend
pkill -f "go run cmd/server/main.go" 2>/dev/null || true
pkill -f "go run cmd/migrate/main.go" 2>/dev/null || true

# Kill any remaining Node/npm/pnpm processes related to the services
pkill -f "nest start" 2>/dev/null || true
pkill -f "next dev" 2>/dev/null || true

echo "✅ Cleanup complete!"

