'use client'

import * as Accordion from '@radix-ui/react-accordion'
import { ChevronDownIcon, GearIcon, PlusIcon } from '@radix-ui/react-icons'
import { Box, Flex, IconButton, Text } from '@radix-ui/themes'
import { Task } from '../lib/types/task'
import { TaskGroup } from '../lib/types/taskgroup'
import { StatusAndStateDots } from './StatusAndStateDots'
import { TaskListItem } from './TaskListItem'

interface TaskGroupAccordionItemProps {
  taskGroup: TaskGroup
  tasks: Task[]
  projectUuid: string
  selectedTaskId?: string | null
  onSettingsClick: (taskGroup: TaskGroup) => void
  onTaskSettingsClick: (task: Task) => void
  onCreateTaskClick?: (taskGroup: TaskGroup) => void
  isReadOnly?: boolean
}

export function TaskGroupAccordionItem({
  taskGroup,
  tasks,
  projectUuid,
  selectedTaskId,
  onSettingsClick,
  onTaskSettingsClick,
  onCreateTaskClick,
  isReadOnly = false,
}: TaskGroupAccordionItemProps) {
  const handleSettingsClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    onSettingsClick(taskGroup)
  }

  const handleCreateTaskClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    onCreateTaskClick?.(taskGroup)
  }

  return (
    <Accordion.Item value={taskGroup.id} style={{ borderBottom: '1px solid var(--gray-6)' }}>
      <Accordion.Header>
        <Flex
          align="center"
          style={{
            width: '100%',
            padding: 'var(--space-3)',
          }}
        >
          <Accordion.Trigger
            style={{
              flex: 1,
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              cursor: 'pointer',
              backgroundColor: 'transparent',
              border: 'none',
              textAlign: 'left',
              padding: 0,
            }}
          >
            <Flex align="center" gap="2" style={{ flex: 1 }}>
              <ChevronDownIcon
                width="16"
                height="16"
                style={{
                  transition: 'transform 0.2s',
                  transform: 'rotate(-90deg)',
                }}
                className="accordion-chevron"
              />
              <Text size="3" weight="medium">
                {taskGroup.name}
              </Text>
              {tasks.length > 0 && (
                <Text size="2" color="gray">
                  ({tasks.length})
                </Text>
              )}
            </Flex>
            <StatusAndStateDots
              status={taskGroup.status}
              state={taskGroup.state}
              size={8}
            />
          </Accordion.Trigger>
          {onCreateTaskClick && !isReadOnly && (
            <IconButton
              variant="ghost"
              size="1"
              onClick={handleCreateTaskClick}
              style={{ cursor: 'pointer', marginLeft: 'var(--space-2)' }}
            >
              <PlusIcon width="14" height="14" />
            </IconButton>
          )}
          <IconButton
            variant="ghost"
            size="1"
            onClick={handleSettingsClick}
            style={{ cursor: 'pointer', marginLeft: 'var(--space-2)' }}
          >
            <GearIcon width="14" height="14" />
          </IconButton>
        </Flex>
      </Accordion.Header>
      <Accordion.Content
        style={{
          padding: '0 var(--space-3) var(--space-3) var(--space-3)',
        }}
      >
        {tasks.length > 0 ? (
          <Box style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
                  {tasks.map((task) => (
                    <TaskListItem
                      key={task.id}
                      task={task}
                      projectUuid={projectUuid}
                      isSelected={selectedTaskId === task.id || selectedTaskId === task.uuid}
                      onSettingsClick={onTaskSettingsClick}
                    />
                  ))}
          </Box>
        ) : (
          <Text size="2" color="gray" style={{ padding: 'var(--space-2)' }}>
            No tasks in this group
          </Text>
        )}
      </Accordion.Content>
    </Accordion.Item>
  )
}

