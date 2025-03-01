package models

// Group represents a user group in the system
type Group struct {
	Base
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Members     []Account `json:"members,omitempty" db:"-"` // Many-to-many relationship
}