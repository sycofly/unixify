package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/home/unixify/internal/auth"
	"github.com/home/unixify/internal/models"
	"github.com/home/unixify/internal/repository"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	db         *gorm.DB
	authService *auth.Service
	repo       *repository.Repository
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB, authService *auth.Service, repo *repository.Repository) *AuthHandler {
	return &AuthHandler{
		db:         db,
		authService: authService,
		repo:       repo,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if result := h.db.Where("username = ?", input.Username).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists
	if result := h.db.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := h.authService.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Username:    input.Username,
		Password:    hashedPassword,
		Email:       input.Email,
		Role:        "user", // Default role for new users
		TOTPEnabled: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if result := h.db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": models.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Role:        user.Role,
			TOTPEnabled: user.TOTPEnabled,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username
	var user models.User
	if result := h.db.Where("username = ?", input.Username).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !h.authService.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// If TOTP is enabled, don't generate a token yet
	if user.TOTPEnabled {
		c.JSON(http.StatusOK, gin.H{
			"message": "TOTP verification required",
			"requires_totp": true,
			"user": models.UserResponse{
				ID:          user.ID,
				Username:    user.Username,
				Email:       user.Email,
				Role:        user.Role,
				TOTPEnabled: user.TOTPEnabled,
			},
		})
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update last login time
	user.LastLogin = time.Now()
	h.db.Save(&user)

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Role:        user.Role,
			TOTPEnabled: user.TOTPEnabled,
		},
	})
}

// VerifyTOTP verifies a TOTP code after initial login
func (h *AuthHandler) VerifyTOTP(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Token    string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username
	var user models.User
	if result := h.db.Where("username = ?", input.Username).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Verify TOTP code
	if !h.authService.VerifyTOTP(user.TOTPSecret, input.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP code"})
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update last login time
	user.LastLogin = time.Now()
	h.db.Save(&user)

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Role:        user.Role,
			TOTPEnabled: user.TOTPEnabled,
		},
	})
}

// SetupTOTP generates a new TOTP secret for a user
func (h *AuthHandler) SetupTOTP(c *gin.Context) {
	// Get user from context
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userResponse, ok := userObj.(*models.UserResponse)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user information"})
		return
	}

	// Find the full user record
	var user models.User
	if result := h.db.Where("id = ?", userResponse.ID).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate TOTP secret and QR code
	totpResponse, err := h.authService.GenerateTOTPSecret(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
		return
	}

	// Store secret in user record but don't enable TOTP yet
	// It will be enabled after the user verifies a valid code
	user.TOTPSecret = totpResponse.Secret
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save TOTP secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "TOTP setup initiated",
		"secret": totpResponse.Secret,
		"qr_code": totpResponse.QRCode,
	})
}

// ActivateTOTP validates a TOTP code and enables TOTP for the user
func (h *AuthHandler) ActivateTOTP(c *gin.Context) {
	// Get user from context
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userResponse, ok := userObj.(*models.UserResponse)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user information"})
		return
	}

	var input models.TOTPVerifyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the full user record
	var user models.User
	if result := h.db.Where("id = ?", userResponse.ID).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify the token against the stored secret
	if !h.authService.VerifyTOTP(user.TOTPSecret, input.Token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TOTP code"})
		return
	}

	// Enable TOTP for the user
	user.TOTPEnabled = true
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable TOTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "TOTP enabled successfully",
	})
}

// DisableTOTP disables TOTP for a user
func (h *AuthHandler) DisableTOTP(c *gin.Context) {
	// Get user from context
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userResponse, ok := userObj.(*models.UserResponse)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user information"})
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the full user record
	var user models.User
	if result := h.db.Where("id = ?", userResponse.ID).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify password for security
	if !h.authService.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Disable TOTP
	user.TOTPEnabled = false
	user.TOTPSecret = ""
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable TOTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "TOTP disabled successfully",
	})
}

// GetProfile returns the user's profile information
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user from context
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userResponse, ok := userObj.(*models.UserResponse)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userResponse,
	})
}

// UpdatePassword updates a user's password
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	// Get user from context
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userResponse, ok := userObj.(*models.UserResponse)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user information"})
		return
	}

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the full user record
	var user models.User
	if result := h.db.Where("id = ?", userResponse.ID).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	if !h.authService.CheckPassword(input.CurrentPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash the new password
	hashedPassword, err := h.authService.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}