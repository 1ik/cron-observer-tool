import { TaskGroup } from '../types/taskgroup'

export const mockTaskGroups: TaskGroup[] = [
  {
    id: '65b2a5e0a7c8d9e1f2a3b4c9',
    uuid: 'group-uuid-001',
    project_id: '65b2a5e0a7c8d9e1f2a3b4c5',
    name: 'Morning Tasks',
    description: 'Tasks that run in the morning',
    status: 'ACTIVE',
    start_time: '09:00',
    end_time: '12:00',
    timezone: 'America/New_York',
    created_at: '2025-01-01T10:00:00Z',
    updated_at: '2025-01-01T10:00:00Z',
  },
  {
    id: '65b2a5e0a7c8d9e1f2a3b4ca',
    uuid: 'group-uuid-002',
    project_id: '65b2a5e0a7c8d9e1f2a3b4c5',
    name: 'Evening Tasks',
    description: 'Tasks that run in the evening',
    status: 'ACTIVE',
    start_time: '18:00',
    end_time: '22:00',
    timezone: 'America/New_York',
    created_at: '2025-01-01T10:00:00Z',
    updated_at: '2025-01-01T10:00:00Z',
  },
]

