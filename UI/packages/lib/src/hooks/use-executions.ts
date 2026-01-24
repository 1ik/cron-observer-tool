import { useQuery } from '@tanstack/react-query';
import { getExecutionsByTaskUUID, getFailedExecutionsStats, getExecutionStats } from '../api';

/**
 * Query keys for executions
 */
export const executionKeys = {
  all: ['executions'] as const,
  lists: () => [...executionKeys.all, 'list'] as const,
  list: (projectId: string, taskUUID: string, date: string, page: number, pageSize: number) => 
    [...executionKeys.lists(), projectId, taskUUID, date, page, pageSize] as const,
};

/**
 * Hook to fetch paginated executions for a specific task
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param taskUUID - Task UUID
 * @param date - Required date filter (YYYY-MM-DD format). Returns executions for that date only
 * @param page - Page number (default: 1)
 * @param pageSize - Page size (default: 100)
 * @param enabled - Whether the query should be enabled (default: true)
 */
export function useExecutionsByTask(
  projectId: string | null | undefined,
  taskUUID: string | null | undefined,
  date: string,
  page: number = 1,
  pageSize: number = 100,
  enabled = true
) {
  return useQuery({
    queryKey: executionKeys.list(projectId || '', taskUUID || '', date, page, pageSize),
    queryFn: async () => {
      // Validate all parameters before making the API call
      if (!projectId || !taskUUID || !date || date.trim() === '') {
        console.warn('useExecutionsByTask: Missing required parameters', { projectId, taskUUID, date });
        return {
          data: [],
          page: 1,
          page_size: pageSize,
          total_count: 0,
          total_pages: 0,
        };
      }
      
      // Ensure date is a valid format (YYYY-MM-DD)
      const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
      const trimmedDate = date.trim();
      if (!dateRegex.test(trimmedDate)) {
        console.error('useExecutionsByTask: Invalid date format', { date: trimmedDate });
        throw new Error(`Invalid date format: ${trimmedDate}. Expected YYYY-MM-DD`);
      }
      
      console.log('useExecutionsByTask: Calling API with', { projectId, taskUUID, date: trimmedDate, page, pageSize });
      return getExecutionsByTaskUUID(projectId, taskUUID, trimmedDate, page, pageSize);
    },
    enabled: enabled && !!projectId && !!taskUUID && !!date && date.trim() !== '',
    retry: false, // Disable retries
    staleTime: 0, // Always refetch when query key changes (taskUUID, date, or projectId)
    refetchOnMount: true, // Refetch when component remounts to ensure fresh data
    refetchOnWindowFocus: false, // Don't refetch when window regains focus
    // Query will automatically refetch when queryKey changes (different taskUUID)
  });
}

/**
 * Hook to fetch failure statistics for a project
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param days - Number of days to look back (default: 7)
 */
export function useFailedExecutionsStats(projectId: string | null, days: number = 7) {
  return useQuery({
    queryKey: ['executions', 'failed-stats', projectId, days],
    queryFn: () => getFailedExecutionsStats(projectId!, days),
    enabled: !!projectId,
    staleTime: 30000, // 30 seconds
  });
}

/**
 * Hook to fetch execution statistics for a project
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param days - Number of days to look back (default: 7)
 */
export function useExecutionStats(projectId: string | null, days: number = 7) {
  return useQuery({
    queryKey: ['execution-stats', projectId, days],
    queryFn: () => getExecutionStats(projectId!, days),
    enabled: !!projectId,
    staleTime: 30000, // 30 seconds
  });
}

