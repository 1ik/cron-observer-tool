#!/bin/sh
set -e

# Load environment variables from .env file if it exists
# Note: Environment variables set by docker-compose take precedence
# This is mainly for variables not set by docker-compose
if [ -f /app/.env ]; then
  # Use a more robust method to load .env file
  # But don't override variables already set (from docker-compose environment section)
  set -a
  . /app/.env
  set +a
fi

# Function to handle shutdown
cleanup() {
  echo "Shutting down services..."
  kill $BACKEND_PID $UI_PID 2>/dev/null || true
  wait
  exit 0
}

# Trap signals for graceful shutdown
trap cleanup SIGTERM SIGINT

# Start backend server
echo "Starting backend server on port ${SERVER_PORT:-8080}..."
/app/server &
BACKEND_PID=$!

# Start UI server
echo "Starting UI server on port ${UI_PORT:-3000}..."
cd /app/ui/apps/web
export PORT=${UI_PORT:-3000}
pnpm start &
UI_PID=$!

# Wait for all background processes
wait $BACKEND_PID $UI_PID

