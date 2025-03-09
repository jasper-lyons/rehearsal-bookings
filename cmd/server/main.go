package main

import (
	_ "fmt"
	"log"
	"net/http"
	"os"

	// GO Embed doesn't support embedding files from outside the module boundary
	// (in this case, anything outside of this directory) but we want to store
	// templates and static files at the route directory so we need to treat them
	// as their own go modules (you'll see the main.go files in the web/static and
	// web/templates directories) so that we can import them here and access the
	// embedded static files!
	// "rehearsal-bookings/web/static"
	da "rehearsal-bookings/pkg/data_access"
	handlers "rehearsal-bookings/pkg/handlers"
	static "rehearsal-bookings/web/static"

	"fmt"

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

	var driver *da.SqliteDriver
	if os.Getenv("APP_ENV") == "production" {
		driver = da.NewSqliteDriver("db/production.db")
	} else {
		driver = da.NewSqliteDriver("db/development.db")
	}

	br := da.NewBookingsRepository(driver)
	sumupApi := handlers.NewApi("https://api.sumup.com", map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("SUMUP_API_KEY")),
	})

	http.Handle("GET /admin/bookings", handlers.AdminBookingsIndex(br))
	http.Handle("GET /admin/bookings/new", handlers.AdminBookingsNew(br))

	http.Handle("POST /bookings/{id}/confirm", handlers.BookingsConfirm(br, sumupApi))
	http.Handle("POST /bookings", handlers.BookingsCreate(br))

	http.Handle("GET /rooms", handlers.RoomsIndex(br))

	http.Handle("POST /sumup/checkouts", handlers.SumupCheckoutCreate(sumupApi))

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.StaticFiles))))
	http.Handle("GET /", handlers.BookingsNew(br))

	server := &http.Server{
		Addr:    ":" + EnvOrDefault("PORT", "8080"),
		Handler: handlers.Logging(http.DefaultServeMux),
	}

	log.Printf("Listening on port %s\n", EnvOrDefault("PORT", "8080"))
	log.Fatal(server.ListenAndServe())
}
