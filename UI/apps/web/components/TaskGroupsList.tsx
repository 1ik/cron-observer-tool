'use client'

import * as Accordion from '@radix-ui/react-accordion'
import { Box, Text } from '@radix-ui/themes'
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
}

export function TaskGroupsList({
  projectId,
  projectUuid,
  taskGroups,
  tasks,
  selectedTaskId,
  onSettingsClick,
}: TaskGroupsListProps) {
  // Group tasks by task_group_id
  const tasksByGroup = new Map<string, Task[]>()
  const ungroupedTasks: Task[] = []

  tasks.forEach((task) => {
    if (task.task_group_id) {
      const groupTasks = tasksByGroup.get(task.task_group_id) || []
      groupTasks.push(task)
      tasksByGroup.set(task.task_group_id, groupTasks)
    } else {
      ungroupedTasks.push(task)
    }
  })


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
        <Accordion.Root type="multiple" style={{ width: '100%' }}>
          {taskGroups.map((group) => (
            <TaskGroupAccordionItem
              key={group.id}
              taskGroup={group}
              tasks={tasksByGroup.get(group.id) || []}
              projectUuid={projectUuid}
              selectedTaskId={selectedTaskId}
              onSettingsClick={onSettingsClick}
            />
          ))}
        </Accordion.Root>
      )}

      {ungroupedTasks.length > 0 && (
        <Box p="3" style={{ borderTop: '1px solid var(--gray-6)' }}>
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
