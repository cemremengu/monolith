package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"monolith"

	"github.com/lmittmann/tint"
)

type Config struct {
	Level slog.Level
}

func New(cfg Config) *slog.Logger {
	var handler slog.Handler

	if monolith.IsDevEnv() {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      cfg.Level,
			TimeFormat: "15:04:05",
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Level,
		})
	}

	return slog.New(handler)
}

func NewWithWriter(cfg Config, w io.Writer) *slog.Logger {
	var handler slog.Handler

	if monolith.IsDevEnv() {
		handler = tint.NewHandler(w, &tint.Options{
			Level:      cfg.Level,
			TimeFormat: "15:04:05",
		})
	} else {
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: cfg.Level,
		})
	}

	return slog.New(handler)
}

func isDev(env string) bool {
	env = strings.ToLower(env)
	return env == "" || env == "dev" || env == "development" || env == "local"
}

func GetLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
