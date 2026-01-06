'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import * as Dialog from '@radix-ui/react-dialog'
import * as Label from '@radix-ui/react-label'
import { Box, Button, Flex, Heading, Tabs, Text, TextArea, TextField } from '@radix-ui/themes'
import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { Project, UpdateProjectRequest } from '../lib/types/project'
import { UpdateProjectFormData, updateProjectSchema } from '../lib/validations/project'
import { StyledDialogContent } from './StyledDialogContent'

interface ProjectSettingsDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  project: Project
  onSubmit: (data: UpdateProjectRequest) => void
}

export function ProjectSettingsDialog({
  open,
  onOpenChange,
  project,
  onSubmit,
}: ProjectSettingsDialogProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateProjectFormData>({
    resolver: zodResolver(updateProjectSchema),
    defaultValues: {
      name: project.name,
      description: project.description || '',
      execution_endpoint: project.execution_endpoint || '',
      alert_emails: project.alert_emails || '',
    },
  })

  // Reset form when dialog opens or project changes
  useEffect(() => {
    if (open) {
      reset({
        name: project.name,
        description: project.description || '',
        execution_endpoint: project.execution_endpoint || '',
        alert_emails: project.alert_emails || '',
      })
    }
  }, [open, project, reset])

  const onFormSubmit = (data: UpdateProjectFormData) => {
    // Transform form data to API request format
    const requestData: UpdateProjectRequest = {
      name: data.name,
      description: data.description || undefined,
      execution_endpoint: data.execution_endpoint || undefined,
      alert_emails: data.alert_emails || undefined,
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
      <StyledDialogContent maxWidth="600px">
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
              Project Settings
            </Heading>
          </Dialog.Title>

          <Dialog.Description asChild>
            <Text size="3" color="gray">
              Configure project settings and alert preferences.
            </Text>
          </Dialog.Description>
        </Box>

        {/* Content - Scrollable with Tabs */}
        <Box
          p="5"
          style={{
            flex: 1,
            overflow: 'hidden',
            display: 'flex',
            flexDirection: 'column',
            minHeight: 0,
          }}
        >
          <form onSubmit={handleSubmit(onFormSubmit)} style={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: 0 }}>
            <Tabs.Root defaultValue="details" style={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: 0 }}>
              <Tabs.List mb="4">
                <Tabs.Trigger value="details">Details</Tabs.Trigger>
                <Tabs.Trigger value="alerts">Alerts</Tabs.Trigger>
              </Tabs.List>

              {/* Details Tab */}
              <Tabs.Content value="details" style={{ flex: 1, overflowY: 'auto' }}>
                <Flex direction="column" gap="4">
                  {/* Name */}
                  <Flex direction="column" gap="2">
                    <Label.Root htmlFor="project-name">
                      <Text size="3" weight="medium">
                        Name <Text color="red">*</Text>
                      </Text>
                    </Label.Root>
                    <TextField.Root
                      id="project-name"
                      {...register('name')}
                      placeholder="Enter project name"
                      size="3"
                      maxLength={255}
                      color={errors.name ? 'red' : undefined}
                    />
                    {errors.name && (
                      <Text size="2" color="red">
                        {errors.name.message}
                      </Text>
                    )}
                  </Flex>

                  {/* Description */}
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
                      maxLength={1000}
                      color={errors.description ? 'red' : undefined}
                    />
                    {errors.description && (
                      <Text size="2" color="red">
                        {errors.description.message}
                      </Text>
                    )}
                  </Flex>

                  {/* Execution Endpoint */}
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
                      color={errors.execution_endpoint ? 'red' : undefined}
                    />
                    <Text size="1" color="gray">
                      URL where task execution notifications will be sent
                    </Text>
                    {errors.execution_endpoint && (
                      <Text size="2" color="red">
                        {errors.execution_endpoint.message}
                      </Text>
                    )}
                  </Flex>
                </Flex>
              </Tabs.Content>

              {/* Alerts Tab */}
              <Tabs.Content value="alerts" style={{ flex: 1, overflowY: 'auto' }}>
                <Flex direction="column" gap="4">
                  {/* Alert Emails */}
                  <Flex direction="column" gap="2">
                    <Label.Root htmlFor="project-alert-emails">
                      <Text size="3" weight="medium">
                        Alert Emails
                      </Text>
                    </Label.Root>
                    <TextArea
                      id="project-alert-emails"
                      {...register('alert_emails')}
                      placeholder="email1@example.com, email2@example.com"
                      rows={6}
                      size="3"
                      color={errors.alert_emails ? 'red' : undefined}
                    />
                    <Text size="1" color="gray">
                      Comma-separated email addresses for receiving alerts (task failures, errors, etc.)
                    </Text>
                    {errors.alert_emails && (
                      <Text size="2" color="red">
                        {errors.alert_emails.message}
                      </Text>
                    )}
                  </Flex>
                </Flex>
              </Tabs.Content>
            </Tabs.Root>
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
            <Button type="submit" variant="solid" onClick={handleSubmit(onFormSubmit)}>
              Save Settings
            </Button>
          </Flex>
        </Box>
      </StyledDialogContent>
    </Dialog.Root>
  )
}
