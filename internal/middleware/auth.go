package middleware

import (
	"net/http"

	"monolith/internal/config"
	"monolith/internal/service/account"
	"monolith/internal/service/auth"

	"github.com/labstack/echo/v4"
)

func SessionAuth(authService *auth.Service, accountService *account.Service, securityConfig config.SecurityConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(securityConfig.LoginCookieName)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			}

			session, err := authService.GetSessionByToken(c.Request().Context(), cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Failed to retrieve session"})
			}

			account, err := accountService.GetAccountByID(c.Request().Context(), session.AccountID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Failed to retrieve account"})
			}

			user := &auth.AuthUser{
				AccountID: account.ID,
				Email:     account.Email,
				IsAdmin:   account.IsAdmin,
				SessionID: session.ID,
			}

			c.Set("user", user)

			return next(c)
		}
	}
}

func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAdmin, ok := c.Get("is_admin").(bool)
			if !ok || !isAdmin {
				// Return 404 instead of 403 for admin-only routes for security reasons
				return echo.ErrNotFound
			}
			return next(c)
		}
	}
}
