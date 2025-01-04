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

// Actual handlers
func BookingsConfirm(br *da.BookingsRepository[da.StorageDriver], sumupApi Api) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		form, err := ExtractForm[SumupCheckoutProcessSuccessForm](r)
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

		fmt.Println(fmt.Sprintf("booking-%d", booking.Id))
		fmt.Println(form.CheckoutReference)

		if fmt.Sprintf("booking-%d", booking.Id) != form.CheckoutReference {
			return Error(fmt.Errorf("Cannot confirm booking"), http.StatusInternalServerError)
		}

		// TODO: Shift param serialisation into api model
		url := fmt.Sprintf(
			"/v2.1/merchants/%s/transactions?id=%s",
			os.Getenv("SUMUP_MERCHANT_CODE"),
			form.TransactionId,
		)
		responseBody, err := sumupApi.Get(url)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		var sumupTransaction SumupTransaction
		json.Unmarshal([]byte(responseBody), &sumupTransaction)

		if sumupTransaction.Status != "SUCCESSFUL" {
			return Error(fmt.Errorf("Cannot confirm booking, payment unsuccessful"), http.StatusInternalServerError)
		}

		booking.Status ="paid"
		br.Update([]da.Booking { *booking })

		return JSON(booking)
	})
}
