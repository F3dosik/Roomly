package service

import (
	"context"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSlotService_GetFreeSlots(t *testing.T) {
	roomID := uuid.New()
	slotID1 := uuid.New()
	slotID2 := uuid.New()
	date := time.Now().UTC()

	repo := &mockRepository{
		GetFreeSlotsFn: func(ctx context.Context, rID uuid.UUID, d time.Time) ([]*domain.Slot, error) {
			require.Equal(t, roomID, rID)
			require.Equal(t, date.Year(), d.Year())
			require.Equal(t, date.Month(), d.Month())
			require.Equal(t, date.Day(), d.Day())
			return []*domain.Slot{
				{
					ID:       slotID1,
					RoomID:   roomID,
					StartsAt: date.Add(9 * time.Hour),
					EndsAt:   date.Add(10 * time.Hour),
				},
				{
					ID:       slotID2,
					RoomID:   roomID,
					StartsAt: date.Add(10 * time.Hour),
					EndsAt:   date.Add(11 * time.Hour),
				},
			}, nil
		},
	}
	ss := NewSlotService(repo)

	slots, err := ss.GetFreeSlots(context.Background(), roomID, date)
	require.NoError(t, err)
	require.Len(t, slots, 2)
	require.Equal(t, slotID1, slots[0].ID)
	require.Equal(t, slotID2, slots[1].ID)
}
