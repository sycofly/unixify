package models

import (
	"time"
)

// Account represents a user account in the system
type Account struct {
	Base
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	FullName  string    `json:"full_name" db:"full_name"`
	Password  string    `json:"-" db:"password"` // Never expose in JSON
	LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
	Active    bool      `json:"active" db:"active"`
	Groups    []Group   `json:"groups,omitempty" db:"-"` // Many-to-many relationship
}