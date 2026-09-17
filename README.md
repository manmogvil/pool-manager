# Lottery Pool Manager

A full-stack web application born from a simple weekly tradition — a group of friends pooling money for lottery draws (once a week, a football pool or La Primitiva). One of our friends worked at a lottery administration, which made the whole process more natural. Combined with the goal of learning **Go** and understanding the **full deployment lifecycle** of a web application, this project came to life. A practical way to digitalize what we used to do on paper.

---

## Tech Stack

| Layer | Technologies |
|-------|-------------|
| **Backend** | Go, Chi, PostgreSQL, pgx, JWT, bcrypt |
| **Frontend** | Vue 3, Vuetify 3, Vite |
| **Infrastructure** | Docker, GitHub Actions, Vercel, Render, Supabase |

---

## Key Features

- **Authentication & Roles** — JWT-based auth with admin/user roles. New users start inactive until activated by admin.
- **Dashboard** — Stats overview with total collected, spent, prizes won, pending payments, contributions by game, and monthly trend chart.
- **Responsive Design** — Mobile-first layout with collapsable sidebar, horizontal scroll tables, and adaptive dialogs.
- **Lotería API Integration** — Fetch draw results and check ticket prizes against the official Spanish lottery API, with dual API key fallback.
- **Full CRUD** — Users, contributions, games, draws, and tickets with role-based permissions.

---

## Architecture

```
├── backend/
│   ├── handlers/        # HTTP handlers + tests
│   ├── store/           # PostgreSQL queries + integration tests
│   ├── models/          # Data structures
│   ├── services/        # Business logic & external APIs
│   ├── middleware/       # JWT auth & admin guards
│   ├── routes/          # Route definitions
│   ├── migrations/      # SQL schema (001-005)
│   └── cmd/migrate/     # Standalone migration binary
└── frontend/
    └── src/
        ├── views/           # Page components
        ├── components/      # Reusable UI components
        ├── composables/     # useAuth()
        ├── services/        # API client
        └── router/          # Vue Router + guards
```

---

## Getting Started

**Prerequisites:** Go 1.27+, Docker, Node.js 18+

```bash
# 1. Start PostgreSQL
docker compose up -d

# 2. Start backend
cd backend && go run main.go

# 3. Start frontend
cd frontend && npm install && npm run dev
```

App runs at `http://localhost:5173` — API at `http://localhost:8080`

### Commands

| Command | Context | Description |
|---------|---------|-------------|
| `docker compose up -d` | Local | Start PostgreSQL |
| `go run main.go` | Local dev | Start API server |
| `go run ./cmd/migrate` | Local setup | Run database migrations |
| `./start.sh` | Production (Docker) | Run migrations + server |

---

## API Endpoints

- **Auth** — `/auth/*` — Register, login, user management (CRUD + activate/deactivate)
- **Contributions** — `/contributions/*` — CRUD + filter by user or period
- **Games** — `/games/*` — CRUD (admin only)
- **Draws** — `/draws/*` — CRUD + fetch results from Lotería API
- **Tickets** — `/tickets/*` — CRUD + check prizes against Lotería API

Full route definitions in `backend/routes/routes.go`.

---

## Testing

```bash
cd backend
go test ./...              # All tests
go test ./handlers/tests/  # Handler tests (mocked, no DB)
go test ./store/...        # Store tests (needs PostgreSQL)
```

- **Handler tests** — Mocked store interfaces, no database required
- **Store tests** — Real PostgreSQL, auto-creates and drops test DB
- **Utils tests** — Pure functions, no dependencies

---

## Deployment

| Service | Stack | Purpose |
|---------|-------|---------|
| Frontend | **Vercel** | Vue.js SPA |
| Backend | **Render** | Go API (Docker) |
| Database | **Supabase** | PostgreSQL (managed) |

---

## What I Learned

- **Go**: Structs, interfaces, error handling, packages, Chi router
- **PostgreSQL**: Relationships, constraints, migrations, aggregation queries
- **REST API**: Proper HTTP methods, status codes, JWT auth, role-based access
- **Testing**: Unit tests with mocks, integration tests with real DB
- **Docker**: Containerization, multi-stage builds, separate binaries
- **Vue 3**: Composition API, composables, Vuetify, responsive design
- **CI/CD**: GitHub Actions for automated testing
- **External APIs**: Third-party integration with rate limit handling
