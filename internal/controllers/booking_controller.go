package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
	"github.com/stripe/stripe-go/v76"
	stripeRefund "github.com/stripe/stripe-go/v76/refund"
)

type CreateBookingInput struct {
	FieldID   string    `json:"field_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

// GET /api/v1/bookings
func GetBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := config.DB.Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// POST /api/v1/bookings
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

	booking := models.Booking{
		ID:        uuid.New(),
		UserID:    uid,
		FieldID:   uuid.MustParse(input.FieldID),
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Status:    "pending",
	}

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// controllers/booking.go
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

func CancelBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	bookingID := c.Param("id")

	var booking models.Booking
	if err := config.DB.Preload("Payments").
		First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// kalau belum ada payment
	if len(booking.Payments) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No payment found for this booking"})
		return
	}

	// ambil payment terakhir
	payment := booking.Payments[len(booking.Payments)-1]

	if payment.Status != "succeeded" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking cannot be cancelled"})
		return
	}

	// Refund di Stripe (pakai PaymentIntent ID dari StripeRefID)
	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(payment.StripeRefID),
	}
	ref, err := stripeRefund.New(refundParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update DB
	booking.Status = "cancelled"
	payment.Status = "refunded"

	if err := config.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking cancelled & refunded",
		"refund":  ref,
	})
}
