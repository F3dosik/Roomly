package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestHandler(us service.UserService, rs service.RoomService) *Handler {
	logger := zap.NewNop().Sugar()
	return New("test-secret", us, rs, logger)
}

func TestHandler_DummyLogin(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockFn         func(ctx context.Context, role domain.Role) (string, error)
		expectedStatus int
	}{
		{
			name: "success user",
			body: `{"role":"user"}`,
			mockFn: func(ctx context.Context, role domain.Role) (string, error) {
				require.Equal(t, domain.Role("user"), role)
				return "jwt-token", nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid role",
			body: `{"role":"invalid"}`,
			mockFn: func(ctx context.Context, role domain.Role) (string, error) {
				return "", service.ErrInvalidRole
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "bad json",
			body:           `invalid-json`,
			mockFn:         nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &userServiceMock{
				DummyLoginFn: tt.mockFn,
			}

			h := newTestHandler(mock, nil)

			req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if rec.Code == http.StatusOK {
				var resp Token
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				require.NotEmpty(t, resp.Token)
			}
		})
	}
}

func TestHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockFn         func(ctx context.Context, email, password string, role domain.Role) (*domain.User, error)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"email":"test@mail.com","password":"123456","role":"user"}`,
			mockFn: func(ctx context.Context, email, password string, role domain.Role) (*domain.User, error) {
				return &domain.User{
					ID:    uuid.New(),
					Email: email,
					Role:  role,
				}, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid email",
			body: `{"email":"bad","password":"123456","role":"user"}`,
			mockFn: func(ctx context.Context, email, password string, role domain.Role) (*domain.User, error) {
				return &domain.User{}, service.ErrInvalidEmail
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "bad json",
			body:           `{bad}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &userServiceMock{
				RegisterFn: tt.mockFn,
			}

			h := newTestHandler(mock, nil)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if rec.Code == http.StatusCreated {
				require.Contains(t, rec.Body.String(), "email")
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockFn         func(ctx context.Context, email, password string) (string, error)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"email":"test@mail.com","password":"123456"}`,
			mockFn: func(ctx context.Context, email, password string) (string, error) {
				return "jwt-token", nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			body: `{"email":"test@mail.com","password":"wrong"}`,
			mockFn: func(ctx context.Context, email, password string) (string, error) {
				return "", service.ErrInvalidCredentials
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "bad json",
			body:           `bad`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &userServiceMock{
				LoginFn: tt.mockFn,
			}

			h := newTestHandler(mock, nil)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if rec.Code == http.StatusOK {
				var resp Token
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				require.NotEmpty(t, resp.Token)
			}
		})
	}
}
