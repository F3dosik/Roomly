package service

import (
	"context"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
)

type SlotService interface {
	GetFreeSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error)
}

type slotService struct {
	repository domain.Repository
}

func NewSlotService(repo domain.Repository) SlotService {
	return &slotService{repository: repo}
}

func (s *slotService) GetFreeSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	slots, err := s.repository.GetFreeSlots(ctx, roomID, date)
	if err != nil {
		return nil, err
	}
	return slots, nil
}
