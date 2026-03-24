package service

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/google/uuid"
)

type RoomService interface {
	GetRooms(ctx context.Context) ([]*domain.Room, error)
	CreateRoom(ctx context.Context, room *domain.Room) error
	CreateSchedule(
		ctx context.Context, roomID uuid.UUID,
		daysOfWeek []domain.DayOfWeek,
		startTime, endTime string,
	) (*domain.Schedule, error)
}

type roomService struct {
	repository domain.Repository
}

func NewRoomService(repo domain.Repository) RoomService {
	return &roomService{repository: repo}
}

func (s *roomService) GetRooms(ctx context.Context) ([]*domain.Room, error) {
	rooms, err := s.repository.GetRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("get rooms: %w", err)
	}

	return rooms, nil
}

func (s *roomService) CreateRoom(ctx context.Context, room *domain.Room) error {
	if err := s.repository.CreateRoom(ctx, room); err != nil {
		return fmt.Errorf("create room: %w", err)
	}
	return nil
}

func (s *roomService) CreateSchedule(
	ctx context.Context, roomID uuid.UUID,
	daysOfWeek []domain.DayOfWeek,
	startTime, endTime string,
) (*domain.Schedule, error) {
	if err := validateDaysOfWeek(daysOfWeek); err != nil {
		return nil, err
	}
	start, end, err := validateAndParseScheduleTime(startTime, endTime)
	if err != nil {
		return nil, err
	}

	schedule := domain.Schedule{
		RoomID:     roomID,
		DaysOfWeek: daysOfWeek,
		StartTime:  start,
		EndTime:    end,
	}

	if err := s.repository.CreateSchedule(ctx, &schedule); err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}

	return &schedule, err

}
