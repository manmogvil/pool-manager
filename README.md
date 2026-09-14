# Lottery Pool Manager

A personal web application designed to manage a group of people who participate together in weekly lottery draws. Built as a learning project to practice **Go backend development**, **PostgreSQL**, and **REST API design**.

---

## Overview

This application simplifies the management of lottery pools by tracking participants, contributions, tickets, draws, and prizes. The project focuses on building a clean, maintainable backend architecture while learning Go best practices.

---

## Tech Stack

**Backend**
- Go
- Chi (HTTP router)
- PostgreSQL
- pgx (database driver)
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
│   │   ├── check_ticket.go     # Ticket verification handler
│   │   ├── loteria_api.go      # Loteria API handler
│   │   └── tests/              # Handler unit tests
│   ├── services/               # External API clients
│   │   └── loteria_api.go      # Loteria API client
│   ├── store/                  # Database operations
│   │   └── integration_test.go # Store integration tests
│   ├── models/                 # Data structures
│   ├── utils/                  # Shared utilities
│   │   └── helpers_test.go     # Utils unit tests
│   └── migrations/             # SQL schema files
├── frontend/                   # Vue 3 application
└── docker-compose.yml          # PostgreSQL setup
```

---

## Features

### Participants
- Register new participants
- List all participants
- Update participant information
- Soft delete (deactivate) participants

### Contributions
- Track monthly payments per game
- View payment history by participant
- Filter contributions by period
- Support for multiple payment methods (CASH, BIZUM)

### Lottery Games
- Define different lottery types
- Configure draw days and ticket prices
- Enable/disable games

### Draws
- Create draws for specific dates
- Record official results
- Fetch results from Loteria API
- View draws by game

### Tickets
- Register tickets with played numbers
- Check prizes against Loteria API
- Track prizes and matched numbers
- View tickets by draw

### Loteria API Integration
- Fetch results by date range (Monday batch)
- Check individual tickets against API
- Dual API key support with automatic fallback
- Rate limit handling (429/403)

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |
| GET | /participants | List participants |
| POST | /participants | Create participant |
| PUT | /participants/{id} | Update participant |
| PUT | /participants/{id}/activate | Activate participant |
| DELETE | /participants/{id} | Deactivate participant |
| GET | /contributions | List contributions |
| POST | /contributions | Create contribution |
| GET | /contributions/period | Get contributions by period |
| GET | /contributions/participant/{id} | Get contributions by participant |
| GET | /games | List lottery games |
| POST | /games | Create lottery game |
| PUT | /games/{id} | Update lottery game |
| DELETE | /games/{id} | Delete lottery game |
| GET | /draws | List draws |
| POST | /draws | Create draw |
| PUT | /draws/{id}/results | Update draw results |
| PUT | /draws/{id}/process | Mark draw as processed |
| DELETE | /draws/{id} | Delete draw |
| GET | /draws/pending | Get pending draws |
| GET | /draws/game/{id} | Get draws by game |
| GET | /tickets | List tickets |
| POST | /tickets | Create ticket |
| PUT | /tickets/{id}/prize | Update ticket prize |
| DELETE | /tickets/{id} | Delete ticket |
| GET | /tickets/draw/{id} | Get tickets by draw |
| POST | /loteria-api/fetch-results | Fetch results from API |
| POST | /check-ticket | Check ticket against API |

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
LOTERIA_API_KEY=your_api_key_1,your_api_key_2

# Server
BASE_URL=http://localhost:8080
PORT=8080

# CORS
CORS_ORIGIN=http://localhost:5173
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

**Run a specific test:**
```bash
go test ./handlers/tests/... -run TestTicketHandler_Create_Success -v
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
│   ├── mock_*_test.go              # Mocks per module (6 files)
│   └── *_test.go                   # Handler tests (56 tests)
└── store/
    ├── setup_test.go               # Test DB setup + cleanup
    ├── *_store_test.go             # Store tests (45 tests)
    ├── participants_store_test.go
    ├── games_store_test.go
    ├── draws_store_test.go
    ├── tickets_store_test.go
    └── contributions_store_test.go
```

### How it works

- **Handler tests** use mocked store interfaces — no database needed
- **Store tests** create a temporary `lottery_pool_test` database, run migrations, execute tests, and drop it automatically
- **Utils tests** are pure functions with no dependencies

### Running

```bash
go test ./...                  # All tests
go test ./store/... -v         # Store tests only
go test ./handlers/tests/...   # Handler tests only
go test ./... -run TestTicket  # Match by name
```

---

## Learning Goals

This project was built to practice:

- **Go fundamentals**: structs, interfaces, error handling, packages
- **REST API design**: proper HTTP methods, status codes, JSON responses
- **Database design**: relationships, constraints, migrations
- **Code organization**: clean architecture, separation of concerns
- **Input validation**: request validation, error handling
- **Testing**: unit tests (mocked store), integration tests (real PostgreSQL), Go `testing` package
- **CI/CD**: GitHub Actions pipeline for automated test execution
- **Docker**: containerization, database setup
- **External APIs**: integrating with third-party services
- **Environment variables**: secure configuration management

---

## Project Structure

- **models/**: Data structures representing database entities
- **store/**: Database operations (PostgreSQL) + integration tests
- **handlers/**: HTTP request handlers
  - **tests/**: Handler unit tests with mocked store
- **services/**: External API clients (Loteria API)
- **routes/**: API route definitions
- **utils/**: Shared utility functions
- **migrations/**: SQL schema files (used by app and tests)
- **frontend/**: Vue 3 application (Vuetify)
