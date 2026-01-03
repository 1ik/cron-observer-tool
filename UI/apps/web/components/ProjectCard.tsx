'use client'

import Link from 'next/link'
import { Card, Box, Flex, Text, Heading, Badge } from '@radix-ui/themes'
import { Project } from '../lib/types/project'

interface ProjectCardProps {
  project: Project
}

export function ProjectCard({ project }: ProjectCardProps) {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  const maskApiKey = (apiKey: string) => {
    if (apiKey.length <= 12) return apiKey
    return `${apiKey.substring(0, 8)}...${apiKey.substring(apiKey.length - 4)}`
  }

  return (
    <Link href={`/projects/${project.uuid}`} style={{ textDecoration: 'none' }}>
      <Card size="3" style={{ height: '100%', cursor: 'pointer' }}>
        <Flex direction="column" gap="3" p="4">
          <Flex justify="between" align="start">
            <Heading size="5" weight="bold">
              {project.name}
            </Heading>
            <Badge color="blue" variant="soft" radius="full">
              Active
            </Badge>
          </Flex>

          {project.description && (
            <Text size="3" color="gray">
              {project.description}
            </Text>
          )}

          <Flex direction="column" gap="2" mt="2">
            {project.api_key && (
              <Flex direction="column" gap="1">
                <Text size="2" weight="medium" color="gray">
                  API Key
                </Text>
                <Text size="2" style={{ fontFamily: 'monospace' }}>
                  {maskApiKey(project.api_key)}
                </Text>
              </Flex>
            )}

            <Flex direction="column" gap="1">
              <Text size="2" weight="medium" color="gray">
                Created
              </Text>
              <Text size="2" color="gray">
                {formatDate(project.created_at)}
              </Text>
            </Flex>
          </Flex>
        </Flex>
      </Card>
    </Link>
  )
}

