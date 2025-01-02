package main

import (
	"flag"
	"fmt"
	"net/url"
	"log"
	"os"
	"strings"
	"time"
	"path/filepath"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrationsPath = "./migrations" // Project-specific path to migrations
)

func normaliseUrl(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse DATABASE_URL=%w", err)
	}

	if u.Scheme == "postgres" {
		u.Scheme = "potgresql"
	}

	return u.String(), nil
}

func main() {
	var command string
	flag.StringVar(&command, "command", "", "Migration command (up/down/version/create)")

	var url string
	flag.StringVar(&url, "url", "", "Url of the database, will override DATABASE_URL")
	flag.Parse()

	if command == "" {
		flag.Usage()
		os.Exit(1)
	}

	if url == "" {
		url = os.Getenv("DATABASE_URL")
	}

	if url == "" {
		log.Fatal("DATABASE_URL environment variable or url param are required")
	}

	normalisedUrl, err := normaliseUrl(url)
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New(
		"file://./migrations",
		normalisedUrl,
	)
	if err != nil {
		log.Fatal("Failed to create migrator: %w", err)
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations: ", err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to rollback migrations: ", err)
		}
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get version: ", err)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
	case "create":
		if flag.NArg() < 1 {
			log.Fatal("Migration name argument required")
		}
		name := strings.Join(flag.Args(), "_")
		timestamp := time.Now().Unix()
		baseName := fmt.Sprintf("%d_%s", timestamp, name)

		upFile := filepath.Join(migrationsPath, baseName + ".up.sql")
		downFile := filepath.Join(migrationsPath, baseName + ".down.sql")

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

	// Don't forget to close the migrator
	defer func() {
		if _, err := m.Close(); err != nil {
			log.Fatal("Failed to close migrator: ", err)
		}
	}()
}
