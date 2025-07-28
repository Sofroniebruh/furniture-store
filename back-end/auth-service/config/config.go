package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var (
	DB_URL            string
	JWT_SECRET        []byte
	ACCESS_TOKEN_TTL  = time.Minute * 15
	REFRESH_TOKEN_TTL = time.Hour * 24 * 7
)

type userId string

const UserIdKey userId = "userId"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	DB_URL = os.Getenv("DATABASE_URL")
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
}
