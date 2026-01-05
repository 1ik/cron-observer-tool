'use client'

import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as Dialog from '@radix-ui/react-dialog'
import * as Label from '@radix-ui/react-label'
import { Box, Flex, Text, Heading, Button, TextField, TextArea } from '@radix-ui/themes'
import { CreateProjectRequest } from '../lib/types/project'
import { createProjectSchema, CreateProjectFormData } from '../lib/validations/project'
import { StyledDialogContent } from './StyledDialogContent'

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
  } =   useForm<CreateProjectFormData>({
    resolver: zodResolver(createProjectSchema),
    defaultValues: {
      name: '',
      description: '',
      execution_endpoint: '',
    },
  })

  const onFormSubmit = (data: CreateProjectFormData) => {
    // Transform form data to API request format (convert empty strings to undefined)
    const requestData: CreateProjectRequest = {
      name: data.name.trim(),
      description: data.description?.trim() || undefined,
      execution_endpoint: data.execution_endpoint?.trim() || undefined,
    }
    
    // Debug logging
    console.log('Form data:', data)
    console.log('Request data:', requestData)
    
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
      <StyledDialogContent maxWidth="500px">
        {/* Header - Sticky */}
        <Box
          p="5"
          style={{
            flexShrink: 0,
            borderBottom: '1px solid var(--gray-6)',
          }}
        >
          <Dialog.Title asChild>
            <Heading size="5" mb="2">
              Create New Project
            </Heading>
          </Dialog.Title>

          <Dialog.Description asChild>
            <Text size="3" color="gray">
              Create a new project to organize your tasks and task groups.
            </Text>
          </Dialog.Description>
        </Box>

        {/* Content - Scrollable */}
        <Box
          p="5"
          style={{
            flex: 1,
            overflowY: 'auto',
            minHeight: 0,
          }}
        >
          <form id="create-project-form" onSubmit={handleSubmit(onFormSubmit)}>
            <Flex direction="column" gap="4">
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

              <Flex direction="column" gap="2">
                <Label.Root htmlFor="project-execution-endpoint">
                  <Text size="3" weight="medium">
                    Execution Endpoint
                  </Text>
                </Label.Root>
                <TextField.Root
                  id="project-execution-endpoint"
                  {...register('execution_endpoint')}
                  placeholder="https://api.example.com/execute"
                  size="3"
                  type="url"
                  color={errors.execution_endpoint ? 'red' : undefined}
                />
                <Text size="1" color="gray">
                  URL where task executions will be sent via POST request
                </Text>
                {errors.execution_endpoint && (
                  <Text size="2" color="red">
                    {errors.execution_endpoint.message}
                  </Text>
                )}
              </Flex>
            </Flex>
          </form>
        </Box>

        {/* Footer - Sticky */}
        <Box
          p="5"
          style={{
            flexShrink: 0,
            borderTop: '1px solid var(--gray-6)',
          }}
        >
          <Flex gap="3" justify="end">
            <Dialog.Close asChild>
              <Button type="button" variant="soft" onClick={handleCancel}>
                Cancel
              </Button>
            </Dialog.Close>
            <Button
              type="button"
              variant="solid"
              onClick={handleSubmit(onFormSubmit)}
            >
              Create Project
            </Button>
          </Flex>
        </Box>
      </StyledDialogContent>
    </Dialog.Root>
  )
}
