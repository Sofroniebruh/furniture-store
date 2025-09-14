package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUrl                 string
	Port                  string
	JWTSecret             string
	StripeSecretKey       string
	StripePublishableKey  string
	StripeWebhookSecret   string
	UserRelatedServiceURL string
}

type ContextKey string

const UserIdKey ContextKey = "userId"

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DbUrl:                 getEnv("DATABASE_URL", ""),
		JWTSecret:             getEnv("JWT_SECRET", ""),
		StripeSecretKey:       getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishableKey:  getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		StripeWebhookSecret:   getEnv("STRIPE_WEBHOOK_SECRET", ""),
		UserRelatedServiceURL: getEnv("USER_RELATED_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
