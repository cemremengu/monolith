package api

import (
	"errors"
	"net/http"

	"monolith/internal/service/account"
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login"})
	}

	session, tokenErr := h.sessionService.CreateSession(c.Request().Context(), &sessionService.CreateSessionRequest{
		AccountID: user.ID,
		ClientIP:  c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	})

	if tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate session token"})
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

func (h *AuthHandler) RotateToken(c echo.Context) error {
	sessionTokenCookie, err := c.Cookie("session_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session token not found"})
	}

	session, err := h.sessionService.RotateSessionToken(c.Request().Context(), &sessionService.RotateSessionTokenRequest{
		UnhashedToken: sessionTokenCookie.Value,
		ClientIP:      c.RealIP(),
		UserAgent:     c.Request().UserAgent(),
	})
	if err != nil {
		// nuke cookies on any refresh error
		h.sessionService.ClearAuthCookies(c)

		switch {
		case errors.Is(err, sessionService.ErrSessionExpired):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		case errors.Is(err, authService.ErrUserNotFound):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to refresh session"})
		}
	}

	h.sessionService.SetSessionCookies(c, session)

	return c.JSON(http.StatusOK, map[string]string{"message": "Session refreshed successfully"})
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req account.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	account, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrPasswordTooShort):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, authService.ErrUserAlreadyExists):
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create account"})
		}
	}

	response := map[string]any{
		"account": account,
	}

	return c.JSON(http.StatusCreated, response)
}
