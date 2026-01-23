'use client'

import { CaretRightIcon } from '@radix-ui/react-icons'
import { Badge, Box, Flex, Separator, Text } from '@radix-ui/themes'
import { useState } from 'react'
import { Execution } from '../lib/types/execution'

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

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    })
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
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
    <Box
      style={{
        borderBottom: '1px solid var(--gray-4)',
        backgroundColor: 'var(--color-panel)',
        transition: 'background-color 0.15s',
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.backgroundColor = 'var(--gray-2)'
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.backgroundColor = 'var(--color-panel)'
      }}
    >
      {/* Main execution row */}
      <Box
        onClick={() => hasLogs && setIsExpanded(!isExpanded)}
        style={{
          cursor: hasLogs ? 'pointer' : 'default',
          padding: 'var(--space-3)',
        }}
      >
        <Flex direction="column" gap="2">
          {/* First row: Task name, status, and expand icon */}
          <Flex justify="between" align="center" gap="3">
            <Flex align="center" gap="2" style={{ flex: 1, minWidth: 0 }}>
              {hasLogs && (
                <Box
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    width: '16px',
                    height: '16px',
                    flexShrink: 0,
                    transition: 'transform 0.2s ease',
                    transform: isExpanded ? 'rotate(90deg)' : 'rotate(0deg)',
                  }}
                >
                  <CaretRightIcon
                    width="12"
                    height="12"
                    style={{
                      color: 'var(--gray-10)',
                    }}
                  />
                </Box>
              )}
              {!hasLogs && <Box style={{ width: '16px', flexShrink: 0 }} />}
              <Text
                size="2"
                weight="medium"
                style={{
                  fontFamily: 'var(--font-mono)',
                  color: 'var(--gray-12)',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                }}
              >
                {execution.task_name}
              </Text>
            </Flex>
            <Badge color={getStatusColor(execution.status)} variant="soft" size="1" radius="full">
              {execution.status}
            </Badge>
          </Flex>

          {/* Second row: Timestamp and duration */}
          <Flex gap="4" wrap="wrap" align="center" pl="5">
            <Text
              size="1"
              color="gray"
              style={{
                fontFamily: 'var(--font-mono)',
                fontSize: '11px',
              }}
            >
              {formatDate(execution.started_at)}
            </Text>
            {execution.duration_ms && (
              <>
                <Text size="1" color="gray" style={{ fontSize: '11px' }}>
                  •
                </Text>
                <Text
                  size="1"
                  color="gray"
                  style={{
                    fontFamily: 'var(--font-mono)',
                    fontSize: '11px',
                  }}
                >
                  {formatDuration(execution.duration_ms)}
                </Text>
              </>
            )}
            {execution.response_status && (
              <>
                <Text size="1" color="gray" style={{ fontSize: '11px' }}>
                  •
                </Text>
                <Text
                  size="1"
                  color="gray"
                  style={{
                    fontFamily: 'var(--font-mono)',
                    fontSize: '11px',
                  }}
                >
                  HTTP {execution.response_status}
                </Text>
              </>
            )}
          </Flex>

          {/* Error message */}
          {execution.error_message && (
            <Box pl="5">
              <Text
                size="1"
                color="red"
                style={{
                  fontFamily: 'var(--font-mono)',
                  fontSize: '11px',
                }}
              >
                ERROR: {execution.error_message}
              </Text>
            </Box>
          )}
        </Flex>
      </Box>

      {/* Expanded logs section */}
      {isExpanded && hasLogs && (
        <>
          <Separator size="4" />
          <Box
            style={{
              padding: 'var(--space-3)',
              backgroundColor: 'var(--gray-1)',
              borderTop: '1px solid var(--gray-4)',
            }}
          >
            <Text
              size="1"
              weight="medium"
              color="gray"
              mb="3"
              style={{
                fontFamily: 'var(--font-mono)',
                fontSize: '10px',
                textTransform: 'uppercase',
                letterSpacing: '0.5px',
              }}
            >
              Logs ({execution.logs!.length})
            </Text>
            <Flex direction="column" gap="1">
              {execution.logs!.map((log, index) => (
                <Box
                  key={log.id || index}
                  p="2"
                  style={{
                    backgroundColor: 'var(--color-panel-solid)',
                    border: '1px solid var(--gray-4)',
                    borderRadius: '2px',
                    fontFamily: 'var(--font-mono)',
                    fontSize: '11px',
                    lineHeight: '1.5',
                  }}
                >
                  <Flex align="start" gap="2" mb="1">
                    <Text
                      size="1"
                      color="gray"
                      style={{
                        fontFamily: 'var(--font-mono)',
                        fontSize: '10px',
                        flexShrink: 0,
                      }}
                    >
                      {formatTime(log.timestamp)}
                    </Text>
                    {log.level && (
                      <Badge
                        color={getLogLevelColor(log.level)}
                        variant="soft"
                        size="1"
                        radius="full"
                        style={{
                          fontFamily: 'var(--font-mono)',
                          fontSize: '9px',
                          flexShrink: 0,
                        }}
                      >
                        {log.level}
                      </Badge>
                    )}
                  </Flex>
                  <Text
                    size="1"
                    style={{
                      wordBreak: 'break-word',
                      fontFamily: 'var(--font-mono)',
                      fontSize: '11px',
                      color: 'var(--gray-12)',
                      whiteSpace: 'pre-wrap',
                    }}
                  >
                    {log.message}
                  </Text>
                </Box>
              ))}
            </Flex>
          </Box>
        </>
      )}
    </Box>
  )
}
