package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *postgresRepository) GetFreeSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	var roomExists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)
	`, roomID).Scan(&roomExists)
	if err != nil {
		return nil, fmt.Errorf("query roomID: %w", err)
	}
	if !roomExists {
		return nil, domain.ErrRoomNotFound
	}

	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.starts_at, s.ends_at
		FROM slots s
		WHERE s.room_id = $1
			AND s.starts_at::date= $2
			AND NOT EXISTS (
				SELECT 1 FROM bookings b
				WHERE b.slot_id = s.id
					AND b.status = 'active'
			)
	`, roomID, date)
	if err != nil {
		return nil, fmt.Errorf("query slots: %w", err)
	}
	defer rows.Close()

	var slots []*domain.Slot
	for rows.Next() {
		var slot domain.Slot
		slot.RoomID = roomID
		err := rows.Scan(
			&slot.ID, &slot.StartsAt, &slot.EndsAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan slot: %w", err)
		}
		slots = append(slots, &slot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return slots, nil
}

func (r *postgresRepository) GetAllSchedules(ctx context.Context) ([]*domain.Schedule, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, room_id, weekdays, start_time, end_time
        FROM schedules
    `)
	if err != nil {
		return nil, fmt.Errorf("get all schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*domain.Schedule
	for rows.Next() {
		var s domain.Schedule
		if err := rows.Scan(&s.ID, &s.RoomID, &s.DaysOfWeek, &s.StartTime, &s.EndTime); err != nil {
			return nil, fmt.Errorf("scan schedule: %w", err)
		}
		schedules = append(schedules, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return schedules, nil
}

func (r *postgresRepository) GetLastSlotDate(ctx context.Context, roomID uuid.UUID) (*time.Time, error) {
	var lastDate *time.Time
	err := r.pool.QueryRow(ctx, `
        SELECT MAX(starts_at) FROM slots WHERE room_id = $1
    `, roomID).Scan(&lastDate)
	if err != nil {
		return nil, fmt.Errorf("get last slot date: %w", err)
	}
	return lastDate, nil
}

func (r *postgresRepository) GenerateSlots(ctx context.Context, slots []*domain.Slot) error {
	batch := &pgx.Batch{}
	for _, slot := range slots {
		batch.Queue(`
            INSERT INTO slots (room_id, starts_at, ends_at)
            VALUES ($1, $2, $3)
            ON CONFLICT DO NOTHING
        `, slot.RoomID, slot.StartsAt, slot.EndsAt)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for range slots {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("insert slot: %w", err)
		}
	}
	return nil
}
