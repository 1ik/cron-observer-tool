'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import * as Dialog from '@radix-ui/react-dialog'
import * as Label from '@radix-ui/react-label'
import { Box, Button, Flex, Heading, Select, Text, TextArea, TextField } from '@radix-ui/themes'
import { useEffect } from 'react'
import { Controller, useForm, useWatch } from 'react-hook-form'
import { TIMEZONES } from '../lib/constants/timezones'
import { useCronDescription } from '../lib/hooks/use-cron-description'
import { Task, TaskStatus, UpdateTaskRequest } from '../lib/types/task'
import { UpdateTaskFormData, updateTaskSchema } from '../lib/validations/task'
import { StyledDialogContent } from './StyledDialogContent'

interface TaskSettingsDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  task: Task
  onSubmit: (data: UpdateTaskRequest) => void
  isReadOnly?: boolean
}

export function TaskSettingsDialog({
  open,
  onOpenChange,
  task,
  onSubmit,
  isReadOnly = false,
}: TaskSettingsDialogProps) {
  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateTaskFormData>({
    resolver: zodResolver(updateTaskSchema),
    mode: 'onChange', // Enable onChange mode to watch cron expression changes
    defaultValues: {
      name: task.name,
      description: task.description || '',
      status: task.status,
      schedule_config: {
        ...task.schedule_config,
        timezone: task.schedule_config?.timezone || 'Asia/Dhaka',
      },
      task_group_id: task.task_group_id || '',
      timeout_seconds: task.timeout_seconds !== undefined && task.timeout_seconds !== null ? task.timeout_seconds : undefined,
    },
  })

  // Watch cron expression to show description
  const cronExpression = useWatch({
    control,
    name: 'schedule_config.cron_expression',
  })
  const cronDescription = useCronDescription(cronExpression)

  // Reset form when task changes
  useEffect(() => {
    if (open) {
      // Debug: log the task to see what timeout_seconds value we're getting
      console.log('TaskSettingsDialog - task.timeout_seconds:', task.timeout_seconds, typeof task.timeout_seconds)
      
      reset({
        name: task.name,
        description: task.description || '',
        status: task.status,
        schedule_config: {
          ...task.schedule_config,
          timezone: task.schedule_config?.timezone || 'Asia/Dhaka',
        },
        task_group_id: task.task_group_id || '',
        timeout_seconds: task.timeout_seconds !== undefined && task.timeout_seconds !== null ? task.timeout_seconds : undefined,
      })
    }
  }, [task, open, reset])

  const onFormSubmit = (data: UpdateTaskFormData) => {
    // Transform form data to API request format (convert empty strings to undefined)
    const requestData: UpdateTaskRequest = {
      name: data.name,
      description: data.description || undefined,
      schedule_type: 'RECURRING', // Default value since schedule_type is no longer shown in UI
      status: data.status,
      schedule_config: data.schedule_config,
      task_group_id: data.task_group_id || undefined,
      timeout_seconds: data.timeout_seconds || undefined,
      metadata: data.metadata,
    }
    onSubmit(requestData)
  }

  const handleCancel = () => {
    reset()
    onOpenChange(false)
  }

  const getStatusDotColor = (s: TaskStatus) => {
    switch (s) {
      case 'ACTIVE':
        return 'var(--green-9)'
      case 'DISABLED':
        return 'var(--gray-9)'
      default:
        return 'var(--gray-9)'
    }
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
              Edit Task Settings
            </Heading>
          </Dialog.Title>

          <Dialog.Description asChild>
            <Text size="3" color="gray">
              Update the task configuration and settings.
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
                readOnly
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
                disabled={isReadOnly}
              />
              {errors.description && (
                <Text size="2" color="red">
                  {errors.description.message}
                </Text>
              )}
            </Flex>

            {/* Status */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-status">
                <Text size="3" weight="medium">
                  Status
                </Text>
              </Label.Root>
              <Controller
                name="status"
                control={control}
                render={({ field }) => (
                  <Select.Root
                    value={field.value}
                    onValueChange={field.onChange}
                    disabled={isReadOnly}
                  >
                    <Select.Trigger id="task-status" style={{ width: '100%' }}>
                      <Flex align="center" gap="2">
                        <Box
                          style={{
                            width: '8px',
                            height: '8px',
                            borderRadius: '50%',
                            backgroundColor: getStatusDotColor(field.value as TaskStatus),
                            flexShrink: 0,
                          }}
                        />
                        <Text>{field.value}</Text>
                      </Flex>
                    </Select.Trigger>
                    <Select.Content>
                      <Select.Item value="ACTIVE">
                        <Flex align="center" gap="2">
                          <Box
                            style={{
                              width: '8px',
                              height: '8px',
                              borderRadius: '50%',
                              backgroundColor: getStatusDotColor('ACTIVE'),
                            }}
                          />
                          <Text>ACTIVE</Text>
                        </Flex>
                      </Select.Item>
                      <Select.Item value="DISABLED">
                        <Flex align="center" gap="2">
                          <Box
                            style={{
                              width: '8px',
                              height: '8px',
                              borderRadius: '50%',
                              backgroundColor: getStatusDotColor('DISABLED'),
                            }}
                          />
                          <Text>DISABLED</Text>
                        </Flex>
                      </Select.Item>
                    </Select.Content>
                  </Select.Root>
                )}
              />
              {errors.status && (
                <Text size="2" color="red">
                  {errors.status.message}
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
                      disabled={isReadOnly}
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
                disabled={isReadOnly}
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

              {/* Timeout */}
              <Flex direction="column" gap="2">
                <Label.Root htmlFor="task-timeout">
                  <Text size="3" weight="medium">
                    Execution Timeout (seconds)
                  </Text>
                </Label.Root>
                <Controller
                  name="timeout_seconds"
                  control={control}
                  render={({ field }) => (
                    <TextField.Root
                      id="task-timeout"
                      type="number"
                      value={field.value !== undefined && field.value !== null ? String(field.value) : ''}
                      onChange={(e) => {
                        const value = e.target.value
                        field.onChange(value === '' ? undefined : Number(value))
                      }}
                      onBlur={field.onBlur}
                      placeholder="Optional - leave empty for no timeout"
                      size="3"
                      min={1}
                      disabled={isReadOnly}
                      color={errors.timeout_seconds ? 'red' : undefined}
                    />
                  )}
                />
                <Text size="1" color="gray">
                  Maximum time allowed for execution. If exceeded, execution will be marked as failed.
                </Text>
                {errors.timeout_seconds && (
                  <Text size="2" color="red">
                    {errors.timeout_seconds.message}
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
                {isReadOnly ? 'Close' : 'Cancel'}
              </Button>
            </Dialog.Close>
            {!isReadOnly && (
              <Button
                type="submit"
                variant="solid"
                onClick={handleSubmit(onFormSubmit)}
              >
                Save Changes
              </Button>
            )}
          </Flex>
        </Box>
      </StyledDialogContent>
    </Dialog.Root>
  )
}

