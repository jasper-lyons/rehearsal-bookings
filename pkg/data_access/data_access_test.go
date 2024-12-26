package data_access

import (
	"testing"
	"fmt"
)

func TestBookings(t *testing.T) {
	driver := NewSqliteDriver("../../db/development.db")
	br := NewBookingsRepository(driver)

	bookings, err := br.Create(
		[]Booking {
			{ CustomerName: "Ben", CustomerEmail: "ben@name.com", RoomName: "Room 1", StartTime: FuckTheError(Time("2025-01-01 19:00:00")), EndTime: FuckTheError(Time("2025-01-01 21:00:00")) },
		},
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(bookings)

	bookings, err = br.Find(bookings[0].Id)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(bookings)

	bookings, err = br.All()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(bookings)

	bookings, err = br.Where("start_time < '2025-01-02'")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(bookings)

	bookings[0].CustomerName = "Tom"
	bookings, err = br.Update(bookings)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(bookings)

	bookings, err = br.Delete(bookings)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(bookings)
}
