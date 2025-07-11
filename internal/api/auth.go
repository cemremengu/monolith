package api

import (
	"errors"
	"net/http"

	authService "monolith/internal/service/auth"
	sessionService "monolith/internal/service/session"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService    *authService.Service
	sessionService *sessionService.Service
}

func NewAuthHandler(authService *authService.Service, sessionService *sessionService.Service) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionService: sessionService,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req authService.LoginRequest
	if err := c.Bind(&req); err != nil {
		return APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Err:     err,
		}
	}

	user, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidCredentials) {
			return APIError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid email or password",
				Err:     err,
			}
		}
		return APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to login",
			Err:     err,
		}
	}

	session, tokenErr := h.sessionService.CreateSession(c.Request().Context(), &sessionService.CreateSessionRequest{
		AccountID: user.ID,
		ClientIP:  c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	})

	if tokenErr != nil {
		return APIError{
			Code: http.StatusInternalServerError,

			Message: "Failed to create session",
			Err:     tokenErr,
		}
	}

	h.sessionService.SetSessionCookies(c, session)

	response := map[string]any{
		"user": user,
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
