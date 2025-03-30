package handlers

import (
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"time"
)

type TimeSlot struct {
	StartTime     time.Time
	EndTime       time.Time
	Room1Bookings []da.Booking
	Room2Bookings []da.Booking
}

type DailyBookings struct {
	Date      string
	TimeSlots []TimeSlot
}

func AdminDailyCalendar(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		dateParam := r.URL.Query().Get("date")
		var selectedDate time.Time
		if dateParam == "" {
			selectedDate = time.Now()
		} else {
			var err error
			selectedDate, err = time.Parse("2006-01-02", dateParam)
			if err != nil {
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return nil
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

		// Generate time slots
		var timeSlots []TimeSlot
		startHour := 10 // Start at 10:00 AM
		endHour := 22   // End at 10:00 PM
		for hour := startHour; hour < endHour; hour++ {
			startTime := startOfDay.Add(time.Duration(hour) * time.Hour)
			endTime := startTime.Add(time.Hour)

			// Filter bookings for this time slot and room
			var room1Bookings, room2Bookings []da.Booking
			for _, booking := range bookings {
				if booking.StartTime.Before(endTime) && booking.EndTime.After(startTime) {
					if booking.RoomName == "Room 1" {
						room1Bookings = append(room1Bookings, booking)
					} else if booking.RoomName == "Room 2" {
						room2Bookings = append(room2Bookings, booking)
					}
				}
			}

			// Add the time slot to the list
			timeSlots = append(timeSlots, TimeSlot{
				StartTime:     startTime,
				EndTime:       endTime,
				Room1Bookings: room1Bookings,
				Room2Bookings: room2Bookings,
			})
		}

		// Prepare the data for the template
		dailyBookings := DailyBookings{
			Date:      selectedDate.Format("2006-01-02"),
			TimeSlots: timeSlots,
		}

		return Template("admin-view-daily-calendar.html.tmpl", dailyBookings)
	})
}
