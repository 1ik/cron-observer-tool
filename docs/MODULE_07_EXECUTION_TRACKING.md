# Module 7: Execution Tracking

## Overview

This module implements the execution tracking engine that manages execution lifecycle, validates status transitions, calculates statistics, and provides execution history queries.

## Execution Tracking Architecture

### Component Structure

```
┌─────────────────────────────────────┐
│   Execution Tracking Service         │
├─────────────────────────────────────┤
│  ┌──────────────────────────────┐  │
│  │   Status Transition Manager  │  │
│  │   - Validate transitions     │  │
│  │   - Update timestamps        │  │
│  └──────────────────────────────┘  │
│  ┌──────────────────────────────┐  │
│  │   Statistics Calculator      │  │
│  │   - Success rates             │  │
│  │   - Average durations        │  │
│  └──────────────────────────────┘  │
│  ┌──────────────────────────────┐  │
│  │   History Query Service      │  │
│  │   - Date-based queries        │  │
│  │   - Filtering & pagination   │  │
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
```

## Status Transition Management

### Valid Status Transitions

```java
public class StatusTransitionValidator {
    private static final Map<ExecutionStatus, Set<ExecutionStatus>> VALID_TRANSITIONS = Map.of(
        ExecutionStatus.PENDING, Set.of(
            ExecutionStatus.RUNNING,
            ExecutionStatus.CANCELLED
        ),
        ExecutionStatus.RUNNING, Set.of(
            ExecutionStatus.FINISHED,
            ExecutionStatus.FAILED,
            ExecutionStatus.CANCELLED
        ),
        ExecutionStatus.FINISHED, Set.of(), // Terminal
        ExecutionStatus.FAILED, Set.of(),   // Terminal
        ExecutionStatus.CANCELLED, Set.of()  // Terminal
    );
    
    public boolean isValidTransition(ExecutionStatus from, ExecutionStatus to) {
        Set<ExecutionStatus> allowed = VALID_TRANSITIONS.get(from);
        return allowed != null && allowed.contains(to);
    }
    
    public void validateTransition(Execution execution, ExecutionStatus newStatus) {
        ExecutionStatus currentStatus = execution.getStatus();
        
        if (!isValidTransition(currentStatus, newStatus)) {
            throw new InvalidStatusTransitionException(
                String.format(
                    "Cannot transition from %s to %s. Allowed transitions: %s",
                    currentStatus,
                    newStatus,
                    VALID_TRANSITIONS.get(currentStatus)
                )
            );
        }
    }
}
```

### Status Update Logic

```java
@Service
public class ExecutionTrackingService {
    
    public Execution updateStatus(Execution execution, ExecutionStatus newStatus) {
        // Validate transition
        statusTransitionValidator.validateTransition(execution, newStatus);
        
        ExecutionStatus oldStatus = execution.getStatus();
        execution.setStatus(newStatus);
        Instant now = Instant.now();
        
        // Handle timestamps based on status
        switch (newStatus) {
            case RUNNING:
                if (execution.getStartedAt() == null) {
                    execution.setStartedAt(now);
                }
                break;
                
            case FINISHED:
            case FAILED:
                if (execution.getCompletedAt() == null) {
                    execution.setCompletedAt(now);
                    
                    // Calculate duration
                    if (execution.getStartedAt() != null) {
                        long durationMs = Duration.between(
                            execution.getStartedAt(),
                            execution.getCompletedAt()
                        ).toMillis();
                        execution.setDurationMs(durationMs);
                    }
                }
                break;
                
            case CANCELLED:
                // Can be cancelled at any time (PENDING or RUNNING)
                if (execution.getCompletedAt() == null) {
                    execution.setCompletedAt(now);
                }
                break;
        }
        
        execution.setUpdatedAt(now);
        
        return executionRepository.save(execution);
    }
}
```

## Execution Statistics

### Task Statistics Calculation

```java
@Service
public class ExecutionStatisticsService {
    
    public TaskStatistics calculateTaskStatistics(Long taskId) {
        List<Execution> executions = executionRepository.findByTaskId(taskId);
        
        if (executions.isEmpty()) {
            return TaskStatistics.empty();
        }
        
        long totalExecutions = executions.size();
        long finishedCount = executions.stream()
            .filter(e -> e.getStatus() == ExecutionStatus.FINISHED)
            .count();
        long failedCount = executions.stream()
            .filter(e -> e.getStatus() == ExecutionStatus.FAILED)
            .count();
        
        double successRate = totalExecutions > 0 
            ? (double) finishedCount / totalExecutions 
            : 0.0;
        
        OptionalDouble avgDuration = executions.stream()
            .filter(e -> e.getDurationMs() != null)
            .mapToLong(Execution::getDurationMs)
            .average();
        
        Optional<Execution> lastExecution = executions.stream()
            .max(Comparator.comparing(Execution::getScheduledAt));
        
        return TaskStatistics.builder()
            .totalExecutions(totalExecutions)
            .finishedCount(finishedCount)
            .failedCount(failedCount)
            .successRate(successRate)
            .averageDurationMs(avgDuration.orElse(0.0))
            .lastExecution(lastExecution.map(Execution::getScheduledAt).orElse(null))
            .build();
    }
}
```

### Statistics Model

```java
public class TaskStatistics {
    private long totalExecutions;
    private long finishedCount;
    private long failedCount;
    private double successRate; // 0.0 to 1.0
    private double averageDurationMs;
    private Instant lastExecution;
    private Instant nextExecution; // Calculated from schedule
    
    // Getters and setters
}
```

## Execution History Queries

### Date-Based Queries

```java
@Service
public class ExecutionHistoryService {
    
    public List<ExecutionDateSummary> getExecutionDates(
            LocalDate startDate, 
            LocalDate endDate, 
            String taskUuid) {
        
        List<Execution> executions;
        if (taskUuid != null) {
            executions = executionRepository.findByTaskUuidAndDateRange(
                taskUuid,
                startDate.atStartOfDay(),
                endDate.atTime(23, 59, 59)
            );
        } else {
            executions = executionRepository.findByDateRange(
                startDate.atStartOfDay(),
                endDate.atTime(23, 59, 59)
            );
        }
        
        // Group by date
        Map<LocalDate, List<Execution>> byDate = executions.stream()
            .collect(Collectors.groupingBy(
                e -> e.getScheduledAt().atZone(ZoneId.systemDefault()).toLocalDate()
            ));
        
        // Build summaries
        return byDate.entrySet().stream()
            .map(entry -> {
                LocalDate date = entry.getKey();
                List<Execution> dateExecutions = entry.getValue();
                
                Map<ExecutionStatus, Long> statusCounts = dateExecutions.stream()
                    .collect(Collectors.groupingBy(
                        Execution::getStatus,
                        Collectors.counting()
                    ));
                
                return ExecutionDateSummary.builder()
                    .date(date)
                    .executionCount(dateExecutions.size())
                    .statusSummary(StatusSummary.builder()
                        .finished(statusCounts.getOrDefault(ExecutionStatus.FINISHED, 0L))
                        .failed(statusCounts.getOrDefault(ExecutionStatus.FAILED, 0L))
                        .pending(statusCounts.getOrDefault(ExecutionStatus.PENDING, 0L))
                        .running(statusCounts.getOrDefault(ExecutionStatus.RUNNING, 0L))
                        .cancelled(statusCounts.getOrDefault(ExecutionStatus.CANCELLED, 0L))
                        .build())
                    .build();
            })
            .sorted(Comparator.comparing(ExecutionDateSummary::getDate).reversed())
            .collect(Collectors.toList());
    }
    
    public List<Execution> getExecutionsByDate(
            LocalDate date, 
            String taskUuid, 
            ExecutionStatus status) {
        
        Instant startOfDay = date.atStartOfDay()
            .atZone(ZoneId.systemDefault())
            .toInstant();
        Instant endOfDay = date.atTime(23, 59, 59)
            .atZone(ZoneId.systemDefault())
            .toInstant();
        
        return executionRepository.findByDateRangeAndFilters(
            startOfDay,
            endOfDay,
            taskUuid,
            status
        );
    }
}
```

### Execution History Model

```java
public class ExecutionDateSummary {
    private LocalDate date;
    private long executionCount;
    private StatusSummary statusSummary;
    
    // Getters and setters
}

public class StatusSummary {
    private long finished;
    private long failed;
    private long pending;
    private long running;
    private long cancelled;
    
    // Getters and setters
}
```

## Log Management

### Log Retrieval

```java
@Service
public class ExecutionLogService {
    
    public List<ExecutionLog> getLogs(
            String executionUuid,
            LogLevel level,
            Instant startTimestamp,
            Instant endTimestamp,
            Integer limit) {
        
        Execution execution = executionRepository.findByExecutionUuid(executionUuid)
            .orElseThrow(() -> new ExecutionNotFoundException(executionUuid));
        
        Specification<ExecutionLog> spec = Specification.where(
            ExecutionLogSpecifications.byExecutionId(execution.getId())
        );
        
        if (level != null) {
            spec = spec.and(ExecutionLogSpecifications.byLevel(level));
        }
        
        if (startTimestamp != null) {
            spec = spec.and(ExecutionLogSpecifications.afterTimestamp(startTimestamp));
        }
        
        if (endTimestamp != null) {
            spec = spec.and(ExecutionLogSpecifications.beforeTimestamp(endTimestamp));
        }
        
        Pageable pageable = PageRequest.of(0, limit != null ? limit : 1000);
        
        return executionLogRepository.findAll(spec, pageable)
            .getContent();
    }
    
    public long getLogCount(String executionUuid) {
        Execution execution = executionRepository.findByExecutionUuid(executionUuid)
            .orElseThrow(() -> new ExecutionNotFoundException(executionUuid));
        
        return executionLogRepository.countByExecutionId(execution.getId());
    }
}
```

## Execution Lifecycle Events

### Event Handling

```java
@Service
public class ExecutionEventService {
    
    @EventListener
    public void handleExecutionStatusChanged(ExecutionStatusChangedEvent event) {
        Execution execution = event.getExecution();
        ExecutionStatus newStatus = event.getNewStatus();
        
        // Trigger notifications if needed
        if (newStatus == ExecutionStatus.FINISHED) {
            notificationService.sendSuccessNotification(execution);
        } else if (newStatus == ExecutionStatus.FAILED) {
            notificationService.sendFailureNotification(execution);
        }
        
        // Update task statistics cache
        taskStatisticsCache.evict(execution.getTaskId());
    }
    
    @EventListener
    public void handleExecutionCreated(ExecutionCreatedEvent event) {
        Execution execution = event.getExecution();
        
        // Log creation
        log.info("Execution created: {} for task: {}", 
            execution.getExecutionUuid(), 
            execution.getTaskUuid());
        
        // Update metrics
        metricsService.incrementExecutionCreated();
    }
}
```

## Query Optimization

### Database Queries

```java
@Repository
public interface ExecutionRepository extends JpaRepository<Execution, Long> {
    
    // Find pending executions for a task
    @Query("SELECT e FROM Execution e " +
           "WHERE e.taskUuid = :taskUuid " +
           "AND e.status = 'PENDING' " +
           "ORDER BY e.scheduledAt ASC")
    List<Execution> findPendingByTaskUuid(@Param("taskUuid") String taskUuid);
    
    // Find executions by date range
    @Query("SELECT e FROM Execution e " +
           "WHERE e.scheduledAt >= :start " +
           "AND e.scheduledAt <= :end " +
           "ORDER BY e.scheduledAt DESC")
    List<Execution> findByDateRange(
        @Param("start") Instant start,
        @Param("end") Instant end
    );
    
    // Find executions with filters
    @Query("SELECT e FROM Execution e " +
           "WHERE e.scheduledAt >= :start " +
           "AND e.scheduledAt <= :end " +
           "AND (:taskUuid IS NULL OR e.taskUuid = :taskUuid) " +
           "AND (:status IS NULL OR e.status = :status) " +
           "ORDER BY e.scheduledAt DESC")
    List<Execution> findByDateRangeAndFilters(
        @Param("start") Instant start,
        @Param("end") Instant end,
        @Param("taskUuid") String taskUuid,
        @Param("status") ExecutionStatus status
    );
    
    // Count executions by status for a task
    @Query("SELECT COUNT(e) FROM Execution e " +
           "WHERE e.taskId = :taskId " +
           "AND e.status = :status")
    long countByTaskIdAndStatus(
        @Param("taskId") Long taskId,
        @Param("status") ExecutionStatus status
    );
    
    // Check if execution exists for one-off task
    boolean existsByTaskIdAndTriggerType(Long taskId, TriggerType triggerType);
}
```

## Caching Strategy

### Statistics Caching

```java
@Service
public class CachedTaskStatisticsService {
    
    @Cacheable(value = "taskStatistics", key = "#taskId")
    public TaskStatistics getTaskStatistics(Long taskId) {
        return executionStatisticsService.calculateTaskStatistics(taskId);
    }
    
    @CacheEvict(value = "taskStatistics", key = "#taskId")
    public void evictTaskStatistics(Long taskId) {
        // Cache evicted when execution status changes
    }
}
```

## Monitoring & Metrics

### Execution Metrics

```java
@Service
public class ExecutionMetricsService {
    
    private final MeterRegistry meterRegistry;
    
    public void incrementExecutionCreated() {
        meterRegistry.counter("executions.created").increment();
    }
    
    public void recordExecutionDuration(long durationMs, ExecutionStatus status) {
        meterRegistry.timer("executions.duration", "status", status.name())
            .record(durationMs, TimeUnit.MILLISECONDS);
    }
    
    public void recordStatusTransition(ExecutionStatus from, ExecutionStatus to) {
        meterRegistry.counter("executions.transitions", 
            "from", from.name(), 
            "to", to.name())
            .increment();
    }
}
```

## Data Retention

### Execution Cleanup

```java
@Service
public class ExecutionCleanupService {
    
    @Scheduled(cron = "0 0 2 * * ?") // Daily at 2 AM
    public void cleanupOldExecutions() {
        LocalDate retentionDate = LocalDate.now().minusDays(90);
        Instant cutoff = retentionDate.atStartOfDay()
            .atZone(ZoneId.systemDefault())
            .toInstant();
        
        // Archive old executions
        List<Execution> oldExecutions = executionRepository
            .findByScheduledAtBefore(cutoff);
        
        // Move to archive table or delete
        // Keep summary statistics
    }
}
```

## Testing

### Unit Tests

- Test status transition validation
- Test statistics calculation
- Test date-based queries
- Test log retrieval

### Integration Tests

- Test execution lifecycle end-to-end
- Test history queries with database
- Test statistics accuracy

## Next Steps

After completing this module:
1. Implement status transition validator
2. Implement statistics calculator
3. Implement history query service
4. Add caching for performance
5. Add monitoring and metrics
6. Proceed to Module 8: Frontend UI

