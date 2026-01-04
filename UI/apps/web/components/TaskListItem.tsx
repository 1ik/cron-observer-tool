'use client'

import { Box, Flex, Text, IconButton } from '@radix-ui/themes'
import { GearIcon } from '@radix-ui/react-icons'
import { useRouter } from 'next/navigation'
import { Task } from '../lib/types/task'

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
  const getStatusDotColor = (status: string) => {
    switch (status) {
      case 'ACTIVE':
        return 'var(--green-9)'
      case 'PAUSED':
        return 'var(--yellow-9)'
      case 'DISABLED':
        return 'var(--gray-9)'
      default:
        return 'var(--gray-9)'
    }
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
          <Box
            style={{
              width: '6px',
              height: '6px',
              borderRadius: '50%',
              backgroundColor: getStatusDotColor(task.status),
              flexShrink: 0,
            }}
          />
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

