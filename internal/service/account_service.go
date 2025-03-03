package service

import (
	"fmt"
	"time"

	"github.com/home/unixify/internal/models"
	"github.com/home/unixify/internal/repository"
	"github.com/home/unixify/internal/validator"
)

// AccountService handles business logic for accounts
type AccountService struct {
	accountRepo *repository.AccountRepository
	groupRepo   *repository.GroupRepository
	auditRepo   *repository.AuditRepository
}

// NewAccountService creates a new account service
func NewAccountService(
	accountRepo *repository.AccountRepository,
	groupRepo *repository.GroupRepository,
	auditRepo *repository.AuditRepository,
) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		groupRepo:   groupRepo,
		auditRepo:   auditRepo,
	}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(account *models.Account, userID uint, username, ipAddress string) error {
	// Validate UID - now just a warning
	if err := validator.ValidateUIDForType(account.UnixUID, account.Type); err != nil {
		// If it's a warning (starts with "WARNING:"), log it but continue
		if len(err.Error()) >= 7 && err.Error()[:7] == "WARNING" {
			// Log the warning but continue
			fmt.Println(err.Error())
		} else {
			// It's a real error, return it
			return err
		}
	}

	// Check if UID already exists using the improved method
	isDuplicate, err := s.accountRepo.IsUIDDuplicate(account.UnixUID, 0)
	if err != nil {
		return err
	}
	if isDuplicate {
		return fmt.Errorf("account with UID %d already exists", account.UnixUID)
	}

	// Check if username already exists
	existingAccount, err := s.accountRepo.FindByUsername(account.Username)
	if err == nil && existingAccount != nil {
		return fmt.Errorf("account with username %s already exists", account.Username)
	}
	
	// Validate primary group if specified
	if account.PrimaryGroupID != 0 {
		primaryGroup, err := s.groupRepo.FindByID(account.PrimaryGroupID)
		if err != nil {
			return fmt.Errorf("primary group with ID %d not found", account.PrimaryGroupID)
		}
		
		// System accounts must have a system group as primary group
		if account.Type == models.AccountTypeSystem && primaryGroup.Type != models.GroupTypeSystem {
			return fmt.Errorf("system accounts must have a system group as primary group")
		}
	} else if account.Type == models.AccountTypeSystem {
		// System accounts must have a primary group
		return fmt.Errorf("system accounts must have a primary group")
	}

	// Create account
	err = s.accountRepo.Create(account)
	if err != nil {
		return err
	}

	// Log audit entry
	auditEntry := &models.AuditEntry{
		Action:     "create",
		EntityID:   account.ID,
		EntityType: "account",
		Details:    fmt.Sprintf("Created account %s with UID %d", account.Username, account.UnixUID),
		UserID:     userID,
		Username:   username,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
	}
	return s.auditRepo.Create(auditEntry)
}

// GetAccount gets an account by ID
func (s *AccountService) GetAccount(id uint) (*models.Account, error) {
	return s.accountRepo.FindByID(id)
}

// GetAccountByUID gets an account by UID
func (s *AccountService) GetAccountByUID(uid int) (*models.Account, error) {
	return s.accountRepo.FindByUID(uid)
}

// GetAccountByUsername gets an account by username
func (s *AccountService) GetAccountByUsername(username string) (*models.Account, error) {
	return s.accountRepo.FindByUsername(username)
}

// GetAllAccounts gets all accounts with optional filtering by type
func (s *AccountService) GetAllAccounts(accountType models.AccountType) ([]models.Account, error) {
	return s.accountRepo.FindAll(accountType)
}

// UpdateAccount updates an account
func (s *AccountService) UpdateAccount(account *models.Account, userID uint, username, ipAddress string) error {
	// Validate UID - now just a warning
	if err := validator.ValidateUIDForType(account.UnixUID, account.Type); err != nil {
		// If it's a warning (starts with "WARNING:"), log it but continue
		if len(err.Error()) >= 7 && err.Error()[:7] == "WARNING" {
			// Log the warning but continue
			fmt.Println(err.Error())
		} else {
			// It's a real error, return it
			return err
		}
	}

	// Get the original account to check if UID is changing
	originalAccount, err := s.accountRepo.FindByID(account.ID)
	if err != nil {
		return err
	}

	// Check if UID already exists using the improved method
	isDuplicate, err := s.accountRepo.IsUIDDuplicate(account.UnixUID, account.ID)
	if err != nil {
		return err
	}
	if isDuplicate {
		return fmt.Errorf("account with UID %d already exists", account.UnixUID)
	}

	// If username is changing, check if new username already exists
	if originalAccount.Username != account.Username {
		existingAccount, err := s.accountRepo.FindByUsername(account.Username)
		if err == nil && existingAccount != nil && existingAccount.ID != account.ID {
			return fmt.Errorf("account with username %s already exists", account.Username)
		}
	}
	
	// Validate primary group if specified
	if account.PrimaryGroupID != 0 {
		primaryGroup, err := s.groupRepo.FindByID(account.PrimaryGroupID)
		if err != nil {
			return fmt.Errorf("primary group with ID %d not found", account.PrimaryGroupID)
		}
		
		// System accounts must have a system group as primary group
		if account.Type == models.AccountTypeSystem && primaryGroup.Type != models.GroupTypeSystem {
			return fmt.Errorf("system accounts must have a system group as primary group")
		}
	} else if account.Type == models.AccountTypeSystem {
		// System accounts must have a primary group
		return fmt.Errorf("system accounts must have a primary group")
	}

	// Update account
	err = s.accountRepo.Update(account)
	if err != nil {
		return err
	}

	// Log audit entry
	auditEntry := &models.AuditEntry{
		Action:     "update",
		EntityID:   account.ID,
		EntityType: "account",
		Details:    fmt.Sprintf("Updated account %s with UID %d", account.Username, account.UnixUID),
		UserID:     userID,
		Username:   username,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
	}
	return s.auditRepo.Create(auditEntry)
}

// DeleteAccount soft deletes an account
func (s *AccountService) DeleteAccount(id uint, userID uint, username, ipAddress string) error {
	// Get account to record username in audit
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete account
	err = s.accountRepo.Delete(id)
	if err != nil {
		return err
	}

	// Log audit entry
	auditEntry := &models.AuditEntry{
		Action:     "delete",
		EntityID:   id,
		EntityType: "account",
		Details:    fmt.Sprintf("Deleted account %s with UID %d", account.Username, account.UnixUID),
		UserID:     userID,
		Username:   username,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
	}
	return s.auditRepo.Create(auditEntry)
}

// GetAccountsInGroup gets all accounts in a specific group
func (s *AccountService) GetAccountsInGroup(groupID uint) ([]models.Account, error) {
	return s.accountRepo.FindByGroupID(groupID)
}

// AssignAccountToGroup assigns an account to a group
func (s *AccountService) AssignAccountToGroup(accountID, groupID uint, userID uint, username, ipAddress string) error {
	// Check if account exists
	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return err
	}

	// Check if group exists
	group, err := s.groupRepo.FindByID(groupID)
	if err != nil {
		return err
	}

	// Validate compatibility between account type and group type
	if !validator.IsValidAccountGroupAssignment(account.Type, group.Type) {
		return fmt.Errorf("account of type %s cannot be assigned to group of type %s", account.Type, group.Type)
	}

	// Assign account to group
	err = s.accountRepo.AssignToGroup(accountID, groupID)
	if err != nil {
		return err
	}

	// Log audit entry
	auditEntry := &models.AuditEntry{
		Action:     "assign",
		EntityID:   accountID,
		EntityType: "account_group",
		Details:    fmt.Sprintf("Assigned account %s (ID: %d) to group %s (ID: %d)", account.Username, accountID, group.Groupname, groupID),
		UserID:     userID,
		Username:   username,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
	}
	return s.auditRepo.Create(auditEntry)
}

// RemoveAccountFromGroup removes an account from a group
func (s *AccountService) RemoveAccountFromGroup(accountID, groupID uint, userID uint, username, ipAddress string) error {
	// Check if account exists
	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return err
	}

	// Check if group exists
	group, err := s.groupRepo.FindByID(groupID)
	if err != nil {
		return err
	}

	// Remove account from group
	err = s.accountRepo.RemoveFromGroup(accountID, groupID)
	if err != nil {
		return err
	}

	// Log audit entry
	auditEntry := &models.AuditEntry{
		Action:     "remove",
		EntityID:   accountID,
		EntityType: "account_group",
		Details:    fmt.Sprintf("Removed account %s (ID: %d) from group %s (ID: %d)", account.Username, accountID, group.Groupname, groupID),
		UserID:     userID,
		Username:   username,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
	}
	return s.auditRepo.Create(auditEntry)
}

// GetAccountGroups gets all groups that an account is a member of
func (s *AccountService) GetAccountGroups(accountID uint) ([]models.Group, error) {
	return s.groupRepo.FindByAccountID(accountID)
}

// SearchAccounts searches for accounts by UID or username
func (s *AccountService) SearchAccounts(query string) ([]models.Account, error) {
	return s.accountRepo.Search(query)
}

// IsUIDDuplicate checks if a UID already exists
func (s *AccountService) IsUIDDuplicate(uid int, excludeID uint) (bool, error) {
	return s.accountRepo.IsUIDDuplicate(uid, excludeID)
}

// GetNextAvailableUID gets the next available UID for a specific account type
func (s *AccountService) GetNextAvailableUID(accountType models.AccountType) (int, error) {
	return s.accountRepo.GetLatestUID(accountType)
}