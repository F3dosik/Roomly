package postgres

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/db"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
)

func (r *postgresRepository) GetSlotByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error) {
	var slot domain.Slot
	err := r.pool.QueryRow(ctx, `
        SELECT id, room_id, starts_at, ends_at
        FROM slots
        WHERE id = $1
    `, slotID).Scan(&slot.ID, &slot.RoomID, &slot.StartsAt, &slot.EndsAt)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, domain.ErrSlotNotFound
		}
		return nil, fmt.Errorf("get slot by id: %w", err)
	}
	return &slot, nil
}

func (r *postgresRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	err := db.WithRetry(ctx, func() error {
		return r.pool.QueryRow(ctx, `
            INSERT INTO bookings (user_id, slot_id, status, conference_link)
            VALUES ($1, $2, 'active', $3)
            RETURNING id, created_at
        `, booking.UserID, booking.SlotID, booking.ConferenceLink).
			Scan(&booking.ID, &booking.CreatedAt)
	})
	if err != nil {
		if db.IsUniqueViolation(err) {
			return domain.ErrSlotAlreadyBooked
		}
		return fmt.Errorf("create booking: %w", err)
	}
	return nil
}

func (r *postgresRepository) ListBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM bookings`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx, `
        SELECT id, user_id, slot_id, status, conference_link, created_at, cancelled_at
        FROM bookings
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.SlotID, &b.Status,
			&b.ConferenceLink, &b.CreatedAt, &b.CancelledAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan booking: %w", err)
		}
		bookings = append(bookings, &b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return bookings, total, nil
}

func (r *postgresRepository) GetMyBookings(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT b.id, b.user_id, b.slot_id, b.status, b.conference_link, b.created_at, b.cancelled_at
        FROM bookings b
        JOIN slots s ON s.id = b.slot_id
        WHERE b.user_id = $1
          AND s.starts_at > NOW()
        ORDER BY s.starts_at ASC
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("get my bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.SlotID, &b.Status,
			&b.ConferenceLink, &b.CreatedAt, &b.CancelledAt,
		); err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		bookings = append(bookings, &b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return bookings, nil
}

func (r *postgresRepository) GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error) {
	var b domain.Booking
	err := r.pool.QueryRow(ctx, `
        SELECT id, user_id, slot_id, status, conference_link, created_at, cancelled_at
        FROM bookings
        WHERE id = $1
    `, bookingID).Scan(
		&b.ID, &b.UserID, &b.SlotID, &b.Status,
		&b.ConferenceLink, &b.CreatedAt, &b.CancelledAt,
	)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("get booking by id: %w", err)
	}
	return &b, nil
}

func (r *postgresRepository) CancelBooking(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error) {
	var b domain.Booking
	err := r.pool.QueryRow(ctx, `
        UPDATE bookings
        SET status = 'cancelled', cancelled_at = NOW()
        WHERE id = $1
        RETURNING id, user_id, slot_id, status, conference_link, created_at, cancelled_at
    `, bookingID).Scan(
		&b.ID, &b.UserID, &b.SlotID, &b.Status,
		&b.ConferenceLink, &b.CreatedAt, &b.CancelledAt,
	)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("cancel booking: %w", err)
	}
	return &b, nil
}
