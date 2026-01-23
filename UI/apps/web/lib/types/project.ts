export interface Project {
  id: string
  uuid: string
  name: string
  description?: string
  api_key?: string
  execution_endpoint?: string
  alert_emails?: string
  project_users?: ProjectUser[]
  created_at: string
  updated_at: string
}

export interface CreateProjectRequest {
  name: string
  description?: string
  execution_endpoint?: string
}

export interface UpdateProjectRequest {
  name?: string
  description?: string
  execution_endpoint?: string
  alert_emails?: string
}

export type ProjectUserRole = 'admin' | 'readonly'

export interface ProjectUser {
  email: string
  role: ProjectUserRole
}

