package service

import (
	"fmt"
	"time"

	"github.com/home/unixify/internal/models"
	"github.com/home/unixify/internal/repository"
	"github.com/home/unixify/internal/validator"
)

// GroupService handles business logic for groups
type GroupService struct {
	groupRepo   *repository.GroupRepository
	accountRepo *repository.AccountRepository
	auditRepo   *repository.AuditRepository
}

// NewGroupService creates a new group service
func NewGroupService(
	groupRepo *repository.GroupRepository,
	accountRepo *repository.AccountRepository,
	auditRepo *repository.AuditRepository,
) *GroupService {
	return &GroupService{
		groupRepo:   groupRepo,
		accountRepo: accountRepo,
		auditRepo:   auditRepo,
	}
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(group *models.Group, userID uint, username, ipAddress string) error {
	// Validate GID - now just a warning
	if err := validator.ValidateGIDForType(group.UnixGID, group.Type); err != nil {
		// If it's a warning (starts with "WARNING:"), log it but continue
		if len(err.Error()) >= 7 && err.Error()[:7] == "WARNING" {
			// Log the warning but continue
			fmt.Println(err.Error())
		} else {
			// It's a real error, return it
			return err
		}
	}

	// Check if GID already exists using the improved method
	isDuplicate, err := s.groupRepo.IsGIDDuplicate(group.UnixGID, 0)
	if err != nil {
		return err
	}
	if isDuplicate {
		return fmt.Errorf("group with GID %d already exists", group.UnixGID)
	}

	// Check if groupname already exists
	existingGroup, err := s.groupRepo.FindByGroupname(group.Groupname)
	if err == nil && existingGroup != nil {
		return fmt.Errorf("group with groupname %s already exists", group.Groupname)
	}

	// Create group
	err = s.groupRepo.Create(group)
	if err != nil {
		return err
	}

	// Log audit entry
	auditEntry := &models.AuditEntry{
		Action:     "create",
		EntityID:   group.ID,
		EntityType: "group",
		Details:    fmt.Sprintf("Created group %s with GID %d", group.Groupname, group.UnixGID),
		UserID:     userID,
		Username:   username,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
	}
	return s.auditRepo.Create(auditEntry)
}

// GetGroup gets a group by ID
func (s *GroupService) GetGroup(id uint) (*models.Group, error) {
	return s.groupRepo.FindByID(id)
}

// GetGroupByGID gets a group by GID
func (s *GroupService) GetGroupByGID(gid int) (*models.Group, error) {
	return s.groupRepo.FindByGID(gid)
}

// GetGroupByGroupname gets a group by groupname
func (s *GroupService) GetGroupByGroupname(groupname string) (*models.Group, error) {
	return s.groupRepo.FindByGroupname(groupname)
}

// GetAllGroups gets all groups with optional filtering by type
func (s *GroupService) GetAllGroups(groupType models.GroupType) ([]models.Group, error) {
	return s.groupRepo.FindAll(groupType)
}

// UpdateGroup updates a group
func (s *GroupService) UpdateGroup(group *models.Group, userID uint, username, ipAddress string) error {
	// Validate GID - now just a warning
	if err := validator.ValidateGIDForType(group.UnixGID, group.Type); err != nil {
		// If it's a warning (starts with "WARNING:"), log it but continue
		if len(err.Error()) >= 7 && err.Error()[:7] == "WARNING" {
			// Log the warning but continue
			fmt.Println(err.Error())
		} else {
			// It's a real error, return it
			return err
		}
	}

	// Get the original group to check if GID is changing
	originalGroup, err := s.groupRepo.FindByID(group.ID)
	if err != nil {
		return err
	}

	// Check if GID already exists using the improved method
	isDuplicate, err := s.groupRepo.IsGIDDuplicate(group.UnixGID, group.ID)
	if err != nil {
		return err
	}
	if isDuplicate {
		return fmt.Errorf("group with GID %d already exists", group.UnixGID)
	}

	// If groupname is changing, check if new groupname already exists
	if originalGroup.Groupname != group.Groupname {
		existingGroup, err := s.groupRepo.FindByGroupname(group.Groupname)
		if err == nil && existingGroup != nil && existingGroup.ID != group.ID {
			return fmt.Errorf("group with groupname %s already exists", group.Groupname)
		}
	}

	// Update group
	err = s.groupRepo.Update(group)
	if err != nil {
		return err
	}

	// Log audit entry
	auditEntry := &models.AuditEntry{
		Action:     "update",
		EntityID:   group.ID,
		EntityType: "group",
		Details:    fmt.Sprintf("Updated group %s with GID %d", group.Groupname, group.UnixGID),
		UserID:     userID,
		Username:   username,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
	}
	return s.auditRepo.Create(auditEntry)
}

// DeleteGroup soft deletes a group
func (s *GroupService) DeleteGroup(id uint, userID uint, username, ipAddress string) error {
	// Get group to record groupname in audit
	group, err := s.groupRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete group
	err = s.groupRepo.Delete(id)
	if err != nil {
		return err
	}

	// Log audit entry
	auditEntry := &models.AuditEntry{
		Action:     "delete",
		EntityID:   id,
		EntityType: "group",
		Details:    fmt.Sprintf("Deleted group %s with GID %d", group.Groupname, group.UnixGID),
		UserID:     userID,
		Username:   username,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
	}
	return s.auditRepo.Create(auditEntry)
}

// GetGroupMembers gets all accounts in a specific group
func (s *GroupService) GetGroupMembers(groupID uint) ([]models.Account, error) {
	return s.groupRepo.GetAccountsInGroup(groupID)
}

// SearchGroups searches for groups by GID or groupname
func (s *GroupService) SearchGroups(query string) ([]models.Group, error) {
	return s.groupRepo.Search(query)
}

// IsGIDDuplicate checks if a GID already exists
func (s *GroupService) IsGIDDuplicate(gid int, excludeID uint) (bool, error) {
	return s.groupRepo.IsGIDDuplicate(gid, excludeID)
}

// GetNextAvailableGID gets the next available GID for a specific group type
func (s *GroupService) GetNextAvailableGID(groupType models.GroupType) (int, error) {
	return s.groupRepo.GetLatestGID(groupType)
}