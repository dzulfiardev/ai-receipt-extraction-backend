package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/domain"
	"github.com/google/uuid"
)

type ReceiptRepository interface {
	Create(receipt *domain.Receipt) error
	FindByID(id int) (*domain.Receipt, error)
	FindByUUID(uuid string) (*domain.Receipt, error)
	FindByUserID(userID int, page, limit int) ([]domain.Receipt, int64, error)
	Update(receipt *domain.Receipt) error
	Delete(id int) error
	GetStatsByUserID(userID int) (map[string]interface{}, error)
}

type receiptRepository struct {
	db *sql.DB
}

// NewReceiptRepository creates a new receipt repository
func NewReceiptRepository(db *sql.DB) ReceiptRepository {
	return &receiptRepository{db: db}
}

// Create creates a new receipt
func (r *receiptRepository) Create(receipt *domain.Receipt) error {
	query := `
		INSERT INTO receipts (
			user_id, store_name, address, phone, date, image_url, original_filename, 
			file_size, status, total_items, total_spending, total_discount, 
			created_at_unix, updated_at_unix
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, uuid, upload_date, created_at, updated_at
	`

	now := time.Now().Unix()
	err := r.db.QueryRow(
		query,
		receipt.UserID,
		receipt.StoreName,
		receipt.Address,
		receipt.Phone,
		receipt.Date,
		receipt.ImageURL,
		receipt.OriginalFilename,
		receipt.FileSize,
		receipt.Status,
		receipt.TotalItems,
		receipt.TotalSpending,
		receipt.TotalDiscount,
		now,
		now,
	).Scan(&receipt.ID, &receipt.UUID, &receipt.UploadDate, &receipt.CreatedAt, &receipt.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create receipt: %w", err)
	}

	receipt.CreatedAtUnix = now
	receipt.UpdatedAtUnix = now

	return nil
}

// FindByID finds receipt by ID
func (r *receiptRepository) FindByID(id int) (*domain.Receipt, error) {
	query := `
		SELECT id, uuid, user_id, store_name, address, phone, date, image_url, 
		       original_filename, file_size, upload_date, status, total_items, 
		       total_spending, total_discount, created_at, updated_at, 
		       created_at_unix, updated_at_unix
		FROM receipts
		WHERE id = $1
	`

	receipt := &domain.Receipt{}
	err := r.db.QueryRow(query, id).Scan(
		&receipt.ID,
		&receipt.UUID,
		&receipt.UserID,
		&receipt.StoreName,
		&receipt.Address,
		&receipt.Phone,
		&receipt.Date,
		&receipt.ImageURL,
		&receipt.OriginalFilename,
		&receipt.FileSize,
		&receipt.UploadDate,
		&receipt.Status,
		&receipt.TotalItems,
		&receipt.TotalSpending,
		&receipt.TotalDiscount,
		&receipt.CreatedAt,
		&receipt.UpdatedAt,
		&receipt.CreatedAtUnix,
		&receipt.UpdatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("receipt not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find receipt: %w", err)
	}

	return receipt, nil
}

// FindByUUID finds receipt by UUID
func (r *receiptRepository) FindByUUID(uuidStr string) (*domain.Receipt, error) {
	query := `
		SELECT id, uuid, user_id, store_name, address, phone, date, image_url, 
		       original_filename, file_size, upload_date, status, total_items, 
		       total_spending, total_discount, created_at, updated_at, 
		       created_at_unix, updated_at_unix
		FROM receipts
		WHERE uuid = $1
	`

	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	receipt := &domain.Receipt{}
	err = r.db.QueryRow(query, uid).Scan(
		&receipt.ID,
		&receipt.UUID,
		&receipt.UserID,
		&receipt.StoreName,
		&receipt.Address,
		&receipt.Phone,
		&receipt.Date,
		&receipt.ImageURL,
		&receipt.OriginalFilename,
		&receipt.FileSize,
		&receipt.UploadDate,
		&receipt.Status,
		&receipt.TotalItems,
		&receipt.TotalSpending,
		&receipt.TotalDiscount,
		&receipt.CreatedAt,
		&receipt.UpdatedAt,
		&receipt.CreatedAtUnix,
		&receipt.UpdatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("receipt not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find receipt: %w", err)
	}

	return receipt, nil
}

// FindByUserID finds receipts by user ID with pagination
func (r *receiptRepository) FindByUserID(userID int, page, limit int) ([]domain.Receipt, int64, error) {
	// Count total
	var total int64
	countQuery := `SELECT COUNT(*) FROM receipts WHERE user_id = $1`
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count receipts: %w", err)
	}

	// Get receipts
	offset := (page - 1) * limit
	query := `
		SELECT id, uuid, user_id, store_name, address, phone, date, image_url, 
		       original_filename, file_size, upload_date, status, total_items, 
		       total_spending, total_discount, created_at, updated_at, 
		       created_at_unix, updated_at_unix
		FROM receipts
		WHERE user_id = $1
		ORDER BY upload_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query receipts: %w", err)
	}
	defer rows.Close()

	var receipts []domain.Receipt
	for rows.Next() {
		var receipt domain.Receipt
		err := rows.Scan(
			&receipt.ID,
			&receipt.UUID,
			&receipt.UserID,
			&receipt.StoreName,
			&receipt.Address,
			&receipt.Phone,
			&receipt.Date,
			&receipt.ImageURL,
			&receipt.OriginalFilename,
			&receipt.FileSize,
			&receipt.UploadDate,
			&receipt.Status,
			&receipt.TotalItems,
			&receipt.TotalSpending,
			&receipt.TotalDiscount,
			&receipt.CreatedAt,
			&receipt.UpdatedAt,
			&receipt.CreatedAtUnix,
			&receipt.UpdatedAtUnix,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan receipt: %w", err)
		}
		receipts = append(receipts, receipt)
	}

	return receipts, total, nil
}

// Update updates receipt
func (r *receiptRepository) Update(receipt *domain.Receipt) error {
	query := `
		UPDATE receipts
		SET store_name = $1, address = $2, phone = $3, date = $4, status = $5, 
		    total_items = $6, total_spending = $7, total_discount = $8, 
		    updated_at = NOW(), updated_at_unix = $9
		WHERE id = $10
		RETURNING updated_at
	`

	now := time.Now().Unix()
	err := r.db.QueryRow(
		query,
		receipt.StoreName,
		receipt.Address,
		receipt.Phone,
		receipt.Date,
		receipt.Status,
		receipt.TotalItems,
		receipt.TotalSpending,
		receipt.TotalDiscount,
		now,
		receipt.ID,
	).Scan(&receipt.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update receipt: %w", err)
	}

	receipt.UpdatedAtUnix = now
	return nil
}

// Delete deletes receipt by ID
func (r *receiptRepository) Delete(id int) error {
	query := `DELETE FROM receipts WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete receipt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("receipt not found")
	}

	return nil
}

// Soft Delete receipt by ID
func (r *receiptRepository) SoftDelete(id int) error {
	query := `
		UPDATE receipts
		SET status = 'deleted', updated_at = NOW(), updated_at_unix = $1
		WHERE id = $2
	`

	now := time.Now().Unix()
	result, err := r.db.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete receipt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("receipt not found")
	}

	return nil
}

// GetStatsByUserID gets spending statistics by user ID
func (r *receiptRepository) GetStatsByUserID(userID int) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_receipts,
			COALESCE(SUM(total_spending), 0) as total_spending,
			COALESCE(SUM(total_discount), 0) as total_discount,
			COALESCE(AVG(total_spending), 0) as avg_spending
		FROM receipts
		WHERE user_id = $1 AND status = 'completed'
	`

	var totalReceipts int
	var totalSpending, totalDiscount, avgSpending float64

	err := r.db.QueryRow(query, userID).Scan(
		&totalReceipts,
		&totalSpending,
		&totalDiscount,
		&avgSpending,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	stats := map[string]interface{}{
		"total_receipts":   totalReceipts,
		"total_spending":   totalSpending,
		"total_discount":   totalDiscount,
		"average_spending": avgSpending,
		"net_spending":     totalSpending - totalDiscount,
	}

	return stats, nil
}
