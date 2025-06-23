package main

import (
	"log/slog"
	"os"

	"monolith/internal/database"
	"monolith/internal/logger"
	"monolith/internal/server"
	"monolith/migrations"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Debug("No .env file found or error loading .env file", slog.Any("error", err))
	}

	log := logger.New(logger.Config{
		Level: logger.GetLevel(os.Getenv("LOG_LEVEL")),
		Env:   os.Getenv("ENV"),
	})

	slog.SetDefault(log)

	db, err := database.New()
	if err != nil {
		log.Error("Failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	migrations.Up(stdlib.OpenDBFromPool(db.Pool))

	srv := server.New(db, log)
	srv.Setup()

	if err := srv.Start(); err != nil {
		log.Error("Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
}
