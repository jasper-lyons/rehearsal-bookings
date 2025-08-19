package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

type AdminUpdateCodesForm struct {
	CodeName  string `json:"code_name"`
	CodeValue string `json:"code_value"`
	Notes     string `json:"code_notes"` // New field for notes
}

// All bookings are created with a "hold" status
func AdminUpdateCodes(cr *da.CodesRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {

		form, err := ExtractForm[AdminUpdateCodesForm](r)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		code, err := cr.Find(form.CodeName)
		if err != nil {
			return Error(err, 404)
		}

		code.CodeValue = form.CodeValue
		code.Notes = form.Notes

		codes, err := cr.Update([]da.Code{*code})
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		return JSON(codes[0])
	})
}
