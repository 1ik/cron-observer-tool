import { Box, Flex, Heading, Text } from '@radix-ui/themes'

export default function Home() {
  return (
    <Box p="9">
      <Flex direction="column" align="center" gap="4">
        <Heading size="9" align="center">
          Cron Observer
        </Heading>
        <Text size="4" color="gray" align="center">
          Task scheduling and execution tracking system
        </Text>
      </Flex>
    </Box>
  )
}

