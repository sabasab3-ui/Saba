package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	// Server
	Port     int
	Host     string
	LogLevel string

	// Database
	DatabaseURL string
	DBPath      string

	// JWT
	JWTSecret             string
	JWTExpirationMinutes  int
	JWTRefreshExpiration  time.Duration

	// OpenAI
	OpenAIAPIKey  string
	OpenAIModel   string

	// External Services
	SABAGatewayURL string
	RedisURL       string

	// Features
	EnablePayments bool
	Environment    string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Port:                    getIntEnv("PORT", 8080),
		Host:                    getEnv("HOST", "0.0.0.0"),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
		DatabaseURL:             getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/saba"),
		DBPath:                  getEnv("DB_PATH", "./data"),
		JWTSecret:               getEnv("JWT_SECRET", ""),
		JWTExpirationMinutes:    getIntEnv("JWT_EXPIRATION_MINUTES", 60),
		OpenAIAPIKey:            getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:             getEnv("OPENAI_MODEL", "gpt-4"),
		SABAGatewayURL:          getEnv("SABA_GATEWAY_URL", "http://localhost:8080"),
		RedisURL:                getEnv("REDIS_URL", "redis://localhost:6379"),
		EnablePayments:          getBoolEnv("ENABLE_PAYMENTS", false),
		Environment:             getEnv("ENVIRONMENT", "development"),
	}

	cfg.JWTRefreshExpiration = time.Duration(cfg.JWTExpirationMinutes*2) * time.Minute

	// Validate critical settings
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required configuration is set.
func (c *Config) Validate() error {
	if c.JWTSecret == "" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	if c.OpenAIAPIKey == "" && c.Environment == "production" {
		return fmt.Errorf("OPENAI_API_KEY must be set in production")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid PORT: %d", c.Port)
	}

	return nil
}

// IsProduction returns true if running in production environment.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// Helper functions

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getIntEnv(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getBoolEnv(key string, defaultVal bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultVal
}
