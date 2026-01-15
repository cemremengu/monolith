package middleware

import (
	"monolith/internal/service/auth"

	"github.com/labstack/echo/v4"
)

func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*auth.AuthUser)
			if !ok || !user.IsAdmin {
				// Return 404 instead of 403 for admin-only routes for security reasons
				return echo.ErrNotFound
			}
			return next(c)
		}
	}
}
