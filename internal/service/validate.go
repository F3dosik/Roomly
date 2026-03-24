package service

import (
	"fmt"
	"net/mail"
	"time"

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

func validateDaysOfWeek(days []domain.DayOfWeek) error {
	if len(days) == 0 {
		return ErrEmptyDaysOfWeek
	}

	for _, d := range days {
		if d < 1 || d > 7 {
			return fmt.Errorf("%w: %d", ErrInvalidDaysOfWeek, d)
		}
	}

	return nil
}

func validateAndParseScheduleTime(start, end string) (time.Time, time.Time, error) {
	startT, err := time.Parse("15:04", start)
	if err != nil {
		return time.Time{}, time.Time{}, ErrInvalidTimeFormat
	}

	endT, err := time.Parse("15:04", end)
	if err != nil {
		return time.Time{}, time.Time{}, ErrInvalidTimeFormat
	}

	if !startT.Before(endT) {
		return time.Time{}, time.Time{}, ErrInvalidTimeRange
	}

	return startT, endT, nil
}

func isHalfHour(t time.Time) bool {
	min := t.Minute()
	return min == 0 || min == 30
}
