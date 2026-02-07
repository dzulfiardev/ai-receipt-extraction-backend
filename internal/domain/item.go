package domain

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID            int       `json:"id" db:"id"`
	UUID          uuid.UUID `json:"uuid" db:"uuid"`
	ReceiptID     int       `json:"receipt_id" db:"receipt_id"`
	Name          string    `json:"name" db:"name"`
	UnitPrice     int       `json:"unit_price" db:"unit_price"`
	Quantity      int       `json:"quantity" db:"quantity"`
	Price         int       `json:"price" db:"price"`
	Total         int       `json:"total" db:"total"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix" db:"created_at_unix"`
}

// CreateItemRequest represents item creation request
type CreateItemRequest struct {
	Name      string `json:"name" validate:"required"`
	UnitPrice int    `json:"unit_price" validate:"required,min=0"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
	Price     int    `json:"price" validate:"required,min=0"`
	Total     int    `json:"total" validate:"required,min=0"`
}
