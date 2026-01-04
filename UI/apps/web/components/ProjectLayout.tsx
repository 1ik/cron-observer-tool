'use client'

import { useState } from 'react'
import { Box, Flex, Text } from '@radix-ui/themes'
import Link from 'next/link'
import { Execution } from '../lib/types/execution'
import { Project } from '../lib/types/project'
import { Task, UpdateTaskRequest } from '../lib/types/task'
import { TaskGroup, UpdateTaskGroupRequest } from '../lib/types/taskgroup'
import { ExecutionsList } from './ExecutionsList'
import { ResizableSplitter } from './ResizableSplitter'
import { TaskGroupsList } from './TaskGroupsList'
import { TaskGroupSettingsDialog } from './TaskGroupSettingsDialog'
import { TaskSettingsDialog } from './TaskSettingsDialog'

interface ProjectLayoutProps {
  project: Project
  taskGroups: TaskGroup[]
  tasks: Task[]
  executions: Execution[]
  selectedTaskId?: string | null
}

export function ProjectLayout({
  project,
  taskGroups,
  tasks,
  executions,
  selectedTaskId,
}: ProjectLayoutProps) {
  const [selectedTaskGroup, setSelectedTaskGroup] = useState<TaskGroup | null>(null)
  const [isTaskGroupSettingsOpen, setIsTaskGroupSettingsOpen] = useState(false)
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [isTaskSettingsOpen, setIsTaskSettingsOpen] = useState(false)

  const handleTaskGroupSettingsClick = (taskGroup: TaskGroup) => {
    setSelectedTaskGroup(taskGroup)
    setIsTaskGroupSettingsOpen(true)
  }

  const handleTaskGroupSettingsSubmit = (data: UpdateTaskGroupRequest) => {
    if (!selectedTaskGroup) return
    // TODO: Implement API call to update task group
    console.log('Updating task group:', selectedTaskGroup.id, data)
    // TODO: Add success toast/notification
    setIsTaskGroupSettingsOpen(false)
  }

  const handleTaskSettingsClick = (task: Task) => {
    setSelectedTask(task)
    setIsTaskSettingsOpen(true)
  }

  const handleTaskSettingsSubmit = (data: UpdateTaskRequest) => {
    if (!selectedTask) return
    // TODO: Implement API call to update task
    console.log('Updating task:', selectedTask.id, data)
    // TODO: Add success toast/notification
    setIsTaskSettingsOpen(false)
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
      
      {/* Header with breadcrumb */}
      <Box
        p="4"
        style={{
          flexShrink: 0,
          borderBottom: '1px solid var(--gray-6)',
        }}
      >
        <Flex align="center" gap="2">
          <Link href="/projects" style={{ textDecoration: 'none' }}>
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
            />
          }
          rightContent={<ExecutionsList executions={executions} />}
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
        />
      )}

      {selectedTask && (
        <TaskSettingsDialog
          open={isTaskSettingsOpen}
          onOpenChange={setIsTaskSettingsOpen}
          task={selectedTask}
          onSubmit={handleTaskSettingsSubmit}
        />
      )}
    </Flex>
  )
}

