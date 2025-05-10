package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"time"
	"text/template"
	"bytes"
)

type AdminBookingsDailyView struct {
	Date     string
	Bookings []AdminBooking
}

type AdminBooking struct {
	da.Booking
	BookingCodesMessage string
}

func NewAdminBooking(booking da.Booking, codesMessage string) AdminBooking {
	return AdminBooking {
		Booking: booking,
		BookingCodesMessage: codesMessage,
	}
}

type CustomerBookingCodesSMSData struct {
	Room string
	StartTime string 
	EndTime string 
	FrontDoorCode string
	RoomDoorCode string
	Cymbals int64
	FirstMessage bool
}

func AdminViewDailyBookings(br *da.BookingsRepository[da.StorageDriver], codes da.Codes) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		dateParam := r.URL.Query().Get("date")
		var selectedDate time.Time
		if dateParam == "" {
			// TODO: redirect to todays date with the date param default to today
			selectedDate = time.Now()
		} else {
			var err error
			selectedDate, err = time.Parse("2006-01-02", dateParam)
			if err != nil {
				return Error(err, http.StatusBadRequest)
			}
		}

		// Fetch all bookings for the selected date
		startOfDay := selectedDate.Truncate(24 * time.Hour)
		endOfDay := startOfDay.Add(24 * time.Hour)

		bookings, err := br.Where("status IN ('paid', 'unpaid', 'hold') and start_time >= ? and end_time <= ?", startOfDay, endOfDay)
		if err != nil {
			http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
			return nil
		}

		adminBookings := make([]AdminBooking, len(bookings))
		customerCodesMessageTemplate := `
Hey, we are looking forward to seeing you at Bad Habit for your rehearsal today!

Here are the details and information about your booking:
• Booking time: {{ .StartTime }}-{{ .EndTime }} 
• Location: {{ .Room }}
• Front door access code: {{ .FrontDoorCode }}#
• Room door access code: {{ .RoomDoorCode }}# {{ if .FirstMessage }}(Please note that the keypad to use for accessing rooms is the one found on the wall rather than on the door){{ end }}
{{ if eq .Cymbals 1 }}

You asked for Cymbals so they'll be left in the room :)
{{ end }}
{{ if .FirstMessage }}
• ROOM 1. This is the room directly in front of you as you walk into the studio.
• ROOM 2. This is the room on the right after you walk into the studio.
{{ end }}
Any questions or concerns, please get in touch!
		`
		for i, booking := range bookings {
			messageTmpl, err := template.New("codes-sms").Parse(customerCodesMessageTemplate)
			if err != nil {
				return Error(err, http.StatusInternalServerError)
			}

			previousBookings, err := br.Where("customer_phone = ? and start_time < ?", booking.CustomerPhone, booking.StartTime)
			if err != nil {
				return Error(err, http.StatusInternalServerError)
			}

			frontDoorCode, err := codes.FrontDoorCodeFor(booking.StartTime.Weekday())
			if err != nil {
				return Error(err, http.StatusInternalServerError)
			}

			roomDoorCode, err := codes.GetCode(booking.RoomName)
			if err != nil {
				return Error(err, http.StatusInternalServerError)
			}

			customerCodesMessageData := CustomerBookingCodesSMSData {
				Room: booking.RoomName,
				StartTime: booking.StartTime.Format("15:04"),
				EndTime: booking.EndTime.Format("15:04"),
				FrontDoorCode: frontDoorCode,
				RoomDoorCode: roomDoorCode,
				Cymbals: booking.Cymbals,
				FirstMessage: previousBookings == nil,
			}

			var messageContent bytes.Buffer
			if err := messageTmpl.Execute(&messageContent, customerCodesMessageData); err != nil {
				return Error(err, http.StatusInternalServerError)
			}
			adminBookings[i] = AdminBooking {
				Booking: booking,
				BookingCodesMessage: messageContent.String(),
			}
		}

		return Template("admin-view-daily-bookings.html.tmpl", AdminBookingsDailyView {
			Date: selectedDate.Format("2006-01-02"),
			Bookings: adminBookings,
		})
	})
}
