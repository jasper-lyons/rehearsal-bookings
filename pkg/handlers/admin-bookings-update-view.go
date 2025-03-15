package handlers

import (
	"log"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
)

type AdminBookingsUpdateData struct {
	Bookings []da.Booking
}

// Actual handlers
func AdminBookingsUpdateView(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		bookingId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return Error(err, http.StatusNotFound)
		}

		booking, err := br.Find(bookingId)
		if err != nil {
			return Error(err, 404)
		}
		log.Printf("Found booking: %v", []da.Booking{*booking})
		return Template("admin-bookings-update.html.tmpl", AdminBookingsUpdateData{Bookings: []da.Booking{*booking}})
	})
}
