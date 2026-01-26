# Cron Observer Deployment Guide

This guide explains how to build, push, and deploy the Cron Observer Docker container.

## Building and Pushing to Docker Hub

### Prerequisites

- Docker installed and running
- Docker Hub account
- Logged in to Docker Hub: `docker login -u <your-dockerhub-username>`

### Build and Push

Run the build script from the project root:

```bash
./scripts/build-and-push.sh [version-tag]
```

**Examples:**
```bash
# Build and push as latest
./scripts/build-and-push.sh

# Build and push with version tag
./scripts/build-and-push.sh v1.0.0
```

The script will:
1. Build the Docker image as `<your-dockerhub-username>/cronoverver:latest` (and version tag if provided)
2. Ask if you want to push to Docker Hub
3. Push the image(s) to Docker Hub

**Note:** Update `DOCKER_USERNAME` in `scripts/build-and-push.sh` with your Docker Hub username.

## Server Deployment

### Prerequisites

- Docker installed on the server
- MongoDB accessible (either on the server or remote)
- Reverse proxy configured (nginx/traefik) for subdomain deployment

### Step 1: Create .env File

Create a `.env` file on your server with the following variables:

```env
# Server Ports
SERVER_PORT=8080
UI_PORT=3000

# Database Configuration
# When using docker-compose, use the service name 'mongodb' instead of 'localhost'
DATABASE_URI=mongodb://mongodb:27017
DATABASE_NAME=cronobserver

# Authentication
JWT_SECRET=your-secure-jwt-secret-here
SUPER_ADMINS=admin@example.com

# Gmail Configuration (optional, for email alerts)
GMAIL_USER=your-email@gmail.com
GMAIL_APP_PASSWORD=your-app-password

# Example Client Configuration (if running separately)
CRON_OBSERVER_URL=http://localhost:8080
CRON_OBSERVER_API_KEY=your-project-api-key-here

# NextAuth Configuration
NEXTAUTH_SECRET=your-nextauth-secret-here
NEXTAUTH_URL=https://cron.example.com
AUTH_TRUST_HOST=true
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Docker Hub (for build script)
DOCKER_PAT_TOKEN=your-docker-hub-pat-token
```

### Step 2: Deploy with Docker Compose (Recommended)

The easiest way to deploy is using Docker Compose, which ensures both MongoDB and the app are on the same network:

```bash
# Pull the latest image
docker pull <your-dockerhub-username>/cronoverver:latest

# Start all services (MongoDB + App)
docker-compose up -d
```

**Important:** In your `.env` file, set:
```env
DATABASE_URI=mongodb://mongodb:27017
```

This uses the Docker Compose service name `mongodb` instead of `localhost`.

### Alternative: Run Container Manually

If you prefer to run containers separately, you need to connect them to the same Docker network:

```bash
# Create a network
docker network create cron-observer-network

# Start MongoDB
docker run -d \
  --name cronobserver-mongodb \
  --network cron-observer-network \
  -p 27017:27017 \
  -v mongo-data:/data/db \
  mongo:7.0

# Start the app (connect to same network)
docker run -d \
  --name cron-observer \
  --network cron-observer-network \
  -p 8080:8080 \
  -p 3000:3000 \
  -v $(pwd)/.env:/app/.env:ro \
  --env-file .env \
  -e DATABASE_URI=mongodb://cronobserver-mongodb:27017 \
  --restart unless-stopped \
  <your-dockerhub-username>/cronoverver:latest
```

### Step 3: Configure Reverse Proxy (for subdomain deployment)

Configure your reverse proxy for subdomain deployment (e.g., `https://cron.example.com`):

#### Nginx Configuration

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name cron.example.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name cron.example.com;

    # SSL configuration (adjust paths as needed)
    ssl_certificate /path/to/ssl/cert.pem;
    ssl_certificate_key /path/to/ssl/key.pem;

    # Backend API (Go server) - must come before / location
    location /api {
        proxy_pass http://127.0.0.1:8080/api;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Connection "";
    }

    # UI (Next.js frontend) - catch-all for everything else
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

This configuration:
- Routes `/api/*` → Backend API on port 8080 (for Go backend)
- Routes `/*` → UI on port 3000 (Next.js handles `/api/auth/*` internally)

### Step 4: Verify Deployment

1. **Check container status:**
   ```bash
   docker ps | grep cron-observer
   ```

2. **Check logs:**
   ```bash
   docker logs cron-observer
   ```

3. **Access the application:**
   - UI: `https://cron.example.com`
   - API: `https://cron.example.com/api/v1/health`

## Container Architecture

The container runs two services in parallel:

- **Backend (Go)**: Port 8080 - Main API server
- **UI (Next.js)**: Port 3000 - Frontend application

All services start automatically via the entrypoint script and handle graceful shutdown.

## Environment Variables

### Required Variables

- `DATABASE_URI` - MongoDB connection string
- `DATABASE_NAME` - MongoDB database name
- `JWT_SECRET` - Secret for JWT token signing
- `NEXTAUTH_SECRET` - Secret for NextAuth.js
- `NEXTAUTH_URL` - Public URL for NextAuth callbacks (e.g., `https://cron.example.com`)
- `AUTH_TRUST_HOST` - Set to `true` when behind a reverse proxy
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret

### Optional Variables

- `SERVER_PORT` - Backend port (default: 8080)
- `UI_PORT` - UI port (default: 3000)
- `SUPER_ADMINS` - Comma-separated list of super admin emails
- `GMAIL_USER` - Gmail address for alerts
- `GMAIL_APP_PASSWORD` - Gmail app password
- `CRON_OBSERVER_API_KEY` - API key for example client

## Troubleshooting

### Container won't start

Check logs:
```bash
docker logs cron-observer
```

Common issues:
- Missing required environment variables
- MongoDB not accessible
- Port conflicts

### Services not accessible

1. Verify ports are exposed: `docker ps`
2. Check firewall rules
3. Verify reverse proxy configuration

### API calls failing

1. Check reverse proxy routes `/api` to port 8080
2. Verify backend is running: `docker logs cron-observer | grep backend`
3. Check browser console for CORS or network errors

## Updating Deployment

To update to a new version:

```bash
# Stop and remove old container
docker stop cron-observer
docker rm cron-observer

# Pull latest image (replace with your Docker Hub username)
docker pull <your-dockerhub-username>/cronoverver:latest

# Run new container (same command as initial deployment)
docker run -d \
  --name cron-observer \
  -p 8080:8080 \
  -p 3000:3000 \
  -v $(pwd)/.env:/app/.env:ro \
  --restart unless-stopped \
  <your-dockerhub-username>/cronoverver:latest
```

