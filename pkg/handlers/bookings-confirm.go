package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	da "rehearsal-bookings/pkg/data_access"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

type SumupCheckoutProcessSuccessForm struct {
	CheckoutId        string    `json:"id"`
	Amount            float32   `json:"amount"`
	CheckoutReference string    `json:"checkout_reference"`
	PaidAt            time.Time `json:"paid_at"`
	Status            string    `json:"status"`
	TransactionCode   string    `json:"transation_code"`
	TransactionId     string    `json:"transaction_id"`
}

type GetSumupTransactionRequest struct {
	Amount            float32 `json:"amount"`
	CheckoutReference string  `json:"checkout_reference"`
	Currency          string  `json:"currency"`
	MerchantCode      string  `json:"merchant_code"`
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

type StripePaymentConfirmationForm struct {
	PaymentId     string  `json:"payment_id"`
	PaymentStatus string  `json:"payment_status"`
	PaymentAmount float64 `json:"payment_amount"`
	PaymentMethod string  `json:"payment_method"`
}

func ConfirmStripePayment(br *da.BookingsRepository[da.StorageDriver], r *http.Request) (string, Handler) {
	// Extract the payment confirmation form
	form, err := ExtractForm[StripePaymentConfirmationForm](r)
	if err != nil {
		return "", Error(err, http.StatusInternalServerError)
	}

	// Set your Stripe API key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Retrieve the payment intent directly using the stripe-go package
	pi, err := paymentintent.Get(form.PaymentId, nil)
	if err != nil {
		return "", Error(fmt.Errorf("Error retrieving payment intent: %v", err), http.StatusInternalServerError)
	}

	// Verify payment status is successful (Stripe uses "succeeded" status)
	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return "", Error(fmt.Errorf("Cannot confirm booking, payment status is %s, not succeeded", pi.Status), http.StatusInternalServerError)
	}

	// Optional: Verify the amount matches what was sent by the client
	// Note: Stripe stores amount in cents/smallest currency unit
	expectedAmount := int64(form.PaymentAmount)
	if pi.Amount != expectedAmount {
		return "", Error(fmt.Errorf("Payment amount mismatch: expected %d, got %d", expectedAmount, pi.Amount), http.StatusInternalServerError)
	}

	// Return the payment ID to be stored with the booking
	return pi.ID, nil
}

func ConfirmPayment(br *da.BookingsRepository[da.StorageDriver], api Api, r *http.Request) (string, Handler) {
	if os.Getenv("FEATURE_FLAG_PAYMENTS_PROVIDER") == "sumup" {
		return ConfirmSumupPayment(br, api, r)
	} else if os.Getenv("FEATURE_FLAG_PAYMENTS_PROVIDER") == "stripe" {
		return ConfirmStripePayment(br, r)
	} else {
		return "", Error(errors.New("No payment provider configured"), http.StatusInternalServerError)
	}
}

func BookingsConfirm(br *da.BookingsRepository[da.StorageDriver], api Api) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
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
		booking.TransactionId = "stripe-" + transactionId
		booking.PaidAt = time.Now()

		br.Update([]da.Booking{*booking})

		if booking.CustomerEmail != "test@test.com" {
			err = SendCustomerBookingConfirmationEmail(booking)
			if err != nil {
				fmt.Printf("Failed sending customer confirmation email for booking %i: %v\n", booking.Id, err)
			}
		} else {
			fmt.Println("Found test email address, skipping confirmation email")
		}

		err = SendAdminBookingNotificationEmail(booking)
		if err != nil {
			fmt.Printf("Failed sending admin booking notification email for booking %i: %v\n", booking.Id, err)
		}

		return JSON(booking)
	})
}
