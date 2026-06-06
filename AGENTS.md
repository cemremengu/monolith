# General Guidelines

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:

- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:

- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:

- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:

```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

# Technology Stack

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

- `task server` - Start backend with hot reload using wgo
- `task test:go` - Run Go tests
- `task dev:lint:go` - Run lint
- `task dev:fmt:go` - Format Go code
- `task dev:update` - Update Go dependencies

## Frontend Development

- `task web` - Start frontend dev server with hot reload
- `cd web && npm run build` - Build frontend for production
- `task dev:lint:web` - Run Oxlint on frontend code
- `task dev:fmt:web` - Format frontend code
- `task test:web` - Run frontend tests
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

- **Identity columns**: Always use `GENERATED ALWAYS AS IDENTITY` instead of `SERIAL` or `BIGSERIAL`. For UUIDs, use `UUID` type with `gen_random_uuid()` as the default value
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
- Keep code comments minimum and relevant to the code itself. Do not add comments that are not directly related to the code or that explain obvious things.

## Golang Standards

- Always use camelCase for JSON struct tags
- Always use `any` instead of `interface{}`
- Use singular package names (e.g., `user` instead of `users`)
- Use pgxscan for scanning rows into structs
- db tags are not required for struct fields, they are automatically inferred from struct field names
- Prefer table-driven tests for unit tests when applicable
- For JSONB columns, use the actual Go types directly in row structs instead of `[]byte`. pgx v5 handles JSONB serialization/deserialization automatically.
- Use restrictive permissions (0700 for dirs, 0600 for files) instead of world-readable (0755/0644) unless the content is genuinely public

## Frontend Standarts

- Prefer TypeScript `Types` over `Interfaces`
- Use **Zod** for schema validation
- Check `package.json` for available packages
- Use **shadcn/ui** for core components
- Create new components under `components` directory. `components/ui` directory is only reserved for components from **shadcn/ui**
- Avoid barrel files (i.e., `index.ts` files that re-export other modules)
- Prefer named exports over default exports
- When width and height are the same (e.g., `w-5 h-5`), use the shorthand `size-5` instead for consistency. Do not apply this to components under `components/ui` as they follow shadcn/ui conventions.
- For forms, the convention is to use react-hook-form + Controller with Field/FieldLabel/FieldError from components/ui/field

### Core UI Components

Refer to **shadcn/ui** for available components and documentation: `https://ui.shadcn.com/docs/components`

Before adding a new component, always check the `components/ui` directory to verify whether it already exists.
Do not overwrite existing components.
Never modify components under `components/ui` directly unless explicitly asked to do so.
