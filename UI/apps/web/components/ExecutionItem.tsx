'use client'

import { useState } from 'react'
import { ChevronDownIcon } from '@radix-ui/react-icons'
import { Box, Card, Flex, Text, Badge } from '@radix-ui/themes'
import { Execution, LogEntry } from '../lib/types/execution'

interface ExecutionItemProps {
  execution: Execution
}

export function ExecutionItem({ execution }: ExecutionItemProps) {
  const [isExpanded, setIsExpanded] = useState(false)

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'SUCCESS':
        return 'green'
      case 'FAILED':
        return 'red'
      case 'RUNNING':
        return 'blue'
      case 'PENDING':
        return 'yellow'
      case 'CANCELLED':
        return 'gray'
      default:
        return 'gray'
    }
  }

  const getLogLevelColor = (level: string) => {
    switch (level) {
      case 'ERROR':
        return 'red'
      case 'WARN':
        return 'yellow'
      case 'INFO':
        return 'blue'
      case 'DEBUG':
        return 'gray'
      default:
        return 'gray'
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  }

  const formatDuration = (ms?: number) => {
    if (!ms) return 'N/A'
    if (ms < 1000) return `${ms}ms`
    if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
    return `${(ms / 60000).toFixed(1)}m`
  }

  const hasLogs = execution.logs && execution.logs.length > 0

  return (
    <Card size="2">
      <Flex direction="column" gap="2">
        {/* Header - clickable to expand */}
        <Box
          onClick={() => hasLogs && setIsExpanded(!isExpanded)}
          style={{
            cursor: hasLogs ? 'pointer' : 'default',
          }}
        >
          <Flex direction="column" gap="2" p="3">
            <Flex justify="between" align="center">
              <Flex align="center" gap="2" style={{ flex: 1 }}>
                {hasLogs && (
                  <ChevronDownIcon
                    width="16"
                    height="16"
                    style={{
                      transition: 'transform 0.2s',
                      transform: isExpanded ? 'rotate(0deg)' : 'rotate(-90deg)',
                      color: 'var(--gray-11)',
                    }}
                  />
                )}
                <Text size="3" weight="medium">
                  {execution.task_name}
                </Text>
              </Flex>
              <Badge color={getStatusColor(execution.status)} variant="soft" radius="full">
                {execution.status}
              </Badge>
            </Flex>

            <Flex gap="4" wrap="wrap">
              <Text size="2" color="gray">
                Started: {formatDate(execution.started_at)}
              </Text>
              {execution.completed_at && (
                <Text size="2" color="gray">
                  Completed: {formatDate(execution.completed_at)}
                </Text>
              )}
              {execution.duration_ms && (
                <Text size="2" color="gray">
                  Duration: {formatDuration(execution.duration_ms)}
                </Text>
              )}
            </Flex>

            {execution.response_status && (
              <Text size="2" color="gray">
                Status Code: {execution.response_status}
              </Text>
            )}

            {execution.error_message && (
              <Box
                p="2"
                style={{
                  backgroundColor: 'var(--red-2)',
                  borderRadius: 'var(--radius-2)',
                }}
              >
                <Text size="2" color="red">
                  Error: {execution.error_message}
                </Text>
              </Box>
            )}
          </Flex>
        </Box>

        {/* Expanded logs section */}
        {isExpanded && hasLogs && (
          <Box
            p="3"
            style={{
              borderTop: '1px solid var(--gray-6)',
              backgroundColor: 'var(--gray-2)',
            }}
          >
            <Text size="2" weight="medium" color="gray" mb="2">
              Logs ({execution.logs!.length})
            </Text>
            <Flex direction="column" gap="1">
              {execution.logs!.map((log) => (
                <Box
                  key={log.id}
                  p="2"
                  style={{
                    backgroundColor: 'var(--gray-1)',
                    borderRadius: 'var(--radius-1)',
                    fontFamily: 'monospace',
                  }}
                >
                  <Flex align="center" gap="2" mb="1">
                    <Text size="1" color="gray">
                      {formatTime(log.timestamp)}
                    </Text>
                    <Badge
                      color={getLogLevelColor(log.level)}
                      variant="soft"
                      size="1"
                    >
                      {log.level}
                    </Badge>
                  </Flex>
                  <Text size="1" style={{ wordBreak: 'break-word' }}>
                    {log.message}
                  </Text>
                </Box>
              ))}
            </Flex>
          </Box>
        )}
      </Flex>
    </Card>
  )
}
