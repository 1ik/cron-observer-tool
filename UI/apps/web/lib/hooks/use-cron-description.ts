import { useMemo } from 'react'
import cronstrue from 'cronstrue'

/**
 * Hook to get human-readable description of a cron expression
 * @param cronExpression - The cron expression to describe
 * @returns Human-readable description or error message
 */
export function useCronDescription(cronExpression: string | undefined): string | null {
  return useMemo(() => {
    if (!cronExpression || cronExpression.trim() === '') {
      return null
    }

    try {
      // Use throwExceptionOnParseError: false to catch errors gracefully
      const description = cronstrue.toString(cronExpression, {
        throwExceptionOnParseError: false,
        verbose: false,
        use24HourTimeFormat: true,
      })
      // cronstrue may return an error message string when parsing fails
      // Check if it's an error message (starts with "Error" or contains "invalid")
      if (description && !description.toLowerCase().includes('error') && !description.toLowerCase().includes('invalid')) {
        return description
      }
      return null
    } catch (error) {
      // If parsing fails, return null (don't show error message)
      return null
    }
  }, [cronExpression])
}

