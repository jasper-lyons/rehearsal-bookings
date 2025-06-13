package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

		// Create a Title caser for English language
		titleCaser := cases.Title(language.English)

		for i := range users {
			// Use the caser to transform the username
			users[i].UserName = titleCaser.String(users[i].UserName)
		}

		return Template("admin-bookings-new.html.tmpl", NewBookingUsers{Users: users})
	})
}
