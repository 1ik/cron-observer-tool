import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { z } from 'zod';
import { createTask, deleteTask, getTasksByProject, triggerTask, updateTask, updateTaskStatus } from '../api';
import { schemas } from '../api-client';

// Infer the Task type from the Zod schema
export type Task = z.infer<typeof schemas.models_Task>;

// Create restricted status type (exclude internal orchestration statuses)
type UserTaskStatus = 'ACTIVE' | 'DISABLED';

// CreateTaskRequest type matching the API function signature
type CreateTaskRequest = {
  project_id: string;
  task_group_id?: string;
  name: string;
  description?: string;
  schedule_type: 'RECURRING' | 'ONEOFF';
  status?: UserTaskStatus;
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
  timeout_seconds?: number;
};

// UpdateTaskRequest type matching the API function signature
type UpdateTaskRequest = {
  name: string;
  description?: string;
  schedule_type: 'RECURRING' | 'ONEOFF';
  status?: UserTaskStatus;
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
  task_group_id?: string;
  timeout_seconds?: number;
};

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

/**
 * Hook to trigger a task manually
 * @param projectId - The ID of the project
 * @param taskUUID - The UUID of the task
 */
export function useTriggerTask(projectId: string, taskUUID: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      return triggerTask(projectId, taskUUID);
    },
    onSuccess: () => {
      // Invalidate executions to show the new triggered execution
      queryClient.invalidateQueries({ queryKey: ['executions'] });
    },
  });
}

/**
 * Hook to delete a task (async)
 * @param projectId - The ID of the project
 */
export function useDeleteTask(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (taskUUID: string) => {
      return deleteTask(projectId, taskUUID);
    },
    onSuccess: () => {
      // Invalidate and refetch tasks for this project
      queryClient.invalidateQueries({ queryKey: taskKeys.list(projectId) });
      // Also invalidate all task queries
      queryClient.invalidateQueries({ queryKey: taskKeys.all });
    },
    retry: false, // Disable retries
  });
}

export type { CreateTaskRequest, UpdateTaskRequest };


