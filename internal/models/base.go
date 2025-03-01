package models

import (
	"time"
)

// Base contains common fields used in all models
type Base struct {
	ID        string     `json:"id" db:"id"` 
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	UpdatedBy string     `json:"updated_by" db:"updated_by"`
	DeletedBy *string    `json:"deleted_by,omitempty" db:"deleted_by"`
}