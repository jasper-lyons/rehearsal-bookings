package main

import (
	"log"
	"os"
	"fmt"
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

	driver.Query("PRAGMA vdbe_debug=ON;")
	driver.Query("PRAGAM journal_mode=WAL;")
	driver.Query("PRAGMA wal_autocheckpoint=1000;")
	driver.Query("PRAGMA busy_timeout=5000;") // 5000 milliseconds
	driver.Query("PRAGMA wal_checkpoint(PASSIVE);")
	driver.Query("PRAGMA synchronous=NORMAL;")

	br := da.NewBookingsRepository(driver)

	imminentBookings, err := br.Where(`start_time <= datetime(current_timestamp, '+1 hours')`)
	if err != nil {
		log.Fatal(err)
	}

	for _, booking := range imminentBookings {
		fmt.Println("Sending Access Codes to booking %s", booking.Id)
	}
}
