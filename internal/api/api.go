package api

import (
	"net/http"

	"monolith"
	mw "monolith/internal/middleware"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes configures all the application routes.
func (hs *HTTPServer) RegisterRoutes() {
	userHandler := NewUserHandler(hs.userService)
	authHandler := NewAuthHandler(hs.loginService, hs.authService)
	accountHandler := NewAccountHandler(hs.accountService)
	authSessionHandler := NewSessionHandler(hs.authService)

	api := hs.echo.Group("/api")

	// Public auth routes
	api.POST("/login", authHandler.Login)
	api.POST("/logout", authHandler.Logout)
	api.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, monolith.GetVersionInfo())
	})
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Protected routes
	protected := api.Group("", mw.SessionAuth(hs.authService, hs.accountService, hs.config.Security))

	protected.GET("/account/sessions", authSessionHandler.GetSessions)
	protected.DELETE("/account/sessions/:sessionId", authSessionHandler.RevokeSession)
	protected.POST("/account/sessions/rotate", authSessionHandler.RotateSession)

	protected.GET("/account/profile", accountHandler.Profile)
	protected.PATCH("/account/preferences", accountHandler.UpdatePreferences)
	protected.POST("/account/register", accountHandler.Register)

	// Admin-only routes
	admin := protected.Group("", mw.AdminOnly())
	admin.GET("/users", userHandler.GetUsers)
	admin.GET("/users/:id", userHandler.GetUser)
	admin.POST("/users", userHandler.CreateUser)
	admin.PUT("/users/:id", userHandler.UpdateUser)
	admin.DELETE("/users/:id", userHandler.DeleteUser)
}
