package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string
type Weekday uint8
type BookingStatus string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

const (
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

type Room struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Capacity    *int
	CreatedAt   time.Time
}

type Schedule struct {
	RoomID    uuid.UUID
	Weekdays  []Weekday
	StartTime time.Time
	EndTime   time.Time
	CreatedAt time.Time
}

type Slot struct {
	ID       uuid.UUID
	RoomID   uuid.UUID
	StartsAt time.Time
	EndsAt   time.Time
}

type Booking struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SlotID         uuid.UUID
	Status         BookingStatus
	ConferenceLink *string
	CreatedAt      time.Time
	CancelledAt    *time.Time
}
