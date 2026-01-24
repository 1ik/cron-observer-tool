'use client'

import * as Accordion from '@radix-ui/react-accordion'
import { Box, Text } from '@radix-ui/themes'
import { useEffect, useMemo, useState } from 'react'
import { Task } from '../lib/types/task'
import { TaskGroup } from '../lib/types/taskgroup'
import { TaskGroupAccordionItem } from './TaskGroupAccordionItem'
import { TaskListItem } from './TaskListItem'

interface TaskGroupsListProps {
  projectId: string
  projectUuid: string
  taskGroups: TaskGroup[]
  tasks: Task[]
  selectedTaskId?: string | null
  onSettingsClick: (taskGroup: TaskGroup) => void
  onTaskSettingsClick: (task: Task) => void
  onCreateTaskClick?: (taskGroup: TaskGroup) => void
  isReadOnly?: boolean
  taskFailuresMap?: Map<string, number>
}

export function TaskGroupsList({
  projectId,
  projectUuid,
  taskGroups,
  tasks,
  selectedTaskId,
  onSettingsClick,
  onTaskSettingsClick,
  onCreateTaskClick,
  isReadOnly = false,
  taskFailuresMap = new Map(),
}: TaskGroupsListProps) {
  // Create a Set of task group IDs for quick lookup (memoized)
  const taskGroupIds = useMemo(() => new Set(taskGroups.map(tg => tg.id)), [taskGroups])

  // Group tasks by task_group_id
  const tasksByGroup = useMemo(() => {
    const map = new Map<string, Task[]>()
    tasks.forEach((task) => {
      if (task.task_group_id && taskGroupIds.has(task.task_group_id)) {
        const groupTasks = map.get(task.task_group_id) || []
        groupTasks.push(task)
        map.set(task.task_group_id, groupTasks)
      }
    })
    return map
  }, [tasks, taskGroupIds])

  const ungroupedTasks = useMemo(() => {
    return tasks.filter((task) => !task.task_group_id || !taskGroupIds.has(task.task_group_id))
  }, [tasks, taskGroupIds])

  // Find the group ID of the selected task
  const selectedTaskGroupId = useMemo(() => {
    if (!selectedTaskId) return null
    const selectedTask = tasks.find((t) => t.id === selectedTaskId || t.uuid === selectedTaskId)
    if (selectedTask?.task_group_id && taskGroupIds.has(selectedTask.task_group_id)) {
      return selectedTask.task_group_id
    }
    return null
  }, [selectedTaskId, tasks, taskGroupIds])

  // Controlled accordion state - open the group containing the selected task
  const [openGroups, setOpenGroups] = useState<string[]>([])

  // Update open groups when selected task changes
  useEffect(() => {
    if (selectedTaskGroupId) {
      setOpenGroups((prev) => {
        // If the group is not already open, add it
        if (!prev.includes(selectedTaskGroupId)) {
          return [...prev, selectedTaskGroupId]
        }
        return prev
      })
    }
  }, [selectedTaskGroupId])

  return (
    <Box
      style={{
        width: '100%',
        height: '100%',
        overflowY: 'auto',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      {taskGroups.length > 0 && (
        <Accordion.Root
          type="multiple"
          value={openGroups}
          onValueChange={setOpenGroups}
          style={{ width: '100%' }}
        >
          {taskGroups.map((group) => (
            <TaskGroupAccordionItem
              key={group.id}
              taskGroup={group}
              tasks={tasksByGroup.get(group.id) || []}
              projectUuid={projectUuid}
              selectedTaskId={selectedTaskId}
              onSettingsClick={onSettingsClick}
              onTaskSettingsClick={onTaskSettingsClick}
              onCreateTaskClick={onCreateTaskClick}
              isReadOnly={isReadOnly}
              taskFailuresMap={taskFailuresMap}
            />
          ))}
        </Accordion.Root>
      )}

      {ungroupedTasks.length > 0 && (
        <Box p="3" style={{ borderTop: taskGroups.length > 0 ? '1px solid var(--gray-6)' : 'none' }}>
          <Text size="2" weight="medium" color="gray" mb="2">
            Tasks
          </Text>
          <Box style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
            {ungroupedTasks.map((task) => (
              <TaskListItem
                key={task.id}
                task={task}
                projectUuid={projectUuid}
                isSelected={selectedTaskId === task.id || selectedTaskId === task.uuid}
                onSettingsClick={onTaskSettingsClick}
                taskFailuresMap={taskFailuresMap}
              />
            ))}
          </Box>
        </Box>
      )}

      {taskGroups.length === 0 && ungroupedTasks.length === 0 && (
        <Box p="4">
          <Text size="3" color="gray" align="center">
            No task groups or tasks yet
          </Text>
        </Box>
      )}
    </Box>
  )
}
