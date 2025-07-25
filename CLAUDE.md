# Instructions

You run in an environment where `ast-grep` is available; whenever a search requires syntax-aware or structural matching, default to `ast-grep -p '<pattern>' --lang go` (or set `--lang` appropriately for example use `ts` for `TypeScript`) and avoid falling back to text-only tools like `rg` or `grep` unless I explicitly request a plain-text search.

Keep code comments minimum and relevant to the code itself. Do not add comments that are not directly related to the code or that explain obvious things.

# Tech Stack

## Backend

- Go + Echo framework
- PostgreSQL + pgx driver
- Goose for migrations
- Testify for unit tests

## Frontend

- React 19 + TypeScript
- Vite
- TanStack Router
- Tailwind v4 + shadcn/ui
- TanStack Query
- i18next for internationalization
- Zod for validation
- Zustand for state management
- React Hook Form for forms
- Lucide React for icons
- Sonner for notifications

## Architecture

- Monolith with embedded frontend (single binary deployment)

# Development Commands

## Backend Development

- `task run:server` - Start backend with hot reload using wgo
- `task test` - Run Go tests
- `task lint` - Run lint
- `task update` - Update Go dependencies

## Frontend Development

- `task run:web` - Start frontend dev server with hot reload
- `cd web && npm run build` - Build frontend for production
- `cd web && npm run lint` - Run ESLint
- `cd web && npm run lint:fix` - Fix ESLint issues automatically
- `cd web && npm run router:generate` - Generate TanStack Router route tree

## Build Commands

- `task build:web` - Build frontend only
- `task build:linux` - Build Linux binary with embedded frontend
- `task build:win` - Build Windows binary with embedded frontend
- `task build:docker` - Build Docker image

# Architecture Overview

## Service Layer Pattern

The backend follows a service-oriented architecture with clear separation:

- **Handlers** (`internal/api/*.go`) - HTTP request/response handling, validation
- **Services** (`internal/service/*/`) - Business logic and data operations
- **Database** (`internal/database/`) - Connection pooling and utilities

## Dependency Injection

Services are instantiated in `cmd/monolith/main.go` and injected into handlers:

```go
userService := user.NewService(db)
accountService := account.NewService(db)
authService := auth.NewService(db)
sessionService := session.NewService(db)
```

## Frontend Architecture

- **Feature-based structure** under `web/src/features/`
- **Components** organized under `web/src/components/`
- **TanStack Router** for file-based routing in `web/src/routes/`
- **Types** for API responses and shared types in `web/src/types/`
- **Context providers** for global state management in `web/src/context/`
- **Hooks** for reusable logic in `web/src/hooks/`
- **Utilities** in `web/src/utils/`
- **Zustand** for global state management (`web/src/hooks/`)

# Development Standards & Best Practices

## PostgreSQL Database Design

### Column Types

- **Identity columns**: Always use `GENERATED ALWAYS AS IDENTITY` instead of `SERIAL` or `BIGSERIAL`
- **Timestamps**: Always use `TIMESTAMPTZ` for datetime columns to handle timezone-aware dates
- **Text data**: Always use `TEXT` instead of `VARCHAR(N)` for variable-length strings
- **JSON data**: Always use `JSONB` instead of `JSON` for better performance and indexing capabilities

## General Naming Conventions

- **Use snake_case** for table and column names
- **Use singular nouns** for table names (e.g., `account` not `accounts`)
- **Use descriptive but concise names** that clearly indicate the column's purpose

## Index Naming Standard

The standard naming convention for indexes in PostgreSQL follows this pattern:

```
{tablename}_{columnname(s)}_{suffix}
```

### Index Suffixes

| Suffix  | Constraint/Index Type   |
| ------- | ----------------------- |
| `pkey`  | Primary Key constraint  |
| `key`   | Unique constraint       |
| `excl`  | Exclusion constraint    |
| `idx`   | Any other kind of index |
| `fkey`  | Foreign key             |
| `check` | Check constraint        |

### Examples

```sql
-- Primary key
account_pkey

-- Unique constraint
user_email_key

-- Foreign key
order_customer_id_fkey

-- Regular index
product_name_idx

-- Check constraint
user_age_check

-- Multi-column index
order_customer_id_created_at_idx
```

### Additional Best Practices

- Include `created_at` and `updated_at` timestamp columns where necessary
- Use proper constraints (NOT NULL, CHECK, UNIQUE) where appropriate
- Create indexes on frequently queried columns
- Keep names under 63 characters (PostgreSQL limit)
- Avoid reserved keywords
- Use consistent abbreviations across your schema
- Consider prefixing related tables with a common identifier for organization

## Golang Standards

- Always use camelCase for JSON struct tags
- Always use `any` instead of `interface{}`
- Use singular package names (e.g., `user` instead of `users`)
- Use pgxscan for scanning rows into structs
- db tags are not required for struct fields, they are automatically inferred from struct field names

## Frontend Standarts

- Prefer TypeScript `Types` over `Interfaces`
- Use **Zod** for schema validation
- Check `package.json` for available packages
- Use **shadcn/ui** for core components
- Create new components under `components` directory. `components/ui` directory is only reserved for components from **shadcn/ui**

### Core Components

Here are the components available in **shadcn/ui**. More docs at `https://ui.shadcn.com/docs/components`.

First check `components/ui` directory to see if a component exists before adding it. Do not overwrite.

- npx shadcn@latest add accordion
- npx shadcn@latest add alert
- npx shadcn@latest add alert-dialog
- npx shadcn@latest add aspect-ratio
- npx shadcn@latest add avatar
- npx shadcn@latest add badge
- npx shadcn@latest add breadcrumb
- npx shadcn@latest add button
- npx shadcn@latest add calendar
- npx shadcn@latest add card
- npx shadcn@latest add carousel
- npx shadcn@latest add chart
- npx shadcn@latest add checkbox
- npx shadcn@latest add collapsible
- npx shadcn@latest add combobox
- npx shadcn@latest add command
- npx shadcn@latest add context-menu
- npx shadcn@latest add date-picker
- npx shadcn@latest add dialog
- npx shadcn@latest add drawer
- npx shadcn@latest add dropdown-menu
- npx shadcn@latest add form
- npx shadcn@latest add hover-card
- npx shadcn@latest add input
- npx shadcn@latest add input-otp
- npx shadcn@latest add label
- npx shadcn@latest add menubar
- npx shadcn@latest add navigation-menu
- npx shadcn@latest add pagination
- npx shadcn@latest add popover
- npx shadcn@latest add progress
- npx shadcn@latest add radio-group
- npx shadcn@latest add resizable
- npx shadcn@latest add scroll-area
- npx shadcn@latest add select
- npx shadcn@latest add separator
- npx shadcn@latest add sheet
- npx shadcn@latest add sidebar
- npx shadcn@latest add skeleton
- npx shadcn@latest add slider
- npx shadcn@latest add sonner
- npx shadcn@latest add switch
- npx shadcn@latest add table
- npx shadcn@latest add tabs
- npx shadcn@latest add textarea
- npx shadcn@latest add toast
- npx shadcn@latest add toggle
- npx shadcn@latest add toggle-group
- npx shadcn@latest add tooltip
