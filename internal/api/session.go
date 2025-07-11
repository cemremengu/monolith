package api

import (
	"errors"
	"net/http"

	authService "monolith/internal/service/auth"
	sessionService "monolith/internal/service/session"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	sessionService *sessionService.Service
}

func NewSessionHandler(sessionService *sessionService.Service) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

func (h *SessionHandler) GetSessions(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	currentSessionID, ok := c.Get("session_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid session ID"})
	}

	sessions, err := h.sessionService.GetUserSessions(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidUserID) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch sessions"})
	}

	for i := range sessions {
		sessions[i].IsCurrent = sessions[i].ID == currentSessionID
	}

	return c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) RevokeSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}
	sessionIDParam := c.Param("sessionId")

	if sessionIDParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Session ID is required"})
	}

	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid session ID format"})
	}

	err = h.sessionService.RevokeSession(c.Request().Context(), userID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidUserID):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, sessionService.ErrSessionNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to revoke session"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session revoked successfully"})
}

func (h *SessionHandler) RotateToken(c echo.Context) error {
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
