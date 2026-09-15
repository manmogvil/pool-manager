package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"lottery-pool-manager/models"
)

func (s *PostgreSQLStore) CreateUser(name, email, passwordHash, role string) (models.User, error) {
	var u models.User

	err := s.pool.QueryRow(context.Background(),
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, email, password_hash, role, active, created_at`,
		name, email, passwordHash, role,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.Active, &u.CreatedAt)

	if err != nil {
		return models.User{}, fmt.Errorf("error creating user: %w", err)
	}

	return u, nil
}

func (s *PostgreSQLStore) GetUserByEmail(email string) (models.User, error) {
	var u models.User

	err := s.pool.QueryRow(context.Background(),
		`SELECT id, name, email, password_hash, role, active, created_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.Active, &u.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return models.User{}, fmt.Errorf("error querying user: %w", err)
	}

	return u, nil
}

func (s *PostgreSQLStore) GetUserByID(id int) (models.User, error) {
	var u models.User

	err := s.pool.QueryRow(context.Background(),
		`SELECT id, name, email, password_hash, role, active, created_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.Active, &u.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return models.User{}, fmt.Errorf("error querying user: %w", err)
	}

	return u, nil
}

func (s *PostgreSQLStore) setUserActive(id int, active bool) error {
	tag, err := s.pool.Exec(context.Background(),
		`UPDATE users SET active = $1 WHERE id = $2`, active, id)
	if err != nil {
		return fmt.Errorf("error updating user active: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user %d not found", id)
	}
	return nil
}

func (s *PostgreSQLStore) ActivateUser(id int) error {
	return s.setUserActive(id, true)
}

func (s *PostgreSQLStore) DeactivateUser(id int) error {
	return s.setUserActive(id, false)
}

func (s *PostgreSQLStore) UpdateUser(id int, name, email, role string) (models.User, error) {
	var u models.User

	err := s.pool.QueryRow(context.Background(),
		`UPDATE users SET name = $1, email = $2, role = $3
		 WHERE id = $4
		 RETURNING id, name, email, password_hash, role, active, created_at`,
		name, email, role, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.Active, &u.CreatedAt)

	if err == pgx.ErrNoRows {
		return models.User{}, fmt.Errorf("user %d not found", id)
	}
	if err != nil {
		return models.User{}, fmt.Errorf("error updating user: %w", err)
	}

	return u, nil
}

func (s *PostgreSQLStore) GetAllUsers() ([]models.User, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, name, email, password_hash, role, active, created_at
		 FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.Active, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}
