import { useQuery } from '@tanstack/react-query';
import { getExecutionsByTaskUUID } from '../api';

/**
 * Query keys for executions
 */
export const executionKeys = {
  all: ['executions'] as const,
  lists: () => [...executionKeys.all, 'list'] as const,
  list: (projectId: string, taskUUID: string, date: string) => [...executionKeys.lists(), projectId, taskUUID, date] as const,
};

/**
 * Hook to fetch executions for a specific task
 * @param projectId - Project ID (MongoDB ObjectID)
 * @param taskUUID - Task UUID
 * @param date - Required date filter (YYYY-MM-DD format). Returns executions for that date only
 * @param enabled - Whether the query should be enabled (default: true)
 */
export function useExecutionsByTask(projectId: string | null | undefined, taskUUID: string | null | undefined, date: string, enabled = true) {
  return useQuery({
    queryKey: executionKeys.list(projectId || '', taskUUID || '', date),
    queryFn: async () => {
      // Validate all parameters before making the API call
      if (!projectId || !taskUUID || !date || date.trim() === '') {
        console.warn('useExecutionsByTask: Missing required parameters', { projectId, taskUUID, date });
        return [];
      }
      
      // Ensure date is a valid format (YYYY-MM-DD)
      const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
      const trimmedDate = date.trim();
      if (!dateRegex.test(trimmedDate)) {
        console.error('useExecutionsByTask: Invalid date format', { date: trimmedDate });
        throw new Error(`Invalid date format: ${trimmedDate}. Expected YYYY-MM-DD`);
      }
      
      console.log('useExecutionsByTask: Calling API with', { projectId, taskUUID, date: trimmedDate });
      return getExecutionsByTaskUUID(projectId, taskUUID, trimmedDate);
    },
    enabled: enabled && !!projectId && !!taskUUID && !!date && date.trim() !== '',
    retry: false, // Disable retries
    staleTime: 0, // Always refetch when query key changes (taskUUID, date, or projectId)
    refetchOnMount: true, // Refetch when component remounts to ensure fresh data
    refetchOnWindowFocus: false, // Don't refetch when window regains focus
    // Query will automatically refetch when queryKey changes (different taskUUID)
  });
}

