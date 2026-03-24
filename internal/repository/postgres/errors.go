package postgres

import "errors"

var (
	ErrEmailAlreadyExist = errors.New("login already exist")
	ErrUserNotFound      = errors.New("user not found")
	ErrScheduleExists    = errors.New("schedule already exist")
	ErrRoomNotFound      = errors.New("room not found")
)
