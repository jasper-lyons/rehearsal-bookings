package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrationsPath = "./migrations"
)

func getDefaultDatabaseURL() string {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		return "sqlite3://db/production.db"
	}
	return "sqlite3://db/development.db"
}

func normaliseUrl(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}

	if u.Scheme == "postgres" {
		u.Scheme = "postgresql"
	}

	return u.String(), nil
}

func main() {
	// Create a new FlagSet to handle command-line arguments
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)

	// Define url flag
	var dbURL string
	fs.StringVar(&dbURL, "url", "", "Database URL (overrides environment defaults)")

	// Check if we have any arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <command> [options]")
		fmt.Println("Commands: up, down, version, create")
		os.Exit(1)
	}

	// First argument is the command
	command := os.Args[1]

	// Parse remaining arguments
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	// Determine database URL
	if dbURL == "" {
		// Check DATABASE_URL environment variable first
		dbURL = os.Getenv("DATABASE_URL")
		if dbURL == "" {
			// If no DATABASE_URL, use default based on APP_ENV
			dbURL = getDefaultDatabaseURL()
		}
	}

	// Ensure database directory exists
	dbDir := filepath.Dir(strings.TrimPrefix(dbURL, "sqlite3://"))
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	normalisedUrl, err := normaliseUrl(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New(
		"file://./migrations",
		normalisedUrl,
	)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	// Don't forget to close the migrator
	defer func() {
		if _, err := m.Close(); err != nil {
			log.Fatalf("Failed to close migrator: %v", err)
		}
	}()

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations: ", err)
		}
		fmt.Println("Successfully applied migrations")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to rollback migrations: ", err)
		}
		fmt.Println("Successfully rolled back migrations")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get version: ", err)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)

	case "create":
		if fs.NArg() < 1 {
			log.Fatal("Migration name argument required")
		}
		name := strings.Join(fs.Args(), "_")
		timestamp := time.Now().Unix()
		baseName := fmt.Sprintf("%d_%s", timestamp, name)

		upFile := filepath.Join(migrationsPath, baseName+".up.sql")
		downFile := filepath.Join(migrationsPath, baseName+".down.sql")

		// Ensure migrations directory exists
		if err := os.MkdirAll(migrationsPath, 0755); err != nil {
			log.Fatalf("Failed to create migrations directory: %v", err)
		}

		// Create empty up file
		if err := os.WriteFile(upFile, []byte(""), 0644); err != nil {
			log.Fatalf("Failed to create up migration: %v", err)
		}

		// Create empty down file
		if err := os.WriteFile(downFile, []byte(""), 0644); err != nil {
			log.Fatalf("Failed to create down migration: %v", err)
		}

		fmt.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
