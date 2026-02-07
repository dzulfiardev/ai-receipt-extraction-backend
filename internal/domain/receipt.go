package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ReceiptStatus string

const (
	StatusPending    ReceiptStatus = "pending"
	StatusProcessing ReceiptStatus = "processing"
	StatusCompleted  ReceiptStatus = "completed"
	StatusFailed     ReceiptStatus = "failed"
)

type Receipt struct {
	ID               int            `json:"id" db:"id"`
	UUID             uuid.UUID      `json:"uuid" db:"uuid"`
	UserID           int            `json:"user_id" db:"user_id"`
	StoreName        sql.NullString `json:"store_name" db:"store_name"`
	Address          sql.NullString `json:"address" db:"address"`
	Phone            sql.NullInt64  `json:"phone" db:"phone"`
	Date             sql.NullTime   `json:"date" db:"date"`
	ImageURL         string         `json:"image_url" db:"image_url"`
	OriginalFilename string         `json:"original_filename" db:"original_filename"`
	FileSize         int            `json:"file_size" db:"file_size"`
	UploadDate       time.Time      `json:"upload_date" db:"upload_date"`
	Status           ReceiptStatus  `json:"status" db:"status"`
	TotalItems       int            `json:"total_items" db:"total_items"`
	TotalSpending    float64        `json:"total_spending" db:"total_spending"`
	TotalDiscount    float64        `json:"total_discount" db:"total_discount"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
	CreatedAtUnix    int64          `json:"created_at_unix" db:"created_at_unix"`
	UpdatedAtUnix    int64          `json:"updated_at_unix" db:"updated_at_unix"`
}

// ReceiptWithItems represents receipt with its items
type ReceiptWithItems struct {
	Receipt
	Items []Item `json:"items"`
}

// CreateReceiptRequest represents receipt creation request
type CreateReceiptRequest struct {
	StoreName     string              `json:"store_name"`
	Address       string              `json:"address"`
	Phone         *int64              `json:"phone"`
	Date          *string             `json:"date"`
	TotalItems    int                 `json:"total_items"`
	TotalSpending float64             `json:"total_spending"`
	TotalDiscount float64             `json:"total_discount"`
	Items         []CreateItemRequest `json:"items"`
}
