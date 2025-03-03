package models

import (
	"time"
)

// Account type constants
type AccountType string

const (
	AccountTypePeople   AccountType = "people"
	AccountTypeSystem   AccountType = "system"
	AccountTypeDatabase AccountType = "database"
	AccountTypeService  AccountType = "service"
)

// Group type constants
type GroupType string

const (
	GroupTypePeople   GroupType = "people"
	GroupTypeSystem   GroupType = "system"
	GroupTypeDatabase GroupType = "database"
	GroupTypeService  GroupType = "service"
)

// Account represents a UNIX account (user)
type Account struct {
	ID             uint        `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      time.Time   `json:"deleted_at" gorm:"index"`
	Username       string      `json:"username" gorm:"unique"`
	UnixUID        int         `json:"uid" gorm:"column:unixuid;unique"` // Using uid for JSON but unixuid for column name
	Type           AccountType `json:"type" gorm:"index"` // people, system, database, service
	PrimaryGroupID uint        `json:"primary_group_id" gorm:"index"`
	PrimaryGroup   *Group      `json:"primary_group" gorm:"-"` // Ignore this field for GORM but keep in JSON
	Active         bool        `json:"active" gorm:"default:true"`
	Firstname      string      `json:"firstname"` // First name of the account owner
	Surname        string      `json:"surname"`   // Last name/surname of the account owner
}

// Group represents a UNIX group
type Group struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"index"`
	Groupname   string    `json:"groupname" gorm:"unique"`
	UnixGID     int       `json:"gid" gorm:"column:unixgid;unique"` // Using gid for JSON but unixgid for column name
	Type        GroupType `json:"type" gorm:"index"` // people, system, database, service
	Description string    `json:"description"`
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedBy   string    `json:"created_by"` // Username of the person who created this group
}

// Membership represents the association between accounts and groups
type Membership struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	AccountID uint      `json:"account_id" gorm:"index"`
	GroupID   uint      `json:"group_id" gorm:"index"`
}

// AccountGroup represents the association between accounts and groups
// This is an alias for Membership to maintain compatibility with existing code
type AccountGroup struct {
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
	EntityID    uint      `json:"entity_id"`      // Alias for ResourceID
	EntityType  string    `json:"entity_type"`    // Alias for ResourceType
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Details     string    `json:"details"`
	Section     string    `json:"section"`
	IPAddress   string    `json:"ip_address"`
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

// RegisteredUser represents a registered user with enhanced security
type RegisteredUser struct {
	ID           uint      `json:"id" gorm:"primaryKey" sql:"AUTO_INCREMENT"`
	Username     string    `json:"username" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash []byte    `json:"-" gorm:"column:password_hash;not null"` // Stored encrypted, never expose in JSON
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name" gorm:"not null"`
	Department   string    `json:"department"`
	Role         string    `json:"role" gorm:"not null;default:'user'"`
	TOTPEnabled  bool      `json:"totp_enabled" gorm:"default:false"`
	TOTPSecret   string    `json:"-" gorm:"column:totp_secret"` // Store securely, never expose in JSON
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	LastLogin    time.Time `json:"last_login"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterUserRequest represents a request to register a new user
type RegisterUserRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Department  string `json:"department"`
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