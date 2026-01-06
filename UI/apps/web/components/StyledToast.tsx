'use client'

import * as Toast from '@radix-ui/react-toast'
import { Box, Text } from '@radix-ui/themes'
import { ReactNode } from 'react'

interface StyledToastRootProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  type?: 'success' | 'error'
  title?: string
  children: ReactNode
}

interface StyledToastViewportProps {
  className?: string
}

/**
 * Styled Toast.Root component with theme-aware styling
 * Uses Radix UI Themes tokens for consistent theming
 */
export function StyledToastRoot({
  open,
  onOpenChange,
  type = 'success',
  title,
  children,
}: StyledToastRootProps) {
  return (
    <Toast.Root
      open={open}
      onOpenChange={onOpenChange}
      style={{
        backgroundColor: type === 'success' ? 'var(--green-9)' : 'var(--red-9)',
        padding: 'var(--space-3)',
        borderRadius: 'var(--radius-3)',
        boxShadow: 'var(--shadow-6)',
        minWidth: '300px',
        border: `1px solid ${type === 'success' ? 'var(--green-7)' : 'var(--red-7)'}`,
      }}
    >
      <Box>
        {title && (
          <Text
            size="3"
            weight="medium"
            style={{
              color: type === 'success' ? 'var(--green-contrast)' : 'var(--red-contrast)',
              marginBottom: 'var(--space-1)',
              display: 'block',
            }}
          >
            {title}
          </Text>
        )}
        <Text
          size="2"
          style={{
            color: type === 'success' ? 'var(--green-contrast)' : 'var(--red-contrast)',
            display: 'block',
          }}
        >
          {children}
        </Text>
      </Box>
    </Toast.Root>
  )
}

/**
 * Styled Toast.Viewport component with theme-aware positioning
 */
export function StyledToastViewport({ className }: StyledToastViewportProps) {
  return (
    <Toast.Viewport
      className={className}
      style={{
        position: 'fixed',
        bottom: 'var(--space-4)',
        right: 'var(--space-4)',
        zIndex: 9999,
        gap: 'var(--space-2)',
      }}
    />
  )
}

/**
 * Styled Toast.Provider component wrapper
 */
export function StyledToastProvider({
  children,
  swipeDirection = 'right',
}: {
  children: ReactNode
  swipeDirection?: 'right' | 'left' | 'up' | 'down'
}) {
  return <Toast.Provider swipeDirection={swipeDirection}>{children}</Toast.Provider>
}

