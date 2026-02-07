package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            int       `json:"id" db:"id"`
	UUID          uuid.UUID `json:"uuid" db:"uuid"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	FullName      string    `json:"full_name" db:"full_name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	CreatedAtUnix int64     `json:"created_at_unix" db:"created_at_unix"`
	UpdatedAtUnix int64     `json:"updated_at_unix" db:"updated_at_unix"`
}

// CreateUserRequest represents user registration request
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents user response (without sensitive data)
type UserResponse struct {
	ID            int       `json:"id"`
	UUID          string    `json:"uuid"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		UUID:          u.UUID.String(),
		Email:         u.Email,
		FullName:      u.FullName,
		CreatedAt:     u.CreatedAt,
		CreatedAtUnix: u.CreatedAtUnix,
	}
}
