package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var command string
	var steps string

	flag.StringVar(&command, "command", "", "Migration command (up, down, reset, drop, version, force)")
	flag.StringVar(&steps, "steps", "", "Number of steps for up/down migration")
	flag.Parse()

	if command == "" {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create migration instance
	m, err := migrate.New(
		"file://migrations",
		cfg.GetDatabaseURL(),
	)

	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	defer m.Close()

	// Execute migration command
	switch command {
	case "up":
		if steps != "" {
			stepsInt, err := strconv.Atoi(steps)
			if err != nil {
				log.Fatalf("Invalid steps number: %v", err)
			}
			err = m.Steps(stepsInt)
		} else {
			err = m.Up()
		}

		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}

		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
		} else {
			fmt.Println("Migration up completed successfully")
		}

	case "down":
		if steps != "" {
			stepsInt, err := strconv.Atoi(steps)
			if err != nil {
				log.Fatalf("Invalid steps number: %v", err)
			}
			err = m.Steps(-stepsInt)
		} else {
			err = m.Steps(-1)
		}

		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
		} else {
			fmt.Println("Migration down completed successfully")
		}

	case "reset":
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration reset failed: %v", err)
		}
		fmt.Println("Migration reset completed successfully")

	case "drop":
		err = m.Drop()
		if err != nil {
			log.Fatalf("Migration drop failed: %v", err)
		}
		fmt.Println("Database dropped successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d\n", version)
		if dirty {
			fmt.Println("Database is in dirty state")
		}

	case "force":
		if steps == "" {
			log.Fatal("Version number is required for force command")
		}
		version, err := strconv.Atoi(steps)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		err = m.Force(version)
		if err != nil {
			log.Fatalf("Migration force failed: %v", err)
		}
		fmt.Printf("Forced migration to version %d\n", version)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Migration Runner")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go -command=<command> [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  up          Apply all pending migrations")
	fmt.Println("  up -steps=N Apply N migrations")
	fmt.Println("  down        Rollback one migration")
	fmt.Println("  down -steps=N Rollback N migrations")
	fmt.Println("  reset       Rollback all migrations")
	fmt.Println("  drop        Drop the entire database")
	fmt.Println("  version     Show current migration version")
	fmt.Println("  force -steps=N Force database to specific version")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go -command=up")
	fmt.Println("  go run cmd/migrate/main.go -command=up -steps=1")
	fmt.Println("  go run cmd/migrate/main.go -command=down")
	fmt.Println("  go run cmd/migrate/main.go -command=down -steps=2")
	fmt.Println("  go run cmd/migrate/main.go -command=reset")
	fmt.Println("  go run cmd/migrate/main.go -command=version")
	fmt.Println("  go run cmd/migrate/main.go -command=force -steps=3")
}
