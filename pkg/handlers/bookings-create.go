
package handlers

import (
	"net/http"
	"time"
	"strconv"
	da "rehearsal-bookings/pkg/data_access"
)

// Actual handlers
func BookingsCreate(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		if err := r.ParseForm(); err != nil {
			return Error(err, http.StatusBadRequest)
		}

		startTime, err := time.Parse("2006-01-02 15:04", r.FormValue("date") + " " + r.FormValue("start_time"))
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}
		
		hours, err := strconv.ParseInt(r.FormValue("duration"), 10, 0)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		booking := da.Booking {
			CustomerName: r.FormValue("name"),
			CustomerEmail: r.FormValue("email"),
			RoomName: r.FormValue("room"),
			StartTime: startTime,
			EndTime: startTime.Add(time.Hour * time.Duration(hours)),
		}

		_, err = br.Create([]da.Booking { booking })
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		return Redirect("/")
	})
}
