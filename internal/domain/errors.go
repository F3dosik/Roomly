package domain

import "errors"

var (
	ErrRoomNotFound      = errors.New("room not found")
	ErrScheduleExists    = errors.New("schedule already exists")
	ErrSlotNotFound      = errors.New("slot not found")
	ErrSlotAlreadyBooked = errors.New("slot is already booked")
	ErrBookingInPast     = errors.New("cannot book slot in the past")
	ErrInternal          = errors.New("internal error")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrForbidden         = errors.New("forbidden")
)
