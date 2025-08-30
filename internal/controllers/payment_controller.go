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
					Currency: stripe.String("usd"), // Use USD for Stripe compatibility
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(booking.Field.Name + " - Field Booking"),
					},
					UnitAmount: stripe.Int64(100), // Fixed $1.00 USD for testing (minimum amount)
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("https://bookmyfield-production.up.railway.app/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://bookmyfield-production.up.railway.app/cancel"),
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

// GetPayments godoc
// @Summary Get all payments for admin
// @Description Get all payments with booking and user details (admin only).
// @Tags payments
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Payment
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /payments [get]
func GetPayments(c *gin.Context) {
	var payments []models.Payment
	if err := config.DB.Preload("Booking.User").Preload("Booking.Field").Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetMyPayments godoc
// @Summary Get user's payments
// @Description Get all payments for the authenticated user.
// @Tags payments
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Payment
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /payments/me [get]
func GetMyPayments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payments []models.Payment
	if err := config.DB.Preload("Booking.User").Preload("Booking.Field").
		Joins("JOIN bookings ON payments.booking_id = bookings.id").
		Where("bookings.user_id = ?", userID).
		Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentByID godoc
// @Summary Get payment by ID
// @Description Get payment details by payment ID.
// @Tags payments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} models.Payment
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /payments/{id} [get]
func GetPaymentByID(c *gin.Context) {
	paymentID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payment models.Payment
	query := config.DB.Preload("Booking.User").Preload("Booking.Field").
		Joins("JOIN bookings ON payments.booking_id = bookings.id").
		Where("payments.id = ?", paymentID)

	// Check user role
	role, roleExists := c.Get("role")
	if !roleExists || role != "admin" {
		// If not admin, only allow access to own payments
		query = query.Where("bookings.user_id = ?", userID)
	}

	if err := query.First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}
