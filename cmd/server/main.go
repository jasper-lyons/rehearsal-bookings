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
	"github.com/joho/godotenv"

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

	driver := da.NewSqliteDriver("db/development.db")
	br := da.NewBookingsRepository(driver)

	http.Handle("GET /", handlers.BookingsNew(br))
	http.Handle("GET /bookings", handlers.BookingsIndex(br))
	http.Handle("POST /bookings", handlers.BookingsCreate(br))
	http.Handle("POST /sumup/checkouts", handlers.SumupCheckoutCreate())

	server := &http.Server {
		Addr: EnvOrDefault("PORT", ":8080"),
		Handler: handlers.Logging(http.DefaultServeMux),
	}

	log.Printf("Listening on port %s\n", EnvOrDefault("PORT", ":8080"))
	log.Fatal(server.ListenAndServe())
}
