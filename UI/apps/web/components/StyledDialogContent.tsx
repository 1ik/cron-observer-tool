'use client'

import { Box, Dialog } from '@radix-ui/themes'
import { ReactNode } from 'react'

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
  return (
    <Dialog.Content
      maxWidth={maxWidth}
      style={{
        maxHeight,
        width: '90vw',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
      }}
    >
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

