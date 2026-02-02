// TaskStatus includes all possible statuses from the API
// Note: PENDING_DELETE and DELETE_FAILED are internal orchestration statuses
// that are filtered out by the backend and should not appear in the UI
export type TaskStatus = 'ACTIVE' | 'DISABLED' | 'PENDING_DELETE' | 'DELETE_FAILED'
export type TaskState = 'RUNNING' | 'NOT_RUNNING' // System-controlled: based on time window
export type ScheduleType = 'RECURRING' | 'ONEOFF'
export type FrequencyUnit = 's' | 'm' | 'h'

export interface Frequency {
  value: number
  unit: FrequencyUnit
}

export interface TimeRange {
  start: string // "HH:MM"
  end: string // "HH:MM"
  frequency: Frequency
}

export interface ScheduleConfig {
  cron_expression?: string
  timezone: string
  time_range?: TimeRange
  days_of_week?: number[]
  exclusions?: number[]
}

export interface HTTPTriggerConfig {
  url: string
  method: string
  headers?: Record<string, string>
  body?: unknown
  timeout?: number
}

export interface TriggerConfig {
  type: 'HTTP'
  http: HTTPTriggerConfig
}

export interface Task {
  id: string
  uuid: string
  project_id: string
  task_group_id?: string
  name: string
  description?: string
  schedule_type: ScheduleType
  status: TaskStatus
  state: TaskState // System-controlled: based on time window
  schedule_config: ScheduleConfig
  trigger_config?: TriggerConfig // Deprecated: Tasks now use project's execution_endpoint
  timeout_seconds?: number // Optional timeout in seconds
  metadata?: Record<string, unknown>
  created_at: string
  updated_at: string
}

export interface CreateTaskRequest {
  project_id: string
  task_group_id?: string
  name: string
  description?: string
  schedule_type: ScheduleType
  status?: TaskStatus
  schedule_config: ScheduleConfig
  timeout_seconds?: number
  metadata?: Record<string, unknown>
}

export interface UpdateTaskRequest {
  task_group_id?: string
  name: string
  description?: string
  schedule_type: ScheduleType
  status?: TaskStatus
  schedule_config: ScheduleConfig
  timeout_seconds?: number
  metadata?: Record<string, unknown>
}

