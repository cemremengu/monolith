package middleware

import (
	"net/http"

	"monolith/internal/auth"

	"github.com/labstack/echo/v4"
)

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			}

			claims, err := auth.ValidateToken(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("is_admin", claims.IsAdmin)

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
