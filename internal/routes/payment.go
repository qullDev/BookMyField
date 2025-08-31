package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
	"github.com/qullDev/BookMyField/internal/middlewares"
)

func PaymentRoutes(api *gin.RouterGroup) {
	// Webhook endpoint (no authentication required - Di panggil di stripe)
	// Harus di luar tanpa middalware auth
	api.POST("/payments/stripe-webhook", controllers.StripeWebhook)

	// Test webhook endpoint untuk development (no signature validation)
	api.POST("/payments/stripe-webhook-test", controllers.StripeWebhookTest)

	payment := api.Group("/payments")
	{
		// Create checkout session (requires authentication)
		payment.POST("/create-checkout-session", middlewares.AuthMiddleware(), controllers.CreateCheckoutSession)

		// Get all payments (admin only)
		payment.GET("/", middlewares.AuthMiddleware(), middlewares.AdminOnly(), controllers.GetPayments)

		// Get user's payments (requires authentication)
		payment.GET("/me", middlewares.AuthMiddleware(), controllers.GetMyPayments)

		// Get payment by ID (requires authentication, user can only access own payments)
		payment.GET("/:id", middlewares.AuthMiddleware(), controllers.GetPaymentByID)
	}
}
