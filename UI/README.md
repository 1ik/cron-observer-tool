# Cron Observer UI

Frontend monorepo for Cron Observer built with pnpm workspaces.

## Structure

```
UI/
├── apps/
│   └── web/              # Next.js 14 main application
├── packages/
│   ├── ui/               # Shared UI components (Radix UI)
│   └── lib/              # Shared utilities and hooks
└── package.json          # Root workspace configuration
```

## Getting Started

### Prerequisites

- Node.js >= 18.0.0
- pnpm >= 8.0.0

### Installation

```bash
# Install all dependencies
pnpm install
```

### Development

```bash
# Run all apps in development mode
pnpm dev

# Run specific app
pnpm --filter @cron-observer/web dev
```

### Building

```bash
# Build all packages and apps
pnpm build

# Build specific package/app
pnpm --filter @cron-observer/ui build
pnpm --filter @cron-observer/web build
```

### Type Checking

```bash
# Type check all packages
pnpm type-check
```

## Packages

### `@cron-observer/ui`

Shared UI components built on Radix UI primitives.

**Usage:**
```typescript
import { Button } from '@cron-observer/ui';
```

### `@cron-observer/lib`

Shared utilities and React hooks.

**Usage:**
```typescript
import { someUtil } from '@cron-observer/lib';
```

### `@cron-observer/web`

Main Next.js application.

## Workspace Commands

- `pnpm dev` - Run all apps in dev mode
- `pnpm build` - Build all packages and apps
- `pnpm lint` - Lint all packages
- `pnpm type-check` - Type check all packages
- `pnpm clean` - Clean all build artifacts

## Adding New Packages

1. Create directory in `packages/` or `apps/`
2. Add `package.json` with workspace name (e.g., `@cron-observer/package-name`)
3. Create `tsconfig.json` that extends root config
4. Install dependencies: `pnpm install`

## TypeScript Configuration

All packages extend the base `tsconfig.json` at the root. Package-specific overrides can be added in each package's `tsconfig.json`.

