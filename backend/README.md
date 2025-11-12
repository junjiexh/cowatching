# Cowatching Backend

A Go backend service for the Cowatching application using PostgreSQL, Chi router, and sqlc.

## Tech Stack

- **Language**: Go 1.21+
- **Router**: Chi v5
- **Database**: PostgreSQL 15
- **ORM**: sqlc (type-safe SQL)
- **Database Driver**: pgx/v5
- **Configuration**: godotenv

## Project Structure

```
backend/
├── cmd/
│   └── api/              # Application entrypoint
│       └── main.go
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection and utilities
│   ├── handlers/         # HTTP request handlers
│   └── models/           # Business models
├── db/
│   ├── migrations/       # SQL migration files
│   └── queries/          # SQL queries for sqlc
├── sqlc.yaml            # sqlc configuration
├── Dockerfile           # Docker build configuration
└── Makefile            # Common tasks
```

## Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15
- sqlc (for code generation)

### Install sqlc

```bash
# On macOS
brew install sqlc

# On Linux
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Or download binary from https://github.com/sqlc-dev/sqlc/releases
```

### Generate Database Code

After installing sqlc, generate the database access code:

```bash
cd backend
sqlc generate
```

This will create type-safe Go code in `internal/database/db/` based on your SQL queries in `db/queries/`.

### Environment Variables

Copy `.env.example` to `.env` in the root directory:

```bash
cp ../.env.example ../.env
```

Configure the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cowatching
SERVER_PORT=8080
```

## Running Locally

### With Docker Compose (Recommended)

From the project root:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Backend API on port 8080
- Frontend on port 3000

### Without Docker

1. Start PostgreSQL:
```bash
# Make sure PostgreSQL is running
psql -U postgres -c "CREATE DATABASE cowatching;"
```

2. Run migrations:
```bash
psql -U postgres -d cowatching -f db/migrations/001_init_schema.sql
```

3. Generate sqlc code:
```bash
sqlc generate
```

4. Run the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

### Health Check

```
GET /health
```

Response:
```json
{
  "status": "healthy"
}
```

### API v1

Base URL: `/api/v1`

```
GET  /api/v1/status  - API status check
```

### Planned Endpoints

- `POST   /api/v1/users` - Create user
- `GET    /api/v1/users/:id` - Get user by ID
- `POST   /api/v1/rooms` - Create room
- `GET    /api/v1/rooms/:id` - Get room by ID
- `GET    /api/v1/rooms/:code` - Get room by code
- `POST   /api/v1/rooms/:id/join` - Join room
- `POST   /api/v1/videos` - Create video
- `PATCH  /api/v1/videos/:id` - Update video playback state

## Development

### Running Tests

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

## Database

### Migrations

Database migrations are located in `db/migrations/`. To apply migrations:

```bash
psql -U postgres -d cowatching -f db/migrations/001_init_schema.sql
```

For production, consider using a migration tool like:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)

### Adding New Queries

1. Add SQL queries to files in `db/queries/`
2. Run `sqlc generate` to generate Go code
3. Use the generated code in your handlers

Example:
```sql
-- db/queries/users.sql
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
```

Then use in code:
```go
user, err := db.GetUserByEmail(ctx, email)
```

## Docker

### Build Image

```bash
docker build -t cowatching-backend .
```

### Run Container

```bash
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=yourpassword \
  cowatching-backend
```

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## License

MIT
