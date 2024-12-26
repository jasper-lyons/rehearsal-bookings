package main

import (
	"net/http"
	"log"
	"fmt"
	"os"
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
		fmt.Fprintf(rw, "Hello!")
	})

	log.Fatal(http.ListenAndServe(EnvOrDefault("PORT", ":8080"), nil))
}
