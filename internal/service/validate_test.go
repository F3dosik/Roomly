package service

import (
	"errors"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/stretchr/testify/require"
)

// validateEmail tests
func TestValidateEmail_Valid(t *testing.T) {
	err := validateEmail("test@example.com")
	require.NoError(t, err)
}

func TestValidateEmail_Empty(t *testing.T) {
	err := validateEmail("")
	require.Error(t, err)
	require.Equal(t, ErrEmptyEmail, err)
}

func TestValidateEmail_Invalid(t *testing.T) {
	err := validateEmail("not-an-email")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidEmail))
}

// validatePassword tests
func TestValidatePassword_Valid(t *testing.T) {
	err := validatePassword("password123")
	require.NoError(t, err)
}

func TestValidatePassword_TooShort(t *testing.T) {
	err := validatePassword("short")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooShort, err)
}

func TestValidatePassword_MinLength(t *testing.T) {
	err := validatePassword("12345678")
	require.NoError(t, err)
}

// validateRole tests
func TestValidateRole_Admin(t *testing.T) {
	err := validateRole(domain.RoleAdmin)
	require.NoError(t, err)
}

func TestValidateRole_User(t *testing.T) {
	err := validateRole(domain.RoleUser)
	require.NoError(t, err)
}

func TestValidateRole_Invalid(t *testing.T) {
	err := validateRole("invalid_role")
	require.Error(t, err)
	require.Equal(t, ErrInvalidRole, err)
}

// validateDaysOfWeek tests
func TestValidateDaysOfWeek_Valid(t *testing.T) {
	err := validateDaysOfWeek([]domain.DayOfWeek{1, 2, 3})
	require.NoError(t, err)
}

func TestValidateDaysOfWeek_Empty(t *testing.T) {
	err := validateDaysOfWeek([]domain.DayOfWeek{})
	require.Error(t, err)
	require.Equal(t, ErrEmptyDaysOfWeek, err)
}

func TestValidateDaysOfWeek_TooLow(t *testing.T) {
	err := validateDaysOfWeek([]domain.DayOfWeek{0})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidDaysOfWeek))
}

func TestValidateDaysOfWeek_TooHigh(t *testing.T) {
	err := validateDaysOfWeek([]domain.DayOfWeek{8})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidDaysOfWeek))
}

func TestValidateDaysOfWeek_AllDays(t *testing.T) {
	err := validateDaysOfWeek([]domain.DayOfWeek{1, 2, 3, 4, 5, 6, 7})
	require.NoError(t, err)
}

// validateAndParseScheduleTime tests
func TestValidateAndParseScheduleTime_Valid(t *testing.T) {
	start, end, err := validateAndParseScheduleTime("09:00", "18:00")
	require.NoError(t, err)
	require.Equal(t, 9, start.Hour())
	require.Equal(t, 0, start.Minute())
	require.Equal(t, 18, end.Hour())
	require.Equal(t, 0, end.Minute())
}

func TestValidateAndParseScheduleTime_InvalidStartFormat(t *testing.T) {
	_, _, err := validateAndParseScheduleTime("25:00", "18:00")
	require.Error(t, err)
	require.Equal(t, ErrInvalidTimeFormat, err)
}

func TestValidateAndParseScheduleTime_InvalidEndFormat(t *testing.T) {
	_, _, err := validateAndParseScheduleTime("09:00", "25:00")
	require.Error(t, err)
	require.Equal(t, ErrInvalidTimeFormat, err)
}

func TestValidateAndParseScheduleTime_StartAfterEnd(t *testing.T) {
	_, _, err := validateAndParseScheduleTime("18:00", "09:00")
	require.Error(t, err)
	require.Equal(t, ErrInvalidTimeRange, err)
}

func TestValidateAndParseScheduleTime_SameTime(t *testing.T) {
	_, _, err := validateAndParseScheduleTime("09:00", "09:00")
	require.Error(t, err)
	require.Equal(t, ErrInvalidTimeRange, err)
}

func TestValidateAndParseScheduleTime_ValidMinutes(t *testing.T) {
	start, end, err := validateAndParseScheduleTime("09:30", "18:45")
	require.NoError(t, err)
	require.Equal(t, 9, start.Hour())
	require.Equal(t, 30, start.Minute())
	require.Equal(t, 18, end.Hour())
	require.Equal(t, 45, end.Minute())
}

// isHalfHour tests
func TestIsHalfHour_OnTheHour(t *testing.T) {
	time := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	require.True(t, isHalfHour(time))
}

func TestIsHalfHour_HalfPast(t *testing.T) {
	time := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
	require.True(t, isHalfHour(time))
}

func TestIsHalfHour_Quarter(t *testing.T) {
	time := time.Date(2024, 1, 1, 9, 15, 0, 0, time.UTC)
	require.False(t, isHalfHour(time))
}

func TestIsHalfHour_45Minutes(t *testing.T) {
	time := time.Date(2024, 1, 1, 9, 45, 0, 0, time.UTC)
	require.False(t, isHalfHour(time))
}
