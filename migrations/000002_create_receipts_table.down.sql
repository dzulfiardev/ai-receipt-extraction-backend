-- Drop indexes
DROP INDEX IF EXISTS idx_receipts_status;
DROP INDEX IF EXISTS idx_receipts_date;
DROP INDEX IF EXISTS idx_receipts_upload_date;
DROP INDEX IF EXISTS idx_receipts_user_id;
DROP INDEX IF EXISTS idx_receipts_uuid;

-- Drop table
DROP TABLE IF EXISTS receipts CASCADE;