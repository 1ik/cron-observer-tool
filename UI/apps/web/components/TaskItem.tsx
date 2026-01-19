'use client'

import { Box, Flex, Text, Badge } from '@radix-ui/themes'
import { Task } from '../lib/types/task'

interface TaskItemProps {
  task: Task
  compact?: boolean
}

export function TaskItem({ task, compact = false }: TaskItemProps) {
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

  const getScheduleDisplay = () => {
    if (task.schedule_config.cron_expression) {
      return task.schedule_config.cron_expression
    }
    if (task.schedule_config.time_range) {
      const { start, end, frequency } = task.schedule_config.time_range
      return `${start}-${end} every ${frequency.value}${frequency.unit}`
    }
    return 'No schedule'
  }

  if (compact) {
    return (
      <Box
        p="2"
        style={{
          backgroundColor: 'var(--gray-2)',
          borderRadius: 'var(--radius-2)',
        }}
      >
        <Flex justify="between" align="center">
          <Text size="2" weight="medium">
            {task.name}
          </Text>
          <Badge color={getStatusColor(task.status)} variant="soft" size="1">
            {task.status}
          </Badge>
        </Flex>
      </Box>
    )
  }

  return (
    <Box
      p="3"
      style={{
        border: '1px solid var(--gray-6)',
        borderRadius: 'var(--radius-2)',
      }}
    >
      <Flex direction="column" gap="2">
        <Flex justify="between" align="center">
          <Text size="3" weight="medium">
            {task.name}
          </Text>
          <Badge color={getStatusColor(task.status)} variant="soft" radius="full">
            {task.status}
          </Badge>
        </Flex>

        {task.description && (
          <Text size="2" color="gray">
            {task.description}
          </Text>
        )}

        <Text size="1" color="gray">
          Schedule: {getScheduleDisplay()}
        </Text>
      </Flex>
    </Box>
  )
}

