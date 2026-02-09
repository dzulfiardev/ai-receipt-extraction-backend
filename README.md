# AI Receipt Extraction Backend

Backend API for AI-powered receipt extraction and management system built with Go.

## Prerequisites

- Go 1.25.1 or higher
- PostgreSQL database
- MinIO (for file storage)

## Getting Started

### 1. Install Dependencies

Install all required Go dependencies:

```bash
make deps
```

Or manually:

```bash
go mod download
go mod tidy
```

### 2. Setup Environment Variables

Create a `.env` file in the root directory of the project with the following configuration:

```env
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=root
DB_NAME=receipt_db
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRE_HOURS=24

# File Storage Configuration (MinIO)
STORAGE_TYPE=minio
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=receipts
MINIO_USE_SSL=false
```

**Important:** Make sure to change the `JWT_SECRET` to a secure random string in production!

### 3. Database Migration

Run all pending migrations to set up the database schema:

```bash
make migrate-up
```

This will create the necessary tables:

- `users` - User accounts
- `receipts` - Receipt records
- `items` - Receipt items
- `sessions` - User sessions

## Available Commands

### Running the Server

Start the API server:

```bash
make run
```

The server will start on `http://localhost:8080` (or the port specified in your `.env` file).

### Building

Build the production binary:

```bash
make build
```

The binary will be created in the `bin/` directory.

### Migration Commands

- **Run migrations:** `make migrate-up`
- **Rollback last migration:** `make migrate-down`
- **Drop all tables:** `make migrate-drop` ⚠️ **Warning: This will delete all data!**
- **Check migration version:** `make migrate-version`
- **Force migration version:** `make migrate-force VERSION=<version_number>`
- **Create new migration:** `make migrate-create NAME=<migration_name>`

### Other Commands

- **Install dependencies:** `make deps`
- **Clean build artifacts:** `make clean`
- **Show all commands:** `make help`

## Project Structure

```
.
├── cmd/
│   ├── api/          # API server entry point
│   └── migrate/      # Database migration tool
├── internal/
│   ├── config/       # Configuration management
│   ├── database/     # Database connection
│   ├── domain/       # Domain models
│   └── utils/        # Utility functions
├── migrations/       # SQL migration files
├── .env             # Environment variables (create this)
├── go.mod           # Go module file
├── Makefile         # Build and run commands
└── README.md        # This file
```

## License

MIT
