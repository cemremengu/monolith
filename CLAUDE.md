# Development Standards & Best Practices

## PostgreSQL Database Design

### Column Types

- **Identity columns**: Always use `GENERATED ALWAYS AS IDENTITY` instead of `SERIAL` or `BIGSERIAL`
- **Timestamps**: Always use `TIMESTAMPTZ` for datetime columns to handle timezone-aware dates
- **Text data**: Always use `TEXT` instead of `VARCHAR(N)` for variable-length strings
- **JSON data**: Always use `JSONB` instead of `JSON` for better performance and indexing capabilities

### Naming Conventions

- Use snake_case for table and column names
- Use singular nouns for table names (e.g., `account` not `accounts`)
- Use descriptive but concise names that clearly indicate the column's purpose

### Additional Best Practices

- Include `created_at` and `updated_at` timestamp columns where necessary
- Use proper constraints (NOT NULL, CHECK, UNIQUE) where appropriate
- Create indexes on frequently queried columns

## Golang Standards

- Always use camelCase for JSON struct tags
- Prefer composition over inheritance
- Always use `any` instead of `interface{}`

## Frontend Standarts

- Prefer TypeScript `Types` over `Interfaces`
- Use **Zod** for schema validation
- Check `package.json` for available packages
- Use **shadcn-ui** for core components
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
