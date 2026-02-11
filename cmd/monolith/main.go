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

	"monolith/internal/api"
	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/logger"
	"monolith/internal/service/account"
	"monolith/internal/service/auth"
	"monolith/internal/service/login"
	"monolith/migrations"

	"github.com/jackc/pgx/v5/stdlib"
)

const shutdownTimeout = 10

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

	go func() {
		if startErr := srv.Start(); startErr != nil && !errors.Is(startErr, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", startErr)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()
	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		slog.Error("Server shutdown failed", "error", shutdownErr)
	}
}
