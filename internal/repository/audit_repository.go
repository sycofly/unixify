package repository

import (
	"github.com/home/unixify/internal/models"
	"gorm.io/gorm"
)

// AuditRepository handles database operations for audit logs
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{
		db: db,
	}
}

// Create creates a new audit entry
func (r *AuditRepository) Create(entry *models.AuditEntry) error {
	return r.db.Create(entry).Error
}

// FindAll retrieves all audit entries with optional filtering
func (r *AuditRepository) FindAll(entityType, action string, entityID, userID uint) ([]models.AuditEntry, error) {
	var entries []models.AuditEntry
	query := r.db

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if entityID != 0 {
		query = query.Where("entity_id = ?", entityID)
	}

	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Order("timestamp DESC").Find(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// FindByID finds an audit entry by ID
func (r *AuditRepository) FindByID(id uint) (*models.AuditEntry, error) {
	var entry models.AuditEntry
	err := r.db.First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}