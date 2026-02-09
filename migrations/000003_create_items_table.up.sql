-- Items table
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    receipt_id INTEGER NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    unit_price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price INTEGER NOT NULL,
    total INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_at_unix INTEGER NOT NULL
);

-- Indexes
CREATE INDEX idx_items_uuid ON items(uuid);
CREATE INDEX idx_items_receipt_id ON items(receipt_id);
CREATE INDEX idx_items_name ON items(name);

-- Comments
COMMENT ON TABLE items IS 'Individual items/products from receipts';
COMMENT ON COLUMN items.unit_price IS 'Price per single unit';
COMMENT ON COLUMN items.price IS 'Price shown on receipt (may differ from unit_price)';
COMMENT ON COLUMN items.total IS 'Total price for this item (price * quantity)';