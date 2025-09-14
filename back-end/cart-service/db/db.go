package db

import (
	"log"

	"cart-service/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init() error {
	var err error
	DB, err = sqlx.Connect("postgres", config.LoadConfig().DbUrl)
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
	-- Create cart_items table
	CREATE TABLE IF NOT EXISTS cart_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		product_id UUID NOT NULL,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, product_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);

	-- Create orders table
	CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		stripe_payment_id VARCHAR(255),
		total_amount DECIMAL(10,2) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled', 'refunded')),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Create order_items table
	CREATE TABLE IF NOT EXISTS order_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id UUID NOT NULL,
		product_id UUID NOT NULL,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		price DECIMAL(10,2) NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
	CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items(product_id);
	CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
	CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
	CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);

	-- Create function to update updated_at timestamp
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	-- Create triggers to automatically update updated_at
	CREATE OR REPLACE TRIGGER update_cart_items_updated_at
		BEFORE UPDATE ON cart_items
		FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

	CREATE OR REPLACE TRIGGER update_orders_updated_at
		BEFORE UPDATE ON orders
		FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`

	if _, err := DB.Exec(migrationSQL); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations completed successfully")
}
