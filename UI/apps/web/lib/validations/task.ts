import { z } from 'zod'

export const taskStatusSchema = z.enum(['ACTIVE', 'PAUSED', 'DISABLED'])
export const scheduleTypeSchema = z.enum(['RECURRING', 'ONEOFF'])
export const frequencyUnitSchema = z.enum(['s', 'm', 'h'])

export const frequencySchema = z.object({
  value: z.number().int().min(1, 'Frequency value must be at least 1'),
  unit: frequencyUnitSchema,
})

export const timeRangeSchema = z.object({
  start: z.string().regex(/^([0-1][0-9]|2[0-3]):[0-5][0-9]$/, 'Start time must be in HH:MM format (24-hour format)'),
  end: z.string().regex(/^([0-1][0-9]|2[0-3]):[0-5][0-9]$/, 'End time must be in HH:MM format (24-hour format)'),
  frequency: frequencySchema,
})

export const scheduleConfigSchema = z.object({
  cron_expression: z.string().trim().optional().or(z.literal('')),
  timezone: z.string().min(1, 'Timezone is required').trim(),
  time_range: timeRangeSchema.optional(),
  days_of_week: z.array(z.number().int().min(0).max(6)).optional(),
  exclusions: z.array(z.number().int().min(0).max(6)).optional(),
})

export const httpTriggerConfigSchema = z.object({
  url: z.string().url('URL must be a valid URL').trim(),
  method: z.string().min(1, 'HTTP method is required').trim(),
  headers: z.record(z.string(), z.string()).optional(),
  body: z.unknown().optional(),
  timeout: z.number().int().min(1).max(300).optional(),
})

export const triggerConfigSchema = z.object({
  type: z.literal('HTTP'),
  http: httpTriggerConfigSchema,
})

export const updateTaskSchema = z.object({
  name: z
    .string()
    .min(1, 'Task name is required')
    .max(255, 'Task name must be 255 characters or less')
    .trim(),
  description: z
    .string()
    .max(1000, 'Description must be 1000 characters or less')
    .trim()
    .optional()
    .or(z.literal('')),
  schedule_type: scheduleTypeSchema,
  status: taskStatusSchema.optional(),
  schedule_config: scheduleConfigSchema,
  trigger_config: triggerConfigSchema,
  task_group_id: z.string().trim().optional().or(z.literal('')),
  metadata: z.record(z.string(), z.unknown()).optional(),
})

export type UpdateTaskFormData = z.infer<typeof updateTaskSchema>

