package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type CreateCheckoutForm struct {
	Amount float32 `json:"amount"`
	CheckoutReference string `json:"checkout_reference"`
}

type CreateSumupCheckoutRequest struct {
	Amount float32 `json:"amount"`
	CheckoutReference string `json:"checkout_reference"`
	Currency string `json:"currency"`
	MerchantCode string `json:"merchant_code"`
	ValidUntil time.Time `json:"valid_until"`
}

type SumupCheckout struct {
	Id string `json:"id"`
	Amount float32 `json:"amount"`
	CheckoutReference string `json:"checkout_reference"`
	Status string `json:"status"`
}

// Actual handlers
func SumupCheckoutCreate(sumupApi Api) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		form, err := ExtractForm[CreateCheckoutForm](r)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		checkoutReference := form.CheckoutReference
		if os.Getenv("APP_ENV") != "production" {
			checkoutReference = checkoutReference + "-dev"
		}

		response, err := sumupApi.Post("/v0.1/checkouts", CreateSumupCheckoutRequest {
			Amount: form.Amount,
			CheckoutReference: form.CheckoutReference,
			Currency: "GBP",
			MerchantCode: os.Getenv("SUMUP_MERCHANT_CODE"),
			ValidUntil: time.Now().Add(time.Minute * time.Duration(20)),
		})
		fmt.Println(response)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		if response.Status != 201 {
			return Error(fmt.Errorf("Error %d creating checkout with reference %s", response.Status, form.CheckoutReference), http.StatusInternalServerError)
		}

		var sumupCheckout SumupCheckout
		json.Unmarshal([]byte(response.Body), &sumupCheckout)
		return JSON(sumupCheckout)
	})
}
