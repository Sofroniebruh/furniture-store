package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"user-related-service/config"
)

var DB *sqlx.DB

func Init() error {
	var err error
	DB, err = sqlx.Connect("postgres", config.DB_URL)
	runMigrations()
	return err
}

func runMigrations() {
	migrationSQL := `
	-- Create histories table
	CREATE TABLE IF NOT EXISTS histories (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		product_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	-- Create indexes for faster queries
	CREATE INDEX IF NOT EXISTS idx_histories_user_id ON histories(user_id);
	CREATE INDEX IF NOT EXISTS idx_histories_product_id ON histories(product_id);
	CREATE INDEX IF NOT EXISTS idx_histories_created_at ON histories(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_histories_user_created ON histories(user_id, created_at DESC);

	-- Create wishlists table
	CREATE TABLE IF NOT EXISTS wishlists (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		product_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(user_id, product_id)
	);

	-- Create indexes for faster queries
	CREATE INDEX IF NOT EXISTS idx_wishlists_user_id ON wishlists(user_id);
	CREATE INDEX IF NOT EXISTS idx_wishlists_product_id ON wishlists(product_id);
	CREATE INDEX IF NOT EXISTS idx_wishlists_created_at ON wishlists(created_at DESC);
	`

	if _, err := DB.Exec(migrationSQL); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations completed successfully")
}
