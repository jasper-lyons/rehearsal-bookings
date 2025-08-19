package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

type AdminCodeNames struct {
	CodeNames []da.Code `json:"code_name"`
}

// Actual handlers
func AdminUpdateCodesView(cr *da.CodesRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		codes, err := cr.All()
		if err != nil {
			return Error(err, 500)
		}
		return Template("admin-codes-update.html.tmpl", AdminCodeNames{CodeNames: codes})
	})
}
