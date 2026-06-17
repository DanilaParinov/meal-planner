package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config содержит все конфигурационные параметры приложения
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

	// Test users
	TestUsers map[string]string // api_key -> device_id
}

// Load загружает конфигурацию из .env файла и переменных окружения
func Load() (*Config, error) {
	// Загружаем .env файл (не критично если не существует)
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

		TestUsers: make(map[string]string),
	}

	// Загружаем тестовых пользователей
	// Пример: TEST_USER_1_API_KEY=abc123, TEST_USER_1_DEVICE_ID=device-1
	for i := 1; i <= 10; i++ {
		apiKeyEnv := fmt.Sprintf("TEST_USER_%d_API_KEY", i)
		deviceIDEnv := fmt.Sprintf("TEST_USER_%d_DEVICE_ID", i)

		apiKey := os.Getenv(apiKeyEnv)
		deviceID := os.Getenv(deviceIDEnv)

		if apiKey != "" && deviceID != "" {
			cfg.TestUsers[apiKey] = deviceID
		}
	}

	return cfg, nil
}

// GetDatabaseURL формирует строку подключения к БД
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

// IsDevelopment проверяет, находимся ли мы в режиме разработки
func (c *Config) IsDevelopment() bool {
	return c.ServerEnv == "development" || c.ServerEnv == "dev"
}

// getEnv получает переменную окружения с fallback значением
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает целое число из переменной окружения
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	return defaultValue
}
