'use client'

import { useTasksByProject, useUpdateTaskStatus } from '@cron-observer/lib'
import { CalendarIcon, PauseIcon, PlayIcon } from '@radix-ui/react-icons'
import * as Popover from '@radix-ui/react-popover'
import { Box, Button, Flex, IconButton, Spinner, Text, Tooltip } from '@radix-ui/themes'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import { useMemo, useState } from 'react'
import { DayPicker } from 'react-day-picker'
import 'react-day-picker/dist/style.css'
import { Execution } from '../lib/types/execution'
import { ExecutionItem } from './ExecutionItem'
import { StyledToastProvider, StyledToastRoot, StyledToastViewport } from './StyledToast'

interface ExecutionsListProps {
  executions: Execution[]
  isLoading?: boolean
  selectedTaskId?: string | null
  projectId?: string | null
}

export function ExecutionsList({ executions, isLoading = false, selectedTaskId, projectId }: ExecutionsListProps) {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  
  // Get tasks to find the selected task's status
  const { data: tasks = [] } = useTasksByProject(projectId || '')
  const selectedTask = selectedTaskId
    ? tasks.find((t) => t.id === selectedTaskId || t.uuid === selectedTaskId)
    : null
  const selectedTaskUUID = selectedTask?.uuid || selectedTask?.id || null
  const currentTaskStatus = selectedTask?.status || 'ACTIVE'
  const isTaskPaused = currentTaskStatus === 'PAUSED'
  
  // Toast state
  const [toastOpen, setToastOpen] = useState(false)
  const [toastMessage, setToastMessage] = useState('')
  const [toastType, setToastType] = useState<'success' | 'error'>('success')
  
  // Update task status mutation
  const updateStatusMutation = useUpdateTaskStatus(projectId || '', selectedTaskUUID || '')

  // Extract unique dates from executions
  const availableDates = useMemo(() => {
    const dateSet = new Set<string>()
    executions.forEach((execution) => {
      const date = new Date(execution.started_at).toISOString().split('T')[0]
      dateSet.add(date)
    })
    
    // For demo purposes, if we have fewer than 5 dates, add some demo dates
    const dates = Array.from(dateSet).sort().reverse()
    if (dates.length < 5) {
      // Add some demo dates going back
      const mostRecent = dates[0] ? new Date(dates[0]) : new Date()
      for (let i = dates.length; i < 6; i++) {
        const date = new Date(mostRecent)
        date.setDate(date.getDate() - (i - dates.length))
        dates.push(date.toISOString().split('T')[0])
      }
      return dates.sort().reverse()
    }
    return dates
  }, [executions])

  // Get selected date from URL or default to most recent
  const selectedDateString = searchParams.get('date') || availableDates[0] || ''
  // Parse date string (YYYY-MM-DD) to Date object at local midnight to avoid timezone issues
  const selectedDate = selectedDateString ? (() => {
    const [year, month, day] = selectedDateString.split('-').map(Number)
    return new Date(year, month - 1, day)
  })() : new Date()

  const handleDateChange = (date: Date | undefined) => {
    if (!date) return
    
    // Format date in local timezone (YYYY-MM-DD) to avoid timezone shift
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const dateString = `${year}-${month}-${day}`
    
    const params = new URLSearchParams(searchParams.toString())
    params.set('date', dateString)
    const queryString = params.toString()
    router.push(`${pathname}${queryString ? `?${queryString}` : ''}`, { scroll: false })
  }

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  }

  const handleToggle = async () => {
    if (!selectedTaskUUID || !projectId) {
      return
    }
    
    const newStatus = isTaskPaused ? 'ACTIVE' : 'PAUSED'
    
    try {
      await updateStatusMutation.mutateAsync(newStatus)
      setToastMessage(newStatus === 'PAUSED' ? 'Task paused successfully' : 'Task resumed successfully')
      setToastType('success')
      setToastOpen(true)
    } catch (error) {
      setToastMessage(error instanceof Error ? error.message : 'Failed to update task status')
      setToastType('error')
      setToastOpen(true)
    }
  }
  
  const isLoadingStatus = updateStatusMutation.isPending

  return (
    <Flex
      direction="column"
      style={{
        width: '100%',
        height: '100%',
        overflow: 'hidden',
      }}
    >
      {/* Sticky header */}
      <Box
        py="3"
        pl="3"
        pr="4"
        style={{
          flexShrink: 0,
          borderBottom: '1px solid var(--gray-6)',
        }}
      >
        <Flex justify="between" align="center" gap="3">
          <Flex align="center" gap="2">
            <Text size="2" weight="medium" color="gray">
              Executions
            </Text>
            <Popover.Root>
              <Popover.Trigger asChild>
                <Button
                  variant="soft"
                  size="2"
                  style={{
                    minWidth: '140px',
                    justifyContent: 'space-between',
                    cursor: 'pointer',
                  }}
                >
                  <Flex align="center" gap="2">
                    <CalendarIcon 
                      width="14" 
                      height="14" 
                      style={{ color: 'var(--gray-11)' }}
                    />
                    <Text size="2" color="gray">
                      {formatDate(selectedDate)}
                    </Text>
                  </Flex>
                </Button>
              </Popover.Trigger>
              <Popover.Content
                side="bottom"
                align="start"
                style={{
                  padding: 'var(--space-3)',
                  backgroundColor: 'var(--color-panel-solid)',
                  border: '1px solid var(--gray-6)',
                  borderRadius: 'var(--radius-3)',
                  boxShadow: 'var(--shadow-6)',
                  zIndex: 100,
                }}
              >
                <Box style={{ fontFamily: 'var(--font-family)', color: 'var(--gray-12)' }}>
                  <DayPicker
                    mode="single"
                    selected={selectedDate}
                    onSelect={handleDateChange}
                    className="rdp"
                    classNames={{
                      months: 'rdp-months',
                      month: 'rdp-month',
                      caption: 'rdp-caption',
                      caption_label: 'rdp-caption_label',
                      nav: 'rdp-nav',
                      button_previous: 'rdp-button_previous',
                      button_next: 'rdp-button_next',
                      table: 'rdp-table',
                      head_row: 'rdp-head_row',
                      head_cell: 'rdp-head_cell',
                      row: 'rdp-row',
                      cell: 'rdp-cell',
                      day: 'rdp-day',
                      day_disabled: 'rdp-day_disabled',
                      day_selected: 'rdp-day_selected',
                      day_today: 'rdp-day_today',
                    }}
                    styles={{
                      root: {
                        color: 'var(--gray-12)',
                      },
                      month: {
                        color: 'var(--gray-12)',
                      },
                      caption_label: {
                        color: 'var(--gray-12)',
                        fontSize: 'var(--font-size-3)',
                        fontWeight: 'var(--font-weight-medium)',
                      },
                      day: {
                        color: 'var(--gray-12)',
                        borderRadius: 'var(--radius-2)',
                      },
                      day_selected: {
                        backgroundColor: 'var(--accent-9)',
                        color: 'var(--accent-contrast)',
                        fontWeight: 'var(--font-weight-medium)',
                      },
                      day_today: {
                        fontWeight: 'var(--font-weight-bold)',
                      },
                      day_disabled: {
                        opacity: 0.3,
                        cursor: 'not-allowed',
                      },
                      head_cell: {
                        color: 'var(--gray-11)',
                        fontSize: 'var(--font-size-1)',
                        fontWeight: 'var(--font-weight-medium)',
                      },
                      button_previous: {
                        color: 'var(--gray-12)',
                      },
                      button_next: {
                        color: 'var(--gray-12)',
                      },
                    }}
                  />
                </Box>
              </Popover.Content>
            </Popover.Root>
          </Flex>
          {selectedTaskUUID && projectId && (
            <Tooltip content={isTaskPaused ? "Resume task" : "Pause task"}>
              <IconButton
                variant="outline"
                size="2"
                onClick={handleToggle}
                disabled={isLoadingStatus}
                style={{ cursor: isLoadingStatus ? 'wait' : 'pointer', marginRight: 'var(--space-1)' }}
              >
                {isLoadingStatus ? (
                  <Spinner size="2" />
                ) : isTaskPaused ? (
                  <PlayIcon width="16" height="16" />
                ) : (
                  <PauseIcon width="16" height="16" />
                )}
              </IconButton>
            </Tooltip>
          )}
        </Flex>
      </Box>

      {/* Scrollable body */}
      <Box
        style={{
          flex: 1,
          overflowY: 'auto',
          padding: 'var(--space-3)',
        }}
      >
        {isLoading ? (
          <Flex justify="center" align="center" style={{ minHeight: '200px' }}>
            <Flex direction="column" gap="3" align="center">
              <Spinner size="3" />
              <Text size="2" color="gray">
                Loading executions...
              </Text>
            </Flex>
          </Flex>
        ) : executions.length === 0 ? (
          <Flex justify="center" align="center" style={{ minHeight: '200px' }}>
            <Text
              size="2"
              color="gray"
              style={{
                fontFamily: 'var(--font-mono)',
              }}
            >
              No executions found
            </Text>
          </Flex>
        ) : (
          <Box
            style={{
              border: '1px solid var(--gray-4)',
              borderRadius: 'var(--radius-2)',
              overflow: 'hidden',
              backgroundColor: 'var(--color-panel)',
            }}
          >
            {executions.map((execution, index) => (
              <ExecutionItem key={execution.id || index} execution={execution} />
            ))}
          </Box>
        )}
      </Box>
      
      {/* Toast notifications */}
      <StyledToastProvider swipeDirection="right">
        <StyledToastRoot
          open={toastOpen}
          onOpenChange={setToastOpen}
          type={toastType}
          title={toastType === 'success' ? 'Success' : 'Error'}
        >
          {toastMessage}
        </StyledToastRoot>
        <StyledToastViewport />
      </StyledToastProvider>
    </Flex>
  )
}
