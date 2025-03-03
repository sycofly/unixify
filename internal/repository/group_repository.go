package repository

import (
	"errors"
	"fmt"

	"github.com/home/unixify/internal/models"
	"gorm.io/gorm"
)

// GroupRepository handles database operations for groups
type GroupRepository struct {
	db *gorm.DB
}

// NewGroupRepository creates a new group repository
func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

// Create creates a new group
func (r *GroupRepository) Create(group *models.Group) error {
	return r.db.Create(group).Error
}

// IsGIDDuplicate checks if a GID already exists
func (r *GroupRepository) IsGIDDuplicate(gid int, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Group{}).Where("unixgid = ?", gid)
	
	// Exclude current group if updating
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// GetLatestGID returns the next available GID for a specific group type
func (r *GroupRepository) GetLatestGID(groupType models.GroupType) (int, error) {
	var group models.Group
	
	err := r.db.Where("type = ?", groupType).Order("unixgid DESC").First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No groups of this type found, return the minimum GID for this type
			switch groupType {
			case models.GroupTypePeople:
				return 1000, nil
			case models.GroupTypeSystem:
				return 9000, nil
			case models.GroupTypeService:
				return 60001, nil
			case models.GroupTypeDatabase:
				return 70000, nil
			default:
				return 1000, nil
			}
		}
		return 0, err
	}
	
	// Return the next available GID (current highest + 1)
	return group.UnixGID + 1, nil
}

// FindByID finds a group by ID
func (r *GroupRepository) FindByID(id uint) (*models.Group, error) {
	var group models.Group
	err := r.db.First(&group, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("group with ID %d not found", id)
		}
		return nil, err
	}
	return &group, nil
}

// FindByGID finds a group by GID
func (r *GroupRepository) FindByGID(gid int) (*models.Group, error) {
	var group models.Group
	err := r.db.Where("gid = ?", gid).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("group with GID %d not found", gid)
		}
		return nil, err
	}
	return &group, nil
}

// FindByGroupname finds a group by groupname
func (r *GroupRepository) FindByGroupname(groupname string) (*models.Group, error) {
	var group models.Group
	err := r.db.Where("groupname = ?", groupname).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("group with groupname %s not found", groupname)
		}
		return nil, err
	}
	return &group, nil
}

// FindAll finds all groups with optional filtering by type
func (r *GroupRepository) FindAll(groupType models.GroupType) ([]models.Group, error) {
	var groups []models.Group
	query := r.db

	if groupType != "" {
		query = query.Where("type = ?", groupType)
	}

	err := query.Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// Update updates a group
func (r *GroupRepository) Update(group *models.Group) error {
	return r.db.Save(group).Error
}

// Delete soft deletes a group
func (r *GroupRepository) Delete(id uint) error {
	return r.db.Delete(&models.Group{}, id).Error
}

// FindByAccountID finds all groups that an account is a member of
func (r *GroupRepository) FindByAccountID(accountID uint) ([]models.Group, error) {
	var groups []models.Group
	err := r.db.Joins("JOIN account_groups ON account_groups.group_id = groups.id").
		Where("account_groups.account_id = ?", accountID).
		Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// GetAccountsInGroup returns all accounts in a group
func (r *GroupRepository) GetAccountsInGroup(groupID uint) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Joins("JOIN account_groups ON account_groups.account_id = accounts.id").
		Where("account_groups.group_id = ?", groupID).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// Search searches for groups by GID or groupname
func (r *GroupRepository) Search(query string) ([]models.Group, error) {
	var groups []models.Group
	err := r.db.Where("groupname LIKE ? OR CAST(gid AS TEXT) LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}