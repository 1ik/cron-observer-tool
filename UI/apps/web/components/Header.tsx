'use client'

import Link from 'next/link'
import { Box, Flex, IconButton, Text } from '@radix-ui/themes'
import { GearIcon, BellIcon } from '@radix-ui/react-icons'
import * as DropdownMenu from '@radix-ui/react-dropdown-menu'
import * as Avatar from '@radix-ui/react-avatar'

export function Header() {
  return (
    <Box height="60px" style={{ flexShrink: 0 }}>
      <Flex justify="between" align="center" height="100%" px="4">
        <Link href="/" style={{ textDecoration: 'none' }}>
          <Text size="6" weight="bold">Cron Observer</Text>
        </Link>
        <Flex gap="4" align="center">
          <IconButton variant="outline" size="3">
            <GearIcon width="20" height="20" />
          </IconButton>
          <IconButton variant="outline" size="3">
            <BellIcon width="20" height="20" />
          </IconButton>
          <DropdownMenu.Root>
            <DropdownMenu.Trigger asChild>
              <IconButton variant="outline" size="3">
                <Avatar.Root>
                  <Avatar.Fallback>U</Avatar.Fallback>
                </Avatar.Root>
              </IconButton>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
              <DropdownMenu.Item>Profile</DropdownMenu.Item>
              <DropdownMenu.Item>Settings</DropdownMenu.Item>
              <DropdownMenu.Separator />
              <DropdownMenu.Item>Sign out</DropdownMenu.Item>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </Flex>
      </Flex>
    </Box>
  )
}

