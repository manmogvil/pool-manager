# Data Models

## Overview

The application manages lottery pools through 5 interconnected models:

```
Participant ──┐
              ├──> Contribution ──> LotteryGame
              │
LotteryGame ──┘
    │
    └──> Draw ──> Ticket
```

---

## Participant

Represents a person in the lottery pool.

| Field | Type | Description |
|-------|------|-------------|
| id | int | Primary key |
| name | string | Participant name |
| email | string | Unique email |
| active | boolean | Soft delete flag |
| created_at | timestamp | Registration date |

---

## LotteryGame

Different lottery types (EuroMillions, Primitiva, etc.).

| Field | Type | Description |
|-------|------|-------------|
| id | int | Primary key |
| name | string | Game name (unique) |
| draw_days | string | Days when draws happen |
| ticket_price | decimal | Default ticket price |
| active | boolean | Soft delete flag |

---

## Draw

A specific lottery draw/sorteo.

| Field | Type | Description |
|-------|------|-------------|
| id | int | Primary key |
| game_id | int (nullable) | FK to lottery_games |
| draw_date | date | When the draw happens |
| result_numbers | string (nullable) | Winning numbers |
| result_stars | string (nullable) | Lucky stars (EuroMillions) |
| processed | boolean | Results checked? |
| created_at | timestamp | Creation date |

---

## Ticket

A ticket purchased for a specific draw.

| Field | Type | Description |
|-------|------|-------------|
| id | int | Primary key |
| draw_id | int (nullable) | FK to draws |
| numbers | string | Played numbers |
| stars | string | Lucky stars played |
| cost | decimal | Ticket cost |
| purchased_at | timestamp | Purchase time |
| prize_tier | string (nullable) | Prize category |
| prize_amount | decimal (nullable) | Prize won |
| matched_numbers | int (nullable) | Numbers matched |
| matched_stars | int (nullable) | Stars matched |

---

## Contribution

Monthly payment from a participant for a specific game.

| Field | Type | Description |
|-------|------|-------------|
| id | int | Primary key |
| participant_id | int | FK to participants |
| game_id | int | FK to lottery_games |
| month | int | Billing month (1-12) |
| year | int | Billing year |
| amount | decimal | Payment amount |
| paid | boolean | Payment status |
| payment_date | date (nullable) | When paid |
| payment_method | string | CASH or BIZUM |
| comments | text | Extra notes |
| created_at | timestamp | Creation date |

---

## Lifecycle

### Step 1: Participants
Add the people who participate in the lottery pool. Each participant has a name and email.

### Step 2: Games
Create the lottery games the pool will play (EuroMillions, Primitiva, etc.). Define which days the draws happen and the default ticket price.

### Step 3: Contributions
Participants pay monthly. Each payment is linked to a specific game. For example, Ana pays 4€ for EuroMillions in September. You can track who has paid and who hasn't.

### Step 4: Draw Creation
When a draw is about to happen, create it in the system. For example, "EuroMillions Friday 12/09/2026". At this point, the result is unknown.

### Step 5: Ticket Purchase
Buy tickets for that draw. Enter the numbers played and the cost. You can buy multiple tickets for the same draw.

### Step 6: Results
After the draw takes place, enter the official winning numbers. The system marks the draw as processed.

### Step 7: Prize Check
Compare each ticket against the winning numbers. Update how many numbers matched and the prize won. The system can then calculate total prizes and distribute them among participants.

---

## Deletion Rules

| Parent | Child | Behavior |
|--------|-------|----------|
| LotteryGame | Draw | ON DELETE SET NULL |
| LotteryGame | Ticket | Blocks deletion if tickets exist |
| Draw | Ticket | Blocks deletion if tickets exist |
