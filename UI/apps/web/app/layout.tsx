import { Box, Flex, Separator, Theme } from '@radix-ui/themes'
import type { Metadata } from 'next'
import { ThemeProvider } from 'next-themes'
import { Header } from '../components/Header'
import { QueryProvider } from '../providers/QueryProvider'
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
      <body>
        <QueryProvider>
          <ThemeProvider attribute="class">
            <Theme accentColor="yellow" grayColor="gray" panelBackground="solid" radius="small" scaling="90%">
              <Flex direction="column" height="100vh" overflow="hidden">
                <Header />
                <Separator />
                <Box style={{ flex: 1, overflowY: 'auto' }}>
                  {children}
                </Box>
              </Flex>
            </Theme>
          </ThemeProvider>
        </QueryProvider>
      </body>
    </html>
  )
}

  