import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { z } from 'zod';
import { createTask, getTasksByProject, updateTask, updateTaskStatus } from '../api';
import { schemas } from '../api-client';

// Infer the Task type from the Zod schema
export type Task = z.infer<typeof schemas.models_Task>;

// CreateTaskRequest type inferred from schema
type CreateTaskRequest = z.infer<typeof schemas.models_CreateTaskRequest>;

// UpdateTaskRequest type inferred from schema
type UpdateTaskRequest = z.infer<typeof schemas.models_UpdateTaskRequest>;

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
    retry: false, // Disable retries
    refetchOnMount: false, // Don't refetch when component remounts
    refetchOnWindowFocus: false, // Don't refetch when window regains focus
    refetchOnReconnect: false, // Don't refetch on reconnect
  });
}

/**
 * Hook to create a new task
 * @param projectId - The ID of the project
 */
export function useCreateTask(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateTaskRequest) => {
      return createTask(projectId, data);
    },
    onSuccess: () => {
      // Invalidate and refetch tasks for this project
      queryClient.invalidateQueries({ queryKey: taskKeys.list(projectId) });
    },
  });
}

/**
 * Hook to update a task (full update)
 * @param projectId - The ID of the project
 */
export function useUpdateTask(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ taskUUID, data }: { taskUUID: string; data: UpdateTaskRequest }) => {
      return updateTask(projectId, taskUUID, data);
    },
    onSuccess: () => {
      // Invalidate and refetch tasks for this project
      queryClient.invalidateQueries({ queryKey: taskKeys.list(projectId) });
    },
    retry: false, // Disable retries
  });
}

/**
 * Hook to update task status (pause/play)
 * @param projectId - The ID of the project
 * @param taskUUID - The UUID of the task
 */
export function useUpdateTaskStatus(projectId: string, taskUUID: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (status: 'ACTIVE' | 'DISABLED') => {
      return updateTaskStatus(projectId, taskUUID, status);
    },
    onSuccess: () => {
      // Invalidate and refetch tasks for this project
      queryClient.invalidateQueries({ queryKey: taskKeys.list(projectId) });
      // Also invalidate task details if needed
      queryClient.invalidateQueries({ queryKey: taskKeys.all });
    },
  });
}

export type { CreateTaskRequest, UpdateTaskRequest };


