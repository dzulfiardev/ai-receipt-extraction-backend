package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID            int       `json:"id" db:"id"`
	UUID          uuid.UUID `json:"uuid" db:"uuid"`
	UserID        int       `json:"user_id" db:"user_id"`
	TokenHash     string    `json:"-" db:"token_hash"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix" db:"created_at_unix"`
}

// IsExpired checks if session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
