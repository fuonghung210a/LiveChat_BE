package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Connection string
		Host       string
		Port       string
		Database   string
		Username   string
		Password   string
	}
}

func LoadENV() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg := &Config{}

	// Server config
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// Database config
	cfg.Database.Connection = getEnv("DATABASE_CONNECTION", "mysql")
	cfg.Database.Host = getEnv("DATABASE_HOST", "127.0.0.1")
	cfg.Database.Port = getEnv("DATABASE_PORT", "3306")
	cfg.Database.Database = getEnv("DATABASE_DATABASE", "ginapp")
	cfg.Database.Username = getEnv("DATABASE_USERNAME", "root")
	cfg.Database.Password = getEnv("DATABASE_PASSWORD", "")

	return cfg
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
