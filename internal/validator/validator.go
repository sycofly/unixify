package validator

import (
	"fmt"

	"github.com/home/unixify/internal/models"
)

// UID range constants
const (
	MinUserUID       = 1000
	MaxUserUID       = 60000
	MinSystemUID     = 1
	MaxSystemUID     = 999
	MinServiceUID    = 60001
	MaxServiceUID    = 65535
	MinDatabaseUID   = 70000
	MaxDatabaseUID   = 79999
)

// GID range constants
const (
	MinUserGID       = 1000
	MaxUserGID       = 60000
	MinSystemGID     = 1
	MaxSystemGID     = 999
	MinServiceGID    = 60001
	MaxServiceGID    = 65535
	MinDatabaseGID   = 70000
	MaxDatabaseGID   = 79999
)

// ValidateUID validates a UID based on its type
func ValidateUID(uid int) error {
	if uid < 0 {
		return fmt.Errorf("UID cannot be negative")
	}

	// For now, just check that it's in a valid range for any type
	if uid > MaxDatabaseUID {
		return fmt.Errorf("UID %d is out of valid range", uid)
	}

	return nil
}

// ValidateUIDForType validates a UID based on its specific type
func ValidateUIDForType(uid int, accountType models.AccountType) error {
	if uid < 0 {
		return fmt.Errorf("UID cannot be negative")
	}

	// Check if UID is outside of recommended range and return a warning-type error
	switch accountType {
	case models.AccountTypePeople:
		if uid < MinUserUID || uid > MaxUserUID {
			return fmt.Errorf("WARNING: People account UID %d is outside recommended range (%d-%d)", uid, MinUserUID, MaxUserUID)
		}
	case models.AccountTypeSystem:
		if uid < MinSystemUID || uid > MaxSystemUID {
			return fmt.Errorf("WARNING: System account UID %d is outside recommended range (%d-%d)", uid, MinSystemUID, MaxSystemUID)
		}
	case models.AccountTypeService:
		if uid < MinServiceUID || uid > MaxServiceUID {
			return fmt.Errorf("WARNING: Service account UID %d is outside recommended range (%d-%d)", uid, MinServiceUID, MaxServiceUID)
		}
	case models.AccountTypeDatabase:
		if uid < MinDatabaseUID || uid > MaxDatabaseUID {
			return fmt.Errorf("WARNING: Database account UID %d is outside recommended range (%d-%d)", uid, MinDatabaseUID, MaxDatabaseUID)
		}
	default:
		return fmt.Errorf("invalid account type: %s", accountType)
	}

	return nil
}

// ValidateGID validates a GID based on its type
func ValidateGID(gid int) error {
	if gid < 0 {
		return fmt.Errorf("GID cannot be negative")
	}

	// For now, just check that it's in a valid range for any type
	if gid > MaxDatabaseGID {
		return fmt.Errorf("GID %d is out of valid range", gid)
	}

	return nil
}

// ValidateGIDForType validates a GID based on its specific type
func ValidateGIDForType(gid int, groupType models.GroupType) error {
	if gid < 0 {
		return fmt.Errorf("GID cannot be negative")
	}

	// Check if GID is outside of recommended range and return a warning-type error
	switch groupType {
	case models.GroupTypePeople:
		if gid < MinUserGID || gid > MaxUserGID {
			return fmt.Errorf("WARNING: People group GID %d is outside recommended range (%d-%d)", gid, MinUserGID, MaxUserGID)
		}
	case models.GroupTypeSystem:
		if gid < MinSystemGID || gid > MaxSystemGID {
			return fmt.Errorf("WARNING: System group GID %d is outside recommended range (%d-%d)", gid, MinSystemGID, MaxSystemGID)
		}
	case models.GroupTypeService:
		if gid < MinServiceGID || gid > MaxServiceGID {
			return fmt.Errorf("WARNING: Service group GID %d is outside recommended range (%d-%d)", gid, MinServiceGID, MaxServiceGID)
		}
	case models.GroupTypeDatabase:
		if gid < MinDatabaseGID || gid > MaxDatabaseGID {
			return fmt.Errorf("WARNING: Database group GID %d is outside recommended range (%d-%d)", gid, MinDatabaseGID, MaxDatabaseGID)
		}
	default:
		return fmt.Errorf("invalid group type: %s", groupType)
	}

	return nil
}

// IsValidAccountGroupAssignment checks if an account can be assigned to a group based on their types
func IsValidAccountGroupAssignment(accountType models.AccountType, groupType models.GroupType) bool {
	switch accountType {
	case models.AccountTypePeople:
		return groupType == models.GroupTypePeople || groupType == models.GroupTypeDatabase
	case models.AccountTypeSystem:
		return groupType == models.GroupTypeSystem
	case models.AccountTypeDatabase:
		return groupType == models.GroupTypeDatabase
	case models.AccountTypeService:
		return groupType == models.GroupTypeService
	default:
		return false
	}
}