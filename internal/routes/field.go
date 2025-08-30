package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
	"github.com/qullDev/BookMyField/internal/middlewares"
)

func FieldRoutes(rg *gin.RouterGroup) {
	// Public endpoints (no authentication required)
	fields := rg.Group("/fields")
	{
		fields.GET("/", controllers.GetFields)
		fields.GET("/:id", controllers.GetFieldByID)
	}

	// Admin-only endpoints (authentication + admin role required)
	admin := rg.Group("/fields/admin", middlewares.AuthMiddleware(), middlewares.AdminOnly())
	{
		admin.POST("/", controllers.CreateField)
		admin.DELETE("/:id", controllers.DeleteField)
		admin.PUT("/:id", controllers.UpdateField)
	}
}
