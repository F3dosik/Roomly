package postgres

import (
	"context"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/db"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
)

func (r *postgresRepository) CreateSchedule(ctx context.Context, schedule *domain.Schedule) error {
	err := db.WithRetry(ctx, func() error {
		return r.pool.QueryRow(ctx, `
			INSERT INTO schedules (room_id, weekdays, start_time, end_time)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at
		`, schedule.RoomID, schedule.DaysOfWeek, schedule.StartTime, schedule.EndTime).
			Scan(&schedule.ID, &schedule.CreatedAt)
	})
	if err != nil {
		switch {
		case db.IsFKViolation(err):
			return domain.ErrRoomNotFound
		case db.IsUniqueViolation(err):
			return domain.ErrScheduleExists
		}
		return err
	}

	return nil
}
