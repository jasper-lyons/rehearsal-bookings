package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

type AdminBookingsIndexView struct {
	Bookings []da.Booking
}

func AdminBookingsIndex(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		bookings, err := br.All()
		if err != nil {
			return Error(err, 500)
		}

		return Template("admin-bookings-index.html.tmpl", AdminBookingsIndexView{Bookings: bookings})
	})
}
