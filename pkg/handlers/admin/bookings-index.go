package admin

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	handlers "rehearsal-bookings/pkg/handlers"
)

type BookingsIndexView struct {
	Bookings []da.Booking
}

func BookingsIndex(br *da.BookingsRepository[da.StorageDriver]) handlers.Handler {
	return handlers.Handler(func (w http.ResponseWriter, r *http.Request) handlers.Handler {
		bookings, err := br.All()
		if err != nil {
			return handlers.Error(err, 500)
		}

		return handlers.Template("admin/bookings-index.html.tmpl", BookingsIndexView{Bookings: bookings})
	})
}
