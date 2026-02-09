package main

import (
	"fmt"
	"log"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/config"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"env":    cfg.Environment,
		})
	})

	// API v1 routes
	v1 := e.Group("/api/v1")

	{
		v1.GET("/ping", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"message": "pong"})
		})
	}

	// Start server
	address := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Server starting on %s", address)
	if err := e.Start(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
