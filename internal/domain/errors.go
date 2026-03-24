package domain

import "errors"

var (
	ErrRoomNotFound   = errors.New("room not found")
	ErrScheduleExists = errors.New("schedule already exists")
	ErrUserNotFound   = errors.New("user not found")
)
