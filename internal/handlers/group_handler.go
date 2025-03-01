package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/home/unixify/internal/models"
)

// groupInput represents the input for group creation/update
type groupInput struct {
	GID         int              `json:"gid" binding:"required"`
	Groupname   string           `json:"groupname" binding:"required"`
	Description string           `json:"description"`
	Type        models.GroupType `json:"type" binding:"required"`
}

// GetAllGroups handles GET /api/groups
func (h *Handler) GetAllGroups(c *gin.Context) {
	// Get group type from query parameter (optional)
	groupType := models.GroupType(c.Query("type"))

	// Get all groups
	groups, err := h.services.Group.GetAllGroups(groupType)
	if err != nil {
		h.logger.Errorf("Failed to get groups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroup handles GET /api/groups/:id
func (h *Handler) GetGroup(c *gin.Context) {
	// Parse group ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get group
	group, err := h.services.Group.GetGroup(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get group: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupByGID handles GET /api/groups/gid/:gid
func (h *Handler) GetGroupByGID(c *gin.Context) {
	// Parse GID
	gid, err := strconv.Atoi(c.Param("gid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GID"})
		return
	}

	// Get group
	group, err := h.services.Group.GetGroupByGID(gid)
	if err != nil {
		h.logger.Errorf("Failed to get group by GID: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupByGroupname handles GET /api/groups/groupname/:groupname
func (h *Handler) GetGroupByGroupname(c *gin.Context) {
	// Get groupname
	groupname := c.Param("groupname")

	// Get group
	group, err := h.services.Group.GetGroupByGroupname(groupname)
	if err != nil {
		h.logger.Errorf("Failed to get group by groupname: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateGroup handles POST /api/groups
func (h *Handler) CreateGroup(c *gin.Context) {
	h.logger.Infof("CreateGroup: Received request")
	
	// Parse input
	var input groupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Errorf("CreateGroup: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("CreateGroup: Input received: %+v", input)

	// Create group
	group := &models.Group{
		GID:         input.GID,
		Groupname:   input.Groupname,
		Description: input.Description,
		Type:        input.Type,
	}

	h.logger.Infof("CreateGroup: Group object created: %+v", group)

	// Get user info for audit
	userID := uint(0) // In a real app, this would be from the auth middleware
	username := "admin" // In a real app, this would be from the auth middleware
	ipAddress := c.ClientIP()

	// Create group
	err := h.services.Group.CreateGroup(group, userID, username, ipAddress)
	if err != nil {
		h.logger.Errorf("CreateGroup: Failed to create group: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("CreateGroup: Group created successfully with ID: %d", group.ID)
	c.JSON(http.StatusCreated, group)
}

// UpdateGroup handles PUT /api/groups/:id
func (h *Handler) UpdateGroup(c *gin.Context) {
	// Parse group ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get existing group
	group, err := h.services.Group.GetGroup(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get group for update: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Parse input
	var input groupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update group fields
	group.GID = input.GID
	group.Groupname = input.Groupname
	group.Description = input.Description
	group.Type = input.Type

	// Get user info for audit
	userID := uint(0) // In a real app, this would be from the auth middleware
	username := "admin" // In a real app, this would be from the auth middleware
	ipAddress := c.ClientIP()

	// Update group
	err = h.services.Group.UpdateGroup(group, userID, username, ipAddress)
	if err != nil {
		h.logger.Errorf("Failed to update group: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteGroup handles DELETE /api/groups/:id
func (h *Handler) DeleteGroup(c *gin.Context) {
	// Parse group ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get user info for audit
	userID := uint(0) // In a real app, this would be from the auth middleware
	username := "admin" // In a real app, this would be from the auth middleware
	ipAddress := c.ClientIP()

	// Delete group
	err = h.services.Group.DeleteGroup(uint(id), userID, username, ipAddress)
	if err != nil {
		h.logger.Errorf("Failed to delete group: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

// GetGroupMembers handles GET /api/groups/:id/accounts
func (h *Handler) GetGroupMembers(c *gin.Context) {
	// Parse group ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get accounts
	accounts, err := h.services.Group.GetGroupMembers(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get group members: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get group members"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// SearchGroups handles GET /api/search/groups
func (h *Handler) SearchGroups(c *gin.Context) {
	// Get search query
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Search groups
	groups, err := h.services.Group.SearchGroups(query)
	if err != nil {
		h.logger.Errorf("Failed to search groups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}