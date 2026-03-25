package postgres

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/db"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
)

func (r *postgresRepository) GetRooms(ctx context.Context) ([]*domain.Room, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, description, capacity, created_at
		FROM rooms ORDER BY created_at
	`)
	if err != nil {
		return nil, fmt.Errorf("query rooms: %w", err)
	}
	defer rows.Close()

	var rooms []*domain.Room
	for rows.Next() {
		var room domain.Room
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Description,
			&room.Capacity,
			&room.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan room: %w", err)
		}
		rooms = append(rooms, &room)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return rooms, nil
}
func (r *postgresRepository) CreateRoom(ctx context.Context, room *domain.Room) error {
	err := db.WithRetry(ctx, func() error {
		return r.pool.QueryRow(ctx, `
			INSERT INTO rooms (name, description, capacity)
			VALUES ($1, $2, $3)
			RETURNING id, created_at
		`, room.Name, room.Description, room.Capacity).Scan(&room.ID, &room.CreatedAt)
	})
	if err != nil {
		return err
	}

	return nil
}
