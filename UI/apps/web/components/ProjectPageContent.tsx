'use client'

import { useProjects, useTaskGroupsByProject, useTasksByProject } from '@cron-observer/lib'
import { Box, Flex, Heading, Spinner, Text } from '@radix-ui/themes'
import { mockExecutions } from '../lib/mocks/executions'
import { Execution } from '../lib/types/execution'
import { ProjectLayout } from './ProjectLayout'

interface ProjectPageContentProps {
  projectId: string
  selectedTaskId?: string | null
}

export function ProjectPageContent({ projectId, selectedTaskId }: ProjectPageContentProps) {
  // Fetch projects to find the one we need
  const { data: projects = [], isLoading: isLoadingProjects } = useProjects()
  
  // Find project by id or uuid
  const project = projects.find((p) => p.id === projectId || p.uuid === projectId)
  
  // Use the project's ObjectID (id) for API calls, not the UUID
  // The backend expects MongoDB ObjectID format (hex string), not UUID
  const projectObjectId = project?.id
  
  // Only fetch tasks and task groups if we have the ObjectID (after projects load)
  const { data: tasks = [], isLoading: isLoadingTasks, error: tasksError } = useTasksByProject(projectObjectId || '')
  const { data: taskGroups = [], isLoading: isLoadingTaskGroups, error: taskGroupsError } = useTaskGroupsByProject(projectObjectId || '')

  // Loading state
  if (isLoadingProjects || isLoadingTasks || isLoadingTaskGroups) {
    return (
      <Flex justify="center" align="center" style={{ minHeight: '400px' }}>
        <Flex direction="column" gap="3" align="center">
          <Spinner size="3" />
          <Text size="3" color="gray">
            Loading project data...
          </Text>
        </Flex>
      </Flex>
    )
  }

  // Project not found
  if (!project) {
    return (
      <Box p="8">
        <Flex direction="column" align="center" gap="4">
          <Heading size="5" color="red">
            Project Not Found
          </Heading>
          <Text size="3" color="gray" align="center">
            The project you&apos;re looking for doesn&apos;t exist.
          </Text>
        </Flex>
      </Box>
    )
  }

  // Filter executions for this project (using mock data for now)
  const projectExecutions: Execution[] = mockExecutions.filter((e) => {
    const task = tasks.find((t) => t.id === e.task_id || t.uuid === e.task_uuid)
    return task !== undefined
  })
  
  // Map API response types to component types
  // The API returns types that match our component types, but we need to ensure
  // all required fields are present
  const mappedTasks = (tasks || []).map((task) => ({
      id: task.id || task.uuid || '',
      uuid: task.uuid || task.id || '',
      project_id: task.project_id || projectId,
      task_group_id: task.task_group_id,
      name: task.name || '',
      description: task.description,
      schedule_type: task.schedule_type || 'RECURRING',
      status: task.status || 'ACTIVE',
      schedule_config: task.schedule_config || { timezone: 'UTC' },
      trigger_config: (task.trigger_config?.type === 'HTTP' && task.trigger_config?.http)
        ? { type: 'HTTP' as const, http: task.trigger_config.http }
        : { type: 'HTTP' as const, http: { url: '', method: 'GET' } },
      metadata: task.metadata,
      created_at: task.created_at || new Date().toISOString(),
      updated_at: task.updated_at || new Date().toISOString(),
    }))

  const mappedTaskGroups = taskGroups.map((tg) => ({
    id: tg.id || tg.uuid || '',
    uuid: tg.uuid || tg.id || '',
    project_id: tg.project_id || projectId,
    name: tg.name || '',
    description: tg.description,
    status: tg.status || 'ACTIVE',
    start_time: tg.start_time,
    end_time: tg.end_time,
    timezone: tg.timezone,
    created_at: tg.created_at || new Date().toISOString(),
    updated_at: tg.updated_at || new Date().toISOString(),
    }))


  return (
    <Box style={{ height: '100%', width: '100%' }}>
      <ProjectLayout
        project={{
          id: project.id || project.uuid || '',
          uuid: project.uuid || project.id || '',
          name: project.name || '',
          description: project.description,
          api_key: project.api_key,
          execution_endpoint: project.execution_endpoint,
          alert_emails: (project as Record<string, unknown>).alert_emails && typeof (project as Record<string, unknown>).alert_emails === 'string' 
            ? (project as Record<string, unknown>).alert_emails as string 
            : undefined,
          created_at: project.created_at || new Date().toISOString(),
          updated_at: project.updated_at || new Date().toISOString(),
        }}
        taskGroups={mappedTaskGroups}
        tasks={mappedTasks}
        executions={projectExecutions}
        selectedTaskId={selectedTaskId}
      />
    </Box>
  )
}

