package config

import (
	"os"
	"strconv"
)

// Config represents the application configuration
type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
}

// ServerConfig contains server-related settings
type ServerConfig struct {
	GRPCPort string
}

// MongoDBConfig contains MongoDB connection settings
type MongoDBConfig struct {
	URI      string
	Database string
	Host     string
	Port     string
	User     string
	Password string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Default configuration
	config := &Config{
		Server: ServerConfig{
			GRPCPort: getEnv("GRPC_PORT", "50056"),
		},
		MongoDB: MongoDBConfig{
			Host:     getEnv("MONGO_HOST", "mongo-logger"),
			Port:     getEnv("MONGO_PORT", "27017"),
			User:     getEnv("MONGO_USERNAME", "logger_user"),
			Password: getEnv("MONGO_PASSWORD", "password"),
			Database: getEnv("MONGO_DATABASE", "logs"),
			URI:      getEnv("MONGO_URI", "mongodb://logger_user:password@mongo-logger:27017/logs"),
		},
	}

	return config, nil
}

// Helper functions to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
