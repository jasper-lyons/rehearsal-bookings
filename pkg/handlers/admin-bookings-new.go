package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

type AdminBookingsNewView struct {
	CRSFToken string
}

// Actual handlers
func AdminBookingsNew(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		return Template("admin-bookings-new.html.tmpl", nil)
	})
}
