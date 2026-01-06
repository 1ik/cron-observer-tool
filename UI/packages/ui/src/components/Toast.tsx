'use client'

import * as ToastPrimitive from '@radix-ui/react-toast'
import { Box, Text } from '@radix-ui/themes'
import { createContext, useContext, useState, useCallback, ReactNode, useRef, useEffect } from 'react'

// ============================================================================
// Styled Toast Components
// ============================================================================

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
    <ToastPrimitive.Root
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
    </ToastPrimitive.Root>
  )
}

/**
 * Styled Toast.Viewport component with theme-aware positioning
 */
export function StyledToastViewport({ className }: StyledToastViewportProps) {
  return (
    <ToastPrimitive.Viewport
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
  return <ToastPrimitive.Provider swipeDirection={swipeDirection}>{children}</ToastPrimitive.Provider>
}

// ============================================================================
// Imperative Toast API
// ============================================================================

interface ToastItem {
  id: string
  type: 'success' | 'error'
  message: string
  title?: string
  open: boolean
}

interface ToastContextType {
  success: (message: string, title?: string) => void
  error: (message: string, title?: string) => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

interface ToastProviderProps {
  children: ReactNode
  swipeDirection?: 'right' | 'left' | 'up' | 'down'
  duration?: number // Auto-dismiss duration in milliseconds
}

/**
 * ToastProvider - Provides imperative toast API via context
 * 
 * Usage:
 * ```tsx
 * // In root layout
 * <ToastProvider>
 *   <App />
 * </ToastProvider>
 * 
 * // In components
 * const toast = useToast()
 * toast.success('Task paused successfully')
 * toast.error('Failed to update task')
 * ```
 */
export function ToastProvider({ 
  children, 
  swipeDirection = 'right',
  duration = 5000 
}: ToastProviderProps) {
  const [toasts, setToasts] = useState<ToastItem[]>([])
  const toastIdCounter = useRef(0)

  const showToast = useCallback((type: 'success' | 'error', message: string, title?: string) => {
    const id = `toast-${toastIdCounter.current++}`
    const newToast: ToastItem = { 
      id, 
      type, 
      message, 
      title, 
      open: true 
    }
    
    setToasts((prev) => [...prev, newToast])
    
    // Auto-dismiss after duration
    setTimeout(() => {
      setToasts((prev) => 
        prev.map((t) => t.id === id ? { ...t, open: false } : t)
      )
      // Remove from array after close animation
      setTimeout(() => {
        setToasts((prev) => prev.filter((t) => t.id !== id))
      }, 300)
    }, duration)
  }, [duration])

  const success = useCallback((message: string, title?: string) => {
    showToast('success', message, title)
  }, [showToast])

  const error = useCallback((message: string, title?: string) => {
    showToast('error', message, title)
  }, [showToast])

  const handleToastClose = useCallback((id: string, open: boolean) => {
    if (!open) {
      setToasts((prev) => prev.filter((t) => t.id !== id))
    }
  }, [])

  return (
    <ToastContext.Provider value={{ success, error }}>
      <StyledToastProvider swipeDirection={swipeDirection}>
        {children}
        {toasts.map((toast) => (
          <StyledToastRoot
            key={toast.id}
            open={toast.open}
            onOpenChange={(open) => handleToastClose(toast.id, open)}
            type={toast.type}
            title={toast.title || (toast.type === 'success' ? 'Success' : 'Error')}
          >
            {toast.message}
          </StyledToastRoot>
        ))}
        <StyledToastViewport />
      </StyledToastProvider>
    </ToastContext.Provider>
  )
}

/**
 * Hook to access imperative toast API
 * 
 * @example
 * ```tsx
 * const toast = useToast()
 * 
 * const handleSuccess = () => {
 *   toast.success('Operation completed successfully')
 * }
 * 
 * const handleError = () => {
 *   toast.error('Something went wrong', 'Error')
 * }
 * ```
 */
export function useToast() {
  const context = useContext(ToastContext)
  if (!context) {
    throw new Error('useToast must be used within ToastProvider')
  }
  return context
}

