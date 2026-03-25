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
	"github.com/avito-internships/test-backend-1-F3dosik/internal/scheduler"
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

	// воркер генерации слотов
	gen := scheduler.New(repo, logger)
	go gen.Run(ctx)

	// сервисы
	us := service.NewUserService(repo, cfg.JWTSecret)
	rs := service.NewRoomService(repo, gen)
	ss := service.NewSlotService(repo)
	bs := service.NewBookingService(repo)

	// хендлеры
	h := handler.New(cfg.JWTSecret, us, rs, ss, bs, logger)

	// сервер
	srv := server.New(cfg, h, logger)
	srv.Run(ctx)
}
