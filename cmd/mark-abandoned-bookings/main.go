package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"

	da "rehearsal-bookings/pkg/data_access"
)

func EnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var driver *da.SqliteDriver
	if os.Getenv("APP_ENV") == "production" {
		driver = da.NewSqliteDriver("db/production.db")
	} else {
		driver = da.NewSqliteDriver("db/development.db")
	}

	br := da.NewBookingsRepository(driver)

	heldBookings, err := br.Where(`status = 'held' and expiration < current_timestamp`)
	if err != nil {
		log.Fatal(err)
	}
	for _, booking := range heldBookings {
		booking.Status = "abandoned"
	}

	br.Update(heldBookings)
}
