'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import * as Dialog from '@radix-ui/react-dialog'
import * as Label from '@radix-ui/react-label'
import { Box, Button, Flex, Heading, Select, Text, TextArea, TextField } from '@radix-ui/themes'
import { useEffect } from 'react'
import { Controller, useForm, useWatch } from 'react-hook-form'
import { TIMEZONES } from '../lib/constants/timezones'
import { useCronDescription } from '../lib/hooks/use-cron-description'
import { CreateTaskRequest } from '../lib/types/task'
import { CreateTaskFormData, createTaskSchema } from '../lib/validations/task'
import { StyledDialogContent } from './StyledDialogContent'

interface CreateTaskDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  projectId: string
  taskGroupId?: string
  onSubmit: (data: CreateTaskRequest) => void
}

export function CreateTaskDialog({
  open,
  onOpenChange,
  projectId,
  taskGroupId,
  onSubmit,
}: CreateTaskDialogProps) {
  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateTaskFormData>({
    mode: 'onChange', // Enable onChange mode to watch cron expression changes
    resolver: zodResolver(createTaskSchema),
    defaultValues: {
      name: '',
      description: '',
      schedule_config: {
        timezone: 'Asia/Dhaka',
        cron_expression: '',
      },
    },
  })

  // Watch cron expression to show description
  const cronExpression = useWatch({
    control,
    name: 'schedule_config.cron_expression',
  })
  const cronDescription = useCronDescription(cronExpression)

  // Reset form when dialog opens
  useEffect(() => {
    if (open) {
      reset({
        name: '',
        description: '',
        schedule_config: {
          timezone: 'Asia/Dhaka',
          cron_expression: '',
        },
      })
    }
  }, [open, reset])

  const onFormSubmit = (data: CreateTaskFormData) => {
    // Transform form data to API request format
    const requestData: CreateTaskRequest = {
      project_id: projectId,
      task_group_id: taskGroupId,
      name: data.name,
      description: data.description || undefined,
      schedule_type: 'RECURRING', // Default value since schedule_type is no longer shown in UI
      schedule_config: {
        timezone: data.schedule_config.timezone,
        cron_expression: data.schedule_config.cron_expression || undefined,
        time_range: data.schedule_config.time_range,
        days_of_week: data.schedule_config.days_of_week,
        exclusions: data.schedule_config.exclusions,
      },
      timeout_seconds: data.timeout_seconds || undefined,
      metadata: data.metadata,
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
              Create New Task
            </Heading>
          </Dialog.Title>

          <Dialog.Description asChild>
            <Text size="3" color="gray">
              Create a new scheduled task to automate your workflows.
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
          <Flex direction="column" gap="4" asChild>
            <form onSubmit={handleSubmit(onFormSubmit)}>
              {/* Name */}
              <Flex direction="column" gap="2">
                <Label.Root htmlFor="task-name">
                  <Text size="3" weight="medium">
                    Name <Text color="red">*</Text>
                  </Text>
                </Label.Root>
                <TextField.Root
                  id="task-name"
                  {...register('name')}
                  placeholder="Enter task name"
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
                <Label.Root htmlFor="task-description">
                  <Text size="3" weight="medium">
                    Description
                  </Text>
                </Label.Root>
                <TextArea
                  id="task-description"
                  {...register('description')}
                  placeholder="Enter task description (optional)"
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

              {/* Schedule Config - Timezone */}
              <Flex direction="column" gap="2">
                <Label.Root htmlFor="task-timezone">
                  <Text size="3" weight="medium">
                    Timezone <Text color="red">*</Text>
                  </Text>
                </Label.Root>
                <Controller
                  name="schedule_config.timezone"
                  control={control}
                  render={({ field }) => {
                    const selectedTimezone = TIMEZONES.find((tz) => tz.value === field.value)
                    return (
                      <Select.Root
                        value={field.value}
                        onValueChange={field.onChange}
                      >
                        <Select.Trigger
                          id="task-timezone"
                          style={{ width: '100%' }}
                          color={errors.schedule_config?.timezone ? 'red' : undefined}
                        >
                          <Text>{selectedTimezone ? selectedTimezone.label : 'Select timezone'}</Text>
                        </Select.Trigger>
                        <Select.Content>
                          {TIMEZONES.map((tz) => (
                            <Select.Item key={tz.value} value={tz.value}>
                              {tz.label}
                            </Select.Item>
                          ))}
                        </Select.Content>
                      </Select.Root>
                    )
                  }}
                />
                {errors.schedule_config?.timezone && (
                  <Text size="2" color="red">
                    {errors.schedule_config.timezone.message}
                  </Text>
                )}
              </Flex>

              {/* Schedule Config - Cron Expression */}
              <Flex direction="column" gap="2">
                <Label.Root htmlFor="task-cron-expression">
                  <Text size="3" weight="medium">
                    Cron Expression
                  </Text>
                </Label.Root>
                <TextField.Root
                  id="task-cron-expression"
                  {...register('schedule_config.cron_expression')}
                  placeholder="e.g., 0 2 * * *"
                  size="3"
                  color={errors.schedule_config?.cron_expression ? 'red' : undefined}
                />
                {cronDescription && (
                  <Text size="2" color="gray">
                    {cronDescription}
                  </Text>
                )}
                {errors.schedule_config?.cron_expression && (
                  <Text size="2" color="red">
                    {errors.schedule_config.cron_expression.message}
                  </Text>
                )}
              </Flex>

            </form>
          </Flex>
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
              Create Task
            </Button>
          </Flex>
        </Box>
      </StyledDialogContent>
    </Dialog.Root>
  )
}

