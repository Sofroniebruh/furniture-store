package db

import (
	"auth-service/config"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init() error {
	var err error
	DB, err = sqlx.Connect("postgres", config.DB_URL)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Successfully connected to database")
	
	runMigrations()
	return nil
}

func runMigrations() {
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		roles TEXT[] DEFAULT ARRAY['user'],
		password VARCHAR(255) NOT NULL,
		refresh_token TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		is_email_verified BOOLEAN DEFAULT FALSE
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users(refresh_token);
	`
	
	if _, err := DB.Exec(migrationSQL); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	
	log.Println("Auth service migrations completed successfully")
}
