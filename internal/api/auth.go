package api

import (
	"errors"
	"net/http"

	authService "monolith/internal/service/auth"
	loginService "monolith/internal/service/login"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	authService *authService.Service
}

func NewSessionHandler(authService *authService.Service) *SessionHandler {
	return &SessionHandler{
		authService: authService,
	}
}

func (h *SessionHandler) GetSessions(c echo.Context) error {
	user, ok := c.Get("user").(*authService.AuthUser)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	sessions, err := h.authService.GetUserSessions(c.Request().Context(), user.UserID)
	if err != nil {
		if errors.Is(err, loginService.ErrInvalidUserID) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID").SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve sessions").SetInternal(err)
	}

	for i := range sessions {
		sessions[i].IsCurrent = sessions[i].ID == user.SessionID
	}

	return c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) RevokeSession(c echo.Context) error {
	user, ok := c.Get("user").(*authService.AuthUser)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	sessionIDParam := c.Param("sessionId")

	if sessionIDParam == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Session ID is required")
	}

	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid session ID format").SetInternal(err)
	}

	err = h.authService.RevokeSession(c.Request().Context(), user.UserID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, loginService.ErrInvalidUserID):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID").SetInternal(err)
		case errors.Is(err, authService.ErrSessionNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Session not found").SetInternal(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke session").SetInternal(err)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session revoked successfully"})
}

func (h *SessionHandler) RotateSession(c echo.Context) error {
	sessionTokenCookie, err := c.Cookie("session_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session token cookie not found").SetInternal(err)
	}

	session, err := h.authService.RotateSession(c.Request().Context(), &authService.RotateSessionRequest{
		UnhashedToken: sessionTokenCookie.Value,
		ClientIP:      c.RealIP(),
		UserAgent:     c.Request().UserAgent(),
	})
	if err != nil {
		h.authService.ClearAuthCookies(c)

		switch {
		case errors.Is(err, authService.ErrSessionExpired):
			return echo.NewHTTPError(http.StatusUnauthorized, "Session expired").SetInternal(err)
		case errors.Is(err, loginService.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusUnauthorized, "User not found").SetInternal(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to rotate session").SetInternal(err)
		}
	}

	h.authService.SetSessionCookies(c, session)

	return c.JSON(http.StatusOK, map[string]string{"message": "Session rotated successfully"})
}
