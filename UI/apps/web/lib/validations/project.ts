import { z } from 'zod'

export const createProjectSchema = z.object({
  name: z
    .string()
    .min(1, 'Project name is required')
    .max(255, 'Project name must be 255 characters or less')
    .trim(),
  description: z
    .string()
    .max(1000, 'Description must be 1000 characters or less')
    .trim()
    .optional()
    .or(z.literal('')),
})

export type CreateProjectFormData = z.infer<typeof createProjectSchema>

