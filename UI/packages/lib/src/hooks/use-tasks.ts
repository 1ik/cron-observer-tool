import { useQuery } from '@tanstack/react-query';
import type { z } from 'zod';
import { getTasksByProject } from '../api';
import { schemas } from '../api-client';

// Infer the Task type from the Zod schema
export type Task = z.infer<typeof schemas.models_Task>;

/**
 * Query keys for tasks
 */
export const taskKeys = {
  all: ['tasks'] as const,
  lists: () => [...taskKeys.all, 'list'] as const,
  list: (projectId: string) => [...taskKeys.lists(), { projectId }] as const,
  details: () => [...taskKeys.all, 'detail'] as const,
  detail: (id: string) => [...taskKeys.details(), id] as const,
};

/**
 * Hook to fetch all tasks for a project
 */
export function useTasksByProject(projectId: string) {
  return useQuery({
    queryKey: taskKeys.list(projectId),
    queryFn: async () => {
      const response = await getTasksByProject(projectId);
      return response;
    },
    enabled: !!projectId, // Only fetch if projectId is provided
  });
}


