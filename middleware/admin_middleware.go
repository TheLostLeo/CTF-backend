package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the authenticated user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated (should be called after AuthMiddleware)
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Check if user is admin
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		// User is authenticated and is admin, continue
		c.Next()
	}
}
