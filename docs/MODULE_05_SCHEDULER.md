# Module 5: Scheduler Engine

## Overview

This module implements the scheduler engine that evaluates cron expressions, handles timezone-aware scheduling, applies exclusion rules, and creates execution records with PENDING status.

## Scheduler Architecture

### Component Structure

```
┌─────────────────────────────────────┐
│      Scheduler Service              │
├─────────────────────────────────────┤
│  ┌──────────────────────────────┐  │
│  │   Cron Evaluator             │  │
│  │   - Parse cron expressions   │  │
│  │   - Calculate next run time  │  │
│  └──────────────────────────────┘  │
│  ┌──────────────────────────────┐  │
│  │   Timezone Handler            │  │
│  │   - Convert timezones         │  │
│  │   - Handle DST                │  │
│  └──────────────────────────────┘  │
│  ┌──────────────────────────────┐  │
│  │   Exclusion Evaluator         │  │
│  │   - Check holidays            │  │
│  │   - Apply exclusions          │  │
│  └──────────────────────────────┘  │
│  ┌──────────────────────────────┐  │
│  │   Execution Creator           │  │
│  │   - Create execution records  │  │
│  │   - Generate UUIDs            │  │
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
```

## Scheduler Service Design

### Service Responsibilities

1. **Poll Active Tasks**: Periodically check active tasks
2. **Evaluate Schedules**: Determine if task should run now
3. **Check Exclusions**: Verify no exclusions apply
4. **Create Executions**: Create execution records with PENDING status
5. **Handle One-off Tasks**: Process one-time scheduled tasks

### Execution Flow

```
1. Scheduler wakes up (every minute or configurable interval)
   ↓
2. Fetch all ACTIVE tasks
   ↓
3. For each task:
   a. Evaluate cron expression (if RECURRING)
   b. Check if current time matches schedule
   c. Check timezone
   d. Check time range (if specified)
   e. Check exclusions
   f. If all checks pass → Create execution record
   ↓
4. Sleep until next evaluation cycle
```

## Cron Expression Evaluation

### Supported Cron Format

Standard 5-field cron expression:
```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │
* * * * *
```

### Examples

- `0 10 * * *` - Every day at 10:00 AM
- `0 10 * * 0-4` - Every day Sunday-Thursday at 10:00 AM
- `0 12,18 * * *` - Every day at 12:00 PM and 6:00 PM
- `0 */2 * * *` - Every 2 hours
- `0 10-14 * * 0-4` - Every hour from 10 AM to 2 PM, Sunday-Thursday

### Cron Evaluation Library

**Recommended Libraries**:
- **Java**: `cron-utils` or Quartz CronExpression
- **Node.js**: `node-cron` or `cron-parser`
- **Python**: `croniter` or `schedule`
- **Go**: `robfig/cron/v3`

### Implementation Example (Pseudocode)

```java
public class CronEvaluator {
    public boolean shouldRunNow(String cronExpression, 
                                String timezone, 
                                LocalDateTime now) {
        // Parse cron expression
        CronExpression cron = CronExpression.parse(cronExpression);
        
        // Convert now to task timezone
        ZonedDateTime zonedNow = now.atZone(ZoneId.of(timezone));
        
        // Get next execution time
        ZonedDateTime nextRun = cron.nextTimeAfter(zonedNow.minusMinutes(1));
        
        // Check if next run is within current minute
        return nextRun != null && 
               nextRun.isBefore(zonedNow.plusMinutes(1)) &&
               nextRun.isAfter(zonedNow.minusSeconds(30));
    }
}
```

## Timezone Handling

### Timezone Requirements

- All scheduling must be timezone-aware
- Support IANA timezone identifiers (e.g., "America/New_York")
- Handle Daylight Saving Time (DST) transitions
- Convert system time to task timezone for evaluation

### Timezone Conversion

```java
public class TimezoneHandler {
    public ZonedDateTime convertToTaskTimezone(LocalDateTime systemTime, 
                                               String taskTimezone) {
        ZoneId taskZone = ZoneId.of(taskTimezone);
        ZoneId systemZone = ZoneId.systemDefault();
        
        return systemTime
            .atZone(systemZone)
            .withZoneSameInstant(taskZone);
    }
    
    public boolean isInTimeRange(ZonedDateTime time, 
                                 String startTime, 
                                 String endTime) {
        LocalTime current = time.toLocalTime();
        LocalTime start = LocalTime.parse(startTime);
        LocalTime end = LocalTime.parse(endTime);
        
        return !current.isBefore(start) && !current.isAfter(end);
    }
}
```

### DST Considerations

- Cron expressions should evaluate correctly across DST boundaries
- Execution times may shift by 1 hour during DST transitions
- Document this behavior for users

## Time Range Evaluation

### Time Range Support

For recurring tasks, support execution windows:
```json
{
  "time_range": {
    "start": "10:00",
    "end": "14:30"
  }
}
```

### Implementation

```java
public boolean isWithinTimeRange(ZonedDateTime taskTime, 
                                 ScheduleConfig config) {
    if (config.getTimeRange() == null) {
        return true; // No time range restriction
    }
    
    LocalTime current = taskTime.toLocalTime();
    LocalTime start = LocalTime.parse(config.getTimeRange().getStart());
    LocalTime end = LocalTime.parse(config.getTimeRange().getEnd());
    
    return !current.isBefore(start) && !current.isAfter(end);
}
```

## Day of Week Filtering

### Day of Week Support

```json
{
  "days_of_week": [0, 1, 2, 3, 4]  // Sunday-Thursday
}
```

### Implementation

```java
public boolean isAllowedDayOfWeek(ZonedDateTime taskTime, 
                                  ScheduleConfig config) {
    if (config.getDaysOfWeek() == null || 
        config.getDaysOfWeek().isEmpty()) {
        return true; // No day restriction
    }
    
    int dayOfWeek = taskTime.getDayOfWeek().getValue() % 7; // 0=Sunday
    return config.getDaysOfWeek().contains(dayOfWeek);
}
```

## Exclusion Evaluation

### Exclusion Types

1. **Single Date**: Specific date exclusion
2. **Date Range**: Range of dates excluded
3. **Recurring Yearly**: Yearly recurring exclusion (e.g., holidays)

### Exclusion Check

```java
public class ExclusionEvaluator {
    public boolean isExcluded(LocalDate date, 
                              Task task, 
                              List<Exclusion> exclusions) {
        // Check global exclusions
        List<Exclusion> globalExclusions = exclusions.stream()
            .filter(e -> e.isAppliesToAllTasks())
            .collect(Collectors.toList());
        
        if (isExcludedByList(date, globalExclusions)) {
            return true;
        }
        
        // Check task-specific exclusions
        List<Exclusion> taskExclusions = task.getExclusions();
        return isExcludedByList(date, taskExclusions);
    }
    
    private boolean isExcludedByList(LocalDate date, 
                                     List<Exclusion> exclusions) {
        for (Exclusion exclusion : exclusions) {
            switch (exclusion.getType()) {
                case SINGLE_DATE:
                    if (date.equals(exclusion.getStartDate())) {
                        return true;
                    }
                    break;
                case DATE_RANGE:
                    if (!date.isBefore(exclusion.getStartDate()) &&
                        !date.isAfter(exclusion.getEndDate())) {
                        return true;
                    }
                    break;
                case RECURRING_YEARLY:
                    if (isYearlyMatch(date, exclusion)) {
                        return true;
                    }
                    break;
            }
        }
        return false;
    }
    
    private boolean isYearlyMatch(LocalDate date, Exclusion exclusion) {
        RecurringPattern pattern = exclusion.getRecurringPattern();
        return date.getMonthValue() == pattern.getMonth() &&
               date.getDayOfMonth() == pattern.getDay();
    }
}
```

## Execution Record Creation

### Creating Execution Records

When scheduler determines a task should run:

```java
public class ExecutionCreator {
    public Execution createExecution(Task task, 
                                    ZonedDateTime scheduledTime) {
        Execution execution = new Execution();
        execution.setExecutionUuid(UUID.randomUUID().toString());
        execution.setTaskId(task.getId());
        execution.setTaskUuid(task.getUuid());
        execution.setScheduledAt(scheduledTime.toInstant());
        execution.setStatus(ExecutionStatus.PENDING);
        execution.setTriggerType(TriggerType.SCHEDULED);
        execution.setCreatedAt(Instant.now());
        execution.setUpdatedAt(Instant.now());
        
        return executionRepository.save(execution);
    }
}
```

### UUID Generation

- Use standard UUID v4 (random UUID)
- Ensure uniqueness (database constraint)
- Store as VARCHAR(36) in database

## One-off Task Handling

### One-off Task Evaluation

```java
public boolean shouldRunOneOffTask(Task task, ZonedDateTime now) {
    if (task.getScheduleType() != ScheduleType.ONEOFF) {
        return false;
    }
    
    ScheduleConfig config = task.getScheduleConfig();
    ZonedDateTime executeAt = ZonedDateTime.parse(config.getExecuteAt())
        .withZoneSameInstant(ZoneId.of(config.getTimezone()));
    
    // Check if execute_at is within current minute
    return executeAt.isBefore(now.plusMinutes(1)) &&
           executeAt.isAfter(now.minusSeconds(30)) &&
           !executionExistsForOneOff(task);
}

private boolean executionExistsForOneOff(Task task) {
    // Check if execution already created for this one-off task
    return executionRepository.existsByTaskIdAndTriggerType(
        task.getId(), 
        TriggerType.SCHEDULED
    );
}
```

## Scheduler Service Implementation

### Main Scheduler Loop

```java
@Service
public class SchedulerService {
    private static final int POLL_INTERVAL_SECONDS = 60;
    
    @Scheduled(fixedDelay = POLL_INTERVAL_SECONDS * 1000)
    public void evaluateSchedules() {
        List<Task> activeTasks = taskRepository.findByStatus(TaskStatus.ACTIVE);
        
        ZonedDateTime now = ZonedDateTime.now();
        
        for (Task task : activeTasks) {
            try {
                if (shouldExecuteTask(task, now)) {
                    createExecution(task, now);
                }
            } catch (Exception e) {
                log.error("Error evaluating task: " + task.getUuid(), e);
            }
        }
    }
    
    private boolean shouldExecuteTask(Task task, ZonedDateTime now) {
        // Convert to task timezone
        ZonedDateTime taskTime = timezoneHandler.convertToTaskTimezone(
            now.toLocalDateTime(), 
            task.getScheduleConfig().getTimezone()
        );
        
        if (task.getScheduleType() == ScheduleType.RECURRING) {
            return evaluateRecurringTask(task, taskTime);
        } else {
            return evaluateOneOffTask(task, taskTime);
        }
    }
    
    private boolean evaluateRecurringTask(Task task, ZonedDateTime taskTime) {
        // Check cron expression
        if (!cronEvaluator.shouldRunNow(
            task.getScheduleConfig().getCronExpression(),
            task.getScheduleConfig().getTimezone(),
            taskTime.toLocalDateTime()
        )) {
            return false;
        }
        
        // Check time range
        if (!timezoneHandler.isInTimeRange(
            taskTime,
            task.getScheduleConfig().getTimeRange()
        )) {
            return false;
        }
        
        // Check day of week
        if (!isAllowedDayOfWeek(taskTime, task.getScheduleConfig())) {
            return false;
        }
        
        // Check exclusions
        LocalDate date = taskTime.toLocalDate();
        if (exclusionEvaluator.isExcluded(date, task, getAllExclusions())) {
            return false;
        }
        
        // Check if execution already created for this schedule
        if (executionAlreadyCreated(task, taskTime)) {
            return false;
        }
        
        return true;
    }
}
```

## Scheduler Configuration

### Configuration Properties

```properties
# Scheduler configuration
scheduler.enabled=true
scheduler.poll-interval-seconds=60
scheduler.max-concurrent-evaluations=100
scheduler.timezone=UTC
```

### Scheduler Threading

- **Single Thread**: Simple, sequential evaluation
- **Thread Pool**: Parallel evaluation for many tasks
- **Async Processing**: Non-blocking evaluation

**Recommendation**: Start with single thread, optimize later if needed.

## Error Handling

### Scheduler Errors

- **Task Evaluation Errors**: Log and continue with next task
- **Database Errors**: Retry with exponential backoff
- **Invalid Cron Expressions**: Log error, mark task as invalid
- **Timezone Errors**: Use UTC as fallback, log warning

### Error Logging

```java
try {
    if (shouldExecuteTask(task, now)) {
        createExecution(task, now);
    }
} catch (InvalidCronExpressionException e) {
    log.error("Invalid cron expression for task: " + task.getUuid(), e);
    // Optionally disable task
} catch (DatabaseException e) {
    log.error("Database error creating execution for task: " + task.getUuid(), e);
    // Retry logic
} catch (Exception e) {
    log.error("Unexpected error evaluating task: " + task.getUuid(), e);
}
```

## Performance Considerations

### Optimization Strategies

1. **Batch Queries**: Fetch all active tasks in one query
2. **Caching**: Cache exclusion lists
3. **Indexing**: Ensure database indexes on status, schedule_type
4. **Lazy Loading**: Load task exclusions only when needed
5. **Parallel Processing**: Evaluate multiple tasks concurrently

### Monitoring

- Track scheduler execution time
- Monitor number of executions created
- Alert on scheduler failures
- Track task evaluation performance

## Testing

### Unit Tests

- Test cron expression evaluation
- Test timezone conversions
- Test exclusion logic
- Test time range filtering
- Test day of week filtering

### Integration Tests

- Test scheduler with database
- Test execution creation
- Test exclusion application
- Test one-off task handling

## Next Steps

After completing this module:
1. Implement cron evaluator
2. Implement timezone handler
3. Implement exclusion evaluator
4. Implement scheduler service
5. Add error handling and logging
6. Proceed to Module 6: SDK/API for External Systems

