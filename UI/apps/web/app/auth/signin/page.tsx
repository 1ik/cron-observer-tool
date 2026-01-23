'use client'

import { Button, Flex, Heading, Text } from '@radix-ui/themes'
import { signIn } from 'next-auth/react'

export default function SignInPage() {
  const handleSignIn = async () => {
    try {
      await signIn('google', { callbackUrl: '/' })
    } catch (error) {
      console.error('Sign in error:', error)
    }
  }

  return (
    <Flex
      direction="column"
      align="center"
      justify="center"
      style={{ minHeight: '100vh', padding: '2rem' }}
      gap="4"
    >
      <Heading size="8">Welcome to Cron Observer</Heading>
      <Text size="4" color="gray">
        Sign in with your Google account to continue
      </Text>
      <Button size="3" onClick={handleSignIn}>
        Sign in with Google
      </Button>
    </Flex>
  )
}

