package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"lottery-pool-manager/models"
)

func (s *PostgreSQLStore) CreateParticipant(name, email string) (models.Participant, error) {
	var p models.Participant

	err := s.pool.QueryRow(context.Background(),
		`INSERT INTO participants (name, email)
		 VALUES ($1, $2)
		 RETURNING id, name, email, active, created_at`,
		name, email,
	).Scan(&p.ID, &p.Name, &p.Email, &p.Active, &p.CreatedAt)

	if err != nil {
		return models.Participant{}, fmt.Errorf("error creating participant: %w", err)
	}

	return p, nil
}

func (s *PostgreSQLStore) GetAllParticipants() ([]models.Participant, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, name, email, active, created_at FROM participants ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("error querying participants: %w", err)
	}
	defer rows.Close()

	participants := make([]models.Participant, 0)
	for rows.Next() {
		var p models.Participant
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.Active, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning participant: %w", err)
		}
		participants = append(participants, p)
	}

	return participants, nil
}

func (s *PostgreSQLStore) GetParticipantByID(id int) (models.Participant, error) {
	var p models.Participant

	err := s.pool.QueryRow(context.Background(),
		`SELECT id, name, email, active, created_at FROM participants WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Email, &p.Active, &p.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Participant{}, fmt.Errorf("participant %d not found", id)
	}
	if err != nil {
		return models.Participant{}, fmt.Errorf("error querying participant: %w", err)
	}

	return p, nil
}

func (s *PostgreSQLStore) UpdateParticipant(id int, name, email string) (models.Participant, error) {
	var p models.Participant

	err := s.pool.QueryRow(context.Background(),
		`UPDATE participants
		 SET name = $1, email = $2
		 WHERE id = $3
		 RETURNING id, name, email, active, created_at`,
		name, email, id,
	).Scan(&p.ID, &p.Name, &p.Email, &p.Active, &p.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Participant{}, fmt.Errorf("participant %d not found", id)
	}
	if err != nil {
		return models.Participant{}, fmt.Errorf("error updating participant: %w", err)
	}

	return p, nil
}

func (s *PostgreSQLStore) DeactivateParticipant(id int) (models.Participant, error) {
	var p models.Participant

	err := s.pool.QueryRow(context.Background(),
		`UPDATE participants
		 SET active = false
		 WHERE id = $1
		 RETURNING id, name, email, active, created_at`,
		id,
	).Scan(&p.ID, &p.Name, &p.Email, &p.Active, &p.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Participant{}, fmt.Errorf("participant %d not found", id)
	}
	if err != nil {
		return models.Participant{}, fmt.Errorf("error deactivating participant: %w", err)
	}

	return p, nil
}

func (s *PostgreSQLStore) ActivateParticipant(id int) (models.Participant, error) {
	var p models.Participant

	err := s.pool.QueryRow(context.Background(),
		`UPDATE participants
		 SET active = true
		 WHERE id = $1
		 RETURNING id, name, email, active, created_at`,
		id,
	).Scan(&p.ID, &p.Name, &p.Email, &p.Active, &p.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.Participant{}, fmt.Errorf("participant %d not found", id)
	}
	if err != nil {
		return models.Participant{}, fmt.Errorf("error activating participant: %w", err)
	}

	return p, nil
}
