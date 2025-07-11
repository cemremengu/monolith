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
		return c.JSON(http.StatusUnauthorized, APIError{Message: "Invalid user ID"})
	}

	currentSessionID, ok := c.Get("session_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, APIError{Message: "Invalid session ID"})
	}

	sessions, err := h.sessionService.GetUserSessions(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidUserID) {
			return c.JSON(http.StatusBadRequest, APIError{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, APIError{Message: "Failed to retrieve sessions"})
	}

	for i := range sessions {
		sessions[i].IsCurrent = sessions[i].ID == currentSessionID
	}

	return c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) RevokeSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, APIError{Message: "Invalid user ID"})
	}
	sessionIDParam := c.Param("sessionId")

	if sessionIDParam == "" {
		return c.JSON(http.StatusBadRequest, APIError{Message: "Session ID is required"})
	}

	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Message: "Invalid session ID format"})
	}

	err = h.sessionService.RevokeSession(c.Request().Context(), userID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidUserID):
			return c.JSON(http.StatusBadRequest, APIError{Message: err.Error()})
		case errors.Is(err, sessionService.ErrSessionNotFound):
			return c.JSON(http.StatusNotFound, APIError{Message: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, APIError{Message: "Failed to revoke session"})
		}
	}

	return c.JSON(http.StatusOK, APIError{Message: "Session revoked successfully"})
}

func (h *SessionHandler) RotateSession(c echo.Context) error {
	sessionTokenCookie, err := c.Cookie("session_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, APIError{Message: "Authentication required"})
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
			return c.JSON(http.StatusUnauthorized, APIError{Message: err.Error()})
		case errors.Is(err, authService.ErrUserNotFound):
			return c.JSON(http.StatusUnauthorized, APIError{Message: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, APIError{Message: "Failed to rotate session"})
		}
	}

	h.sessionService.SetSessionCookies(c, session)

	return c.JSON(http.StatusOK, APIError{Message: "Session rotated successfully"})
}
