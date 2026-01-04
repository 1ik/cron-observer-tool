'use client'

import { PauseIcon, PlayIcon } from '@radix-ui/react-icons'
import { Box, Flex, IconButton, Text } from '@radix-ui/themes'
import { useState } from 'react'
import { Execution } from '../lib/types/execution'
import { ExecutionItem } from './ExecutionItem'

interface ExecutionsListProps {
  executions: Execution[]
}

export function ExecutionsList({ executions }: ExecutionsListProps) {
  const [isPaused, setIsPaused] = useState(false)

  const handleToggle = () => {
    setIsPaused(!isPaused)
    // TODO: Implement actual pause/play logic for executions
  }

  return (
    <Flex
      direction="column"
      style={{
        width: '100%',
        height: '100%',
        overflow: 'hidden',
      }}
    >
      {/* Sticky header */}
      <Box
        py="3"
        pl="3"
        pr="4"
        style={{
          flexShrink: 0,
          borderBottom: '1px solid var(--gray-6)',
        }}
      >
        <Flex justify="between" align="center">
          <Text size="2" weight="medium" color="gray">
            Executions
          </Text>
          <IconButton
            variant="ghost"
            size="2"
            onClick={handleToggle}
            style={{ cursor: 'pointer' }}
          >
            {isPaused ? (
              <PlayIcon width="16" height="16" />
            ) : (
              <PauseIcon width="16" height="16" />
            )}
          </IconButton>
        </Flex>
      </Box>

      {/* Scrollable body */}
      <Box
        style={{
          flex: 1,
          overflowY: 'auto',
          padding: 'var(--space-3)',
        }}
      >
        {executions.length === 0 ? (
          <Box>
            <Text size="3" color="gray" align="center">
              No executions yet
            </Text>
          </Box>
        ) : (
          <Flex direction="column" gap="3">
            {executions.map((execution) => (
              <ExecutionItem key={execution.id} execution={execution} />
            ))}
          </Flex>
        )}
      </Box>
    </Flex>
  )
}
