'use client'

import { useCreateTask, useCreateTaskGroup, useUpdateProject, useUpdateTask, useUpdateTaskGroup } from '@cron-observer/lib'
import * as DropdownMenu from '@radix-ui/react-dropdown-menu'
import { CaretDownIcon, GearIcon, PlusIcon } from '@radix-ui/react-icons'
import { Box, Button, Flex, IconButton, Text, Tooltip } from '@radix-ui/themes'
import Link from 'next/link'
import { useState } from 'react'
import { useProjectRole } from '../contexts/ProjectRoleContext'
import { Execution } from '../lib/types/execution'
import { Project, UpdateProjectRequest } from '../lib/types/project'
import { CreateTaskRequest, Task, UpdateTaskRequest } from '../lib/types/task'
import { CreateTaskGroupRequest, TaskGroup, UpdateTaskGroupRequest } from '../lib/types/taskgroup'
import { CreateTaskDialog } from './CreateTaskDialog'
import { CreateTaskGroupDialog } from './CreateTaskGroupDialog'
import { ExecutionsList } from './ExecutionsList'
import { ProjectSettingsDialog } from './ProjectSettingsDialog'
import { ResizableSplitter } from './ResizableSplitter'
import { TaskGroupSettingsDialog } from './TaskGroupSettingsDialog'
import { TaskGroupsList } from './TaskGroupsList'
import { TaskSettingsDialog } from './TaskSettingsDialog'

interface ProjectLayoutProps {
  project: Project
  taskGroups: TaskGroup[]
  tasks: Task[]
  executions: Execution[]
  selectedTaskId?: string | null
  isLoadingExecutions?: boolean
  paginationData?: {
    page: number
    page_size: number
    total_count: number
    total_pages: number
  } | null
}

export function ProjectLayout({
  project,
  taskGroups,
  tasks,
  executions,
  selectedTaskId,
  isLoadingExecutions = false,
  paginationData,
}: ProjectLayoutProps) {
  const { canEdit, isReadOnly } = useProjectRole()
  const [selectedTaskGroup, setSelectedTaskGroup] = useState<TaskGroup | null>(null)
  const [isTaskGroupSettingsOpen, setIsTaskGroupSettingsOpen] = useState(false)
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [isTaskSettingsOpen, setIsTaskSettingsOpen] = useState(false)
  const [isCreateTaskDialogOpen, setIsCreateTaskDialogOpen] = useState(false)
  const [isCreateTaskGroupDialogOpen, setIsCreateTaskGroupDialogOpen] = useState(false)
  const [selectedTaskGroupForCreate, setSelectedTaskGroupForCreate] = useState<TaskGroup | null>(null)
  const [isProjectSettingsOpen, setIsProjectSettingsOpen] = useState(false)

  const createTaskGroupMutation = useCreateTaskGroup(project.id)
  const createTaskMutation = useCreateTask(project.id)
  const updateProjectMutation = useUpdateProject(project.id)
  const updateTaskGroupMutation = useUpdateTaskGroup(project.id)
  const updateTaskMutation = useUpdateTask(project.id)

  const handleTaskGroupSettingsClick = (taskGroup: TaskGroup) => {
    setSelectedTaskGroup(taskGroup)
    setIsTaskGroupSettingsOpen(true)
  }

  const handleTaskGroupSettingsSubmit = (data: UpdateTaskGroupRequest) => {
    if (!selectedTaskGroup) return
    
    // Ensure we use the UUID field, not the ID field
    const groupUuid = selectedTaskGroup.uuid || selectedTaskGroup.id
    if (!groupUuid) {
      console.error('Task group missing both uuid and id:', selectedTaskGroup)
      return
    }
    
    console.log('Updating task group:', { groupUuid, id: selectedTaskGroup.id, uuid: selectedTaskGroup.uuid, data })
    
    updateTaskGroupMutation.mutate(
      { groupUuid, data: data as any },
      {
        onSuccess: () => {
          // Dialog will close automatically via onOpenChange
          // Task groups will be refetched automatically via query invalidation
          setIsTaskGroupSettingsOpen(false)
        },
        onError: (error: Error) => {
          // TODO: Add error toast/notification
          console.error('Failed to update task group:', error)
        },
      }
    )
  }

  const handleTaskSettingsClick = (task: Task) => {
    setSelectedTask(task)
    setIsTaskSettingsOpen(true)
  }

  const handleTaskSettingsSubmit = (data: UpdateTaskRequest) => {
    if (!selectedTask) return
    
    // Ensure we use the UUID field, not the ID field
    const taskUuid = selectedTask.uuid || selectedTask.id
    if (!taskUuid) {
      console.error('Task missing both uuid and id:', selectedTask)
      return
    }
    
    updateTaskMutation.mutate(
      { taskUUID: taskUuid, data: data as any },
      {
        onSuccess: () => {
          // Dialog will close automatically via onOpenChange
          // Tasks will be refetched automatically via query invalidation
          setIsTaskSettingsOpen(false)
        },
        onError: (error: Error) => {
          // TODO: Add error toast/notification
          console.error('Failed to update task:', error)
        },
      }
    )
  }

  const handleCreateTaskGroupSubmit = (data: CreateTaskGroupRequest) => {
    createTaskGroupMutation.mutate(data as any, {
      onSuccess: () => {
        // Dialog will close automatically via onOpenChange
        // Task groups will be refetched automatically via query invalidation
      },
      onError: (error: Error) => {
        // TODO: Add error toast/notification
        console.error('Failed to create task group:', error)
      },
    })
  }

  const handleCreateTaskClick = (taskGroup: TaskGroup) => {
    setSelectedTaskGroupForCreate(taskGroup)
    setIsCreateTaskDialogOpen(true)
  }

  const handleCreateTaskSubmit = (data: CreateTaskRequest) => {
    // Set task_group_id from selected task group
    const requestData: CreateTaskRequest = {
      ...data,
      task_group_id: selectedTaskGroupForCreate?.id,
    }
    createTaskMutation.mutate(requestData as any, {
      onSuccess: () => {
        // Dialog will close automatically via onOpenChange
        // Tasks will be refetched automatically via query invalidation
        setSelectedTaskGroupForCreate(null)
      },
      onError: (error: Error) => {
        // TODO: Add error toast/notification
        console.error('Failed to create task:', error)
      },
    })
  }

  const handleProjectSettingsSubmit = (data: UpdateProjectRequest) => {
    updateProjectMutation.mutate(data as any, {
      onSuccess: () => {
        // Dialog will close automatically via onOpenChange
        // Projects will be refetched automatically via query invalidation
        setIsProjectSettingsOpen(false)
      },
      onError: (error: Error) => {
        // TODO: Add error toast/notification
        console.error('Failed to update project:', error)
      },
    })
  }

  return (
    <Flex direction="column" height="100%" width="100%" overflow="hidden">
      {/* Top separator */}
      <Box
        style={{
          width: '100%',
          height: '1px',
          backgroundColor: 'var(--gray-6)',
          flexShrink: 0,
        }}
      />
      
      {/* Header with breadcrumb and create button */}
      <Box
        p="4"
        style={{
          flexShrink: 0,
          borderBottom: '1px solid var(--gray-6)',
        }}
      >
        <Flex align="center" justify="between" gap="3">
          <Flex align="center" gap="2">
            <Link href="/projects" style={{ textDecoration: 'none', cursor: 'pointer' }}>
              <Text size="2" color="gray">
                Projects
              </Text>
            </Link>
            <Text size="2" color="gray">
              &gt;
            </Text>
            <Text size="2" color="gray">
              {project.name}
            </Text>
          </Flex>
          
          <Flex align="center" gap="2">
            {/* Settings Icon */}
            <IconButton
              size="2"
              variant="ghost"
              onClick={() => setIsProjectSettingsOpen(true)}
              style={{ cursor: 'pointer' }}
            >
              <GearIcon width="16" height="16" />
            </IconButton>
            
            {/* GitHub-style split button - only show for users with edit permissions */}
            {canEdit && (
              <DropdownMenu.Root>
                <Flex style={{ display: 'inline-flex', gap: 0 }}>
                  <Button
                    size="2"
                    variant="outline"
                    onClick={() => setIsCreateTaskDialogOpen(true)}
                    style={{
                      borderTopRightRadius: 0,
                      borderBottomRightRadius: 0,
                      borderRight: '1px solid var(--gray-6)',
                      marginRight: '-1px',
                      backgroundColor: 'transparent',
                      cursor: 'pointer',
                    }}
                  >
                    <PlusIcon width="14" height="14" />
                  </Button>
                  <DropdownMenu.Trigger asChild>
                    <Button
                      size="2"
                      variant="outline"
                      style={{
                        borderTopLeftRadius: 0,
                        borderBottomLeftRadius: 0,
                        paddingLeft: 'var(--space-1)',
                        paddingRight: 'var(--space-2)',
                        backgroundColor: 'transparent',
                      }}
                    >
                      <CaretDownIcon width="14" height="14" />
                    </Button>
                  </DropdownMenu.Trigger>
                </Flex>
                <DropdownMenu.Content
                  align="end"
                  style={{
                    backgroundColor: 'var(--color-panel)',
                    border: '1px solid var(--gray-6)',
                    borderRadius: 'var(--radius-3)',
                    boxShadow: 'var(--shadow-4)',
                    padding: 'var(--space-1)',
                    minWidth: '180px',
                    zIndex: 1000,
                  }}
                >
                  <DropdownMenu.Item
                    onSelect={() => setIsCreateTaskDialogOpen(true)}
                    style={{
                      borderRadius: 'var(--radius-2)',
                      padding: 'var(--space-2) var(--space-3)',
                      cursor: 'pointer',
                    }}
                  >
                    Create Task
                  </DropdownMenu.Item>
                  <DropdownMenu.Item
                    onSelect={() => setIsCreateTaskGroupDialogOpen(true)}
                    style={{
                      borderRadius: 'var(--radius-2)',
                      padding: 'var(--space-2) var(--space-3)',
                      cursor: 'pointer',
                    }}
                  >
                    Create Task Group
                  </DropdownMenu.Item>
                </DropdownMenu.Content>
              </DropdownMenu.Root>
            )}
          </Flex>
        </Flex>
      </Box>

      {/* Body with 2-column resizable layout */}
      <Box style={{ flex: 1, overflow: 'hidden' }}>
        <ResizableSplitter
          leftContent={
            <TaskGroupsList
              projectId={project.id}
              projectUuid={project.uuid}
              taskGroups={taskGroups}
              tasks={tasks}
              selectedTaskId={selectedTaskId}
              onSettingsClick={handleTaskGroupSettingsClick}
              onTaskSettingsClick={handleTaskSettingsClick}
              onCreateTaskClick={canEdit ? handleCreateTaskClick : undefined}
              isReadOnly={isReadOnly}
            />
          }
          rightContent={
            <ExecutionsList
              executions={executions}
              isLoading={isLoadingExecutions}
              selectedTaskId={selectedTaskId}
              projectId={project.id}
              paginationData={paginationData}
            />
          }
          initialLeftWidth={20}
          minLeftWidth={20}
          minRightWidth={20}
        />
      </Box>

      {selectedTaskGroup && (
        <TaskGroupSettingsDialog
          open={isTaskGroupSettingsOpen}
          onOpenChange={setIsTaskGroupSettingsOpen}
          taskGroup={selectedTaskGroup}
          onSubmit={handleTaskGroupSettingsSubmit}
          isReadOnly={isReadOnly}
        />
      )}

      {selectedTask && (
        <TaskSettingsDialog
          open={isTaskSettingsOpen}
          onOpenChange={setIsTaskSettingsOpen}
          task={selectedTask}
          onSubmit={handleTaskSettingsSubmit}
          isReadOnly={isReadOnly}
        />
      )}

      {canEdit && (
        <>
          <CreateTaskGroupDialog
            open={isCreateTaskGroupDialogOpen}
            onOpenChange={setIsCreateTaskGroupDialogOpen}
            projectId={project.id}
            onSubmit={handleCreateTaskGroupSubmit}
          />

          <CreateTaskDialog
            open={isCreateTaskDialogOpen}
            onOpenChange={setIsCreateTaskDialogOpen}
            projectId={project.id}
            taskGroupId={selectedTaskGroupForCreate?.id}
            onSubmit={handleCreateTaskSubmit}
          />
        </>
      )}

      <ProjectSettingsDialog
        open={isProjectSettingsOpen}
        onOpenChange={setIsProjectSettingsOpen}
        project={project}
        onSubmit={handleProjectSettingsSubmit}
        isReadOnly={isReadOnly}
      />
    </Flex>
  )
}

