'use client'

import { useFailedExecutionsStats } from '@cron-observer/lib'
import { Box, Badge, Flex, Spinner, Text } from '@radix-ui/themes'

interface FailedExecutionsStatsProps {
  projectId: string | null
  days?: number
}

export function FailedExecutionsStats({ projectId, days = 7 }: FailedExecutionsStatsProps) {
  const { data, isLoading, error } = useFailedExecutionsStats(projectId, days)

  if (!projectId) {
    return null
  }

  if (isLoading) {
    return (
      <Box p="3" style={{ borderBottom: '1px solid var(--gray-6)' }}>
        <Flex align="center" gap="2">
          <Spinner size="2" />
          <Text size="2" color="gray">
            Loading failure statistics...
          </Text>
        </Flex>
      </Box>
    )
  }

  if (error) {
    return (
      <Box p="3" style={{ borderBottom: '1px solid var(--gray-6)' }}>
        <Text size="2" color="red">
          Failed to load statistics
        </Text>
      </Box>
    )
  }

  if (!data || data.stats.length === 0) {
    return (
      <Box p="3" style={{ borderBottom: '1px solid var(--gray-6)' }}>
        <Flex align="center" gap="2">
          <Text size="2" weight="medium">
            Failure Statistics (Last {days} days):
          </Text>
          <Text size="2" color="gray">
            No failures
          </Text>
        </Flex>
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
    <Box p="3" style={{ borderBottom: '1px solid var(--gray-6)' }}>
      <Flex direction="column" gap="2">
        <Flex align="center" gap="2">
          <Text size="2" weight="medium">
            Failure Statistics (Last {days} days):
          </Text>
          <Badge color="red" size="2">
            {data.total} total
          </Badge>
        </Flex>
        <Flex wrap="wrap" gap="2" style={{ flexWrap: 'wrap' }}>
          {data.stats.map((stat) => (
            <Flex key={stat.date} align="center" gap="1">
              <Text size="2" color="gray">
                {formatDate(stat.date)}:
              </Text>
              <Badge color="red" size="1">
                {stat.count}
              </Badge>
            </Flex>
          ))}
        </Flex>
      </Flex>
    </Box>
  )
}

