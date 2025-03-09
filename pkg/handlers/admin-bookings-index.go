package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"sort"
)

type AdminBookingsIndexView struct {
	GroupedBookings []GroupedBookings
}

type GroupedBookings struct {
	Date     string
	Bookings []da.Booking
}

func AdminBookingsIndex(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		bookings, err := br.All()
		if err != nil {
			return Error(err, 500)
		}

		sort.Slice(bookings, func(i, j int) bool {
			return bookings[i].StartTime.Before(bookings[j].StartTime)
		})

		// Group bookings by date
		groupedBookings := groupBookingsByDate(bookings)

		return Template("admin-bookings-index.html.tmpl", AdminBookingsIndexView{GroupedBookings: groupedBookings})
	})
}

// function to group bookings by date - called from AdminBookingsIndex
func groupBookingsByDate(bookings []da.Booking) []GroupedBookings {
	grouped := make(map[string][]da.Booking)
	for _, booking := range bookings {
		date := booking.StartTime.Format("2006-01-02")
		grouped[date] = append(grouped[date], booking)
	}

	var groupedBookings []GroupedBookings
	for date, bookings := range grouped {
		groupedBookings = append(groupedBookings, GroupedBookings{Date: date, Bookings: bookings})
	}

	// Sort grouped bookings by date
	sort.Slice(groupedBookings, func(i, j int) bool {
		return groupedBookings[i].Date < groupedBookings[j].Date
	})

	return groupedBookings
}
