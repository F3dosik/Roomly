package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
)

func TestHandler_GetRooms(t *testing.T) {
	tests := []struct {
		name           string
		mockFn         func(ctx context.Context) ([]*domain.Room, error)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name: "success — returns rooms",
			mockFn: func(ctx context.Context) ([]*domain.Room, error) {
				return []*domain.Room{
					{ID: uuid.New(), Name: "Room A"},
					{ID: uuid.New(), Name: "Room B"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				require.Contains(t, body, "Room A")
				require.Contains(t, body, "Room B")
			},
		},
		{
			name: "success — empty list",
			mockFn: func(ctx context.Context) ([]*domain.Room, error) {
				return []*domain.Room{}, nil
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				rooms := resp["rooms"].([]any)
				require.Empty(t, rooms)
			},
		},
		{
			name: "internal error",
			mockFn: func(ctx context.Context) ([]*domain.Room, error) {
				return nil, domain.ErrInternal
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &roomServiceMock{GetRoomsFn: tt.mockFn}
			h := newTestHandler(nil, mock, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
			req.Header.Set("Authorization", "Bearer "+userToken(t))
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, rec.Body.String())
			}
		})
	}
}

func TestHandler_CreateRoom(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		token          func(t *testing.T) string
		mockFn         func(ctx context.Context, room *domain.Room) error
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:  "success",
			body:  `{"name":"Boardroom","description":"Big room","capacity":10}`,
			token: adminToken,
			mockFn: func(ctx context.Context, room *domain.Room) error {
				room.ID = uuid.New()
				return nil
			},
			expectedStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				require.Contains(t, body, "Boardroom")
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
			body:  `{"name":"Room"}`,
			token: userToken,
			mockFn: func(ctx context.Context, room *domain.Room) error {
				return nil
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:  "internal error",
			body:  `{"name":"Room"}`,
			token: adminToken,
			mockFn: func(ctx context.Context, room *domain.Room) error {
				return domain.ErrInternal
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &roomServiceMock{CreateRoomFn: tt.mockFn}
			h := newTestHandler(nil, mock, nil, nil)

			req := httptest.NewRequest(http.MethodPost, "/rooms/create", strings.NewReader(tt.body))
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
