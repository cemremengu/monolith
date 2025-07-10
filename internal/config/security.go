package config

import (
	"os"
	"strconv"
	"time"
)

type SecurityConfig struct {
	SecretKey                            string
	LoginMaximumLifetimeDuration         time.Duration
	LoginMaximumInactiveLifetimeDuration time.Duration
	LoginCookieName                      string
	TokenRotationIntervalMinutes         int
}

const (
	defaultTokenRotationIntervalMinutes = 10                  // 10 minutes
	defaultLoginMaximumLifetime         = 30 * 24 * time.Hour // 30 days
	defaultLoginInactiveLifetime        = 7 * 24 * time.Hour  // 7 days
	defaultLoginCookieName              = "session_token"     // Default cookie name for session tokens
)

func NewSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		SecretKey:                            getEnvOrDefault("SECRET_KEY", "aTiONDsHeAngUaTeRvESteRUmbayaNCI"),
		LoginMaximumLifetimeDuration:         parseDurationOrDefault("LOGIN_MAXIMUM_LIFETIME_DURATION", defaultLoginMaximumLifetime),
		LoginMaximumInactiveLifetimeDuration: parseDurationOrDefault("LOGIN_MAXIMUM_INACTIVE_LIFETIME_DURATION", defaultLoginInactiveLifetime),
		LoginCookieName:                      getEnvOrDefault("LOGIN_COOKIE_NAME", defaultLoginCookieName),
		TokenRotationIntervalMinutes:         parseIntOrDefault("TOKEN_ROTATION_INTERVAL_MINUTES", defaultTokenRotationIntervalMinutes),
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

func parseIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
