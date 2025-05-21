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
			expiryDuration := time.Now().AddDate(0, -1, 0) // Default expiry duration (1 month)
			if code.CodeName == "Room 1" || code.CodeName == "Room 2" {
				expiryDuration = time.Now().Add(-72 * time.Hour) // 72 hours expiry for Room 1 and Room 2
			}

			adminCodes = append(adminCodes, AdminCode{
				Code:    code,
				Expired: expiryDuration.After(code.UpdatedAt),
			})
		}

		return Template("admin-codes-view.html.tmpl", AdminCodesView{Codes: adminCodes})
	})
}
