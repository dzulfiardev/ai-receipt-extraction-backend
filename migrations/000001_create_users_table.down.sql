-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_uuid;

-- Drop table
DROP TABLE IF EXISTS users CASCADE;