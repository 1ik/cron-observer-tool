'use client'

import { useExecutionsByTask, useProjects, useTaskGroupsByProject, useTasksByProject } from '@cron-observer/lib'
import { Box, Flex, Heading, Spinner, Text } from '@radix-ui/themes'
import { useRouter, useSearchParams } from 'next/navigation'
import { useEffect, useMemo } from 'react'
import { Execution } from '../lib/types/execution'
import { ProjectLayout } from './ProjectLayout'

interface ProjectPageContentProps {
  projectId: string
  selectedTaskId?: string | null
}

export function ProjectPageContent({ projectId, selectedTaskId }: ProjectPageContentProps) {
  const searchParams = useSearchParams()
  const router = useRouter()
  
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

  // Find the selected task to get its UUID (memoized to ensure proper updates)
  const selectedTask = useMemo(() => {
    if (!selectedTaskId) return null
    return tasks.find((t) => t.id === selectedTaskId || t.uuid === selectedTaskId) || null
  }, [selectedTaskId, tasks])

  const selectedTaskUUID = useMemo(() => {
    return selectedTask?.uuid || selectedTask?.id || null
  }, [selectedTask])

  // Get current date string (always computed synchronously)
  const getCurrentDateString = () => {
    const today = new Date()
    const year = today.getFullYear()
    const month = String(today.getMonth() + 1).padStart(2, '0')
    const day = String(today.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  // Get selected date from URL params or default to current date
  // Memoize to prevent unnecessary recalculations and query key changes
  // Always ensure we have a valid date string
  const selectedDate = useMemo(() => {
    const dateParam = searchParams.get('date')
    if (dateParam && dateParam.trim() !== '') {
      return dateParam.trim()
    }
    return getCurrentDateString()
  }, [searchParams])

  // Ensure date is always in URL - update URL if missing
  useEffect(() => {
    const dateParam = searchParams.get('date')
    if (!dateParam || dateParam.trim() === '') {
      const currentDate = getCurrentDateString()
      const params = new URLSearchParams(searchParams.toString())
      params.set('date', currentDate)
      router.replace(`?${params.toString()}`, { scroll: false })
    }
  }, [searchParams, router])

  // Fetch executions for the selected task (only if a task is selected)
  // Ensure selectedDate is always a valid date string before passing to the query
  const validDate = selectedDate && selectedDate.trim() !== '' ? selectedDate : getCurrentDateString()
  
  const { data: executionsData = [], isLoading: isLoadingExecutions } = useExecutionsByTask(
    projectObjectId || null,
    selectedTaskUUID,
    validDate,
    !!projectObjectId && !!selectedTaskUUID && !!validDate && validDate.trim() !== ''
  )

  // Loading state (excluding executions loading - that will be shown in ExecutionsList)
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

  // Map API executions to Execution type
  const projectExecutions: Execution[] = executionsData.map((execution: any) => {
    // Transform logs if they exist
    const transformedLogs = execution.logs && Array.isArray(execution.logs) && execution.logs.length > 0
      ? execution.logs.map((log: any, index: number) => {
          // Normalize log level: backend sends lowercase ("info", "warn", "error"), frontend expects uppercase
          const levelUpper = (log.level || 'info').toUpperCase()
          // Ensure it's a valid log level, default to INFO if not
          const validLevel = ['INFO', 'WARN', 'ERROR', 'DEBUG'].includes(levelUpper) 
            ? levelUpper as 'INFO' | 'WARN' | 'ERROR' | 'DEBUG'
            : 'INFO'
          
          return {
            id: `${execution.id || execution.uuid}-log-${index}`,
            timestamp: log.timestamp || new Date().toISOString(),
            level: validLevel,
            message: log.message || '',
            metadata: log.metadata,
          }
        })
      : undefined

    return {
      id: execution.id || execution.uuid || '',
      task_id: execution.task_id || '',
      task_uuid: execution.task_uuid || '',
      task_name: selectedTask?.name || '',
      status: execution.status || 'PENDING',
      started_at: execution.started_at || new Date().toISOString(),
      completed_at: execution.ended_at || undefined,
      duration_ms: execution.duration_ms,
      error_message: execution.error || undefined,
      logs: transformedLogs,
      created_at: execution.created_at || new Date().toISOString(),
    }
  })
  
  // Map API response types to component types
  // The API returns types that match our component types, but we need to ensure
  // all required fields are present
  const mappedTasks = (tasks || []).map((task) => ({
      id: task.id || task.uuid || '',
      uuid: task.uuid || task.id || '',
      project_id: task.project_id || projectId,
      task_group_id: task.task_group_id || undefined, // Preserve task_group_id, convert empty string to undefined
      name: task.name || '',
      description: task.description,
      schedule_type: task.schedule_type || 'RECURRING',
      status: task.status || 'ACTIVE',
      state: task.state || 'NOT_RUNNING', // System-controlled state
      schedule_config: task.schedule_config || { timezone: 'UTC' },
      trigger_config: task.trigger_config && task.trigger_config.http ? {
        type: 'HTTP' as const,
        http: task.trigger_config.http,
      } : undefined,
      metadata: task.metadata,
      created_at: task.created_at || new Date().toISOString(),
      updated_at: task.updated_at || new Date().toISOString(),
    }))

  const mappedTaskGroups = taskGroups.map((tg) => {
    // Handle status: default to 'ACTIVE' if not set
    const status = tg.status || 'ACTIVE'
    
    return {
      id: tg.id || tg.uuid || '',
      uuid: tg.uuid || '', // Always use uuid field, don't fall back to id
      project_id: tg.project_id || projectId,
      name: tg.name || '',
      description: tg.description,
      status: status as 'ACTIVE' | 'DISABLED',
      state: tg.state || 'NOT_RUNNING', // System-controlled state
      start_time: tg.start_time,
      end_time: tg.end_time,
      timezone: tg.timezone,
      created_at: tg.created_at || new Date().toISOString(),
      updated_at: tg.updated_at || new Date().toISOString(),
    }
  })


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
        isLoadingExecutions={isLoadingExecutions}
      />
    </Box>
  )
}

