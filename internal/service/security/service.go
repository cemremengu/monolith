package security

import (
	"strings"

	"github.com/labstack/echo/v4"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetDeviceInfo(c echo.Context) string {
	userAgent := c.Request().Header.Get("User-Agent")
	if userAgent == "" {
		return "Unknown Device"
	}

	switch {
	case strings.Contains(userAgent, "iPhone"):
		return "iPhone"
	case strings.Contains(userAgent, "Android"):
		return "Android Device"
	case strings.Contains(userAgent, "Mobile"):
		return "Mobile Device"
	case strings.Contains(userAgent, "Chrome"):
		return "Chrome Browser"
	case strings.Contains(userAgent, "Firefox"):
		return "Firefox Browser"
	case strings.Contains(userAgent, "Safari"):
		return "Safari Browser"
	default:
		return "Desktop Browser"
	}
}

func (s *Service) GetClientIP(c echo.Context) string {
	if xff := c.Request().Header.Get("X-Forwarded-For"); xff != "" {
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xri := c.Request().Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	return c.RealIP()
}
