# Module 10: Deployment & DevOps

## Overview

This module covers deployment strategies, containerization, environment configuration, CI/CD pipelines, monitoring, and operational considerations for Cron Observer.

## Deployment Architecture

### Deployment Options

#### Option 1: Docker Compose (Development/Simple Production)

```
┌─────────────────────────────────────┐
│         Docker Compose              │
├─────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐       │
│  │ Backend  │  │ Frontend  │       │
│  │ Container│  │ Container │       │
│  └──────────┘  └──────────┘       │
│  ┌──────────┐                      │
│  │PostgreSQL │                      │
│  │ Container │                      │
│  └──────────┘                      │
└─────────────────────────────────────┘
```

#### Option 2: Kubernetes (Production/Scale)

```
┌─────────────────────────────────────┐
│         Kubernetes Cluster          │
├─────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐        │
│  │ Backend  │  │ Frontend  │        │
│  │   Pods   │  │   Pods    │        │
│  └──────────┘  └──────────┘        │
│  ┌──────────┐                      │
│  │PostgreSQL │                      │
│  │  Service  │                      │
│  └──────────┘                      │
└─────────────────────────────────────┘
```

## Docker Configuration

### Backend Dockerfile

```dockerfile
# backend/Dockerfile
FROM openjdk:17-jdk-slim AS build

WORKDIR /app

# Copy build files
COPY gradle/ ./gradle/
COPY build.gradle settings.gradle ./
COPY src/ ./src/

# Build application
RUN ./gradlew build -x test

# Runtime stage
FROM openjdk:17-jre-slim

WORKDIR /app

# Copy built JAR
COPY --from=build /app/build/libs/*.jar app.jar

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s \
  CMD curl -f http://localhost:8080/actuator/health || exit 1

# Run application
ENTRYPOINT ["java", "-jar", "app.jar"]
```

### Frontend Dockerfile

```dockerfile
# frontend/Dockerfile
# Build stage
FROM node:18-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine

COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: cron_observer
      POSTGRES_USER: cron_observer_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U cron_observer_user"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      SPRING_DATASOURCE_URL: jdbc:postgresql://postgres:5432/cron_observer
      SPRING_DATASOURCE_USERNAME: cron_observer_user
      SPRING_DATASOURCE_PASSWORD: ${DB_PASSWORD}
      SPRING_PROFILES_ACTIVE: prod
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/actuator/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  postgres_data:
```

## Environment Configuration

### Configuration Files

#### application-prod.properties

```properties
# Database
spring.datasource.url=jdbc:postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME}
spring.datasource.username=${DB_USER}
spring.datasource.password=${DB_PASSWORD}

# JPA
spring.jpa.hibernate.ddl-auto=validate
spring.jpa.show-sql=false
spring.jpa.properties.hibernate.format_sql=false

# Flyway
spring.flyway.enabled=true
spring.flyway.baseline-on-migrate=true

# Server
server.port=8080
server.error.include-stacktrace=never
server.error.include-message=never

# Logging
logging.level.root=INFO
logging.level.com.cronobserver=DEBUG
logging.file.name=/var/log/cron-observer/application.log

# Scheduler
scheduler.enabled=true
scheduler.poll-interval-seconds=60

# API
api.rate-limit.enabled=true
api.rate-limit.requests-per-minute=100
```

### Environment Variables

```bash
# .env.example
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cron_observer
DB_USER=cron_observer_user
DB_PASSWORD=change-me

API_KEY_SECRET=change-me-secret-key
JWT_SECRET=change-me-jwt-secret

LOG_LEVEL=INFO
SERVER_PORT=8080
```

## Kubernetes Deployment

### Backend Deployment

```yaml
# k8s/backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cron-observer-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: cron-observer-backend
  template:
    metadata:
      labels:
        app: cron-observer-backend
    spec:
      containers:
      - name: backend
        image: cron-observer/backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: SPRING_DATASOURCE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        - name: SPRING_DATASOURCE_USERNAME
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: username
        - name: SPRING_DATASOURCE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /actuator/health
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /actuator/health/readiness
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: cron-observer-backend
spec:
  selector:
    app: cron-observer-backend
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Frontend Deployment

```yaml
# k8s/frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cron-observer-frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: cron-observer-frontend
  template:
    metadata:
      labels:
        app: cron-observer-frontend
    spec:
      containers:
      - name: frontend
        image: cron-observer/frontend:latest
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: cron-observer-frontend
spec:
  selector:
    app: cron-observer-frontend
  ports:
  - port: 80
    targetPort: 80
  type: LoadBalancer
```

## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run backend tests
        run: |
          cd backend
          ./gradlew test
      
      - name: Run frontend tests
        run: |
          cd frontend
          npm ci
          npm test

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build backend
        run: |
          cd backend
          ./gradlew build
      
      - name: Build frontend
        run: |
          cd frontend
          npm ci
          npm run build
      
      - name: Build Docker images
        run: |
          docker build -t cron-observer/backend:latest ./backend
          docker build -t cron-observer/frontend:latest ./frontend
      
      - name: Push to registry
        if: github.ref == 'refs/heads/main'
        run: |
          echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin
          docker push cron-observer/backend:latest
          docker push cron-observer/frontend:latest

  deploy:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to Kubernetes
        run: |
          kubectl apply -f k8s/
          kubectl rollout restart deployment/cron-observer-backend
          kubectl rollout restart deployment/cron-observer-frontend
```

## Monitoring & Logging

### Health Checks

```java
// Health check endpoint
@RestController
@RequestMapping("/actuator/health")
public class HealthController {
    
    @Autowired
    private DataSource dataSource;
    
    @GetMapping
    public ResponseEntity<HealthStatus> health() {
        HealthStatus status = new HealthStatus();
        status.setStatus("UP");
        status.setDatabase(checkDatabase());
        status.setScheduler(checkScheduler());
        return ResponseEntity.ok(status);
    }
    
    private boolean checkDatabase() {
        try {
            dataSource.getConnection().isValid(5);
            return true;
        } catch (SQLException e) {
            return false;
        }
    }
}
```

### Logging Configuration

```xml
<!-- logback-spring.xml -->
<configuration>
    <appender name="FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>/var/log/cron-observer/application.log</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
            <fileNamePattern>/var/log/cron-observer/application-%d{yyyy-MM-dd}.log</fileNamePattern>
            <maxHistory>30</maxHistory>
        </rollingPolicy>
        <encoder>
            <pattern>%d{yyyy-MM-dd HH:mm:ss} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>
    
    <root level="INFO">
        <appender-ref ref="FILE" />
    </root>
</configuration>
```

### Metrics (Prometheus)

```java
@Configuration
public class MetricsConfig {
    
    @Bean
    public MeterRegistryCustomizer<MeterRegistry> metricsCommonTags() {
        return registry -> registry.config()
            .commonTags("application", "cron-observer");
    }
}
```

## Database Backup

### Backup Strategy

```bash
#!/bin/bash
# scripts/backup-database.sh

BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/cron_observer_$DATE.sql"

pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > $BACKUP_FILE

# Compress
gzip $BACKUP_FILE

# Keep only last 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
```

### Scheduled Backup (Cron)

```bash
# Add to crontab
0 2 * * * /path/to/backup-database.sh
```

## Security Considerations

### Secrets Management

- **Kubernetes**: Use Secrets
- **Docker Compose**: Use environment files (not in git)
- **Production**: Use secret management services (AWS Secrets Manager, HashiCorp Vault)

### Network Security

- Use HTTPS in production
- Restrict database access
- Use firewall rules
- Implement rate limiting

### Application Security

- Input validation
- SQL injection prevention (use parameterized queries)
- XSS prevention (frontend)
- CSRF protection
- API authentication

## Scaling Considerations

### Horizontal Scaling

- **Backend**: Stateless, can scale horizontally
- **Frontend**: Stateless, can scale horizontally
- **Database**: Use read replicas for queries
- **Scheduler**: Only one instance should run scheduler (use leader election)

### Leader Election (Scheduler)

```java
@Component
public class SchedulerLeaderElection {
    
    @Scheduled(fixedDelay = 30000) // Every 30 seconds
    public void electLeader() {
        // Use database lock or distributed lock (Redis, etc.)
        if (acquireLock()) {
            schedulerService.evaluateSchedules();
        }
    }
    
    private boolean acquireLock() {
        // Implement distributed lock
        // Return true if this instance is the leader
    }
}
```

## Disaster Recovery

### Backup & Restore

1. **Database Backups**: Daily automated backups
2. **Configuration Backups**: Version control
3. **Restore Procedures**: Documented and tested

### High Availability

- Multiple backend instances
- Database replication
- Load balancer
- Health checks and auto-restart

## Documentation

### Deployment Documentation

- Setup instructions
- Configuration guide
- Troubleshooting guide
- Runbook for operations

## Next Steps

After completing this module:
1. Set up Docker configuration
2. Create Kubernetes manifests
3. Set up CI/CD pipeline
4. Configure monitoring
5. Set up backups
6. Document deployment procedures
7. **Project is ready for implementation!**

---

## Implementation Checklist

- [ ] Module 1: Project Structure ✅
- [ ] Module 2: Data Models ✅
- [ ] Module 3: Database Setup ✅
- [ ] Module 4: API Endpoints ✅
- [ ] Module 5: Scheduler Engine ✅
- [ ] Module 6: SDK/API ✅
- [ ] Module 7: Execution Tracking ✅
- [ ] Module 8: Frontend UI ✅
- [ ] Module 9: Testing Strategy ✅
- [ ] Module 10: Deployment ✅

**All modules documented. Ready to begin implementation!**

