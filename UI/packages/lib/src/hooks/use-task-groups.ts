import { useQuery } from '@tanstack/react-query';
import type { z } from 'zod';
import { getTaskGroupsByProject } from '../api';
import { schemas } from '../api-client';

// Infer the TaskGroup type from the Zod schema
export type TaskGroup = z.infer<typeof schemas.models_TaskGroup>;

/**
 * Query keys for task groups
 */
export const taskGroupKeys = {
  all: ['taskGroups'] as const,
  lists: () => [...taskGroupKeys.all, 'list'] as const,
  list: (projectId: string) => [...taskGroupKeys.lists(), { projectId }] as const,
  details: () => [...taskGroupKeys.all, 'detail'] as const,
  detail: (id: string) => [...taskGroupKeys.details(), id] as const,
};

/**
 * Hook to fetch all task groups for a project
 */
export function useTaskGroupsByProject(projectId: string) {
  return useQuery({
    queryKey: taskGroupKeys.list(projectId),
    queryFn: async () => {
      const response = await getTaskGroupsByProject(projectId);
      return response;
    },
    enabled: !!projectId, // Only fetch if projectId is provided
  });
}


