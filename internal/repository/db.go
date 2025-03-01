package repository

import (
	"fmt"

	"github.com/home/unixify/internal/config"
	"github.com/home/unixify/internal/models"
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

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.Account{},
		&models.Group{},
		&models.AccountGroup{},
		&models.AuditEntry{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// Repositories is a holder for all repositories
type Repositories struct {
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