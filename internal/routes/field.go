package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
	"github.com/qullDev/BookMyField/internal/middlewares"
)

func FieldRoutes(rg *gin.RouterGroup) {
	field := rg.Group("/fields")
	field.Use(middlewares.AuthMiddleware()) // semua butuh login
	{
		field.GET("/", controllers.GetFields)         // semua user bisa lihat
		field.POST("/", controllers.CreateField)      // admin only
		field.PUT("/:id", controllers.UpdateField)    // admin only
		field.DELETE("/:id", controllers.DeleteField) // admin only
	}
}
