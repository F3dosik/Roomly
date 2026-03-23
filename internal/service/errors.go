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
)
