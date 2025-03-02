package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/home/unixify/internal/models"
)

// AuthMiddleware creates a middleware for authenticating requests
func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString, err := s.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Verify the token
		token, err := s.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user information from the token
		user, err := s.GetUserFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract user information"})
			c.Abort()
			return
		}

		// Store user info in the context for later use in handlers
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("role", user.Role)

		c.Next()
	}
}

// RoleMiddleware checks if the user has the required role
func (s *Service) RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userObj, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}

		user, ok := userObj.(*models.UserResponse)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user information"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasPermission := false
		for _, role := range roles {
			if user.Role == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}