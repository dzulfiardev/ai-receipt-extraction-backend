-- Sessions table
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_at_unix INTEGER NOT NULL
);

-- Indexes
CREATE INDEX idx_sessions_uuid ON sessions(uuid);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Comments
COMMENT ON TABLE sessions IS 'User authentication sessions (JWT token tracking)';
COMMENT ON COLUMN sessions.token_hash IS 'Hashed JWT token for blacklist functionality';