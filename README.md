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

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL (or use Docker Compose)

### Database Setup

1. Start PostgreSQL with Docker Compose:

```bash
docker-compose up -d
```

2. Run migrations:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/my_db" up
```

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
3. Open your browser to `http://localhost:8080`
4. Initial credentials for the admin user are:
   - Username: `admin`
   - Password: `admin123`

The application will serve the React frontend and provide API endpoints for user management.

## API Endpoints

- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create new user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

## Development

For development, you can run the frontend and backend separately:

1. Backend: `go run cmd/main.go`
2. Frontend: `cd web && npm run dev`

This allows hot reloading for the frontend while developing.
