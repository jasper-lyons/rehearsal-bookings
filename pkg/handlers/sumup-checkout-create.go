package handlers

import (
	"net/http"
	"encoding/json"
	"fmt"
	"os"
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
}

type SumupCheckout struct {
	Id string `json:"id"`
	Amount float32 `json:"amount"`
	CheckoutReference string `json:"checkout_reference"`
}

// Actual handlers
func SumupCheckoutCreate(sumupApi Api) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		form, err := ExtractForm[CreateCheckoutForm](r)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		responseBody, err := sumupApi.Post("/v0.1/checkouts", CreateSumupCheckoutRequest {
			Amount: form.Amount,
			CheckoutReference: form.CheckoutReference,
			Currency: "GBP",
			MerchantCode: os.Getenv("SUMUP_MERCHANT_CODE"),
		})
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		fmt.Println(responseBody)

		var sumupCheckout SumupCheckout
		json.Unmarshal([]byte(responseBody), &sumupCheckout)
		return JSON(sumupCheckout)
	})
}
