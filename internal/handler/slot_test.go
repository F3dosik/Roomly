package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetFreeSlots(t *testing.T) {
	roomID := uuid.New()

	tests := []struct {
		name           string
		query          string
		token          func(t *testing.T) string
		mockFn         func(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:  "success — returns slots",
			query: "?date=2024-06-10",
			token: userToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
				require.Equal(t, roomID, rID)
				require.Equal(t, 2024, date.Year())
				require.Equal(t, time.June, date.Month())
				require.Equal(t, 10, date.Day())
				return []*domain.Slot{
					{
						ID:       uuid.New(),
						RoomID:   rID,
						StartsAt: time.Date(2024, 6, 10, 9, 0, 0, 0, time.UTC),
						EndsAt:   time.Date(2024, 6, 10, 9, 30, 0, 0, time.UTC),
					},
					{
						ID:       uuid.New(),
						RoomID:   rID,
						StartsAt: time.Date(2024, 6, 10, 9, 30, 0, 0, time.UTC),
						EndsAt:   time.Date(2024, 6, 10, 10, 0, 0, 0, time.UTC),
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				slots := resp["slots"].([]any)
				require.Len(t, slots, 2)
			},
		},
		{
			name:  "success — empty list",
			query: "?date=2024-06-10",
			token: userToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
				return []*domain.Slot{}, nil
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				slots := resp["slots"].([]any)
				require.Empty(t, slots)
			},
		},
		{
			name:           "missing date param",
			query:          "",
			token:          userToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid date format",
			query:          "?date=10-06-2024",
			token:          userToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "room not found",
			query: "?date=2024-06-10",
			token: userToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
				return nil, domain.ErrRoomNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:  "internal error",
			query: "?date=2024-06-10",
			token: userToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
				return nil, domain.ErrInternal
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unauthorized",
			query:          "?date=2024-06-10",
			token:          nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &slotServiceMock{GetFreeSlotsFn: tt.mockFn}
			h := newTestHandler(nil, nil, mock, nil)

			url := "/rooms/" + roomID.String() + "/slots/list" + tt.query
			req := httptest.NewRequest(http.MethodGet, url, nil)
			if tt.token != nil {
				req.Header.Set("Authorization", "Bearer "+tt.token(t))
			}

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, rec.Body.String())
			}
		})
	}
}
