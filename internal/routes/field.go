package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
	"github.com/qullDev/BookMyField/internal/middlewares"
)

func FieldRoutes(rg *gin.RouterGroup) {
	field := rg.Group("/fields", middlewares.AuthMiddleware())
	{
		field.GET("/", controllers.GetFields)
		field.GET("/:id", controllers.GetFieldByID)
	}

	admin := field.Group("/admin", middlewares.AdminOnly())
	{
		admin.POST("/", controllers.CreateField)
		admin.DELETE("/:id", controllers.DeleteField)
		admin.PUT("/:id", controllers.UpdateField)
	}
}
