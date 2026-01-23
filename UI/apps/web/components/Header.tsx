'use client'

import * as Avatar from '@radix-ui/react-avatar'
import { BellIcon, GearIcon } from '@radix-ui/react-icons'
import * as Popover from '@radix-ui/react-popover'
import { Box, Button, Flex, IconButton, Separator, Text } from '@radix-ui/themes'
import { signOut, useSession } from 'next-auth/react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'

export function Header() {
  const { data: session, status } = useSession()
  const router = useRouter()

  const handleSignOut = async () => {
    await signOut({ redirectTo: '/auth/signin' })
  }

  const handleSignIn = () => {
    router.push('/auth/signin')
  }

  const isLoading = status === 'loading'
  const isAuthenticated = !!session

  // Get user initials for avatar
  const getInitials = () => {
    if (session?.user?.name) {
      return session.user.name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    }
    if (session?.user?.email) {
      return session.user.email[0].toUpperCase()
    }
    return 'U'
  }

  return (
    <Box height="60px" style={{ flexShrink: 0 }}>
      <Flex justify="between" align="center" height="100%" px="4">
        <Link href="/" style={{ textDecoration: 'none', cursor: 'pointer' }}>
          <Text size="6" weight="bold">Cron Observer</Text>
        </Link>
        <Flex gap="4" align="center">
          {isLoading ? (
            <Text size="2" color="gray">Loading...</Text>
          ) : isAuthenticated ? (
            <>
              <IconButton variant="outline" size="3" style={{ cursor: 'pointer' }}>
                <GearIcon width="20" height="20" />
              </IconButton>
              <IconButton variant="outline" size="3" style={{ cursor: 'pointer' }}>
                <BellIcon width="20" height="20" />
              </IconButton>
              <Popover.Root>
                <Popover.Trigger asChild>
                  <button
                    style={{
                      background: 'none',
                      border: 'none',
                      padding: 0,
                      cursor: 'pointer',
                    }}
                  >
                    <Avatar.Root
                      style={{
                        display: 'inline-flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        verticalAlign: 'middle',
                        overflow: 'hidden',
                        userSelect: 'none',
                        width: '40px',
                        height: '40px',
                        borderRadius: '100%',
                        backgroundColor: 'var(--gray-5)',
                      }}
                    >
                      <Avatar.Image
                        src={session.user?.image || ''}
                        alt={session.user?.name || 'User avatar'}
                        style={{
                          width: '100%',
                          height: '100%',
                          objectFit: 'cover',
                          borderRadius: 'inherit',
                        }}
                      />
                      <Avatar.Fallback
                        style={{
                          width: '100%',
                          height: '100%',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          backgroundColor: 'var(--gray-9)',
                          color: 'white',
                          fontSize: '16px',
                          fontWeight: 500,
                        }}
                      >
                        {getInitials()}
                      </Avatar.Fallback>
                    </Avatar.Root>
                  </button>
                </Popover.Trigger>
                <Popover.Content
                  side="bottom"
                  align="end"
                  style={{
                    padding: 'var(--space-2)',
                    backgroundColor: 'var(--color-panel-solid)',
                    border: '1px solid var(--gray-6)',
                    borderRadius: 'var(--radius-3)',
                    boxShadow: 'var(--shadow-6)',
                    zIndex: 100,
                    minWidth: '200px',
                  }}
                >
                  <Box p="3">
                    <Box mb="3">
                      <Text size="2" weight="bold" as="div">
                        {session.user?.name || 'User'}
                      </Text>
                      <Text size="1" color="gray" as="div">
                        {session.user?.email}
                      </Text>
                    </Box>
                    <Separator mb="2" />
                    <Flex direction="column" gap="1">
                      <Button
                        variant="ghost"
                        size="2"
                        style={{ justifyContent: 'flex-start' }}
                      >
                        Profile
                      </Button>
                      <Button
                        variant="ghost"
                        size="2"
                        style={{ justifyContent: 'flex-start' }}
                      >
                        Settings
                      </Button>
                      <Separator my="1" />
                      <Button
                        variant="ghost"
                        size="2"
                        color="red"
                        style={{ justifyContent: 'flex-start' }}
                        onClick={handleSignOut}
                      >
                        Sign out
                      </Button>
                    </Flex>
                  </Box>
                </Popover.Content>
              </Popover.Root>
            </>
          ) : (
            <Button size="2" onClick={handleSignIn}>
              Sign in
            </Button>
          )}
        </Flex>
      </Flex>
    </Box>
  )
}

