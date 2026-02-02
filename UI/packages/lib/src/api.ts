import { createApiClient } from './api-client';

// Default API base URL - can be overridden via environment variable
// IMPORTANT: Use literal process.env.NEXT_PUBLIC_API_BASE_URL so Next.js can inline it at build time.
// Dynamic access (e.g. globalThis.process?.env) is NOT replaced and breaks local dev.
const getDefaultApiBaseUrl = (): string => {
  // Next.js replaces this literal at build time with the value from .env.local / .env
  if (process.env.NEXT_PUBLIC_API_BASE_URL) {
    return process.env.NEXT_PUBLIC_API_BASE_URL;
  }

  // Fallback: Check window.__ENV__ (for custom runtime injection)
  if (typeof window !== 'undefined') {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const windowEnv = (window as any).__ENV__;
    if (windowEnv?.NEXT_PUBLIC_API_BASE_URL) {
      return windowEnv.NEXT_PUBLIC_API_BASE_URL;
    }
    // Relative path for production when API is same-origin or behind proxy
    return '/api/v1';
  }

  // Server-side fallback when NEXT_PUBLIC_API_BASE_URL is not set
  return 'http://localhost:8080/api/v1';
};

const DEFAULT_API_BASE_URL = getDefaultApiBaseUrl();

// Token cache for Authorization header
let cachedToken: string | null = null;
let tokenFetchPromise: Promise<string | null> | null = null;

/**
 * Get the authentication token from the Next.js API
 * Uses caching to avoid repeated fetches
 */
async function getAuthToken(): Promise<string | null> {
  // Return cached token if available
  if (cachedToken) {
    return cachedToken;
  }

  // If a fetch is already in progress, wait for it
  if (tokenFetchPromise) {
    return tokenFetchPromise;
  }

  // Only fetch token on client-side
  if (typeof window === 'undefined') {
    return null;
  }

  tokenFetchPromise = (async () => {
    try {
      const response = await fetch('/api/auth/token');
      if (response.ok) {
        const data = await response.json();
        cachedToken = data.token;
        return cachedToken;
      }
    } catch (error) {
      console.error('Failed to get auth token:', error);
    }
    return null;
  })();

  const token = await tokenFetchPromise;
  tokenFetchPromise = null;
  return token;
}

/**
 * Clear the cached token (call on logout)
 */
export function clearAuthToken(): void {
  cachedToken = null;
  tokenFetchPromise = null;
}

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

    // Add axios interceptor to include Authorization header
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (apiClient as any).axios.interceptors.request.use(async (config: any) => {
      const token = await getAuthToken();
      if (token) {
        config.headers = config.headers || {};
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
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
    project_users?: Array<{ email: string; role: 'admin' | 'readonly' }>;
  }
) {
  const client = getApiClient();

  // Create a clean object that matches the schema exactly
  const requestBody: {
    name?: string;
    description?: string;
    execution_endpoint?: string;
    alert_emails?: string;
    project_users?: Array<{ email: string; role: 'admin' | 'readonly' }>;
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
  if (project.project_users !== undefined) {
    requestBody.project_users = project.project_users;
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
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  const response = await fetch(
    `${baseUrl}/projects/${projectId}/tasks/${taskUUID}/status`,
    {
      method: 'PATCH',
      headers,
      body: JSON.stringify({ status }),
    }
  );
  
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to update task status' }));
    throw new Error(error.error || `HTTP ${response.status}: Failed to update task status`);
  }
  
  return response.json();
}

/**
 * Trigger a task manually
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param taskUUID - Task UUID
 * @returns Promise resolving to the trigger response with execution UUID
 */
export async function triggerTask(projectId: string, taskUUID: string) {
  if (!projectId || !taskUUID) {
    throw new Error('Missing required parameters: projectId and taskUUID are required');
  }
  
  // Use fetch directly since the generated client may have issues with POST requests without body
  const baseUrl = getDefaultApiBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  
  const response = await fetch(
    `${baseUrl}/projects/${projectId}/tasks/${taskUUID}/trigger`,
    {
      method: 'POST',
      headers,
    }
  );
  
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to trigger task' }));
    throw new Error(error.error || `HTTP ${response.status}: Failed to trigger task`);
  }
  
  return response.json();
}

/**
 * Delete a task (instant deletion)
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param taskUUID - Task UUID
 * @returns Promise resolving to the delete response with status
 */
export async function deleteTask(projectId: string, taskUUID: string) {
  if (!projectId || !taskUUID) {
    throw new Error('Missing required parameters: projectId and taskUUID are required');
  }
  
  // Use fetch directly to ensure path parameters are correctly substituted
  // Zodios may have issues with DELETE requests without body
  const baseUrl = getDefaultApiBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  
  // Construct URL with actual parameter values
  const url = `${baseUrl}/projects/${encodeURIComponent(projectId)}/tasks/${encodeURIComponent(taskUUID)}`;
  
  const response = await fetch(url, {
    method: 'DELETE',
    headers,
  });
  
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to delete task' }));
    throw new Error(error.error || `HTTP ${response.status}: Failed to delete task`);
  }
  
  return response.json();
}

// ============================================================================
// Executions API
// ============================================================================

/**
 * Get paginated executions for a specific task
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param taskUUID - Task UUID
 * @param date - Required date filter (YYYY-MM-DD format). Returns executions for that date only
 * @param page - Page number (default: 1)
 * @param pageSize - Page size (default: 100)
 * @returns Promise resolving to paginated executions response
 */
export async function getExecutionsByTaskUUID(
  projectId: string,
  taskUUID: string,
  date: string,
  page: number = 1,
  pageSize: number = 100
) {
  if (!date || typeof date !== 'string' || date.trim() === '') {
    throw new Error('Date parameter is required and must be a non-empty string (YYYY-MM-DD format)');
  }
  
  const client = getApiClient();
  const trimmedDate = date.trim();
  
  // Ensure all required params are present
  if (!projectId || !taskUUID || !trimmedDate) {
    throw new Error(`Missing required parameters: projectId=${!!projectId}, taskUUID=${!!taskUUID}, date=${!!trimmedDate}`);
  }
  
  // Validate pagination parameters
  if (page < 1) {
    throw new Error('Page must be greater than 0');
  }
  if (pageSize < 1 || pageSize > 100) {
    throw new Error('Page size must be between 1 and 100');
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
      page,
      page_size: pageSize,
    },
  });
}

/**
 * Get failure statistics for a project
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param days - Number of days to look back (default: 7)
 * @returns Promise resolving to failure statistics response
 */
export async function getFailedExecutionsStats(
  projectId: string,
  days: number = 7
): Promise<{ stats: Array<{ date: string; count: number }>; total: number }> {
  if (!projectId) {
    throw new Error('Missing required parameter: projectId is required');
  }

  const baseUrl = getDefaultApiBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(
    `${baseUrl}/projects/${projectId}/executions/failed-stats?days=${days}`,
    {
      method: 'GET',
      headers,
    }
  );

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to fetch failure statistics' }));
    throw new Error(error.error || `HTTP ${response.status}: Failed to fetch failure statistics`);
  }

  return response.json();
}

/**
 * Get execution statistics for a project
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param days - Number of days to look back (default: 7)
 * @returns Promise resolving to execution statistics response
 */
export async function getExecutionStats(
  projectId: string,
  days: number = 7
): Promise<{ stats: Array<{ date: string; failures: number; success: number; total: number }> }> {
  if (!projectId) {
    throw new Error('Missing required parameter: projectId is required');
  }

  const baseUrl = getDefaultApiBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(
    `${baseUrl}/projects/${projectId}/executions/stats?days=${days}`,
    {
      method: 'GET',
      headers,
    }
  );

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to fetch execution statistics' }));
    throw new Error(error.error || `HTTP ${response.status}: Failed to fetch execution statistics`);
  }

  return response.json();
}

/**
 * Get task failures by date for a project
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param date - Date in YYYY-MM-DD format
 * @returns Promise resolving to task failures response with calculated_at timestamp
 */
export async function getTaskFailuresByDate(
  projectId: string,
  date: string
): Promise<{ date: string; tasks: Array<{ taskId: string; failures: number }>; total: number; calculated_at: string }> {
  if (!projectId) {
    throw new Error('Missing required parameter: projectId is required');
  }
  if (!date) {
    throw new Error('Missing required parameter: date is required');
  }

  const baseUrl = getDefaultApiBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(
    `${baseUrl}/projects/${projectId}/failures?date=${date}`,
    {
      method: 'GET',
      headers,
    }
  );

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to fetch task failures' }));
    throw new Error(error.error || `HTTP ${response.status}: Failed to fetch task failures`);
  }

  return response.json();
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
    name: string;
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

