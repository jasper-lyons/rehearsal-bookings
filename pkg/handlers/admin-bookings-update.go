package handlers

import (
	"fmt"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
	"time"
)

type AdminUpdateBookingsForm struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Room          string `json:"room"`
	Date          string `json:"date"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	Cymbals       int64  `json:"cymbals"`
	RevisedPrice  string `json:"revised_price"`
	Status        string `json:"status"`
	BookingNotes  string `json:"booking_notes"`
	PaymentMethod string `json:"payment_method"`
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

		form, err := ExtractForm[AdminUpdateBookingsForm](r)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

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

		price, err := BookingPrice(form.Type, startTime, endTime, int(form.Cymbals))
		if err != nil {
			fmt.Println("price")
			return Error(err, http.StatusBadRequest)
		}

		var discount_amount float64 = 0
		if form.RevisedPrice != "" {
			revised_price, err := strconv.ParseFloat(form.RevisedPrice, 64)
			if err != nil {
				fmt.Println("discount price")
				return Error(err, http.StatusBadRequest)
			}
			discount_amount = float64(price - revised_price)
			price = revised_price

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
		booking.DiscountAmount = discount_amount
		booking.Cymbals = form.Cymbals
		booking.BookingNotes = form.BookingNotes
		booking.PaymentMethod = form.PaymentMethod

		bookings, err := br.Update([]da.Booking{*booking})
		if err != nil {
			return Error(err, http.StatusInternalServerError)

		}

		return JSON(bookings[0])
	})
}
