package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"monolith/internal/account"
	"monolith/internal/api"
	"monolith/internal/auth"
	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/logger"
	"monolith/internal/login"
	"monolith/migrations"

	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.NewConfig()

	slog.SetDefault(logger.New(logger.Config{
		Level: cfg.Logging.Level,
	}))

	db, err := database.New(cfg.Database.URL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic("Database connection error")
	}
	defer db.Close()

	migrations.Up(stdlib.OpenDBFromPool(db.PgxPool()))

	accountService := account.NewService(db)
	loginService := login.NewService(db, accountService)
	authService := auth.NewService(db, cfg.Security)

	srv := api.NewHTTPServer(db, cfg, accountService, loginService, authService)
	srv.Setup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startSessionCleanup(ctx, authService, time.Hour)

	if startErr := srv.Start(ctx); startErr != nil && !errors.Is(startErr, http.ErrServerClosed) {
		slog.Error("Server failed to start", "error", startErr)
	}
}

func startSessionCleanup(ctx context.Context, authService *auth.Service, interval time.Duration) {
	go func() {
		if err := authService.CleanupSessions(ctx); err != nil {
			slog.Warn("Failed to cleanup sessions", "error", err)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := authService.CleanupSessions(ctx); err != nil {
					slog.Warn("Failed to cleanup sessions", "error", err)
				}
			}
		}
	}()
}
