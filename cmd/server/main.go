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
	"github.com/stripe/stripe-go/v81"
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

	var paymentsApi handlers.Api
	if os.Getenv("FEATURE_FLAG_PAYMENTS_PROVIDER") == "sumup" {
		paymentsApi = handlers.NewApi("https://api.sumup.com", map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("SUMUP_API_KEY")),
		})
	} else if os.Getenv("FEATURE_FLAG_PAYMENTS_PROVIDER") == "stripe" {
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		paymentsApi = handlers.NewApi("todo: stripe api url", map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("STRIPE_API_KEY")),
		})
	}

	basicauth := handlers.CreateBasicAuthMiddleware(
		os.Getenv("ADMIN_USERNAME"),
		os.Getenv("ADMIN_PASSWORD"),
	)

	br := da.NewBookingsRepository(driver)

	// Adming methods
	http.Handle("POST /admin/bookings", handlers.AdminBookingsCreate(br))
	http.Handle("PUT /admin/bookings/{id}/update", handlers.AdminBookingsUpdate(br))
	http.Handle("DELETE /admin/bookings/{id}", basicauth(handlers.AdminBookingsDelete(br)))
	http.Handle("PUT /admin/bookings/{id}/paid", basicauth(handlers.AdminBookingsStatusUpdate(br)))
	http.Handle("PUT /admin/bookings/{id}/cancel", basicauth(handlers.AdminBookingsStatusUpdate(br)))

	// Admin panel forms
	http.Handle("GET /admin/bookings/new", basicauth(handlers.AdminBookingsNew(br)))
	http.Handle("GET /admin/bookings/{id}/edit", basicauth(handlers.AdminBookingsUpdateView(br)))

	// Admin booking views
	http.Handle("GET /admin/bookings", basicauth(handlers.AdminViewDailyBookings(br)))
	http.Handle("GET /admin/availability", basicauth(handlers.AdminViewDailyAvailability(br)))
	http.Handle("GET /admin/bookings/all", basicauth(handlers.AdminViewAllBookings(br)))
	http.Handle("GET /admin/bookings/future", basicauth(handlers.AdminBookingsFutureBookings(br)))
	http.Handle("GET /admin/bookings/past", basicauth(handlers.AdminBookingsPastBookings(br)))
	http.Handle("GET /admin", handlers.Redirect("/admin/bookings"))

	http.Handle("POST /bookings/{id}/confirm", handlers.BookingsConfirm(br, paymentsApi))
	http.Handle("POST /bookings", handlers.BookingsCreate(br))

	http.Handle("GET /rooms", handlers.RoomsIndex(br))
	http.Handle("GET /price-calculator", handlers.CalculatePrice(br))

	if os.Getenv("FEATURE_FLAG_PAYMENTS_PROVIDER") == "sumup" {
		http.Handle("POST /sumup/checkouts", handlers.SumupCheckoutCreate(paymentsApi))
	} else if os.Getenv("FEATURE_FLAG_PAYMENTS_PROVIDER") == "stripe" {
		// stripe api endpoints
		http.Handle("POST /stripe/payment-intents", handlers.StripePaymentIntentCreate(br))
	}

	http.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(http.FS(static.StaticFiles))))
	http.Handle("GET /", handlers.BookingsNew(br))

	server := &http.Server{
		Addr:    ":" + EnvOrDefault("PORT", "8080"),
		Handler: handlers.Logging(http.DefaultServeMux),
	}

	log.Printf("Listening on port %s\n", EnvOrDefault("PORT", "8080"))
	log.Fatal(server.ListenAndServe())
}
