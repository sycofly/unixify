package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/home/unixify/internal/models"
)

// accountInput represents the input for account creation/update
type accountInput struct {
	UID            int                `json:"uid" binding:"required"`
	Username       string             `json:"username" binding:"required"`
	Type           models.AccountType `json:"type" binding:"required"`
	PrimaryGroupID uint               `json:"primary_group_id"`
}

// GetAllAccounts handles GET /api/accounts
func (h *Handler) GetAllAccounts(c *gin.Context) {
	// Get account type from query parameter (optional)
	accountType := models.AccountType(c.Query("type"))

	// Get all accounts
	accounts, err := h.services.Account.GetAllAccounts(accountType)
	if err != nil {
		h.logger.Errorf("Failed to get accounts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// GetAccount handles GET /api/accounts/:id
func (h *Handler) GetAccount(c *gin.Context) {
	// Parse account ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get account
	account, err := h.services.Account.GetAccount(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get account: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetAccountByUID handles GET /api/accounts/uid/:uid
func (h *Handler) GetAccountByUID(c *gin.Context) {
	// Parse UID
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UID"})
		return
	}

	// Get account
	account, err := h.services.Account.GetAccountByUID(uid)
	if err != nil {
		h.logger.Errorf("Failed to get account by UID: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetAccountByUsername handles GET /api/accounts/username/:username
func (h *Handler) GetAccountByUsername(c *gin.Context) {
	// Get username
	username := c.Param("username")

	// Get account
	account, err := h.services.Account.GetAccountByUsername(username)
	if err != nil {
		h.logger.Errorf("Failed to get account by username: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// CreateAccount handles POST /api/accounts
func (h *Handler) CreateAccount(c *gin.Context) {
	// Parse input
	var input accountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create account
	account := &models.Account{
		UID:            input.UID,
		Username:       input.Username,
		Type:           input.Type,
		PrimaryGroupID: input.PrimaryGroupID,
	}

	// Get user info for audit
	userID := uint(0) // In a real app, this would be from the auth middleware
	username := "admin" // In a real app, this would be from the auth middleware
	ipAddress := c.ClientIP()

	// Create account
	err := h.services.Account.CreateAccount(account, userID, username, ipAddress)
	if err != nil {
		h.logger.Errorf("Failed to create account: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// UpdateAccount handles PUT /api/accounts/:id
func (h *Handler) UpdateAccount(c *gin.Context) {
	h.logger.Infof("UpdateAccount: Received request")
	
	// Parse account ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Errorf("UpdateAccount: Invalid account ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}
	h.logger.Infof("UpdateAccount: Updating account with ID: %d", id)

	// Get existing account
	account, err := h.services.Account.GetAccount(uint(id))
	if err != nil {
		h.logger.Errorf("UpdateAccount: Failed to get account for update: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	h.logger.Infof("UpdateAccount: Found existing account: %+v", account)

	// Parse input
	var input accountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Errorf("UpdateAccount: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Infof("UpdateAccount: Input received: %+v", input)

	// Log old values for debugging
	h.logger.Infof("UpdateAccount: Old account values - UID: %d, Username: %s, Type: %s, PrimaryGroupID: %d",
		account.UID, account.Username, account.Type, account.PrimaryGroupID)

	// Update account fields
	account.UID = input.UID
	account.Username = input.Username
	account.Type = input.Type
	account.PrimaryGroupID = input.PrimaryGroupID

	// Log new values for debugging
	h.logger.Infof("UpdateAccount: New account values - UID: %d, Username: %s, Type: %s, PrimaryGroupID: %d",
		account.UID, account.Username, account.Type, account.PrimaryGroupID)

	// Get user info for audit
	userID := uint(0) // In a real app, this would be from the auth middleware
	username := "admin" // In a real app, this would be from the auth middleware
	ipAddress := c.ClientIP()

	// Update account
	err = h.services.Account.UpdateAccount(account, userID, username, ipAddress)
	if err != nil {
		h.logger.Errorf("UpdateAccount: Failed to update account: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Infof("UpdateAccount: Account updated successfully with ID: %d", account.ID)

	// Get the updated account to ensure we return the most current data
	updatedAccount, err := h.services.Account.GetAccount(uint(id))
	if err != nil {
		h.logger.Warnf("UpdateAccount: Failed to get updated account: %v", err)
		// Still return the local account object if we can't fetch the updated one
		c.JSON(http.StatusOK, account)
		return
	}
	
	h.logger.Infof("UpdateAccount: Returning updated account: %+v", updatedAccount)
	c.JSON(http.StatusOK, updatedAccount)
}

// DeleteAccount handles DELETE /api/accounts/:id
func (h *Handler) DeleteAccount(c *gin.Context) {
	// Parse account ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get user info for audit
	userID := uint(0) // In a real app, this would be from the auth middleware
	username := "admin" // In a real app, this would be from the auth middleware
	ipAddress := c.ClientIP()

	// Delete account
	err = h.services.Account.DeleteAccount(uint(id), userID, username, ipAddress)
	if err != nil {
		h.logger.Errorf("Failed to delete account: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// GetAccountGroups handles GET /api/accounts/:id/groups
func (h *Handler) GetAccountGroups(c *gin.Context) {
	// Parse account ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get groups
	groups, err := h.services.Account.GetAccountGroups(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get account groups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// SearchAccounts handles GET /api/search/accounts
func (h *Handler) SearchAccounts(c *gin.Context) {
	// Get search query
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Search accounts
	accounts, err := h.services.Account.SearchAccounts(query)
	if err != nil {
		h.logger.Errorf("Failed to search accounts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}