package handlers

import (
	"net/http"
	_ "encoding/json"
	_ "bytes"
)

type CreateCheckoutRequest struct {
	
}

type SumupCheckout struct {
	Amount float32 `json:"amount"`
	CheckoutReference string `json:"checkout_reference"`
}

// Actual handlers
func SumupCheckoutCreate() Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		return JSON(SumupCheckout {
			Amount: 10,
			CheckoutReference: "test",
		})
	})
}
