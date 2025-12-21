# Module 9: Testing Strategy

## Overview

This module defines the comprehensive testing strategy for Cron Observer, including unit tests, integration tests, API tests, and end-to-end tests.

## Testing Philosophy

### Testing Pyramid

```
        /\
       /  \
      / E2E \          (Few, critical user flows)
     /───────\
    /         \
   / Integration \    (API endpoints, service integration)
  /───────────────\
 /                 \
/   Unit Tests      \  (Many, fast, isolated)
/─────────────────────\
```

### Principles

1. **Unit Tests**: Fast, isolated, test single functions/classes
2. **Integration Tests**: Test component interactions, database operations
3. **API Tests**: Test HTTP endpoints, request/response handling
4. **E2E Tests**: Test complete user workflows

## Unit Testing

### Backend Unit Tests

#### Test Structure

```
backend/src/test/java/com/cronobserver/
├── services/
│   ├── TaskServiceTest.java
│   ├── SchedulerServiceTest.java
│   └── ExecutionTrackingServiceTest.java
├── utils/
│   ├── CronEvaluatorTest.java
│   └── StatusTransitionValidatorTest.java
└── controllers/
    └── TaskControllerTest.java
```

#### Example: Cron Evaluator Test

```java
@ExtendWith(MockitoExtension.class)
class CronEvaluatorTest {
    
    private CronEvaluator cronEvaluator;
    
    @BeforeEach
    void setUp() {
        cronEvaluator = new CronEvaluator();
    }
    
    @Test
    void testDailyAt10AM() {
        String cron = "0 10 * * *";
        LocalDateTime testTime = LocalDateTime.of(2025, 1, 15, 10, 0);
        
        assertTrue(cronEvaluator.shouldRunNow(cron, "UTC", testTime));
    }
    
    @Test
    void testWeekdayOnly() {
        String cron = "0 10 * * 1-5"; // Monday-Friday
        LocalDateTime monday = LocalDateTime.of(2025, 1, 13, 10, 0); // Monday
        LocalDateTime saturday = LocalDateTime.of(2025, 1, 18, 10, 0); // Saturday
        
        assertTrue(cronEvaluator.shouldRunNow(cron, "UTC", monday));
        assertFalse(cronEvaluator.shouldRunNow(cron, "UTC", saturday));
    }
    
    @Test
    void testInvalidCronExpression() {
        String invalidCron = "invalid";
        
        assertThrows(InvalidCronExpressionException.class, () -> {
            cronEvaluator.shouldRunNow(invalidCron, "UTC", LocalDateTime.now());
        });
    }
}
```

#### Example: Status Transition Validator Test

```java
@ExtendWith(MockitoExtension.class)
class StatusTransitionValidatorTest {
    
    private StatusTransitionValidator validator;
    
    @BeforeEach
    void setUp() {
        validator = new StatusTransitionValidator();
    }
    
    @Test
    void testValidTransitions() {
        assertTrue(validator.isValidTransition(
            ExecutionStatus.PENDING, 
            ExecutionStatus.RUNNING
        ));
        
        assertTrue(validator.isValidTransition(
            ExecutionStatus.RUNNING, 
            ExecutionStatus.FINISHED
        ));
    }
    
    @Test
    void testInvalidTransitions() {
        assertFalse(validator.isValidTransition(
            ExecutionStatus.FINISHED, 
            ExecutionStatus.RUNNING
        ));
        
        assertFalse(validator.isValidTransition(
            ExecutionStatus.PENDING, 
            ExecutionStatus.FINISHED
        ));
    }
    
    @Test
    void testTerminalStates() {
        assertFalse(validator.isValidTransition(
            ExecutionStatus.FINISHED, 
            ExecutionStatus.RUNNING
        ));
        
        assertFalse(validator.isValidTransition(
            ExecutionStatus.FAILED, 
            ExecutionStatus.RUNNING
        ));
    }
}
```

### Frontend Unit Tests

#### Test Structure

```
frontend/src/
├── components/
│   ├── TaskCard.test.tsx
│   └── ExecutionLogs.test.tsx
├── services/
│   └── taskService.test.ts
└── utils/
    └── dateUtils.test.ts
```

#### Example: TaskCard Component Test

```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { TaskCard } from './TaskCard';

describe('TaskCard', () => {
  const mockTask = {
    uuid: '123',
    name: 'Test Task',
    status: 'ACTIVE',
    schedule_type: 'RECURRING'
  };
  
  it('renders task name', () => {
    render(<TaskCard task={mockTask} />);
    expect(screen.getByText('Test Task')).toBeInTheDocument();
  });
  
  it('calls onPause when pause button clicked', () => {
    const onPause = jest.fn();
    render(<TaskCard task={mockTask} onPause={onPause} />);
    
    fireEvent.click(screen.getByText('Pause'));
    expect(onPause).toHaveBeenCalledWith('123');
  });
  
  it('shows paused state correctly', () => {
    const pausedTask = { ...mockTask, status: 'PAUSED' };
    render(<TaskCard task={pausedTask} />);
    expect(screen.getByText('PAUSED')).toBeInTheDocument();
  });
});
```

## Integration Testing

### Backend Integration Tests

#### Test Database Setup

```java
@SpringBootTest
@Testcontainers
@Transactional
class TaskServiceIntegrationTest {
    
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:14")
            .withDatabaseName("testdb")
            .withUsername("test")
            .withPassword("test");
    
    @Autowired
    private TaskService taskService;
    
    @Autowired
    private TaskRepository taskRepository;
    
    @Test
    void testCreateTask() {
        CreateTaskRequest request = new CreateTaskRequest();
        request.setName("Test Task");
        request.setScheduleType(ScheduleType.RECURRING);
        // ... set other fields
        
        Task task = taskService.createTask(request);
        
        assertNotNull(task.getId());
        assertNotNull(task.getUuid());
        assertEquals("Test Task", task.getName());
        
        // Verify in database
        Optional<Task> saved = taskRepository.findByUuid(task.getUuid());
        assertTrue(saved.isPresent());
    }
    
    @Test
    void testPauseTask() {
        Task task = createTestTask();
        taskService.pauseTask(task.getUuid());
        
        Task updated = taskRepository.findByUuid(task.getUuid()).orElseThrow();
        assertEquals(TaskStatus.PAUSED, updated.getStatus());
    }
}
```

#### Scheduler Integration Test

```java
@SpringBootTest
@Testcontainers
class SchedulerServiceIntegrationTest {
    
    @Autowired
    private SchedulerService schedulerService;
    
    @Autowired
    private TaskRepository taskRepository;
    
    @Autowired
    private ExecutionRepository executionRepository;
    
    @Test
    void testCreateExecutionForScheduledTask() {
        // Create a task scheduled to run now
        Task task = createTaskWithCron("0 * * * *"); // Every hour
        
        // Mock current time to match schedule
        LocalDateTime now = LocalDateTime.of(2025, 1, 15, 10, 0);
        
        schedulerService.evaluateSchedules(now);
        
        // Verify execution created
        List<Execution> executions = executionRepository
            .findByTaskUuid(task.getUuid());
        
        assertEquals(1, executions.size());
        assertEquals(ExecutionStatus.PENDING, executions.get(0).getStatus());
    }
}
```

### API Integration Tests

```java
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
class TaskControllerIntegrationTest {
    
    @Autowired
    private MockMvc mockMvc;
    
    @Autowired
    private ObjectMapper objectMapper;
    
    @Test
    void testCreateTask() throws Exception {
        CreateTaskRequest request = new CreateTaskRequest();
        request.setName("API Test Task");
        // ... set fields
        
        mockMvc.perform(post("/api/v1/tasks")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
            .andExpect(status().isCreated())
            .andExpect(jsonPath("$.data.name").value("API Test Task"))
            .andExpect(jsonPath("$.data.uuid").exists());
    }
    
    @Test
    void testGetTask() throws Exception {
        Task task = createTestTask();
        
        mockMvc.perform(get("/api/v1/tasks/" + task.getUuid()))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.data.uuid").value(task.getUuid()))
            .andExpect(jsonPath("$.data.name").value(task.getName()));
    }
    
    @Test
    void testUpdateExecutionStatus() throws Exception {
        Execution execution = createPendingExecution();
        
        StatusUpdateRequest request = new StatusUpdateRequest();
        request.setStatus("RUNNING");
        
        mockMvc.perform(put("/api/v1/sdk/executions/" + execution.getExecutionUuid() + "/status")
                .header("X-API-Key", "test-api-key")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.data.status").value("RUNNING"));
    }
}
```

## End-to-End Testing

### E2E Test Scenarios

#### Scenario 1: Complete Task Lifecycle

```java
@SpringBootTest
@Testcontainers
class TaskLifecycleE2ETest {
    
    @Test
    void testCompleteTaskExecutionFlow() {
        // 1. Create task
        Task task = taskService.createTask(createTaskRequest());
        assertNotNull(task.getUuid());
        
        // 2. Trigger execution manually
        Execution execution = taskService.triggerTask(task.getUuid());
        assertEquals(ExecutionStatus.PENDING, execution.getStatus());
        
        // 3. External system updates status to RUNNING
        executionService.updateStatus(
            execution.getExecutionUuid(), 
            ExecutionStatus.RUNNING
        );
        
        // 4. External system appends logs
        executionService.appendLogs(execution.getExecutionUuid(), Arrays.asList(
            new LogEntry("INFO", "Task started"),
            new LogEntry("INFO", "Processing data")
        ));
        
        // 5. External system updates status to FINISHED
        executionService.updateStatus(
            execution.getExecutionUuid(), 
            ExecutionStatus.FINISHED,
            Map.of("records_processed", 150)
        );
        
        // 6. Verify execution is complete
        Execution completed = executionService.getExecution(execution.getExecutionUuid());
        assertEquals(ExecutionStatus.FINISHED, completed.getStatus());
        assertNotNull(completed.getCompletedAt());
        assertNotNull(completed.getDurationMs());
        
        // 7. Verify logs
        List<ExecutionLog> logs = executionService.getLogs(execution.getExecutionUuid());
        assertEquals(2, logs.size());
    }
}
```

### Frontend E2E Tests (Playwright/Cypress)

```typescript
// e2e/task-management.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Task Management', () => {
  test('create and view task', async ({ page }) => {
    // Navigate to tasks page
    await page.goto('/tasks');
    
    // Click create task button
    await page.click('text=Create Task');
    
    // Fill form
    await page.fill('input[name="name"]', 'E2E Test Task');
    await page.selectOption('select[name="schedule_type"]', 'RECURRING');
    await page.fill('input[name="cron_expression"]', '0 10 * * *');
    
    // Submit
    await page.click('button[type="submit"]');
    
    // Verify task appears in list
    await expect(page.locator('text=E2E Test Task')).toBeVisible();
  });
  
  test('view execution history with date navigation', async ({ page }) => {
    await page.goto('/executions');
    
    // Click on a date in left panel
    await page.click('text=2025-01-15');
    
    // Verify executions for that date are shown
    await expect(page.locator('.execution-card')).toHaveCount(5);
    
    // Click on an execution
    await page.click('.execution-card:first-child');
    
    // Verify execution detail page
    await expect(page.locator('text=Execution Detail')).toBeVisible();
    
    // Expand logs
    await page.click('text=Logs');
    
    // Verify logs are visible
    await expect(page.locator('.log-entry')).toHaveCount(2);
  });
});
```

## Test Data Management

### Test Fixtures

```java
public class TestFixtures {
    
    public static Task createTestTask() {
        Task task = new Task();
        task.setUuid(UUID.randomUUID().toString());
        task.setName("Test Task");
        task.setScheduleType(ScheduleType.RECURRING);
        task.setStatus(TaskStatus.ACTIVE);
        // ... set other fields
        return task;
    }
    
    public static Execution createPendingExecution(Task task) {
        Execution execution = new Execution();
        execution.setExecutionUuid(UUID.randomUUID().toString());
        execution.setTaskId(task.getId());
        execution.setTaskUuid(task.getUuid());
        execution.setStatus(ExecutionStatus.PENDING);
        execution.setScheduledAt(Instant.now());
        return execution;
    }
}
```

### Database Seeding

```java
@Component
public class TestDataSeeder {
    
    @Autowired
    private TaskRepository taskRepository;
    
    public void seedTestData() {
        // Create test tasks
        for (int i = 0; i < 10; i++) {
            Task task = TestFixtures.createTestTask();
            task.setName("Test Task " + i);
            taskRepository.save(task);
        }
    }
}
```

## Performance Testing

### Load Testing

```java
@Test
void testConcurrentStatusUpdates() throws InterruptedException {
    Execution execution = createPendingExecution();
    int threadCount = 100;
    ExecutorService executor = Executors.newFixedThreadPool(threadCount);
    CountDownLatch latch = new CountDownLatch(threadCount);
    
    for (int i = 0; i < threadCount; i++) {
        executor.submit(() -> {
            try {
                executionService.updateStatus(
                    execution.getExecutionUuid(),
                    ExecutionStatus.RUNNING
                );
            } finally {
                latch.countDown();
            }
        });
    }
    
    latch.await(10, TimeUnit.SECONDS);
    
    // Verify only one update succeeded (idempotent or first wins)
    Execution updated = executionService.getExecution(execution.getExecutionUuid());
    assertEquals(ExecutionStatus.RUNNING, updated.getStatus());
}
```

## Test Coverage Goals

### Coverage Targets

- **Unit Tests**: 80%+ code coverage
- **Integration Tests**: All critical paths
- **API Tests**: All endpoints
- **E2E Tests**: Main user workflows

### Coverage Tools

- **Backend**: JaCoCo (Java) / Coverage.py (Python) / Istanbul (Node.js)
- **Frontend**: Jest Coverage / Vitest Coverage

## Continuous Integration

### CI Pipeline

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up JDK
        uses: actions/setup-java@v3
        with:
          java-version: '17'
      - name: Run tests
        run: ./gradlew test
      - name: Generate coverage
        run: ./gradlew jacocoTestReport
  
  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: npm ci
      - name: Run tests
        run: npm test
      - name: Run E2E tests
        run: npm run test:e2e
```

## Mocking Strategy

### External Dependencies

- **Database**: Use Testcontainers for real database testing
- **Time**: Use time mocking libraries (e.g., `java.time` test utilities)
- **External APIs**: Use WireMock or MockServer

### Example: Time Mocking

```java
@Test
void testSchedulerWithMockedTime() {
    // Mock current time
    LocalDateTime mockedTime = LocalDateTime.of(2025, 1, 15, 10, 0);
    
    // Use time provider
    TimeProvider timeProvider = Mockito.mock(TimeProvider.class);
    when(timeProvider.now()).thenReturn(mockedTime);
    
    schedulerService.setTimeProvider(timeProvider);
    
    // Test scheduler behavior at specific time
    schedulerService.evaluateSchedules();
    
    // Verify executions created
}
```

## Next Steps

After completing this module:
1. Set up testing frameworks
2. Write unit tests for core components
3. Write integration tests for services
4. Write API tests for endpoints
5. Write E2E tests for critical flows
6. Set up CI/CD pipeline
7. Proceed to Module 10: Deployment & DevOps

