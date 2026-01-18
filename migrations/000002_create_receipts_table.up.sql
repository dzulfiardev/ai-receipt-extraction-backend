-- Receipts table
CREATE TABLE receipts (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    store_name VARCHAR(255),
    address VARCHAR(255),
    phone BIGINT,
    date DATE,
    image_url TEXT NOT NULL,
    original_filename VARCHAR(255),
    file_size INTEGER,
    upload_date TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) DEFAULT 'pending',
    total_items INTEGER DEFAULT 0,
    total_spending DECIMAL(15, 2),
    total_discount DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at_unix INTEGER NOT NULL,
    updated_at_unix INTEGER NOT NULL
);

-- Indexes
CREATE INDEX idx_receipts_uuid ON receipts(uuid);
CREATE INDEX idx_receipts_user_id ON receipts(user_id);
CREATE INDEX idx_receipts_upload_date ON receipts(upload_date DESC);
CREATE INDEX idx_receipts_date ON receipts(date DESC);
CREATE INDEX idx_receipts_status ON receipts(status);

-- Comments
COMMENT ON TABLE receipts IS 'Payment receipts uploaded by users';
COMMENT ON COLUMN receipts.status IS 'Status:  pending, processing, completed, failed';
COMMENT ON COLUMN receipts.total_discount IS 'Total discount amount from receipt';