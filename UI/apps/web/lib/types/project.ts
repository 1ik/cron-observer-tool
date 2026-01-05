export interface Project {
  id: string
  uuid: string
  name: string
  description?: string
  api_key?: string
  execution_endpoint?: string
  created_at: string
  updated_at: string
}

export interface CreateProjectRequest {
  name: string
  description?: string
  execution_endpoint?: string
}

