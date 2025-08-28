package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
)

func GetFields(c *gin.Context) {
	var fields []models.Field

	if err := config.DB.Find(&fields).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve fields"})
		return
	}

	c.JSON(http.StatusOK, fields)

}

func CreateField(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var input struct {
		Name     string  `json:"name" binding:"required"`
		Location string  `json:"location" binding:"required"`
		Price    float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	field := models.Field{
		ID:       uuid.New(),
		Name:     input.Name,
		Location: input.Location,
		Price:    input.Price,
	}

	if err := config.DB.Create(&field).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, field)
}

func UpdateField(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")

	var field models.Field
	if err := config.DB.First(&field, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Field not found"})
	}

	var input struct {
		Name     string  `json:"name"`
		Location string  `json:"location"`
		Price    float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	field.Name = input.Name
	field.Location = input.Location
	field.Price = input.Price

	if err := config.DB.Save(&field).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, field)

}

func DeleteField(c *gin.Context) {

	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	if err := config.DB.Delete(&models.Field{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Field deleted successfully"})
}
