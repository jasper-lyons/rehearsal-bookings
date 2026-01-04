package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	da "rehearsal-bookings/pkg/data_access"
)

var validRooms = map[string]bool{
	"Room 1":   true,
	"Room 2":   true,
	"Rec Room": true,
}

func EnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// Parse rooms from comma-separated string, removing quotes
func parseRooms(roomsStr string) ([]string, error) {
	rooms := strings.Split(roomsStr, ",")
	for i := range rooms {
		rooms[i] = strings.Trim(strings.TrimSpace(rooms[i]), `"`)
		if !validRooms[rooms[i]] {
			return nil, fmt.Errorf("invalid room name '%s'. Valid rooms are: Room 1, Room 2, Rec Room", rooms[i])
		}
	}
	return rooms, nil
}

// Parse time range from HH:MM-HH:MM format
func parseTimeRange(timeStr string) (start, end time.Time, err error) {
	parts := strings.Split(timeStr, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid time format, expected HH:MM-HH:MM")
	}

	startTime, err := time.Parse("15:04", strings.TrimSpace(parts[0]))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time: %v", err)
	}

	endTime, err := time.Parse("15:04", strings.TrimSpace(parts[1]))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time: %v", err)
	}

	return startTime, endTime, nil
}

// Parse dates from comma-separated YYYY/MM/DD format
func parseDates(datesStr string) ([]time.Time, error) {
	dateStrs := strings.Split(datesStr, ",")
	dates := make([]time.Time, len(dateStrs))

	for i, dateStr := range dateStrs {
		date, err := time.Parse("2006/01/02", strings.TrimSpace(dateStr))
		if err != nil {
			return nil, fmt.Errorf("invalid date format '%s', expected YYYY/MM/DD: %v", dateStr, err)
		}
		dates[i] = date
	}

	return dates, nil
}

// Check for booking conflicts using overlap detection
func checkConflicts(br *da.BookingsRepository[da.StorageDriver], rooms []string, dates []time.Time, startTime, endTime time.Time) ([]da.Booking, error) {
	var conflicts []da.Booking

	for _, room := range rooms {
		for _, date := range dates {
			// Combine date with time
			newStart := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
			newEnd := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

			// Query for overlapping bookings in this room
			// Two ranges overlap if: existing_start < new_end AND existing_end > new_start
			query := `room_name = ? AND status != 'cancelled' AND status != 'abandoned' AND start_time < ? AND end_time > ?`

			bookings, err := br.Where(query, room, newEnd, newStart)
			if err != nil {
				return nil, fmt.Errorf("error checking conflicts in room %s: %v", room, err)
			}

			conflicts = append(conflicts, bookings...)
		}
	}

	return conflicts, nil
}

func main() {
	// Manually parse customer flag from anywhere in args
	var customerFlag string
	var filteredArgs []string
	
	i := 0
	for i < len(os.Args) {
		if i == 0 {
			filteredArgs = append(filteredArgs, os.Args[i]) // Keep program name
		} else if os.Args[i] == "--customer" {
			if i+1 >= len(os.Args) {
				log.Fatal("--customer flag requires a value")
			}
			customerFlag = os.Args[i+1]
			i++ // Skip the flag value
		} else {
			filteredArgs = append(filteredArgs, os.Args[i])
		}
		i++
	}

	if len(filteredArgs) != 4 {
		fmt.Printf("Usage: %s [--customer \"Name, email, phone\"] \"Room1,Room2\" HH:MM-HH:MM YYYY/MM/DD,YYYY/MM/DD\n", os.Args[0])
		fmt.Printf("Example: %s \"Room 1,Room 2\" 16:00-19:00 2026/01/25,2026/02/22\n", os.Args[0])
		fmt.Printf("Example: %s \"Room 1\" 10:00-17:00 2026/01/29 --customer \"Mila Lyons, mila.lyons@gmail.com, 07463728168\"\n", os.Args[0])
		fmt.Printf("Valid rooms: Room 1, Room 2, Rec Room\n")
		os.Exit(1)
	}

	roomsStr := filteredArgs[1]
	timeRangeStr := filteredArgs[2]
	datesStr := filteredArgs[3]

	// Parse customer details
	customerName := "Jasper Lyons"
	customerEmail := "jasper.lyons@gmail.com" 
	customerPhone := "07506845146"

	if customerFlag != "" {
		customerParts := strings.Split(customerFlag, ",")
		if len(customerParts) != 3 {
			log.Fatal("Customer flag must contain exactly 3 comma-separated values: \"Name, email, phone\"")
		}
		customerName = strings.TrimSpace(customerParts[0])
		customerEmail = strings.TrimSpace(customerParts[1])
		customerPhone = strings.TrimSpace(customerParts[2])

		if customerName == "" || customerEmail == "" || customerPhone == "" {
			log.Fatal("Customer name, email, and phone cannot be empty")
		}
	}

	// Rest of the function stays exactly the same from here...
	// Parse arguments
	rooms, err := parseRooms(roomsStr)
	if err != nil {
		log.Fatal(err)
	}

	startTime, endTime, err := parseTimeRange(timeRangeStr)
	if err != nil {
		log.Fatal(err)
	}

	dates, err := parseDates(datesStr)
	if err != nil {
		log.Fatal(err)
	}

	// Load environment
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Set up database (following the same pattern as the existing main.go)
	var driver *da.SqliteDriver
	if os.Getenv("APP_ENV") == "production" {
		driver = da.NewSqliteDriver("db/production.db")
	} else {
		driver = da.NewSqliteDriver("db/development.db")
	}

	driver.Query("PRAGMA vdbe_debug=ON;")
	driver.Query("PRAGMA journal_mode=WAL;")
	driver.Query("PRAGMA wal_autocheckpoint=1000;")
	driver.Query("PRAGMA busy_timeout=5000;") // 5000 milliseconds
	driver.Query("PRAGMA wal_checkpoint(PASSIVE);")
	driver.Query("PRAGMA synchronous=NORMAL;")

	br := da.NewBookingsRepository(driver)

	// Check for conflicts
	conflicts, err := checkConflicts(br, rooms, dates, startTime, endTime)
	if err != nil {
		log.Fatal(err)
	}

	if len(conflicts) > 0 {
		fmt.Println("Booking conflicts detected:")
		for _, conflict := range conflicts {
			fmt.Printf("- %s from %s to %s in %s (ID: %d) - https://rehearsal-bookings.badhabitstudios.co.uk/admin/bookings/%d/edit\n",
				conflict.CustomerName,
				conflict.StartTime.Format("2006/01/02 15:04"),
				conflict.EndTime.Format("2006/01/02 15:04"),
				conflict.RoomName,
				conflict.Id,
				conflict.Id,
			)
		}
		os.Exit(1)
	}

	// Create bookings
	var bookings []da.Booking
	originalCommand := strings.Join(os.Args, " ")

	for _, room := range rooms {
		for _, date := range dates {
			start := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
			end := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

			booking := da.Booking{
				Type:           "band",
				CustomerName:   customerName,
				CustomerEmail:  customerEmail,
				CustomerPhone:  customerPhone,
				RoomName:       room,
				StartTime:      start,
				EndTime:        end,
				Status:         "unpaid",
				Price:          0,
				DiscountAmount: 0,
				BookingNotes:   fmt.Sprintf("Created with %s", originalCommand),
				PaymentMethod:  "internal",
			}

			bookings = append(bookings, booking)
		}
	}

	// Insert all bookings in a transaction
	_, err = br.Create(bookings)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create bookings: %v", err))
	}

	fmt.Printf("Successfully created %d bookings:\n", len(bookings))
	for _, booking := range bookings {
		fmt.Printf("- %s from %s to %s in %s\n",
			booking.CustomerName,
			booking.StartTime.Format("2006/01/02 15:04"),
			booking.EndTime.Format("2006/01/02 15:04"),
			booking.RoomName,
		)
	}
}
