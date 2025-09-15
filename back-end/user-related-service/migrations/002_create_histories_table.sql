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