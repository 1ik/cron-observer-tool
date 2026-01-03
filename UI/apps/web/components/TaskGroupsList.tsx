'use client'

import { Box, Flex, Heading, Text } from '@radix-ui/themes'
import { TaskGroup } from '../lib/types/taskgroup'
import { Task } from '../lib/types/task'
import { TaskGroupItem } from './TaskGroupItem'
import { TaskItem } from './TaskItem'

interface TaskGroupsListProps {
  projectId: string
  taskGroups: TaskGroup[]
  tasks: Task[]
}

export function TaskGroupsList({
  projectId,
  taskGroups,
  tasks,
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
    <Box p="4">
      <Heading size="6" mb="4">
        Task Groups & Tasks
      </Heading>

      {/* Task Groups */}
      {taskGroups.length > 0 && (
        <Flex direction="column" gap="3" mb="6">
          {taskGroups.map((group) => (
            <TaskGroupItem
              key={group.id}
              taskGroup={group}
              tasks={tasksByGroup.get(group.id) || []}
            />
          ))}
        </Flex>
      )}

      {/* Ungrouped Tasks */}
      {ungroupedTasks.length > 0 && (
        <Box>
          <Heading size="5" mb="3">
            Tasks
          </Heading>
          <Flex direction="column" gap="2">
            {ungroupedTasks.map((task) => (
              <TaskItem key={task.id} task={task} />
            ))}
          </Flex>
        </Box>
      )}

      {/* Empty state */}
      {taskGroups.length === 0 && ungroupedTasks.length === 0 && (
        <Box>
          <Text size="3" color="gray" align="center">
            No task groups or tasks yet
          </Text>
        </Box>
      )}
    </Box>
  )
}

