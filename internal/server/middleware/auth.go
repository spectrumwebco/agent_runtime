package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/server/auth"
)

// AuthMiddleware is a middleware that authenticates a request
func AuthMiddleware(authHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Check if the header is in the correct format
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must be in the format 'Bearer {token}'",
			})
			return
		}

		// Extract the token
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			return
		}

		// Check if the token is valid
		if token != authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// Set the user in the context
		c.Set("user", "user")

		c.Next()
	}
}

// AdminMiddleware is a middleware that checks if the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetString("user")
		if user != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			return
		}

		c.Next()
	}
}

// AuthorizeMiddleware is a middleware that authorizes a request
func AuthorizeMiddleware(enforcer *auth.Enforcer) gin.HandlerFunc {
	return enforcer.Authorize()
}
