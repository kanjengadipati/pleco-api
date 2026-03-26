package middleware

import (
	"go-auth-app/config"
	"go-auth-app/models"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatus(401)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.AbortWithStatus(401)
			return
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		claims := token.Claims.(jwt.MapClaims)

		// convert float64 → uint
		userIDFloat := claims["user_id"].(float64)
		userID := uint(userIDFloat)

		c.Set("user_id", userID)
		c.Set("role", claims["role"])
		c.Next()
	}
}

func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Bind request
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Parse token
	token, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Ambil claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid token claims"})
		return
	}

	// 🔍 Ambil user dari DB
	var user models.User
	config.DB.Where("id = ?", claims["user_id"]).First(&user)

	if user.ID == 0 {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	// ✅ VALIDASI refresh token harus sama dengan DB
	if user.RefreshToken != body.RefreshToken {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	// 🔐 Generate access token baru
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	accessTokenString, _ := newAccessToken.SignedString(jwtKey)

	// (Optional tapi recommended 🔥)
	// 🔄 Rotate refresh token
	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, _ := newRefreshToken.SignedString(jwtKey)

	// Simpan refresh token baru ke DB
	user.RefreshToken = refreshTokenString
	config.DB.Save(&user)

	// Response
	c.JSON(200, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	})
}
