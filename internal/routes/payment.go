package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
	"github.com/qullDev/BookMyField/internal/middlewares"
)

func PaymentRoutes(api *gin.RouterGroup) {
	payment := api.Group("/payments")
	{
		// Checkout endpoint requires authentication
		payment.POST("/create-checkout-session", middlewares.AuthMiddleware(), controllers.CreateCheckoutSession)

		// Webhook endpoint tidak perlu authentication (dipanggil oleh Stripe)
		payment.POST("/stripe-webhook", controllers.StripeWebhook)
	}
}
