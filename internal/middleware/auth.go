package middleware

import (
	"strings"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates JWT token
func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.ErrorResponse(c, 401, "Missing authorization header")
			}

			// Extract token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return utils.ErrorResponse(c, 401, "Invalid authorization header format")
			}

			tokenString := parts[1]

			// Validate token
			claims, err := utils.ValidateJWT(tokenString, jwtSecret)
			if err != nil {
				return utils.ErrorResponse(c, 401, "Invalid or expired token")
			}

			// Set user info in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)

			return next(c)
		}
	}
}

// GetUserID extracts user ID from context
func GetUserID(c echo.Context) int {
	return c.Get("user_id").(int)
}
