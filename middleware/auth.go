package middleware

import (
	"gin/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bearer <token>
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized header request"})
			c.Abort()
			return
		}

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access token"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(authHeader)
		log.Printf("Claims: %v\n", claims)
		if err != nil {
			log.Printf("Token validation error: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
