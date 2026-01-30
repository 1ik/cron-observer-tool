# Cron Observer

An open-source task scheduling and execution tracking system where external systems execute tasks and report status/logs via SDK.

## Project Status

### Implemented Features
- ✅ **Backend**: Complete REST API with OpenAPI documentation
- ✅ **Scheduler**: Cron-based task scheduling with timezone support
- ✅ **Execution Tracking**: Real-time status updates, log management, pagination
- ✅ **Statistics**: Pre-aggregated failure stats, execution analytics
- ✅ **Frontend**: Modern React/Next.js UI with Radix UI components
- ✅ **Alerting**: Email notifications for execution failures via Gmail
- ✅ **Event-Driven Architecture**: Event bus for decoupled services
- ✅ **Task Timeouts**: Configurable execution timeouts with automatic failure handling

### Current Architecture
- **Backend**: Go (Gin framework) with MongoDB
- **Frontend**: Next.js 14+ with React Query, Radix UI
- **Event System**: In-memory event bus (ready for message queue upgrade)
- **Email**: Gmail SMTP integration
- **Documentation**: Auto-generated OpenAPI/Swagger specs

See [backend/README.md](backend/README.md) for development setup and API documentation.

## Key Features

- **Task Scheduling**: Complex cron-like scheduling with timezone support
- **Task Groups**: Group tasks together with time windows and coordinated control
- **External Execution**: Tasks executed by external systems, not by Cron Observer
- **Status Tracking**: Real-time execution status updates via SDK
- **Log Management**: Append-only logs with timestamps and levels
- **Execution History**: Complete history with date-based navigation and pagination
- **Execution Statistics**: Pre-aggregated stats (failures, success, totals) with 6-hour refresh
- **Task Failure Alerts**: Automatic email notifications to project users on execution failures
- **Task Timeouts**: Configurable execution timeouts with automatic failure handling
- **UUID-Based**: Tasks and executions use UUIDs for external reference
- **OpenAPI Specification**: Auto-generated API documentation (swagger.json/yaml)
- **Event-Driven Architecture**: Decoupled services using event bus pattern

## Generating OpenAPI Documentation

To generate both backend OpenAPI docs and frontend TypeScript API client:

```bash
./scripts/generate-openapi.sh
```

This script:
1. Generates backend Swagger/OpenAPI documentation from Go code annotations
2. Converts Swagger 2.0 to OpenAPI 3.0 format
3. Generates TypeScript API client for the frontend

See [backend/README.md](backend/README.md) for more details.

## Architecture

```
Cron Observer (Scheduler & Tracker)
    ↓ Creates execution records (PENDING)
    ↓
External Systems (Execute actual work)
    ↓ Report status & logs via SDK
    ↓
Cron Observer (Tracks & Displays)
```

## Documentation Structure

- **docs/MODULE_XX_*.md**: Phase-by-phase implementation guides
- Each module is self-contained and can be implemented independently
- **backend/docs/DURABLE_TASK_DELETE.md**: Design and flow for the durable, asynchronous task deletion pipeline (cron stop → hard delete → event/ack)

## Development Approach

This project follows a **modular, phase-by-phase** development approach:

1. **Plan First**: Each module has complete documentation before implementation
2. **Incremental**: Build and test each module before moving to the next
3. **Documentation-Driven**: Code follows documented specifications

## License

[To be determined - Open source license]

## Contributing

[Contributing guidelines to be added]

---