package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"lottery-pool-manager/models"
)

func (s *PostgreSQLStore) CreateTicket(drawID int, numbers, stars string, cost float64) (models.Ticket, error) {
	var t models.Ticket

	err := s.pool.QueryRow(context.Background(),
		`INSERT INTO tickets (draw_id, numbers, stars, cost)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars`,
		drawID, numbers, stars, cost,
	).Scan(&t.ID, &t.DrawID, &t.Numbers, &t.Stars, &t.Cost, &t.PurchasedAt,
		&t.PrizeTier, &t.PrizeAmount, &t.MatchedNumbers, &t.MatchedStars)

	if err != nil {
		return models.Ticket{}, fmt.Errorf("error creating ticket: %w", err)
	}

	return t, nil
}

func (s *PostgreSQLStore) GetAllTickets() ([]models.Ticket, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars
		 FROM tickets ORDER BY purchased_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("error querying tickets: %w", err)
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.DrawID, &t.Numbers, &t.Stars, &t.Cost, &t.PurchasedAt,
			&t.PrizeTier, &t.PrizeAmount, &t.MatchedNumbers, &t.MatchedStars); err != nil {
			return nil, fmt.Errorf("error scanning ticket: %w", err)
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}

func (s *PostgreSQLStore) GetTicketByID(id int) (models.Ticket, error) {
	var t models.Ticket

	err := s.pool.QueryRow(context.Background(),
		`SELECT id, draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars
		 FROM tickets WHERE id = $1`, id,
	).Scan(&t.ID, &t.DrawID, &t.Numbers, &t.Stars, &t.Cost, &t.PurchasedAt,
		&t.PrizeTier, &t.PrizeAmount, &t.MatchedNumbers, &t.MatchedStars)

	if err == pgx.ErrNoRows {
		return models.Ticket{}, fmt.Errorf("ticket %d not found", id)
	}
	if err != nil {
		return models.Ticket{}, fmt.Errorf("error querying ticket: %w", err)
	}

	return t, nil
}

func (s *PostgreSQLStore) GetTicketsByDraw(drawID int) ([]models.Ticket, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars
		 FROM tickets WHERE draw_id = $1 ORDER BY purchased_at`, drawID)
	if err != nil {
		return nil, fmt.Errorf("error querying tickets: %w", err)
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.DrawID, &t.Numbers, &t.Stars, &t.Cost, &t.PurchasedAt,
			&t.PrizeTier, &t.PrizeAmount, &t.MatchedNumbers, &t.MatchedStars); err != nil {
			return nil, fmt.Errorf("error scanning ticket: %w", err)
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}

func (s *PostgreSQLStore) UpdateTicketPrize(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error) {
	var t models.Ticket

	err := s.pool.QueryRow(context.Background(),
		`UPDATE tickets
		 SET prize_tier = $1, prize_amount = $2, matched_numbers = $3, matched_stars = $4
		 WHERE id = $5
		 RETURNING id, draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars`,
		prizeTier, prizeAmount, matchedNumbers, matchedStars, id,
	).Scan(&t.ID, &t.DrawID, &t.Numbers, &t.Stars, &t.Cost, &t.PurchasedAt,
		&t.PrizeTier, &t.PrizeAmount, &t.MatchedNumbers, &t.MatchedStars)

	if err == pgx.ErrNoRows {
		return models.Ticket{}, fmt.Errorf("ticket %d not found", id)
	}
	if err != nil {
		return models.Ticket{}, fmt.Errorf("error updating ticket prize: %w", err)
	}

	return t, nil
}

func (s *PostgreSQLStore) DeleteTicket(id int) error {
	tag, err := s.pool.Exec(context.Background(),
		`DELETE FROM tickets WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting ticket: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ticket %d not found", id)
	}

	return nil
}
