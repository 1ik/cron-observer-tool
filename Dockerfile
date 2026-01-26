# Stage 1: Build Backend (Go)
FROM golang:1.23-alpine AS backend-builder

WORKDIR /build

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Stage 2: Build UI (Next.js)
FROM node:20-alpine AS ui-builder

# Enable corepack for pnpm (built into Node.js 20+)
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /build

# Copy UI package files and root tsconfig.json
COPY UI/package.json UI/pnpm-lock.yaml UI/pnpm-workspace.yaml UI/tsconfig.json ./
COPY UI/packages/ ./packages/
COPY UI/apps/ ./apps/

# Install dependencies
RUN pnpm install --frozen-lockfile

# Build UI
RUN pnpm build

# Stage 3: Runtime Container
FROM node:20-alpine

# Enable corepack for pnpm (built into Node.js 20+)
RUN corepack enable && corepack prepare pnpm@latest --activate

# Install dumb-init and tzdata for timezone support
RUN apk add --no-cache dumb-init tzdata

# Set timezone to Asia/Dhaka (UTC+6)
ENV TZ=Asia/Dhaka
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /build/server /app/server

# Copy UI - copy entire UI directory structure (needed for Next.js SSR)
COPY --from=ui-builder /build ./ui

# Install UI production dependencies only (remove dev dependencies)
WORKDIR /app/ui
RUN CI=true pnpm install --prod --frozen-lockfile --ignore-scripts

# Copy entrypoint script
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

WORKDIR /app

# Expose ports
EXPOSE 8080 3000

# Use dumb-init to handle signals properly
ENTRYPOINT ["dumb-init", "--"]
CMD ["/app/docker-entrypoint.sh"]

