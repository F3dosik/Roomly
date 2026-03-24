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

func (s *userService) generateToken(userID uuid.UUID, role domain.Role) (string, error) {
	var token string
	var err error
	switch role {
	case domain.RoleUser:
		token, err = jwt.GenerateToken(userID, domain.RoleUser, s.jwtSecret)
	case domain.RoleAdmin:
		token, err = jwt.GenerateToken(userID, domain.RoleAdmin, s.jwtSecret)
	default:
		return "", ErrInvalidRole
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) DummyLogin(ctx context.Context, role domain.Role) (string, error) {
	var token string
	var err error
	if role == domain.RoleAdmin {
		token, err = s.generateToken(adminUUID, domain.RoleAdmin)
	} else {
		token, err = s.generateToken(userUUID, domain.RoleUser)
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

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	if err := validateEmail(email); err != nil {
		return "", err
	}

	if err := validatePassword(password); err != nil {
		return "", err
	}

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

	token, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil

}
