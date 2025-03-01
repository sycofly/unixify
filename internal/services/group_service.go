package services

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/home/unixify/feature/ai_assisted/internal/models"
)

// GroupService handles operations on groups
type GroupService struct {
	db *sqlx.DB
}

// NewGroupService creates a new group service
func NewGroupService(db *sqlx.DB) *GroupService {
	return &GroupService{db: db}
}

// FindAll retrieves groups with pagination and optional search
func (s *GroupService) FindAll(page, pageSize int, search string, currentUserID string) ([]models.Group, int, error) {
	offset := (page - 1) * pageSize
	
	// Base query with soft delete filter
	countQuery := `SELECT COUNT(*) FROM groups WHERE deleted_at IS NULL`
	query := `SELECT * FROM groups WHERE deleted_at IS NULL`
	
	// Add search condition if provided
	args := []interface{}{}
	if search != "" {
		countQuery += ` AND (name ILIKE $1 OR description ILIKE $1)`
		query += ` AND (name ILIKE $1 OR description ILIKE $1)`
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
	
	// Get groups
	groups := []models.Group{}
	err = s.db.Select(&groups, query, args...)
	if err != nil {
		return nil, 0, err
	}
	
	return groups, total, nil
}

// FindByID retrieves a group by ID
func (s *GroupService) FindByID(id string) (*models.Group, error) {
	group := models.Group{}
	err := s.db.Get(&group, "SELECT * FROM groups WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	
	// Get group members
	members := []models.Account{}
	err = s.db.Select(&members, `
		SELECT a.* FROM accounts a
		INNER JOIN account_groups ag ON a.id = ag.account_id
		WHERE ag.group_id = $1 AND a.deleted_at IS NULL
	`, id)
	
	if err != nil {
		return nil, err
	}
	
	group.Members = members
	return &group, nil
}

// Create creates a new group
func (s *GroupService) Create(group *models.Group, currentUserID string) error {
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	group.CreatedBy = currentUserID
	group.UpdatedBy = currentUserID

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert group
	_, err = tx.NamedExec(`
		INSERT INTO groups (
			id, name, description, created_at, updated_at, created_by, updated_by
		) VALUES (
			:id, :name, :description, :created_at, :updated_at, :created_by, :updated_by
		)
	`, group)
	
	if err != nil {
		return err
	}

	// Save group members
	for _, member := range group.Members {
		_, err = tx.Exec(`
			INSERT INTO account_groups (account_id, group_id, created_at, updated_at, created_by, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, member.ID, group.ID, time.Now(), time.Now(), currentUserID, currentUserID)
		
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Update updates an existing group
func (s *GroupService) Update(group *models.Group, currentUserID string) error {
	group.UpdatedAt = time.Now()
	group.UpdatedBy = currentUserID

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update group
	_, err = tx.NamedExec(`
		UPDATE groups SET 
			name = :name,
			description = :description,
			updated_at = :updated_at,
			updated_by = :updated_by
		WHERE id = :id
	`, group)
	
	if err != nil {
		return err
	}

	// Delete existing member associations
	_, err = tx.Exec("DELETE FROM account_groups WHERE group_id = $1", group.ID)
	if err != nil {
		return err
	}

	// Insert new member associations
	for _, member := range group.Members {
		_, err = tx.Exec(`
			INSERT INTO account_groups (account_id, group_id, created_at, updated_at, created_by, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, member.ID, group.ID, time.Now(), time.Now(), currentUserID, currentUserID)
		
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Delete soft-deletes a group
func (s *GroupService) Delete(id string, currentUserID string) error {
	now := time.Now()
	_, err := s.db.Exec(`
		UPDATE groups SET 
			deleted_at = $1,
			deleted_by = $2
		WHERE id = $3
	`, now, currentUserID, id)
	
	return err
}