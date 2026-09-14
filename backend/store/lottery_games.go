package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"lottery-pool-manager/models"
)

func (s *PostgreSQLStore) CreateLotteryGame(name, drawDays string, ticketPrice float64) (models.LotteryGame, error) {
	var game models.LotteryGame

	err := s.pool.QueryRow(context.Background(),
		`INSERT INTO lottery_games (name, draw_days, ticket_price)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, draw_days, ticket_price, active`,
		name, drawDays, ticketPrice,
	).Scan(&game.ID, &game.Name, &game.DrawDays, &game.TicketPrice, &game.Active)

	if err != nil {
		return models.LotteryGame{}, fmt.Errorf("error creating lottery game: %w", err)
	}

	return game, nil
}

func (s *PostgreSQLStore) GetAllLotteryGames() ([]models.LotteryGame, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, name, draw_days, ticket_price, active FROM lottery_games ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("error querying lottery games: %w", err)
	}
	defer rows.Close()

	games := make([]models.LotteryGame, 0)
	for rows.Next() {
		var g models.LotteryGame
		if err := rows.Scan(&g.ID, &g.Name, &g.DrawDays, &g.TicketPrice, &g.Active); err != nil {
			return nil, fmt.Errorf("error scanning lottery game: %w", err)
		}
		games = append(games, g)
	}

	return games, nil
}

func (s *PostgreSQLStore) GetLotteryGameByID(id int) (models.LotteryGame, error) {
	var g models.LotteryGame

	err := s.pool.QueryRow(context.Background(),
		`SELECT id, name, draw_days, ticket_price, active FROM lottery_games WHERE id = $1`, id,
	).Scan(&g.ID, &g.Name, &g.DrawDays, &g.TicketPrice, &g.Active)

	if err == pgx.ErrNoRows {
		return models.LotteryGame{}, fmt.Errorf("lottery game %d not found", id)
	}
	if err != nil {
		return models.LotteryGame{}, fmt.Errorf("error querying lottery game: %w", err)
	}

	return g, nil
}

func (s *PostgreSQLStore) UpdateLotteryGame(id int, name, drawDays string, ticketPrice float64, active bool) (models.LotteryGame, error) {
	var g models.LotteryGame

	err := s.pool.QueryRow(context.Background(),
		`UPDATE lottery_games
		 SET name = $1, draw_days = $2, ticket_price = $3, active = $4
		 WHERE id = $5
		 RETURNING id, name, draw_days, ticket_price, active`,
		name, drawDays, ticketPrice, active, id,
	).Scan(&g.ID, &g.Name, &g.DrawDays, &g.TicketPrice, &g.Active)

	if err == pgx.ErrNoRows {
		return models.LotteryGame{}, fmt.Errorf("lottery game %d not found", id)
	}
	if err != nil {
		return models.LotteryGame{}, fmt.Errorf("error updating lottery game: %w", err)
	}

	return g, nil
}

func (s *PostgreSQLStore) DeleteLotteryGame(id int) error {
	var hasTickets bool

	err := s.pool.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM tickets t
			JOIN draws d ON t.draw_id = d.id
			WHERE d.game_id = $1
		)`, id,
	).Scan(&hasTickets)

	if err != nil {
		return fmt.Errorf("error checking tickets for game: %w", err)
	}

	if hasTickets {
		return fmt.Errorf("cannot delete game %d: tickets exist for this game", id)
	}

	tag, err := s.pool.Exec(context.Background(),
		`DELETE FROM lottery_games WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting lottery game: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("lottery game %d not found", id)
	}

	return nil
}

func (s *PostgreSQLStore) GetLotteryGameByName(name string) (models.LotteryGame, error) {
	var g models.LotteryGame

	err := s.pool.QueryRow(context.Background(),
		`SELECT id, name, draw_days, ticket_price, active FROM lottery_games WHERE name = $1`, name,
	).Scan(&g.ID, &g.Name, &g.DrawDays, &g.TicketPrice, &g.Active)

	if err == pgx.ErrNoRows {
		return models.LotteryGame{}, fmt.Errorf("lottery game %s not found", name)
	}
	if err != nil {
		return models.LotteryGame{}, fmt.Errorf("error querying lottery game: %w", err)
	}

	return g, nil
}
