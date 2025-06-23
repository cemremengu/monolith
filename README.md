# Monolith App Template

A simple template application built with:

## Backend

- Go with Echo framework
- PostgreSQL database
- pgx for database connection
- scany for query result scanning
- Goose for database migrations

## Frontend

- React with TypeScript
- Vite as build tool
- shadcn/ui components
- TanStack Router for routing
- Tailwind CSS for styling
- Zod for schema validation
- Zustand for state management

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js 18+
- PostgreSQL (or use Docker Compose)

### Database Setup

1. Start PostgreSQL with Docker Compose:

```bash
docker-compose up -d
```

2. Migrations:

Migations will applied automatically on server start.

### Backend Setup

1. Install dependencies:

```bash
go mod tidy
```

2. Run the server:

```bash
go run cmd/main.go
```

### Frontend Setup

1. Navigate to the web directory:

```bash
cd web
```

2. Install dependencies:

```bash
npm install
```

3. Build the frontend:

```bash
npm run build
```

### Running the Application

1. Make sure PostgreSQL is running and migrations are applied
2. Start the Go server: `go run cmd/main.go`
3. Open your browser to `http://localhost:3001`
4. Initial credentials for the admin user are:
   - Username: `admin`
   - Password: `admin123`

## Development

For development you can use the provided `Taskfile` to run tasks easily. Besure to have the following dependecies installed:

- [Task](https://taskfile.dev/installation/)
- [wgo](https://github.com/bokwoon95/wgo)
- [golangci-lint](https://golangci-lint.run/welcome/install/) (optional, for linting)

For development, you can run the frontend and backend separately:

1. Backend: `task run:server`
2. Frontend: `task run:web`

This allows hot reloading for the frontend while developing.
