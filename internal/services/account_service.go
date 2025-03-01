package services

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/home/unixify/feature/ai_assisted/internal/models"
)

// AccountService handles operations on accounts
type AccountService struct {
	db *sqlx.DB
}

// NewAccountService creates a new account service
func NewAccountService(db *sqlx.DB) *AccountService {
	return &AccountService{db: db}
}

// FindAll retrieves accounts with pagination and optional search
func (s *AccountService) FindAll(page, pageSize int, search string, currentUserID string) ([]models.Account, int, error) {
	offset := (page - 1) * pageSize
	
	// Base query with soft delete filter
	countQuery := `SELECT COUNT(*) FROM accounts WHERE deleted_at IS NULL`
	query := `SELECT * FROM accounts WHERE deleted_at IS NULL`
	
	// Add search condition if provided
	args := []interface{}{}
	if search != "" {
		countQuery += ` AND (username ILIKE $1 OR email ILIKE $1 OR full_name ILIKE $1)`
		query += ` AND (username ILIKE $1 OR email ILIKE $1 OR full_name ILIKE $1)`
		args = append(args, "%"+search+"%")
	}
	
	// Add pagination
	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, pageSize, offset)
	
	// Get total count
	var total int
	err := s.db.Get(&total, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, err
	}
	
	// Get accounts
	accounts := []models.Account{}
	err = s.db.Select(&accounts, query, args...)
	if err != nil {
		return nil, 0, err
	}
	
	return accounts, total, nil
}

// FindByID retrieves an account by ID
func (s *AccountService) FindByID(id string) (*models.Account, error) {
	account := models.Account{}
	err := s.db.Get(&account, "SELECT * FROM accounts WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	
	// Get account groups
	groups := []models.Group{}
	err = s.db.Select(&groups, `
		SELECT g.* FROM groups g
		INNER JOIN account_groups ag ON g.id = ag.group_id
		WHERE ag.account_id = $1 AND g.deleted_at IS NULL
	`, id)
	
	if err != nil {
		return nil, err
	}
	
	account.Groups = groups
	return &account, nil
}

// Create creates a new account
func (s *AccountService) Create(account *models.Account, currentUserID string) error {
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	account.CreatedBy = currentUserID
	account.UpdatedBy = currentUserID

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert account
	_, err = tx.NamedExec(`
		INSERT INTO accounts (
			id, username, email, full_name, password, active, 
			created_at, updated_at, created_by, updated_by
		) VALUES (
			:id, :username, :email, :full_name, :password, :active, 
			:created_at, :updated_at, :created_by, :updated_by
		)
	`, account)
	
	if err != nil {
		return err
	}

	// Save account groups
	for _, group := range account.Groups {
		_, err = tx.Exec(`
			INSERT INTO account_groups (account_id, group_id, created_at, updated_at, created_by, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, account.ID, group.ID, time.Now(), time.Now(), currentUserID, currentUserID)
		
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Update updates an existing account
func (s *AccountService) Update(account *models.Account, currentUserID string) error {
	account.UpdatedAt = time.Now()
	account.UpdatedBy = currentUserID

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update account
	_, err = tx.NamedExec(`
		UPDATE accounts SET 
			username = :username,
			email = :email,
			full_name = :full_name,
			active = :active,
			updated_at = :updated_at,
			updated_by = :updated_by
		WHERE id = :id
	`, account)
	
	if err != nil {
		return err
	}

	// Update password if provided
	if account.Password != "" {
		_, err = tx.Exec(`
			UPDATE accounts SET 
				password = $1,
				updated_at = $2,
				updated_by = $3
			WHERE id = $4
		`, account.Password, account.UpdatedAt, currentUserID, account.ID)
		
		if err != nil {
			return err
		}
	}

	// Delete existing group associations
	_, err = tx.Exec("DELETE FROM account_groups WHERE account_id = $1", account.ID)
	if err != nil {
		return err
	}

	// Insert new group associations
	for _, group := range account.Groups {
		_, err = tx.Exec(`
			INSERT INTO account_groups (account_id, group_id, created_at, updated_at, created_by, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, account.ID, group.ID, time.Now(), time.Now(), currentUserID, currentUserID)
		
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Delete soft-deletes an account
func (s *AccountService) Delete(id string, currentUserID string) error {
	now := time.Now()
	_, err := s.db.Exec(`
		UPDATE accounts SET 
			deleted_at = $1,
			deleted_by = $2
		WHERE id = $3
	`, now, currentUserID, id)
	
	return err
}