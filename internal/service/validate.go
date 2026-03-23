package service

import (
	"fmt"
	"net/mail"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
)

func validateEmail(email string) error {
	if len(email) == 0 {
		return ErrEmptyEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidEmail, err)
	}
	return nil
}

func validateRole(role domain.Role) error {
	switch role {
	case domain.RoleAdmin, domain.RoleUser:
		return nil
	default:
		return ErrInvalidRole

	}
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}
