'use client'

import * as Accordion from '@radix-ui/react-accordion'
import { ChevronDownIcon, GearIcon } from '@radix-ui/react-icons'
import { Box, Flex, IconButton, Text } from '@radix-ui/themes'
import { Task } from '../lib/types/task'
import { TaskGroup } from '../lib/types/taskgroup'
import { TaskListItem } from './TaskListItem'
import { TaskRuntimeStatus, getStatusDotColor, getStatusTooltip } from '../lib/utils/task-status'
import { StatusDot } from './StatusDot'

interface TaskGroupAccordionItemProps {
  taskGroup: TaskGroup
  tasks: Task[]
  projectUuid: string
  selectedTaskId?: string | null
  onSettingsClick: (taskGroup: TaskGroup) => void
  onTaskSettingsClick: (task: Task) => void
}

export function TaskGroupAccordionItem({
  taskGroup,
  tasks,
  projectUuid,
  selectedTaskId,
  onSettingsClick,
  onTaskSettingsClick,
}: TaskGroupAccordionItemProps) {
  // Map TaskGroup status to TaskRuntimeStatus for consistency
  const getTaskGroupRuntimeStatus = (status: string): TaskRuntimeStatus => {
    switch (status) {
      case 'ACTIVE':
        return 'success' // Task groups are either active (success) or paused
      case 'PAUSED':
        return 'paused'
      default:
        return 'not-running'
    }
  }

  const runtimeStatus = getTaskGroupRuntimeStatus(taskGroup.status)

  const handleSettingsClick = (e: React.MouseEvent) => {
    e.stopPropagation()
    onSettingsClick(taskGroup)
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
            <StatusDot
              status={runtimeStatus}
              size={8}
              tooltip={
                taskGroup.status === 'ACTIVE'
                  ? 'Task group is active'
                  : taskGroup.status === 'PAUSED'
                    ? 'Task group is paused'
                    : 'Task group is disabled'
              }
            />
          </Accordion.Trigger>
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

