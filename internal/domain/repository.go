package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, email, password string, role Role) (*User, error)
	GetUser(ctx context.Context, email string) (*User, error)
	UpsertUser(ctx context.Context, id uuid.UUID, email string, role Role) error

	GetRooms(ctx context.Context) ([]*Room, error)
	CreateRoom(ctx context.Context, room *Room) error
	CreateSchedule(ctx context.Context, schedule *Schedule) error

	GetFreeSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*Slot, error)

	GetSlotByID(ctx context.Context, slotID uuid.UUID) (*Slot, error)
	CreateBooking(ctx context.Context, booking *Booking) error
	ListBookings(ctx context.Context, page, pageSize int) ([]*Booking, int, error)
	GetMyBookings(ctx context.Context, userID uuid.UUID) ([]*Booking, error)
	GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*Booking, error)
	CancelBooking(ctx context.Context, bookingID uuid.UUID) (*Booking, error)

	GetAllSchedules(ctx context.Context) ([]*Schedule, error)
	GenerateSlots(ctx context.Context, slots []*Slot) error
	GetLastSlotDate(ctx context.Context, roomID uuid.UUID) (*time.Time, error)
}
