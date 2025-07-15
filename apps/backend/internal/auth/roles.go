package auth

import (
	"errors"
	"strings"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleAgency Role = "agency"
	RoleAgent  Role = "agent"
	RoleOwner  Role = "owner"
	RoleBuyer  Role = "buyer"
)

// Permission represents system permissions
type Permission string

const (
	// Property permissions
	PermissionPropertyCreate Permission = "property:create"
	PermissionPropertyRead   Permission = "property:read"
	PermissionPropertyUpdate Permission = "property:update"
	PermissionPropertyDelete Permission = "property:delete"
	PermissionPropertyList   Permission = "property:list"
	
	// User permissions
	PermissionUserCreate Permission = "user:create"
	PermissionUserRead   Permission = "user:read"
	PermissionUserUpdate Permission = "user:update"
	PermissionUserDelete Permission = "user:delete"
	PermissionUserList   Permission = "user:list"
	
	// Agency permissions
	PermissionAgencyCreate Permission = "agency:create"
	PermissionAgencyRead   Permission = "agency:read"
	PermissionAgencyUpdate Permission = "agency:update"
	PermissionAgencyDelete Permission = "agency:delete"
	PermissionAgencyList   Permission = "agency:list"
	
	// Image permissions
	PermissionImageUpload Permission = "image:upload"
	PermissionImageRead   Permission = "image:read"
	PermissionImageUpdate Permission = "image:update"
	PermissionImageDelete Permission = "image:delete"
	
	// System permissions
	PermissionSystemAdmin     Permission = "system:admin"
	PermissionSystemMonitor   Permission = "system:monitor"
	PermissionSystemSecurity  Permission = "system:security"
	PermissionSystemAnalytics Permission = "system:analytics"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		// Admin has all permissions
		PermissionPropertyCreate, PermissionPropertyRead, PermissionPropertyUpdate, PermissionPropertyDelete, PermissionPropertyList,
		PermissionUserCreate, PermissionUserRead, PermissionUserUpdate, PermissionUserDelete, PermissionUserList,
		PermissionAgencyCreate, PermissionAgencyRead, PermissionAgencyUpdate, PermissionAgencyDelete, PermissionAgencyList,
		PermissionImageUpload, PermissionImageRead, PermissionImageUpdate, PermissionImageDelete,
		PermissionSystemAdmin, PermissionSystemMonitor, PermissionSystemSecurity, PermissionSystemAnalytics,
	},
	RoleAgency: {
		// Agency can manage their properties and agents
		PermissionPropertyCreate, PermissionPropertyRead, PermissionPropertyUpdate, PermissionPropertyDelete, PermissionPropertyList,
		PermissionUserCreate, PermissionUserRead, PermissionUserUpdate, PermissionUserList, // Can manage agents
		PermissionAgencyRead, PermissionAgencyUpdate, // Can update own agency
		PermissionImageUpload, PermissionImageRead, PermissionImageUpdate, PermissionImageDelete,
		PermissionSystemAnalytics, // Can view analytics
	},
	RoleAgent: {
		// Agent can manage properties for their agency
		PermissionPropertyCreate, PermissionPropertyRead, PermissionPropertyUpdate, PermissionPropertyList,
		PermissionUserRead, // Can view other users
		PermissionAgencyRead, // Can view agency info
		PermissionImageUpload, PermissionImageRead, PermissionImageUpdate, PermissionImageDelete,
	},
	RoleOwner: {
		// Owner can manage their own properties
		PermissionPropertyCreate, PermissionPropertyRead, PermissionPropertyUpdate, PermissionPropertyList,
		PermissionUserRead, // Can view agents/agencies
		PermissionAgencyRead, PermissionAgencyList,
		PermissionImageUpload, PermissionImageRead, PermissionImageUpdate, PermissionImageDelete,
	},
	RoleBuyer: {
		// Buyer has read-only access
		PermissionPropertyRead, PermissionPropertyList,
		PermissionUserRead, // Can view agents/agencies
		PermissionAgencyRead, PermissionAgencyList,
		PermissionImageRead,
	},
}

// ResourceContext represents the context for permission checking
type ResourceContext struct {
	UserID     string
	AgencyID   string
	ResourceID string // Property ID, User ID, etc.
	OwnerID    string // Owner of the resource
}

// AuthorizationManager handles role-based access control
type AuthorizationManager struct {
	rolePermissions map[Role][]Permission
}

// NewAuthorizationManager creates a new authorization manager
func NewAuthorizationManager() *AuthorizationManager {
	return &AuthorizationManager{
		rolePermissions: RolePermissions,
	}
}

// HasPermission checks if a role has a specific permission
func (am *AuthorizationManager) HasPermission(role Role, permission Permission) bool {
	permissions, exists := am.rolePermissions[role]
	if !exists {
		return false
	}
	
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	
	return false
}

// CanAccessResource checks if a user can access a specific resource
func (am *AuthorizationManager) CanAccessResource(userRole Role, permission Permission, context *ResourceContext) bool {
	// First check if role has the permission
	if !am.HasPermission(userRole, permission) {
		return false
	}
	
	// Admin can access everything
	if userRole == RoleAdmin {
		return true
	}
	
	// For other roles, check resource-specific rules
	switch permission {
	case PermissionPropertyUpdate, PermissionPropertyDelete:
		// Users can only modify their own properties or properties in their agency
		return am.canModifyProperty(userRole, context)
		
	case PermissionUserUpdate, PermissionUserDelete:
		// Users can modify themselves, agencies can modify their agents
		return am.canModifyUser(userRole, context)
		
	case PermissionAgencyUpdate, PermissionAgencyDelete:
		// Only agency owners can modify their agency
		return am.canModifyAgency(userRole, context)
		
	default:
		// For read operations and other permissions, role permission is sufficient
		return true
	}
}

// canModifyProperty checks if user can modify a property
func (am *AuthorizationManager) canModifyProperty(userRole Role, context *ResourceContext) bool {
	switch userRole {
	case RoleAgency, RoleAgent:
		// Agency and agents can modify properties in their agency
		return context.AgencyID != "" && context.AgencyID == getPropertyAgencyID(context.ResourceID)
		
	case RoleOwner:
		// Owner can modify their own properties
		return context.UserID == context.OwnerID
		
	default:
		return false
	}
}

// canModifyUser checks if user can modify another user
func (am *AuthorizationManager) canModifyUser(userRole Role, context *ResourceContext) bool {
	// Users can always modify themselves
	if context.UserID == context.ResourceID {
		return true
	}
	
	switch userRole {
	case RoleAgency:
		// Agency can modify agents in their agency
		return context.AgencyID != "" && context.AgencyID == getUserAgencyID(context.ResourceID)
		
	default:
		return false
	}
}

// canModifyAgency checks if user can modify an agency
func (am *AuthorizationManager) canModifyAgency(userRole Role, context *ResourceContext) bool {
	switch userRole {
	case RoleAgency:
		// Agency can modify their own agency
		return context.AgencyID == context.ResourceID
		
	default:
		return false
	}
}

// ValidateRole checks if a role string is valid
func ValidateRole(roleStr string) (Role, error) {
	role := Role(strings.ToLower(roleStr))
	
	switch role {
	case RoleAdmin, RoleAgency, RoleAgent, RoleOwner, RoleBuyer:
		return role, nil
	default:
		return "", errors.New("invalid role")
	}
}

// GetRoleHierarchy returns roles in hierarchical order (highest to lowest)
func GetRoleHierarchy() []Role {
	return []Role{RoleAdmin, RoleAgency, RoleAgent, RoleOwner, RoleBuyer}
}

// IsHigherRole checks if role1 has higher privileges than role2
func IsHigherRole(role1, role2 Role) bool {
	hierarchy := GetRoleHierarchy()
	
	role1Index := -1
	role2Index := -1
	
	for i, r := range hierarchy {
		if r == role1 {
			role1Index = i
		}
		if r == role2 {
			role2Index = i
		}
	}
	
	return role1Index != -1 && role2Index != -1 && role1Index < role2Index
}

// GetPermissionsForRole returns all permissions for a role
func (am *AuthorizationManager) GetPermissionsForRole(role Role) []Permission {
	permissions, exists := am.rolePermissions[role]
	if !exists {
		return []Permission{}
	}
	
	// Return a copy to prevent modification
	result := make([]Permission, len(permissions))
	copy(result, permissions)
	return result
}

// Helper functions (would be implemented with actual database queries)

// getPropertyAgencyID returns the agency ID for a property
func getPropertyAgencyID(propertyID string) string {
	// TODO: Implement database query to get property's agency
	// This is a placeholder
	return ""
}

// getUserAgencyID returns the agency ID for a user
func getUserAgencyID(userID string) string {
	// TODO: Implement database query to get user's agency
	// This is a placeholder
	return ""
}

// RoleMiddlewareConfig contains configuration for role-based middleware
type RoleMiddlewareConfig struct {
	RequiredRole       Role
	RequiredPermission Permission
	AllowOwnerAccess   bool // Allow resource owner access even without role permission
}

// CheckAccess is a helper function for middleware to check access
func (am *AuthorizationManager) CheckAccess(userRole Role, userID, agencyID string, config RoleMiddlewareConfig, resourceID string) bool {
	context := &ResourceContext{
		UserID:     userID,
		AgencyID:   agencyID,
		ResourceID: resourceID,
	}
	
	// Check if user has required role
	if config.RequiredRole != "" && userRole != config.RequiredRole {
		// Check if user role is higher in hierarchy
		if !IsHigherRole(userRole, config.RequiredRole) {
			return false
		}
	}
	
	// Check if user has required permission
	if config.RequiredPermission != "" {
		return am.CanAccessResource(userRole, config.RequiredPermission, context)
	}
	
	return true
}