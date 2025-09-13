package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"products-service/config"
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
	CREATE TABLE IF NOT EXISTS colors (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(100) UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
		price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
		picture_urls TEXT[] DEFAULT '{}',
		event VARCHAR(100),
		model VARCHAR(100),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS product_colors (
		product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		color_id UUID NOT NULL REFERENCES colors(id) ON DELETE CASCADE,
		PRIMARY KEY (product_id, color_id)
	);

	CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
	CREATE INDEX IF NOT EXISTS idx_products_event ON products(event);
	CREATE INDEX IF NOT EXISTS idx_products_model ON products(model);
	CREATE INDEX IF NOT EXISTS idx_products_stock ON products(stock);
	CREATE INDEX IF NOT EXISTS idx_colors_name ON colors(name);

	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	CREATE OR REPLACE TRIGGER update_products_updated_at
		BEFORE UPDATE ON products
		FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`
	
	if _, err := DB.Exec(migrationSQL); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	
	log.Println("Products service migrations completed successfully")
}
