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
	log := logger.New(logger.Config{
		Level: logger.GetLevel(os.Getenv("LOG_LEVEL")),
		Env:   os.Getenv("ENV"),
	})

	if err := godotenv.Load(); err != nil {
		log.Debug("No .env file found or error loading .env file", "error", err)
	}

	slog.SetDefault(log)

	db, err := database.New()
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic("Database connection error")
	}
	defer db.Close()

	migrations.Up(stdlib.OpenDBFromPool(db.Pool))

	srv := server.New(db, log)
	srv.Setup()

	if startErr := srv.Start(); startErr != nil {
		log.Error("Server failed to start", "error", startErr)
		panic("Server startup error")
	}
}
