package service

import (
	"context"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
)

type mockRepository struct {
	CreateUserFn      func(ctx context.Context, email, password string, role domain.Role) (*domain.User, error)
	GetUserFn         func(ctx context.Context, email string) (*domain.User, error)
	UpsertUserFn      func(ctx context.Context, id uuid.UUID, email string, role domain.Role) error
	GetRoomsFn        func(ctx context.Context) ([]*domain.Room, error)
	CreateRoomFn      func(ctx context.Context, room *domain.Room) error
	CreateScheduleFn  func(ctx context.Context, schedule *domain.Schedule) error
	GetFreeSlotsFn    func(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error)
	GetSlotByIDFn     func(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error)
	CreateBookingFn   func(ctx context.Context, booking *domain.Booking) error
	ListBookingsFn    func(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error)
	GetMyBookingsFn   func(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error)
	GetBookingByIDFn  func(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error)
	CancelBookingFn   func(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error)
	GetAllSchedulesFn func(ctx context.Context) ([]*domain.Schedule, error)
	GenerateSlotsFn   func(ctx context.Context, slots []*domain.Slot) error
	GetLastSlotDateFn func(ctx context.Context, roomID uuid.UUID) (*time.Time, error)
}

func (m *mockRepository) CreateUser(ctx context.Context, email, password string, role domain.Role) (*domain.User, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, email, password, role)
	}
	return nil, nil
}

func (m *mockRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	if m.GetUserFn != nil {
		return m.GetUserFn(ctx, email)
	}
	return nil, nil
}

func (m *mockRepository) UpsertUser(ctx context.Context, id uuid.UUID, email string, role domain.Role) error {
	if m.UpsertUserFn != nil {
		return m.UpsertUserFn(ctx, id, email, role)
	}
	return nil
}

func (m *mockRepository) GetRooms(ctx context.Context) ([]*domain.Room, error) {
	if m.GetRoomsFn != nil {
		return m.GetRoomsFn(ctx)
	}
	return nil, nil
}

func (m *mockRepository) CreateRoom(ctx context.Context, room *domain.Room) error {
	if m.CreateRoomFn != nil {
		return m.CreateRoomFn(ctx, room)
	}
	return nil
}

func (m *mockRepository) CreateSchedule(ctx context.Context, schedule *domain.Schedule) error {
	if m.CreateScheduleFn != nil {
		return m.CreateScheduleFn(ctx, schedule)
	}
	return nil
}

func (m *mockRepository) GetFreeSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	if m.GetFreeSlotsFn != nil {
		return m.GetFreeSlotsFn(ctx, roomID, date)
	}
	return nil, nil
}

func (m *mockRepository) GetSlotByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error) {
	if m.GetSlotByIDFn != nil {
		return m.GetSlotByIDFn(ctx, slotID)
	}
	return nil, nil
}

func (m *mockRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	if m.CreateBookingFn != nil {
		return m.CreateBookingFn(ctx, booking)
	}
	return nil
}

func (m *mockRepository) ListBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
	if m.ListBookingsFn != nil {
		return m.ListBookingsFn(ctx, page, pageSize)
	}
	return nil, 0, nil
}

func (m *mockRepository) GetMyBookings(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
	if m.GetMyBookingsFn != nil {
		return m.GetMyBookingsFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockRepository) GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error) {
	if m.GetBookingByIDFn != nil {
		return m.GetBookingByIDFn(ctx, bookingID)
	}
	return nil, nil
}

func (m *mockRepository) CancelBooking(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error) {
	if m.CancelBookingFn != nil {
		return m.CancelBookingFn(ctx, bookingID)
	}
	return nil, nil
}

func (m *mockRepository) GetAllSchedules(ctx context.Context) ([]*domain.Schedule, error) {
	if m.GetAllSchedulesFn != nil {
		return m.GetAllSchedulesFn(ctx)
	}
	return nil, nil
}

func (m *mockRepository) GenerateSlots(ctx context.Context, slots []*domain.Slot) error {
	if m.GenerateSlotsFn != nil {
		return m.GenerateSlotsFn(ctx, slots)
	}
	return nil
}

func (m *mockRepository) GetLastSlotDate(ctx context.Context, roomID uuid.UUID) (*time.Time, error) {
	if m.GetLastSlotDateFn != nil {
		return m.GetLastSlotDateFn(ctx, roomID)
	}
	return nil, nil
}
