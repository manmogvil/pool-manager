package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"lottery-pool-manager/models"
)

func (s *PostgreSQLStore) CreateDraw(gameID int, drawDate time.Time) (models.Draw, error) {
	var d models.Draw

	var exists bool
	err := s.pool.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM draws WHERE game_id = $1 AND draw_date = $2)`,
		gameID, drawDate,
	).Scan(&exists)
	if err != nil {
		return models.Draw{}, fmt.Errorf("error checking existing draw: %w", err)
	}
	if exists {
		return models.Draw{}, fmt.Errorf("a draw already exists for this game on %s", drawDate.Format("2006-01-02"))
	}

	err = s.pool.QueryRow(context.Background(),
		`INSERT INTO draws (game_id, draw_date)
		 VALUES ($1, $2)
		 RETURNING id, game_id, draw_date, result_numbers, result_stars, draw_id_api, processed, created_at`,
		gameID, drawDate,
	).Scan(&d.ID, &d.GameID, &d.DrawDate, &d.ResultNumbers, &d.ResultStars, &d.DrawIDAPI, &d.Processed, &d.CreatedAt)

	if err != nil {
		return models.Draw{}, fmt.Errorf("error creating draw: %w", err)
	}

	return d, nil
}

func (s *PostgreSQLStore) GetAllDraws() ([]models.Draw, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, game_id, draw_date, result_numbers, result_stars, draw_id_api, processed, created_at
		 FROM draws ORDER BY draw_date DESC`)
	if err != nil {
		return nil, fmt.Errorf("error querying draws: %w", err)
	}
	defer rows.Close()

	draws := make([]models.Draw, 0)
	for rows.Next() {
		var d models.Draw
		if err := rows.Scan(&d.ID, &d.GameID, &d.DrawDate, &d.ResultNumbers, &d.ResultStars, &d.DrawIDAPI, &d.Processed, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning draw: %w", err)
		}
		draws = append(draws, d)
	}

	return draws, nil
}

func (s *PostgreSQLStore) GetDrawByID(id int) (models.Draw, error) {
	var d models.Draw

	err := s.pool.QueryRow(context.Background(),
		`SELECT id, game_id, draw_date, result_numbers, result_stars, draw_id_api, processed, created_at
		 FROM draws WHERE id = $1`, id,
	).Scan(&d.ID, &d.GameID, &d.DrawDate, &d.ResultNumbers, &d.ResultStars, &d.DrawIDAPI, &d.Processed, &d.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Draw{}, fmt.Errorf("draw %d not found", id)
	}
	if err != nil {
		return models.Draw{}, fmt.Errorf("error querying draw: %w", err)
	}

	return d, nil
}

func (s *PostgreSQLStore) GetDrawsByGame(gameID int) ([]models.Draw, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, game_id, draw_date, result_numbers, result_stars, draw_id_api, processed, created_at
		 FROM draws WHERE game_id = $1 ORDER BY draw_date DESC`, gameID)
	if err != nil {
		return nil, fmt.Errorf("error querying draws: %w", err)
	}
	defer rows.Close()

	draws := make([]models.Draw, 0)
	for rows.Next() {
		var d models.Draw
		if err := rows.Scan(&d.ID, &d.GameID, &d.DrawDate, &d.ResultNumbers, &d.ResultStars, &d.DrawIDAPI, &d.Processed, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning draw: %w", err)
		}
		draws = append(draws, d)
	}

	return draws, nil
}

func (s *PostgreSQLStore) GetPendingDraws() ([]models.Draw, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, game_id, draw_date, result_numbers, result_stars, draw_id_api, processed, created_at
		 FROM draws WHERE processed = false ORDER BY draw_date`)
	if err != nil {
		return nil, fmt.Errorf("error querying pending draws: %w", err)
	}
	defer rows.Close()

	draws := make([]models.Draw, 0)
	for rows.Next() {
		var d models.Draw
		if err := rows.Scan(&d.ID, &d.GameID, &d.DrawDate, &d.ResultNumbers, &d.ResultStars, &d.DrawIDAPI, &d.Processed, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning draw: %w", err)
		}
		draws = append(draws, d)
	}

	return draws, nil
}

func (s *PostgreSQLStore) UpdateDrawResults(id int, resultNumbers *string, resultStars *string) (models.Draw, error) {
	var d models.Draw

	err := s.pool.QueryRow(context.Background(),
		`UPDATE draws
		 SET result_numbers = $1, result_stars = $2
		 WHERE id = $3
		 RETURNING id, game_id, draw_date, result_numbers, result_stars, draw_id_api, processed, created_at`,
		resultNumbers, resultStars, id,
	).Scan(&d.ID, &d.GameID, &d.DrawDate, &d.ResultNumbers, &d.ResultStars, &d.DrawIDAPI, &d.Processed, &d.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Draw{}, fmt.Errorf("draw %d not found", id)
	}
	if err != nil {
		return models.Draw{}, fmt.Errorf("error updating draw results: %w", err)
	}

	return d, nil
}

func (s *PostgreSQLStore) UpdateDrawDrawIDAPI(id int, drawIDAPI string) (models.Draw, error) {
	var d models.Draw

	err := s.pool.QueryRow(context.Background(),
		`UPDATE draws
		 SET draw_id_api = $1
		 WHERE id = $2
		 RETURNING id, game_id, draw_date, result_numbers, result_stars, draw_id_api, processed, created_at`,
		drawIDAPI, id,
	).Scan(&d.ID, &d.GameID, &d.DrawDate, &d.ResultNumbers, &d.ResultStars, &d.DrawIDAPI, &d.Processed, &d.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Draw{}, fmt.Errorf("draw %d not found", id)
	}
	if err != nil {
		return models.Draw{}, fmt.Errorf("error updating draw draw_id_api: %w", err)
	}

	return d, nil
}

func (s *PostgreSQLStore) MarkDrawAsProcessed(id int) error {
	tag, err := s.pool.Exec(context.Background(),
		`UPDATE draws SET processed = true WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error marking draw as processed: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("draw %d not found", id)
	}

	return nil
}

func (s *PostgreSQLStore) DeleteDraw(id int) error {
	var hasTickets bool

	err := s.pool.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM tickets WHERE draw_id = $1)`, id,
	).Scan(&hasTickets)

	if err != nil {
		return fmt.Errorf("error checking tickets for draw: %w", err)
	}

	if hasTickets {
		return fmt.Errorf("cannot delete draw %d: tickets exist for this draw", id)
	}

	tag, err := s.pool.Exec(context.Background(),
		`DELETE FROM draws WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting draw: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("draw %d not found", id)
	}

	return nil
}
