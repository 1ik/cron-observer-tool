'use client'

import Link from 'next/link'
import { Box, Flex, Text } from '@radix-ui/themes'
import { Project } from '../lib/types/project'
import { TaskGroup } from '../lib/types/taskgroup'
import { Task } from '../lib/types/task'
import { Execution } from '../lib/types/execution'
import { TaskGroupsList } from './TaskGroupsList'
import { ExecutionsList } from './ExecutionsList'
import { ResizableSplitter } from './ResizableSplitter'

interface ProjectLayoutProps {
  project: Project
  taskGroups: TaskGroup[]
  tasks: Task[]
  executions: Execution[]
}

export function ProjectLayout({
  project,
  taskGroups,
  tasks,
  executions,
}: ProjectLayoutProps) {
  return (
    <Flex direction="column" height="100%" width="100%" overflow="hidden">
      {/* Top separator */}
      <Box
        style={{
          width: '100%',
          height: '1px',
          backgroundColor: 'var(--gray-6)',
          flexShrink: 0,
        }}
      />
      
      {/* Header with breadcrumb */}
      <Box
        p="4"
        style={{
          flexShrink: 0,
          borderBottom: '1px solid var(--gray-6)',
        }}
      >
        <Flex align="center" gap="2">
          <Link href="/projects" style={{ textDecoration: 'none' }}>
            <Text size="2" color="gray">
              Projects
            </Text>
          </Link>
          <Text size="2" color="gray">
            &gt;
          </Text>
          <Text size="2" color="gray">
            {project.name}
          </Text>
        </Flex>
      </Box>

      {/* Body with 2-column resizable layout */}
      <Box style={{ flex: 1, overflow: 'hidden' }}>
        <ResizableSplitter
          leftContent={
            <TaskGroupsList
              projectId={project.id}
              taskGroups={taskGroups}
              tasks={tasks}
            />
          }
          rightContent={<ExecutionsList executions={executions} />}
          initialLeftWidth={20}
          minLeftWidth={20}
          minRightWidth={20}
        />
      </Box>
    </Flex>
  )
}

