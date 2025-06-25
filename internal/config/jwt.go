package config

import (
	"os"
	"time"
)

const defaultAccessTokenDuration = 10 * time.Minute

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	SessionTimeout       time.Duration
}

func NewJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:               getEnvOrDefault("SECRET_KEY", "your-256-bit-secret"),
		AccessTokenDuration:  parseDurationOrDefault("JWT_ACCESS_TOKEN_DURATION", defaultAccessTokenDuration),
		RefreshTokenDuration: parseDurationOrDefault("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour),
		SessionTimeout:       parseDurationOrDefault("SESSION_TIMEOUT", 0),
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
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
