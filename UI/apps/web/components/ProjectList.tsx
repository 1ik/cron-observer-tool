'use client'

import { updateProject, useCreateProject, useProjects } from '@cron-observer/lib'
import { Box, Button, Flex, Grid, Heading, Spinner, Text } from '@radix-ui/themes'
import { useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { CreateProjectRequest, Project, UpdateProjectRequest } from '../lib/types/project'
import { CreateProjectDialog } from './CreateProjectDialog'
import { ProjectCard } from './ProjectCard'
import { ProjectSettingsDialog } from './ProjectSettingsDialog'

export function ProjectList() {
  const { data: projects = [], isLoading, error } = useProjects()
  const createProject = useCreateProject()
  const queryClient = useQueryClient()
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isProjectSettingsOpen, setIsProjectSettingsOpen] = useState(false)
  const [selectedProject, setSelectedProject] = useState<Project | null>(null)
  const [selectedProjectApiId, setSelectedProjectApiId] = useState<string | null>(null)

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

  const handleProjectSettingsClick = (project: Project, apiProjectId: string) => {
    setSelectedProject(project)
    setSelectedProjectApiId(apiProjectId)
    setIsProjectSettingsOpen(true)
  }

  const handleProjectSettingsSubmit = async (data: UpdateProjectRequest) => {
    if (!selectedProject || !selectedProjectApiId) return
    
    try {
      await updateProject(selectedProjectApiId, data)
      
      // Invalidate queries to refetch projects
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      
      setIsProjectSettingsOpen(false)
      setSelectedProject(null)
      setSelectedProjectApiId(null)
      // TODO: Add success toast/notification
    } catch (error) {
      // TODO: Add error toast/notification
      console.error('Failed to update project:', error)
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
        {projects.map((project) => {
          const mappedProject: Project = {
            id: project.id || project.uuid || '',
            uuid: project.uuid || project.id || '',
            name: project.name || '',
            description: project.description,
            api_key: project.api_key,
            execution_endpoint: project.execution_endpoint,
            alert_emails: (project as Record<string, unknown>).alert_emails && typeof (project as Record<string, unknown>).alert_emails === 'string' 
              ? (project as Record<string, unknown>).alert_emails as string 
              : undefined,
            project_users: (project as Record<string, unknown>).project_users as Project['project_users'],
            created_at: project.created_at || new Date().toISOString(),
            updated_at: project.updated_at || new Date().toISOString(),
          }
          return (
            <ProjectCard
              key={project.uuid || project.id}
              project={mappedProject}
              onSettingsClick={(p) => {
                if (project.id) {
                  handleProjectSettingsClick(p, project.id)
                }
              }}
            />
          )
        })}
      </Grid>

      <CreateProjectDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSubmit={handleCreateProject}
      />

      {selectedProject && (
        <ProjectSettingsDialog
          open={isProjectSettingsOpen}
          onOpenChange={setIsProjectSettingsOpen}
          project={selectedProject}
          onSubmit={handleProjectSettingsSubmit}
        />
      )}
    </Box>
  )
}

