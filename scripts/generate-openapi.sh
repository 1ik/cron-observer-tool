#!/bin/bash

# Generate OpenAPI documentation and API client
# This script generates:
# 1. Backend Swagger/OpenAPI docs from Go code annotations
# 2. Frontend TypeScript API client from OpenAPI spec

set -e

echo "🚀 Generating OpenAPI documentation and API client..."

# Get the root directory of the project
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
UI_DIR="$ROOT_DIR/UI"

# Step 1: Generate backend Swagger docs
echo ""
echo "📝 Step 1: Generating backend Swagger documentation..."
cd "$BACKEND_DIR"

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "⚠️  swag tool not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Add GOPATH/bin to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin

# Generate swagger docs
swag init -g ./cmd/server/main.go -o api-docs

if [ $? -ne 0 ]; then
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi

echo "✅ Backend Swagger documentation generated successfully"
echo "   - api-docs/swagger.json"
echo "   - api-docs/swagger.yaml"

# Step 2: Convert Swagger 2.0 to OpenAPI 3.0 and generate frontend client
echo ""
echo "📝 Step 2: Converting to OpenAPI 3.0 and generating frontend API client..."
cd "$UI_DIR"

# Check if swagger2openapi is available
if ! command -v npx &> /dev/null; then
    echo "❌ npx not found. Please install Node.js and npm"
    exit 1
fi

# Convert Swagger 2.0 to OpenAPI 3.0
echo "   Converting Swagger 2.0 to OpenAPI 3.0..."
npx swagger2openapi "$BACKEND_DIR/api-docs/swagger.json" -o "$BACKEND_DIR/api-docs/openapi.json"

if [ $? -ne 0 ]; then
    echo "❌ Failed to convert Swagger to OpenAPI 3.0"
    exit 1
fi

echo "✅ OpenAPI 3.0 specification generated"
echo "   - backend/api-docs/openapi.json"

# Generate frontend API client
echo "   Generating TypeScript API client..."
npx openapi-zod-client "$BACKEND_DIR/api-docs/openapi.json" -o ./packages/lib/src/api-client.ts

if [ $? -ne 0 ]; then
    echo "❌ Failed to generate frontend API client"
    exit 1
fi

echo "✅ Frontend API client generated successfully"
echo "   - UI/packages/lib/src/api-client.ts"

echo ""
echo "🎉 OpenAPI generation complete!"
echo ""
echo "Generated files:"
echo "  Backend:"
echo "    - backend/api-docs/swagger.json"
echo "    - backend/api-docs/swagger.yaml"
echo "    - backend/api-docs/openapi.json"
echo "  Frontend:"
echo "    - UI/packages/lib/src/api-client.ts"

