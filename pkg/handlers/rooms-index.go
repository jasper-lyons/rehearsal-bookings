package handlers

import (
	"errors"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"time"
)

type RoomsIndexForm struct {
	Date string `json:"date"`
}

type Room struct {
	Name         string          `json:"name"`
	Availability map[string]bool `json:"availability"`
}

func NewRoom(name string) Room {
	return Room{
		Name: name,
		Availability: map[string]bool{
			"00:00": true,
			"01:00": true,
			"02:00": true,
			"03:00": true,
			"04:00": true,
			"05:00": true,
			"06:00": true,
			"07:00": true,
			"08:00": true,
			"09:00": true,
			"10:00": true,
			"11:00": true,
			"12:00": true,
			"13:00": true,
			"14:00": true,
			"15:00": true,
			"16:00": true,
			"17:00": true,
			"18:00": true,
			"19:00": true,
			"20:00": true,
			"21:00": true,
			"22:00": true,
			"23:00": true,
		},
	}
}

type RoomsIndexView struct {
	Rooms []Room `json:"rooms"`
}

func RoomsIndex(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		day := r.URL.Query().Get("day")
		if day == "" {
			return Error(errors.New("Missing 'day' query parameter."), http.StatusBadRequest)
		}

		dayStart, err := time.Parse("2006-01-02", day)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		dayEnd := dayStart.Add(time.Hour * 24 * time.Duration(1))

		bookings, err := br.Where("status IN ('paid', 'unpaid', 'hold') and start_time >= ? and end_time <= ?", dayStart, dayEnd)
		if err != nil {
			return Error(err, 500)
		}

		room1 := NewRoom("Room 1")
		for _, booking := range bookings {
			if booking.RoomName == "Room 1" {
				startTime := booking.StartTime
				// This truncates the float64 into an int so we're assuming accurate, whole hour maths...
				hours := int(booking.EndTime.Sub(startTime).Hours())
				for i := range hours {
					bookedHour := startTime.Add(time.Hour * time.Duration(i))
					room1.Availability[bookedHour.Format("15:04")] = false
				}
			}
		}

		room2 := NewRoom("Room 2")
		for _, booking := range bookings {
			if booking.RoomName == "Room 2" {
				startTime := booking.StartTime
				// This truncates the float64 into an int so we're assuming accurate, whole hour maths...
				hours := int(booking.EndTime.Sub(startTime).Hours())
				for i := range hours {
					bookedHour := startTime.Add(time.Hour * time.Duration(i))
					room2.Availability[bookedHour.Format("15:04")] = false
				}
			}
		}

		return JSON(RoomsIndexView{Rooms: []Room{room1, room2}})
	})
}
