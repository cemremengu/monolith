package api

import (
	"errors"
	"net/http"

	"monolith/internal/database"
	"monolith/internal/service/account"
	authService "monolith/internal/service/auth"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	authService    *authService.Service
	accountService *account.Service
}

func NewAccountHandler(db *database.DB) *AccountHandler {
	return &AccountHandler{
		authService:    authService.NewService(db),
		accountService: account.NewService(db),
	}
}

func (h *AccountHandler) Profile(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	account, err := h.accountService.GetAccountByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
	}

	return c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) UpdatePreferences(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	var req account.UpdatePreferencesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	updatedAccount, err := h.accountService.UpdatePreferences(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update preferences"})
	}

	return c.JSON(http.StatusOK, updatedAccount)
}

func (h *AccountHandler) GetSessions(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
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

	err = h.authService.RevokeSession(c.Request().Context(), userID, sessionID)
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
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	var currentSessionID uuid.UUID
	if sessionCookie, cookieErr := c.Cookie("session_id"); cookieErr == nil {
		if parsedID, parseErr := uuid.Parse(sessionCookie.Value); parseErr == nil {
			currentSessionID = parsedID
		}
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
