package config

import (
	"os"
	"strconv"
	"time"

	"github.com/home/unixify/feature/ai_assisted/internal/database"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database database.Config
}

// ServerConfig holds the HTTP server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  time.Duration(getEnvInt("SERVER_READ_TIMEOUT", 10)) * time.Second,
			WriteTimeout: time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second,
			IdleTimeout:  time.Duration(getEnvInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second,
		},
		Database: database.Config{
			Host:         getEnvStr("DB_HOST", "localhost"),
			Port:         getEnvInt("DB_PORT", 5432),
			Username:     getEnvStr("DB_USER", "postgres"),
			Password:     getEnvStr("DB_PASSWORD", "postgres"),
			Database:     getEnvStr("DB_NAME", "unixify"),
			SSLMode:      getEnvStr("DB_SSLMODE", "disable"),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 50),
			MaxLifetime:  time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME", 30)) * time.Minute,
		},
	}
}

// Helper function to get environment variable as string with default value
func getEnvStr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as int with default value
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}