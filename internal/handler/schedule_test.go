package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/service"
)

func TestHandler_CreateSchedule(t *testing.T) {
	roomID := uuid.New()
	scheduleID := uuid.New()

	tests := []struct {
		name           string
		body           string
		token          func(t *testing.T) string
		mockFn         func(ctx context.Context, roomID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:  "success",
			body:  `{"daysOfWeek":[1,2,3],"startTime":"09:00","endTime":"18:00"}`,
			token: adminToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error) {
				require.Equal(t, roomID, rID)
				return &domain.Schedule{
					ID:         scheduleID,
					RoomID:     rID,
					DaysOfWeek: days,
					StartTime:  time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					EndTime:    time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC),
				}, nil
			},
			expectedStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				var resp scheduleResponse
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.Equal(t, scheduleID, resp.ID)
				require.Equal(t, roomID, resp.RoomID)
				require.Equal(t, "09:00", resp.StartTime)
				require.Equal(t, "18:00", resp.EndTime)
			},
		},
		{
			name:           "bad json",
			body:           `{bad}`,
			token:          adminToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "forbidden for user role",
			body:  `{"daysOfWeek":[1],"startTime":"09:00","endTime":"18:00"}`,
			token: userToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error) {
				return nil, nil
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:  "invalid days of week",
			body:  `{"daysOfWeek":[0,8],"startTime":"09:00","endTime":"18:00"}`,
			token: adminToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error) {
				return nil, service.ErrInvalidDaysOfWeek
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "invalid time range",
			body:  `{"daysOfWeek":[1],"startTime":"18:00","endTime":"09:00"}`,
			token: adminToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error) {
				return nil, service.ErrInvalidTimeRange
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "room not found",
			body:  `{"daysOfWeek":[1],"startTime":"09:00","endTime":"18:00"}`,
			token: adminToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error) {
				return nil, domain.ErrRoomNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:  "schedule already exists",
			body:  `{"daysOfWeek":[1],"startTime":"09:00","endTime":"18:00"}`,
			token: adminToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error) {
				return nil, domain.ErrScheduleExists
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:  "internal error",
			body:  `{"daysOfWeek":[1],"startTime":"09:00","endTime":"18:00"}`,
			token: adminToken,
			mockFn: func(ctx context.Context, rID uuid.UUID, days []domain.DayOfWeek, start, end string) (*domain.Schedule, error) {
				return nil, service.ErrInternal
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &roomServiceMock{CreateScheduleFn: tt.mockFn}
			h := newTestHandler(nil, mock)

			url := "/rooms/" + roomID.String() + "/schedule/create"
			req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token(t))

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, rec.Body.String())
			}
		})
	}
}
