package handlers

import (
	"fmt"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
	"time"
)

type UpdateBookingsForm struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Room      string `json:"room"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Cymbals   int64  `json:"cymbals"`
	Price     string `json:"price"`
	Status    string `json:"status"`
}

// All bookings are created with a "hold" status
func AdminBookingsUpdate(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		bookingId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return Error(err, http.StatusNotFound)
		}

		booking, err := br.Find(bookingId)
		if err != nil {
			return Error(err, 404)
		}

		form, err := ExtractForm[UpdateBookingsForm](r)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		fmt.Println(form)
		startTime, err := time.Parse("2006-01-02 15:04", form.Date+" "+form.StartTime)
		if err != nil {
			fmt.Println("start")
			return Error(err, http.StatusBadRequest)
		}

		endTime, err := time.Parse("2006-01-02 15:04", form.Date+" "+form.EndTime)
		if err != nil {
			fmt.Println("end")
			return Error(err, http.StatusBadRequest)
		}

		price, err := strconv.ParseFloat(form.Price, 64)
		if err != nil {
			fmt.Println("price")
			return Error(err, http.StatusBadRequest)
		}

		booking.Type = form.Type
		booking.CustomerName = form.Name
		booking.CustomerEmail = form.Email
		booking.CustomerPhone = form.Phone
		booking.RoomName = form.Room
		booking.StartTime = startTime
		booking.EndTime = endTime
		booking.Status = form.Status
		booking.Expiration = time.Now().Add(time.Minute * 15)
		booking.Price = price
		booking.Cymbals = form.Cymbals

		bookings, err := br.Update([]da.Booking{*booking})
		if err != nil {
			return Error(err, http.StatusInternalServerError)

		}

		fmt.Println(bookings)

		return JSON(bookings[0])
	})
}
