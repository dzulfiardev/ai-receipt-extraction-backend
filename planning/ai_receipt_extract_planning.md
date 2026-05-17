# Copilot Chat Conversation Export: AI agent payment recipe extraction - Back End

**User:** @dzulfiardev
**Thread URL:** https://github.com/copilot/c/e8664193-4f1e-4cf6-ad93-ff7dcd351e6e

## @dzulfiardev

update schema each column id to Use SERIAL, and use echo framework don't use gorm

## Copilot

Perfect! Let's update the schema to use `SERIAL` and switch to Echo framework without GORM. Here are the updated files:

---

## **Step 1: Updated Migration Files**

### `migrations/000001_create_users_table. up.sql`

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at_unix INTEGER NOT NULL,
    updated_at_unix INTEGER NOT NULL
);

-- Indexes
CREATE INDEX idx_users_uuid ON users(uuid);
CREATE INDEX idx_users_email ON users(email);

-- Comments
COMMENT ON TABLE users IS 'User accounts table';
COMMENT ON COLUMN users.uuid IS 'Public UUID for external reference';
COMMENT ON COLUMN users.created_at_unix IS 'Unix timestamp for created_at';
```

### `migrations/000001_create_users_table.down.sql`

```sql
-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_uuid;

-- Drop table
DROP TABLE IF EXISTS users CASCADE;
```

---

### `migrations/000002_create_receipts_table.up.sql`

```sql
-- Receipts table
CREATE TABLE receipts (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    store_name VARCHAR(255),
    address VARCHAR(255),
    phone BIGINT,
    date DATE,
    image_url TEXT NOT NULL,
    original_filename VARCHAR(255),
    file_size INTEGER,
    upload_date TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) DEFAULT 'pending',
    total_items INTEGER DEFAULT 0,
    total_spending DECIMAL(15, 2),
    total_discount DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at_unix INTEGER NOT NULL,
    updated_at_unix INTEGER NOT NULL
);

-- Indexes
CREATE INDEX idx_receipts_uuid ON receipts(uuid);
CREATE INDEX idx_receipts_user_id ON receipts(user_id);
CREATE INDEX idx_receipts_upload_date ON receipts(upload_date DESC);
CREATE INDEX idx_receipts_date ON receipts(date DESC);
CREATE INDEX idx_receipts_status ON receipts(status);

-- Comments
COMMENT ON TABLE receipts IS 'Payment receipts uploaded by users';
COMMENT ON COLUMN receipts.status IS 'Status:  pending, processing, completed, failed';
COMMENT ON COLUMN receipts.total_discount IS 'Total discount amount from receipt';
```

### `migrations/000002_create_receipts_table.down. sql`

```sql
-- Drop indexes
DROP INDEX IF EXISTS idx_receipts_status;
DROP INDEX IF EXISTS idx_receipts_date;
DROP INDEX IF EXISTS idx_receipts_upload_date;
DROP INDEX IF EXISTS idx_receipts_user_id;
DROP INDEX IF EXISTS idx_receipts_uuid;

-- Drop table
DROP TABLE IF EXISTS receipts CASCADE;
```

---

### `migrations/000003_create_items_table. up.sql`

```sql
-- Items table
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    receipt_id INTEGER NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    unit_price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price INTEGER NOT NULL,
    total INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_at_unix INTEGER NOT NULL
);

-- Indexes
CREATE INDEX idx_items_uuid ON items(uuid);
CREATE INDEX idx_items_receipt_id ON items(receipt_id);
CREATE INDEX idx_items_name ON items(name);

-- Comments
COMMENT ON TABLE items IS 'Individual items/products from receipts';
COMMENT ON COLUMN items.unit_price IS 'Price per single unit';
COMMENT ON COLUMN items.price IS 'Price shown on receipt (may differ from unit_price)';
COMMENT ON COLUMN items.total IS 'Total price for this item (price * quantity)';
```

### `migrations/000003_create_items_table.down. sql`

```sql
-- Drop indexes
DROP INDEX IF EXISTS idx_items_name;
DROP INDEX IF EXISTS idx_items_receipt_id;
DROP INDEX IF EXISTS idx_items_uuid;

-- Drop table
DROP TABLE IF EXISTS items CASCADE;
```

---

### `migrations/000004_create_sessions_table.up.sql`

```sql
-- Sessions table
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_at_unix INTEGER NOT NULL
);

-- Indexes
CREATE INDEX idx_sessions_uuid ON sessions(uuid);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Comments
COMMENT ON TABLE sessions IS 'User authentication sessions (JWT token tracking)';
COMMENT ON COLUMN sessions.token_hash IS 'Hashed JWT token for blacklist functionality';
```

### `migrations/000004_create_sessions_table.down.sql`

```sql
-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_token_hash;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_uuid;

-- Drop table
DROP TABLE IF EXISTS sessions CASCADE;

-- Drop extension (only if no other tables use it)
DROP EXTENSION IF EXISTS "pgcrypto";
```

---

## **Step 2: Update Dependencies**

Update `go.mod` to use Echo instead of Gin/GORM:

```bash
# Remove old dependencies (if any)
go mod tidy

# Install Echo framework
go get -u github.com/labstack/echo/v4
go get -u github.com/labstack/echo/v4/middleware

# Install PostgreSQL driver
go get -u github.com/lib/pq

# Install other dependencies
go get -u github. com/joho/godotenv
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/google/uuid

# Migration packages
go get -u github.com/golang-migrate/migrate/v4
go get -u github. com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```

---

## **Step 3: Updated Project Structure**

```
receipt-backend/
├── cmd/
│   ├── api/
│   │   └── main.go                 # API server entry point
│   └── migrate/
│       └── main.go                 # Migration runner
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration
│   ├── database/
│   │   └── postgres. go             # Database connection pool
│   ├── domain/
│   │   ├── user.go                 # User struct
│   │   ├── receipt.go              # Receipt struct
│   │   ├── item.go                 # Item struct
│   │   └── session. go              # Session struct
│   ├── repository/
│   │   ├── user_repository.go      # User data access
│   │   ├── receipt_repository.go   # Receipt data access
│   │   └── item_repository.go      # Item data access
│   ├── service/
│   │   ├── auth_service.go         # Authentication logic
│   │   └── receipt_service.go      # Receipt processing
│   ├── handler/
│   │   ├── auth_handler.go         # Auth endpoints
│   │   └── receipt_handler.go      # Receipt endpoints
│   ├── middleware/
│   │   ├── auth. go                 # JWT validation
│   │   └── logger.go               # Request logging
│   └── utils/
│       ├── jwt.go                  # JWT helper
│       ├── response.go             # API response helper
│       └── validator.go            # Input validation
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_receipts_table. up.sql
│   ├── 000002_create_receipts_table.down.sql
│   ├── 000003_create_items_table.up.sql
│   ├── 000003_create_items_table.down.sql
│   ├── 000004_create_sessions_table.up.sql
│   └── 000004_create_sessions_table.down.sql
├── .env. example
├── .env
├── . gitignore
├── go. mod
├── go.sum
├── Makefile
└── README.md
```

---

## **Step 4: Create Database Connection (using database/sql)**

### `internal/database/postgres.go`

```go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL connection pool
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database:  %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time. Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt. Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return db, nil
}
```

---

## **Step 5: Create Domain Models**

### `internal/domain/user.go`

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            int       `json:"id" db:"id"`
	UUID          uuid.UUID `json:"uuid" db:"uuid"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	FullName      string    `json:"full_name" db:"full_name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time. Time `json:"updated_at" db:"updated_at"`
	CreatedAtUnix int64     `json:"created_at_unix" db:"created_at_unix"`
	UpdatedAtUnix int64     `json:"updated_at_unix" db:"updated_at_unix"`
}

// CreateUserRequest represents user registration request
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents user response (without sensitive data)
type UserResponse struct {
	ID            int       `json:"id"`
	UUID          string    `json:"uuid"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		UUID:          u.UUID. String(),
		Email:         u.Email,
		FullName:      u.FullName,
		CreatedAt:     u.CreatedAt,
		CreatedAtUnix: u.CreatedAtUnix,
	}
}
```

### `internal/domain/receipt.go`

```go
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
	StoreName     string  `json:"store_name"`
	Address       string  `json:"address"`
	Phone         *int64  `json:"phone"`
	Date          *string `json:"date"`
	TotalItems    int     `json:"total_items"`
	TotalSpending float64 `json:"total_spending"`
	TotalDiscount float64 `json:"total_discount"`
	Items         []CreateItemRequest `json:"items"`
}
```

### `internal/domain/item.go`

```go
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
```

### `internal/domain/session.go`

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID            int       `json:"id" db:"id"`
	UUID          uuid.UUID `json:"uuid" db:"uuid"`
	UserID        int       `json:"user_id" db:"user_id"`
	TokenHash     string    `json: "-" db:"token_hash"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	CreatedAtUnix int64     `json:"created_at_unix" db:"created_at_unix"`
}

// IsExpired checks if session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
```

---

## **Step 6: Updated Config**

### `internal/config/config.go`

```go
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string
	Environment string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret      string
	JWTExpireHours int

	// File Storage
	StorageType      string
	MinioEndpoint    string
	MinioAccessKey   string
	MinioSecretKey   string
	MinioBucket      string
	MinioUseSSL      bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load . env file if exists
	_ = godotenv.Load()

	jwtExpire, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	minioUseSSL, _ := strconv.ParseBool(getEnv("MINIO_USE_SSL", "false"))

	config := &Config{
		// Server
		ServerPort:   getEnv("PORT", "8080"),
		Environment: getEnv("ENV", "development"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "receipt_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpireHours: jwtExpire,

		// Storage
		StorageType:    getEnv("STORAGE_TYPE", "minio"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioBucket:    getEnv("MINIO_BUCKET", "receipts"),
		MinioUseSSL:    minioUseSSL,
	}

	return config, nil
}

// GetDatabaseURL returns formatted database connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c. DBSSLMode,
	)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
```

---

## **Step 7: Update . env files**

### `.env.example`

```env
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=receipt_db
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRE_HOURS=24

# File Storage Configuration
STORAGE_TYPE=minio
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=receipts
MINIO_USE_SSL=false
```

---

## **Step 8: Create API Server Entry Point**

### `cmd/api/main.go`

```go
package main

import (
	"fmt"
	"log"

	"github.com/dzulfiardev/receipt-backend/internal/config"
	"github.com/dzulfiardev/receipt-backend/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log. Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dbCfg := database.Config{
		Host:     cfg. DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg. DBName,
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
	e.Use(middleware.Logger())
	e.Use(middleware. Recover())
	e.Use(middleware. CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c. JSON(200, map[string]string{
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
```

---

## **Step 9: Updated Makefile**

### `Makefile`

```makefile
. PHONY: help run build migrate-up migrate-down migrate-drop migrate-version migrate-force migrate-create deps clean

help:  ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+: .*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the API server
	@echo "🚀 Starting API server..."
	@go run cmd/api/main. go

build: ## Build the API server
	@echo "🔨 Building API server..."
	@go build -o bin/api cmd/api/main.go

migrate-up: ## Run all pending migrations
	@echo "🚀 Running migrations..."
	@go run cmd/migrate/main. go -command=up

migrate-down: ## Rollback last migration
	@echo "⏪ Rolling back migration..."
	@go run cmd/migrate/main.go -command=down

migrate-drop: ## Drop all tables (DANGER!)
	@echo "⚠️ Dropping all tables..."
	@go run cmd/migrate/main. go -command=drop

migrate-version: ## Show current migration version
	@go run cmd/migrate/main.go -command=version

migrate-force: ## Force migration to specific version (usage: make migrate-force VERSION=1)
	@go run cmd/migrate/main.go -command=force -version=$(VERSION)

migrate-create: ## Create new migration (usage: make migrate-create NAME=add_column)
	@go run cmd/migrate/main.go -command=create -name=$(NAME)

deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/

docker-up: ## Start Docker containers (Postgres + MinIO)
	@echo "🐳 Starting Docker containers..."
	@docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "🐳 Stopping Docker containers..."
	@docker-compose down

docker-logs: ## Show Docker logs
	@docker-compose logs -f

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

. DEFAULT_GOAL := help
```

---

## **Step 10: Test the Setup**

```bash
# 1. Install dependencies
make deps

# 2. Start PostgreSQL
docker run --name receipt-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=receipt_db \
  -p 5432:5432 \
  -d postgres: 15-alpine

# 3. Run migrations
make migrate-up

# Output:
# 🚀 Running migrations...
# ✅ Migrations applied successfully!

# 4. Check migration version
make migrate-version

# Output:
# Current version: 4

# 5. Run API server
make run

# Output:
# ✅ Database connected successfully
# 🚀 Server starting on : 8080

# 6. Test health endpoint
curl http://localhost:8080/health

# Output:
# {"status":"ok","env":"development"}

# 7. Test ping endpoint
curl http://localhost:8080/api/v1/ping

# Output:
# {"message":"pong"}
```

---

## **Summary of Changes**

### ✅ **Schema Updates:**

- Changed all `INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY` to `SERIAL PRIMARY KEY`
- All tables now use `SERIAL` for auto-incrementing IDs

### ✅ **Framework Changes:**

- Switched from Gin to **Echo v4**
- Removed GORM, using **database/sql** with `lib/pq`
- Direct SQL queries (will be implemented in repositories)

### ✅ **New Files:**

- `internal/database/postgres.go` - Database connection pool
- `internal/domain/*. go` - Plain structs (no GORM tags)
- `cmd/api/main.go` - Echo server entry point

### ✅ **Updated Files:**

- All migration files now use `SERIAL`
- `internal/config/config.go` - Added server and storage config
- `Makefile` - Added `run` and Docker commands

---

**Next Steps:**

1. ✅ Run `make migrate-up` to apply schema changes
2. ✅ Run `make run` to start Echo server
3. 🚀 Ready to implement repositories and handlers!

Would you like me to continue with **implementing the repositories and handlers** using raw SQL queries? 🎯

## @dzulfiardev

yes make implementing the repositories and handlers using raw SQL queries

## Copilot

Perfect! Let's implement the repositories and handlers using raw SQL queries with Echo framework.

---

## **Step 1: Create Utility Helpers**

### `internal/utils/response.go`

```go
package utils

import (
	"github.com/labstack/echo/v4"
)

// Response represents standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse returns success response
func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse returns error response
func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse returns paginated response
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginatedSuccessResponse returns paginated success response
func PaginatedSuccessResponse(c echo.Context, statusCode int, data interface{}, meta PaginationMeta) error {
	return c.JSON(statusCode, PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: meta,
	})
}
```

### `internal/utils/jwt.go`

```go
package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(userID int, email string, secret string, expireHours int) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt. NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates JWT token and returns claims
func ValidateJWT(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token. Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
```

### `internal/utils/validator.go`

```go
package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground validator
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		// Format validation errors
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt. Sprintf("%s: %s", err.Field(), getErrorMsg(err)))
		}
		return fmt.Errorf(strings.Join(errors, ", "))
	}
	return nil
}

// getErrorMsg returns user-friendly error message
func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	default:
		return "Invalid value"
	}
}
```

---

## **Step 2: Create Repositories**

### `internal/repository/user_repository.go`

```go
package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dzulfiardev/receipt-backend/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id int) (*domain.User, error)
	FindByUUID(uuid string) (*domain.User, error)
	Update(user *domain.User) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db:  db}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, created_at_unix, updated_at_unix)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, uuid, created_at, updated_at
	`

	now := time.Now().Unix()
	err := r.db.QueryRow(
		query,
		user. Email,
		user.PasswordHash,
		user.FullName,
		now,
		now,
	).Scan(&user.ID, &user.UUID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAtUnix = now
	user.UpdatedAtUnix = now

	return nil
}

// FindByEmail finds user by email
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, uuid, email, password_hash, full_name, created_at, updated_at, created_at_unix, updated_at_unix
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedAtUnix,
		&user.UpdatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByID finds user by ID
func (r *userRepository) FindByID(id int) (*domain.User, error) {
	query := `
		SELECT id, uuid, email, password_hash, full_name, created_at, updated_at, created_at_unix, updated_at_unix
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user. PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedAtUnix,
		&user.UpdatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByUUID finds user by UUID
func (r *userRepository) FindByUUID(uuidStr string) (*domain.User, error) {
	query := `
		SELECT id, uuid, email, password_hash, full_name, created_at, updated_at, created_at_unix, updated_at_unix
		FROM users
		WHERE uuid = $1
	`

	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid:  %w", err)
	}

	user := &domain.User{}
	err = r.db.QueryRow(query, uid).Scan(
		&user. ID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user. UpdatedAt,
		&user.CreatedAtUnix,
		&user.UpdatedAtUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt. Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// Update updates user information
func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, updated_at = NOW(), updated_at_unix = $3
		WHERE id = $4
		RETURNING updated_at
	`

	now := time.Now().Unix()
	err := r.db.QueryRow(query, user.Email, user. FullName, now, user.ID).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt. Errorf("failed to update user: %w", err)
	}

	user.UpdatedAtUnix = now
	return nil
}
```

---

### `internal/repository/receipt_repository.go`

```go
package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dzulfiardev/receipt-backend/internal/domain"
	"github. com/google/uuid"
)

type ReceiptRepository interface {
	Create(receipt *domain. Receipt) error
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
		receipt. Date,
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
		&receipt. ImageURL,
		&receipt. OriginalFilename,
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
		return nil, fmt. Errorf("invalid uuid: %w", err)
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
		&receipt. TotalDiscount,
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
		return nil, 0, fmt. Errorf("failed to query receipts: %w", err)
	}
	defer rows.Close()

	var receipts []domain. Receipt
	for rows.Next() {
		var receipt domain.Receipt
		err := rows.Scan(
			&receipt.ID,
			&receipt.UUID,
			&receipt. UserID,
			&receipt. StoreName,
			&receipt. Address,
			&receipt.Phone,
			&receipt.Date,
			&receipt.ImageURL,
			&receipt.OriginalFilename,
			&receipt. FileSize,
			&receipt. UploadDate,
			&receipt.Status,
			&receipt.TotalItems,
			&receipt.TotalSpending,
			&receipt.TotalDiscount,
			&receipt.CreatedAt,
			&receipt.UpdatedAt,
			&receipt. CreatedAtUnix,
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
		receipt. Status,
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
		return fmt. Errorf("receipt not found")
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
		"total_spending":  totalSpending,
		"total_discount":   totalDiscount,
		"average_spending": avgSpending,
		"net_spending":    totalSpending - totalDiscount,
	}

	return stats, nil
}
```

---

### `internal/repository/item_repository.go`

```go
package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dzulfiardev/receipt-backend/internal/domain"
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
		item. Price,
		item.Total,
		now,
	).Scan(&item.ID, &item.UUID, &item.CreatedAt)

	if err != nil {
		return fmt. Errorf("failed to create item:  %w", err)
	}

	item.CreatedAtUnix = now
	return nil
}

// CreateBatch creates multiple items in a single transaction
func (r *itemRepository) CreateBatch(items []domain.Item) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt. Errorf("failed to begin transaction: %w", err)
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
			items[i]. Quantity,
			items[i].Price,
			items[i].Total,
			now,
		).Scan(&items[i].ID, &items[i]. UUID, &items[i].CreatedAt)

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
		SELECT id, uuid, receipt_id, name, unit_price, quantity, price, total,
		       created_at, created_at_unix
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
			&item. UUID,
			&item. ReceiptID,
			&item.Name,
			&item. UnitPrice,
			&item.Quantity,
			&item.Price,
			&item. Total,
			&item.CreatedAt,
			&item. CreatedAtUnix,
		)
		if err != nil {
			return nil, fmt. Errorf("failed to scan item:  %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// FindByID finds item by ID
func (r *itemRepository) FindByID(id int) (*domain.Item, error) {
	query := `
		SELECT id, uuid, receipt_id, name, unit_price, quantity, price, total,
		       created_at, created_at_unix
		FROM items
		WHERE id = $1
	`

	item := &domain.Item{}
	err := r.db.QueryRow(query, id).Scan(
		&item. ID,
		&item.UUID,
		&item.ReceiptID,
		&item. Name,
		&item.UnitPrice,
		&item. Quantity,
		&item. Price,
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

	result, err := r.db. Exec(
		query,
		item.Name,
		item.UnitPrice,
		item.Quantity,
		item.Price,
		item.Total,
		item. ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rowsAffected, err := result. RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt. Errorf("item not found")
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
```

---

## **Step 3: Create Services**

### `internal/service/auth_service.go`

```go
package service

import (
	"fmt"

	"github.com/dzulfiardev/receipt-backend/internal/domain"
	"github.com/dzulfiardev/receipt-backend/internal/repository"
	"github.com/dzulfiardev/receipt-backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req domain.CreateUserRequest) (*domain.User, error)
	Login(req domain.LoginRequest) (string, *domain.User, error)
	GetUserByID(id int) (*domain.User, error)
}

type authService struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtExpireHours int
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpireHours int) AuthService {
	return &authService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

// Register registers a new user
func (s *authService) Register(req domain.CreateUserRequest) (*domain.User, error) {
	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain. User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates user and returns JWT token
func (s *authService) Login(req domain.LoginRequest) (string, *domain.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user. PasswordHash), []byte(req.Password)); err != nil {
		return "", nil, fmt. Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user. ID, user.Email, s.jwtSecret, s.jwtExpireHours)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// GetUserByID gets user by ID
func (s *authService) GetUserByID(id int) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}
```

---

### `internal/service/receipt_service.go`

```go
package service

import (
	"database/sql"
	"fmt"

	"github.com/dzulfiardev/receipt-backend/internal/domain"
	"github. com/dzulfiardev/receipt-backend/internal/repository"
)

type ReceiptService interface {
	CreateReceipt(userID int, req domain.CreateReceiptRequest, imageURL string, filename string, fileSize int) (*domain. ReceiptWithItems, error)
	GetReceiptByID(id int, userID int) (*domain.ReceiptWithItems, error)
	GetReceiptsByUserID(userID int, page, limit int) ([]domain.Receipt, int64, error)
	UpdateReceipt(id int, userID int, req domain.CreateReceiptRequest) (*domain.ReceiptWithItems, error)
	DeleteReceipt(id int, userID int) error
	GetStatsByUserID(userID int) (map[string]interface{}, error)
}

type receiptService struct {
	receiptRepo repository. ReceiptRepository
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
		Phone:            sql.NullInt64{Int64: *req.Phone, Valid: req. Phone != nil},
		Date:             date,
		ImageURL:         imageURL,
		OriginalFilename: filename,
		FileSize:         fileSize,
		Status:           domain.StatusCompleted,
		TotalItems:       req.TotalItems,
		TotalSpending:    req.TotalSpending,
		TotalDiscount:    req.TotalDiscount,
	}

	if err := s.receiptRepo. Create(receipt); err != nil {
		return nil, fmt. Errorf("failed to create receipt:  %w", err)
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
			Total:     itemReq. Total,
		}
		items = append(items, item)
	}

	if len(items) > 0 {
		if err := s.itemRepo.CreateBatch(items); err != nil {
			return nil, fmt. Errorf("failed to create items:  %w", err)
		}
	}

	return &domain.ReceiptWithItems{
		Receipt:  *receipt,
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

	return &domain. ReceiptWithItems{
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
	receipt.Address = sql. NullString{String: req. Address, Valid: req.Address != ""}
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
```

---

## **Step 4: Create Middleware**

### `internal/middleware/auth. go`

```go
package middleware

import (
	"strings"

	"github.com/dzulfiardev/receipt-backend/internal/utils"
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
```

---

## **Step 5: Create Handlers**

### `internal/handler/auth_handler.go`

```go
package handler

import (
	"github.com/dzulfiardev/receipt-backend/internal/domain"
	"github.com/dzulfiardev/receipt-backend/internal/middleware"
	"github.com/dzulfiardev/receipt-backend/internal/service"
	"github.com/dzulfiardev/receipt-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *utils.Validator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
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

	user, err := h.authService. Register(req)
	if err != nil {
		return utils. ErrorResponse(c, 400, err.Error())
	}

	return utils. SuccessResponse(c, 201, "User registered successfully", user. ToResponse())
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest

	if err := c.Bind(&req); err != nil {
		return utils. ErrorResponse(c, 400, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	token, user, err := h. authService.Login(req)
	if err != nil {
		return utils.ErrorResponse(c, 401, err.Error())
	}

	return utils.SuccessResponse(c, 200, "Login successful", map[string]interface{}{
		"token": token,
		"user":   user. ToResponse(),
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
```

---

### `internal/handler/receipt_handler.go`

```go
package handler

import (
	"math"
	"strconv"

	"github.com/dzulfiardev/receipt-backend/internal/domain"
	"github.com/dzulfiardev/receipt-backend/internal/middleware"
	"github.com/dzulfiardev/receipt-backend/internal/service"
	"github.com/dzulfiardev/receipt-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

type ReceiptHandler struct {
	receiptService service.ReceiptService
	validator      *utils. Validator
}

// NewReceiptHandler creates a new receipt handler
func NewReceiptHandler(receiptService service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{
		receiptService: receiptService,
		validator:      utils. NewValidator(),
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
	userID := middleware. GetUserID(c)

	id, err := strconv. Atoi(c.Param("id"))
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
	limit, _ := strconv. Atoi(c.QueryParam("limit"))

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

	if err := h.receiptService. DeleteReceipt(id, userID); err != nil {
		return utils.ErrorResponse(c, 404, err.Error())
	}

	return utils.SuccessResponse(c, 200, "Receipt deleted successfully", nil)
}

// GetStats gets spending statistics
func (h *ReceiptHandler) GetStats(c echo.Context) error {
	userID := middleware.GetUserID(c)

	stats, err := h.receiptService. GetStatsByUserID(userID)
	if err != nil {
		return utils.ErrorResponse(c, 500, err.Error())
	}

	return utils.SuccessResponse(c, 200, "Statistics retrieved successfully", stats)
}
```

---

## **Step 6: Update Main Server**

### `cmd/api/main.go`

```go
package main

import (
	"fmt"
	"log"

	"github.com/dzulfiardev/receipt-backend/internal/config"
	"github. com/dzulfiardev/receipt-backend/internal/database"
	"github.com/dzulfiardev/receipt-backend/internal/handler"
	"github.com/dzulfiardev/receipt-backend/internal/middleware"
	"github.com/dzulfiardev/receipt-backend/internal/repository"
	"github.com/dzulfiardev/receipt-backend/internal/service"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log. Fatalf("Failed to load config:  %v", err)
	}

	// Connect to database
	dbCfg := database.Config{
		Host:     cfg. DBHost,
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
	e := echo. New()

	// Global middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware. CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c. JSON(200, map[string]string{
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
			receipts.DELETE("/:id", receiptHandler. DeleteReceipt)
			receipts.GET("/stats", receiptHandler.GetStats)
		}
	}

	// Start server
	address := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Server starting on %s", address)
	log.Printf("📝 Environment: %s", cfg.Environment)
	log.Printf("🔐 JWT Expiration: %d hours", cfg. JWTExpireHours)

	if err := e.Start(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

---

## **Step 7: Install Additional Dependencies**

```bash
go get -u github.com/go-playground/validator/v10
go mod tidy
```

---

## **Step 8: Test the API**

### 1. Run Migrations

```bash
make migrate-up
```

### 2. Start Server

```bash
make run
```

### 3. Test Endpoints

#### **Register User**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "email": "test@example.com",
    "full_name": "Test User",
    "created_at": "2026-01-16T10:00:00Z",
    "created_at_unix": 1736937600
  }
}
```

#### **Login**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":  "test@example.com",
    "password": "password123"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.. .",
    "user": {
      "id": 1,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "email": "test@example.com",
      "full_name": "Test User",
      "created_at": "2026-01-16T10:00:00Z",
      "created_at_unix": 1736937600
    }
  }
}
```

#### **Get Current User**

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### **Create Receipt**

```bash
curl -X POST http://localhost:8080/api/v1/receipts \
  -H "Authorization:  Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "store_name": "Indomaret",
    "address": "Jl. Sudirman No. 123",
    "phone": 628123456789,
    "total_items": 3,
    "total_spending": 75000,
    "total_discount": 5000,
    "items": [
      {
        "name": "Indomie Goreng",
        "unit_price": 3500,
        "quantity": 10,
        "price":  3500,
        "total":  35000
      },
      {
        "name": "Aqua 600ml",
        "unit_price":  4000,
        "quantity": 5,
        "price": 4000,
        "total":  20000
      },
      {
        "name": "Teh Botol",
        "unit_price": 5000,
        "quantity": 4,
        "price": 5000,
        "total": 20000
      }
    ]
  }'
```

#### **Get All Receipts (with pagination)**

```bash
curl -X GET "http://localhost:8080/api/v1/receipts?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### **Get Receipt by ID**

```bash
curl -X GET http://localhost:8080/api/v1/receipts/1 \
  -H "Authorization:  Bearer YOUR_TOKEN_HERE"
```

#### **Get Statistics**

```bash
curl -X GET http://localhost:8080/api/v1/receipts/stats \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Response:**

```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "average_spending": 75000,
    "net_spending": 70000,
    "total_discount": 5000,
    "total_receipts": 1,
    "total_spending": 75000
  }
}
```

#### **Delete Receipt**

```bash
curl -X DELETE http://localhost:8080/api/v1/receipts/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## **Summary**

✅ **Repositories** - Raw SQL queries with `database/sql`  
✅ **Services** - Business logic layer  
✅ **Handlers** - HTTP request handlers with Echo  
✅ **Middleware** - JWT authentication  
✅ **Utils** - Response helpers, JWT, validation  
✅ **Complete API** - Register, Login, CRUD Receipts, Statistics

---

**Complete File Structure:**

```
receipt-backend/
├── cmd/
│   ├── api/main.go ✅
│   └── migrate/main.go ✅
├── internal/
│   ├── config/config.go ✅
│   ├── database/postgres.go ✅
│   ├── domain/
│   │   ├── user.go ✅
│   │   ├── receipt.go ✅
│   │   ├── item.go ✅
│   │   └── session.go ✅
│   ├── repository/
│   │   ├── user_repository.go ✅
│   │   ├── receipt_repository.go ✅
│   │   └── item_repository.go ✅
│   ├── service/
│   │   ├── auth_service.go ✅
│   │   └── receipt_service.go ✅
│   ├── handler/
│   │   ├── auth_handler. go ✅
│   │   └── receipt_handler.go ✅
│   ├── middleware/
│   │   └── auth.go ✅
│   └── utils/
│       ├── response.go ✅
│       ├── jwt.go ✅
│       └── validator.go ✅
├── migrations/ ✅
├── . env ✅
├── go.mod ✅
├── Makefile ✅
└── README.md ✅
```

---

**Next Steps:**

1. ✅ Test all endpoints
2. 🚀 Add file upload (S3/MinIO) for receipt images
3. 🚀 Add Gemini API integration for AI extraction
4. 🚀 Add Docker Compose for complete stack

Would you like me to implement **file upload with MinIO** or **Gemini API integration** next? 🎯
