package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/Michaelpalacce/gobi/pkg/digest"
	"github.com/gin-gonic/gin"
)

func BasicAuth(userService *services.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if the Authorization header is present
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if the Authorization header starts with "Basic "
		if !strings.HasPrefix(authHeader, "Basic ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		// Decode the base64-encoded username and password
		decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid base64 encoding"})
			c.Abort()
			return
		}

		// Split the decoded string into username and password
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			c.Abort()
			return
		}

		user, err := userService.GetUserByName(credentials[0])

		// Check if the user exists
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching user by that name"})
			c.Abort()
			return
		}

		// Check if the provided credentials match the valid credentials
		if credentials[0] != user.Username || digest.SHA256(credentials[1]) != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		c.Set("user", user)

		// Continue with the next middleware or route handler
		c.Next()
	}
}
