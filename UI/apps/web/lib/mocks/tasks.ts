import { Task } from '../types/task'

export const mockTasks: Task[] = [
  {
    id: '65b2a5e0a7c8d9e1f2a3b4c1',
    uuid: 'task-uuid-001',
    project_id: '65b2a5e0a7c8d9e1f2a3b4c5',
    name: 'Daily Backup',
    description: 'Backup database daily at 2 AM',
    schedule_type: 'RECURRING',
    status: 'ACTIVE',
    schedule_config: {
      cron_expression: '0 2 * * *',
      timezone: 'America/New_York',
    },
    trigger_config: {
      type: 'HTTP',
      http: {
        url: 'https://api.example.com/backup',
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: { type: 'full' },
        timeout: 300,
      },
    },
    created_at: '2025-01-01T10:00:00Z',
    updated_at: '2025-01-01T10:00:00Z',
  },
  {
    id: '65b2a5e0a7c8d9e1f2a3b4c2',
    uuid: 'task-uuid-002',
    project_id: '65b2a5e0a7c8d9e1f2a3b4c5',
    task_group_id: '65b2a5e0a7c8d9e1f2a3b4c9',
    name: 'Send Morning Report',
    description: 'Send daily morning report',
    schedule_type: 'RECURRING',
    status: 'ACTIVE',
    schedule_config: {
      timezone: 'America/New_York',
      time_range: {
        start: '09:00',
        end: '10:00',
        frequency: { value: 15, unit: 'm' },
      },
      days_of_week: [1, 2, 3, 4, 5],
    },
    trigger_config: {
      type: 'HTTP',
      http: {
        url: 'https://api.example.com/reports/morning',
        method: 'GET',
        timeout: 60,
      },
    },
    created_at: '2025-01-01T10:00:00Z',
    updated_at: '2025-01-01T10:00:00Z',
  },
  {
    id: '65b2a5e0a7c8d9e1f2a3b4c3',
    uuid: 'task-uuid-003',
    project_id: '65b2a5e0a7c8d9e1f2a3b4c5',
    name: 'Cleanup Logs',
    description: 'Clean up old log files',
    schedule_type: 'RECURRING',
    status: 'PAUSED',
    schedule_config: {
      cron_expression: '0 0 * * 0',
      timezone: 'UTC',
    },
    trigger_config: {
      type: 'HTTP',
      http: {
        url: 'https://api.example.com/cleanup',
        method: 'POST',
        timeout: 120,
      },
    },
    created_at: '2025-01-01T10:00:00Z',
    updated_at: '2025-01-01T10:00:00Z',
  },
  {
    id: '65b2a5e0a7c8d9e1f2a3b4c4',
    uuid: 'task-uuid-004',
    project_id: '65b2a5e0a7c8d9e1f2a3b4c5',
    task_group_id: '65b2a5e0a7c8d9e1f2a3b4c9',
    name: 'Process Orders',
    description: 'Process pending orders every 30 minutes',
    schedule_type: 'RECURRING',
    status: 'ACTIVE',
    schedule_config: {
      timezone: 'America/New_York',
      time_range: {
        start: '09:00',
        end: '12:00',
        frequency: { value: 30, unit: 'm' },
      },
      days_of_week: [1, 2, 3, 4, 5],
    },
    trigger_config: {
      type: 'HTTP',
      http: {
        url: 'https://api.example.com/orders/process',
        method: 'POST',
        timeout: 180,
      },
    },
    created_at: '2025-01-01T10:00:00Z',
    updated_at: '2025-01-01T10:00:00Z',
  },
]

