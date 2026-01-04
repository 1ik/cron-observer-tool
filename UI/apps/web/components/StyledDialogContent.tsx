'use client'

import * as Dialog from '@radix-ui/react-dialog'
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
  const style: CSSProperties = {
    position: 'fixed',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    maxWidth,
    width: '90vw',
    maxHeight,
    overflowY: 'auto',
    backgroundColor: 'var(--color-panel-solid)',
    border: '1px solid var(--gray-6)',
    borderRadius: 'var(--radius-3)',
    padding: 'var(--space-5)',
    boxShadow: 'var(--shadow-6)',
    zIndex: 51,
  }

  return (
    <Dialog.Content style={style}>
      {children}
    </Dialog.Content>
  )
}

