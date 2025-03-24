package main

import (
	"log"
	"os"
	"fmt"
	"github.com/joho/godotenv"
	"encoding/json"
	"time"

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

	imminentBookings := br.Where(`start_time <= datetime(current_timestamp, '+1 hours')`)
	for _, booking := range imminentBookings {
		fmt.Println("Sending Access Codes to booking %s", bookings.id)
	}

	br.Update(heldBookings)
}
