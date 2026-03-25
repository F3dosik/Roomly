package service

import (
	"errors"
)

var (
	ErrEmptyEmail         = errors.New("login is empty")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrPasswordTooShort   = errors.New("password too short")
	ErrEmailAlreadyExist  = errors.New("login already exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidRole        = errors.New("invalid role")
	ErrEmptyDaysOfWeek    = errors.New("daysOfWeek cannot be empty")
	ErrInvalidDaysOfWeek  = errors.New("invalid dayOfWeek")
	ErrInvalidTimeFormat  = errors.New("invalid time format")
	ErrInvalidTimeRange   = errors.New("startTime must be before endTime")
)
