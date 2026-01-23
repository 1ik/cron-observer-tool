'use client'

import { Box, Flex, Tooltip } from '@radix-ui/themes'
import { TaskState } from '../lib/types/task'
import { TaskGroupState } from '../lib/types/taskgroup'
import { TaskRuntimeStatus, getStatusDotColor, getStatusTooltip } from '../lib/utils/task-status'

interface StatusAndStateDotsProps {
  status: TaskRuntimeStatus
  state?: TaskState | TaskGroupState
  size?: number
  tooltip?: string
}

/**
 * Component that displays both status dot (green/yellow/gray) and state dot (blue/gray)
 * Status dot shows user-controlled status (ACTIVE/DISABLED)
 * State dot shows system-controlled state (RUNNING/NOT_RUNNING)
 */
export function StatusAndStateDots({ status, state, size = 6, tooltip }: StatusAndStateDotsProps) {
  const statusColor = getStatusDotColor(status)
  const statusTooltipText = tooltip || getStatusTooltip(status)

  // State dot colors
  const getStateDotColor = (state?: TaskState | TaskGroupState): string => {
    if (state === 'RUNNING') {
      return 'var(--blue-9)'
    }
    return 'var(--gray-9)'
  }

  const getStateTooltip = (state?: TaskState | TaskGroupState): string => {
    if (state === 'RUNNING') {
      return 'Running (within time window)'
    }
    return 'Not running (outside time window)'
  }

  const stateColor = getStateDotColor(state)
  const stateTooltipText = getStateTooltip(state)

  return (
    <Flex align="center" gap="1">
      {/* Status dot (green/yellow/gray) */}
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
      {/* State dot (blue/gray) - only show if state is provided */}
      {state !== undefined && (
        <Tooltip content={stateTooltipText}>
          <Box
            style={{
              width: `${size}px`,
              height: `${size}px`,
              borderRadius: '50%',
              backgroundColor: stateColor,
              flexShrink: 0,
              cursor: 'help',
            }}
          />
        </Tooltip>
      )}
    </Flex>
  )
}

