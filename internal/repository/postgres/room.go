package postgres

import (
	"context"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
)

func (r *postgresRepository) GetList(ctx context.Context) ([]*domain.Room, error) {
	var rooms []*domain.Room
	rows, err := r.pool.Query(ctx, `
		SELECT * FROM rooms
	`)
}
