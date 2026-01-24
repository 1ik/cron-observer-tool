'use client'

import { useExecutionStats } from '@cron-observer/lib'
import { Box, Flex, Spinner, Table, Text } from '@radix-ui/themes'

interface ExecutionStatsTableProps {
  projectId: string | null
  days?: number
}

export function ExecutionStatsTable({ projectId, days = 7 }: ExecutionStatsTableProps) {
  const { data, isLoading, error } = useExecutionStats(projectId, days)

  if (!projectId) {
    return null
  }

  if (isLoading) {
    return (
      <Box p="5" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '200px' }}>
        <Flex direction="column" gap="3" align="center">
          <Spinner size="3" />
          <Text size="2" color="gray">
            Loading execution statistics...
          </Text>
        </Flex>
      </Box>
    )
  }

  if (error) {
    return (
      <Box p="5">
        <Text size="2" color="red">
          Failed to load execution statistics
        </Text>
      </Box>
    )
  }

  if (!data || !data.stats || data.stats.length === 0) {
    return (
      <Box p="5">
        <Text size="2" color="gray">
          No execution statistics available
        </Text>
      </Box>
    )
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr + 'T00:00:00Z')
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  }

  return (
    <Box p="3" style={{ overflow: 'auto' }}>
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.ColumnHeaderCell>Date</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell style={{ textAlign: 'right' }}>Failures</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell style={{ textAlign: 'right' }}>Success</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell style={{ textAlign: 'right' }}>Total</Table.ColumnHeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {data.stats.map((stat) => (
            <Table.Row key={stat.date}>
              <Table.Cell>
                <Text size="2" style={{ fontFamily: 'var(--font-mono)' }}>
                  {formatDate(stat.date)}
                </Text>
              </Table.Cell>
              <Table.Cell style={{ textAlign: 'right' }}>
                <Text size="2" color="red" weight="medium">
                  {stat.failures}
                </Text>
              </Table.Cell>
              <Table.Cell style={{ textAlign: 'right' }}>
                <Text size="2" color="green" weight="medium">
                  {stat.success}
                </Text>
              </Table.Cell>
              <Table.Cell style={{ textAlign: 'right' }}>
                <Text size="2" weight="medium">
                  {stat.total}
                </Text>
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table.Root>
    </Box>
  )
}

