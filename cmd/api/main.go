package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
	"github.com/qullDev/BookMyField/internal/routes"
	"github.com/qullDev/BookMyField/internal/seed"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	config.ConnectDatabse()
	config.InitRedis()
	config.InitStripe()

	err := config.DB.AutoMigrate(&models.User{}, &models.Field{}, &models.Booking{})
	if err != nil {
		log.Fatal("Error migrating database:", err.Error())
		return
	}

	// seed data
	seed.SeedAdminUser()
	seed.SeedFields()
	seed.SeedRegularUser()

	// Route
	api_v1 := r.Group("/api/v1")
	{
		routes.AuthRoutes(api_v1)
		routes.BookingsRoutes(api_v1)
		routes.FieldRoutes(api_v1)
		routes.PaymentRoutes(api_v1)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error running server:", err.Error())
	}
}
