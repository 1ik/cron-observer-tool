import { createApiClient } from './api-client';

// Default API base URL - can be overridden via environment variable
const getDefaultApiBaseUrl = (): string => {
  // Check for environment variable first (for both client and server)
  try {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const env = (globalThis as any).process?.env;
    if (env?.NEXT_PUBLIC_API_BASE_URL) {
      return env.NEXT_PUBLIC_API_BASE_URL;
    }
  } catch {
    // process not available, continue
  }

  // In browser, check for window.env or use default backend URL
  if (typeof window !== 'undefined') {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const windowEnv = (window as any).__ENV__;
    if (windowEnv?.NEXT_PUBLIC_API_BASE_URL) {
      return windowEnv.NEXT_PUBLIC_API_BASE_URL;
    }
    // Default to backend server on port 8080
    return 'http://localhost:8080/api/v1';
  }

  // Server-side default
  return 'http://localhost:8080/api/v1';
};

const DEFAULT_API_BASE_URL = getDefaultApiBaseUrl();

// Create a singleton API client instance
let apiClient: ReturnType<typeof createApiClient> | null = null;

/**
 * Get or create the API client instance
 */
function getApiClient(): ReturnType<typeof createApiClient> {
  if (!apiClient) {
    apiClient = createApiClient(DEFAULT_API_BASE_URL, {
      validate: false, // Disable Zod validation
    });
  }
  return apiClient;
}

// ============================================================================
// Projects API
// ============================================================================

/**
 * Get all projects
 * @returns Promise resolving to an array of projects
 */
export async function getProjects() {
  const client = getApiClient();
  return client.getProjects();
}

/**
 * Create a new project
 * @param project - Project creation data
 * @returns Promise resolving to the created project
 */
export async function createProject(project: {
  name: string;
  description?: string;
  execution_endpoint?: string;
}) {
  const client = getApiClient();
  
  // Create a clean object that matches the schema exactly
  const requestBody: {
    name: string;
    description?: string;
    execution_endpoint?: string;
  } = {
    name: project.name,
  };
  
  // Only include optional fields if they have values
  if (project.description) {
    requestBody.description = project.description;
  }
  if (project.execution_endpoint) {
    requestBody.execution_endpoint = project.execution_endpoint;
  }
  
  // Debug: log what we're sending
  console.log('createProject - original project:', project);
  console.log('createProject - requestBody:', requestBody);
  console.log('createProject - requestBody.name:', requestBody.name);
  
  // Ensure name is a string and not undefined
  if (!requestBody.name || typeof requestBody.name !== 'string') {
    throw new Error(`Invalid project data: name must be a non-empty string, got: ${JSON.stringify(requestBody)}`);
  }
  
  // When there are NO path parameters, Zodios expects the body as the first argument directly
  // (unlike endpoints with path params which use: client.method(body, { params: { ... } }))
  return client.postProjects(requestBody);
}

/**
 * Update an existing project
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param project - Project update data
 * @returns Promise resolving to the updated project
 */
export async function updateProject(
  projectId: string,
  project: {
    name?: string;
    description?: string;
    execution_endpoint?: string;
    alert_emails?: string;
  }
) {
  const client = getApiClient();

  // Create a clean object that matches the schema exactly
  const requestBody: {
    name?: string;
    description?: string;
    execution_endpoint?: string;
    alert_emails?: string;
  } = {};

  // Only include fields that are provided
  if (project.name !== undefined) {
    requestBody.name = project.name.trim();
  }
  if (project.description !== undefined) {
    requestBody.description = project.description.trim() || '';
  }
  if (project.execution_endpoint !== undefined) {
    requestBody.execution_endpoint = project.execution_endpoint.trim() || '';
  }
  if (project.alert_emails !== undefined) {
    requestBody.alert_emails = project.alert_emails.trim() || '';
  }

  return client.putProjectsProject_id(requestBody, { params: { project_id: projectId } });
}

// ============================================================================
// Tasks API
// ============================================================================

/**
 * Get all tasks for a project
 * @param projectId - Project ID
 * @returns Promise resolving to an array of tasks
 */
export async function getTasksByProject(projectId: string) {
  const client = getApiClient();
  return client.getProjectsProject_idtasks({ params: { project_id: projectId } });
}

/**
 * Create a new task in a project
 * @param projectId - Project ID
 * @param task - Task creation data
 * @returns Promise resolving to the created task
 */
export async function createTask(
  projectId: string,
  task: {
    project_id: string;
    task_group_id?: string;
    name: string;
    description?: string;
    schedule_type: 'RECURRING' | 'ONEOFF';
    status?: 'ACTIVE' | 'DISABLED';
    schedule_config: {
      cron_expression?: string;
      timezone: string;
      time_range?: {
        start: string;
        end: string;
        frequency: { value: number; unit: 's' | 'm' | 'h' };
      };
      days_of_week?: number[];
      exclusions?: number[];
    };
    trigger_config: {
      type: 'HTTP';
      http: {
        url: string;
        method: string;
        headers?: Record<string, string>;
        body?: unknown;
        timeout?: number;
      };
    };
    metadata?: Record<string, unknown>;
  }
) {
  const client = getApiClient();
  return client.postProjectsProject_idtasks(task, { params: { project_id: projectId } });
}

/**
 * Update a task (full update)
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param taskUUID - Task UUID
 * @param task - Task update request
 * @returns Promise resolving to the updated task
 */
export async function updateTask(
  projectId: string,
  taskUUID: string,
  task: {
    name?: string;
    description?: string;
    schedule_type?: 'RECURRING' | 'ONEOFF';
    status?: 'ACTIVE' | 'DISABLED';
    schedule_config?: {
      cron_expression?: string;
      timezone: string;
      time_range?: {
        start: string;
        end: string;
        frequency: {
          value: number;
          unit: 's' | 'm' | 'h';
        };
      };
      days_of_week?: number[];
      exclusions?: number[];
    };
    metadata?: Record<string, unknown>;
    task_group_id?: string;
  }
) {
  const client = getApiClient();
  return client.putProjectsProject_idtasksTask_uuid(task, {
    params: { project_id: projectId, task_uuid: taskUUID },
  });
}

/**
 * Update task status
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param taskUUID - Task UUID
 * @param status - New status (ACTIVE or DISABLED)
 * @returns Promise resolving to the updated task
 */
export async function updateTaskStatus(projectId: string, taskUUID: string, status: 'ACTIVE' | 'DISABLED') {
  const client = getApiClient();
  
  if (!projectId || !taskUUID || !status) {
    throw new Error('Missing required parameters: projectId, taskUUID, and status are required');
  }
  
  if (status !== 'ACTIVE' && status !== 'DISABLED') {
    throw new Error('Status must be either ACTIVE or DISABLED');
  }
  
  // Use fetch directly since the endpoint might not be in the generated client yet
  const baseUrl = getDefaultApiBaseUrl();
  const response = await fetch(
    `${baseUrl}/projects/${projectId}/tasks/${taskUUID}/status`,
    {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ status }),
    }
  );
  
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to update task status' }));
    throw new Error(error.error || `HTTP ${response.status}: Failed to update task status`);
  }
  
  return response.json();
}

// ============================================================================
// Executions API
// ============================================================================

/**
 * Get executions for a specific task
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param taskUUID - Task UUID
 * @param date - Required date filter (YYYY-MM-DD format). Returns executions for that date only
 * @returns Promise resolving to an array of executions
 */
export async function getExecutionsByTaskUUID(projectId: string, taskUUID: string, date: string) {
  if (!date || typeof date !== 'string' || date.trim() === '') {
    throw new Error('Date parameter is required and must be a non-empty string (YYYY-MM-DD format)');
  }
  
  const client = getApiClient();
  const trimmedDate = date.trim();
  
  // Ensure all required params are present
  if (!projectId || !taskUUID || !trimmedDate) {
    throw new Error(`Missing required parameters: projectId=${!!projectId}, taskUUID=${!!taskUUID}, date=${!!trimmedDate}`);
  }
  
  // openapi-zod-client generates incomplete types (doesn't include query params in method signature)
  // Zodios supports query parameters via a separate 'queries' property
  // We use 'as any' because the generated types are incomplete
  return (client.getProjectsProject_idtasksTask_uuidexecutions as any)({
    params: {
      project_id: projectId,
      task_uuid: taskUUID,
    },
    queries: {
      date: trimmedDate,
    },
  });
}

// ============================================================================
// Task Groups API
// ============================================================================

/**
 * Get all task groups for a project
 * @param projectId - Project ID
 * @returns Promise resolving to an array of task groups
 */
export async function getTaskGroupsByProject(projectId: string) {
  const client = getApiClient();
  return client.getProjectsProject_idtaskGroups({ params: { project_id: projectId } });
}

/**
 * Create a new task group in a project
 * @param projectId - Project ID
 * @param taskGroup - Task group creation data
 * @returns Promise resolving to the created task group
 */
export async function createTaskGroup(
  projectId: string,
  taskGroup: {
    name: string;
    project_id: string;
    description?: string;
    status?: 'ACTIVE' | 'DISABLED';
    start_time?: string;
    end_time?: string;
    timezone?: string;
  }
) {
  const client = getApiClient();
  return client.postProjectsProject_idtaskGroups(taskGroup, { params: { project_id: projectId } });
}

export async function updateTaskGroup(
  projectId: string,
  groupUuid: string,
  taskGroup: {
    name?: string;
    description?: string;
    status?: 'ACTIVE' | 'DISABLED';
    start_time?: string;
    end_time?: string;
    timezone?: string;
  }
) {
  const client = getApiClient();
  return client.putProjectsProject_idtaskGroupsGroup_uuid(taskGroup, {
    params: { project_id: projectId, group_uuid: groupUuid },
  });
}

