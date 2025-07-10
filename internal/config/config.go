package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Security SecurityConfig
	Database DatabaseConfig
	Server   ServerConfig
	Logging  LoggingConfig
}

type SecurityConfig struct {
	SecretKey                            string
	LoginMaximumLifetimeDuration         time.Duration
	LoginMaximumInactiveLifetimeDuration time.Duration
	LoginCookieName                      string
	TokenRotationIntervalMinutes         int
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

type LoggingConfig struct {
	Level slog.Level
	Env   string
}

const (
	defaultTokenRotationIntervalMinutes = 10                  // 10 minutes
	defaultLoginMaximumLifetime         = 30 * 24 * time.Hour // 30 days
	defaultLoginInactiveLifetime        = 7 * 24 * time.Hour  // 7 days
	defaultLoginCookieName              = "session_token"     // Default cookie name for session tokens
	defaultDatabaseURL                  = "postgres://postgres:postgres@localhost:5432/my_db"
	defaultPort                         = "3001"
	defaultLogLevel                     = slog.LevelInfo
	defaultEnv                          = "development"
)

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load .env file: %v\n", err)
	}

	return &Config{
		Security: SecurityConfig{
			SecretKey:                            getEnvOrDefault("SECRET_KEY", "aTiONDsHeAngUaTeRvESteRUmbayaNCI"),
			LoginMaximumLifetimeDuration:         parseDurationOrDefault("LOGIN_MAXIMUM_LIFETIME_DURATION", defaultLoginMaximumLifetime),
			LoginMaximumInactiveLifetimeDuration: parseDurationOrDefault("LOGIN_MAXIMUM_INACTIVE_LIFETIME_DURATION", defaultLoginInactiveLifetime),
			LoginCookieName:                      getEnvOrDefault("LOGIN_COOKIE_NAME", defaultLoginCookieName),
			TokenRotationIntervalMinutes:         parseIntOrDefault("TOKEN_ROTATION_INTERVAL_MINUTES", defaultTokenRotationIntervalMinutes),
		},
		Database: DatabaseConfig{
			URL: getEnvOrDefault("DATABASE_URL", defaultDatabaseURL),
		},
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", defaultPort),
		},
		Logging: LoggingConfig{
			Level: parseLogLevelOrDefault("LOG_LEVEL", defaultLogLevel),
			Env:   getEnvOrDefault("ENV", defaultEnv),
		},
	}
}

func NewSecurityConfig() *SecurityConfig {
	config := NewConfig()
	return &config.Security
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func parseIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func parseLogLevelOrDefault(key string, defaultValue slog.Level) slog.Level {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "debug":
			return slog.LevelDebug
		case "info":
			return slog.LevelInfo
		case "warn":
			return slog.LevelWarn
		case "error":
			return slog.LevelError
		}
	}
	return defaultValue
}
