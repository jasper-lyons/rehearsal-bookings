package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"sort"
)

type AdminBookingsSearchView struct {
	Bookings []da.Booking
}

func AdminBookingsSearch(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		bookings, err := br.All()
		if err != nil {
			return Error(err, 500)
		}

		sort.Slice(bookings, func(i, j int) bool {
			return bookings[i].StartTime.Before(bookings[j].StartTime)
		})

		return Template("admin-bookings-search-view.html.tmpl", AdminBookingsSearchView{Bookings: bookings})
	})
}
