'use client'

import { GearIcon } from '@radix-ui/react-icons'
import { Box, Flex, IconButton, Text } from '@radix-ui/themes'
import { useRouter } from 'next/navigation'
import { Task } from '../lib/types/task'
import { getTaskRuntimeStatus } from '../lib/utils/task-status'
import { StateDot } from './StateDot'
import { StatusDot } from './StatusDot'

interface TaskListItemProps {
  task: Task
  projectUuid: string
  isSelected: boolean
  onSettingsClick?: (task: Task) => void
}

export function TaskListItem({ task, projectUuid, isSelected, onSettingsClick }: TaskListItemProps) {
  const router = useRouter()

  const handleClick = () => {
    router.push(`/projects/${projectUuid}/tasks/${task.uuid}`)
  }

  const handleSettingsClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    onSettingsClick?.(task)
  }

  return (
    <Box
      onClick={handleClick}
      p="2"
      style={{
        cursor: 'pointer',
        backgroundColor: isSelected ? 'var(--gray-4)' : 'transparent',
        borderRadius: 'var(--radius-2)',
        border: isSelected ? '1px solid var(--gray-7)' : '1px solid transparent',
        transition: 'background-color 0.2s, border-color 0.2s',
      }}
      onMouseEnter={(e) => {
        if (!isSelected) {
          e.currentTarget.style.backgroundColor = 'var(--gray-3)'
        }
      }}
      onMouseLeave={(e) => {
        if (!isSelected) {
          e.currentTarget.style.backgroundColor = 'transparent'
        }
      }}
    >
      <Flex align="center" gap="2" justify="between">
        <Text
          size="2"
          weight={isSelected ? 'bold' : 'regular'}
          style={{
            flex: 1,
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}
        >
          {task.name}
        </Text>
        <Flex align="center" gap="2">
          {onSettingsClick && (
            <IconButton
              variant="ghost"
              size="1"
              onClick={handleSettingsClick}
              style={{ cursor: 'pointer' }}
            >
              <GearIcon width="14" height="14" />
            </IconButton>
          )}
          {task.task_group_id ? (
            // Task belongs to a group: show only state dot (inherits from group)
            <StateDot state={task.state} size={6} />
          ) : (
            // Task has no group: show status dot (user-controlled status)
            <StatusDot status={getTaskRuntimeStatus(task)} size={6} />
          )}
        </Flex>
      </Flex>
      {task.description && (
        <Text
          size="1"
          color="gray"
          style={{
            marginTop: 'var(--space-1)',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}
        >
          {task.description}
        </Text>
      )}
    </Box>
  )
}

