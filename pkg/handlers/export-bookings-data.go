package handlers

import (
	"encoding/csv"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

func ExportData(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		username, password, ok := r.BasicAuth()
		if !ok || username != "your-username" || password != "your-password" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return nil
		}
		// Set the headers for CSV download
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=export.csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Example data
		data := [][]string{
			{"Name", "Email"},
			{"Alice", "alice@example.com"},
			{"Bob", "bob@example.com"},
		}

		for _, row := range data {
			writer.Write(row)
		}

		return nil
	})
}
