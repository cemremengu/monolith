package main

import (
	"log/slog"
	"os"

	"monolith/internal/database"
	"monolith/internal/handlers"
	"monolith/internal/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
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

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	userHandler := handlers.NewUserHandler(db)

	api := e.Group("/api")
	api.GET("/users", userHandler.GetUsers)
	api.GET("/users/:id", userHandler.GetUser)
	api.POST("/users", userHandler.CreateUser)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("Server starting", slog.String("port", port))
	if err := e.Start(":" + port); err != nil {
		log.Error("Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
}
