export type ExecutionStatus = 'PENDING' | 'RUNNING' | 'SUCCESS' | 'FAILED' | 'CANCELLED'

export interface Execution {
  id: string
  task_id: string
  task_uuid: string
  task_name: string
  status: ExecutionStatus
  started_at: string
  completed_at?: string
  duration_ms?: number
  error_message?: string
  response_status?: number
  response_body?: unknown
  created_at: string
}

