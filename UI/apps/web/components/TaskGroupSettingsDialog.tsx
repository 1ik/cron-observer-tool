'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import * as Dialog from '@radix-ui/react-dialog'
import { ChevronDownIcon } from '@radix-ui/react-icons'
import * as Label from '@radix-ui/react-label'
import * as Select from '@radix-ui/react-select'
import { Box, Button, Flex, Heading, Text, TextArea, TextField } from '@radix-ui/themes'
import { useEffect } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { TaskGroup, TaskGroupStatus, UpdateTaskGroupRequest } from '../lib/types/taskgroup'
import { UpdateTaskGroupFormData, updateTaskGroupSchema } from '../lib/validations/taskgroup'
import { StyledDialogContent } from './StyledDialogContent'

interface TaskGroupSettingsDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  taskGroup: TaskGroup
  onSubmit: (data: UpdateTaskGroupRequest) => void
}

export function TaskGroupSettingsDialog({
  open,
  onOpenChange,
  taskGroup,
  onSubmit,
}: TaskGroupSettingsDialogProps) {
  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateTaskGroupFormData>({
    resolver: zodResolver(updateTaskGroupSchema),
    defaultValues: {
      name: taskGroup.name,
      description: taskGroup.description || '',
      status: taskGroup.status,
      start_time: taskGroup.start_time || '',
      end_time: taskGroup.end_time || '',
      timezone: taskGroup.timezone || '',
    },
  })

  // Reset form when taskGroup changes
  useEffect(() => {
    if (open) {
      reset({
        name: taskGroup.name,
        description: taskGroup.description || '',
        status: taskGroup.status,
        start_time: taskGroup.start_time || '',
        end_time: taskGroup.end_time || '',
        timezone: taskGroup.timezone || '',
      })
    }
  }, [taskGroup, open, reset])

  const onFormSubmit = (data: UpdateTaskGroupFormData) => {
    // Transform form data to API request format (convert empty strings to undefined)
    const requestData: UpdateTaskGroupRequest = {
      name: data.name,
      description: data.description || undefined,
      status: data.status,
      start_time: data.start_time || undefined,
      end_time: data.end_time || undefined,
      timezone: data.timezone || undefined,
    }
    onSubmit(requestData)
  }

  const handleCancel = () => {
    reset()
    onOpenChange(false)
  }

  const getStatusDotColor = (s: TaskGroupStatus) => {
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
        <Dialog.Title asChild>
          <Heading size="5" mb="4">
            Edit Task Group Settings
          </Heading>
        </Dialog.Title>

        <Dialog.Description asChild>
          <Text size="3" color="gray" mb="4">
            Update the task group configuration and settings.
          </Text>
        </Dialog.Description>

        <Flex direction="column" gap="4" asChild>
          <form onSubmit={handleSubmit(onFormSubmit)}>
            {/* Name */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-group-name">
                <Text size="3" weight="medium">
                  Name <Text color="red">*</Text>
                </Text>
              </Label.Root>
              <TextField.Root
                id="task-group-name"
                {...register('name')}
                placeholder="Enter task group name"
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
              <Label.Root htmlFor="task-group-description">
                <Text size="3" weight="medium">
                  Description
                </Text>
              </Label.Root>
              <TextArea
                id="task-group-description"
                {...register('description')}
                placeholder="Enter task group description (optional)"
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
              <Label.Root htmlFor="task-group-status">
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
                <Select.Trigger id="task-group-status" style={{ width: '100%' }}>
                  <Flex align="center" gap="2">
                    <Box
                      style={{
                        width: '8px',
                        height: '8px',
                        borderRadius: '50%',
                        backgroundColor: getStatusDotColor(field.value),
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

            {/* Start Time */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-group-start-time">
                <Text size="3" weight="medium">
                  Start Time
                </Text>
              </Label.Root>
              <TextField.Root
                id="task-group-start-time"
                {...register('start_time')}
                placeholder="HH:MM"
                size="3"
                pattern="^([01]\d|2[0-3]):([0-5]\d)$"
                color={errors.start_time ? 'red' : undefined}
              />
              <Text size="1" color="gray">
                Format: HH:MM (24-hour format)
              </Text>
              {errors.start_time && (
                <Text size="2" color="red">
                  {errors.start_time.message}
                </Text>
              )}
            </Flex>

            {/* End Time */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-group-end-time">
                <Text size="3" weight="medium">
                  End Time
                </Text>
              </Label.Root>
              <TextField.Root
                id="task-group-end-time"
                {...register('end_time')}
                placeholder="HH:MM"
                size="3"
                pattern="^([01]\d|2[0-3]):([0-5]\d)$"
                color={errors.end_time ? 'red' : undefined}
              />
              <Text size="1" color="gray">
                Format: HH:MM (24-hour format)
              </Text>
              {errors.end_time && (
                <Text size="2" color="red">
                  {errors.end_time.message}
                </Text>
              )}
            </Flex>

            {/* Timezone */}
            <Flex direction="column" gap="2">
              <Label.Root htmlFor="task-group-timezone">
                <Text size="3" weight="medium">
                  Timezone
                </Text>
              </Label.Root>
              <TextField.Root
                id="task-group-timezone"
                {...register('timezone')}
                placeholder="e.g., America/New_York"
                size="3"
                color={errors.timezone ? 'red' : undefined}
              />
              <Text size="1" color="gray">
                IANA timezone identifier (e.g., America/New_York, UTC)
              </Text>
              {errors.timezone && (
                <Text size="2" color="red">
                  {errors.timezone.message}
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
                Save Changes
              </Button>
            </Flex>
          </form>
        </Flex>
      </StyledDialogContent>
    </Dialog.Root>
  )
}
