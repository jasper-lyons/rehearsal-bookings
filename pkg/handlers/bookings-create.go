package handlers

import (
	"fmt"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"time"
)

type CreateBookingsForm struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Room      string `json:"room"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Duration  int    `json:"duration"`
	Cymbals   bool   `json:"cymbals"`
}

// All bookings are created with a "hold" status
func BookingsCreate(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		form, err := ExtractForm[CreateBookingsForm](r)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		startTime, err := time.Parse("2006-01-02 15:04", form.Date+" "+form.StartTime)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		endTime, err := time.Parse("2006-01-02 15:04", form.Date+" "+form.EndTime)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		price, err := BookingPrice(form.Type, startTime, endTime, form.Cymbals)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		booking := da.Booking{
			Type:          form.Type,
			CustomerName:  form.Name,
			CustomerEmail: form.Email,
			CustomerPhone: form.Phone,
			RoomName:      form.Room,
			StartTime:     startTime,
			EndTime:       endTime,
			Status:        "hold",
			Expiration:    time.Now().Add(time.Minute * 15),
			Price:         price,
			Cymbals:       form.Cymbals,
		}

		bookings, err := br.Create([]da.Booking{booking})
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		fmt.Println(bookings)

		return JSON(bookings[0])
	})
}
