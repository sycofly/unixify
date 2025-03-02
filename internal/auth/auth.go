package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/home/unixify/internal/config"
	"github.com/home/unixify/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Service contains the authentication-related logic
type Service struct {
	config config.Config
}

// NewService creates a new authentication service
func NewService(cfg config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// HashPassword creates a bcrypt hash of the password
func (s *Service) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword checks if the provided password matches the hash
func (s *Service) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT token for the given user
func (s *Service) GenerateToken(user *models.User) (string, error) {
	// Set token expiration to 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := jwt.MapClaims{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"role":      user.Role,
		"exp":       expirationTime.Unix(),
		"issued_at": time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString([]byte(s.config.Server.Secret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// VerifyToken validates a JWT token
func (s *Service) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Server.Secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return token, nil
}

// ExtractTokenFromHeader extracts the JWT token from the Authorization header
func (s *Service) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}
	
	return parts[1], nil
}

// GetUserFromToken extracts user information from a JWT token
func (s *Service) GetUserFromToken(token *jwt.Token) (*models.UserResponse, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	
	userID, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}
	
	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("invalid username in token")
	}
	
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email in token")
	}
	
	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in token")
	}
	
	user := &models.UserResponse{
		ID:       uint(userID),
		Username: username,
		Email:    email,
		Role:     role,
	}
	
	return user, nil
}