package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
}

// ServerConfig contains server-related settings
type ServerConfig struct {
	Host     string
	Port     int
	GRPCPort string
	HTTPPort string
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

// LoadConfig reads configuration from environment variables or a config file
func LoadConfig(path string) (*Config, error) {
	// Default configuration
	config := &Config{
		Server: ServerConfig{
			GRPCPort: getEnv("GRPC_PORT", "50055"),
			// HTTPPort: getEnv("HTTP_PORT", "8085"),
		},
		MongoDB: MongoDBConfig{
			Host:     getEnv("MONGO_HOST", "localhost"),
			Port:     getEnv("MONGO_PORT", "27019"),
			User:     getEnv("MONGO_USER", ""),
			Password: getEnv("MONGO_PASSWORD", ""),
			Database: getEnv("MONGO_DATABASE", "checkout"),
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27019/checkout"),
		},
	}

	// If a config file path is provided, load it
	if path != "" {
		viper.SetConfigFile(path)
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}

		// Override config with values from the config file
		if viper.IsSet("server.host") {
			config.Server.Host = viper.GetString("server.host")
		}
		if viper.IsSet("server.port") {
			config.Server.Port = viper.GetInt("server.port")
		}
		if viper.IsSet("server.grpcPort") {
			config.Server.GRPCPort = viper.GetString("server.grpcPort")
		}
		if viper.IsSet("server.httpPort") {
			config.Server.HTTPPort = viper.GetString("server.httpPort")
		}
		if viper.IsSet("mongodb.uri") {
			config.MongoDB.URI = viper.GetString("mongodb.uri")
		}
		if viper.IsSet("mongodb.database") {
			config.MongoDB.Database = viper.GetString("mongodb.database")
		}
		if viper.IsSet("mongodb.host") {
			config.MongoDB.Host = viper.GetString("mongodb.host")
		}
		if viper.IsSet("mongodb.port") {
			config.MongoDB.Port = viper.GetString("mongodb.port")
		}
		if viper.IsSet("mongodb.user") {
			config.MongoDB.User = viper.GetString("mongodb.user")
		}
		if viper.IsSet("mongodb.password") {
			config.MongoDB.Password = viper.GetString("mongodb.password")
		}
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
