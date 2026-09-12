# Pool Manager

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

**Frontend**
- Vue 3
- PrimeVue
- Vite

**Infrastructure**
- Docker
- Docker Compose

---

## Architecture

```
lottery-pool-manager/
├── main.go                 # Server entry point
├── routes/                 # Route definitions
├── handlers/               # HTTP request handlers
├── store/                  # Database operations
├── models/                 # Data structures
├── utils/                  # Shared utilities
├── migrations/             # SQL schema files
├── frontend/               # Vue 3 application
└── docker-compose.yml      # PostgreSQL setup
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
- Track processing status
- View draws by game

### Tickets
- Register tickets with played numbers
- Track prizes and matched numbers
- View tickets by draw

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |
| GET | /participants | List participants |
| POST | /participants | Create participant |
| PUT | /participants/{id} | Update participant |
| DELETE | /participants/{id} | Deactivate participant |
| GET | /contributions | List contributions |
| POST | /contributions | Create contribution |
| GET | /games | List lottery games |
| POST | /games | Create lottery game |
| GET | /draws | List draws |
| POST | /draws | Create draw |
| PUT | /draws/{id}/results | Update draw results |
| GET | /tickets | List tickets |
| POST | /tickets | Create ticket |

---

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose

### Setup

1. Clone the repository
```bash
git clone https://github.com/username/lottery-pool-manager.git
cd lottery-pool-manager
```

2. Start PostgreSQL
```bash
docker compose up -d
```

3. Run the application
```bash
go run main.go
```

4. Access the API at `http://localhost:8080`

---

## Learning Goals

This project was built to practice:

- **Go fundamentals**: structs, interfaces, error handling, packages
- **REST API design**: proper HTTP methods, status codes, JSON responses
- **Database design**: relationships, constraints, migrations
- **Code organization**: clean architecture, separation of concerns
- **Input validation**: request validation, error handling
- **Docker**: containerization, database setup

---

## Project Structure

- **models/**: Data structures representing database entities
- **store/**: Database operations (PostgreSQL and in-memory)
- **handlers/**: HTTP request handlers
- **routes/**: API route definitions
- **utils/**: Shared utility functions
- **migrations/**: SQL schema files
- **frontend/**: Vue 3 application (in progress)

---

## Future Improvements

- [ ] Prize calculation logic
- [ ] User authentication
- [ ] Notification system
- [ ] Background workers for result processing
- [ ] Prometheus metrics
- [ ] Grafana dashboards
