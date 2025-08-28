package controllers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ======================= REGISTER =======================
func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// cek email unique
	var existing models.User
	if err := config.DB.First(&existing, "email = ?", input.Email).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		ID:       uuid.New(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// ======================= LOGIN =======================
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// generate access token
	accessToken, exp, err := config.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// generate refresh token
	refreshToken := uuid.NewString()
	err = config.RedisClient.Set(config.Ctx, "refresh:"+refreshToken, user.ID.String(), time.Hour*24*7).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"expires_in":    exp,
		"refresh_token": refreshToken,
	})
}

// ======================= LOGOUT =======================
func Logout(c *gin.Context) {
	// ambil access token dari header
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	// hapus refresh token dari Redis
	config.RedisClient.Del(config.Ctx, "refresh:"+body.RefreshToken)

	// blacklist access token sampai expired
	claims, err := config.ParseAccessToken(tokenString)
	if err == nil {
		exp := claims.ExpiresAt.Time.Sub(time.Now())
		config.RedisClient.Set(config.Ctx, "blacklist:"+tokenString, "1", exp)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// ======================= REFRESH =======================
func Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	// cek refresh token di Redis
	userID, err := config.RedisClient.Get(config.Ctx, "refresh:"+body.RefreshToken).Result()
	if err != nil || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token revoked or expired"})
		return
	}

	// ambil role user
	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// generate access token baru
	accessToken, exp, _ := config.GenerateAccessToken(user.ID.String(), user.Role)

	// rotate refresh token (hapus lama, bikin baru)
	config.RedisClient.Del(config.Ctx, "refresh:"+body.RefreshToken)
	newRefresh := uuid.NewString()
	config.RedisClient.Set(config.Ctx, "refresh:"+newRefresh, user.ID.String(), time.Hour*24*7)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"expires_in":    exp,
		"refresh_token": newRefresh,
	})
}
