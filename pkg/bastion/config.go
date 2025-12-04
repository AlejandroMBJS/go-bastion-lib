package bastion

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application.
type Config struct {
	// Server configuration
	Port         int           `env:"PORT"`
	Env          string        `env:"ENV"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT"`

	// Templating configuration
	TemplateRoot string `env:"TEMPLATE_ROOT"`

	// Security features
	EnableCSRF      bool `env:"ENABLE_CSRF"`
	EnableJWT       bool `env:"ENABLE_JWT"`
	EnableRateLimit bool `env:"ENABLE_RATE_LIMIT"`

	// JWT configuration
	JWTSecret     string        `env:"JWT_SECRET"`
	JWTAccessTTL  time.Duration `env:"JWT_ACCESS_TTL"`
	JWTRefreshTTL time.Duration `env:"JWT_REFRESH_TTL"`

	// Rate limiting configuration
	RateLimitRequests int           `env:"RATE_LIMIT_REQUESTS"`
	RateLimitWindow   time.Duration `env:"RATE_LIMIT_WINDOW"`

	// Security headers
	EnableSecurityHeaders bool   `env:"ENABLE_SECURITY_HEADERS"`
	CSP                   string `env:"CSP"`

	// Database configuration (agnostic)
	DBDriver string `env:"DB_DRIVER"`
	DSN      string `env:"DSN"`

	// Logging
	LogLevel string `env:"LOG_LEVEL"`
}

// DefaultConfig returns a Config instance with safe defaults.
func DefaultConfig() Config {
	return Config{
		Port:         9876,
		Env:          "development",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,

		TemplateRoot: "templates", // Default template root

		EnableCSRF:      false,
		EnableJWT:       false,
		EnableRateLimit: false,

		JWTSecret:     "",
		JWTAccessTTL:  15 * time.Minute,
		JWTRefreshTTL: 7 * 24 * time.Hour,

		RateLimitRequests: 100,
		RateLimitWindow:   1 * time.Minute,

		EnableSecurityHeaders: true,
		CSP:                   "default-src 'self';",

		DBDriver: "",
		DSN:      "",

		LogLevel: "info",
	}
}

// LoadConfigFromEnv loads configuration from environment variables.
// Values from environment override the defaults.
func LoadConfigFromEnv() Config {
	cfg := DefaultConfig()

	// Helper function to get environment variable
	getEnv := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}

	// Helper function to get boolean environment variable
	getEnvBool := func(key string, defaultValue bool) bool {
		if value := os.Getenv(key); value != "" {
			if b, err := strconv.ParseBool(value); err == nil {
				return b
			}
		}
		return defaultValue
	}

	// Helper function to get int environment variable
	getEnvInt := func(key string, defaultValue int) int {
		if value := os.Getenv(key); value != "" {
			if i, err := strconv.Atoi(value); err == nil {
				return i
			}
		}
		return defaultValue
	}

	// Helper function to get duration environment variable
	getEnvDuration := func(key string, defaultValue time.Duration) time.Duration {
		if value := os.Getenv(key); value != "" {
			if d, err := time.ParseDuration(value); err == nil {
				return d
			}
		}
		return defaultValue
	}

	// Load configuration from environment variables
	cfg.Port = getEnvInt("PORT", cfg.Port)
	cfg.Env = getEnv("ENV", cfg.Env)
	cfg.ReadTimeout = getEnvDuration("READ_TIMEOUT", cfg.ReadTimeout)
	cfg.WriteTimeout = getEnvDuration("WRITE_TIMEOUT", cfg.WriteTimeout)
	cfg.IdleTimeout = getEnvDuration("IDLE_TIMEOUT", cfg.IdleTimeout)

	cfg.EnableCSRF = getEnvBool("ENABLE_CSRF", cfg.EnableCSRF)
	cfg.EnableJWT = getEnvBool("ENABLE_JWT", cfg.EnableJWT)
	cfg.EnableRateLimit = getEnvBool("ENABLE_RATE_LIMIT", cfg.EnableRateLimit)

	cfg.JWTSecret = getEnv("JWT_SECRET", cfg.JWTSecret)
	cfg.JWTAccessTTL = getEnvDuration("JWT_ACCESS_TTL", cfg.JWTAccessTTL)
	cfg.JWTRefreshTTL = getEnvDuration("JWT_REFRESH_TTL", cfg.JWTRefreshTTL)

	cfg.RateLimitRequests = getEnvInt("RATE_LIMIT_REQUESTS", cfg.RateLimitRequests)
	cfg.RateLimitWindow = getEnvDuration("RATE_LIMIT_WINDOW", cfg.RateLimitWindow)

	cfg.EnableSecurityHeaders = getEnvBool("ENABLE_SECURITY_HEADERS", cfg.EnableSecurityHeaders)
	cfg.CSP = getEnv("CSP", cfg.CSP)

	cfg.DBDriver = getEnv("DB_DRIVER", cfg.DBDriver)
	cfg.DSN = getEnv("DSN", cfg.DSN)

	cfg.LogLevel = strings.ToLower(getEnv("LOG_LEVEL", cfg.LogLevel))

	return cfg
}
