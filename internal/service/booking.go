package service

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
)

type BookingService interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, slotID uuid.UUID, createConferenceLink bool) (*domain.Booking, error)
	ListBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error)
	GetMyBookings(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error)
	CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error)
}

type bookingService struct {
	repository domain.Repository
}

func NewBookingService(repo domain.Repository) BookingService {
	return &bookingService{repository: repo}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID, slotID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
	slot, err := s.repository.GetSlotByID(ctx, slotID)
	if err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	if slot.StartsAt.Before(time.Now().UTC()) {
		return nil, domain.ErrBookingInPast
	}

	booking := &domain.Booking{
		UserID: userID,
		SlotID: slotID,
		Status: domain.BookingStatusActive,
	}

	if createConferenceLink {
		link := generateConferenceLink()
		booking.ConferenceLink = &link
	}

	if err := s.repository.CreateBooking(ctx, booking); err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	return booking, nil
}

func generateConferenceLink() string {
	return fmt.Sprintf("https://meet.example.com/%s", uuid.New().String())
}

func (s *bookingService) ListBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
	bookings, total, err := s.repository.ListBookings(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}
	return bookings, total, nil
}

func (s *bookingService) GetMyBookings(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
	bookings, err := s.repository.GetMyBookings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get my bookings: %w", err)
	}
	return bookings, nil
}

func (s *bookingService) CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error) {
	booking, err := s.repository.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("cancel booking: %w", err)
	}

	if booking.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if booking.Status == domain.BookingStatusCancelled {
		return booking, nil
	}

	booking, err = s.repository.CancelBooking(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("cancel booking: %w", err)
	}

	return booking, nil
}
