package controllers

import (
	"go-auth-app/dto"
	"go-auth-app/models"
	"go-auth-app/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	AuthService *services.AuthService
}

func (a *AuthController) Register(c *gin.Context) {
	var input dto.RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err := a.AuthService.Register(&user, input.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered"})
}

func (a *AuthController) Login(c *gin.Context) {
	var input dto.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	tokens, err := a.AuthService.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tokens)
}

func (a *AuthController) Logout(c *gin.Context) {

	// ✅ ambil user dari middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// ✅ ambil device id dari header
	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device id required"})
		return
	}

	// ✅ call service
	err := a.AuthService.Logout(userID.(uint), deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
	})
}

func (a *AuthController) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	log.Println("refresh token", body.RefreshToken)

	tokens, err := a.AuthService.RefreshToken(body.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tokens)
}

func (a *AuthController) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// ambil user dari repository
	user, err := a.AuthService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// response (hindari kirim password!)
	c.JSON(200, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
