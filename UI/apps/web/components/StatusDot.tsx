'use client'

import { Box, Tooltip } from '@radix-ui/themes'
import { TaskRuntimeStatus, getStatusDotColor, getStatusTooltip } from '../lib/utils/task-status'

interface StatusDotProps {
  status: TaskRuntimeStatus
  size?: number
  tooltip?: string // Optional custom tooltip text
}

export function StatusDot({ status, size = 6, tooltip }: StatusDotProps) {
  const color = getStatusDotColor(status)
  const tooltipText = tooltip || getStatusTooltip(status)

  return (
    <Tooltip content={tooltipText}>
      <Box
        style={{
          width: `${size}px`,
          height: `${size}px`,
          borderRadius: '50%',
          backgroundColor: color,
          flexShrink: 0,
          cursor: 'help',
        }}
      />
    </Tooltip>
  )
}

