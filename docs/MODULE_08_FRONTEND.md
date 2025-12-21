# Module 8: Frontend UI

## Overview

This module defines the frontend user interface for Cron Observer, including task management, execution history, date-based navigation, and dashboard components.

## Frontend Architecture

### Technology Stack Options

#### Option A: React + TypeScript
- **Framework**: React 18+
- **Language**: TypeScript
- **State Management**: React Query / Zustand
- **UI Library**: Material-UI / Ant Design / Tailwind CSS
- **Routing**: React Router
- **HTTP Client**: Axios / Fetch API

#### Option B: Vue.js + TypeScript
- **Framework**: Vue 3
- **Language**: TypeScript
- **State Management**: Pinia
- **UI Library**: Vuetify / Element Plus
- **Routing**: Vue Router
- **HTTP Client**: Axios

#### Option C: Next.js (React)
- **Framework**: Next.js 14+
- **Language**: TypeScript
- **Features**: SSR, API routes
- **UI Library**: Tailwind CSS + shadcn/ui

**Recommendation**: React + TypeScript for ecosystem and community support.

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── common/
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Select.tsx
│   │   │   ├── Modal.tsx
│   │   │   └── LoadingSpinner.tsx
│   │   ├── task/
│   │   │   ├── TaskList.tsx
│   │   │   ├── TaskCard.tsx
│   │   │   ├── TaskForm.tsx
│   │   │   └── TaskDetail.tsx
│   │   ├── execution/
│   │   │   ├── ExecutionList.tsx
│   │   │   ├── ExecutionCard.tsx
│   │   │   ├── ExecutionDetail.tsx
│   │   │   └── ExecutionLogs.tsx
│   │   └── dashboard/
│   │       ├── Dashboard.tsx
│   │       ├── StatsCard.tsx
│   │       └── RecentExecutions.tsx
│   ├── pages/
│   │   ├── TasksPage.tsx
│   │   ├── TaskDetailPage.tsx
│   │   ├── ExecutionsPage.tsx
│   │   ├── ExecutionDetailPage.tsx
│   │   └── DashboardPage.tsx
│   ├── services/
│   │   ├── api.ts
│   │   ├── taskService.ts
│   │   ├── executionService.ts
│   │   └── dateNavigationService.ts
│   ├── hooks/
│   │   ├── useTasks.ts
│   │   ├── useExecutions.ts
│   │   └── useDateNavigation.ts
│   ├── utils/
│   │   ├── dateUtils.ts
│   │   ├── formatters.ts
│   │   └── validators.ts
│   ├── types/
│   │   ├── task.ts
│   │   ├── execution.ts
│   │   └── api.ts
│   ├── App.tsx
│   └── main.tsx
├── package.json
├── tsconfig.json
└── vite.config.ts (or webpack.config.js)
```

## UI Components

### 1. Dashboard Page

**Route**: `/`

**Components**:
- Overview statistics cards
- Recent executions feed
- Upcoming executions
- Task status summary

**Layout**:
```
┌─────────────────────────────────────────┐
│  Dashboard                                │
├─────────────────────────────────────────┤
│  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐│
│  │Total │  │Active│  │Paused│  │Failed││
│  │Tasks │  │      │  │      │  │      ││
│  └──────┘  └──────┘  └──────┘  └──────┘│
├─────────────────────────────────────────┤
│  Recent Executions                       │
│  ┌─────────────────────────────────────┐│
│  │ [Execution list]                     ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
```

### 2. Task List Page

**Route**: `/tasks`

**Components**:
- Task list table/cards
- Filter controls
- Search bar
- Create task button
- Bulk actions

**Features**:
- Filter by status, schedule type
- Search by name
- Sort by various columns
- Pagination

**Layout**:
```
┌─────────────────────────────────────────┐
│  Tasks                    [+ Create Task]│
├─────────────────────────────────────────┤
│  [Search] [Status Filter] [Type Filter] │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐│
│  │ Name      │ Type    │ Status │ ... ││
│  ├─────────────────────────────────────┤│
│  │ Task 1    │ Recurring│ Active │ ... ││
│  │ Task 2    │ One-off  │ Paused │ ... ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
```

### 3. Task Detail Page

**Route**: `/tasks/{task_uuid}`

**Components**:
- Task information
- Schedule configuration
- Execution history table
- Manual trigger button
- Pause/Resume toggle
- Task statistics

**Layout**:
```
┌─────────────────────────────────────────┐
│  Task: DSE News scrape                   │
│  [Pause] [Resume] [Trigger] [Edit]       │
├─────────────────────────────────────────┤
│  Information                             │
│  ┌─────────────────────────────────────┐│
│  │ Name: DSE News scrape                ││
│  │ UUID: 550e8400-...                   ││
│  │ Status: Active                       ││
│  │ Schedule: Every day 10:00 AM         ││
│  └─────────────────────────────────────┘│
├─────────────────────────────────────────┤
│  Statistics                              │
│  ┌─────────────────────────────────────┐│
│  │ Total: 150 | Success: 95% | Avg: 3s ││
│  └─────────────────────────────────────┘│
├─────────────────────────────────────────┤
│  Execution History                       │
│  ┌─────────────────────────────────────┐│
│  │ Date      │ Status │ Duration │ ... ││
│  ├─────────────────────────────────────┤│
│  │ 2025-01-15│Finished│ 3s       │ ... ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
```

### 4. Execution History Page with Date Navigation

**Route**: `/executions`

**Components**:
- Left panel: Date list
- Right panel: Execution list for selected date
- Execution cards with status

**Layout**:
```
┌─────────────────────────────────────────┐
│  Execution History                      │
├──────────┬─────────────────────────────┤
│          │  Executions for 2025-01-15  │
│ 2025-01-15│ ┌─────────────────────────┐│
│ (5 execs) │ │ DSE News scrape           ││
│          │ │ ✓ Finished - 3s           ││
│ 2025-01-14│ │ 10:01:10                  ││
│ (3 execs) │ └─────────────────────────┘│
│          │ ┌─────────────────────────┐│
│ 2025-01-13│ │ Another Task            ││
│ (2 execs) │ │ ✓ Finished - 2s         ││
│          │ │ 10:01:10                 ││
│ 2025-01-12│ └─────────────────────────┘│
│ (1 exec)  │                             │
└──────────┴─────────────────────────────┘
```

### 5. Execution Detail Page

**Route**: `/executions/{execution_uuid}`

**Components**:
- Execution information
- Status timeline
- Logs viewer (expandable)
- Result data display

**Layout** (Based on Image):
```
┌─────────────────────────────────────────┐
│  Execution: DSE News scrape             │
│  2025-01-15 10:01:10                    │
├─────────────────────────────────────────┤
│  Status: ✓ Finished (3s)                │
│  ┌─────────────────────────────────────┐│
│  │ Execution ID: 223sds-232-...        ││
│  │                                      ││
│  │ ▼ Logs                              ││
│  │ ┌─────────────────────────────────┐││
│  │ │ 2025-12-21 10:01:00 Task started │││
│  │ │ 2025-12-21 10:01:00 Fetch started│││
│  │ └─────────────────────────────────┘││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
```

## Component Specifications

### TaskList Component

```typescript
interface TaskListProps {
  tasks: Task[];
  onTaskClick: (task: Task) => void;
  onPause: (taskUuid: string) => void;
  onResume: (taskUuid: string) => void;
  onDelete: (taskUuid: string) => void;
}

const TaskList: React.FC<TaskListProps> = ({ tasks, ... }) => {
  return (
    <div className="task-list">
      {tasks.map(task => (
        <TaskCard
          key={task.uuid}
          task={task}
          onClick={() => onTaskClick(task)}
          onPause={() => onPause(task.uuid)}
          onResume={() => onResume(task.uuid)}
          onDelete={() => onDelete(task.uuid)}
        />
      ))}
    </div>
  );
};
```

### DateNavigation Component

```typescript
interface DateNavigationProps {
  dates: ExecutionDateSummary[];
  selectedDate: LocalDate | null;
  onDateSelect: (date: LocalDate) => void;
}

const DateNavigation: React.FC<DateNavigationProps> = ({ 
  dates, 
  selectedDate, 
  onDateSelect 
}) => {
  return (
    <div className="date-navigation">
      <h3>Execution Dates</h3>
      {dates.map(dateSummary => (
        <div
          key={dateSummary.date}
          className={`date-item ${selectedDate === dateSummary.date ? 'selected' : ''}`}
          onClick={() => onDateSelect(dateSummary.date)}
        >
          <div className="date">{formatDate(dateSummary.date)}</div>
          <div className="count">{dateSummary.executionCount} executions</div>
          <div className="summary">
            ✓ {dateSummary.statusSummary.finished} | 
            ✗ {dateSummary.statusSummary.failed}
          </div>
        </div>
      ))}
    </div>
  );
};
```

### ExecutionLogs Component

```typescript
interface ExecutionLogsProps {
  executionUuid: string;
  logs: ExecutionLog[];
}

const ExecutionLogs: React.FC<ExecutionLogsProps> = ({ executionUuid, logs }) => {
  const [expanded, setExpanded] = useState(false);
  const [filterLevel, setFilterLevel] = useState<LogLevel | null>(null);
  
  const filteredLogs = filterLevel 
    ? logs.filter(log => log.level === filterLevel)
    : logs;
  
  return (
    <div className="execution-logs">
      <div className="logs-header">
        <button onClick={() => setExpanded(!expanded)}>
          {expanded ? '▼' : '▶'} Logs ({logs.length})
        </button>
        <select onChange={(e) => setFilterLevel(e.target.value || null)}>
          <option value="">All Levels</option>
          <option value="INFO">INFO</option>
          <option value="WARN">WARN</option>
          <option value="ERROR">ERROR</option>
          <option value="DEBUG">DEBUG</option>
        </select>
      </div>
      
      {expanded && (
        <div className="logs-content">
          {filteredLogs.map(log => (
            <div key={log.id} className={`log-entry log-${log.level.toLowerCase()}`}>
              <span className="log-timestamp">
                {formatTimestamp(log.timestamp)}
              </span>
              <span className="log-level">{log.level}</span>
              <span className="log-message">{log.message}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
```

## API Integration

### API Service

```typescript
// services/api.ts
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

class ApiClient {
  private baseUrl: string;
  
  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }
  
  async get<T>(endpoint: string): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`);
    if (!response.ok) {
      throw new Error(`API error: ${response.statusText}`);
    }
    const data = await response.json();
    return data.data;
  }
  
  async post<T>(endpoint: string, body: any): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    });
    if (!response.ok) {
      throw new Error(`API error: ${response.statusText}`);
    }
    const data = await response.json();
    return data.data;
  }
  
  // Similar for PUT, DELETE, PATCH
}

export const apiClient = new ApiClient(API_BASE_URL);
```

### Task Service

```typescript
// services/taskService.ts
export const taskService = {
  async getAllTasks(filters?: TaskFilters): Promise<Task[]> {
    const params = new URLSearchParams();
    if (filters?.status) params.append('status', filters.status);
    if (filters?.schedule_type) params.append('schedule_type', filters.schedule_type);
    if (filters?.search) params.append('search', filters.search);
    
    return apiClient.get<Task[]>(`/tasks?${params}`);
  },
  
  async getTask(uuid: string): Promise<Task> {
    return apiClient.get<Task>(`/tasks/${uuid}`);
  },
  
  async createTask(task: CreateTaskRequest): Promise<Task> {
    return apiClient.post<Task>('/tasks', task);
  },
  
  async pauseTask(uuid: string): Promise<void> {
    await apiClient.post(`/tasks/${uuid}/pause`, {});
  },
  
  async resumeTask(uuid: string): Promise<void> {
    await apiClient.post(`/tasks/${uuid}/resume`, {});
  },
  
  async triggerTask(uuid: string): Promise<Execution> {
    return apiClient.post<Execution>(`/tasks/${uuid}/trigger`, {});
  }
};
```

## State Management

### React Query Setup

```typescript
// hooks/useTasks.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { taskService } from '../services/taskService';

export const useTasks = (filters?: TaskFilters) => {
  return useQuery({
    queryKey: ['tasks', filters],
    queryFn: () => taskService.getAllTasks(filters)
  });
};

export const useTask = (uuid: string) => {
  return useQuery({
    queryKey: ['task', uuid],
    queryFn: () => taskService.getTask(uuid),
    enabled: !!uuid
  });
};

export const usePauseTask = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (uuid: string) => taskService.pauseTask(uuid),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    }
  });
};
```

## Styling

### CSS Approach

**Option 1**: Tailwind CSS (Utility-first)
```tsx
<div className="flex items-center justify-between p-4 bg-white rounded-lg shadow">
  <h2 className="text-xl font-bold">Task Name</h2>
  <button className="px-4 py-2 bg-blue-500 text-white rounded">
    Pause
  </button>
</div>
```

**Option 2**: CSS Modules
```tsx
// TaskCard.module.css
.card {
  display: flex;
  padding: 1rem;
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

// TaskCard.tsx
import styles from './TaskCard.module.css';
<div className={styles.card}>...</div>
```

**Option 3**: Styled Components
```tsx
const Card = styled.div`
  display: flex;
  padding: 1rem;
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
`;
```

**Recommendation**: Tailwind CSS for rapid development and consistency.

## Responsive Design

### Breakpoints

- **Mobile**: < 768px
- **Tablet**: 768px - 1024px
- **Desktop**: > 1024px

### Mobile Considerations

- Collapsible date navigation
- Stacked execution cards
- Bottom navigation for main pages
- Touch-friendly buttons

## Real-time Updates (Future)

### WebSocket Integration

```typescript
// hooks/useExecutionUpdates.ts
export const useExecutionUpdates = (executionUuid: string) => {
  const [execution, setExecution] = useState<Execution | null>(null);
  
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws/executions/${executionUuid}`);
    
    ws.onmessage = (event) => {
      const update = JSON.parse(event.data);
      setExecution(update);
    };
    
    return () => ws.close();
  }, [executionUuid]);
  
  return execution;
};
```

## Testing

### Component Testing

```typescript
// TaskCard.test.tsx
import { render, screen } from '@testing-library/react';
import { TaskCard } from './TaskCard';

test('renders task name', () => {
  const task = { uuid: '123', name: 'Test Task', status: 'ACTIVE' };
  render(<TaskCard task={task} />);
  expect(screen.getByText('Test Task')).toBeInTheDocument();
});
```

## Next Steps

After completing this module:
1. Set up frontend project structure
2. Implement core components
3. Integrate with API
4. Add routing
5. Style components
6. Test components
7. Proceed to Module 9: Testing Strategy

