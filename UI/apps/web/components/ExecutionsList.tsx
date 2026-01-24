'use client'

import { useTasksByProject, useTriggerTask } from '@cron-observer/lib'
import { useToast } from '@cron-observer/ui'
import { CalendarIcon, ChevronLeftIcon, ChevronRightIcon } from '@radix-ui/react-icons'
import * as Popover from '@radix-ui/react-popover'
import { Box, Button, Flex, IconButton, Spinner, Text, Tooltip } from '@radix-ui/themes'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import { useMemo } from 'react'
import { DayPicker } from 'react-day-picker'
import 'react-day-picker/dist/style.css'
import { Execution } from '../lib/types/execution'
import { ExecutionItem } from './ExecutionItem'

interface ExecutionsListProps {
  executions: Execution[]
  isLoading?: boolean
  selectedTaskId?: string | null
  projectId?: string | null
  paginationData?: {
    page: number
    page_size: number
    total_count: number
    total_pages: number
  } | null
}

export function ExecutionsList({ executions, isLoading = false, selectedTaskId, projectId, paginationData }: ExecutionsListProps) {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  
  // Get current page from URL or pagination data
  const currentPage = useMemo(() => {
    const pageParam = searchParams.get('page')
    if (pageParam) {
      const parsed = parseInt(pageParam, 10)
      return isNaN(parsed) || parsed < 1 ? 1 : parsed
    }
    return paginationData?.page || 1
  }, [searchParams, paginationData])

  const pageSize = paginationData?.page_size || 100
  const totalCount = paginationData?.total_count || 0
  const totalPages = paginationData?.total_pages || 0

  const handlePageChange = (newPage: number) => {
    if (newPage < 1 || newPage > totalPages) return
    
    const params = new URLSearchParams(searchParams.toString())
    params.set('page', newPage.toString())
    router.push(`${pathname}?${params.toString()}`, { scroll: false })
  }
  
  // Get tasks to find the selected task's status
  const { data: tasks = [] } = useTasksByProject(projectId || '')
  const selectedTask = selectedTaskId
    ? tasks.find((t) => t.id === selectedTaskId || t.uuid === selectedTaskId)
    : null
  const selectedTaskUUID = selectedTask?.uuid || selectedTask?.id || null
  
  // Toast hook for imperative API
  const toast = useToast()
  
  // Trigger task mutation
  const triggerTaskMutation = useTriggerTask(projectId || '', selectedTaskUUID || '')

  // Generate available dates (last 30 days) for date picker
  // Since we're using pagination, we can't extract dates from executions
  const availableDates = useMemo(() => {
    const dates: string[] = []
    const today = new Date()
    for (let i = 0; i < 30; i++) {
      const date = new Date(today)
      date.setDate(date.getDate() - i)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      dates.push(`${year}-${month}-${day}`)
    }
    return dates
  }, [])

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
    // Reset to page 1 when date changes
    params.set('page', '1')
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

  const handleTrigger = async () => {
    if (!selectedTaskUUID || !projectId) {
      return
    }
    
    try {
      await triggerTaskMutation.mutateAsync()
      toast.success('Task triggered successfully')
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : 'Failed to trigger task',
        'Error'
      )
    }
  }
  
  const isTriggering = triggerTaskMutation.isPending

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
            <Tooltip content="Trigger task manually">
              <IconButton
                variant="outline"
                size="2"
                onClick={handleTrigger}
                disabled={isTriggering}
                style={{ cursor: isTriggering ? 'wait' : 'pointer', marginRight: 'var(--space-1)' }}
              >
                {isTriggering ? (
                  <Spinner size="2" />
                ) : (
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 16 16"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                    style={{ stroke: 'currentColor', strokeWidth: 1.5, fill: 'none' }}
                  >
                    <path
                      d="M8.5 2L4.5 9H8.5L7.5 14L11.5 7H7.5L8.5 2Z"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                )}
              </IconButton>
            </Tooltip>
          )}
        </Flex>
      </Box>

      {/* Scrollable body with sticky pagination */}
      <Flex
        direction="column"
        style={{
          flex: 1,
          overflow: 'hidden',
          minHeight: 0,
        }}
      >
        {/* Scrollable content */}
        <Box
          style={{
            flex: 1,
            overflowY: 'auto',
            padding: 'var(--space-3)',
            minHeight: 0,
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

        {/* Sticky pagination footer */}
        {paginationData && totalPages > 0 && (
          <Box
            style={{
              flexShrink: 0,
              borderTop: '1px solid var(--gray-6)',
              padding: 'var(--space-3)',
              backgroundColor: 'var(--color-panel)',
            }}
          >
            <Flex justify="between" align="center" gap="3">
              <Text size="2" color="gray">
                Showing {((currentPage - 1) * pageSize) + 1} to {Math.min(currentPage * pageSize, totalCount)} of {totalCount} executions
              </Text>
              <Flex align="center" gap="2">
                <IconButton
                  variant="soft"
                  size="2"
                  onClick={() => handlePageChange(currentPage - 1)}
                  disabled={currentPage <= 1 || isLoading}
                  style={{ cursor: currentPage <= 1 ? 'not-allowed' : 'pointer' }}
                >
                  <ChevronLeftIcon width="16" height="16" />
                </IconButton>
                <Text size="2" color="gray">
                  Page {currentPage} of {totalPages}
                </Text>
                <IconButton
                  variant="soft"
                  size="2"
                  onClick={() => handlePageChange(currentPage + 1)}
                  disabled={currentPage >= totalPages || isLoading}
                  style={{ cursor: currentPage >= totalPages ? 'not-allowed' : 'pointer' }}
                >
                  <ChevronRightIcon width="16" height="16" />
                </IconButton>
              </Flex>
            </Flex>
          </Box>
        )}
      </Flex>
    </Flex>
  )
}
