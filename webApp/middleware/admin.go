package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	return func(c *gin.Context) {
		password := c.GetHeader("X-Admin-Password")
		if password != adminPassword {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
