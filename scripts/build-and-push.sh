#!/bin/bash
set -e

# Docker Hub username
DOCKER_USERNAME="anikgtx"
IMAGE_NAME="cronoverver"
FULL_IMAGE_NAME="${DOCKER_USERNAME}/${IMAGE_NAME}"

# Get version tag from argument or use latest
VERSION_TAG="${1:-latest}"

# Load environment variables from .env file if it exists
if [ -f ".env" ]; then
    # Extract variables from .env file (handles quoted and unquoted values)
    DOCKER_PAT_TOKEN=$(grep -E "^DOCKER_PAT_TOKEN=" ".env" | cut -d '=' -f2- | sed 's/^["'\'']//; s/["'\'']$//' | xargs)
    export DOCKER_PAT_TOKEN
fi

echo "Building Docker image: ${FULL_IMAGE_NAME}:${VERSION_TAG}"
echo "Platform: linux/amd64 (for server compatibility)"

# Build the Docker image for linux/amd64 platform
docker build \
  --platform linux/amd64 \
  -t "${FULL_IMAGE_NAME}:${VERSION_TAG}" .

# Also tag as latest if a specific version was provided
if [ "$VERSION_TAG" != "latest" ]; then
  echo "Tagging as latest..."
  docker tag "${FULL_IMAGE_NAME}:${VERSION_TAG}" "${FULL_IMAGE_NAME}:latest"
fi

# Check if user wants to push (skip prompt if --push flag or AUTO_PUSH env var is set)
SHOULD_PUSH=false
if [[ "$*" == *"--push"* ]] || [ "${AUTO_PUSH}" == "true" ]; then
  SHOULD_PUSH=true
elif [ -t 0 ]; then
  # Only prompt if running interactively
  read -p "Push to Docker Hub? (y/n) " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    SHOULD_PUSH=true
  fi
fi

if [ "$SHOULD_PUSH" = true ]; then
  # Authenticate with Docker Hub using PAT token from .env
  if [ -n "$DOCKER_PAT_TOKEN" ]; then
    echo "🔐 Authenticating with Docker Hub using PAT token..."
    echo "$DOCKER_PAT_TOKEN" | docker login -u "${DOCKER_USERNAME}" --password-stdin
    if [ $? -eq 0 ]; then
      echo "✅ Successfully authenticated with Docker Hub"
    else
      echo "❌ Failed to authenticate with Docker Hub"
      exit 1
    fi
  elif docker info 2>/dev/null | grep -q "Username:"; then
    echo "✅ Already authenticated with Docker Hub"
  else
    echo "⚠️  Docker Hub authentication required"
    echo "   Please set DOCKER_PAT_TOKEN in your .env file"
    echo "   Or login manually: docker login -u ${DOCKER_USERNAME}"
    exit 1
  fi

  echo "Pushing ${FULL_IMAGE_NAME}:${VERSION_TAG} to Docker Hub..."
  docker push "${FULL_IMAGE_NAME}:${VERSION_TAG}"

  if [ "$VERSION_TAG" != "latest" ]; then
    echo "Pushing ${FULL_IMAGE_NAME}:latest to Docker Hub..."
    docker push "${FULL_IMAGE_NAME}:latest"
  fi

  echo "✅ Successfully pushed to Docker Hub!"
else
  echo "Skipping push. Image built locally: ${FULL_IMAGE_NAME}:${VERSION_TAG}"
fi

