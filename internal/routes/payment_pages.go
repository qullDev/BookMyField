package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PaymentPagesRoutes(r *gin.Engine) {
	// Success page after payment
	r.GET("/success", func(c *gin.Context) {
		sessionID := c.Query("session_id")

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"message":    "🎉 Payment completed successfully!",
			"session_id": sessionID,
			"next_steps": "Your booking has been confirmed. Check your bookings in the app.",
		})
	})

	// Cancel page when payment is cancelled
	r.GET("/cancel", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "cancelled",
			"message": "❌ Payment was cancelled. You can try again anytime.",
		})
	})
}
