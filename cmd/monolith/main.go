package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"monolith/internal/api"
	"monolith/internal/config"
	"monolith/internal/database"
	"monolith/internal/logger"
	"monolith/internal/service/account"
	"monolith/internal/service/auth"
	"monolith/internal/service/user"
	"monolith/migrations"

	"github.com/jackc/pgx/v5/stdlib"
	_ "go.uber.org/automaxprocs"
)

const shutdownTimeout = 10

func main() {
	cfg := config.NewConfig()

	log := logger.New(logger.Config{
		Level: cfg.Logging.Level,
		Env:   cfg.Logging.Env,
	})

	slog.SetDefault(log)

	db, err := database.New(cfg.Database.URL)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic("Database connection error")
	}
	defer db.Close()

	migrations.Up(stdlib.OpenDBFromPool(db.Pool))

	userService := user.NewService(db)
	accountService := account.NewService(db)
	authService := auth.NewService(db)

	srv := api.NewHTTPServer(db, log, cfg, userService, accountService, authService)
	srv.Setup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if startErr := srv.Start(); startErr != nil && errors.Is(startErr, http.ErrServerClosed) {
			log.Error("Server failed to start", "error", startErr)
			panic("Server startup error")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()
	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		panic(shutdownErr)
	}
}
