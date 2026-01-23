'use client'

import { setTokenProvider } from '@cron-observer/lib'
import { SessionProvider as NextAuthSessionProvider } from 'next-auth/react'
import { useEffect } from 'react'

export function SessionProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    // Set up token provider for API client
    setTokenProvider(async () => {
      try {
        const response = await fetch('/api/auth/token')
        if (response.ok) {
          const data = await response.json()
          return data.token
        }
      } catch (error) {
        console.error('Failed to get token:', error)
      }
      return null
    })
  }, [])

  return <NextAuthSessionProvider>{children}</NextAuthSessionProvider>
}

