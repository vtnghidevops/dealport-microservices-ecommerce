package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config  holds all configuration for the product service
type Config struct {
	Server   ServerConfig   // Server
	Database DatabaseConfig // Database
	Storage  StorageConfig  `mapstructure:"storage"`
}

// ServerConfig
type ServerConfig struct {
	HTTPPort string // HTTP port
	GRPCPort string // gRPC port
}

// DatabaseConfig
type DatabaseConfig struct {
	Host     string // Host
	Port     string // Port
	User     string // User
	Password string // Password
	DBName   string // Database name
	SSLMode  string // SSL mode
}

// StorageConfig contains storage configuration
type StorageConfig struct {
	Provider string      `mapstructure:"provider"` // "local" or "minio"
	MinIO    MinIOConfig `mapstructure:"minio"`
}

// MinIOConfig contains MinIO configuration
type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	BucketName      string `mapstructure:"bucket_name"`
	Location        string `mapstructure:"location"`
	BaseURL         string `mapstructure:"base_url"`
	PresignedTTL    int    `mapstructure:"presigned_ttl"` // In seconds
}

// LoadConfig loads configuration from environment variables
func LoadConfig(path string) (*Config, error) {
	if path != "" {
		_ = godotenv.Load(path) // Load .env file if provided
	}

	cfg := &Config{
		Server: ServerConfig{
			HTTPPort: getEnv("HTTP_PORT", "8082"),
			GRPCPort: getEnv("GRPC_PORT", "50053"),
			// PRODUCT_DSN=host=postgres-products port=5432 user=postgres-products password=password dbname=products sslmode=disable timezone=UTC connect_timeout=5
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres-products"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres-products"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "products"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Storage: StorageConfig{
			Provider: getEnv("STORAGE_PROVIDER", "local"),
			MinIO: MinIOConfig{
				Endpoint:        getEnv("MINIO_ENDPOINT", ""),
				AccessKeyID:     getEnv("MINIO_ACCESS_KEY_ID", ""),
				SecretAccessKey: getEnv("MINIO_SECRET_ACCESS_KEY", ""),
				UseSSL:          getEnv("MINIO_USE_SSL", "false") == "true",
				BucketName:      getEnv("MINIO_BUCKET_NAME", "images"),
				Location:        getEnv("MINIO_LOCATION", "us-east-1"),
				BaseURL:         getEnv("MINIO_BASE_URL", ""),
				PresignedTTL:    getEnvAsInt("MINIO_PRESIGNED_TTL", 3600),
			},
		},
	}
	return cfg, nil
}

// PostgresConnectionString
func (c *Config) PostgresConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// getEnv
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, fmt.Sprintf("%d", defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
