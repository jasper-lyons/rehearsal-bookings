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

	heldBookings, err := br.Where(`expiration < current_timestamp and status = 'hold'`)
	if err != nil {
		log.Fatal(err)
	}

	for i := range heldBookings {
		log.Println("Marking booking %d as abandoned", heldBookings[i].Id)
		heldBookings[i].Status = "abandoned"
	}

	_, err = br.Update(heldBookings)
	if err != nil {
		log.Fatal(err)
	}
}
