package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAuditEntries handles GET /api/audit
func (h *Handler) GetAuditEntries(c *gin.Context) {
	// Get filter parameters
	entityType := c.Query("entity_type")
	action := c.Query("action")
	
	var entityID, userID uint
	entityIDStr := c.Query("entity_id")
	if entityIDStr != "" {
		id, err := strconv.ParseUint(entityIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
			return
		}
		entityID = uint(id)
	}
	
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userID = uint(id)
	}

	// Get audit entries
	entries, err := h.services.Audit.GetAuditEntries(entityType, action, entityID, userID)
	if err != nil {
		h.logger.Errorf("Failed to get audit entries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get audit entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// GetAuditEntry handles GET /api/audit/:id
func (h *Handler) GetAuditEntry(c *gin.Context) {
	// Parse audit entry ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit entry ID"})
		return
	}

	// Get audit entry
	entry, err := h.services.Audit.GetAuditEntry(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get audit entry: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Audit entry not found"})
		return
	}

	c.JSON(http.StatusOK, entry)
}