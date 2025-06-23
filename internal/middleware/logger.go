package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Logger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogLatency:      true,
		LogRemoteIP:     true,
		LogUserAgent:    true,
		LogError:        true,
		LogMethod:       true,
		LogResponseSize: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				slog.Error("http request",
					"error", v.Error,
					"remote_ip", v.RemoteIP,
					"method", v.Method,
					"URI", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"size", v.ResponseSize,
				)
			} else {
				slog.Info("http request",
					"remote_ip", v.RemoteIP,
					"method", v.Method,
					"URI", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"size", v.ResponseSize,
				)
			}

			return nil
		},
	})
}
