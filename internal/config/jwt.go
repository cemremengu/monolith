package config

import (
	"os"
	"strconv"
	"time"
)

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func NewJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:               getEnvOrDefault("SECRET_KEY", "your-256-bit-secret"),
		AccessTokenDuration:  parseDurationOrDefault("JWT_ACCESS_TOKEN_DURATION", 24*time.Hour),
		RefreshTokenDuration: parseDurationOrDefault("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if hours, err := strconv.Atoi(value); err == nil {
			return time.Duration(hours) * time.Hour
		}
	}
	return defaultValue
}