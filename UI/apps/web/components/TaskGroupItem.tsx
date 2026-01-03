'use client'

import { Box, Card, Flex, Text, Heading, Badge } from '@radix-ui/themes'
import { TaskGroup } from '../lib/types/taskgroup'
import { Task } from '../lib/types/task'
import { TaskItem } from './TaskItem'

interface TaskGroupItemProps {
  taskGroup: TaskGroup
  tasks: Task[]
}

export function TaskGroupItem({ taskGroup, tasks }: TaskGroupItemProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'ACTIVE':
        return 'green'
      case 'PAUSED':
        return 'yellow'
      case 'DISABLED':
        return 'gray'
      default:
        return 'gray'
    }
  }

  return (
    <Card size="2">
      <Flex direction="column" gap="3" p="3">
        <Flex justify="between" align="center">
          <Heading size="4">{taskGroup.name}</Heading>
          <Badge color={getStatusColor(taskGroup.status)} variant="soft" radius="full">
            {taskGroup.status}
          </Badge>
        </Flex>

        {taskGroup.description && (
          <Text size="2" color="gray">
            {taskGroup.description}
          </Text>
        )}

        {taskGroup.start_time && taskGroup.end_time && (
          <Text size="2" color="gray">
            {taskGroup.start_time} - {taskGroup.end_time}
            {taskGroup.timezone && ` (${taskGroup.timezone})`}
          </Text>
        )}

        {tasks.length > 0 && (
          <Box mt="2">
            <Text size="2" weight="medium" color="gray" mb="2">
              Tasks ({tasks.length})
            </Text>
            <Flex direction="column" gap="2">
              {tasks.map((task) => (
                <TaskItem key={task.id} task={task} compact />
              ))}
            </Flex>
          </Box>
        )}
      </Flex>
    </Card>
  )
}

