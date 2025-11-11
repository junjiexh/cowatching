# Cowatching

A real-time co-watching platform that allows users to watch videos together in synchronized rooms.

## Architecture

This is a full-stack application with:

- **Frontend**: React with TypeScript (located in `/frontend`)
- **Backend**: Go with Chi router (located in `/backend`)
- **Database**: PostgreSQL 15
- **Deployment**: Docker & Docker Compose

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Router**: Chi v5
- **Database**: PostgreSQL 15
- **ORM**: sqlc (type-safe SQL)
- **Database Driver**: pgx/v5

### Frontend
- **Framework**: React with TypeScript
- **Build Tool**: Vite

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)
- sqlc (for database code generation)

### Running with Docker Compose

1. Clone the repository:
```bash
git clone <repository-url>
cd cowatching
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Start all services:
```bash
docker-compose up -d
```

This will start:
- **PostgreSQL** on `localhost:5432`
- **Backend API** on `localhost:8080`
- **Frontend** on `localhost:3000`

4. Check services are running:
```bash
# Check backend health
curl http://localhost:8080/health

# Check API status
curl http://localhost:8080/api/v1/status
```

### Running Locally (Development)

#### Backend

1. Install dependencies:
```bash
cd backend
go mod download
```

2. Install sqlc for code generation:
```bash
# On macOS
brew install sqlc

# On Linux
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

3. Generate database code:
```bash
sqlc generate
```

4. Start PostgreSQL:
```bash
docker-compose up -d postgres
```

5. Run the backend:
```bash
go run cmd/api/main.go
```

#### Frontend

1. Install dependencies:
```bash
cd frontend
npm install
```

2. Start the development server:
```bash
npm run dev
```

## Project Structure

```
cowatching/
├── backend/                 # Go backend service
│   ├── cmd/
│   │   └── api/            # Application entrypoint
│   ├── internal/
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection
│   │   ├── handlers/       # HTTP handlers
│   │   └── models/         # Business models
│   ├── db/
│   │   ├── migrations/     # SQL migrations
│   │   └── queries/        # SQL queries for sqlc
│   ├── sqlc.yaml           # sqlc configuration
│   ├── Dockerfile
│   ├── Makefile
│   └── README.md
├── frontend/               # React frontend
│   ├── app/
│   ├── public/
│   ├── Dockerfile
│   └── README.md
├── docker-compose.yml      # Docker Compose configuration
├── .env.example           # Environment variables template
└── README.md              # This file
```

## API Endpoints

### Health & Status

- `GET /health` - Health check endpoint
- `GET /` - API root
- `GET /api/v1/status` - API status

### Database Schema

The application includes the following main tables:

- **users** - User accounts
- **rooms** - Watch rooms
- **room_participants** - Room membership
- **videos** - Videos being watched in rooms

See `backend/db/migrations/001_init_schema.sql` for the complete schema.

## Development

### Backend Development

```bash
cd backend

# Run with hot reload (requires air)
make dev

# Run tests
make test

# Format code
make fmt

# Build binary
make build

# Generate sqlc code
make sqlc-generate
```

### Frontend Development

```bash
cd frontend

# Start dev server
npm run dev

# Build for production
npm run build

# Run tests
npm test
```

## Docker Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Rebuild and start
docker-compose up -d --build

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v
```

## Environment Variables

Key environment variables (see `.env.example`):

```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cowatching

# Backend
SERVER_PORT=8080
```

## Database Migrations

To run migrations manually:

```bash
# Using psql
psql -U postgres -d cowatching -f backend/db/migrations/001_init_schema.sql

# Or use the Makefile
cd backend
make migrate-up DATABASE_URL=postgresql://postgres:postgres@localhost:5432/cowatching
```

## Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

## Contributing

1. Create a feature branch from `main`
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## Roadmap

- [ ] User authentication and authorization
- [ ] Real-time video synchronization with WebSockets
- [ ] Room chat functionality
- [ ] Video upload and streaming
- [ ] User presence indicators
- [ ] Playback controls sync
- [ ] Room invitations and sharing

## License

MIT

## Support

For issues and questions, please open an issue on GitHub.
