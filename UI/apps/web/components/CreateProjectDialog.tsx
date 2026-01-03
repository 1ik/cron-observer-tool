'use client'

import { useState } from 'react'
import * as Dialog from '@radix-ui/react-dialog'
import * as Label from '@radix-ui/react-label'
import { Box, Flex, Text, Heading, Button, TextField, TextArea } from '@radix-ui/themes'
import { CreateProjectRequest } from '../lib/types/project'

interface CreateProjectDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (data: CreateProjectRequest) => void
}

export function CreateProjectDialog({
  open,
  onOpenChange,
  onSubmit,
}: CreateProjectDialogProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = (e?: React.FormEvent | React.MouseEvent) => {
    if (e) {
      e.preventDefault()
    }
    setError(null)

    if (!name.trim()) {
      setError('Project name is required')
      return
    }

    onSubmit({
      name: name.trim(),
      description: description.trim() || undefined,
    })

    // Reset form
    setName('')
    setDescription('')
    onOpenChange(false)
  }

  const handleCancel = () => {
    setName('')
    setDescription('')
    setError(null)
    onOpenChange(false)
  }

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Content style={{ maxWidth: '500px' }}>
        <Dialog.Title asChild>
          <Heading size="5" mb="4">
            Create New Project
          </Heading>
        </Dialog.Title>

        <Dialog.Description asChild>
          <Text size="3" color="gray" mb="4">
            Create a new project to organize your tasks and task groups.
          </Text>
        </Dialog.Description>

        <Flex direction="column" gap="4">
          <Flex direction="column" gap="2">
            <Label.Root htmlFor="project-name">
              <Text size="3" weight="medium">
                Project Name <Text color="red">*</Text>
              </Text>
            </Label.Root>
            <TextField.Root
              id="project-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter project name"
              size="3"
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  e.preventDefault()
                  handleSubmit(e as any)
                }
              }}
            />
          </Flex>

          <Flex direction="column" gap="2">
            <Label.Root htmlFor="project-description">
              <Text size="3" weight="medium">
                Description
              </Text>
            </Label.Root>
            <TextArea
              id="project-description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Enter project description (optional)"
              rows={4}
              size="3"
            />
          </Flex>

          {error && (
            <Box>
              <Text size="2" color="red">
                {error}
              </Text>
            </Box>
          )}

          <Flex gap="3" justify="end" mt="4">
            <Dialog.Close asChild>
              <Button type="button" variant="soft" onClick={handleCancel}>
                Cancel
              </Button>
            </Dialog.Close>
            <Button type="button" variant="solid" onClick={handleSubmit}>
              Create Project
            </Button>
          </Flex>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  )
}

