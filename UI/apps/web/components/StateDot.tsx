'use client'

import { Box, Tooltip } from '@radix-ui/themes'
import { TaskState } from '../lib/types/task'
import { TaskGroupState } from '../lib/types/taskgroup'

interface StateDotProps {
  state?: TaskState | TaskGroupState
  size?: number
}

/**
 * Component that displays only the state dot (blue/gray)
 * State dot shows system-controlled state (RUNNING/NOT_RUNNING)
 */
export function StateDot({ state, size = 6 }: StateDotProps) {
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
  )
}

