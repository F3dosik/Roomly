package scheduler

import (
	"context"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"go.uber.org/zap"
)

const (
	slotDuration   = 30 * time.Minute
	windowDays     = 7
	tickerInterval = 1 * time.Hour
)

type SlotGenerator struct {
	repo    domain.Repository
	logger  *zap.SugaredLogger
	trigger chan struct{}
}

func New(repo domain.Repository, logger *zap.SugaredLogger) *SlotGenerator {
	return &SlotGenerator{
		repo:    repo,
		logger:  logger,
		trigger: make(chan struct{}, 1),
	}
}

// Trigger — вызывается после создания расписания
func (g *SlotGenerator) Trigger() {
	select {
	case g.trigger <- struct{}{}:
	default: // уже есть сигнал в канале, пропускаем
	}
}

func (g *SlotGenerator) Run(ctx context.Context) {
	g.generate(ctx)

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-g.trigger:
			g.generate(ctx)
		case <-ticker.C:
			g.generate(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (g *SlotGenerator) generate(ctx context.Context) {
	schedules, err := g.repo.GetAllSchedules(ctx)
	if err != nil {
		g.logger.Errorw("slot generator: get schedules", "error", err)
		return
	}

	g.logger.Infow("slot generator: found schedules", "count", len(schedules))

	if len(schedules) == 0 {
		g.logger.Info("slot generator: no schedules found, skipping")
		return
	}

	now := time.Now().UTC()
	windowEnd := now.AddDate(0, 0, windowDays)

	for _, schedule := range schedules {
		if err := g.generateForSchedule(ctx, schedule, now, windowEnd); err != nil {
			g.logger.Errorw("slot generator: generate for schedule",
				"room_id", schedule.RoomID,
				"error", err,
			)
		}
	}
}

func (g *SlotGenerator) generateForSchedule(ctx context.Context, schedule *domain.Schedule, from, to time.Time) error {
	lastDate, err := g.repo.GetLastSlotDate(ctx, schedule.RoomID)
	if err != nil {
		return err
	}

	startDate := from.Truncate(24 * time.Hour)
	if lastDate != nil && lastDate.After(startDate) {
		startDate = lastDate.Truncate(24*time.Hour).AddDate(0, 0, 1)
	}

	var slots []*domain.Slot
	for d := startDate; d.Before(to); d = d.AddDate(0, 0, 1) {
		weekday := int(d.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		if !containsDay(schedule.DaysOfWeek, domain.DayOfWeek(weekday)) {
			continue
		}

		slots = append(slots, generateDaySlots(schedule, d)...)
	}

	if len(slots) == 0 {
		return nil
	}

	return g.repo.GenerateSlots(ctx, slots)
}

func generateDaySlots(schedule *domain.Schedule, date time.Time) []*domain.Slot {
	var slots []*domain.Slot

	startTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		schedule.StartTime.Hour(), schedule.StartTime.Minute(), 0, 0, time.UTC,
	)
	endTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		schedule.EndTime.Hour(), schedule.EndTime.Minute(), 0, 0, time.UTC,
	)

	for t := startTime; t.Add(slotDuration).Before(endTime) || t.Add(slotDuration).Equal(endTime); t = t.Add(slotDuration) {
		slots = append(slots, &domain.Slot{
			RoomID:   schedule.RoomID,
			StartsAt: t,
			EndsAt:   t.Add(slotDuration),
		})
	}

	return slots
}

func containsDay(days []domain.DayOfWeek, day domain.DayOfWeek) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}
