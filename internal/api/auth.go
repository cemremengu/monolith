package api

import (
	"errors"
	"net/http"

	"monolith/internal/database"
	authService "monolith/internal/service/auth"
	"monolith/internal/service/user"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *authService.Service
	userService *user.Service
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{
		authService: authService.NewService(db),
		userService: user.NewService(db),
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

	if tokenErr := h.authService.GenerateAndSetTokens(c, user.ID, user.Email, user.IsAdmin); tokenErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	response := map[string]any{
		"user": user,
	}

	return c.JSON(http.StatusOK, response)
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
		// nuke cookies on any refresh error
		h.authService.ClearAuthCookies(c)

		switch {
		case errors.Is(err, authService.ErrInvalidRefreshToken):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		case errors.Is(err, authService.ErrSessionExpired):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		case errors.Is(err, authService.ErrUserNotFound):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to refresh tokens"})
		}
	}

	h.authService.SetRefreshCookies(c, newAccessToken, newRefreshToken)
	return c.JSON(http.StatusOK, map[string]string{"message": "Tokens refreshed successfully"})
}
