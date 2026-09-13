package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"lottery-pool-manager/models"
)

func (s *PostgreSQLStore) CreateContribution(c *models.Contribution) (models.Contribution, error) {
	var result models.Contribution

	err := s.pool.QueryRow(context.Background(),
		`INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_date, payment_method, comments)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, participant_id, game_id, month, year, amount, paid, payment_date, payment_method, comments, created_at`,
		c.ParticipantID, c.GameID, c.Month, c.Year, c.Amount, c.Paid, c.PaymentDate, c.PaymentMethod, c.Comments,
	).Scan(&result.ID, &result.ParticipantID, &result.GameID, &result.Month, &result.Year, &result.Amount,
		&result.Paid, &result.PaymentDate, &result.PaymentMethod, &result.Comments, &result.CreatedAt)

	if err != nil {
		return models.Contribution{}, fmt.Errorf("error creating contribution: %w", err)
	}

	return result, nil
}

func (s *PostgreSQLStore) GetAllContributions() ([]models.Contribution, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, participant_id, game_id, month, year, amount, paid, payment_date, payment_method, comments, created_at
		 FROM contributions ORDER BY year DESC, month DESC`)
	if err != nil {
		return nil, fmt.Errorf("error querying contributions: %w", err)
	}
	defer rows.Close()

	contributions := make([]models.Contribution, 0)
	for rows.Next() {
		var c models.Contribution
		if err := rows.Scan(&c.ID, &c.ParticipantID, &c.GameID, &c.Month, &c.Year, &c.Amount,
			&c.Paid, &c.PaymentDate, &c.PaymentMethod, &c.Comments, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning contribution: %w", err)
		}
		contributions = append(contributions, c)
	}

	return contributions, nil
}

func (s *PostgreSQLStore) GetContributionByID(id int) (models.Contribution, error) {
	var c models.Contribution

	err := s.pool.QueryRow(context.Background(),
		`SELECT id, participant_id, game_id, month, year, amount, paid, payment_date, payment_method, comments, created_at
		 FROM contributions WHERE id = $1`, id,
	).Scan(&c.ID, &c.ParticipantID, &c.GameID, &c.Month, &c.Year, &c.Amount,
		&c.Paid, &c.PaymentDate, &c.PaymentMethod, &c.Comments, &c.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Contribution{}, fmt.Errorf("contribution %d not found", id)
	}
	if err != nil {
		return models.Contribution{}, fmt.Errorf("error querying contribution: %w", err)
	}

	return c, nil
}

func (s *PostgreSQLStore) GetContributionsByParticipant(participantID int) ([]models.Contribution, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, participant_id, game_id, month, year, amount, paid, payment_date, payment_method, comments, created_at
		 FROM contributions WHERE participant_id = $1 ORDER BY year DESC, month DESC`, participantID)
	if err != nil {
		return nil, fmt.Errorf("error querying contributions: %w", err)
	}
	defer rows.Close()

	contributions := make([]models.Contribution, 0)
	for rows.Next() {
		var c models.Contribution
		if err := rows.Scan(&c.ID, &c.ParticipantID, &c.GameID, &c.Month, &c.Year, &c.Amount,
			&c.Paid, &c.PaymentDate, &c.PaymentMethod, &c.Comments, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning contribution: %w", err)
		}
		contributions = append(contributions, c)
	}

	return contributions, nil
}

func (s *PostgreSQLStore) GetContributionsByPeriod(month, year int) ([]models.Contribution, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, participant_id, game_id, month, year, amount, paid, payment_date, payment_method, comments, created_at
		 FROM contributions WHERE month = $1 AND year = $2 ORDER BY participant_id`, month, year)
	if err != nil {
		return nil, fmt.Errorf("error querying contributions: %w", err)
	}
	defer rows.Close()

	contributions := make([]models.Contribution, 0)
	for rows.Next() {
		var c models.Contribution
		if err := rows.Scan(&c.ID, &c.ParticipantID, &c.GameID, &c.Month, &c.Year, &c.Amount,
			&c.Paid, &c.PaymentDate, &c.PaymentMethod, &c.Comments, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning contribution: %w", err)
		}
		contributions = append(contributions, c)
	}

	return contributions, nil
}

func (s *PostgreSQLStore) UpdateContribution(id int, c *models.Contribution) (models.Contribution, error) {
	var result models.Contribution

	err := s.pool.QueryRow(context.Background(),
		`UPDATE contributions
		 SET participant_id = $1, game_id = $2, month = $3, year = $4, amount = $5, paid = $6,
		     payment_date = $7, payment_method = $8, comments = $9
		 WHERE id = $10
		 RETURNING id, participant_id, game_id, month, year, amount, paid, payment_date, payment_method, comments, created_at`,
		c.ParticipantID, c.GameID, c.Month, c.Year, c.Amount, c.Paid, c.PaymentDate, c.PaymentMethod, c.Comments, id,
	).Scan(&result.ID, &result.ParticipantID, &result.GameID, &result.Month, &result.Year, &result.Amount,
		&result.Paid, &result.PaymentDate, &result.PaymentMethod, &result.Comments, &result.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Contribution{}, fmt.Errorf("contribution %d not found", id)
	}
	if err != nil {
		return models.Contribution{}, fmt.Errorf("error updating contribution: %w", err)
	}

	return result, nil
}

func (s *PostgreSQLStore) DeleteContribution(id int) error {
	tag, err := s.pool.Exec(context.Background(),
		`DELETE FROM contributions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting contribution: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("contribution %d not found", id)
	}

	return nil
}
