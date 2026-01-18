-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_token_hash;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_uuid;

-- Drop table
DROP TABLE IF EXISTS sessions CASCADE;

-- Drop extension (only if no other tables use it)
DROP EXTENSION IF EXISTS "pgcrypto";