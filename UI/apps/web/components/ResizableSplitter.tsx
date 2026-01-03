'use client'

import { useState, useRef, useEffect } from 'react'
import { Box } from '@radix-ui/themes'

interface ResizableSplitterProps {
  leftContent: React.ReactNode
  rightContent: React.ReactNode
  initialLeftWidth?: number // percentage (0-100)
  minLeftWidth?: number // percentage
  minRightWidth?: number // percentage
}

export function ResizableSplitter({
  leftContent,
  rightContent,
  initialLeftWidth = 40,
  minLeftWidth = 20,
  minRightWidth = 20,
}: ResizableSplitterProps) {
  const [leftWidth, setLeftWidth] = useState(initialLeftWidth)
  const [isDragging, setIsDragging] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const startXRef = useRef<number>(0)
  const startLeftWidthRef = useRef<number>(initialLeftWidth)

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isDragging || !containerRef.current) return

      const container = containerRef.current
      const containerWidth = container.offsetWidth
      const deltaX = e.clientX - startXRef.current
      const deltaPercent = (deltaX / containerWidth) * 100

      let newLeftWidth = startLeftWidthRef.current + deltaPercent

      // Enforce minimum widths
      const maxLeftWidth = 100 - minRightWidth
      newLeftWidth = Math.max(minLeftWidth, Math.min(maxLeftWidth, newLeftWidth))

      setLeftWidth(newLeftWidth)
    }

    const handleMouseUp = () => {
      setIsDragging(false)
    }

    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove)
      document.addEventListener('mouseup', handleMouseUp)
      document.body.style.cursor = 'col-resize'
      document.body.style.userSelect = 'none'
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseup', handleMouseUp)
      document.body.style.cursor = ''
      document.body.style.userSelect = ''
    }
  }, [isDragging, minLeftWidth, minRightWidth])

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault()
    setIsDragging(true)
    startXRef.current = e.clientX
    startLeftWidthRef.current = leftWidth
  }

  return (
    <Box
      ref={containerRef}
      style={{
        display: 'flex',
        width: '100%',
        height: '100%',
        position: 'relative',
      }}
    >
      {/* Left column */}
      <Box
        style={{
          flex: `0 0 ${leftWidth}%`,
          overflowY: 'auto',
          height: '100%',
        }}
      >
        {leftContent}
      </Box>

      {/* Resizable separator */}
      <Box
        onMouseDown={handleMouseDown}
        style={{
          width: '4px',
          flexShrink: 0,
          backgroundColor: 'var(--gray-6)',
          cursor: 'col-resize',
          position: 'relative',
          userSelect: 'none',
        }}
      >
        {/* Hover indicator */}
        <Box
          style={{
            position: 'absolute',
            top: 0,
            left: '-2px',
            right: '-2px',
            bottom: 0,
            backgroundColor: isDragging ? 'var(--gray-8)' : 'transparent',
            transition: isDragging ? 'none' : 'background-color 0.2s',
          }}
        />
      </Box>

      {/* Right column */}
      <Box
        style={{
          flex: `0 0 ${100 - leftWidth}%`,
          overflowY: 'auto',
          height: '100%',
        }}
      >
        {rightContent}
      </Box>
    </Box>
  )
}

