package postgres

import "errors"

var (
	ErrEmailAlreadyExist = errors.New("login alreadyexist")
	ErrUserNotFound      = errors.New("user not found")
)
