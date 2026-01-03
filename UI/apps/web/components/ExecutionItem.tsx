'use client'

import { Box, Card, Flex, Text, Badge } from '@radix-ui/themes'
import { Execution } from '../lib/types/execution'

interface ExecutionItemProps {
  execution: Execution
}

export function ExecutionItem({ execution }: ExecutionItemProps) {
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

  const formatDuration = (ms?: number) => {
    if (!ms) return 'N/A'
    if (ms < 1000) return `${ms}ms`
    if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
    return `${(ms / 60000).toFixed(1)}m`
  }

  return (
    <Card size="2">
      <Flex direction="column" gap="2" p="3">
        <Flex justify="between" align="center">
          <Text size="3" weight="medium">
            {execution.task_name}
          </Text>
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
    </Card>
  )
}

