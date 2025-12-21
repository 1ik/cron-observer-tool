# Module 3: Database Setup & Migrations

## Overview

This module covers database selection, connection configuration, schema implementation, migration strategy, and index optimization.

## Database Selection

### Recommended: PostgreSQL

**Rationale**:
- Excellent JSON support (for schedule_config, metadata, etc.)
- Robust foreign key constraints
- Strong ACID compliance
- Excellent performance for time-series queries (execution history)
- Widely used in production
- Good migration tool support

### Alternative: SQLite (Development/MVP)

**Use Case**: 
- Local development
- MVP/prototype
- Single-user deployments

**Limitations**:
- Limited concurrency
- No advanced JSON functions
- File-based (backup considerations)

## Database Configuration

### Connection Configuration

**Environment Variables**:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cron_observer
DB_USER=cron_observer_user
DB_PASSWORD=secure_password
DB_SSL_MODE=prefer
DB_POOL_SIZE=10
DB_MAX_POOL_SIZE=20
```

**Connection String Format**:
```
postgresql://user:password@host:port/database?sslmode=prefer
```

### Connection Pool Settings

- **Initial Pool Size**: 5 connections
- **Max Pool Size**: 20 connections
- **Min Idle**: 2 connections
- **Connection Timeout**: 30 seconds
- **Idle Timeout**: 10 minutes
- **Max Lifetime**: 1 hour

## Schema Design

### Database Creation

```sql
-- Create database
CREATE DATABASE cron_observer
    WITH ENCODING 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8';

-- Create user (if needed)
CREATE USER cron_observer_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE cron_observer TO cron_observer_user;
```

### Schema Implementation

See Module 2 for complete entity definitions. This section covers the SQL DDL.

## Migration Strategy

### Migration Tool: Flyway (Recommended) or Liquibase

**Flyway Advantages**:
- Version-based migrations
- Simple SQL files
- Easy rollback (with pro version)
- Good Spring Boot integration

**Migration File Naming**:
```
V1__Create_tasks_table.sql
V2__Create_executions_table.sql
V3__Create_execution_logs_table.sql
V4__Create_exclusions_table.sql
V5__Create_task_exclusions_junction.sql
V6__Create_users_table.sql
V7__Create_notification_logs_table.sql
V8__Add_indexes.sql
```

### Migration Files

#### V1__Create_tasks_table.sql

```sql
-- Create enums
CREATE TYPE schedule_type AS ENUM ('RECURRING', 'ONEOFF');
CREATE TYPE task_status AS ENUM ('ACTIVE', 'PAUSED', 'DISABLED');

-- Create tasks table
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    schedule_type schedule_type NOT NULL,
    status task_status NOT NULL DEFAULT 'ACTIVE',
    schedule_config JSONB NOT NULL,
    metadata JSONB,
    notification_config JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT
);

-- Create indexes
CREATE UNIQUE INDEX idx_tasks_uuid ON tasks(uuid);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_schedule_type ON tasks(schedule_type);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for updated_at
CREATE TRIGGER update_tasks_updated_at 
    BEFORE UPDATE ON tasks 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

#### V2__Create_executions_table.sql

```sql
-- Create enums
CREATE TYPE execution_status AS ENUM ('PENDING', 'RUNNING', 'FINISHED', 'FAILED', 'CANCELLED');
CREATE TYPE trigger_type AS ENUM ('SCHEDULED', 'MANUAL');

-- Create executions table
CREATE TABLE executions (
    id BIGSERIAL PRIMARY KEY,
    execution_uuid VARCHAR(36) UNIQUE NOT NULL,
    task_id BIGINT NOT NULL,
    task_uuid VARCHAR(36) NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    status execution_status NOT NULL DEFAULT 'PENDING',
    trigger_type trigger_type NOT NULL,
    duration_ms BIGINT,
    result_data JSONB,
    error_message TEXT,
    updated_by VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_executions_task 
        FOREIGN KEY (task_id) 
        REFERENCES tasks(id) 
        ON DELETE CASCADE
);

-- Create indexes
CREATE UNIQUE INDEX idx_executions_uuid ON executions(execution_uuid);
CREATE INDEX idx_executions_task_id ON executions(task_id);
CREATE INDEX idx_executions_task_uuid ON executions(task_uuid);
CREATE INDEX idx_executions_status ON executions(status);
CREATE INDEX idx_executions_scheduled_at ON executions(scheduled_at);
CREATE INDEX idx_executions_task_scheduled ON executions(task_id, scheduled_at);
CREATE INDEX idx_executions_status_scheduled ON executions(status, scheduled_at);

-- Create trigger for updated_at
CREATE TRIGGER update_executions_updated_at 
    BEFORE UPDATE ON executions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

#### V3__Create_execution_logs_table.sql

```sql
-- Create enum
CREATE TYPE log_level AS ENUM ('INFO', 'WARN', 'ERROR', 'DEBUG');

-- Create execution_logs table
CREATE TABLE execution_logs (
    id BIGSERIAL PRIMARY KEY,
    execution_id BIGINT NOT NULL,
    execution_uuid VARCHAR(36) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    level log_level NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_execution_logs_execution 
        FOREIGN KEY (execution_id) 
        REFERENCES executions(id) 
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_execution_logs_execution_id ON execution_logs(execution_id);
CREATE INDEX idx_execution_logs_execution_uuid ON execution_logs(execution_uuid);
CREATE INDEX idx_execution_logs_timestamp ON execution_logs(timestamp);
CREATE INDEX idx_execution_logs_execution_timestamp ON execution_logs(execution_id, timestamp);
```

#### V4__Create_exclusions_table.sql

```sql
-- Create enum
CREATE TYPE exclusion_type AS ENUM ('SINGLE_DATE', 'DATE_RANGE', 'RECURRING_YEARLY');

-- Create exclusions table
CREATE TABLE exclusions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type exclusion_type NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    recurring_pattern JSONB,
    applies_to_all_tasks BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_exclusions_type ON exclusions(type);
CREATE INDEX idx_exclusions_start_date ON exclusions(start_date);
CREATE INDEX idx_exclusions_applies_to_all ON exclusions(applies_to_all_tasks);

-- Create trigger for updated_at
CREATE TRIGGER update_exclusions_updated_at 
    BEFORE UPDATE ON exclusions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

#### V5__Create_task_exclusions_junction.sql

```sql
-- Create task_exclusions junction table
CREATE TABLE task_exclusions (
    task_id BIGINT NOT NULL,
    exclusion_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (task_id, exclusion_id),
    CONSTRAINT fk_task_exclusions_task 
        FOREIGN KEY (task_id) 
        REFERENCES tasks(id) 
        ON DELETE CASCADE,
    CONSTRAINT fk_task_exclusions_exclusion 
        FOREIGN KEY (exclusion_id) 
        REFERENCES exclusions(id) 
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_task_exclusions_task_id ON task_exclusions(task_id);
CREATE INDEX idx_task_exclusions_exclusion_id ON task_exclusions(exclusion_id);
```

#### V6__Create_users_table.sql

```sql
-- Create enum
CREATE TYPE user_role AS ENUM ('ADMIN', 'USER');

-- Create users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255),
    role user_role NOT NULL DEFAULT 'USER',
    api_key VARCHAR(255) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE UNIQUE INDEX idx_users_api_key ON users(api_key);

-- Create trigger for updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

#### V7__Create_notification_logs_table.sql

```sql
-- Create enums
CREATE TYPE notification_channel_type AS ENUM ('EMAIL', 'WEBHOOK', 'SLACK', 'DISCORD');
CREATE TYPE notification_status AS ENUM ('SENT', 'FAILED');

-- Create notification_logs table
CREATE TABLE notification_logs (
    id BIGSERIAL PRIMARY KEY,
    execution_id BIGINT NOT NULL,
    channel_type notification_channel_type NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status notification_status NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_notification_logs_execution 
        FOREIGN KEY (execution_id) 
        REFERENCES executions(id) 
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_notification_logs_execution_id ON notification_logs(execution_id);
CREATE INDEX idx_notification_logs_sent_at ON notification_logs(sent_at);
CREATE INDEX idx_notification_logs_status ON notification_logs(status);
```

#### V8__Add_indexes_and_optimizations.sql

```sql
-- Additional composite indexes for common queries

-- For finding pending executions by task
CREATE INDEX IF NOT EXISTS idx_executions_task_pending 
    ON executions(task_id, status) 
    WHERE status = 'PENDING';

-- For date-based execution queries
CREATE INDEX IF NOT EXISTS idx_executions_date_range 
    ON executions USING btree (DATE(scheduled_at), status);

-- For execution history queries
CREATE INDEX IF NOT EXISTS idx_executions_history 
    ON executions(task_id, scheduled_at DESC, status);

-- For log queries by execution
CREATE INDEX IF NOT EXISTS idx_execution_logs_by_execution 
    ON execution_logs(execution_id, timestamp ASC);
```

## Database Initialization

### Development Setup

1. **Create Database**:
   ```bash
   createdb cron_observer
   ```

2. **Run Migrations**:
   ```bash
   # Using Flyway CLI
   flyway migrate
   
   # Or using application
   # Migrations run automatically on startup
   ```

3. **Verify Schema**:
   ```sql
   \dt  -- List tables
   \d tasks  -- Describe tasks table
   ```

### Production Setup

1. **Database Provisioning**:
   - Use managed PostgreSQL service (AWS RDS, Azure, etc.)
   - Or self-hosted PostgreSQL cluster

2. **Connection Security**:
   - Use SSL/TLS connections
   - Restrict network access
   - Use strong passwords
   - Rotate credentials regularly

3. **Backup Strategy**:
   - Daily automated backups
   - Point-in-time recovery
   - Test restore procedures

## Index Strategy

### Primary Indexes
- All primary keys (automatic)
- All foreign keys
- All unique constraints

### Performance Indexes
- `executions(status, scheduled_at)`: For finding pending executions
- `executions(task_id, scheduled_at DESC)`: For execution history
- `execution_logs(execution_id, timestamp)`: For log retrieval
- `tasks(status)`: For filtering active tasks

### Partial Indexes (PostgreSQL)
- `executions` with `status = 'PENDING'`: For pending execution queries
- Consider partial indexes for common filters

## Query Optimization

### Common Queries

1. **Find Pending Executions**:
   ```sql
   SELECT * FROM executions 
   WHERE status = 'PENDING' 
   ORDER BY scheduled_at ASC;
   ```

2. **Get Execution History**:
   ```sql
   SELECT * FROM executions 
   WHERE task_id = ? 
   ORDER BY scheduled_at DESC 
   LIMIT 50;
   ```

3. **Get Logs for Execution**:
   ```sql
   SELECT * FROM execution_logs 
   WHERE execution_id = ? 
   ORDER BY timestamp ASC;
   ```

4. **Date-Based Executions**:
   ```sql
   SELECT * FROM executions 
   WHERE DATE(scheduled_at) = ? 
   ORDER BY scheduled_at DESC;
   ```

## Data Retention

### Execution History Retention
- Consider archiving old executions (>90 days)
- Move to archive table or separate database
- Keep summary statistics

### Log Retention
- Execution logs can grow large
- Consider log rotation/archival
- Keep recent logs (30 days) in main table
- Archive older logs

## Next Steps

After completing this module:
1. Set up local database
2. Run initial migrations
3. Verify schema creation
4. Proceed to Module 4: API Endpoints

## Notes

- Use connection pooling in production
- Monitor query performance
- Consider read replicas for scaling
- Plan for data archival strategy
- Use database migrations for all schema changes (never manual DDL)

