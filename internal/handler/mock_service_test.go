package handler

import (
	"context"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/jwt"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const testJWTSecret = "test-secret"

func adminToken(t *testing.T) string {
	t.Helper()
	token, err := jwt.GenerateToken(uuid.New(), domain.RoleAdmin, testJWTSecret)
	require.NoError(t, err)
	return token
}

func userToken(t *testing.T) string {
	t.Helper()
	token, err := jwt.GenerateToken(uuid.New(), domain.RoleUser, testJWTSecret)
	require.NoError(t, err)
	return token
}

func newTestHandler(us service.UserService, rs service.RoomService, ss service.SlotService, bs service.BookingService) *Handler {
	logger := zap.NewNop().Sugar()
	return New("test-secret", us, rs, ss, bs, logger)
}

type userServiceMock struct {
	DummyLoginFn func(ctx context.Context, role domain.Role) (string, error)
	RegisterFn   func(ctx context.Context, email, password string, role domain.Role) (*domain.User, error)
	LoginFn      func(ctx context.Context, email, password string) (string, error)
}

func (m *userServiceMock) DummyLogin(ctx context.Context, role domain.Role) (string, error) {
	return m.DummyLoginFn(ctx, role)
}

func (m *userServiceMock) Register(ctx context.Context, email, password string, role domain.Role) (*domain.User, error) {
	return m.RegisterFn(ctx, email, password, role)
}

func (m *userServiceMock) Login(ctx context.Context, email, password string) (string, error) {
	return m.LoginFn(ctx, email, password)
}

type roomServiceMock struct {
	GetRoomsFn       func(ctx context.Context) ([]*domain.Room, error)
	CreateRoomFn     func(ctx context.Context, room *domain.Room) error
	CreateScheduleFn func(ctx context.Context, roomID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error)
}

func (m *roomServiceMock) GetRooms(ctx context.Context) ([]*domain.Room, error) {
	return m.GetRoomsFn(ctx)
}

func (m *roomServiceMock) CreateRoom(ctx context.Context, room *domain.Room) error {
	return m.CreateRoomFn(ctx, room)
}

func (m *roomServiceMock) CreateSchedule(ctx context.Context, roomID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error) {
	return m.CreateScheduleFn(ctx, roomID, days, start, end)
}

type slotServiceMock struct {
	GetFreeSlotsFn func(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error)
}

func (m *slotServiceMock) GetFreeSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	return m.GetFreeSlotsFn(ctx, roomID, date)
}

type bookingServiceMock struct {
	CreateBookingFn func(ctx context.Context, userID, slotID uuid.UUID, createConferenceLink bool) (*domain.Booking, error)
	ListBookingsFn  func(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error)
	GetMyBookingsFn func(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error)
	CancelBookingFn func(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error)
}

func (m *bookingServiceMock) CreateBooking(ctx context.Context, userID, slotID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
	return m.CreateBookingFn(ctx, userID, slotID, createConferenceLink)
}

func (m *bookingServiceMock) ListBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
	return m.ListBookingsFn(ctx, page, pageSize)
}

func (m *bookingServiceMock) GetMyBookings(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
	return m.GetMyBookingsFn(ctx, userID)
}

func (m *bookingServiceMock) CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error) {
	return m.CancelBookingFn(ctx, userID, bookingID)
}
