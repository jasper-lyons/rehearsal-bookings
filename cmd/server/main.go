package main

import (
	"net/http"
	"log"
	"os"
	// GO Embed doesn't support embedding files from outside the module boundary
	// (in this case, anything outside of this directory) but we want to store
	// templates and static files at the route directory so we need to treat them
	// as their own go modules (you'll see the main.go files in the web/static and
	// web/templates directories) so that we can import them here and access the 
	// embedded static files!
	// "rehearsal-bookings/web/static"
	templates "rehearsal-bookings/web/templates"
)

func EnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}


func main() {
	http.HandleFunc("GET /", func (rw http.ResponseWriter, r *http.Request) {
		templates.Render(rw, "bookings.html.tmpl", nil)
	})

	log.Fatal(http.ListenAndServe(EnvOrDefault("PORT", ":8080"), nil))
}
