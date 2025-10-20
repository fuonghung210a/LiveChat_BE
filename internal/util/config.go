package util

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port    string
		GinMode string
	}
	Database struct {
		Connection      string
		Host            string
		Port            string
		Database        string
		Username        string
		Password        string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime time.Duration
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
		PoolSize int
	}
	JWT struct {
		Secret string
		Expiry time.Duration
	}
	CORS struct {
		AllowedOrigins string
		AllowedMethods string
		AllowedHeaders string
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
	cfg.Server.GinMode = getEnv("GIN_MODE", "debug")

	// Database config
	cfg.Database.Connection = getEnv("DATABASE_CONNECTION", "mysql")
	cfg.Database.Host = getEnv("DATABASE_HOST", "127.0.0.1")
	cfg.Database.Port = getEnv("DATABASE_PORT", "3307")
	cfg.Database.Database = getEnv("DATABASE_DATABASE", "ginapp")
	cfg.Database.Username = getEnv("DATABASE_USERNAME", "root")
	cfg.Database.Password = getEnv("DATABASE_PASSWORD", "")
	cfg.Database.MaxIdleConns = getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 10)
	cfg.Database.MaxOpenConns = getEnvAsInt("DATABASE_MAX_OPEN_CONNS", 100)
	cfg.Database.ConnMaxLifetime = time.Duration(getEnvAsInt("DATABASE_CONN_MAX_LIFETIME", 3600)) * time.Second

	// Redis config
	cfg.Redis.Host = getEnv("REDIS_HOST", "127.0.0.1")
	cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvAsInt("REDIS_DB", 0)
	cfg.Redis.PoolSize = getEnvAsInt("REDIS_POOL_SIZE", 10)

	// JWT config
	cfg.JWT.Secret = getEnv("JWT_SECRET", "default-secret-change-this")
	expiryStr := getEnv("JWT_EXPIRY", "24h")
	if expiry, err := time.ParseDuration(expiryStr); err == nil {
		cfg.JWT.Expiry = expiry
	} else {
		cfg.JWT.Expiry = 24 * time.Hour
	}

	// CORS config
	cfg.CORS.AllowedOrigins = getEnv("CORS_ALLOWED_ORIGINS", "*")
	cfg.CORS.AllowedMethods = getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
	cfg.CORS.AllowedHeaders = getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization")

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

// getEnvAsInt reads an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	return defaultValue
}
