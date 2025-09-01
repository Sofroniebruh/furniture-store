package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	Port                string
	JWTSecret           string
	StripeSecretKey     string
	StripePublishableKey string
	StripeWebhookSecret string
	UserRelatedServiceURL string
}

type ContextKey string

const UserIdKey ContextKey = "userId"

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5433"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", ""),
		DBName:              getEnv("DB_NAME", "furniture_store"),
		Port:                getEnv("PORT", "8083"),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		UserRelatedServiceURL: getEnv("USER_RELATED_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}