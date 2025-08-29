package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
)

func AuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")

	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)
	auth.POST("/logout", controllers.Logout)
	auth.POST("/refresh", controllers.Refresh)
}
