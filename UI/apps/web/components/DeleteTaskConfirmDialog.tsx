'use client'

import { AlertDialog, Box, Button, Flex, Heading, Separator, Text } from '@radix-ui/themes'

interface DeleteTaskConfirmDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  taskName: string
  onConfirm: () => void
  isDeleting?: boolean
}

export function DeleteTaskConfirmDialog({
  open,
  onOpenChange,
  taskName,
  onConfirm,
  isDeleting = false,
}: DeleteTaskConfirmDialogProps) {
  return (
    <AlertDialog.Root open={open} onOpenChange={onOpenChange}>
      <AlertDialog.Content maxWidth="500px">
        <Flex
          direction="column"
          style={{
            height: '100%',
            overflow: 'hidden',
            width: '100%',
          }}
        >
          {/* Header */}
          <Box p="5" style={{ flexShrink: 0 }}>
            <AlertDialog.Title>
              <Text size="5" weight="bold" mb="2" color="red" style={{ display: 'block' }}>
                Delete Task
              </Text>
            </AlertDialog.Title>
            <AlertDialog.Description>
              <Text size="3" color="gray">
                Are you sure you want to delete this task?
              </Text>
            </AlertDialog.Description>
          </Box>

          <Separator size="4" my="0" />

          {/* Content */}
          <Box
            p="5"
            style={{
              flex: 1,
              overflowY: 'auto',
              minHeight: 0,
            }}
          >
            <Text size="3" color="red" weight="medium" mb="3">
              Task: {taskName}
            </Text>
            <Text size="3" color="gray">
              The operation is dangerous and irreversible. This task will be permanently deleted.
            </Text>
          </Box>

          <Separator size="4" my="0" />

          {/* Footer */}
          <Box p="5" style={{ flexShrink: 0 }}>
            <Flex gap="3" justify="end">
              <AlertDialog.Cancel>
                <Button type="button" variant="soft" disabled={isDeleting}>
                  Cancel
                </Button>
              </AlertDialog.Cancel>
              <AlertDialog.Action>
                <Button
                  type="button"
                  variant="solid"
                  color="red"
                  onClick={onConfirm}
                  disabled={isDeleting}
                >
                  {isDeleting ? 'Deleting...' : 'Confirm Delete'}
                </Button>
              </AlertDialog.Action>
            </Flex>
          </Box>
        </Flex>
      </AlertDialog.Content>
    </AlertDialog.Root>
  )
}
