package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/dto"
	"github.com/qullDev/BookMyField/internal/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
)

// CreateCheckoutSession godoc
// @Summary Create a checkout session
// @Description Create a new checkout session for a booking payment.
// @Tags payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dto.CreateCheckoutSessionRequest true "Booking ID for payment"
// @Success 200 {object} dto.CreateCheckoutSessionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /payments/create-checkout-session [post]
func CreateCheckoutSession(c *gin.Context) {
	var req dto.CreateCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if booking exists and belongs to the user
	var booking models.Booking
	if err := config.DB.Preload("Field").
		Where("id = ? AND user_id = ?", req.BookingID, userID).
		First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found or not authorized"})
		return
	}

	// Check if payment already exists for this booking
	var existingPayment models.Payment
	if err := config.DB.Where("booking_id = ? AND status IN (?)", req.BookingID, []string{"pending", "succeeded"}).First(&existingPayment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment already exists for this booking"})
		return
	}

	// Check booking status
	if booking.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is not in pending status"})
		return
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("payment"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("idr"), // Changed to IDR
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(booking.Field.Name),
					},
					UnitAmount: stripe.Int64(int64(booking.Field.Price)), // IDR doesn't need * 100
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://localhost:3000/success"),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
		Metadata: map[string]string{
			"booking_id": booking.ID.String(),
		},
	}

	s, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save Payment record with session ID instead of PaymentIntent ID
	payment := models.Payment{
		ID:          uuid.New(),
		BookingID:   booking.ID,
		Amount:      booking.Field.Price,
		Currency:    "idr", // Changed to IDR
		Status:      "pending",
		StripeRefID: s.ID, // Use session ID for webhook matching
	}
	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
		return
	}

	c.JSON(http.StatusOK, dto.CreateCheckoutSessionResponse{
		SessionID:  s.ID,
		SessionURL: s.URL,
	})
}

// StripeWebhook godoc
// @Summary Stripe webhook
// @Description Handle Stripe webhook events to update payment and booking status.
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /payments/stripe-webhook [post]
func StripeWebhook(c *gin.Context) {
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}
	sigHeader := c.GetHeader("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook signature"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
			// Update payment status
			result := config.DB.Model(&models.Payment{}).
				Where("stripe_ref_id = ?", session.ID).
				Update("status", "succeeded")

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
				return
			}

			// Update booking status to confirmed
			if bookingID, exists := session.Metadata["booking_id"]; exists {
				config.DB.Model(&models.Booking{}).
					Where("id = ?", bookingID).
					Update("status", "confirmed")
			}
		}

	case "checkout.session.expired", "checkout.session.async_payment_failed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
			config.DB.Model(&models.Payment{}).
				Where("stripe_ref_id = ?", session.ID).
				Update("status", "failed")
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
