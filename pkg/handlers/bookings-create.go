
package handlers

import (
	"net/http"
	"time"
	da "rehearsal-bookings/pkg/data_access"
	"fmt"
)

type CreateBookingsForm struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Room string `json:"room"`
	Date string `json:"date"`
	StartTime string `json:"start_time"`
	Duration int `json:"duration"`
}

func BookingPrice(_type string, duration int) (float32, error) {
	switch _type {
		case "solo":
			return 6.5 * float32(duration), nil
		case "band":
			return 12 * float32(duration), nil
	}

	return -1.0, fmt.Errorf("Unknown rehearsal type: %s", _type)
}

// All bookings are created with a "hold" status
func BookingsCreate(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		form, err := ExtractForm[CreateBookingsForm](r)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		startTime, err := time.Parse("2006-01-02 15:04", form.Date + " " + form.StartTime)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		price, err := BookingPrice(form.Type, form.Duration)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}
		
		booking := da.Booking {
			Type: form.Type,
			CustomerName: form.Name,
			CustomerEmail: form.Email,
			CustomerPhone: form.Phone,
			RoomName: form.Room,
			StartTime: startTime,
			EndTime: startTime.Add(time.Hour * time.Duration(form.Duration)),
			Status: "hold",
			Expiration: time.Now().Add(time.Minute * 15),
			Price: price,
		}

		bookings, err := br.Create([]da.Booking { booking })
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		fmt.Println(bookings)

		return JSON(bookings[0])
	})
}
