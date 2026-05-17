package handler

import (
	"github.com/dzulfiardev/receipt-extraction-backend/internal/domain"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/middleware"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/service"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *utils.Validator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   utils.NewValidator(),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var req domain.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 201, "User registered successfully", user.ToResponse())
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	token, user, err := h.authService.Login(req)
	if err != nil {
		return utils.ErrorResponse(c, 401, err.Error())
	}

	return utils.SuccessResponse(c, 200, "Login successful", map[string]interface{}{
		"token": token,
		"user":  user.ToResponse(),
	})
}

// GetMe gets current user info
func (h *AuthHandler) GetMe(c echo.Context) error {
	userID := middleware.GetUserID(c)

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	return utils.SuccessResponse(c, 200, "User retrieved successfully", user.ToResponse())
}
