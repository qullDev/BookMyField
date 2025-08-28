package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
)

func PaymentRoutes(api *gin.RouterGroup) {
	api.POST("/checkout", controllers.CreateCheckoutSession)
	api.POST("/webhook", controllers.StripeWebhook)
}
