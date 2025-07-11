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
		return APIError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid user ID",
		}
	}

	currentSessionID, ok := c.Get("session_id").(uuid.UUID)
	if !ok {
		return APIError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid session ID",
		}
	}

	sessions, err := h.sessionService.GetUserSessions(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidUserID) {
			return APIError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Err:     err,
			}
		}
		return APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve sessions",
			Err:     err,
		}
	}

	for i := range sessions {
		sessions[i].IsCurrent = sessions[i].ID == currentSessionID
	}

	return c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) RevokeSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return APIError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid user ID",
		}
	}
	sessionIDParam := c.Param("sessionId")

	if sessionIDParam == "" {
		return APIError{
			Code:    http.StatusBadRequest,
			Message: "Session ID is required",
		}
	}

	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		return APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid session ID format",
			Err:     err,
		}
	}

	err = h.sessionService.RevokeSession(c.Request().Context(), userID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidUserID):
			return APIError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Err:     err,
			}
		case errors.Is(err, sessionService.ErrSessionNotFound):
			return APIError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
				Err:     err,
			}
		default:
			return APIError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to revoke session",
				Err:     err,
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session revoked successfully"})
}

func (h *SessionHandler) RotateSession(c echo.Context) error {
	sessionTokenCookie, err := c.Cookie("session_token")
	if err != nil {
		return APIError{
			Code:    http.StatusUnauthorized,
			Message: "Authentication required",
			Err:     err,
		}
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
			return APIError{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
				Err:     err,
			}
		case errors.Is(err, authService.ErrUserNotFound):
			return APIError{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
				Err:     err,
			}
		default:
			return APIError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to rotate session",
				Err:     err,
			}
		}
	}

	h.sessionService.SetSessionCookies(c, session)

	return c.JSON(http.StatusOK, map[string]string{"message": "Session rotated successfully"})
}
