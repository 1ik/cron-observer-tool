'use client'

import * as Accordion from '@radix-ui/react-accordion'
import { ChevronDownIcon } from '@radix-ui/react-icons'
import { Box, Flex, Text } from '@radix-ui/themes'
import { Task } from '../lib/types/task'
import { TaskGroup } from '../lib/types/taskgroup'
import { TaskListItem } from './TaskListItem'

interface TaskGroupAccordionItemProps {
  taskGroup: TaskGroup
  tasks: Task[]
  selectedTaskId: string | null
  onTaskSelect: (taskId: string) => void
}

export function TaskGroupAccordionItem({
  taskGroup,
  tasks,
  selectedTaskId,
  onTaskSelect,
}: TaskGroupAccordionItemProps) {
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
    <Accordion.Item value={taskGroup.id} style={{ borderBottom: '1px solid var(--gray-6)' }}>
      <Accordion.Header>
        <Accordion.Trigger
          style={{
            width: '100%',
            padding: 'var(--space-3)',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            cursor: 'pointer',
            backgroundColor: 'transparent',
            border: 'none',
            textAlign: 'left',
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
          <Box
            style={{
              width: '8px',
              height: '8px',
              borderRadius: '50%',
              backgroundColor: getStatusDotColor(taskGroup.status),
              flexShrink: 0,
            }}
          />
        </Accordion.Trigger>
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
                isSelected={selectedTaskId === task.id}
                onSelect={() => onTaskSelect(task.id)}
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

