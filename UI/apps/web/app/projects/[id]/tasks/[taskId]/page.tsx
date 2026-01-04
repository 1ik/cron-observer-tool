import { ProjectPageContent } from '../../../../../components/ProjectPageContent'

interface TaskPageProps {
  params: Promise<{ id: string; taskId: string }>
}

export default async function TaskPage({ params }: TaskPageProps) {
  const { id, taskId } = await params

  return <ProjectPageContent projectId={id} selectedTaskId={taskId} />
}

