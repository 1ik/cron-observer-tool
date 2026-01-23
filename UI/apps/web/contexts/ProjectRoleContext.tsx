'use client'

import { createContext, useContext, ReactNode } from 'react'
import { ProjectUserRole } from '../lib/types/project'

export interface ProjectRoleContextValue {
  /** The current user's role in the project (admin, readonly, or null if super admin / owner) */
  userRole: ProjectUserRole | null
  /** Whether the user can edit (create, update, delete) resources in this project */
  canEdit: boolean
  /** Whether the user can only view resources (readonly role) */
  isReadOnly: boolean
}

const ProjectRoleContext = createContext<ProjectRoleContextValue | undefined>(undefined)

interface ProjectRoleProviderProps {
  children: ReactNode
  userRole: ProjectUserRole | null
}

export function ProjectRoleProvider({ children, userRole }: ProjectRoleProviderProps) {
  // readonly users cannot edit, all other users (admin, null/super admin) can edit
  const isReadOnly = userRole === 'readonly'
  const canEdit = !isReadOnly

  return (
    <ProjectRoleContext.Provider value={{ userRole, canEdit, isReadOnly }}>
      {children}
    </ProjectRoleContext.Provider>
  )
}

export function useProjectRole(): ProjectRoleContextValue {
  const context = useContext(ProjectRoleContext)
  if (context === undefined) {
    // Default to full access if used outside of provider (shouldn't happen in practice)
    return { userRole: null, canEdit: true, isReadOnly: false }
  }
  return context
}

