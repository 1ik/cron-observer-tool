import { Box } from '@radix-ui/themes'
import { ProjectLayout } from '../../../components/ProjectLayout'
import { mockProjects } from '../../../lib/mocks/projects'
import { mockTaskGroups } from '../../../lib/mocks/taskgroups'
import { mockTasks } from '../../../lib/mocks/tasks'
import { mockExecutions } from '../../../lib/mocks/executions'
import { notFound } from 'next/navigation'

interface ProjectPageProps {
  params: Promise<{ id: string }>
}

export default async function ProjectPage({ params }: ProjectPageProps) {
  const { id } = await params

  // Find project by id or uuid
  const project = mockProjects.find(
    (p) => p.id === id || p.uuid === id
  )

  if (!project) {
    notFound()
  }

  // Filter task groups and tasks for this project
  const projectTaskGroups = mockTaskGroups.filter(
    (tg) => tg.project_id === project.id
  )
  const projectTasks = mockTasks.filter(
    (t) => t.project_id === project.id
  )
  const projectExecutions = mockExecutions.filter(
    (e) => {
      const task = projectTasks.find((t) => t.id === e.task_id)
      return task !== undefined
    }
  )

  return (
    <Box style={{ height: '100%', width: '100%' }}>
      <ProjectLayout
        project={project}
        taskGroups={projectTaskGroups}
        tasks={projectTasks}
        executions={projectExecutions}
      />
    </Box>
  )
}

