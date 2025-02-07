package main

import (
	"log"
	"os"
	"fmt"
	handlers "rehearsal-bookings/pkg/handlers"
	"github.com/joho/godotenv"
	"encoding/json"
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

	sumupApi := handlers.NewApi("https://api.sumup.com", map[string]string {
		"Content-Type": "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("SUMUP_API_KEY")),
	})

	response, err := sumupApi.Get("/v0.1/checkouts")
	if err != nil {
		log.Fatal("error! can't list checkouts")
	}


	var checkouts []handlers.SumupCheckout
	json.Unmarshal([]byte(response.Body), &checkouts)

	for _, checkout := range checkouts {
		if checkout.Status == "PENDING" || checkout.Status == "FAILED" {
			fmt.Println("Clearing")
			fmt.Println(checkout)
			sumupApi.Delete("/v1.0/checkouts/%s", checkout.Id)
		}
	}
}
