package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	Redis    RedisConfig
	Database DatabaseConfig
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	GRPCPort string
	HTTPPort string
}

// RedisConfig represents the Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig reads configuration from environment variables or a config file
func LoadConfig(path string) (*Config, error) {
	// Default configuration
	config := &Config{
		Server: ServerConfig{
			GRPCPort: getEnv("GRPC_PORT", "50054"),
			// HTTPPort: getEnv("HTTP_PORT", "8084"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "redis-cart"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		// Database: DatabaseConfig{
		// 	Host:     getEnv("DB_HOST", "postgres-cart"),
		// 	Port:     getEnv("DB_PORT", "5432"),
		// 	User:     getEnv("DB_USER", "postgres-cart"),
		// 	Password: getEnv("DB_PASSWORD", "password"),
		// 	DBName:   getEnv("DB_NAME", "cart"),
		// 	SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		// },
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
		if viper.IsSet("server.grpcPort") {
			config.Server.GRPCPort = viper.GetString("server.grpcPort")
		}
		if viper.IsSet("server.httpPort") {
			config.Server.HTTPPort = viper.GetString("server.httpPort")
		}
		if viper.IsSet("redis.host") {
			config.Redis.Host = viper.GetString("redis.host")
		}
		if viper.IsSet("redis.port") {
			config.Redis.Port = viper.GetString("redis.port")
		}
		if viper.IsSet("redis.password") {
			config.Redis.Password = viper.GetString("redis.password")
		}
		if viper.IsSet("redis.db") {
			config.Redis.DB = viper.GetInt("redis.db")
		}
		if viper.IsSet("database.host") {
			config.Database.Host = viper.GetString("database.host")
		}
		if viper.IsSet("database.port") {
			config.Database.Port = viper.GetString("database.port")
		}
		if viper.IsSet("database.user") {
			config.Database.User = viper.GetString("database.user")
		}
		if viper.IsSet("database.password") {
			config.Database.Password = viper.GetString("database.password")
		}
		if viper.IsSet("database.dbName") {
			config.Database.DBName = viper.GetString("database.dbName")
		}
		if viper.IsSet("database.sslMode") {
			config.Database.SSLMode = viper.GetString("database.sslMode")
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
