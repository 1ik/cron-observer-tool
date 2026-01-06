import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { z } from 'zod';
import { createProject, getProjects, updateProject } from '../api';
import { schemas } from '../api-client';

// Infer the Project type from the Zod schema
export type Project = z.infer<typeof schemas.models_Project>;

/**
 * Query keys for projects
 */
export const projectKeys = {
  all: ['projects'] as const,
  lists: () => [...projectKeys.all, 'list'] as const,
  list: (filters: string) => [...projectKeys.lists(), { filters }] as const,
  details: () => [...projectKeys.all, 'detail'] as const,
  detail: (id: string) => [...projectKeys.details(), id] as const,
};

/**
 * Hook to fetch all projects
 */
export function useProjects() {
  return useQuery({
    queryKey: projectKeys.lists(),
    queryFn: async () => {
      const response = await getProjects();
      return response;
    },
    refetchOnMount: false, // Don't refetch when component remounts
    refetchOnWindowFocus: false, // Don't refetch when window regains focus
    refetchOnReconnect: false, // Don't refetch on reconnect
  });
}

/**
 * Hook to create a new project
 */
export function useCreateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createProject,
    onSuccess: () => {
      // Invalidate and refetch projects list after creating a new project
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() });
    },
  });
}

/**
 * Hook to update an existing project
 */
export function useUpdateProject(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (project: z.infer<typeof schemas.models_UpdateProjectRequest>) => {
      return updateProject(projectId, project);
    },
    onSuccess: () => {
      // Invalidate and refetch projects list
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() });
      // Also invalidate the specific project if we have a detail query key
      queryClient.invalidateQueries({ queryKey: projectKeys.detail(projectId) });
    },
  });
}

