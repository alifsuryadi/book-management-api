package config

import (
	"os"
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
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/book_management?sslmode=disable"),
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