package handlers

import (
	"encoding/csv"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

func ExportData(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
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
