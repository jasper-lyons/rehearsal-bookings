package handlers

import (
	"net/http"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	da "rehearsal-bookings/pkg/data_access"
)

type CreatePaymentIntentForm struct {
	BookingId int `json:"booking_id"`
	CheckoutReference string `json:"checkout_reference"`
	CustomerEmail string `json:"customer_email"`
}

func StripePaymentIntentCreate(br *da.BookingsRepository[da.StorageDriver]) Handler {
	return Handler(func (w http.ResponseWriter, r *http.Request) Handler {
		form, err := ExtractForm[CreatePaymentIntentForm](r)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		booking, err := br.Find(form.BookingId)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		params := &stripe.PaymentIntentParams {
			Amount: stripe.Int64(int64(booking.Price * 100)),
			Currency: stripe.String(string(stripe.CurrencyGBP)),
		}
		result, err := paymentintent.New(params)
		if err != nil {
			return Error(err, http.StatusInternalServerError)
		}

		return JSON(result)
	})
}
