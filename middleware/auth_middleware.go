package middleware

import (
	"go-auth-app/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// ✅ inject ke context
		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)
		c.Set("role", claims["role"])

		c.Next()
	}
}
