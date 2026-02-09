package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/domain"
)

type ItemRepository interface {
	Create(item *domain.Item) error
	CreateBatch(items []domain.Item) error
	FindByReceiptID(receiptID int) ([]domain.Item, error)
	FindByID(id int) (*domain.Item, error)
	Update(item *domain.Item) error
	Delete(id int) error
}

type itemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

// Create creates a new item
func (r *itemRepository) Create(item *domain.Item) error {
	query := `
		INSERT INTO items (receipt_id, name, unit_price, quantity, price, total, created_at_unix)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, created_at
	`

	now := time.Now().Unix()
	err := r.db.QueryRow(
		query,
		item.ReceiptID,
		item.Name,
		item.UnitPrice,
		item.Quantity,
		item.Price,
		item.Total,
		now,
	).Scan(&item.ID, &item.UUID, &item.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create item:  %w", err)
	}

	item.CreatedAtUnix = now
	return nil
}

// CreateBatch creates multiple items in a single transaction
func (r *itemRepository) CreateBatch(items []domain.Item) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO items (receipt_id, name, unit_price, quantity, price, total, created_at_unix)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, created_at
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now().Unix()

	for i := range items {
		err := stmt.QueryRow(
			items[i].ReceiptID,
			items[i].Name,
			items[i].UnitPrice,
			items[i].Quantity,
			items[i].Price,
			items[i].Total,
			now,
		).Scan(&items[i].ID, &items[i].UUID, &items[i].CreatedAt)

		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}

		items[i].CreatedAtUnix = now
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// FindByReceiptID finds all items for a receipt
func (r *itemRepository) FindByReceiptID(receiptID int) ([]domain.Item, error) {
	query := `
		SELECT id, uuid, receipt_id, name, unit_price, quantity, price, total, created_at, created_at_unix
		FROM items
		WHERE receipt_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query, receiptID)

	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}

	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item

		err := rows.Scan(
			&item.ID,
			&item.UUID,
			&item.ReceiptID,
			&item.Name,
			&item.UnitPrice,
			&item.Quantity,
			&item.Price,
			&item.Total,
			&item.CreatedAt,
			&item.CreatedAtUnix,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan item:  %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

// FindByID finds item by ID
func (r *itemRepository) FindByID(id int) (*domain.Item, error) {
	query := `
		SELECT id, uuid, receipt_id, name, unit_price, quantity, price, total, created_at, created_at_unix
		FROM items
		WHERE id = $1
	`

	item := &domain.Item{}
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.UUID,
		&item.ReceiptID,
		&item.Name,
		&item.UnitPrice,
		&item.Quantity,
		&item.Price,
		&item.Total,
		&item.CreatedAt,
		&item.CreatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find item: %w", err)
	}

	return item, nil
}

// Update updates an item
func (r *itemRepository) Update(item *domain.Item) error {
	query := `
		UPDATE items
		SET name = $1, unit_price = $2, quantity = $3, price = $4, total = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(
		query,
		item.Name,
		item.UnitPrice,
		item.Quantity,
		item.Price,
		item.Total,
		item.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// Delete deletes an item
func (r *itemRepository) Delete(id int) error {
	query := `DELETE FROM items WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}
