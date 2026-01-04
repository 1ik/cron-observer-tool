export type TaskGroupStatus = 'ACTIVE' | 'PAUSED' | 'DISABLED'

export interface TaskGroup {
  id: string
  uuid: string
  project_id: string
  name: string
  description?: string
  status: TaskGroupStatus
  start_time?: string // "HH:MM"
  end_time?: string // "HH:MM"
  timezone?: string
  created_at: string
  updated_at: string
}

export interface CreateTaskGroupRequest {
  project_id: string
  name: string
  description?: string
  status?: TaskGroupStatus
  start_time?: string
  end_time?: string
  timezone?: string
}

export interface UpdateTaskGroupRequest {
  name: string
  description?: string
  status?: TaskGroupStatus
  start_time?: string
  end_time?: string
  timezone?: string
}

