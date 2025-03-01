package service

import (
	"github.com/home/unixify/internal/repository"
)

// Deps is a holder for dependencies needed by services
type Deps struct {
	Repos *repository.Repositories
}

// Services is a holder for all services
type Services struct {
	Account *AccountService
	Group   *GroupService
	Audit   *AuditService
}

// NewServices creates new instances of all services
func NewServices(deps Deps) *Services {
	return &Services{
		Account: NewAccountService(deps.Repos.Account, deps.Repos.Group, deps.Repos.Audit),
		Group:   NewGroupService(deps.Repos.Group, deps.Repos.Account, deps.Repos.Audit),
		Audit:   NewAuditService(deps.Repos.Audit),
	}
}