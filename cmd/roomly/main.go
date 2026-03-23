package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/config"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/db"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/handler"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/logger"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/repository/postgres"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/server"
	"github.com/avito-internships/test-backend-1-F3dosik/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	mode := logger.Mode(cfg.LogLevel)
	logger := logger.NewLogger(mode)
	defer func() { _ = logger.Sync() }()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer pool.Close()

	// репозитории
	repo := postgres.New(pool)

	// сервисы
	us := service.NewUserService(repo, cfg.JWTSecret)
	// roomService := service.NewRoomService(repo)
	// slotService := service.NewSlotService(repo)
	// bookingService := service.NewBookingService(repo)
	// authService := service.NewAuthService(repo, cfg.JWTSecret)

	// хендлеры
	h := handler.New(cfg.JWTSecret, us, logger)

	// // планировщик генерации слотов
	// scheduler := scheduler.New(repo)
	// go scheduler.Run(ctx)

	// сервер
	srv := server.New(cfg, h, logger)
	srv.Run(ctx)
}
