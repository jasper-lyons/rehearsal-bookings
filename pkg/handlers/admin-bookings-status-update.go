package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
	"time"
)

type AdminStatusUpdateForm struct {
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
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

		if form.Status == "paid" {
			booking.PaymentMethod = form.PaymentMethod
			booking.PaidAt = time.Now()
		} else {
			booking.CancelledAt = time.Now()
		}
		booking.Status = form.Status

		_, err = br.Update([]da.Booking{*booking})
		if err != nil {
			return Error(err, http.StatusInternalServerError)

		}

		return Redirect("/admin/bookings")
	})
}
