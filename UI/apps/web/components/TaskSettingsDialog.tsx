'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import * as Dialog from '@radix-ui/react-dialog'
import { ChevronDownIcon } from '@radix-ui/react-icons'
import * as Label from '@radix-ui/react-label'
import * as Select from '@radix-ui/react-select'
import { Box, Button, Flex, Heading, Text, TextArea, TextField } from '@radix-ui/themes'
import { useEffect } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { Task, TaskStatus, UpdateTaskRequest } from '../lib/types/task'
import { UpdateTaskFormData, updateTaskSchema } from '../lib/validations/task'
import { StyledDialogContent } from './StyledDialogContent'

interface TaskSettingsDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  task: Task
  onSubmit: (data: UpdateTaskRequest) => void
}

export function TaskSettingsDialog({
  open,
  onOpenChange,
  task,
  onSubmit,
}: TaskSettingsDialogProps) {
  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateTaskFormData>({
    resolver: zodResolver(updateTaskSchema),
    defaultValues: {
      name: task.name,
      description: task.description || '',
      schedule_type: task.schedule_type,
      status: task.status,
      schedule_config: task.schedule_config,
      trigger_config: task.trigger_config,
      task_group_id: task.task_group_id || '',
    },
  })

  // Reset form when task changes
  useEffect(() => {
    if (open) {
      reset({
        name: task.name,
        description: task.description || '',
        schedule_type: task.schedule_type,
        status: task.status,
        schedule_config: task.schedule_config,
        trigger_config: task.trigger_config,
        task_group_id: task.task_group_id || '',
      })
    }
  }, [task, open, reset])

  const onFormSubmit = (data: UpdateTaskFormData) => {
    // Transform form data to API request format (convert empty strings to undefined)
    const requestData: UpdateTaskRequest = {
      name: data.name,
      description: data.description || undefined,
      schedule_type: data.schedule_type,
      status: data.status,
      schedule_config: data.schedule_config,
      trigger_config: data.trigger_config,
      task_group_id: data.task_group_id || undefined,
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
      case 'PAUSED':
        return 'var(--yellow-9)'
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
                        <Select.Value />
                      </Flex>
                      <Select.Icon>
                        <ChevronDownIcon />
                      </Select.Icon>
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
                      <Select.Item value="PAUSED">
                        <Flex align="center" gap="2">
                          <Box
                            style={{
                              width: '8px',
                              height: '8px',
                              borderRadius: '50%',
                              backgroundColor: getStatusDotColor('PAUSED'),
                            }}
                          />
                          <Text>PAUSED</Text>
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

            {/* Schedule Type */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-schedule-type">
                <Text size="3" weight="medium">
                  Schedule Type <Text color="red">*</Text>
                </Text>
              </Label.Root>
              <Controller
                name="schedule_type"
                control={control}
                render={({ field }) => (
                  <Select.Root
                    value={field.value}
                    onValueChange={field.onChange}
                  >
                    <Select.Trigger id="task-schedule-type" style={{ width: '100%' }}>
                      <Select.Value />
                      <Select.Icon>
                        <ChevronDownIcon />
                      </Select.Icon>
                    </Select.Trigger>
                    <Select.Content>
                      <Select.Item value="RECURRING">
                        <Text>RECURRING</Text>
                      </Select.Item>
                      <Select.Item value="ONEOFF">
                        <Text>ONEOFF</Text>
                      </Select.Item>
                    </Select.Content>
                  </Select.Root>
                )}
              />
              {errors.schedule_type && (
                <Text size="2" color="red">
                  {errors.schedule_type.message}
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
              <TextField.Root
                id="task-timezone"
                {...register('schedule_config.timezone')}
                placeholder="e.g., America/New_York, UTC"
                size="3"
                color={errors.schedule_config?.timezone ? 'red' : undefined}
              />
              <Text size="1" color="gray">
                IANA timezone identifier (e.g., America/New_York, UTC)
              </Text>
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
              <Text size="1" color="gray">
                If provided, TimeRange and DaysOfWeek are ignored
              </Text>
              {errors.schedule_config?.cron_expression && (
                <Text size="2" color="red">
                  {errors.schedule_config.cron_expression.message}
                </Text>
              )}
            </Flex>

            {/* Trigger Config - HTTP URL */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-trigger-url">
                <Text size="3" weight="medium">
                  HTTP URL <Text color="red">*</Text>
                </Text>
              </Label.Root>
              <TextField.Root
                id="task-trigger-url"
                {...register('trigger_config.http.url')}
                placeholder="https://api.example.com/endpoint"
                size="3"
                color={errors.trigger_config?.http?.url ? 'red' : undefined}
              />
              {errors.trigger_config?.http?.url && (
                <Text size="2" color="red">
                  {errors.trigger_config.http.url.message}
                </Text>
              )}
            </Flex>

            {/* Trigger Config - HTTP Method */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-trigger-method">
                <Text size="3" weight="medium">
                  HTTP Method <Text color="red">*</Text>
                </Text>
              </Label.Root>
              <Controller
                name="trigger_config.http.method"
                control={control}
                render={({ field }) => (
                  <Select.Root
                    value={field.value}
                    onValueChange={field.onChange}
                  >
                    <Select.Trigger id="task-trigger-method" style={{ width: '100%' }}>
                      <Select.Value />
                      <Select.Icon>
                        <ChevronDownIcon />
                      </Select.Icon>
                    </Select.Trigger>
                    <Select.Content>
                      <Select.Item value="GET">
                        <Text>GET</Text>
                      </Select.Item>
                      <Select.Item value="POST">
                        <Text>POST</Text>
                      </Select.Item>
                      <Select.Item value="PUT">
                        <Text>PUT</Text>
                      </Select.Item>
                      <Select.Item value="PATCH">
                        <Text>PATCH</Text>
                      </Select.Item>
                      <Select.Item value="DELETE">
                        <Text>DELETE</Text>
                      </Select.Item>
                    </Select.Content>
                  </Select.Root>
                )}
              />
              {errors.trigger_config?.http?.method && (
                <Text size="2" color="red">
                  {errors.trigger_config.http.method.message}
                </Text>
              )}
            </Flex>

            {/* Trigger Config - HTTP Timeout */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-trigger-timeout">
                <Text size="3" weight="medium">
                  HTTP Timeout (seconds)
                </Text>
              </Label.Root>
              <TextField.Root
                id="task-trigger-timeout"
                type="number"
                {...register('trigger_config.http.timeout', { valueAsNumber: true })}
                placeholder="300"
                size="3"
                min={1}
                max={300}
                color={errors.trigger_config?.http?.timeout ? 'red' : undefined}
              />
              <Text size="1" color="gray">
                Timeout in seconds (1-300)
              </Text>
              {errors.trigger_config?.http?.timeout && (
                <Text size="2" color="red">
                  {errors.trigger_config.http.timeout.message}
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
            <Button
              type="submit"
              variant="solid"
              onClick={handleSubmit(onFormSubmit)}
            >
              Save Changes
            </Button>
          </Flex>
        </Box>
      </StyledDialogContent>
    </Dialog.Root>
  )
}

