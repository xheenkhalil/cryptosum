package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	Env            string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	AllowedOrigins string
	GeminiAPIKey   string
}

var AppConfig *Config

// LoadConfig parses environment variables and enforces required fields.
func LoadConfig() {
	// Load .env file if it exists (mostly for local development)
	_ = godotenv.Load()

	AppConfig = &Config{
		Port:           getEnvOrDefault("PORT", "8080"),
		Env:            getEnvOrDefault("ENV", "development"),
		DatabaseURL:    getEnvOrDefault("DATABASE_URL", "host=127.0.0.1 user=postgres password=password dbname=postgres port=15432 sslmode=disable"),
		RedisURL:       getEnvOrDefault("REDIS_URL", "127.0.0.1:6379"),
		JWTSecret:      getEnvOrDefault("JWT_SECRET", "super_secret_fallback_key"), // Should fail in production if empty
		AllowedOrigins: getEnvOrDefault("ALLOWED_ORIGINS", "*"),
		GeminiAPIKey:   os.Getenv("GEMINI_API_KEY"),
	}

	// Fail fast in production if critical secrets are missing
	if AppConfig.Env == "production" {
		if os.Getenv("JWT_SECRET") == "" {
			log.Fatal("FATAL: JWT_SECRET is required in production environment")
		}
		if os.Getenv("DATABASE_URL") == "" {
			log.Fatal("FATAL: DATABASE_URL is required in production environment")
		}
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
