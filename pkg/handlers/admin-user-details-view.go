package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
)

type AdminUsers struct {
	Users []da.User
}

func AdminUserDetailsViews(ur *da.UsersRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		users, err := ur.All()
		if err != nil {
			return Error(err, 500)
		}

		return Template("admin-users-view.html.tmpl", AdminUsers{Users: users})
	})
}
