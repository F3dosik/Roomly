package service

import (
	"context"
	"testing"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/scheduler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRoomService_GetRooms(t *testing.T) {
	repo := &mockRepository{
		GetRoomsFn: func(ctx context.Context) ([]*domain.Room, error) {
			return []*domain.Room{
				{ID: uuid.New(), Name: "Room 1"},
			}, nil
		},
	}
	sg := &scheduler.SlotGenerator{}
	rs := NewRoomService(repo, sg)

	rooms, err := rs.GetRooms(context.Background())
	require.NoError(t, err)
	require.Len(t, rooms, 1)
	require.Equal(t, "Room 1", rooms[0].Name)
}

func TestRoomService_CreateRoom(t *testing.T) {
	roomID := uuid.New()
	repo := &mockRepository{
		CreateRoomFn: func(ctx context.Context, room *domain.Room) error {
			require.Equal(t, roomID, room.ID)
			require.Equal(t, "Test Room", room.Name)
			return nil
		},
	}
	sg := &scheduler.SlotGenerator{}
	rs := NewRoomService(repo, sg)

	err := rs.CreateRoom(context.Background(), &domain.Room{ID: roomID, Name: "Test Room"})
	require.NoError(t, err)
}

func TestRoomService_CreateSchedule(t *testing.T) {
	roomID := uuid.New()
	scheduleID := uuid.New()
	repo := &mockRepository{
		CreateScheduleFn: func(ctx context.Context, schedule *domain.Schedule) error {
			schedule.ID = scheduleID
			require.Equal(t, roomID, schedule.RoomID)
			return nil
		},
	}
	sg := &scheduler.SlotGenerator{}
	rs := NewRoomService(repo, sg)

	schedule, err := rs.CreateSchedule(context.Background(), roomID, []domain.DayOfWeek{1}, "09:00", "18:00")
	require.NoError(t, err)
	require.Equal(t, scheduleID, schedule.ID)
	require.Equal(t, roomID, schedule.RoomID)
}
