import { Box } from '@radix-ui/themes'
import { notFound } from 'next/navigation'
import { ProjectLayout } from '../../../../../components/ProjectLayout'
import { mockExecutions } from '../../../../../lib/mocks/executions'
import { mockProjects } from '../../../../../lib/mocks/projects'
import { mockTaskGroups } from '../../../../../lib/mocks/taskgroups'
import { mockTasks } from '../../../../../lib/mocks/tasks'

interface TaskPageProps {
  params: Promise<{ id: string; taskId: string }>
}

export default async function TaskPage({ params }: TaskPageProps) {
  const { id, taskId } = await params

  // Find project by id or uuid
  const project = mockProjects.find(
    (p) => p.id === id || p.uuid === id
  )

  if (!project) {
    notFound()
  }

  // Find task by id or uuid
  const task = mockTasks.find(
    (t) => (t.id === taskId || t.uuid === taskId) && t.project_id === project.id
  )

  if (!task) {
    notFound()
  }

  // Filter task groups and tasks for this project
  const projectTaskGroups = mockTaskGroups.filter(
    (tg) => tg.project_id === project.id
  )
  const projectTasks = mockTasks.filter(
    (t) => t.project_id === project.id
  )
  
  // Filter executions for the selected task
  const projectExecutions = mockExecutions.filter(
    (e) => e.task_id === task.id || e.task_uuid === task.uuid
  )

  return (
    <Box style={{ height: '100%', width: '100%' }}>
      <ProjectLayout
        project={project}
        taskGroups={projectTaskGroups}
        tasks={projectTasks}
        executions={projectExecutions}
        selectedTaskId={task.uuid}
      />
    </Box>
  )
}

