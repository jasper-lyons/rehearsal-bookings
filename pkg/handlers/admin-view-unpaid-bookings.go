package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"sort"
)

type AdminUnpaidBookingsSearch struct {
	Bookings []da.Booking
}

func AdminViewUnpaidBookings(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		bookings, err := br.Where("status = 'unpaid' AND payment_method != 'booking_press'")
		if err != nil {
			return Error(err, 500)
		}

		sort.Slice(bookings, func(i, j int) bool {
			return bookings[i].StartTime.Before(bookings[j].StartTime)
		})

		return Template("admin-view-table-template.html.tmpl", AdminUnpaidBookingsSearch{Bookings: bookings})
	})
}
