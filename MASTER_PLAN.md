# Cron Observer - Master Plan

## Project Overview

Cron Observer is an open-source task scheduling and execution tracking system. Unlike traditional cron schedulers that execute tasks directly, Cron Observer acts as a **scheduling and tracking platform** where:

1. **Tasks are defined** in Cron Observer with unique UUIDs
2. **External systems execute** the actual work
3. **External systems report back** execution status and logs via SDK/API
4. **Cron Observer tracks and displays** all execution history

## Project Goals

### Primary Goals
- Provide a centralized platform for scheduling tasks with complex cron-like patterns
- Enable external systems to execute tasks and report status/logs via SDK
- Offer comprehensive execution history and monitoring
- Support both recurring (cron-based) and one-off (scheduled) tasks
- Provide intuitive UI for task management and execution tracking

### Secondary Goals
- Open-source friendly architecture
- Extensible design for future features
- Self-hostable solution
- API-first approach for programmatic access

## Core Requirements

### Functional Requirements

#### Task Management
- **Task Definition**: Create tasks with UUIDs, names, descriptions
- **Schedule Types**: 
  - Recurring (cron expressions with timezone support)
  - One-off (single execution at future datetime)
- **Task Status**: ACTIVE, PAUSED, DISABLED
- **Task Control**: Pause, resume, manual trigger, cancel

#### Scheduling Engine
- **Cron Support**: Standard cron expression evaluation
- **Timezone Awareness**: IANA timezone support
- **Complex Patterns**: 
  - Day-of-week ranges (Sunday-Thursday)
  - Time windows (10am-2:30pm)
  - Multiple schedules per day
  - Holiday/exclusion support
- **Execution Creation**: Automatically create execution records with PENDING status

#### Execution Tracking
- **Status Lifecycle**: PENDING → RUNNING → FINISHED/FAILED
- **Log Management**: Append-only logs with timestamps and levels
- **Execution History**: Complete history per task, filterable by date/status
- **Date-Based Navigation**: UI left panel showing dates with executions

#### SDK/API for External Systems
- **Status Updates**: Allow external systems to update execution status
- **Log Appending**: Allow external systems to append logs
- **Execution Discovery**: Allow external systems to find pending executions
- **Authentication**: API key/token-based authentication

#### User Interface
- **Task List**: View all tasks with filtering and search
- **Task Detail**: View task configuration and execution history
- **Execution Detail**: View execution status, logs, and results
- **Date Navigation**: Left panel with date list, click to view executions
- **Dashboard**: Overview of tasks, executions, and statistics

### Non-Functional Requirements

#### Performance
- Efficient cron evaluation
- Fast execution history queries
- Scalable to handle many tasks and executions

#### Reliability
- Accurate scheduling
- Reliable execution tracking
- Data consistency

#### Security
- API authentication
- Secure storage of sensitive data (future)
- Rate limiting for API calls

#### Usability
- Intuitive UI/UX
- Clear execution status visualization
- Easy task configuration

## Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Cron Observer                        │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────┐  │
│  │   Scheduler  │───▶│   Database   │◀───│   API    │  │
│  │   Service    │    │              │    │  Server  │  │
│  └──────────────┘    └──────────────┘    └──────────┘  │
│         │                    │                  │        │
│         │                    │                  │        │
│         ▼                    ▼                  ▼        │
│  ┌──────────────────────────────────────────────────┐  │
│  │            Execution Tracking Engine              │  │
│  └──────────────────────────────────────────────────┘  │
│                                                           │
└─────────────────────────────────────────────────────────┘
         │                    │                  │
         │                    │                  │
         ▼                    ▼                  ▼
┌─────────────────────────────────────────────────────────┐
│              External Systems (Consumers)                │
│  - Discover pending executions via API                   │
│  - Execute actual work                                    │
│  - Report status and logs via SDK                        │
└─────────────────────────────────────────────────────────┘
```

### Component Responsibilities

1. **Scheduler Service**: Evaluates cron expressions, creates execution records
2. **API Server**: REST API for task management and SDK endpoints
3. **Database**: Stores tasks, executions, logs, configurations
4. **Execution Tracking Engine**: Manages execution lifecycle and status transitions
5. **UI Frontend**: User interface for task management and monitoring

## Implementation Modules

The implementation is divided into separate modules that can be developed phase by phase:

### Module 1: Project Structure & Backend Outline
**File**: `docs/MODULE_01_PROJECT_STRUCTURE.md`
- Overall project structure
- Backend architecture outline
- Technology stack decisions
- Project directory structure
- Development environment setup

### Module 2: Data Models
**File**: `docs/MODULE_02_DATA_MODELS.md`
- Complete data model definitions
- Entity relationships
- Field specifications
- Enums and constants
- Validation rules

### Module 3: Database Setup & Migrations
**File**: `docs/MODULE_03_DATABASE.md`
- Database selection and rationale
- Connection configuration
- Schema design
- Migration strategy
- Index definitions

### Module 4: Core API Endpoints
**File**: `docs/MODULE_04_API_ENDPOINTS.md`
- REST API design
- Endpoint specifications
- Request/response schemas
- Authentication/authorization
- Error handling

### Module 5: Scheduler Engine
**File**: `docs/MODULE_05_SCHEDULER.md`
- Cron expression evaluation
- Timezone handling
- Execution record creation
- Schedule validation

### Module 6: SDK/API for External Systems
**File**: `docs/MODULE_06_SDK_API.md`
- Status update endpoints
- Log appending endpoints
- Execution discovery endpoints
- SDK client examples

### Module 7: Execution Tracking
**File**: `docs/MODULE_07_EXECUTION_TRACKING.md`
- Execution lifecycle management
- Status transition validation
- Log storage and retrieval
- Statistics calculation

### Module 8: Frontend UI
**File**: `docs/MODULE_08_FRONTEND.md`
- UI component structure
- Task management views
- Execution history views
- Date-based navigation
- Dashboard components

### Module 9: Testing Strategy
**File**: `docs/MODULE_09_TESTING.md`
- Unit testing approach
- Integration testing
- API testing
- E2E testing strategy

### Module 10: Deployment & DevOps
**File**: `docs/MODULE_10_DEPLOYMENT.md`
- Docker containerization
- Environment configuration
- CI/CD pipeline
- Monitoring and logging

## Development Phases

### Phase 1: Foundation (Modules 1-3)
- Project structure setup
- Data models definition
- Database setup and migrations
- Basic project scaffolding

### Phase 2: Core Backend (Modules 4-7)
- API endpoints implementation
- Scheduler engine
- SDK/API for external systems
- Execution tracking

### Phase 3: Frontend (Module 8)
- UI components
- Task management interface
- Execution history views
- Date-based navigation

### Phase 4: Polish & Deploy (Modules 9-10)
- Testing implementation
- Deployment setup
- Documentation
- Open-source preparation

## Key Design Decisions

### Execution Model
- **External Execution**: Tasks are executed by external systems, not by Cron Observer
- **Status Tracking**: External systems report status via SDK
- **Log Appending**: Logs are append-only, cannot be modified
- **UUID-Based**: Tasks and executions use UUIDs for external reference

### Scheduling Model
- **Cron-Based**: Standard cron expressions for recurring tasks
- **Timezone-Aware**: All scheduling respects timezones
- **Execution Records**: Created automatically when scheduled, start as PENDING

### Data Model
- **Relational Database**: For tasks, configurations, relationships
- **Time-Series Consideration**: Execution history may benefit from time-series DB (future optimization)
- **JSON Fields**: Flexible metadata storage using JSON

## Success Criteria

### MVP Success Criteria
- [ ] Can create tasks with UUIDs
- [ ] Can schedule recurring and one-off tasks
- [ ] Scheduler creates execution records with PENDING status
- [ ] External systems can update execution status via API
- [ ] External systems can append logs via API
- [ ] Can view execution history
- [ ] Can pause/resume tasks
- [ ] Can manually trigger tasks

### Full Feature Success Criteria
- [ ] Complex scheduling patterns work correctly
- [ ] Date-based navigation in UI
- [ ] Complete execution tracking
- [ ] Notification system (future)
- [ ] Multi-user support (future)
- [ ] Comprehensive documentation
- [ ] Open-source ready

## Open Source Considerations

- Clear documentation
- Contribution guidelines
- License selection
- Example configurations
- SDK client libraries (future)
- Community-friendly architecture

## Notes

- Authentication is deferred to later phases (manual token for MVP)
- Multi-user support is planned for Phase 3
- Notifications are planned for Phase 2
- Advanced features (dependencies, chaining) are future considerations

---

**Last Updated**: 2025-01-XX
**Version**: 1.0.0
**Status**: Planning Phase

