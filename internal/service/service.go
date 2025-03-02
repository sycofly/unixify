package service

import (
	"github.com/home/unixify/internal/repository"
	"gorm.io/gorm"
)

// Deps is a holder for dependencies needed by services
type Deps struct {
	Repos *repository.Repositories
	DB    *gorm.DB
}

// Services is a holder for all services
type Services struct {
	Account *AccountService
	Group   *GroupService
	Audit   *AuditService
	db      *gorm.DB  // Add DB connection for direct access if needed
}

// NewServices creates new instances of all services
func NewServices(deps Deps) *Services {
	return &Services{
		Account: NewAccountService(deps.Repos.Account, deps.Repos.Group, deps.Repos.Audit),
		Group:   NewGroupService(deps.Repos.Group, deps.Repos.Account, deps.Repos.Audit),
		Audit:   NewAuditService(deps.Repos.Audit),
		db:      deps.DB,
	}
}

// GetDB returns the database connection
func (s *Services) GetDB() *gorm.DB {
	return s.db
}