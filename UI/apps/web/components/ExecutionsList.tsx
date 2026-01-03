'use client'

import { Box, Flex, Heading, Text } from '@radix-ui/themes'
import { Execution } from '../lib/types/execution'
import { ExecutionItem } from './ExecutionItem'

interface ExecutionsListProps {
  executions: Execution[]
}

export function ExecutionsList({ executions }: ExecutionsListProps) {
  return (
    <Box p="4">
      <Heading size="6" mb="4">
        Executions
      </Heading>

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
  )
}

