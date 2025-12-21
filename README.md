# Cron Observer

An open-source task scheduling and execution tracking system where external systems execute tasks and report status/logs via SDK.

## Quick Start

1. Read the [Master Plan](MASTER_PLAN.md) for project overview
2. Follow modules in order:
   - [Module 1: Project Structure](docs/MODULE_01_PROJECT_STRUCTURE.md)
   - [Module 2: Data Models](docs/MODULE_02_DATA_MODELS.md)
   - [Module 3: Database Setup](docs/MODULE_03_DATABASE.md)
   - [Module 4: API Endpoints](docs/MODULE_04_API_ENDPOINTS.md) (Coming soon)
   - [Module 5: Scheduler Engine](docs/MODULE_05_SCHEDULER.md) (Coming soon)
   - [Module 6: SDK/API](docs/MODULE_06_SDK_API.md) (Coming soon)
   - [Module 7: Execution Tracking](docs/MODULE_07_EXECUTION_TRACKING.md) (Coming soon)
   - [Module 8: Frontend](docs/MODULE_08_FRONTEND.md) (Coming soon)
   - [Module 9: Testing](docs/MODULE_09_TESTING.md) (Coming soon)
   - [Module 10: Deployment](docs/MODULE_10_DEPLOYMENT.md) (Coming soon)

## Project Status

🚧 **Planning Phase** - Documentation and architecture design

## Key Features

- **Task Scheduling**: Complex cron-like scheduling with timezone support
- **External Execution**: Tasks executed by external systems, not by Cron Observer
- **Status Tracking**: Real-time execution status updates via SDK
- **Log Management**: Append-only logs with timestamps and levels
- **Execution History**: Complete history with date-based navigation
- **UUID-Based**: Tasks and executions use UUIDs for external reference

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

- **MASTER_PLAN.md**: Complete project overview, goals, requirements
- **docs/MODULE_XX_*.md**: Phase-by-phase implementation guides
- Each module is self-contained and can be implemented independently

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

**Note**: This project is in early planning phase. Implementation will begin after completing all module documentation.

