# Technical Documentation

## Data Models

### User
| Field | Type | Description |
|-------|------|-------------|
| `id` | int | Primary key |
| `name` | string | Display name |
| `email` | string | Unique, used for login |
| `password_hash` | string | bcrypt hash (never exposed in API) |
| `role` | string | `admin` or `user` |
| `active` | bool | New users start `false`. Admin must activate. |
| `created_at` | timestamp | Account creation date |

### Contribution
| Field | Type | Description |
|-------|------|-------------|
| `id` | int | Primary key |
| `user_id` | int | FK → User |
| `game_id` | int | FK → LotteryGame |
| `month` | int | 1-12 |
| `year` | int | ≥ 2020 |
| `amount` | decimal | Money contributed |
| `paid` | bool | `true` = paid, `false` = pending |
| `payment_date` | date | When payment was made |
| `payment_method` | string | `CASH` or `BIZUM` (required) |
| `comments` | text | Optional notes |

### LotteryGame
| Field | Type | Description |
|-------|------|-------------|
| `id` | int | Primary key |
| `name` | string | e.g. "EuroMillones" |
| `draw_days` | string | e.g. "Tuesday, Friday" |
| `ticket_price` | decimal | Cost per ticket |
| `active` | bool | Enabled/disabled |

### Draw
| Field | Type | Description |
|-------|------|-------------|
| `id` | int | Primary key |
| `game_id` | int | FK → LotteryGame |
| `draw_date` | date | When the draw happens |
| `result_numbers` | string | Official winning numbers (nullable) |
| `result_stars` | string | Official winning stars (nullable) |
| `draw_id_api` | string | ID from Lotería API |
| `processed` | bool | `true` = tickets checked against results |

### Ticket
| Field | Type | Description |
|-------|------|-------------|
| `id` | int | Primary key |
| `draw_id` | int | FK → Draw |
| `numbers` | string | Played numbers |
| `stars` | string | Played stars (nullable) |
| `cost` | decimal | Auto-calculated from game's `ticket_price` |
| `prize_tier` | string | Prize category (nullable) |
| `prize_amount` | decimal | Money won (nullable) |
| `matched_numbers` | int | How many numbers matched (nullable) |
| `matched_stars` | int | How many stars matched (nullable) |

---

## Relationships

```
User ─────────┐
              │
              ▼
LotteryGame ──┬──▶ Contribution
              │
              └──▶ Draw ──▶ Ticket
```

- One User can have many Contributions
- One Game can have many Contributions and many Draws
- One Draw can have many Tickets

---

## Workflows

### Contribution Flow
```
1. User creates contribution (amount, game, month, year)
2. paid = true by default (can toggle to false)
3. Admin can edit/delete any contribution
4. User can only view their own
```

### Draw Flow
```
1. Admin creates draw (game + date)
2. Admin fetches results from Lotería API → result_numbers, result_stars populated
3. Admin marks draw as processed
4. Tickets can now be checked against results
```

### Ticket Flow
```
1. Admin creates ticket (draw, numbers, stars)
2. Cost auto-calculated from game's ticket_price
3. Admin clicks "Check" → API verifies against draw results
4. Prize tier, amount, and matched numbers are populated
```

---

## Role Permissions

| Action | Admin | User |
|--------|:-----:|:----:|
| Login | ✅ | ✅ |
| View dashboard | ✅ | ✅ |
| View contributions | ✅ | ✅ |
| Create contribution | ✅ | ✅ |
| Edit/delete contribution | ✅ | ❌ (own only) |
| View games | ✅ | ✅ |
| Create/edit/delete games | ✅ | ❌ |
| View draws | ✅ | ✅ |
| Create draw | ✅ | ❌ |
| Fetch draw results from API | ✅ | ❌ |
| Mark draw as processed | ✅ | ❌ |
| View tickets | ✅ | ✅ |
| Create ticket | ✅ | ❌ |
| Check ticket against API | ✅ | ❌ |
| Edit/delete ticket | ✅ | ❌ |
| View/edit users | ✅ | ❌ |
| Activate/deactivate users | ✅ | ❌ |

---

## API Authentication

```
1. POST /auth/login → receive JWT token
2. Store token in localStorage
3. All protected requests include: Authorization: Bearer <token>
4. Token expires in 24 hours
5. 401 response → redirect to login
```

---

## Lotería API Integration

- **Endpoint**: `GET /loteria-api/fetch-results` (admin only)
- **Dual key support**: If first API key hits 429/403, automatically switches to next key
- **Rate limit**: 50 requests/month, 10/minute per key
- **Check ticket**: `POST /check-ticket` verifies a ticket against draw results from the API

---

## Scheduler (Automatic Draw Result Checking)

Optional background job that automatically fetches draw results from the Lotería API.

### Configuration

```env
SCHEDULER_ENABLED=false        # Disabled by default
SCHEDULER_CRON=0 6 * * 1       # Cron expression (default: every Monday at 06:00)
```

Uses standard cron expressions (5 fields) via `robfig/cron` library.

### Cron Examples

```
0 6 * * 1       → Every Monday at 06:00
0 6 * * 1-5     → Weekdays at 06:00
@every 30m      → Every 30 minutes
0 6,18 * * 1    → Mondays at 06:00 and 18:00
```

### How it works

1. On server start, if `SCHEDULER_ENABLED=true`, scheduler starts with the configured cron expression
2. Cron triggers at the scheduled times, querying for pending draws (processed=false, date <= now)
3. For each pending draw, calls Lotería API to fetch results
4. Updates the draw with `result_numbers` and `result_stars`
5. Errors are logged but don't stop the scheduler
6. Gracefully stops on server shutdown (context cancellation)

### Key points

- **Disabled by default** — no impact on existing functionality
- **Reuses existing services** — same API client and store methods as manual endpoint
- **Optional** — application works identically without it
- **8 test cases** covering start/stop, cron schedule, errors, shutdown, overlap
