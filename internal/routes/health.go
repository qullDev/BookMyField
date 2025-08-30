package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthRoute(api *gin.RouterGroup) {
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

}
