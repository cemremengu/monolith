package api

import (
	"errors"
	"net/http"

	authService "monolith/internal/auth"
	loginService "monolith/internal/login"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	loginService *loginService.Service
	authService  *authService.Service
}

func NewAuthHandler(loginService *loginService.Service, authService *authService.Service) *AuthHandler {
	return &AuthHandler{
		loginService: loginService,
		authService:  authService,
	}
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var req loginService.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").Wrap(err)
	}

	user, err := h.loginService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, loginService.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to login").Wrap(err)
	}

	session, tokenErr := h.authService.CreateSession(c.Request().Context(), &authService.CreateSessionRequest{
		AccountID: user.ID,
		ClientIP:  c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	})

	if tokenErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create session").Wrap(tokenErr)
	}

	h.authService.SetSessionCookies(c, session)

	response := map[string]any{
		"message": "Login successful",
	}

	return c.JSON(http.StatusOK, response)
}

// Logout revokes the current session when present and clears authentication cookies.
func (h *AuthHandler) Logout(c *echo.Context) error {
	revokeErr := h.authService.RevokeSessionFromCookie(c)
	h.authService.ClearAuthCookies(c)
	if revokeErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to logout").Wrap(revokeErr)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
