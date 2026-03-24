package handler

import (
	"context"
	"testing"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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
