'use client'

import { CheckIcon, CopyIcon, GearIcon } from '@radix-ui/react-icons'
import { Card, Flex, Heading, IconButton, Text, Tooltip } from '@radix-ui/themes'
import Link from 'next/link'
import { useState } from 'react'
import { Project } from '../lib/types/project'

interface ProjectCardProps {
  project: Project
  onSettingsClick?: (project: Project) => void
}

export function ProjectCard({ project, onSettingsClick }: ProjectCardProps) {
  const [copiedApiKey, setCopiedApiKey] = useState(false)
  const [copiedEndpoint, setCopiedEndpoint] = useState(false)

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

  const copyToClipboard = async (text: string, type: 'apiKey' | 'endpoint') => {
    try {
      await navigator.clipboard.writeText(text)
      if (type === 'apiKey') {
        setCopiedApiKey(true)
        setTimeout(() => setCopiedApiKey(false), 2000)
      } else {
        setCopiedEndpoint(true)
        setTimeout(() => setCopiedEndpoint(false), 2000)
      }
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  return (
    <Link href={`/projects/${project.uuid}`} style={{ textDecoration: 'none' }}>
      <Card
        size="3"
        style={{
          height: '100%',
          cursor: 'pointer',
          transition: 'background-color 0.2s, box-shadow 0.2s, transform 0.2s',
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.backgroundColor = 'var(--gray-2)'
          e.currentTarget.style.boxShadow = 'var(--shadow-4)'
          e.currentTarget.style.transform = 'translateY(-2px)'
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.backgroundColor = ''
          e.currentTarget.style.boxShadow = ''
          e.currentTarget.style.transform = ''
        }}
      >
        <Flex direction="column" gap="3" p="4">
          <Flex justify="between" align="start">
            <Heading size="5" weight="bold" style={{ flex: 1 }}>
              {project.name}
            </Heading>
            {onSettingsClick && (
              <Tooltip content="Project settings">
                <IconButton
                  size="2"
                  variant="ghost"
                  onClick={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    onSettingsClick(project)
                  }}
                  style={{ cursor: 'pointer', flexShrink: 0 }}
                >
                  <GearIcon width="16" height="16" />
                </IconButton>
              </Tooltip>
            )}
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
                <Flex align="center" gap="2">
                  <Text size="2" style={{ fontFamily: 'monospace', flex: 1 }}>
                    {maskApiKey(project.api_key)}
                  </Text>
                  <Tooltip content={copiedApiKey ? 'Copied!' : 'Copy API key'}>
                    <IconButton
                      size="1"
                      variant="ghost"
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                        copyToClipboard(project.api_key!, 'apiKey')
                      }}
                      style={{ cursor: 'pointer' }}
                    >
                      {copiedApiKey ? (
                        <CheckIcon width="14" height="14" />
                      ) : (
                        <CopyIcon width="14" height="14" />
                      )}
                    </IconButton>
                  </Tooltip>
                </Flex>
              </Flex>
            )}

            {project.execution_endpoint && (
              <Flex direction="column" gap="1">
                <Text size="2" weight="medium" color="gray">
                  Execution Endpoint
                </Text>
                <Flex align="center" gap="2">
                  <Text
                    size="2"
                    style={{
                      fontFamily: 'monospace',
                      flex: 1,
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                    }}
                  >
                    {project.execution_endpoint}
                  </Text>
                  <Tooltip content={copiedEndpoint ? 'Copied!' : 'Copy endpoint URL'}>
                    <IconButton
                      size="1"
                      variant="ghost"
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                        copyToClipboard(project.execution_endpoint!, 'endpoint')
                      }}
                      style={{ cursor: 'pointer' }}
                    >
                      {copiedEndpoint ? (
                        <CheckIcon width="14" height="14" />
                      ) : (
                        <CopyIcon width="14" height="14" />
                      )}
                    </IconButton>
                  </Tooltip>
                </Flex>
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

