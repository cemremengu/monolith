package api

import (
	"net/http"

	"monolith/internal/auth"
	"monolith/internal/database"
	"monolith/internal/service/account"
	authService "monolith/internal/service/auth"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService    *authService.Service
	accountService *account.Service
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{
		authService:    authService.NewService(db),
		accountService: account.NewService(db),
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req auth.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		switch err {
		case authService.ErrPasswordTooShort:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case authService.ErrUserAlreadyExists:
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}
	}

	if tokenErr := h.authService.GenerateAndSetTokens(c, user.ID, user.Email, user.IsAdmin); tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	response := struct {
		User auth.UserAccount `json:"user"`
	}{
		User: *user,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req auth.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		if err == authService.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login"})
	}

	if tokenErr := h.authService.GenerateAndSetTokens(c, user.ID, user.Email, user.IsAdmin); tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	response := struct {
		User auth.UserAccount `json:"user"`
	}{
		User: *user,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.accountService.GetAccountByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	if sessionCookie, cookieErr := c.Cookie("session_id"); cookieErr == nil {
		_ = h.authService.RevokeSession(c.Request().Context(), "", sessionCookie.Value)
	}

	h.authService.ClearAuthCookies(c)
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Refresh token not found"})
	}

	sessionCookie, err := c.Cookie("session_id")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session not found"})
	}

	_, newAccessToken, newRefreshToken, err := h.authService.RefreshTokens(
		c.Request().Context(),
		refreshTokenCookie.Value,
		sessionCookie.Value,
	)
	if err != nil {
		switch err {
		case authService.ErrInvalidRefreshToken:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		case authService.ErrSessionExpired:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		case authService.ErrUserNotFound:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to refresh tokens"})
		}
	}

	h.authService.SetRefreshCookies(c, newAccessToken, newRefreshToken)
	return c.JSON(http.StatusOK, map[string]string{"message": "Tokens refreshed successfully"})
}

func (h *AuthHandler) GetSessions(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	sessions, err := h.authService.GetUserSessions(c.Request().Context(), userID)
	if err != nil {
		if err == authService.ErrInvalidUserID {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch sessions"})
	}

	var currentSessionID string
	if sessionCookie, cookieErr := c.Cookie("session_id"); cookieErr == nil {
		currentSessionID = sessionCookie.Value
	}

	for i := range sessions {
		sessions[i].IsCurrent = sessions[i].SessionID == currentSessionID
	}

	return c.JSON(http.StatusOK, sessions)
}

func (h *AuthHandler) RevokeSession(c echo.Context) error {
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
		switch err {
		case authService.ErrInvalidUserID:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case authService.ErrSessionNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to revoke session"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session revoked successfully"})
}

func (h *AuthHandler) RevokeAllOtherSessions(c echo.Context) error {
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
		if err == authService.ErrInvalidUserID {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to revoke sessions"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message":      "Other sessions revoked successfully",
		"revokedCount": revokedCount,
	})
}
