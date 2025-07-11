package middleware

import (
	"net/http"

	"monolith/internal/config"
	"monolith/internal/service/account"
	"monolith/internal/service/session"

	"github.com/labstack/echo/v4"
)

func SessionAuth(sessionService *session.Service, accountService *account.Service, securityConfig config.SecurityConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(securityConfig.LoginCookieName)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			}

			session, err := sessionService.GetSessionByToken(c.Request().Context(), cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Failed to retrieve session"})
			}

			account, err := accountService.GetAccountByID(c.Request().Context(), session.AccountID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Failed to retrieve account"})
			}

			c.Set("user_id", account.ID)
			c.Set("user_email", account.Email)
			c.Set("is_admin", account.IsAdmin)
			c.Set("session_id", session.ID)

			return next(c)
		}
	}
}

func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAdmin, ok := c.Get("is_admin").(bool)
			if !ok || !isAdmin {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
			}
			return next(c)
		}
	}
}
