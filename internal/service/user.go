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
)

type UserService interface {
	DummyLogin(ctx context.Context, role domain.Role) (string, error)
	Register(ctx context.Context, email, password string, role domain.Role) (*domain.User, error)
	// Login(ctx context.Context, login, password string) (string, error)
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
	var token string
	var err error
	switch role {
	case domain.RoleUser:
		token, err = jwt.GenerateToken(userUUID, domain.RoleUser, s.jwtSecret)
	case domain.RoleAdmin:
		token, err = jwt.GenerateToken(adminUUID, domain.RoleAdmin, s.jwtSecret)
	default:
		return "", ErrInvalidRole
	}
	if err != nil {
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
