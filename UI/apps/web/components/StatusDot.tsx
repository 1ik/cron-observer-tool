'use client'

import { Box, Tooltip } from '@radix-ui/themes'
import { TaskRuntimeStatus, getStatusDotColor, getStatusTooltip } from '../lib/utils/task-status'

interface StatusDotProps {
  status: TaskRuntimeStatus
  size?: number
}

export function StatusDot({ status, size = 6 }: StatusDotProps) {
  const statusColor = getStatusDotColor(status)
  const statusTooltipText = getStatusTooltip(status)

  return (
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
  )
}
