package main

import (
	"log/slog"
	"net/http"
	"os"

	"monolith/internal/database"
	"monolith/internal/handlers"
	"monolith/internal/logger"
	customMiddleware "monolith/internal/middleware"
	"monolith/migrations"
	"monolith/web"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())

	// Middleware
	// the order of the middleware is important in most cases

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			slog.Error("[PANIC RECOVER]", "error", err, "stack", string(stack))
			return err
		},
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(web.Assets()),
		HTML5:      true,
	}))

	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	api := e.Group("/api")

	// Public auth routes
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// Protected routes
	protected := api.Group("", customMiddleware.JWTAuth())
	protected.GET("/auth/me", authHandler.Me)
	protected.POST("/auth/logout", authHandler.Logout)
	protected.GET("/users", userHandler.GetUsers)
	protected.GET("/users/:id", userHandler.GetUser)
	protected.POST("/users", userHandler.CreateUser)
	protected.PUT("/users/:id", userHandler.UpdateUser)
	protected.DELETE("/users/:id", userHandler.DeleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Info("Server starting", slog.String("port", port))
	if err := e.Start(":" + port); err != nil {
		log.Error("Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
}
