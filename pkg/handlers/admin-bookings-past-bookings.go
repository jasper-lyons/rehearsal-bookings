package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"sort"
	"time"
)

type AdminBookingsPast struct {
	Bookings []da.Booking
}

func AdminBookingsSearchPast(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		bookings, err := br.All()
		if err != nil {
			return Error(err, 500)
		}

		now := time.Now().Add(-1 * time.Hour)
		filteredBookings := []da.Booking{}
		for _, booking := range bookings {
			if booking.StartTime.Before(now) {
				filteredBookings = append(filteredBookings, booking)
			}
		}
		sort.Slice(filteredBookings, func(i, j int) bool {
			return filteredBookings[i].StartTime.Before(filteredBookings[j].StartTime)
		})

		return Template("admin-bookings-search-view.html.tmpl", AdminBookingsPast{Bookings: filteredBookings})
	})
}
