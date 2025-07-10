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
	"monolith/internal/database"
	"monolith/internal/logger"
	"monolith/migrations"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	_ "go.uber.org/automaxprocs"
)

const shutdownTimeout = 10

func main() {
	if err := godotenv.Load(); err != nil {
		log.Debug("No .env file found or error loading .env file", "error", err)
	}

	log := logger.New(logger.Config{
		Level: logger.GetLevel(os.Getenv("LOG_LEVEL")),
		Env:   os.Getenv("ENV"),
	})

	slog.SetDefault(log)

	db, err := database.New()
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic("Database connection error")
	}
	defer db.Close()

	migrations.Up(stdlib.OpenDBFromPool(db.Pool))

	srv := api.NewHTTPServer(db, log)
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
