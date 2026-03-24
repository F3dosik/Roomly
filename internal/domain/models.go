package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string
type DayOfWeek int
type BookingStatus string

const (
	Monday    DayOfWeek = 1
	Tuesday   DayOfWeek = 2
	Wednesday DayOfWeek = 3
	Thursday  DayOfWeek = 4
	Friday    DayOfWeek = 5
	Saturday  DayOfWeek = 6
	Sunday    DayOfWeek = 7
)

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
	CreatedAt   *time.Time
}

type Schedule struct {
	ID         uuid.UUID
	RoomID     uuid.UUID
	DaysOfWeek []DayOfWeek
	StartTime  time.Time
	EndTime    time.Time
	CreatedAt  time.Time
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
