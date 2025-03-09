package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
)

func AdminBookingsDelete(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		bookingId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return Error(err, http.StatusNotFound)
		}

		booking, err := br.Find(bookingId)
		if err != nil {
			return Error(err, 404)
		}

		_, err = br.Delete([]da.Booking { *booking })
		if err != nil {
			return Error(err, 500)
		}

		return Redirect("/admin/bookings")
	})
}
