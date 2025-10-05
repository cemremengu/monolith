# Instructions

You are an agent - please keep going until the user’s query is completely resolved, before ending your turn and yielding back to the user. Only terminate your turn when you are sure that the problem is solved.

If you are not sure about file content or codebase structure pertaining to the user’s request, use your tools to read files and gather the relevant information: do NOT guess or make up an answer.

You MUST plan extensively before each function call, and reflect extensively on the outcomes of the previous function calls. DO NOT do this entire process by making function calls only, as this can impair your ability to solve the problem and think insightfully.

If you are unsure about the user’s intent, ask clarifying questions to gather more context. Do not make assumptions about what the user wants.

# Workflow

## High-Level Problem Solving Strategy

1. Understand the problem deeply. Carefully read the issue and think critically about what is required.
2. Investigate the codebase. Explore relevant files, search for key functions, and gather context.
3. Develop a clear, step-by-step plan. Break down the fix into manageable, incremental steps.
4. Implement the fix incrementally. Make small, testable code changes.
5. Debug as needed. Use debugging techniques to isolate and resolve issues.
6. Test frequently. Run tests after each change to verify correctness.
7. Iterate until the root cause is fixed and all tests pass.
8. Reflect and validate comprehensively. After tests pass, think about the original intent, write additional tests to ensure correctness, and remember there are hidden tests that must also pass before the solution is truly complete.

Refer to the detailed sections below for more information on each step.

## 1. Deeply Understand the Problem

Carefully read the issue and think hard about a plan to solve it before coding.

## 2. Codebase Investigation

- Explore relevant files and directories.
- Search for key functions, classes, or variables related to the issue.
- Read and understand relevant code snippets.
- Identify the root cause of the problem.
- Validate and update your understanding continuously as you gather more context.

## 3. Develop a Detailed Plan

- Outline a specific, simple, and verifiable sequence of steps to fix the problem.
- Break down the fix into small, incremental changes.

## 4. Making Code Changes

- Before editing, always read the relevant file contents or section to ensure complete context.
- If a patch is not applied correctly, attempt to reapply it.
- Make small, testable, incremental changes that logically follow from your investigation and plan.
- If you need to use external libraries from GitHub, make sure they have recent commits and compatible with the codebase. If not sure, let the user decide.

## 5. Debugging

- Make code changes only if you have high confidence they can solve the problem
- When debugging, try to determine the root cause rather than addressing symptoms
- Debug for as long as needed to identify the root cause and identify a fix
- Use print statements, logs, or temporary code to inspect program state, including descriptive statements or error messages to understand what's happening
- To test hypotheses, you can also add test statements or functions
- Revisit your assumptions if unexpected behavior occurs.

## 6. Testing

- After each change, verify correctness by running relevant tests.
- If tests fail, analyze failures and revise your patch.
- Write additional tests if needed to capture important behaviors or edge cases.
- Ensure all tests pass before finalizing.

## 7. Final Verification

- Confirm the root cause is fixed.
- Review your solution for logic correctness and robustness.
- Iterate until you are extremely confident the fix is complete and all tests pass.

## 8. Final Reflection and Additional Testing

- Reflect carefully on the original intent of the user and the problem statement.
- Think about potential edge cases or scenarios that may not be covered by existing tests.
- Write additional tests that would need to pass to fully validate the correctness of your solution.
- Run these new tests and ensure they all pass.
- Be aware that there are additional hidden tests that must also pass for the solution to be successful.
- Do not assume the task is complete just because the visible tests pass; continue refining until you are confident the fix is robust and comprehensive.

# Prerequisites

Before working with this repository, ensure you have:

- Go 1.25.1 or higher
- Node.js 18 or higher
- PostgreSQL (or use Docker Compose: `docker-compose up -d`)
- Docker (for sandbox functionality)
- Task (go-task) - for running build commands

# Technology Stack

This repository is a monolithic application built with:

## Backend

- **Language**: Go 1.25.1+
- **Framework**: Echo v4
- **Database**: PostgreSQL with pgx driver
- **ORM/Query**: scany for result scanning
- **Migrations**: Goose (auto-applied on server start)
- **Testing**: Go testing with stretchr/testify
- **Linting**: golangci-lint

## Frontend

- **Language**: TypeScript
- **Framework**: React 18+
- **Build Tool**: Vite
- **Router**: TanStack Router (with code generation)
- **State Management**: Zustand
- **UI Components**: shadcn/ui (Radix UI + Tailwind)
- **Styling**: Tailwind CSS v4
- **Validation**: Zod
- **HTTP Client**: TanStack Query

# Project Structure

- `/cmd/monolith` - Main application entry point
- `/internal` - Private application code
  - `/api` - API handlers
  - `/config` - Configuration loading
  - `/service` - Business logic services
  - `/util` - Utility packages
- `/web` - React frontend application
  - `/src` - Frontend source code
  - `/public` - Static assets
- `/migrations` - Database migration files
- `/taskfiles` - Task definitions for build, test, development

# Build and Test Commands

Use the Taskfile for all common operations:

## Running

```bash
task server        # Run backend with hot reload (wgo)
task web          # Run frontend dev server
task run          # Run both server and web
```

## Testing

```bash
task test:go      # Run Go tests
task test:web     # Run frontend tests
task test:all     # Run all tests
task test:coverage # Run Go tests with coverage report
```

## Linting

```bash
task dev:lint:go    # Run golangci-lint (fmt + run)
task dev:lint:web   # Run ESLint and TypeScript type checking
```

## Building

```bash
task build:server  # Build backend binary
task build:web    # Build frontend for production
```

## Other

```bash
task dev:router:generate  # Generate TanStack Router routes
task dev:update          # Update Go dependencies
```

# Code Style

## Golang

- Always use `any` instead of `interface{}`
- Use stretchr/testify for tests
- Use camelCase JSON struct tags
- Follow the enabled golangci-lint rules
- Use context for cancellation and timeouts
- Keep packages focused and small
- Use meaningful variable names (avoid single-letter except for short loops)

## TypeScript/React

- Use TypeScript strict mode
- Follow ESLint rules configured in the project
- Use functional components with hooks
- Prefer named exports over default exports
- Use TanStack Query for server state
- Use Zustand for client state
- Follow shadcn/ui patterns for components
- Use Tailwind utility classes (avoid custom CSS when possible)
- Use Zod for runtime validation

# Common Patterns

## Backend Patterns

- **Services**: Business logic in `/internal/service/*` packages
- **Configuration**: Environment variables loaded via godotenv
- **Database**: Use pgx with scany for type-safe queries
- **Validation**: Use go-playground/validator for struct validation
- **API**: Echo handlers in service packages, group routes logically

## Frontend Patterns

- **Routing**: File-based routing with TanStack Router (regenerate after route changes)
- **API Calls**: Use TanStack Query hooks for server state
- **Forms**: React Hook Form + Zod validation + shadcn form components
- **State**: Zustand stores for client-side state
- **Styling**: Tailwind utility classes with shadcn/ui components
- **Internationalization**: i18next for translations (en_US and tr_TR supported)

# Important Notes

- **Migrations**: Database migrations are automatically applied on server startup (no manual migration command needed)
- **Hot Reload**: The backend uses `wgo` for hot reloading during development
- **Router Generation**: After adding/modifying routes in `/web/src`, run `task dev:router:generate` to update TanStack Router
- **Sandbox Security**: Code execution in sandbox is isolated using Docker containers with resource limits
- **Environment Variables**: Copy `.env.example` to `.env` and configure as needed
- **Agent Mode**: The application supports agent mode for AI-powered interactions with configurable tools and iterations
