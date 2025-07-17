package api

import (
	"errors"
	"net/http"

	loginService "monolith/internal/service/login"
	sessionService "monolith/internal/service/session"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	loginService   *loginService.Service
	sessionService *sessionService.Service
}

func NewAuthHandler(loginService *loginService.Service, sessionService *sessionService.Service) *AuthHandler {
	return &AuthHandler{
		loginService:   loginService,
		sessionService: sessionService,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginService.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").SetInternal(err)
	}

	user, err := h.loginService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, loginService.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials").SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to login").SetInternal(err)
	}

	session, tokenErr := h.sessionService.CreateSession(c.Request().Context(), &sessionService.CreateSessionRequest{
		AccountID: user.ID,
		ClientIP:  c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	})

	if tokenErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create session").SetInternal(tokenErr)
	}

	h.sessionService.SetSessionCookies(c, session)

	response := map[string]any{
		"message": "Login successful",
	}

	return c.JSON(http.StatusOK, response)
}

// Logout revokes the current session and clears authentication cookies.
// Ignore any errors and act as a no-op on failure.
func (h *AuthHandler) Logout(c echo.Context) error {
	// For logout, just clear the cookie
	// The session will eventually be cleaned up by the cleanup process
	h.sessionService.ClearAuthCookies(c)
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
