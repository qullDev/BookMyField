package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
	"github.com/qullDev/BookMyField/internal/middlewares"
)

func PaymentRoutes(api *gin.RouterGroup) {
	// Checkout endpoint requires authentication
	api.POST("/checkout", middlewares.AuthMiddleware(), controllers.CreateCheckoutSession)

	// Webhook endpoint tidak perlu authentication (dipanggil oleh Stripe)
	api.POST("/webhook", controllers.StripeWebhook)
}
