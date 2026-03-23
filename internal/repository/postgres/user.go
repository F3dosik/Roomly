package postgres

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/db"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
)

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
