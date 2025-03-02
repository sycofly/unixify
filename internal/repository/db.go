package repository

import (
	"fmt"

	"github.com/home/unixify/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Skip auto-migration and rely on the SQL migrations
	// that have already been applied via the migrate tool.
	// This avoids issues with column naming differences.

	return db, nil
}

// Repositories is a holder for all repositories
type Repositories struct {
	Account *AccountRepository
	Group   *GroupRepository
	Audit   *AuditRepository
}

// Repository is an alias for Repositories for backward compatibility
type Repository struct {
	Account *AccountRepository
	Group   *GroupRepository
	Audit   *AuditRepository
}

// NewRepositories creates new instances of all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Account: NewAccountRepository(db),
		Group:   NewGroupRepository(db),
		Audit:   NewAuditRepository(db),
	}
}

// NewRepository creates new instances of all repositories (alias for NewRepositories)
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Account: NewAccountRepository(db),
		Group:   NewGroupRepository(db),
		Audit:   NewAuditRepository(db),
	}
}