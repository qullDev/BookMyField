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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required,min=2"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize email to lowercase
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	// Check if email already exists
	var existing models.User
	if err := config.DB.First(&existing, "email = ?", input.Email).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Name:     strings.TrimSpace(input.Name),
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary Log in a user
// @Description Log in a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginRequest true "User login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize email
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	var user models.User
	if err := config.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate access token
	accessToken, exp, err := config.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate refresh token
	refreshToken := uuid.NewString()
	if config.RedisClient != nil {
		err = config.RedisClient.Set(config.Ctx, "refresh:"+refreshToken, user.ID.String(), time.Hour*24*7).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"expires_in":    exp,
		"refresh_token": refreshToken,
	})
}

// Logout godoc
// @Summary Log out a user
// @Description Log out a user by blacklisting access token and deleting refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.LogoutRequest true "Refresh token to invalidate"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/logout [post]
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
	if config.RedisClient != nil {
		config.RedisClient.Del(config.Ctx, "refresh:"+body.RefreshToken)
	}

	// blacklist access token sampai expired
	claims, err := config.ParseAccessToken(tokenString)
	if err == nil && config.RedisClient != nil {
		exp := claims.ExpiresAt.Time.Sub(time.Now())
		config.RedisClient.Set(config.Ctx, "blacklist:"+tokenString, "1", exp)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	// Check refresh token in Redis
	if config.RedisClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis not available"})
		return
	}

	userID, err := config.RedisClient.Get(config.Ctx, "refresh:"+body.RefreshToken).Result()
	if err != nil || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token revoked or expired"})
		return
	}

	// Get user role and validate user exists
	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		// If user not found, delete the refresh token for security
		config.RedisClient.Del(config.Ctx, "refresh:"+body.RefreshToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new access token
	accessToken, exp, err := config.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Rotate refresh token (delete old, create new)
	config.RedisClient.Del(config.Ctx, "refresh:"+body.RefreshToken)
	newRefresh := uuid.NewString()
	if err := config.RedisClient.Set(config.Ctx, "refresh:"+newRefresh, user.ID.String(), time.Hour*24*7).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"expires_in":    exp,
		"refresh_token": newRefresh,
	})
}
