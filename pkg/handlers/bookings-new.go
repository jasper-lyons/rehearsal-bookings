package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

type BookingsNewView struct {
	CRSFToken string
}

// Actual handlers
func BookingsNew(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		return Template("bookings-new.html.tmpl", nil)
	})
}
