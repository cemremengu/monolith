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
