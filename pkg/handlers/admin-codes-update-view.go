package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

// Actual handlers
func AdminUpdateCodesView(cr *da.CodesRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		return Template("admin-codes-update.html.tmpl", nil)
	})
}
