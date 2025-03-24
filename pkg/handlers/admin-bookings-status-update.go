package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
)

type AdminStatusUpdateForm struct {
	Status string `json:"status"`
}

// All bookings are created with a "hold" status
func AdminBookingsStatusUpdate(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		bookingId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return Error(err, http.StatusNotFound)
		}

		booking, err := br.Find(bookingId)
		if err != nil {
			return Error(err, 404)
		}

		form, err := ExtractForm[AdminStatusUpdateForm](r)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		booking.Status = form.Status

		_, err = br.Update([]da.Booking{*booking})
		if err != nil {
			return Error(err, http.StatusInternalServerError)

		}

		return Redirect("/admin/bookings")
	})
}
