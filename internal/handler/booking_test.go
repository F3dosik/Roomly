package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateBooking(t *testing.T) {
	slotID := uuid.New()
	bookingID := uuid.New()

	tests := []struct {
		name           string
		body           string
		token          func(t *testing.T) string
		mockFn         func(ctx context.Context, userID, slotID uuid.UUID, createConferenceLink bool) (*domain.Booking, error)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:  "success",
			body:  `{"slotId":"` + slotID.String() + `"}`,
			token: userToken,
			mockFn: func(ctx context.Context, userID, sID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
				require.Equal(t, slotID, sID)
				require.False(t, createConferenceLink)
				return &domain.Booking{
					ID:     bookingID,
					UserID: userID,
					SlotID: sID,
					Status: domain.BookingStatusActive,
				}, nil
			},
			expectedStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				booking := resp["booking"].(map[string]any)
				require.Equal(t, bookingID.String(), booking["id"])
				require.Equal(t, "active", booking["status"])
			},
		},
		{
			name:  "success with conference link",
			body:  `{"slotId":"` + slotID.String() + `","createConferenceLink":true}`,
			token: userToken,
			mockFn: func(ctx context.Context, userID, sID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
				require.True(t, createConferenceLink)
				link := "https://meet.example.com/test"
				return &domain.Booking{
					ID:             bookingID,
					UserID:         userID,
					SlotID:         sID,
					Status:         domain.BookingStatusActive,
					ConferenceLink: &link,
				}, nil
			},
			expectedStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				booking := resp["booking"].(map[string]any)
				require.NotNil(t, booking["conferenceLink"])
			},
		},
		{
			name:           "bad json",
			body:           `{bad}`,
			token:          userToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing slotId",
			body:           `{}`,
			token:          userToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "slot not found",
			body:  `{"slotId":"` + slotID.String() + `"}`,
			token: userToken,
			mockFn: func(ctx context.Context, userID, sID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
				return nil, domain.ErrSlotNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:  "slot already booked",
			body:  `{"slotId":"` + slotID.String() + `"}`,
			token: userToken,
			mockFn: func(ctx context.Context, userID, sID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
				return nil, domain.ErrSlotAlreadyBooked
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:  "booking in past",
			body:  `{"slotId":"` + slotID.String() + `"}`,
			token: userToken,
			mockFn: func(ctx context.Context, userID, sID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
				return nil, domain.ErrBookingInPast
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "forbidden for admin",
			body:  `{"slotId":"` + slotID.String() + `"}`,
			token: adminToken,
			mockFn: func(ctx context.Context, userID, sID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
				return nil, nil
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "unauthorized",
			body:           `{"slotId":"` + slotID.String() + `"}`,
			token:          nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "internal error",
			body:  `{"slotId":"` + slotID.String() + `"}`,
			token: userToken,
			mockFn: func(ctx context.Context, userID, sID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
				return nil, domain.ErrInternal
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &bookingServiceMock{CreateBookingFn: tt.mockFn}
			h := newTestHandler(nil, nil, nil, mock)

			req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
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

func TestHandler_ListBookings(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		token          func(t *testing.T) string
		mockFn         func(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:  "success — default pagination",
			query: "",
			token: adminToken,
			mockFn: func(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
				require.Equal(t, 1, page)
				require.Equal(t, 20, pageSize)
				return []*domain.Booking{
					{ID: uuid.New(), Status: domain.BookingStatusActive},
				}, 1, nil
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				bookings := resp["bookings"].([]any)
				require.Len(t, bookings, 1)
				pagination := resp["pagination"].(map[string]any)
				require.Equal(t, float64(1), pagination["total"])
			},
		},
		{
			name:  "success — custom pagination",
			query: "?page=2&pageSize=5",
			token: adminToken,
			mockFn: func(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
				require.Equal(t, 2, page)
				require.Equal(t, 5, pageSize)
				return []*domain.Booking{}, 0, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid page",
			query:          "?page=0",
			token:          adminToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid pageSize — too large",
			query:          "?pageSize=101",
			token:          adminToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid pageSize — negative",
			query:          "?pageSize=-1",
			token:          adminToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "forbidden for user",
			query: "",
			token: userToken,
			mockFn: func(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
				return nil, 0, nil
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "unauthorized",
			query:          "",
			token:          nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "internal error",
			query: "",
			token: adminToken,
			mockFn: func(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
				return nil, 0, domain.ErrInternal
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &bookingServiceMock{ListBookingsFn: tt.mockFn}
			h := newTestHandler(nil, nil, nil, mock)

			req := httptest.NewRequest(http.MethodGet, "/bookings/list"+tt.query, nil)
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

func TestHandler_GetMyBookings(t *testing.T) {
	tests := []struct {
		name           string
		token          func(t *testing.T) string
		mockFn         func(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:  "success",
			token: userToken,
			mockFn: func(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
				return []*domain.Booking{
					{
						ID:        uuid.New(),
						UserID:    userID,
						SlotID:    uuid.New(),
						Status:    domain.BookingStatusActive,
						CreatedAt: time.Now(),
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				bookings := resp["bookings"].([]any)
				require.Len(t, bookings, 1)
			},
		},
		{
			name:  "success — empty list",
			token: userToken,
			mockFn: func(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
				return []*domain.Booking{}, nil
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				bookings := resp["bookings"].([]any)
				require.Empty(t, bookings)
			},
		},
		{
			name:  "forbidden for admin",
			token: adminToken,
			mockFn: func(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
				return nil, nil
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "unauthorized",
			token:          nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "internal error",
			token: userToken,
			mockFn: func(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
				return nil, domain.ErrInternal
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &bookingServiceMock{GetMyBookingsFn: tt.mockFn}
			h := newTestHandler(nil, nil, nil, mock)

			req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
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

func TestHandler_CancelBooking(t *testing.T) {
	bookingID := uuid.New()

	tests := []struct {
		name           string
		bookingID      string
		token          func(t *testing.T) string
		mockFn         func(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:      "success",
			bookingID: bookingID.String(),
			token:     userToken,
			mockFn: func(ctx context.Context, userID, bID uuid.UUID) (*domain.Booking, error) {
				return &domain.Booking{
					ID:     bID,
					UserID: userID,
					Status: domain.BookingStatusCancelled,
				}, nil
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				var resp map[string]any
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				booking := resp["booking"].(map[string]any)
				require.Equal(t, "cancelled", booking["status"])
			},
		},
		{
			name:      "idempotent — already cancelled",
			bookingID: bookingID.String(),
			token:     userToken,
			mockFn: func(ctx context.Context, userID, bID uuid.UUID) (*domain.Booking, error) {
				return &domain.Booking{
					ID:     bID,
					UserID: userID,
					Status: domain.BookingStatusCancelled,
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid bookingId",
			bookingID:      "not-a-uuid",
			token:          userToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "booking not found",
			bookingID: bookingID.String(),
			token:     userToken,
			mockFn: func(ctx context.Context, userID, bID uuid.UUID) (*domain.Booking, error) {
				return nil, domain.ErrBookingNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "forbidden — not own booking",
			bookingID: bookingID.String(),
			token:     userToken,
			mockFn: func(ctx context.Context, userID, bID uuid.UUID) (*domain.Booking, error) {
				return nil, domain.ErrForbidden
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:      "forbidden for admin",
			bookingID: bookingID.String(),
			token:     adminToken,
			mockFn: func(ctx context.Context, userID, bID uuid.UUID) (*domain.Booking, error) {
				return nil, nil
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "unauthorized",
			bookingID:      bookingID.String(),
			token:          nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:      "internal error",
			bookingID: bookingID.String(),
			token:     userToken,
			mockFn: func(ctx context.Context, userID, bID uuid.UUID) (*domain.Booking, error) {
				return nil, domain.ErrInternal
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &bookingServiceMock{CancelBookingFn: tt.mockFn}
			h := newTestHandler(nil, nil, nil, mock)

			url := "/bookings/" + tt.bookingID + "/cancel"
			req := httptest.NewRequest(http.MethodPost, url, nil)
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
