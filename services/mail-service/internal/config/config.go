package config

import (
	"log"
	"os"
	"strconv"

	"mail-service/internal/mailer"
)

// Config holds all the configuration for the mail service
type Config struct {
	Mailer mailer.Mail
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	// Create mail configuration from environment variables
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDR"),
	}

	// Log mail configuration (exclude password)
	log.Printf("Mail configuration: domain=%s, host=%s, port=%d, username=%s, encryption=%s",
		m.Domain, m.Host, m.Port, m.Username, m.Encryption)

	return Config{
		Mailer: m,
	}
}
