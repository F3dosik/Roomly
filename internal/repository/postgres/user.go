package postgres

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/db"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
)

func (r *postgresRepository) UpsertUser(ctx context.Context, id uuid.UUID, email string, role domain.Role) error {
	err := db.WithRetry(ctx, func() error {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO users (id, email, password_hash, role)
			VALUES ($1, $2, '', $3)
			ON CONFLICT (id) DO NOTHING
		`, id, email, role)
		return err
	})

	if err != nil {
		return fmt.Errorf("create fixed user: %w", err)
	}
	return nil
}

func (r *postgresRepository) CreateUser(
	ctx context.Context, email, password string,
	role domain.Role,
) (*domain.User, error) {
	var user domain.User
	err := db.WithRetry(ctx, func() error {
		return r.pool.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, role)
			VALUES ($1, $2, $3)
			RETURNING id, created_at
		`, email, password, role).Scan(&user.ID, &user.CreatedAt)
	})

	if err != nil {
		if db.IsUniqueViolation(err) {
			return nil, ErrEmailAlreadyExist
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	user.Email = email
	user.Role = role

	return &user, nil
}

func (r *postgresRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := db.WithRetry(ctx, func() error {
		return r.pool.QueryRow(ctx, `
			SELECT id, password_hash, role, created_at FROM users
			WHERE email = $1
		`, email).Scan(&user.ID, &user.PasswordHash, &user.Role, &user.CreatedAt)
	})
	if err != nil {
		if db.IsNoRows(err) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	user.Email = email

	return &user, nil
}
