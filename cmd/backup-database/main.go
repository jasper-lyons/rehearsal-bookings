package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

// EnvOrDefault returns environment variable value or fallback if not set
func EnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// ensureDir creates directory if it doesn't exist
func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination directory if it doesn't exist
	err = ensureDir(filepath.Dir(dst))
	if err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// getBackupFiles returns list of backup files in a directory sorted by date
func getBackupFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		// If directory doesn't exist yet, return empty slice
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var filePaths []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "production.") {
			filePaths = append(filePaths, filepath.Join(dir, file.Name()))
		}
	}

	// Sort files by modification time (newest first)
	sort.Slice(filePaths, func(i, j int) bool {
		iInfo, _ := os.Stat(filePaths[i])
		jInfo, _ := os.Stat(filePaths[j])
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	return filePaths, nil
}

// cleanupOldBackups enforces retention policy
func cleanupOldBackups(dir string, keep int) error {
	files, err := getBackupFiles(dir)
	if err != nil {
		return err
	}

	// Remove files beyond retention limit
	for i := keep; i < len(files); i++ {
		// Extra safety check: never delete a file called production.db
		if filepath.Base(files[i]) == "production.db" {
			log.Printf("Skipping deletion of primary database file: %s", files[i])
			continue
		}
		
		err := os.Remove(files[i])
		if err != nil {
			return err
		}
		log.Printf("Removed old backup: %s", files[i])
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	// Determine database path based on environment
	dbPath := "db/development.db"
	if os.Getenv("APP_ENV") == "production" {
		dbPath = "db/production.db"
	}

	// Ensure database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("Database file not found: %s", dbPath)
	}

	// Create backup directories
	dirs := []string{"db/days", "db/weeks", "db/months", "db/years"}
	for _, dir := range dirs {
		if err := ensureDir(dir); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Get current date information
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	
	// Create daily backup
	dailyBackupPath := fmt.Sprintf("db/days/production.%s.db", dateStr)
	log.Printf("Creating daily backup: %s", dailyBackupPath)
	if err := copyFile(dbPath, dailyBackupPath); err != nil {
		log.Fatalf("Failed to create daily backup: %v", err)
	}

	// Handle weekly backup (on Mondays, day 1 in Go's time package)
	if now.Weekday() == time.Monday {
		weeklyBackupPath := fmt.Sprintf("db/weeks/production.%s.db", dateStr)
		log.Printf("Creating weekly backup: %s", weeklyBackupPath)
		if err := copyFile(dbPath, weeklyBackupPath); err != nil {
			log.Fatalf("Failed to create weekly backup: %v", err)
		}
	}

	// Handle monthly backup (on the 1st of each month)
	if now.Day() == 1 {
		monthlyBackupPath := fmt.Sprintf("db/months/production.%s.db", dateStr)
		log.Printf("Creating monthly backup: %s", monthlyBackupPath)
		if err := copyFile(dbPath, monthlyBackupPath); err != nil {
			log.Fatalf("Failed to create monthly backup: %v", err)
		}
	}

	// Handle yearly backup (on January 1st)
	if now.Day() == 1 && now.Month() == time.January {
		yearlyBackupPath := fmt.Sprintf("db/years/production.%s.db", dateStr)
		log.Printf("Creating yearly backup: %s", yearlyBackupPath)
		if err := copyFile(dbPath, yearlyBackupPath); err != nil {
			log.Fatalf("Failed to create yearly backup: %v", err)
		}
	}

	// Apply retention policies
	retentionPolicies := map[string]int{
		"db/days":   7,    // Keep last 7 daily backups
		"db/weeks":  4,    // Keep last 4 weekly backups
		"db/months": 12,   // Keep last 12 monthly backups
		"db/years":  5,    // Keep last 5 yearly backups
	}

	for dir, keep := range retentionPolicies {
		if err := cleanupOldBackups(dir, keep); err != nil {
			log.Printf("Warning: Failed to clean up old backups in %s: %v", dir, err)
		}
	}

	log.Println("Database backup completed successfully")
}
