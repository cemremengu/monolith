package api

import (
	"errors"
	"net/http"

	"monolith/internal/database"
	authService "monolith/internal/service/auth"
	"monolith/internal/service/user"

	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	authService *authService.Service
	userService *user.Service
}

func NewAccountHandler(db *database.DB) *AccountHandler {
	return &AccountHandler{
		authService: authService.NewService(db),
		userService: user.NewService(db),
	}
}

func (h *AccountHandler) Profile(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.userService.GetAccountByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AccountHandler) UpdatePreferences(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	var req user.UpdatePreferencesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	updatedUser, err := h.userService.UpdatePreferences(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update preferences"})
	}

	return c.JSON(http.StatusOK, updatedUser)
}

func (h *AccountHandler) GetSessions(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	sessions, err := h.authService.GetUserSessions(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidUserID) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch sessions"})
	}

	var currentSessionID string
	if sessionCookie, cookieErr := c.Cookie("session_id"); cookieErr == nil {
		currentSessionID = sessionCookie.Value
	}

	for i := range sessions {
		sessions[i].IsCurrent = sessions[i].ID.String() == currentSessionID
	}

	return c.JSON(http.StatusOK, sessions)
}

func (h *AccountHandler) RevokeSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}
	sessionID := c.Param("sessionId")

	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Session ID is required"})
	}

	err := h.authService.RevokeSession(c.Request().Context(), userID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidUserID):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, authService.ErrSessionNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to revoke session"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session revoked successfully"})
}

func (h *AccountHandler) RevokeAllOtherSessions(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	currentSessionID := ""
	if sessionCookie, cookieErr := c.Cookie("session_id"); cookieErr == nil {
		currentSessionID = sessionCookie.Value
	}

	revokedCount, err := h.authService.RevokeAllOtherSessions(c.Request().Context(), userID, currentSessionID)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidUserID) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to revoke sessions"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message":      "Other sessions revoked successfully",
		"revokedCount": revokedCount,
	})
}

func (h *AccountHandler) Register(c echo.Context) error {
	var req user.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrPasswordTooShort):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, authService.ErrUserAlreadyExists):
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}
	}

	response := map[string]any{
		"user": user,
	}

	return c.JSON(http.StatusCreated, response)
}
