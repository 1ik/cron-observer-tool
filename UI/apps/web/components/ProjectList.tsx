'use client'

import { useCreateProject, useProjects } from '@cron-observer/lib'
import { Box, Button, Flex, Grid, Heading, Spinner, Text } from '@radix-ui/themes'
import { useState } from 'react'
import { CreateProjectRequest } from '../lib/types/project'
import { CreateProjectDialog } from './CreateProjectDialog'
import { ProjectCard } from './ProjectCard'

export function ProjectList() {
  const { data: projects = [], isLoading, error } = useProjects()
  const createProject = useCreateProject()
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)

  const handleCreateProject = async (data: CreateProjectRequest) => {
    try {
      await createProject.mutateAsync(data)
      setIsCreateDialogOpen(false)
      // TODO: Add success toast/notification
    } catch (error) {
      // TODO: Add error toast/notification
      console.error('Failed to create project:', error)
    }
  }

  // Error state
  if (error) {
    return (
      <Box p="8">
        <Flex direction="column" align="center" gap="4">
          <Heading size="5" color="red">
            Error Loading Projects
          </Heading>
          <Text size="3" color="gray" align="center">
            {error instanceof Error ? error.message : 'Failed to load projects'}
          </Text>
          <Button size="3" onClick={() => window.location.reload()}>
            Retry
          </Button>
        </Flex>
      </Box>
    )
  }

  // Loading state
  if (isLoading) {
    return (
      <Flex justify="center" align="center" style={{ minHeight: '400px' }}>
        <Flex direction="column" gap="3" align="center">
          <Spinner size="3" />
          <Text size="3" color="gray">
            Loading projects...
          </Text>
        </Flex>
      </Flex>
    )
  }

  // Empty state
  if (projects.length === 0) {
    return (
      <Box p="8">
        <Flex direction="column" align="center" gap="4">
          <Heading size="5">No Projects</Heading>
          <Text size="3" color="gray" align="center">
            Get started by creating your first project.
          </Text>
          <Button size="3" onClick={() => setIsCreateDialogOpen(true)}>
            Create Project
          </Button>
        </Flex>
        <CreateProjectDialog
          open={isCreateDialogOpen}
          onOpenChange={setIsCreateDialogOpen}
          onSubmit={handleCreateProject}
        />
      </Box>
    )
  }

  return (
    <Box>
      <Flex justify="between" align="center" mb="6">
        <Heading size="8">Projects</Heading>
        <Button size="3" onClick={() => setIsCreateDialogOpen(true)}>
          Create Project
        </Button>
      </Flex>

      <Grid columns={{ initial: '1', sm: '2', md: '3', lg: '4' }} gap="4">
        {projects.map((project) => (
          <ProjectCard
            key={project.uuid || project.id}
            project={{
              id: project.id || project.uuid || '',
              uuid: project.uuid || project.id || '',
              name: project.name || '',
              description: project.description,
              api_key: project.api_key,
              created_at: project.created_at || new Date().toISOString(),
              updated_at: project.updated_at || new Date().toISOString(),
            }}
          />
        ))}
      </Grid>

      <CreateProjectDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSubmit={handleCreateProject}
      />
    </Box>
  )
}

