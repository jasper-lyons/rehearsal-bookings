package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"time"
)

type AdminCodesView struct {
	Codes []AdminCode
}

type AdminCode struct {
	da.Code
	Expired bool
}

func AdminViewCodes(cr *da.CodesRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		codes, err := cr.All()
		if err != nil {
			return Error(err, 500)
		}

		var adminCodes []AdminCode

		for _, code := range codes {
			adminCodes = append(adminCodes, AdminCode{
				Code:    code,
				Expired: time.Now().AddDate(0, -1, 0).After(code.UpdatedAt),
			})
		}

		return Template("admin-codes-view.html.tmpl", AdminCodesView{Codes: adminCodes})
	})
}
