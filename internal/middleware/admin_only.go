package middleware

import (
	"github.com/labstack/echo/v4"
)

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
