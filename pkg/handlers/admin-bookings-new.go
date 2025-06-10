package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

type NewBookingUsers struct {
	Users []da.User
}

// Actual handlers
func AdminBookingsNew(br *da.BookingsRepository[da.StorageDriver], ur *da.UsersRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		users, err := ur.All()
		if err != nil {
			return Error(err, 500)
		}

		return Template("admin-bookings-new.html.tmpl", NewBookingUsers{Users: users})
	})
}
