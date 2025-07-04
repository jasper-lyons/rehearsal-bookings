package handlers

import (
	"encoding/csv"
	"errors"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
	"time"
)

func ExportData(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {

		startDay := r.URL.Query().Get("start-day")
		if startDay == "" {
			return Error(errors.New("Missing 'start-day' query parameter."), http.StatusBadRequest)
		}

		startDayFormatted, err := time.Parse("2006-01-02", startDay)
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		endDayFormatted := time.Now()
		endDay := r.URL.Query().Get("end-day")
		if endDay != "" {
			endDayFormatted, err = time.Parse("2006-01-02", endDay)
			if err != nil {
				return Error(err, http.StatusBadRequest)
			}
		}
		bookings, err := br.Where("status != 'hold' and start_time >= ? and end_time <= ?", startDayFormatted, endDayFormatted)
		if err != nil {
			return Error(err, 500)
		}

		// Set the headers for CSV download
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=export.csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Write the header row to the CSV
		writer.Write([]string{
			"Booking ID",
			"Type",
			"Customer Name",
			"Room Name",
			"Start Time",
			"End Time",
			"Status",
			"Cymbals",
			"Price",
			"Discount Amount",
			"Paid At",
			"Payment Method",
			"Transaction ID",
			"Updated At",
			"Created At",
			"Cancelled At",
		})

		// Write each booking to the CSV
		for _, booking := range bookings {
			writer.Write([]string{
				strconv.FormatInt(booking.Id, 10),
				booking.Type,
				booking.CustomerName,
				booking.RoomName,
				booking.StartTime.Format("2006-01-02 15:04:05"),
				booking.EndTime.Format("2006-01-02 15:04:05"),
				booking.Status,
				strconv.FormatInt(booking.Cymbals, 10),
				strconv.FormatFloat(booking.Price, 'f', 2, 64),
				strconv.FormatFloat(booking.DiscountAmount, 'f', 2, 64),
				booking.PaidAt.Format("2006-01-02 15:04:05"),
				booking.PaymentMethod,
				booking.TransactionId,
				booking.UpdatedAt.Format("2006-01-02 15:04:05"),
				booking.CreatedAt.Format("2006-01-02 15:04:05"),
				booking.CancelledAt.Format("2006-01-02 15:04:05"),
			})
		}

		return nil
	})
}
