'use client'

import * as Dialog from '@radix-ui/react-dialog'
import { Box } from '@radix-ui/themes'
import { CSSProperties, ReactNode } from 'react'

interface StyledDialogContentProps {
  children: ReactNode
  maxWidth?: string
  maxHeight?: string
}

export function StyledDialogContent({
  children,
  maxWidth = '600px',
  maxHeight = '90vh',
}: StyledDialogContentProps) {
  const containerStyle: CSSProperties = {
    position: 'fixed',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    maxWidth,
    width: '90vw',
    maxHeight,
    backgroundColor: 'var(--color-panel-solid)',
    border: '1px solid var(--gray-6)',
    borderRadius: 'var(--radius-3)',
    boxShadow: 'var(--shadow-6)',
    zIndex: 51,
    display: 'flex',
    flexDirection: 'column',
    overflow: 'hidden',
  }

  return (
    <Dialog.Content style={containerStyle}>
      <Box
        style={{
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
          overflow: 'hidden',
        }}
      >
        {children}
      </Box>
    </Dialog.Content>
  )
}

