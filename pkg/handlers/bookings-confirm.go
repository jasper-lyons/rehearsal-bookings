package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	"time"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
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

func BookingsConfirm(br *da.BookingsRepository[da.StorageDriver], sumupApi Api) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		form, err := ExtractForm[SumupCheckoutProcessSuccessForm](r)
		fmt.Println(form)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		bookingId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return Error(err, http.StatusNotFound)
		}

		booking, err := br.Find(bookingId)
		if err != nil {
			return Error(err, http.StatusNotFound)
		}

		response, err := sumupApi.Get(
			"/v2.1/merchants/%s/transactions?id=%s",
			os.Getenv("SUMUP_MERCHANT_CODE"),
			form.TransactionId,
		)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		if response.Status != 200 {
			return Error(fmt.Errorf("Error %d retrieving transcation %s", response.Status, form.TransactionId), http.StatusInternalServerError)
		}

		var sumupTransaction SumupTransaction
		json.Unmarshal([]byte(response.Body), &sumupTransaction)

		if sumupTransaction.Status != "SUCCESSFUL" {
			return Error(fmt.Errorf("Cannot confirm booking, payment unsuccessful"), http.StatusInternalServerError)
		}

		booking.Status = "paid"
		booking.TransactionId = form.TransactionId
		br.Update([]da.Booking { *booking })

		return JSON(booking)
	})
}
