-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at_unix INTEGER NOT NULL,
    updated_at_unix INTEGER NOT NULL
);

-- Indexes
CREATE INDEX idx_users_uuid ON users(uuid);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_fullname ON users(full_name);

-- Comments
COMMENT ON TABLE users IS 'User accounts table';
COMMENT ON COLUMN users.uuid IS 'Public UUID for external reference';
COMMENT ON COLUMN users.created_at_unix IS 'Unix timestamp for created_at';