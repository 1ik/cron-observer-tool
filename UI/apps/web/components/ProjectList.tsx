'use client'

import { useState } from 'react'
import { Box, Flex, Grid, Text, Heading, Button, Spinner } from '@radix-ui/themes'
import { ProjectCard } from './ProjectCard'
import { CreateProjectDialog } from './CreateProjectDialog'
import { mockProjects } from '../lib/mocks/projects'
import { CreateProjectRequest } from '../lib/types/project'

export function ProjectList() {
  const [projects] = useState(mockProjects)
  const [isLoading] = useState(false)
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)

  const handleCreateProject = (data: CreateProjectRequest) => {
    // For now, just log the data
    // In the future, this will call the API
    console.log('Creating project:', data)
    // TODO: Add success toast/notification
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
          <ProjectCard key={project.id} project={project} />
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

