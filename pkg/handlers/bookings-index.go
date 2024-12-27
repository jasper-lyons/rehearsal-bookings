package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)


type BookingsIndexView struct {
	Bookings []da.Booking
}

// Actual handlers
func BookingsIndex(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		bookings, err := br.All()
		if err != nil {
			return Error(err, 500)
		}

		return Template("bookings-index.html.tmpl", BookingsIndexView{Bookings: bookings})
	})
}
