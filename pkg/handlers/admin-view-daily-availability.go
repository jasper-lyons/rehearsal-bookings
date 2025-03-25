package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

// Actual handlers
func AdminViewDailyAvailability(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		return Template("admin-view-booking-availabililty.html.tmpl", nil)
	})
}
