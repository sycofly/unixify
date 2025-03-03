package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/home/unixify/internal/mail"
	"github.com/home/unixify/internal/models"
	"github.com/home/unixify/internal/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	UserService *service.UserService
	JWTSecret   string
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		UserService: userService,
		JWTSecret:   jwtSecret,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterUserRequest

	// Bind request
	if err := c.ShouldBindJSON(&req); err \!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Register user
	user, err := h.UserService.RegisterUser(&req)
	if err \!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token for the new user
	token, err := h.generateToken(user)
	if err \!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Prepare response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Check your email for a welcome message\!",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind request
	if err := c.ShouldBindJSON(&req); err \!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user
	user, authenticated, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err \!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	if \!authenticated || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if TOTP is required
	if user.TOTPEnabled {
		// Return response indicating TOTP is required
		c.JSON(http.StatusOK, gin.H{
			"requires_totp": true,
			"username":      user.Username,
		})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user)
	if err \!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return successful login response
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// ForgotPassword handles password reset requests
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	// Bind request
	if err := c.ShouldBindJSON(&req); err \!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send password reset email
	err := h.UserService.SendPasswordResetLink(req.Email)
	if err \!= nil {
		// Don't reveal if the email exists or not for security reasons
		c.JSON(http.StatusOK, gin.H{"message": "If your email is registered, you will receive a password reset link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent successfully. Check your email\!"})
}

// SetupTOTP sets up TOTP for a user
func (h *UserHandler) SetupTOTP(c *gin.Context) {
	// Get username from JWT claims
	claims := c.MustGet("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	// Setup TOTP
	secret, qrCode, err := h.UserService.SetupTOTP(username)
	if err \!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup TOTP"})
		return
	}

	// Return TOTP setup information
	c.JSON(http.StatusOK, gin.H{
		"secret":  secret,
		"qr_code": qrCode,
	})
}

// ActivateTOTP activates TOTP for a user
func (h *UserHandler) ActivateTOTP(c *gin.Context) {
	// Get username from JWT claims
	claims := c.MustGet("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	var req struct {
		Secret string `json:"secret" binding:"required"`
		Token  string `json:"token" binding:"required"`
	}

	// Bind request
	if err := c.ShouldBindJSON(&req); err \!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Activate TOTP
	activated, err := h.UserService.ActivateTOTP(username, req.Secret, req.Token)
	if err \!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate TOTP"})
		return
	}

	if \!activated {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TOTP activated successfully"})
}

// VerifyTOTP verifies a TOTP token
func (h *UserHandler) VerifyTOTP(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Token    string `json:"token" binding:"required"`
	}

	// Bind request
	if err := c.ShouldBindJSON(&req); err \!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	user, err := h.UserService.UserRepo.FindByUsername(req.Username)
	if err \!= nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Verify TOTP
	valid, err := h.UserService.VerifyTOTP(req.Username, req.Token)
	if err \!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify TOTP"})
		return
	}

	if \!valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP token"})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user)
	if err \!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return successful verification response
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// generateToken generates a JWT token for a user
func (h *UserHandler) generateToken(user *models.RegisteredUser) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	return token.SignedString([]byte(h.JWTSecret))
}

// AuthMiddleware is a middleware that validates JWT tokens
func (h *UserHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := authHeader[7:] // Remove "Bearer " prefix

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.JWTSecret), nil
		})

		if err \!= nil || \!token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if \!ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("claims", claims)
		c.Next()
	}
}

// RequireRole is a middleware that requires a specific role
func (h *UserHandler) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(jwt.MapClaims)
		userRole := claims["role"].(string)

		if userRole \!= role && userRole \!= "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
