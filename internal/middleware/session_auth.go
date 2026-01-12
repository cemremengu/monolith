package middleware

import (
	"net/http"

	"monolith/internal/config"
	"monolith/internal/service/auth"

	"github.com/labstack/echo/v4"
)

func SessionAuth(authService *auth.Service, securityConfig config.SecurityConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(securityConfig.LoginCookieName)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			}

			// Get all auth context in a single database query
			authCtx, err := authService.GetAuthContextByToken(c.Request().Context(), cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Failed to retrieve session"})
			}

			user := &auth.AuthUser{
				AccountID: authCtx.AccountID,
				Email:     authCtx.AccountEmail,
				IsAdmin:   authCtx.AccountIsAdmin,
				SessionID: authCtx.SessionID,
			}

			c.Set("user", user)

			return next(c)
		}
	}
}
