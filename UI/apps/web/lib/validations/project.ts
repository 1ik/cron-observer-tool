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
  execution_endpoint: z
    .string()
    .url('Execution endpoint must be a valid URL')
    .trim()
    .optional()
    .or(z.literal('')),
})

export type CreateProjectFormData = z.infer<typeof createProjectSchema>

// Helper function to validate comma-separated emails
const validateCommaSeparatedEmails = (emails: string): boolean => {
  if (!emails.trim()) return true // Empty is valid (optional field)
  
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  const emailList = emails.split(',').map(email => email.trim()).filter(Boolean)
  
  return emailList.every(email => emailRegex.test(email))
}

const projectUserSchema = z.object({
  email: z.string().email('Invalid email address'),
  role: z.enum(['admin', 'readonly']),
})

export const updateProjectSchema = z.object({
  name: z
    .string()
    .min(1, 'Project name is required')
    .max(255, 'Project name must be 255 characters or less')
    .trim()
    .optional(),
  description: z
    .string()
    .max(1000, 'Description must be 1000 characters or less')
    .trim()
    .optional()
    .or(z.literal('')),
  execution_endpoint: z
    .string()
    .url('Execution endpoint must be a valid URL')
    .trim()
    .optional()
    .or(z.literal('')),
  alert_emails: z
    .string()
    .trim()
    .optional()
    .or(z.literal(''))
    .refine(
      (val) => !val || validateCommaSeparatedEmails(val),
      {
        message: 'Please enter valid email addresses separated by commas',
      }
    ),
  project_users: z.array(projectUserSchema).optional(),
})

export type UpdateProjectFormData = z.infer<typeof updateProjectSchema>

