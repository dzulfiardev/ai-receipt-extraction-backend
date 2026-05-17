package handler

import (
	"math"
	"strconv"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/domain"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/middleware"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/service"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

type ReceiptHandler struct {
	receiptService service.ReceiptService
	validator      *utils.Validator
}

// NewReceiptHandler creates a new receipt handler
func NewReceiptHandler(receiptService service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{
		receiptService: receiptService,
		validator:      utils.NewValidator(),
	}
}

// CreateReceipt handles receipt creation
func (h *ReceiptHandler) CreateReceipt(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req domain.CreateReceiptRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	// For now, we'll use a dummy image URL
	// In production, you would upload the image to S3/MinIO first
	imageURL := "https://storage.example.com/receipts/sample. jpg"
	filename := "sample.jpg"
	fileSize := 0

	receipt, err := h.receiptService.CreateReceipt(userID, req, imageURL, filename, fileSize)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 201, "Receipt created successfully", receipt)
}

// GetReceipt gets receipt by ID
func (h *ReceiptHandler) GetReceipt(c echo.Context) error {
	userID := middleware.GetUserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid receipt ID")
	}

	receipt, err := h.receiptService.GetReceiptByID(id, userID)
	if err != nil {
		return utils.ErrorResponse(c, 404, err.Error())
	}

	return utils.SuccessResponse(c, 200, "Receipt retrieved successfully", receipt)
}

// GetReceipts gets all receipts for current user
func (h *ReceiptHandler) GetReceipts(c echo.Context) error {
	userID := middleware.GetUserID(c)

	// Parse pagination params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	receipts, total, err := h.receiptService.GetReceiptsByUserID(userID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, 500, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	meta := utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return utils.PaginatedSuccessResponse(c, 200, receipts, meta)
}

// UpdateReceipt updates receipt
func (h *ReceiptHandler) UpdateReceipt(c echo.Context) error {
	userID := middleware.GetUserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid receipt ID")
	}

	var req domain.CreateReceiptRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	receipt, err := h.receiptService.UpdateReceipt(id, userID, req)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "Receipt updated successfully", receipt)
}

// DeleteReceipt deletes receipt
func (h *ReceiptHandler) DeleteReceipt(c echo.Context) error {
	userID := middleware.GetUserID(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid receipt ID")
	}

	if err := h.receiptService.DeleteReceipt(id, userID); err != nil {
		return utils.ErrorResponse(c, 404, err.Error())
	}

	return utils.SuccessResponse(c, 200, "Receipt deleted successfully", nil)
}

// GetStats gets spending statistics
func (h *ReceiptHandler) GetStats(c echo.Context) error {
	userID := middleware.GetUserID(c)

	stats, err := h.receiptService.GetStatsByUserID(userID)
	if err != nil {
		return utils.ErrorResponse(c, 500, err.Error())
	}

	return utils.SuccessResponse(c, 200, "Statistics retrieved successfully", stats)
}
