package service

import (
	"github.com/home/unixify/internal/models"
	"github.com/home/unixify/internal/repository"
)

// AuditService handles business logic for audit logs
type AuditService struct {
	auditRepo *repository.AuditRepository
}

// NewAuditService creates a new audit service
func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
}

// CreateAuditEntry creates a new audit entry
func (s *AuditService) CreateAuditEntry(entry *models.AuditEntry) error {
	return s.auditRepo.Create(entry)
}

// GetAuditEntries gets all audit entries with optional filtering
func (s *AuditService) GetAuditEntries(entityType, action string, entityID, userID uint) ([]models.AuditEntry, error) {
	return s.auditRepo.FindAll(entityType, action, entityID, userID)
}

// GetAuditEntry gets an audit entry by ID
func (s *AuditService) GetAuditEntry(id uint) (*models.AuditEntry, error) {
	return s.auditRepo.FindByID(id)
}