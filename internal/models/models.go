package models

import (
	"time"
)

// Account represents a UNIX account (user)
type Account struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
	Username  string    `json:"username" gorm:"unique"`
	UID       int       `json:"uid" gorm:"unique"`
	Type      string    `json:"type" gorm:"index"` // people, system, database, service
	Active    bool      `json:"active" gorm:"default:true"`
}

// Group represents a UNIX group
type Group struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
	Groupname string    `json:"groupname" gorm:"unique"`
	GID       int       `json:"gid" gorm:"unique"`
	Type      string    `json:"type" gorm:"index"` // people, system, database, service
	Active    bool      `json:"active" gorm:"default:true"`
}

// Membership represents the association between accounts and groups
type Membership struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	AccountID uint      `json:"account_id" gorm:"index"`
	GroupID   uint      `json:"group_id" gorm:"index"`
}

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Timestamp   time.Time `json:"timestamp"`
	Action      string    `json:"action"`
	ResourceID  uint      `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Details     string    `json:"details"`
	Section     string    `json:"section"`
}

// User represents an authenticated user of the application
type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Username    string    `json:"username" gorm:"unique"`
	Password    string    `json:"-"` // Never expose in JSON
	Email       string    `json:"email" gorm:"unique"`
	Role        string    `json:"role"`
	TOTPEnabled bool      `json:"totp_enabled"`
	TOTPSecret  string    `json:"-"` // Store securely, never expose in JSON
	LastLogin   time.Time `json:"last_login"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TOTPSetupResponse is returned when a user sets up TOTP
type TOTPSetupResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

// TOTPVerifyRequest contains the TOTP code for verification
type TOTPVerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

// AuthResponse is the response after successful authentication
type AuthResponse struct {
	Token       string `json:"token"`
	RequiresTOTP bool   `json:"requires_totp,omitempty"`
	User        UserResponse `json:"user"`
}

// UserResponse contains information about the logged-in user
type UserResponse struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	TOTPEnabled bool   `json:"totp_enabled"`
}