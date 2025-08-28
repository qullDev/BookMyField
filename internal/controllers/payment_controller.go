package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

type CreateCheckoutSessionInput struct {
	BookingID string `json:"booking_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required"` // dalam sen (IDR pakai x100)
}

func CreateCheckoutSession(c *gin.Context) {
	var input CreateCheckoutSessionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"), // kalau mau IDR ganti ke "idr"
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Booking Lapangan"),
					},
					UnitAmount: stripe.Int64(input.Amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(os.Getenv("FRONTEND_URL") + "/success"),
		CancelURL:  stripe.String(os.Getenv("FRONTEND_URL") + "/cancel"),
	}

	s, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"checkout_url": s.URL})
}
