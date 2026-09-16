# Lottery Pool Manager

A personal web application designed to manage a group of people who participate together in weekly lottery draws. Built as a learning project to practice **Go backend development**, **PostgreSQL**, and **REST API design**.

---

## Overview

This application simplifies the management of lottery pools by tracking users, contributions, tickets, draws, and prizes. The project focuses on building a clean, maintainable backend architecture while learning Go best practices.

---

## Tech Stack

**Backend**
- Go
- Chi (HTTP router)
- PostgreSQL
- pgx (database driver)
- JWT authentication (golang-jwt)
- bcrypt password hashing
- Unit & integration tests (Go `testing` package)

**Frontend**
- Vue 3
- Vuetify
- Vite

**Infrastructure**
- Docker
- Docker Compose
- GitHub Actions (CI/CD)

---

## Architecture

```
lottery-pool-manager/
├── backend/
│   ├── main.go                 # Server entry point
│   ├── .env                    # Environment variables (not in git)
│   ├── .env.example            # Environment template
│   ├── routes/                 # Route definitions
│   ├── handlers/               # HTTP request handlers
│   │   ├── auth.go             # Auth handler (register, login, users)
│   │   ├── check_ticket.go     # Ticket verification handler
│   │   ├── loteria_api.go      # Loteria API handler
│   │   └── tests/              # Handler unit tests
│   ├── middleware/             # Auth & admin middleware
│   ├── services/               # Business logic & external API clients
│   │   ├── auth.go             # Password hashing, JWT generation
│   │   └── loteria_api.go      # Loteria API client
│   ├── store/                  # Database operations
│   │   └── *_store_test.go     # Store integration tests
│   ├── models/                 # Data structures
│   ├── utils/                  # Shared utilities
│   │   └── helpers_test.go     # Utils unit tests
│   └── migrations/             # SQL schema files (001-006)
├── frontend/                   # Vue 3 application
└── docker-compose.yml          # PostgreSQL setup
```

---

## Features

### Users & Authentication
- User registration and login (JWT)
- Role-based access control (admin / user)
- Activate/deactivate users (admin only)
- Edit user profile (admin only)
- Users created inactive by default (admin must activate)

### Contributions
- Track monthly payments per game
- View payment history by user
- Filter contributions by period
- Support for multiple payment methods (CASH, BIZUM)

### Lottery Games
- Define different lottery types
- Configure draw days and ticket prices
- Enable/disable games (admin only)

### Draws
- Create draws for specific dates
- Record official results
- Fetch results from Loteria API (admin only)
- View draws by game

### Tickets
- Register tickets with played numbers
- Check prizes against Loteria API
- Track prizes and matched numbers
- Cost auto-calculated from game ticket price

### Loteria API Integration
- Fetch results by date range (Monday batch)
- Check individual tickets against API
- Dual API key support with automatic fallback
- Rate limit handling (429/403)

---

## API Endpoints

### Auth
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /auth/register | No | Register new user (inactive) |
| POST | /auth/login | No | Login, returns JWT |
| GET | /auth/me | Yes | Get current user |
| GET | /auth/users | Admin | List all users |
| PUT | /auth/users/{id} | Admin | Update user |
| PUT | /auth/users/{id}/activate | Admin | Activate user |
| PUT | /auth/users/{id}/deactivate | Admin | Deactivate user |

### Contributions
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /contributions | No | List contributions |
| GET | /contributions/{id} | No | Get contribution by ID |
| GET | /contributions/user/{id} | No | Get contributions by user |
| GET | /contributions/period | No | Get contributions by period |
| POST | /contributions | Yes | Create contribution |
| PUT | /contributions/{id} | Yes | Update contribution |
| DELETE | /contributions/{id} | Yes | Delete contribution |

### Games
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /games | No | List lottery games |
| GET | /games/{id} | No | Get game by ID |
| POST | /games | Admin | Create game |
| PUT | /games/{id} | Admin | Update game |
| DELETE | /games/{id} | Admin | Delete game |

### Draws
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /draws | No | List draws |
| GET | /draws/{id} | No | Get draw by ID |
| GET | /draws/game/{id} | No | Get draws by game |
| GET | /draws/pending | No | Get pending draws |
| POST | /draws | Yes | Create draw |
| PUT | /draws/{id}/results | Yes | Update draw results |
| PUT | /draws/{id}/process | Yes | Mark draw as processed |
| DELETE | /draws/{id} | Yes | Delete draw |

### Tickets
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /tickets | No | List tickets |
| GET | /tickets/{id} | No | Get ticket by ID |
| GET | /tickets/draw/{id} | No | Get tickets by draw |
| POST | /tickets | Yes | Create ticket |
| PUT | /tickets/{id}/prize | Yes | Update ticket prize |
| DELETE | /tickets/{id} | Yes | Delete ticket |

### External APIs
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /loteria-api/fetch-results | Yes | Fetch results from API |
| POST | /check-ticket | Yes | Check ticket against API |

---

## Environment Variables

Create a `.env` file in `backend/` based on `.env.example`:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=lottery
DB_PASSWORD=lottery123
DB_NAME=lottery_pool

# Loteria API (comma-separated for fallback)
LOTERIA_API_KEY=your_api_key

# Server
BASE_URL=http://localhost:8080
PORT=8080

# CORS
CORS_ORIGIN=http://localhost:5173

# Auth
JWT_SECRET=your-secret-key-change-in-production
```

---

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Node.js 18+

### 1. Start PostgreSQL

```bash
docker compose up -d
```

Verify it's running:
```bash
docker ps --filter "name=lottery-postgres"
```

### 2. Run the Backend

```bash
cd backend
go run main.go
```

The API starts at `http://localhost:8080`.

### 3. Run the Frontend

```bash
cd frontend
npm install
npm run dev
```

The app starts at `http://localhost:5173`.

### 4. Run Tests

**All tests:**
```bash
cd backend
go test ./...
```

**With verbose output:**
```bash
go test ./... -v
```

**Specific packages:**
```bash
go test ./utils/...              # Utils tests
go test ./handlers/tests/...     # Handler tests (mocked store)
go test ./store/...              # Store tests (requires PostgreSQL)
```

**Force re-run (no cache):**
```bash
go test ./... -count=1
```

**Run tests matching a pattern:**
```bash
go test ./handlers/tests/... -run TestTicketHandler -v
```

> **Note:** Store tests (`./store/...`) require the PostgreSQL container running. Handler tests use mocks and don't need a database.

---

## Testing

### Structure

```
backend/
├── utils/
│   └── helpers_test.go              # Utils tests (6 tests)
├── handlers/tests/
│   ├── helpers_test.go              # Shared test helpers
│   ├── mock_*_test.go              # Mocks per module
│   └── *_test.go                   # Handler tests
│       ├── auth_test.go            # 18 auth tests
│       ├── check_ticket_test.go    # 6 check ticket tests
│       ├── loteria_api_test.go     # 3 loteria API tests
│       ├── contributions_test.go   # Contribution tests
│       ├── games_test.go           # Game tests
│       ├── draws_test.go           # Draw tests
│       └── tickets_test.go         # Ticket tests
└── store/
    ├── setup_test.go               # Test DB setup + cleanup
    └── *_store_test.go             # Store integration tests
```

### How it works

- **Handler tests** use mocked store interfaces — no database needed
- **Store tests** create a temporary `lottery_pool_test` database, run migrations, execute tests, and drop it automatically
- **Utils tests** are pure functions with no dependencies

### Seed Data

Migrations include seed data (`006_seed_data.sql`):
- 5 users (1 admin active, 3 regular active, 1 inactive)
- All passwords: `password123`
- 3 lottery games, 6 draws, 6 contributions, 6 tickets

---

## Migrations

| File | Description |
|------|-------------|
| `001_create_users.sql` | Users table with auth fields |
| `002_create_lottery_games.sql` | Lottery games |
| `003_create_contributions.sql` | Contributions (user_id FK) |
| `004_create_draws.sql` | Draws with draw_id_api |
| `005_create_tickets.sql` | Tickets |
| `006_seed_data.sql` | Seed data for development |

---

## Deployment

### Architecture

| Service | Technology | Purpose |
|---------|-----------|---------|
| Frontend | Vercel | Vue.js SPA hosting |
| Backend | Render | Go API server |
| Database | Supabase | PostgreSQL (managed) |

### Commands

| Command | Context | Description |
|---------|---------|-------------|
| `go run main.go` | Local dev | Start server only |
| `go run ./cmd/migrate` | Local setup | Run migrations only |
| `./server` | Production | Start server only |
| `./migrate` | Production | Run migrations only |
| `./start.sh` | Docker | Run migrations + server |

### Production Environment Variables

| Variable | Source | Description |
|----------|--------|-------------|
| `DB_HOST` | Supabase | `db.xxx.supabase.co` |
| `DB_PORT` | Supabase | `5432` |
| `DB_USER` | Supabase | `postgres` |
| `DB_PASSWORD` | Supabase | Project password |
| `DB_NAME` | Supabase | `postgres` |
| `JWT_SECRET` | Generate | `openssl rand -hex 32` |
| `CORS_ORIGIN` | Vercel URL | `https://app.vercel.app` |
| `LOTERIA_API_KEY` | Your keys | Comma-separated |
| `PORT` | Render | `8080` |

### Deployment Steps

1. **Supabase** — Create project, copy connection string
2. **Render** — Create Web Service, set Docker, add env vars
3. **Vercel** — Import repo, set `VITE_API_URL`

---

## Learning Goals

This project was built to practice:

- **Go fundamentals**: structs, interfaces, error handling, packages
- **REST API design**: proper HTTP methods, status codes, JSON responses
- **Database design**: relationships, constraints, migrations
- **Code organization**: clean architecture, separation of concerns
- **Input validation**: request validation, error handling
- **Authentication**: JWT tokens, bcrypt hashing, role-based access
- **Testing**: unit tests (mocked store), integration tests (real PostgreSQL)
- **CI/CD**: GitHub Actions pipeline for automated test execution
- **Docker**: containerization, database setup
- **External APIs**: integrating with third-party services
- **Environment variables**: secure configuration management
