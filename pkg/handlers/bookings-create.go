package handlers

import (
	"errors"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"time"
	"regexp"
	"strings"
)

type CreateBookingsForm struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Room         string `json:"room"`
	Date         string `json:"date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Cymbals      int64  `json:"cymbals"`
	BookingNotes string `json:"booking_notes"`
}

func NormalizePhoneNumber(phoneNumber string) (string, error) {
	// Remove all spaces, hyphens, parentheses, and dots
	reg := regexp.MustCompile(`[\s\-\(\)\.]`)
	phoneNumber = reg.ReplaceAllString(phoneNumber, "")
	
	// Check if the string is empty after cleaning
	if phoneNumber == "" {
		return "", errors.New("phone number is empty")
	}
	
	// Case: Already in international format with + (e.g., +447XXXXXXXXX)
	if strings.HasPrefix(phoneNumber, "+447") && len(phoneNumber) == 13 {
		return phoneNumber, nil
	}
	
	// Case: International format without + (e.g., 447XXXXXXXXX)
	if strings.HasPrefix(phoneNumber, "447") && len(phoneNumber) == 12 {
		return "+" + phoneNumber, nil
	}
	
	// Case: Domestic format starting with 07 (e.g., 07XXXXXXXXX)
	if strings.HasPrefix(phoneNumber, "07") && len(phoneNumber) == 11 {
		return "+44" + phoneNumber[1:], nil
	}
	
	// Case: Just the number part without any prefix (e.g., 7XXXXXXXXX)
	if strings.HasPrefix(phoneNumber, "7") && len(phoneNumber) == 10 {
		return "+44" + phoneNumber, nil
	}
	
	// If it's a UK number starting with +44 but not a mobile (not +447...)
	if strings.HasPrefix(phoneNumber, "+44") && !strings.HasPrefix(phoneNumber, "+447") {
		return "", errors.New("not a UK mobile number: does not start with +447")
	}
	
	// For non-UK international numbers that start with +
	if strings.HasPrefix(phoneNumber, "+") && !strings.HasPrefix(phoneNumber, "+44") {
		return "", errors.New("not a UK number: does not start with +44")
	}
	
	// If we reach here, the number doesn't match any of our expected formats
	return "", errors.New("invalid UK mobile number format")
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

		price, err := BookingPrice(form.Type, startTime, endTime, int(form.Cymbals))
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		phone, err := NormalizePhoneNumber(form.Phone)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		// check that no bookings overlap!
		bookings, err := br.Where("(status in ('paid', 'unpaid', 'hold')) and room_name = ? and (end_time > ? and start_time < ?)", form.Room, startTime, endTime)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		if len(bookings) > 0 {
			return Error(errors.New("Can't book, there are overlapping bookings!"), http.StatusBadRequest)
		}

		booking := da.Booking {
			Type:          form.Type,
			CustomerName:  form.Name,
			CustomerEmail: form.Email,
			CustomerPhone: phone, 
			RoomName:      form.Room,
			StartTime:     startTime,
			EndTime:       endTime,
			Status:        "hold",
			Expiration:    time.Now().Add(time.Minute * 15),
			Price:         price,
			Cymbals:       form.Cymbals,
			BookingNotes:  form.BookingNotes,
			PaymentMethod: "online",
		}

		bookings, err = br.Create([]da.Booking{booking})
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		return JSON(bookings[0])
	})
}
