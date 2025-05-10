package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

type AdminCodesView struct {
	Codes []da.Code
}

func AdminViewCodes(cr *da.CodesRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		codes, err := cr.All()
		if err != nil {
			return Error(err, 500)
		}

		return Template("admin-codes-view.html.tmpl", AdminCodesView{Codes: codes})
	})
}
