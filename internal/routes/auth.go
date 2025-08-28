package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/qullDev/BookMyField/internal/controllers"
)

func AuthRoutes(api *gin.RouterGroup) {
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
	api.POST("/logout", controllers.Logout)
	api.POST("/refresh", controllers.Refresh)
}
