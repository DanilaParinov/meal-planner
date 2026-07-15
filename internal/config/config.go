package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration parameters
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server
	ServerPort string
	ServerEnv  string
}

// Load loads configuration from the .env file and environment variables
func Load() (*Config, error) {
	// Load .env file (not critical if it doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "changeme_user"),
		DBPassword: getEnv("DB_PASSWORD", "changeme_password"),
		DBName:    getEnv("DB_NAME", "meal_planner"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),

		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerEnv:  getEnv("SERVER_ENV", "development"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate checks that configuration values are sane and safe to run with
func (c *Config) validate() error {
	if _, err := strconv.Atoi(c.ServerPort); err != nil {
		return fmt.Errorf("SERVER_PORT %q must be a number", c.ServerPort)
	}
	if _, err := strconv.Atoi(c.DBPort); err != nil {
		return fmt.Errorf("DB_PORT %q must be a number", c.DBPort)
	}

	validSSLModes := map[string]bool{
		"disable":     true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if !validSSLModes[c.DBSSLMode] {
		return fmt.Errorf("DB_SSLMODE %q must be one of: disable, require, verify-ca, verify-full", c.DBSSLMode)
	}

	if !c.IsDevelopment() && (c.DBUser == "changeme_user" || c.DBPassword == "changeme_password") {
		return fmt.Errorf("DB_USER/DB_PASSWORD must be set to non-default values outside development (SERVER_ENV=%s)", c.ServerEnv)
	}

	return nil
}

// GetDatabaseURL builds the database connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

// IsDevelopment checks whether we're running in development mode
func (c *Config) IsDevelopment() bool {
	return c.ServerEnv == "development" || c.ServerEnv == "dev"
}

// getEnv retrieves an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an integer value from an environment variable
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	return defaultValue
}
