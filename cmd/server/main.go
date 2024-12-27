package main

import (
	"net/http"
	"log"
	"os"
	_ "fmt"
	// GO Embed doesn't support embedding files from outside the module boundary
	// (in this case, anything outside of this directory) but we want to store
	// templates and static files at the route directory so we need to treat them
	// as their own go modules (you'll see the main.go files in the web/static and
	// web/templates directories) so that we can import them here and access the 
	// embedded static files!
	// "rehearsal-bookings/web/static"
	da "rehearsal-bookings/pkg/data_access"
	handlers "rehearsal-bookings/pkg/handlers"
)

func EnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}


func main() {
	driver := da.NewSqliteDriver("db/development.db")
	br := da.NewBookingsRepository(driver)

	http.Handle("GET /new", handlers.BookingsNew(br))
	http.Handle("GET /", handlers.BookingsIndex(br))
	http.Handle("POST /", handlers.BookingsCreate(br))

	log.Fatal(http.ListenAndServe(EnvOrDefault("PORT", ":8080"), nil))
}
