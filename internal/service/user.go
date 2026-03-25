package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/jwt"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/repository/postgres"
	"github.com/google/uuid"
)

var (
	adminUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userUUID  = uuid.MustParse("00000000-0000-0000-0000-000000000002")

	adminEmail = "admin@dummy.local"
	userEmail  = "user@dummy.local"
)

type UserService interface {
	DummyLogin(ctx context.Context, role domain.Role) (string, error)
	Register(ctx context.Context, email, password string, role domain.Role) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	repository domain.Repository
	jwtSecret  string
}

func NewUserService(repo domain.Repository, jwtSecret string) UserService {
	return &userService{
		repository: repo,
		jwtSecret:  jwtSecret,
	}
}

func (s *userService) DummyLogin(ctx context.Context, role domain.Role) (string, error) {
	if err := validateRole(role); err != nil {
		return "", err
	}
	var id uuid.UUID
	var email string
	if role == domain.RoleAdmin {
		id = adminUUID
		email = adminEmail
	} else {
		id = userUUID
		email = userEmail
	}
	token, err := jwt.GenerateToken(id, role, s.jwtSecret)
	if err != nil {
		return "", err
	}

	if err := s.repository.UpsertUser(ctx, id, email, role); err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) Register(ctx context.Context, email, password string, role domain.Role) (*domain.User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	if err := validateRole(role); err != nil {
		return nil, err
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.CreateUser(ctx, email, hash, role)
	if err != nil {
		if errors.Is(err, postgres.ErrEmailAlreadyExist) {
			return nil, ErrEmailAlreadyExist
		}
		return nil, fmt.Errorf("register: %w", err)
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repository.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("login: %w", err)
	}

	if !checkPassword(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil

}
