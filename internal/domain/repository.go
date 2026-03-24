package domain

import (
	"context"
)

type Repository interface {
	// CreateRoom(ctx context.Context, room *Room) (*Room, error)
	// CreateSchedule(ctx context.Context, schedule *Schedule) (*Schedule, error)
	// GetAllBookings(ctx context.Context) ([]*Booking, error)

	// GetRooms(ctx context.Context) ([]*Room, error)
	// GetSlots(ctx context.Context, room_id uuid.UUID, date time.Time) ([]*Slot, error)
	// CreateBooking(ctx context.Context, user_id uuid.UUID) (*Booking, error)
	// DeleteBooking(ctx context.Context, booking_id *Booking) error
	// GetBookings(ctx context.Context, user_id uuid.UUID) ([]*Booking, error)

	CreateUser(ctx context.Context, email, password string, role Role) (*User, error)
	GetUser(ctx context.Context, email string) (*User, error)
}
