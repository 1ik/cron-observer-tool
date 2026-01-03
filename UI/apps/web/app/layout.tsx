import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Cron Observer',
  description: 'Task scheduling and execution tracking system',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}

