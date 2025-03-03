package repository

import (
	"errors"
	"time"

	"github.com/home/unixify/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles database operations for RegisteredUser
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser creates a new user with encrypted password
func (r *UserRepository) CreateUser(user *models.RegisteredUser, password string) error {
	// Convert plain text password to bytea
	user.PasswordHash = []byte(password)
	
	// Database trigger will handle the encryption

	// Create the user
	result := r.DB.Create(user)
	return result.Error
}

// FindByUsername finds a user by their username
func (r *UserRepository) FindByUsername(username string) (*models.RegisteredUser, error) {
	var user models.RegisteredUser
	result := r.DB.Where("username = ?", username).First(&user)
	if result.Error \!= nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, result.Error // Database error
	}
	return &user, nil
}

// FindByEmail finds a user by their email
func (r *UserRepository) FindByEmail(email string) (*models.RegisteredUser, error) {
	var user models.RegisteredUser
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error \!= nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, result.Error // Database error
	}
	return &user, nil
}

// VerifyPassword verifies the user's password directly with the database
func (r *UserRepository) VerifyPassword(username, password string) (bool, error) {
	// Get the stored hash for this user
	var storedHash []byte
	result := r.DB.Model(&models.RegisteredUser{}).
		Where("username = ?", username).
		Select("password_hash").
		Take(&storedHash)
	
	if result.Error \!= nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil // User not found
		}
		return false, result.Error // Database error
	}

	// For PostgreSQL with pgcrypto, we can verify directly with SQL
	var passwordMatches bool
	result = r.DB.Raw("SELECT (password_hash = crypt(?, password_hash)) FROM registered_users WHERE username = ?", 
		password, username).Scan(&passwordMatches)
	
	if result.Error \!= nil {
		return false, result.Error
	}
	
	return passwordMatches, nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(username string) error {
	result := r.DB.Model(&models.RegisteredUser{}).
		Where("username = ?", username).
		Update("last_login", time.Now())
	return result.Error
}

// EnableTOTP enables TOTP for a user
func (r *UserRepository) EnableTOTP(username, secret string) error {
	result := r.DB.Model(&models.RegisteredUser{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"totp_enabled": true,
			"totp_secret": secret,
		})
	return result.Error
}

// DisableTOTP disables TOTP for a user
func (r *UserRepository) DisableTOTP(username string) error {
	result := r.DB.Model(&models.RegisteredUser{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"totp_enabled": false,
			"totp_secret": "",
		})
	return result.Error
}
