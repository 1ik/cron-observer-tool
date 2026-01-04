import { Task, ScheduleConfig } from '../types/task'

export type TaskRuntimeStatus = 'running' | 'success' | 'not-running' | 'paused' | 'disabled'

/**
 * Checks if a task is currently within its schedule window
 * For now, this is a simplified check. In production, this would need
 * to evaluate cron expressions and time ranges against the current time.
 */
export function isTaskWithinWindow(task: Task): boolean {
  // If task is not ACTIVE, it's not within window
  if (task.status !== 'ACTIVE') {
    return false
  }

  const config = task.schedule_config

  // If there's a cron expression, we'd need to evaluate it
  // For now, we'll assume tasks with cron expressions are "within window"
  // if they're ACTIVE (this is a simplification for the UI)
  if (config.cron_expression) {
    // TODO: Evaluate cron expression against current time
    // For demo purposes, assume ACTIVE tasks with cron are within window
    return true
  }

  // If there's a time range, check if current time is within it
  if (config.time_range) {
    const now = new Date()
    const timezone = config.timezone || 'UTC'
    
    // Convert current time to task timezone
    const nowInTimezone = new Date(now.toLocaleString('en-US', { timeZone: timezone }))
    const currentHour = nowInTimezone.getHours()
    const currentMinute = nowInTimezone.getMinutes()
    const currentTimeMinutes = currentHour * 60 + currentMinute

    // Parse time range
    const [startHour, startMin] = config.time_range.start.split(':').map(Number)
    const [endHour, endMin] = config.time_range.end.split(':').map(Number)
    const startTimeMinutes = startHour * 60 + startMin
    const endTimeMinutes = endHour * 60 + endMin

    // Check if current time is within range
    if (startTimeMinutes <= endTimeMinutes) {
      // Normal case: start < end (e.g., 09:00 - 17:00)
      return currentTimeMinutes >= startTimeMinutes && currentTimeMinutes < endTimeMinutes
    } else {
      // Overnight case: start > end (e.g., 22:00 - 06:00)
      return currentTimeMinutes >= startTimeMinutes || currentTimeMinutes < endTimeMinutes
    }
  }

  // If no specific schedule constraints, assume it's within window if ACTIVE
  return true
}

/**
 * Determines the runtime status of a task
 */
export function getTaskRuntimeStatus(task: Task): TaskRuntimeStatus {
  if (task.status === 'PAUSED') {
    return 'paused'
  }

  if (task.status === 'DISABLED') {
    return 'disabled'
  }

  if (task.status === 'ACTIVE') {
    // Check if task is within its schedule window
    if (isTaskWithinWindow(task)) {
      // For now, we'll use 'success' for active tasks within window
      // 'running' would require execution data to know if it's currently executing
      return 'success'
    } else {
      return 'not-running'
    }
  }

  return 'not-running'
}

/**
 * Gets the color for a task status dot
 */
export function getStatusDotColor(status: TaskRuntimeStatus): string {
  switch (status) {
    case 'running':
      return 'var(--blue-9)'
    case 'success':
      return 'var(--green-9)'
    case 'not-running':
      return 'var(--gray-9)'
    case 'paused':
      return 'var(--yellow-9)'
    case 'disabled':
      return 'var(--gray-9)'
    default:
      return 'var(--gray-9)'
  }
}

/**
 * Gets the tooltip text for a task status dot
 */
export function getStatusTooltip(status: TaskRuntimeStatus): string {
  switch (status) {
    case 'running':
      return 'Task is running'
    case 'success':
      return 'Success'
    case 'not-running':
      return 'Task is not running'
    case 'paused':
      return 'Task is paused'
    case 'disabled':
      return 'Task is disabled'
    default:
      return 'Unknown status'
  }
}

