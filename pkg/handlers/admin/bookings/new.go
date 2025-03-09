package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	handlers "rehearsal-bookings/pkg/handlers"
)

type AdminBookingsNewView struct {
	CRSFToken string
}

// Actual handlers
func AdminBookingsNew(br *da.BookingsRepository[da.StorageDriver]) handlers.Handler {
	return handlers.Handler(func(w http.ResponseWriter, r *http.Request) handlers.Handler {
		return handlers.Template("admin/bookings/new.html.tmpl", nil)
	})
}
