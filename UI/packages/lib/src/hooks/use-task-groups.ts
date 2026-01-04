import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { z } from 'zod';
import { createTaskGroup, getTaskGroupsByProject } from '../api';
import { schemas } from '../api-client';

// CreateTaskGroupRequest type inferred from schema
type CreateTaskGroupRequest = z.infer<typeof schemas.models_CreateTaskGroupRequest>;

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

/**
 * Hook to create a new task group
 * @param projectId - The ID of the project
 */
export function useCreateTaskGroup(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateTaskGroupRequest) => {
      return createTaskGroup(projectId, data);
    },
    onSuccess: () => {
      // Invalidate and refetch task groups for this project
      queryClient.invalidateQueries({ queryKey: taskGroupKeys.list(projectId) });
    },
  });
}


