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
		HandleError:     true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			fields := []any{
				"remote_ip", v.RemoteIP,
				"method", v.Method,
				"URI", v.URI,
				"status", v.Status,
				"latency", v.Latency,
				"size", v.ResponseSize,
			}

			if v.Error != nil {
				fields = append(fields, "error", v.Error)
			}

			switch {
			case v.Status >= 500:
				slog.Error("http request", fields...)
			case v.Status >= 400:
				slog.Warn("http request", fields...)
			default:
				slog.Info("http request", fields...)
			}

			return nil
		},
	})
}
