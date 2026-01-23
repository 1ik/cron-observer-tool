'use client'

import { CheckIcon, ChevronDownIcon, PlusIcon, TrashIcon } from '@radix-ui/react-icons'
import * as Select from '@radix-ui/react-select'
import { Box, Button, Flex, IconButton, Text, TextField } from '@radix-ui/themes'
import { useState } from 'react'
import { FieldErrors } from 'react-hook-form'
import { ProjectUser, ProjectUserRole } from '../lib/types/project'

interface ProjectUsersTabProps {
  projectUsers: ProjectUser[]
  onUsersChange: (users: ProjectUser[]) => void
  errors?: FieldErrors
}

interface ProjectUserRow extends ProjectUser {
  isNew?: boolean
  tempId?: string
}

export function ProjectUsersTab({ projectUsers, onUsersChange, errors }: ProjectUsersTabProps) {
  const [users, setUsers] = useState<ProjectUserRow[]>(projectUsers || [])

  const handleAddUser = () => {
    const newUser: ProjectUserRow = {
      email: '',
      role: 'readonly',
      isNew: true,
      tempId: `temp-${Date.now()}`,
    }
    const updatedUsers = [newUser, ...users]
    setUsers(updatedUsers)
  }

  const handleDeleteUser = (index: number) => {
    const updatedUsers = users.filter((_, i) => i !== index)
    setUsers(updatedUsers)
    onUsersChange(updatedUsers.map(({ isNew, tempId, ...user }) => user))
  }

  const handleEmailChange = (index: number, email: string) => {
    const updatedUsers = [...users]
    updatedUsers[index] = { ...updatedUsers[index], email }
    setUsers(updatedUsers)
    onUsersChange(updatedUsers.map(({ isNew, tempId, ...user }) => user))
  }

  const handleRoleChange = (index: number, role: ProjectUserRole) => {
    const updatedUsers = [...users]
    updatedUsers[index] = { ...updatedUsers[index], role }
    setUsers(updatedUsers)
    onUsersChange(updatedUsers.map(({ isNew, tempId, ...user }) => user))
  }

  return (
    <Flex direction="column" gap="4">
      {/* Header */}
      <Flex justify="between" align="center">
        <Box>
          <Text size="3" weight="medium">
            Project Users
          </Text>
          <Text size="2" color="gray">
            Manage user access and roles for this project
          </Text>
        </Box>
        <Button type="button" size="2" variant="soft" onClick={handleAddUser}>
          <PlusIcon />
          Add User
        </Button>
      </Flex>

      {/* Table */}
      {users.length > 0 ? (
        <Box
          style={{
            border: '1px solid var(--gray-6)',
            borderRadius: 'var(--radius-3)',
            overflow: 'hidden',
          }}
        >
          {/* Table Header */}
          <Flex
            p="3"
            gap="3"
            style={{
              backgroundColor: 'var(--gray-2)',
              borderBottom: '1px solid var(--gray-6)',
            }}
          >
            <Box style={{ flex: '1' }}>
              <Text size="2" weight="medium" color="gray">
                Email
              </Text>
            </Box>
            <Box style={{ width: '200px' }}>
              <Text size="2" weight="medium" color="gray">
                Role
              </Text>
            </Box>
            <Box style={{ width: '60px' }}>
              <Text size="2" weight="medium" color="gray">
                Actions
              </Text>
            </Box>
          </Flex>

          {/* Table Rows */}
          {users.map((user, index) => (
            <Flex
              key={user.tempId || user.email || index}
              p="3"
              gap="3"
              align="center"
              style={{
                borderBottom: index < users.length - 1 ? '1px solid var(--gray-6)' : 'none',
                backgroundColor: user.isNew ? 'var(--accent-2)' : 'transparent',
              }}
            >
              {/* Email Column */}
              <Box style={{ flex: '1' }}>
                {user.isNew ? (
                  <TextField.Root
                    type="email"
                    placeholder="user@example.com"
                    value={user.email}
                    onChange={(e) => handleEmailChange(index, e.target.value)}
                    size="2"
                  />
                ) : (
                  <Text size="2">{user.email}</Text>
                )}
              </Box>

              {/* Role Column */}
              <Box style={{ width: '200px' }}>
                <Select.Root
                  value={user.role}
                  onValueChange={(value) => handleRoleChange(index, value as ProjectUserRole)}
                >
                  <Select.Trigger
                    style={{
                      display: 'inline-flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                      gap: '8px',
                      padding: '6px 12px',
                      fontSize: '14px',
                      lineHeight: '20px',
                      fontWeight: '400',
                      border: '1px solid var(--gray-7)',
                      borderRadius: 'var(--radius-2)',
                      backgroundColor: 'var(--color-surface)',
                      color: 'var(--gray-12)',
                      cursor: 'pointer',
                      width: '100%',
                    }}
                  >
                    <Select.Value />
                    <Select.Icon>
                      <ChevronDownIcon />
                    </Select.Icon>
                  </Select.Trigger>
                  <Select.Portal>
                    <Select.Content
                      position="popper"
                      style={{
                        backgroundColor: 'var(--color-panel-solid)',
                        border: '1px solid var(--gray-7)',
                        borderRadius: 'var(--radius-3)',
                        boxShadow: '0 8px 16px rgba(0, 0, 0, 0.15)',
                        minWidth: '200px',
                        zIndex: 9999,
                      }}
                    >
                      <Select.Viewport style={{ padding: '4px' }}>
                        <Select.Item
                          value="admin"
                          style={{
                            display: 'flex',
                            alignItems: 'center',
                            padding: '8px 12px',
                            fontSize: '14px',
                            lineHeight: '20px',
                            borderRadius: 'var(--radius-2)',
                            cursor: 'pointer',
                            position: 'relative',
                            paddingLeft: '32px',
                          }}
                        >
                          <Select.ItemIndicator
                            style={{
                              position: 'absolute',
                              left: '8px',
                              display: 'inline-flex',
                              alignItems: 'center',
                            }}
                          >
                            <CheckIcon />
                          </Select.ItemIndicator>
                          <Select.ItemText>Admin</Select.ItemText>
                        </Select.Item>
                        <Select.Item
                          value="readonly"
                          style={{
                            display: 'flex',
                            alignItems: 'center',
                            padding: '8px 12px',
                            fontSize: '14px',
                            lineHeight: '20px',
                            borderRadius: 'var(--radius-2)',
                            cursor: 'pointer',
                            position: 'relative',
                            paddingLeft: '32px',
                          }}
                        >
                          <Select.ItemIndicator
                            style={{
                              position: 'absolute',
                              left: '8px',
                              display: 'inline-flex',
                              alignItems: 'center',
                            }}
                          >
                            <CheckIcon />
                          </Select.ItemIndicator>
                          <Select.ItemText>Read Only</Select.ItemText>
                        </Select.Item>
                      </Select.Viewport>
                    </Select.Content>
                  </Select.Portal>
                </Select.Root>
              </Box>

              {/* Actions Column */}
              <Box style={{ width: '60px' }}>
                <IconButton
                  type="button"
                  size="2"
                  variant="ghost"
                  color="red"
                  onClick={() => handleDeleteUser(index)}
                >
                  <TrashIcon />
                </IconButton>
              </Box>
            </Flex>
          ))}
        </Box>
      ) : (
        <Box
          p="8"
          style={{
            border: '1px solid var(--gray-6)',
            borderRadius: 'var(--radius-3)',
            textAlign: 'center',
          }}
        >
          <Text size="2" color="gray">
            No users added yet. Click &ldquo;Add User&rdquo; to get started.
          </Text>
        </Box>
      )}

      {errors?.project_users && (
        <Text size="2" color="red">
          {errors.project_users.message as string}
        </Text>
      )}
    </Flex>
  )
}

