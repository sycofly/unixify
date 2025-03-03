package repository

import (
	"errors"
	"fmt"

	"github.com/home/unixify/internal/models"
	"gorm.io/gorm"
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// Create creates a new account
func (r *AccountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

// IsUIDDuplicate checks if a UID already exists
func (r *AccountRepository) IsUIDDuplicate(uid int, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Account{}).Where("unixuid = ?", uid)
	
	// Exclude current account if updating
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// GetLatestUID returns the next available UID for a specific account type
func (r *AccountRepository) GetLatestUID(accountType models.AccountType) (int, error) {
	var account models.Account
	
	err := r.db.Where("type = ?", accountType).Order("unixuid DESC").First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No accounts of this type found, return the minimum UID for this type
			switch accountType {
			case models.AccountTypePeople:
				return 1000, nil
			case models.AccountTypeSystem:
				return 9000, nil
			case models.AccountTypeService:
				return 60001, nil
			case models.AccountTypeDatabase:
				return 70000, nil
			default:
				return 1000, nil
			}
		}
		return 0, err
	}
	
	// Return the next available UID (current highest + 1)
	return account.UnixUID + 1, nil
}

// FindByID finds an account by ID
func (r *AccountRepository) FindByID(id uint) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account with ID %d not found", id)
		}
		return nil, err
	}
	
	// Manually load primary group if set
	if account.PrimaryGroupID > 0 {
		var group models.Group
		if err := r.db.First(&group, account.PrimaryGroupID).Error; err != nil {
			// Log the error but don't fail
			fmt.Printf("Failed to load primary group for account %d: %v\n", account.ID, err)
		} else {
			account.PrimaryGroup = &group
		}
	}
	
	return &account, nil
}

// FindByUID finds an account by UID
func (r *AccountRepository) FindByUID(uid int) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("unixuid = ?", uid).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account with UID %d not found", uid)
		}
		return nil, err
	}
	return &account, nil
}

// FindByUsername finds an account by username
func (r *AccountRepository) FindByUsername(username string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("username = ?", username).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account with username %s not found", username)
		}
		return nil, err
	}
	return &account, nil
}

// FindAll finds all accounts with optional filtering by type
func (r *AccountRepository) FindAll(accountType models.AccountType) ([]models.Account, error) {
	var accounts []models.Account
	query := r.db

	if accountType != "" {
		query = query.Where("type = ?", accountType)
	}

	// Get all accounts without preloading initially
	err := query.Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	// Manually handle the relationship by fetching the primary groups
	// This avoids using GORM preload which is failing
	for i := range accounts {
		if accounts[i].PrimaryGroupID > 0 {
			var group models.Group
			if err := r.db.First(&group, accounts[i].PrimaryGroupID).Error; err != nil {
				// Just log the error and continue
				fmt.Printf("Failed to load primary group for account %d: %v\n", accounts[i].ID, err)
				continue
			}
			accounts[i].PrimaryGroup = &group
		}
	}
	
	return accounts, nil
}

// Update updates an account
func (r *AccountRepository) Update(account *models.Account) error {
	// Use a transaction to ensure atomicity
	tx := r.db.Begin()
	
	// Save the account - this updates the record
	if err := tx.Save(account).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update account: %w", err)
	}
	
	// Verify the update was successful by reloading the account
	var updatedAccount models.Account
	if err := tx.First(&updatedAccount, account.ID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to verify account update: %w", err)
	}
	
	// Ensure type was updated correctly
	if updatedAccount.Type != account.Type {
		tx.Rollback()
		return fmt.Errorf("account type was not updated correctly: expected %s, got %s", account.Type, updatedAccount.Type)
	}
	
	// Commit the transaction
	return tx.Commit().Error
}

// Delete soft deletes an account
func (r *AccountRepository) Delete(id uint) error {
	return r.db.Delete(&models.Account{}, id).Error
}

// FindByGroupID finds all accounts in a specific group
func (r *AccountRepository) FindByGroupID(groupID uint) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Joins("JOIN account_groups ON account_groups.account_id = accounts.id").
		Where("account_groups.group_id = ?", groupID).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// AssignToGroup assigns an account to a group
func (r *AccountRepository) AssignToGroup(accountID, groupID uint) error {
	accountGroup := models.AccountGroup{
		AccountID: accountID,
		GroupID:   groupID,
	}
	
	// Check if the relation already exists
	var count int64
	err := r.db.Model(&models.AccountGroup{}).
		Where("account_id = ? AND group_id = ?", accountID, groupID).
		Count(&count).Error
	if err != nil {
		return err
	}
	
	if count > 0 {
		return fmt.Errorf("account is already a member of this group")
	}
	
	return r.db.Create(&accountGroup).Error
}

// RemoveFromGroup removes an account from a group
func (r *AccountRepository) RemoveFromGroup(accountID, groupID uint) error {
	return r.db.Where("account_id = ? AND group_id = ?", accountID, groupID).
		Delete(&models.AccountGroup{}).Error
}

// Search searches for accounts by UID or username
func (r *AccountRepository) Search(query string) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("username LIKE ? OR CAST(unixuid AS TEXT) LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}