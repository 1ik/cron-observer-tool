import { ProjectPageContent } from '../../../components/ProjectPageContent'

interface ProjectPageProps {
  params: Promise<{ id: string }>
}

export default async function ProjectPage({ params }: ProjectPageProps) {
  const { id } = await params

  return <ProjectPageContent projectId={id} />
}

