package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	stripeRefund "github.com/stripe/stripe-go/v76/refund"
)

type CreateBookingInput struct {
	FieldID   string    `json:"field_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Notes     string    `json:"notes,omitempty"`
}

// GetBookings godoc
// @Summary Get all bookings (Admin only)
// @Description Get a list of all bookings. Requires admin privileges.
// @Tags bookings
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Booking
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /bookings [get]
func GetBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := config.DB.
		Preload("User").
		Preload("Field").
		Preload("Payments").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking for a field.
// @Tags bookings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dto.CreateBookingRequest true "Booking data"
// @Success 201 {object} models.Booking
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /bookings [post]
func CreateBooking(c *gin.Context) {
	var input CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id in token"})
		return
	}

	fieldID, err := uuid.Parse(input.FieldID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field_id format"})
		return
	}

	// Validate field exists
	var field models.Field
	if err := config.DB.First(&field, "id = ?", fieldID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Field not found"})
		return
	}

	// Validate time
	if input.StartTime.After(input.EndTime) || input.StartTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking time"})
		return
	}

	// Check for conflicting bookings
	var conflictCount int64
	config.DB.Model(&models.Booking{}).
		Where("field_id = ? AND status != ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
			fieldID, "cancelled", input.StartTime, input.StartTime, input.EndTime, input.EndTime).
		Count(&conflictCount)

	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Field is already booked for this time slot"})
		return
	}

	booking := models.Booking{
		ID:        uuid.New(),
		UserID:    uid,
		FieldID:   fieldID,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Status:    "pending",
		Notes:     input.Notes,
	}

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// Preload relasi untuk response
	if err := config.DB.
		Preload("User").
		Preload("Field").
		Preload("Payments").
		First(&booking, "id = ?", booking.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve booking details"})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// GetMyBookings godoc
// @Summary Get my bookings
// @Description Get a list of bookings for the currently authenticated user.
// @Tags bookings
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Booking
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /bookings/me [get]
func GetMyBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var bookings []models.Booking
	if err := config.DB.
		Preload("Field").
		Preload("Payments").
		Where("user_id = ?", userID).
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking by its ID. If payment exists, it will be refunded.
// @Tags bookings
// @Security BearerAuth
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} dto.CancelBookingResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /bookings/{id}/cancel [delete]
func CancelBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	bookingID := c.Param("id")

	var booking models.Booking
	if err := config.DB.Preload("Payments").
		First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Check if booking is already cancelled
	if booking.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking already cancelled"})
		return
	}

	// Start database transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// If no payments, simply cancel the booking
	if len(booking.Payments) == 0 {
		booking.Status = "cancelled"
		if err := tx.Save(&booking).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
		return
	}

	// Find the latest successful payment
	var latestPayment *models.Payment
	for i := len(booking.Payments) - 1; i >= 0; i-- {
		if booking.Payments[i].Status == "succeeded" {
			latestPayment = &booking.Payments[i]
			break
		}
	}

	if latestPayment == nil {
		// No successful payment, just cancel
		booking.Status = "cancelled"
		if err := tx.Save(&booking).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
		return
	}

	// Refund via Stripe - need to get PaymentIntent ID from the session
	var sessionDetails *stripe.CheckoutSession
	sessionDetails, err := session.Get(latestPayment.StripeRefID, nil)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment session"})
		return
	}

	if sessionDetails.PaymentIntent == nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment intent not found"})
		return
	}

	// Create refund
	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(sessionDetails.PaymentIntent.ID),
	}
	ref, err := stripeRefund.New(refundParams)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process refund: " + err.Error()})
		return
	}

	// Update database
	booking.Status = "cancelled"
	latestPayment.Status = "refunded"

	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}
	if err := tx.Save(latestPayment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":       "Booking cancelled and payment refunded successfully",
		"refund_id":     ref.ID,
		"refund_status": ref.Status,
	})
}
