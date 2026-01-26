'use client'

import { Box, Flex, Tooltip } from '@radix-ui/themes'
import { TaskState, TaskStatus } from '../lib/types/task'
import { TaskGroupState, TaskGroupStatus } from '../lib/types/taskgroup'

interface StatusAndStateDotsProps {
  status: TaskStatus | TaskGroupStatus
  state?: TaskState | TaskGroupState
  size?: number
  tooltip?: string
}

/**
 * Component that displays status dot (green/gray) and optionally state dot (blue)
 * Status dot shows user-controlled status: Green = ACTIVE, Grey = DISABLED
 * State dot shows system-controlled state: Blue = RUNNING (only shown when RUNNING)
 */
export function StatusAndStateDots({ status, state, size = 6, tooltip }: StatusAndStateDotsProps) {
  // Status dot: Green for ACTIVE, Grey for DISABLED
  const getStatusColor = (status: TaskStatus | TaskGroupStatus): string => {
    if (status === 'ACTIVE') {
      return 'var(--green-9)'
    }
    return 'var(--gray-9)'
  }

  const getStatusTooltip = (status: TaskStatus | TaskGroupStatus): string => {
    if (status === 'ACTIVE') {
      return 'Active'
    }
    return 'Disabled'
  }

  const statusColor = getStatusColor(status)
  const statusTooltipText = tooltip || getStatusTooltip(status)

  // State dot: Only show when RUNNING (blue)
  // Temporarily disabled - blue marker for running within time window
  // const shouldShowStateDot = state === 'RUNNING'

  return (
    <Flex align="center" gap="1">
      {/* Status dot (green/gray) */}
      <Tooltip content={statusTooltipText}>
        <Box
          style={{
            width: `${size}px`,
            height: `${size}px`,
            borderRadius: '50%',
            backgroundColor: statusColor,
            flexShrink: 0,
            cursor: 'help',
          }}
        />
      </Tooltip>
      {/* State dot (blue) - only show when RUNNING */}
      {/* Temporarily disabled - blue marker for running within time window */}
      {/* {shouldShowStateDot && (
        <Tooltip content="Running (within time window)">
          <Box
            style={{
              width: `${size}px`,
              height: `${size}px`,
              borderRadius: '50%',
              backgroundColor: 'var(--blue-9)',
              flexShrink: 0,
              cursor: 'help',
            }}
          />
        </Tooltip>
      )} */}
    </Flex>
  )
}

