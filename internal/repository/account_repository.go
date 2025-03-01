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

// FindByID finds an account by ID
func (r *AccountRepository) FindByID(id uint) (*models.Account, error) {
	var account models.Account
	err := r.db.Preload("PrimaryGroup").First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account with ID %d not found", id)
		}
		return nil, err
	}
	return &account, nil
}

// FindByUID finds an account by UID
func (r *AccountRepository) FindByUID(uid int) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("uid = ?", uid).First(&account).Error
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

	// Preload primary group for each account
	err := query.Preload("PrimaryGroup").Find(&accounts).Error
	if err != nil {
		return nil, err
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
	err := r.db.Where("username LIKE ? OR CAST(uid AS TEXT) LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}