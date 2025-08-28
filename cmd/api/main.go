package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
	"github.com/qullDev/BookMyField/internal/routes"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	config.ConnectDatabse()

	err := config.DB.AutoMigrate(&models.User{}, &models.Field{}, &models.Booking{})
	if err != nil {
		log.Fatal("Error migrating database:", err.Error())
		return
	}

	api := r.Group("/api/v1")
	{
		routes.AuthRoutes(api)
		routes.BookingsRoutes(api)
		routes.FieldRoutes(api)
	}

	r.Run(":8080")
}
