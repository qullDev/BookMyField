package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
	"github.com/qullDev/BookMyField/internal/middlewares"
)

func BookingsRoutes(api *gin.RouterGroup) {

	booking := api.Group("/bookings")
	booking.Use(middlewares.AuthMiddleware())
	{
		booking.GET("/", middlewares.AdminOnly(),controllers.GetBookings)     // semua booking (admin only)
		booking.GET("/me", controllers.GetMyBookings) // hanya booking user sendiri
		booking.POST("/", controllers.CreateBooking)
		booking.DELETE("/:id", controllers.CancelBooking)
		booking.DELETE("/:id/cancel", controllers.CancelBooking)
	}
}
