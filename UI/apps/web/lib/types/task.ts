export type TaskStatus = 'ACTIVE' | 'PAUSED' | 'DISABLED'
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
  schedule_config: ScheduleConfig
  trigger_config: TriggerConfig
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
  trigger_config: TriggerConfig
  metadata?: Record<string, unknown>
}

