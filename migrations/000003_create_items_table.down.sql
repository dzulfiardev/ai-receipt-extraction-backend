-- Drop indexes
DROP INDEX IF EXISTS idx_items_name;
DROP INDEX IF EXISTS idx_items_receipt_id;
DROP INDEX IF EXISTS idx_items_uuid;

-- Drop table
DROP TABLE IF EXISTS items CASCADE;