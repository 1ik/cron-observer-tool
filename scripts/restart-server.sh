#!/bin/bash
set -e

# Script to restart services on the server
# This script should be run on the server where the application is deployed

echo "🔄 Restarting cron-observer services..."

# Pull the latest image
echo "📥 Pulling latest Docker image..."
docker pull anikgtx/cronoverver:latest

# Restart services using docker-compose
echo "🔄 Restarting containers..."
docker-compose down
docker-compose up -d

echo "✅ Services restarted successfully!"
echo ""
echo "📊 Container status:"
docker-compose ps

echo ""
echo "📋 Recent logs (last 20 lines):"
docker-compose logs --tail=20

