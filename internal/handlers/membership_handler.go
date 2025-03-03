package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// membershipInput represents the input for assigning/removing accounts to/from groups
type membershipInput struct {
	AccountID uint `json:"account_id" binding:"required"`
	GroupID   uint `json:"group_id" binding:"required"`
}

// AssignAccountToGroup handles POST /api/memberships
func (h *Handler) AssignAccountToGroup(c *gin.Context) {
	// Parse input
	var input membershipInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info for audit
	userID := uint(0) // In a real app, this would be from the auth middleware
	username := "admin" // In a real app, this would be from the auth middleware
	ipAddress := c.ClientIP()

	// Assign account to group
	err := h.services.Account.AssignAccountToGroup(input.AccountID, input.GroupID, userID, username, ipAddress)
	if err != nil {
		h.logger.Errorf("Failed to assign account to group: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account assigned to group successfully"})
}

// RemoveAccountFromGroup handles DELETE /api/memberships
func (h *Handler) RemoveAccountFromGroup(c *gin.Context) {
	// Parse input
	var input membershipInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info for audit
	userID := uint(0) // In a real app, this would be from the auth middleware
	username := "admin" // In a real app, this would be from the auth middleware
	ipAddress := c.ClientIP()

	// Remove account from group
	err := h.services.Account.RemoveAccountFromGroup(input.AccountID, input.GroupID, userID, username, ipAddress)
	if err != nil {
		h.logger.Errorf("Failed to remove account from group: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account removed from group successfully"})
}