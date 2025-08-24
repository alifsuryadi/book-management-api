package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Environment string
	JWTSecret   string
	BasicAuth   struct {
		Username string
		Password string
	}
}

func Load() *Config {
	// Load .env file if it exists (ignore errors for production deployments)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/book_management?sslmode=disable"),
		Environment: getEnv("ENVIRONMENT", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
	}

	cfg.BasicAuth.Username = getEnv("BASIC_AUTH_USERNAME", "admin")
	cfg.BasicAuth.Password = getEnv("BASIC_AUTH_PASSWORD", "password")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}