package handlers

import (
	"errors"
	"fmt"
	"net/http"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
	"time"
)

type PriceRequest struct {
	Type      string `json:"type"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Cymbals   int    `json:"cymbals"`
}

type PriceResponse struct {
	Price float64 `json:"price"`
	Error string  `json:"error,omitempty"`
}

func BookingPrice(_type string, startTime time.Time, endTime time.Time, cymbals int) (float64, error) {
	duration := endTime.Sub(startTime).Hours()
	var price float64

	switch _type {
	case "solo":
		price = 6.5 * duration
	case "band":
		if duration > 9 {
			price = 100.00
		} else if duration > 3 {
			price = 10.00 * duration
		} else {
			price = 12.00 * duration
		}
	default:
		return -1.0, fmt.Errorf("Unknown rehearsal type: %s", _type)
	}

	if cymbals == 1 {
		price += 3.0
	}

	return price, nil
}

func CalculatePrice(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		startTime, err := time.Parse("2006-01-02 15:04", r.URL.Query().Get("startTime"))
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		endTime, err := time.Parse("2006-01-02 15:04", r.URL.Query().Get("endTime"))
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}

		SessionType := r.URL.Query().Get("type")
		if SessionType == "" {
			return Error(errors.New("Missing 'Type' query parameter."), http.StatusBadRequest)
		}

		cymbals64, err := strconv.ParseInt(r.URL.Query().Get("cymbals"), 10, 64)
		if err != nil {
			return Error(errors.New("Invalid 'cymbals' query parameter."), http.StatusBadRequest)
		}

		fmt.Println(SessionType)
		fmt.Println(startTime)
		fmt.Println(endTime)
		fmt.Println(cymbals64)

		price, err := BookingPrice(SessionType, startTime, endTime, int(cymbals64))
		if err != nil {
			return Error(err, http.StatusBadRequest)
		}
		fmt.Println(price)

		return JSON(PriceResponse{Price: price})
	})
}
