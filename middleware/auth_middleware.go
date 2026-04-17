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

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		userIDValue, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		roleValue, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// ✅ inject ke context
		userID := uint(userIDValue)
		c.Set("user_id", userID)
		c.Set("role", roleValue)

		c.Next()
	}
}
