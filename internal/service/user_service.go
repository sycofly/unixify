package service

import (
	"errors"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/home/unixify/internal/mail"
	"github.com/home/unixify/internal/models"
	"github.com/home/unixify/internal/repository"
)

// UserService handles business logic for users
type UserService struct {
	UserRepo *repository.UserRepository
	Mailer   *mail.Mailer
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		Mailer:   mail.NewMailer(),
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(req *models.RegisterUserRequest) (*models.RegisteredUser, error) {
	// Check if username already exists
	existingUser, err := s.UserRepo.FindByUsername(req.Username)
	if err \!= nil {
		return nil, err
	}
	if existingUser \!= nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.UserRepo.FindByEmail(req.Email)
	if err \!= nil {
		return nil, err
	}
	if existingUser \!= nil {
		return nil, errors.New("email already registered")
	}

	// Create new user
	user := &models.RegisteredUser{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Department: req.Department,
		Role:      "user", // Default role
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Special case for admin user
	if req.Email == "peter.frederikson@gmail.com" {
		user.Role = "admin"
	}

	// Create user in database (password will be encrypted by repository)
	err = s.UserRepo.CreateUser(user, req.Password)
	if err \!= nil {
		return nil, err
	}

	// Send welcome email
	go func() {
		// Send a welcome email to the user
		emailErr := s.Mailer.SendWelcomeEmail(user.Email, user.Username)
		if emailErr \!= nil {
			// Just log the error, don't interrupt registration
			// logger.Error("Failed to send welcome email", zap.Error(emailErr))
		}
	}()

	return user, nil
}

// AuthenticateUser authenticates a user
func (s *UserService) AuthenticateUser(username, password string) (*models.RegisteredUser, bool, error) {
	// Find user by username
	user, err := s.UserRepo.FindByUsername(username)
	if err \!= nil {
		return nil, false, err
	}
	if user == nil {
		return nil, false, nil // User not found
	}

	// Verify password
	passwordMatches, err := s.UserRepo.VerifyPassword(username, password)
	if err \!= nil {
		return nil, false, err
	}
	if \!passwordMatches {
		return nil, false, nil // Password doesn't match
	}

	// Update last login timestamp
	err = s.UserRepo.UpdateLastLogin(username)
	if err \!= nil {
		// Non-critical error, just log it
		// logger.Error("Failed to update last login time", zap.Error(err))
	}

	return user, true, nil
}

// SendPasswordResetLink sends a password reset link to the user's email
func (s *UserService) SendPasswordResetLink(email string) error {
	// Find user by email
	user, err := s.UserRepo.FindByEmail(email)
	if err \!= nil {
		return err
	}
	if user == nil {
		return errors.New("email not found")
	}

	// Create a reset token (in a real app, store this with expiration)
	token := "reset_" + time.Now().Format("20060102150405")

	// Send password reset email
	return s.Mailer.SendPasswordResetEmail(user.Email, user.Username, token)
}

// SetupTOTP sets up TOTP for a user
func (s *UserService) SetupTOTP(username string) (string, string, error) {
	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Unixify",
		AccountName: username,
	})
	if err \!= nil {
		return "", "", err
	}

	// Get the secret in base32 format
	secret := key.Secret()
	
	// Get the QR code URL
	qrCode := key.URL()

	return secret, qrCode, nil
}

// ActivateTOTP activates TOTP for a user
func (s *UserService) ActivateTOTP(username, secret, token string) (bool, error) {
	// Verify the token
	valid := totp.Validate(token, secret)
	if \!valid {
		return false, nil
	}

	// Enable TOTP for the user
	err := s.UserRepo.EnableTOTP(username, secret)
	if err \!= nil {
		return false, err
	}

	return true, nil
}

// VerifyTOTP verifies a TOTP token
func (s *UserService) VerifyTOTP(username, token string) (bool, error) {
	// Find user
	user, err := s.UserRepo.FindByUsername(username)
	if err \!= nil || user == nil {
		return false, err
	}

	// Verify the token
	return totp.Validate(token, user.TOTPSecret), nil
}
