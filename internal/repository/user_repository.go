package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id int) (*domain.User, error)
	FindByUUID(uuid string) (*domain.User, error)
	Update(user *domain.User) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, created_at_unix, updated_at_unix)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, uuid, created_at, updated_at
	`

	now := time.Now().Unix()
	err := r.db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		now,
		now,
	).Scan(&user.ID, &user.UUID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAtUnix = now
	user.UpdatedAtUnix = now

	return nil
}

// FindByEmail finds user by email
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, uuid, email, password_hash, full_name, created_at, updated_at, created_at_unix, updated_at_unix
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedAtUnix,
		&user.UpdatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByID finds user by ID
func (r *userRepository) FindByID(id int) (*domain.User, error) {
	query := `
		SELECT id, uuid, email, password_hash, full_name, created_at, updated_at, created_at_unix, updated_at_unix
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedAtUnix,
		&user.UpdatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByUUID finds user by UUID
func (r *userRepository) FindByUUID(uuidStr string) (*domain.User, error) {
	query := `
		SELECT id, uuid, email, password_hash, full_name, created_at, updated_at, created_at_unix, updated_at_unix
		FROM users
		WHERE uuid = $1
	`

	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid:  %w", err)
	}

	user := &domain.User{}
	err = r.db.QueryRow(query, uid).Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedAtUnix,
		&user.UpdatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// Update updates user information
func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, updated_at = NOW(), updated_at_unix = $3
		WHERE id = $4
		RETURNING updated_at
	`

	now := time.Now().Unix()
	err := r.db.QueryRow(query, user.Email, user.FullName, now, user.ID).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.UpdatedAtUnix = now
	return nil
}
