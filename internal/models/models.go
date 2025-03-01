package models

import (
	"time"

	"gorm.io/gorm"
)

// AccountType represents the type of user account
type AccountType string

const (
	AccountTypePeople   AccountType = "people"
	AccountTypeSystem   AccountType = "system"
	AccountTypeDatabase AccountType = "database"
	AccountTypeService  AccountType = "service"
)

// GroupType represents the type of group
type GroupType string

const (
	GroupTypePeople   GroupType = "people"
	GroupTypeSystem   GroupType = "system"
	GroupTypeDatabase GroupType = "database"
	GroupTypeService  GroupType = "service"
)

// Account represents a UNIX user account
type Account struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UID           int            `json:"uid" gorm:"uniqueIndex;not null"`
	Username      string         `json:"username" gorm:"uniqueIndex;not null"`
	Type          AccountType    `json:"type" gorm:"type:varchar(20);not null"`
	PrimaryGroupID uint          `json:"primary_group_id" gorm:"default:null"`
	PrimaryGroup  *Group         `json:"primary_group,omitempty" gorm:"foreignKey:PrimaryGroupID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Groups        []*Group       `json:"groups,omitempty" gorm:"many2many:account_groups;"`
}

// Group represents a UNIX group
type Group struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	GID         int            `json:"gid" gorm:"uniqueIndex;not null"`
	Groupname   string         `json:"groupname" gorm:"uniqueIndex;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Type        GroupType      `json:"type" gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Accounts    []*Account     `json:"accounts,omitempty" gorm:"many2many:account_groups;"`
}

// AccountGroup represents the many-to-many relationship between accounts and groups
type AccountGroup struct {
	AccountID uint      `gorm:"primaryKey"`
	GroupID   uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Action    string    `json:"action" gorm:"not null"`
	EntityID  uint      `json:"entity_id"`
	EntityType string   `json:"entity_type" gorm:"not null"`
	Details   string    `json:"details"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IPAddress string    `json:"ip_address"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
}

// TableName sets the table name for AccountGroup model
func (AccountGroup) TableName() string {
	return "account_groups"
}