package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	"time"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
	"errors"
)

type SumupCheckoutProcessSuccessForm struct {
	CheckoutId string `json:"id"`
	Amount float32 `json:"amount"`
	CheckoutReference string `json:"checkout_reference"`
	PaidAt time.Time `json:"paid_at"`
	Status string `json:"status"`
	TransactionCode string `json:"transation_code"`
	TransactionId string `json:"transaction_id"`
}

type GetSumupTransactionRequest struct {
	Amount float32 `json:"amount"`
	CheckoutReference string `json:"checkout_reference"`
	Currency string `json:"currency"`
	MerchantCode string `json:"merchant_code"`
}

type SumupTransaction struct {
	Status string `json:"status"`
}

func ConfirmSumupPayment(br *da.BookingsRepository[da.StorageDriver], sumupApi Api, r *http.Request) (string, Handler) {
	form, err := ExtractForm[SumupCheckoutProcessSuccessForm](r)
	if err != nil {
		return "", Error(err, http.StatusInternalServerError)
	}

	response, err := sumupApi.Get(
		"/v2.1/merchants/%s/transactions?id=%s",
		os.Getenv("SUMUP_MERCHANT_CODE"),
		form.TransactionId,
	)
	if err != nil {
		return "", Error(err, http.StatusInternalServerError)
	}

	if response.Status != 200 {
		return "", Error(fmt.Errorf("Error %d retrieving transcation %s", response.Status, form.TransactionId), http.StatusInternalServerError)
	}

	var sumupTransaction SumupTransaction
	json.Unmarshal([]byte(response.Body), &sumupTransaction)

	if sumupTransaction.Status != "SUCCESSFUL" {
		return "", Error(fmt.Errorf("Cannot confirm booking, payment unsuccessful"), http.StatusInternalServerError)
	}

	return form.TransactionId, nil
}

func ConfirmStripePayment(br *da.BookingsRepository[da.StorageDriver], api Api, r *http.Request) (string, Handler) {
	return "", nil
}

func ConfirmPayment(br *da.BookingsRepository[da.StorageDriver], api Api, r *http.Request) (string, Handler) {
	if os.Getenv("FEATURE_FLAG_PAYMENTS_PROVIDER") == "sumup" {
		return ConfirmSumupPayment(br, api, r)
	} else if os.Getenv("FEATURE_FLAG_PAYMENTS_PROVIDER") == "stripe" {
		return ConfirmStripePayment(br, api, r)
	} else {
		return "", Error(errors.New("No payment provider configured"), http.StatusInternalServerError)
	}
}

func BookingsConfirm(br *da.BookingsRepository[da.StorageDriver], api Api) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		transactionId, errorHandler := ConfirmPayment(br, api, r)
		if errorHandler != nil {
			return errorHandler
		}

		bookingId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return Error(err, http.StatusNotFound)
		}

		booking, err := br.Find(bookingId)
		if err != nil {
			// if we've confirmed the payment and the booking doesn't exist... we should probably refund the user?
			return Error(err, http.StatusNotFound)
		}

		booking.Status = "paid"
		booking.TransactionId = transactionId 
		br.Update([]da.Booking { *booking })

		return JSON(booking)
	})
}
