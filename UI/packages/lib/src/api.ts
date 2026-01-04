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
    apiClient = createApiClient(DEFAULT_API_BASE_URL);
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
}) {
  const client = getApiClient();
  return client.postProjects({ body: project });
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

