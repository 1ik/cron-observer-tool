'use client'

import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as Dialog from '@radix-ui/react-dialog'
import * as Label from '@radix-ui/react-label'
import { Box, Flex, Text, Heading, Button, TextField, TextArea } from '@radix-ui/themes'
import { CreateProjectRequest } from '../lib/types/project'
import { createProjectSchema, CreateProjectFormData } from '../lib/validations/project'

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
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateProjectFormData>({
    resolver: zodResolver(createProjectSchema),
    defaultValues: {
      name: '',
      description: '',
    },
  })

  const onFormSubmit = (data: CreateProjectFormData) => {
    // Transform form data to API request format (convert empty strings to undefined)
    const requestData: CreateProjectRequest = {
      name: data.name,
      description: data.description || undefined,
    }
    onSubmit(requestData)
    reset()
    onOpenChange(false)
  }

  const handleCancel = () => {
    reset()
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

        <Flex direction="column" gap="4" asChild>
          <form onSubmit={handleSubmit(onFormSubmit)}>
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="project-name">
                <Text size="3" weight="medium">
                  Project Name <Text color="red">*</Text>
                </Text>
              </Label.Root>
              <TextField.Root
                id="project-name"
                {...register('name')}
                placeholder="Enter project name"
                size="3"
                color={errors.name ? 'red' : undefined}
              />
              {errors.name && (
                <Text size="2" color="red">
                  {errors.name.message}
                </Text>
              )}
            </Flex>

            <Flex direction="column" gap="2">
              <Label.Root htmlFor="project-description">
                <Text size="3" weight="medium">
                  Description
                </Text>
              </Label.Root>
              <TextArea
                id="project-description"
                {...register('description')}
                placeholder="Enter project description (optional)"
                rows={4}
                size="3"
                color={errors.description ? 'red' : undefined}
              />
              {errors.description && (
                <Text size="2" color="red">
                  {errors.description.message}
                </Text>
              )}
            </Flex>

            <Flex gap="3" justify="end" mt="4">
              <Dialog.Close asChild>
                <Button type="button" variant="soft" onClick={handleCancel}>
                  Cancel
                </Button>
              </Dialog.Close>
              <Button type="submit" variant="solid">
                Create Project
              </Button>
            </Flex>
          </form>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  )
}
