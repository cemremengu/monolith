package config

import (
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
	LDAP     LDAPConfig
}

type LDAPConfig struct {
	Enabled           bool
	Host              string
	Port              int
	BindDN            string
	BindPassword      string
	BaseDN            string
	SearchFilter      string
	UsernameAttribute string
	EmailAttribute    string
	NameAttribute     string
	StartTLS          bool
	SkipTLSVerify     bool
	AutoProvision     bool
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
}

const (
	defaultTokenRotationIntervalMinutes = 10
	defaultLoginMaximumLifetime         = 30 * 24 * time.Hour
	defaultLoginInactiveLifetime        = 7 * 24 * time.Hour
	defaultLoginCookieName              = "session_token"
	defaultDatabaseURL                  = "postgres://postgres:postgres@localhost:5432/my_db"
	defaultPort                         = "3001"
	defaultLogLevel                     = slog.LevelInfo
	defaultLDAPPort                     = 389
	defaultLDAPSearchFilter             = "(|(uid=%s)(mail=%s))"
	defaultLDAPUsernameAttribute        = "uid"
	defaultLDAPEmailAttribute           = "mail"
	defaultLDAPNameAttribute            = "cn"
)

func NewConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Security: SecurityConfig{
			SecretKey:                            getEnvOrDefault("SECRET_KEY", "VrrkYCoULJMS6hyCfPrf6ThBJqkpWrKn6O2IIMj1Z3s="),
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
		},
		LDAP: LDAPConfig{
			Enabled:           parseBoolOrDefault("LDAP_ENABLED", false),
			Host:              getEnvOrDefault("LDAP_HOST", ""),
			Port:              parseIntOrDefault("LDAP_PORT", defaultLDAPPort),
			BindDN:            getEnvOrDefault("LDAP_BIND_DN", ""),
			BindPassword:      getEnvOrDefault("LDAP_BIND_PASSWORD", ""),
			BaseDN:            getEnvOrDefault("LDAP_BASE_DN", ""),
			SearchFilter:      getEnvOrDefault("LDAP_SEARCH_FILTER", defaultLDAPSearchFilter),
			UsernameAttribute: getEnvOrDefault("LDAP_USERNAME_ATTRIBUTE", defaultLDAPUsernameAttribute),
			EmailAttribute:    getEnvOrDefault("LDAP_EMAIL_ATTRIBUTE", defaultLDAPEmailAttribute),
			NameAttribute:     getEnvOrDefault("LDAP_NAME_ATTRIBUTE", defaultLDAPNameAttribute),
			StartTLS:          parseBoolOrDefault("LDAP_START_TLS", false),
			SkipTLSVerify:     parseBoolOrDefault("LDAP_SKIP_TLS_VERIFY", false),
			AutoProvision:     parseBoolOrDefault("LDAP_AUTO_PROVISION", true),
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

func parseBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
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
