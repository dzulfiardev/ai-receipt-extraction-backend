package service

import (
	"database/sql"
	"fmt"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/domain"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/repository"
)

type ReceiptService interface {
	CreateReceipt(userID int, req domain.CreateReceiptRequest, imageURL string, filename string, fileSize int) (*domain.ReceiptWithItems, error)
	GetReceiptByID(id int, userID int) (*domain.ReceiptWithItems, error)
	GetReceiptsByUserID(userID int, page, limit int) ([]domain.Receipt, int64, error)
	UpdateReceipt(id int, userID int, req domain.CreateReceiptRequest) (*domain.ReceiptWithItems, error)
	DeleteReceipt(id int, userID int) error
	GetStatsByUserID(userID int) (map[string]interface{}, error)
}

type receiptService struct {
	receiptRepo repository.ReceiptRepository
	itemRepo    repository.ItemRepository
}

// NewReceiptService creates a new receipt service
func NewReceiptService(receiptRepo repository.ReceiptRepository, itemRepo repository.ItemRepository) ReceiptService {
	return &receiptService{
		receiptRepo: receiptRepo,
		itemRepo:    itemRepo,
	}
}

// CreateReceipt creates a new receipt with items
func (s *receiptService) CreateReceipt(userID int, req domain.CreateReceiptRequest, imageURL string, filename string, fileSize int) (*domain.ReceiptWithItems, error) {
	// Parse date if provided
	var date sql.NullTime
	if req.Date != nil && *req.Date != "" {
		// You can parse the date here if needed
		// For now, we'll leave it as null
	}

	// Create receipt
	receipt := &domain.Receipt{
		UserID:           userID,
		StoreName:        sql.NullString{String: req.StoreName, Valid: req.StoreName != ""},
		Address:          sql.NullString{String: req.Address, Valid: req.Address != ""},
		Phone:            sql.NullInt64{Int64: *req.Phone, Valid: req.Phone != nil},
		Date:             date,
		ImageURL:         imageURL,
		OriginalFilename: filename,
		FileSize:         fileSize,
		Status:           domain.StatusCompleted,
		TotalItems:       req.TotalItems,
		TotalSpending:    req.TotalSpending,
		TotalDiscount:    req.TotalDiscount,
	}

	if err := s.receiptRepo.Create(receipt); err != nil {
		return nil, fmt.Errorf("failed to create receipt:  %w", err)
	}

	// Create items
	var items []domain.Item
	for _, itemReq := range req.Items {
		item := domain.Item{
			ReceiptID: receipt.ID,
			Name:      itemReq.Name,
			UnitPrice: itemReq.UnitPrice,
			Quantity:  itemReq.Quantity,
			Price:     itemReq.Price,
			Total:     itemReq.Total,
		}
		items = append(items, item)
	}

	if len(items) > 0 {
		if err := s.itemRepo.CreateBatch(items); err != nil {
			return nil, fmt.Errorf("failed to create items:  %w", err)
		}
	}

	return &domain.ReceiptWithItems{
		Receipt: *receipt,
		Items:   items,
	}, nil
}

// GetReceiptByID gets receipt by ID with items
func (s *receiptService) GetReceiptByID(id int, userID int) (*domain.ReceiptWithItems, error) {
	receipt, err := s.receiptRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if receipt.UserID != userID {
		return nil, fmt.Errorf("unauthorized access")
	}

	// Get items
	items, err := s.itemRepo.FindByReceiptID(receipt.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return &domain.ReceiptWithItems{
		Receipt: *receipt,
		Items:   items,
	}, nil
}

// GetReceiptsByUserID gets all receipts for a user with pagination
func (s *receiptService) GetReceiptsByUserID(userID int, page, limit int) ([]domain.Receipt, int64, error) {
	return s.receiptRepo.FindByUserID(userID, page, limit)
}

// UpdateReceipt updates receipt and items
func (s *receiptService) UpdateReceipt(id int, userID int, req domain.CreateReceiptRequest) (*domain.ReceiptWithItems, error) {
	// Get existing receipt
	receipt, err := s.receiptRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if receipt.UserID != userID {
		return nil, fmt.Errorf("unauthorized access")
	}

	// Update receipt
	receipt.StoreName = sql.NullString{String: req.StoreName, Valid: req.StoreName != ""}
	receipt.Address = sql.NullString{String: req.Address, Valid: req.Address != ""}
	receipt.Phone = sql.NullInt64{Int64: *req.Phone, Valid: req.Phone != nil}
	receipt.TotalItems = req.TotalItems
	receipt.TotalSpending = req.TotalSpending
	receipt.TotalDiscount = req.TotalDiscount

	if err := s.receiptRepo.Update(receipt); err != nil {
		return nil, fmt.Errorf("failed to update receipt: %w", err)
	}

	// Get updated items
	items, err := s.itemRepo.FindByReceiptID(receipt.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return &domain.ReceiptWithItems{
		Receipt: *receipt,
		Items:   items,
	}, nil
}

// DeleteReceipt deletes receipt
func (s *receiptService) DeleteReceipt(id int, userID int) error {
	// Get receipt to check ownership
	receipt, err := s.receiptRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check ownership
	if receipt.UserID != userID {
		return fmt.Errorf("unauthorized access")
	}

	return s.receiptRepo.Delete(id)
}

// GetStatsByUserID gets spending statistics
func (s *receiptService) GetStatsByUserID(userID int) (map[string]interface{}, error) {
	return s.receiptRepo.GetStatsByUserID(userID)
}
