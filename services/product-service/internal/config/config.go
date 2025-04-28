package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

// Config (Cấu hình) holds all configuration for the product service (dịch vụ sản phẩm)
type Config struct {
	Server   ServerConfig   // Server (máy chủ)
	Database DatabaseConfig // Database (cơ sở dữ liệu)
}

// ServerConfig (Cấu hình máy chủ)
type ServerConfig struct {
	HTTPPort string // HTTP port (cổng HTTP)
	GRPCPort string // gRPC port (cổng gRPC)
}

// DatabaseConfig (Cấu hình cơ sở dữ liệu)
type DatabaseConfig struct {
	Host     string // Host (máy chủ)
	Port     string // Port (cổng)
	User     string // User (người dùng)
	Password string // Password (mật khẩu)
	DBName   string // Database name (tên cơ sở dữ liệu)
	SSLMode  string // SSL mode (chế độ SSL)
}

// LoadConfig (Tải cấu hình) loads configuration from environment variables (biến môi trường)
func LoadConfig(path string) (*Config, error) {
	if path != "" {
		_ = godotenv.Load(path) // Load .env file if provided (nếu có)
	}

	cfg := &Config{
		Server: ServerConfig{
			HTTPPort: getEnv("HTTP_PORT", "8082"),
			GRPCPort: getEnv("GRPC_PORT", "50053"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "postgres-products"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "products"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
	return cfg, nil
}

// PostgresConnectionString (Chuỗi kết nối Postgres)
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

// getEnv (lấy biến môi trường)
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}