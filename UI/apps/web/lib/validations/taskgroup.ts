import { z } from 'zod'

export const taskGroupStatusSchema = z.enum(['ACTIVE', 'PAUSED', 'DISABLED'])

export const updateTaskGroupSchema = z.object({
  name: z
    .string()
    .min(1, 'Task group name is required')
    .max(255, 'Task group name must be 255 characters or less')
    .trim(),
  description: z
    .string()
    .max(1000, 'Description must be 1000 characters or less')
    .trim()
    .optional()
    .or(z.literal('')),
  status: taskGroupStatusSchema,
  start_time: z
    .string()
    .regex(/^([0-1][0-9]|2[0-3]):[0-5][0-9]$/, 'Start time must be in HH:MM format (24-hour format)')
    .optional()
    .or(z.literal('')),
  end_time: z
    .string()
    .regex(/^([0-1][0-9]|2[0-3]):[0-5][0-9]$/, 'End time must be in HH:MM format (24-hour format)')
    .optional()
    .or(z.literal('')),
  timezone: z.string().trim().optional().or(z.literal('')),
})

export type UpdateTaskGroupFormData = z.infer<typeof updateTaskGroupSchema>

export const createTaskGroupSchema = z.object({
  name: z
    .string()
    .min(1, 'Task group name is required')
    .max(255, 'Task group name must be 255 characters or less')
    .trim(),
  description: z
    .string()
    .max(1000, 'Description must be 1000 characters or less')
    .trim()
    .optional()
    .or(z.literal('')),
  start_time: z
    .string()
    .regex(/^([0-1][0-9]|2[0-3]):[0-5][0-9]$/, 'Start time must be in HH:MM format (24-hour format)')
    .optional()
    .or(z.literal('')),
  end_time: z
    .string()
    .regex(/^([0-1][0-9]|2[0-3]):[0-5][0-9]$/, 'End time must be in HH:MM format (24-hour format)')
    .optional()
    .or(z.literal('')),
  timezone: z.string().trim().optional().or(z.literal('')),
})

export type CreateTaskGroupFormData = z.infer<typeof createTaskGroupSchema>

