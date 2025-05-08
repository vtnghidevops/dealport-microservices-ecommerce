package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// ServerConfig holds all server related configuration
type ServerConfig struct {
	Port        string
	HTTPPort    string
	Environment string
}

// DatabaseConfig holds all database related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AuthConfig holds authentication related configuration
type AuthConfig struct {
	AuthServiceURL string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig(path string) (*Config, error) {
	// Load .env file if it exists
	if path != "" {
		err := godotenv.Load(path)
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Initialize config
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "50052"),
			// HTTPPort:    getEnv("HTTP_PORT", "9001"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("USER_DB_HOST", "postgres-users"),
			Port:     getEnv("USER_DB_PORT", "5432"), // user-service
			User:     getEnv("USER_DB_USER", "postgres-users"),
			Password: getEnv("USER_DB_PASSWORD", "password"),
			DBName:   getEnv("USER_DB_NAME", "users"),
			SSLMode:  getEnv("USER_DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			AuthServiceURL: getEnv("AUTH_SERVICE_URL", "auth-service:50051"),
		},
	}

	return cfg, nil
}

// PostgresConnectionString returns the connection string for the postgres database
func (c *Config) PostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
