import Link from 'next/link'
import { Box, Flex, Heading, Text, Button } from '@radix-ui/themes'

export default function Home() {
  return (
    <Box p="9">
      <Flex direction="column" align="center" gap="6">
        <Heading size="9" align="center">
          Cron Observer
        </Heading>
        <Text size="4" color="gray" align="center">
          Task scheduling and execution tracking system
        </Text>
        <Flex gap="3" mt="4">
          <Link href="/projects" style={{ textDecoration: 'none' }}>
            <Button size="4">View Projects</Button>
          </Link>
        </Flex>
      </Flex>
    </Box>
  )
}

