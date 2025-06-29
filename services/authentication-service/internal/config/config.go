package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	OTP        OTPConfig
	MailClient MailClientConfig
	RabbitMQ   RabbitMQConfig
}

// ServerConfig holds all server related configuration
type ServerConfig struct {
	Port        string
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

// JWTConfig holds all JWT related configuration
type JWTConfig struct {
	AccessSecret    string
	RefreshSecret   string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

// OTPConfig holds all OTP related configuration
type OTPConfig struct {
	Length        int
	Expiry        time.Duration
	MaxAttempts   int
	BypassEnabled bool
	BypassCode    string
}

// MailClientConfig holds all mail client related configuration
type MailClientConfig struct {
	BaseURL string
}

// RabbitMQConfig holds all RabbitMQ related configuration
type RabbitMQConfig struct {
	URL string
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
			Port:        getEnv("PORT", "50051"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres-auth"),
			Port:     getEnv("DB_PORT", "5432"), // authen
			User:     getEnv("DB_USER", "postgres-auth"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "auth"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:    getEnv("JWT_ACCESS_SECRET", "default-sercet-key"),
			RefreshSecret:   getEnv("JWT_REFRESH_SECRET", "default-refresh-secret-key"),
			AccessDuration:  time.Duration(getEnvAsInt("JWT_ACCESS_DURATION", 15)) * time.Minute,
			RefreshDuration: time.Duration(getEnvAsInt("JWT_REFRESH_DURATION", 24*7)) * time.Hour,
		},
		OTP: OTPConfig{
			Length:        getEnvAsInt("OTP_LENGTH", 6),
			Expiry:        time.Duration(getEnvAsInt("OTP_EXPIRY", 15)) * time.Minute,
			MaxAttempts:   getEnvAsInt("OTP_MAX_ATTEMPTS", 3),
			BypassEnabled: getOTPBypassEnabled(getEnv("ENVIRONMENT", "development")),
			BypassCode:    getOTPBypassCode(getEnv("ENVIRONMENT", "development")),
		},
		MailClient: MailClientConfig{
			BaseURL: getEnv("MAIL_SERVICE_URL", "http://mail-service:9002"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672"),
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

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getOTPBypassEnabled(environment string) bool {
	// Only enable OTP bypass in staging for performance testing
	return environment == "staging"
}

func getOTPBypassCode(environment string) string {
	// Only provide bypass code for staging
	if environment == "staging" {
		return "123456" // Simple bypass code for staging
	}
	return ""
}
