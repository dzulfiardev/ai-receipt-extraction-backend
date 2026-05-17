package main

import (
	"fmt"
	"log"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/config"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/database"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/handler"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/middleware"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/repository"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/service"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dbCfg := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	db, err := database.NewPostgresDB(dbCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	receiptRepo := repository.NewReceiptRepository(db)
	itemRepo := repository.NewItemRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpireHours)
	receiptService := service.NewReceiptService(receiptRepo, itemRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	receiptHandler := handler.NewReceiptHandler(receiptService)

	// Create Echo instance
	e := echo.New()

	// Global middleware
	// e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"env":    cfg.Environment,
		})
	})

	// API v1 routes
	v1 := e.Group("/api/v1")

	// Public routes (no authentication)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes (require authentication)
	protected := v1.Group("")
	protected.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	{
		// Auth routes
		protected.GET("/auth/me", authHandler.GetMe)

		// Receipt routes
		receipts := protected.Group("/receipts")
		{
			receipts.POST("", receiptHandler.CreateReceipt)
			receipts.GET("", receiptHandler.GetReceipts)
			receipts.GET("/:id", receiptHandler.GetReceipt)
			receipts.PUT("/:id", receiptHandler.UpdateReceipt)
			receipts.DELETE("/:id", receiptHandler.DeleteReceipt)
			receipts.GET("/stats", receiptHandler.GetStats)
		}
	}

	// Start server
	address := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Server starting on %s", address)
	log.Printf("📝 Environment: %s", cfg.Environment)
	log.Printf("🔐 JWT Expiration: %d hours", cfg.JWTExpireHours)

	if err := e.Start(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
