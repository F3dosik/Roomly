package service

import (
	"context"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBookingService_CreateBooking(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()
	bookingID := uuid.New()
	futureTime := time.Now().UTC().Add(24 * time.Hour)

	repo := &mockRepository{
		GetSlotByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Slot, error) {
			require.Equal(t, slotID, id)
			return &domain.Slot{
				ID:       slotID,
				RoomID:   uuid.New(),
				StartsAt: futureTime,
				EndsAt:   futureTime.Add(time.Hour),
			}, nil
		},
		CreateBookingFn: func(ctx context.Context, booking *domain.Booking) error {
			booking.ID = bookingID
			require.Equal(t, userID, booking.UserID)
			require.Equal(t, slotID, booking.SlotID)
			require.Equal(t, domain.BookingStatusActive, booking.Status)
			return nil
		},
	}
	bs := NewBookingService(repo)

	booking, err := bs.CreateBooking(context.Background(), userID, slotID, false)
	require.NoError(t, err)
	require.Equal(t, bookingID, booking.ID)
	require.Equal(t, userID, booking.UserID)
	require.Nil(t, booking.ConferenceLink)
}

func TestBookingService_CreateBooking_WithConferenceLink(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()
	bookingID := uuid.New()
	futureTime := time.Now().UTC().Add(24 * time.Hour)

	repo := &mockRepository{
		GetSlotByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Slot, error) {
			return &domain.Slot{
				ID:       slotID,
				RoomID:   uuid.New(),
				StartsAt: futureTime,
				EndsAt:   futureTime.Add(time.Hour),
			}, nil
		},
		CreateBookingFn: func(ctx context.Context, booking *domain.Booking) error {
			booking.ID = bookingID
			require.NotNil(t, booking.ConferenceLink)
			return nil
		},
	}
	bs := NewBookingService(repo)

	booking, err := bs.CreateBooking(context.Background(), userID, slotID, true)
	require.NoError(t, err)
	require.NotNil(t, booking.ConferenceLink)
}

func TestBookingService_CreateBooking_InPast(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()
	pastTime := time.Now().UTC().Add(-1 * time.Hour)

	repo := &mockRepository{
		GetSlotByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Slot, error) {
			return &domain.Slot{
				ID:       slotID,
				RoomID:   uuid.New(),
				StartsAt: pastTime,
				EndsAt:   pastTime.Add(time.Hour),
			}, nil
		},
	}
	bs := NewBookingService(repo)

	_, err := bs.CreateBooking(context.Background(), userID, slotID, false)
	require.Error(t, err)
	require.Equal(t, domain.ErrBookingInPast, err)
}

func TestBookingService_ListBookings(t *testing.T) {
	bookingID := uuid.New()
	userID := uuid.New()

	repo := &mockRepository{
		ListBookingsFn: func(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
			require.Equal(t, 1, page)
			require.Equal(t, 10, pageSize)
			return []*domain.Booking{
				{
					ID:     bookingID,
					UserID: userID,
					Status: domain.BookingStatusActive,
				},
			}, 1, nil
		},
	}
	bs := NewBookingService(repo)

	bookings, total, err := bs.ListBookings(context.Background(), 1, 10)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, bookings, 1)
	require.Equal(t, bookingID, bookings[0].ID)
}

func TestBookingService_GetMyBookings(t *testing.T) {
	userID := uuid.New()
	bookingID := uuid.New()

	repo := &mockRepository{
		GetMyBookingsFn: func(ctx context.Context, id uuid.UUID) ([]*domain.Booking, error) {
			require.Equal(t, userID, id)
			return []*domain.Booking{
				{
					ID:     bookingID,
					UserID: userID,
					Status: domain.BookingStatusActive,
				},
			}, nil
		},
	}
	bs := NewBookingService(repo)

	bookings, err := bs.GetMyBookings(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, bookings, 1)
	require.Equal(t, bookingID, bookings[0].ID)
}

func TestBookingService_CancelBooking(t *testing.T) {
	userID := uuid.New()
	bookingID := uuid.New()

	repo := &mockRepository{
		GetBookingByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
			require.Equal(t, bookingID, id)
			return &domain.Booking{
				ID:     bookingID,
				UserID: userID,
				Status: domain.BookingStatusActive,
			}, nil
		},
		CancelBookingFn: func(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
			return &domain.Booking{
				ID:     bookingID,
				UserID: userID,
				Status: domain.BookingStatusCancelled,
			}, nil
		},
	}
	bs := NewBookingService(repo)

	booking, err := bs.CancelBooking(context.Background(), userID, bookingID)
	require.NoError(t, err)
	require.Equal(t, domain.BookingStatusCancelled, booking.Status)
}

func TestBookingService_CancelBooking_Forbidden(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	bookingID := uuid.New()

	repo := &mockRepository{
		GetBookingByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
			return &domain.Booking{
				ID:     bookingID,
				UserID: otherUserID,
				Status: domain.BookingStatusActive,
			}, nil
		},
	}
	bs := NewBookingService(repo)

	_, err := bs.CancelBooking(context.Background(), userID, bookingID)
	require.Error(t, err)
	require.Equal(t, domain.ErrForbidden, err)
}

func TestBookingService_CancelBooking_AlreadyCancelled(t *testing.T) {
	userID := uuid.New()
	bookingID := uuid.New()

	repo := &mockRepository{
		GetBookingByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
			return &domain.Booking{
				ID:     bookingID,
				UserID: userID,
				Status: domain.BookingStatusCancelled,
			}, nil
		},
	}
	bs := NewBookingService(repo)

	booking, err := bs.CancelBooking(context.Background(), userID, bookingID)
	require.NoError(t, err)
	require.Equal(t, domain.BookingStatusCancelled, booking.Status)
}
