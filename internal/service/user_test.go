package service

import (
	"context"
	"testing"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserService_DummyLogin(t *testing.T) {
	repo := &mockRepository{
		UpsertUserFn: func(ctx context.Context, id uuid.UUID, email string, role domain.Role) error {
			require.Equal(t, adminUUID, id)
			require.Equal(t, adminEmail, email)
			require.Equal(t, domain.RoleAdmin, role)
			return nil
		},
	}
	us := NewUserService(repo, "test-secret")

	token, err := us.DummyLogin(context.Background(), domain.RoleAdmin)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestUserService_Register(t *testing.T) {
	userID := uuid.New()
	repo := &mockRepository{
		CreateUserFn: func(ctx context.Context, email, password string, role domain.Role) (*domain.User, error) {
			require.Equal(t, "test@example.com", email)
			require.Equal(t, domain.RoleUser, role)
			return &domain.User{ID: userID, Email: email, Role: role}, nil
		},
	}
	us := NewUserService(repo, "test-secret")

	user, err := us.Register(context.Background(), "test@example.com", "password123", domain.RoleUser)
	require.NoError(t, err)
	require.Equal(t, userID, user.ID)
	require.Equal(t, "test@example.com", user.Email)
}

func TestUserService_Register_EmailExists(t *testing.T) {
	repo := &mockRepository{
		CreateUserFn: func(ctx context.Context, email, password string, role domain.Role) (*domain.User, error) {
			return nil, postgres.ErrEmailAlreadyExist
		},
	}
	us := NewUserService(repo, "test-secret")

	_, err := us.Register(context.Background(), "test@example.com", "password123", domain.RoleUser)
	require.Error(t, err)
	require.Equal(t, ErrEmailAlreadyExist, err)
}

func TestUserService_Login(t *testing.T) {
	hashedPassword, _ := hashPassword("password123")
	repo := &mockRepository{
		GetUserFn: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:           uuid.New(),
				Email:        email,
				PasswordHash: hashedPassword,
				Role:         domain.RoleUser,
			}, nil
		},
	}
	us := NewUserService(repo, "test-secret")

	token, err := us.Login(context.Background(), "test@example.com", "password123")
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	repo := &mockRepository{
		GetUserFn: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, postgres.ErrUserNotFound
		},
	}
	us := NewUserService(repo, "test-secret")

	_, err := us.Login(context.Background(), "test@example.com", "password123")
	require.Error(t, err)
	require.Equal(t, ErrInvalidCredentials, err)
}
