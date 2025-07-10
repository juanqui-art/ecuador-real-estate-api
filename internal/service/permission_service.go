package service

import (
	"fmt"
	"log"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// PermissionService handles role-based access control
type PermissionService struct {
	userRepo     *repository.UserRepository
	agencyRepo   *repository.AgencyRepository
	propertyRepo *repository.PropertyRepository
	logger       *log.Logger
}

// NewPermissionService creates a new permission service
func NewPermissionService(
	userRepo *repository.UserRepository,
	agencyRepo *repository.AgencyRepository,
	propertyRepo *repository.PropertyRepository,
	logger *log.Logger,
) *PermissionService {
	return &PermissionService{
		userRepo:     userRepo,
		agencyRepo:   agencyRepo,
		propertyRepo: propertyRepo,
		logger:       logger,
	}
}

// Permission represents a specific permission
type Permission string

const (
	// User permissions
	PermissionCreateUser    Permission = "create_user"
	PermissionViewUser      Permission = "view_user"
	PermissionUpdateUser    Permission = "update_user"
	PermissionDeleteUser    Permission = "delete_user"
	PermissionViewAllUsers  Permission = "view_all_users"
	PermissionManageUsers   Permission = "manage_users"

	// Agency permissions
	PermissionCreateAgency    Permission = "create_agency"
	PermissionViewAgency      Permission = "view_agency"
	PermissionUpdateAgency    Permission = "update_agency"
	PermissionDeleteAgency    Permission = "delete_agency"
	PermissionViewAllAgencies Permission = "view_all_agencies"
	PermissionManageAgencies  Permission = "manage_agencies"

	// Property permissions
	PermissionCreateProperty    Permission = "create_property"
	PermissionViewProperty      Permission = "view_property"
	PermissionUpdateProperty    Permission = "update_property"
	PermissionDeleteProperty    Permission = "delete_property"
	PermissionViewAllProperties Permission = "view_all_properties"
	PermissionManageProperties  Permission = "manage_properties"
	PermissionAssignProperty    Permission = "assign_property"
	PermissionTransferProperty  Permission = "transfer_property"

	// System permissions
	PermissionViewStats     Permission = "view_stats"
	PermissionSystemAdmin   Permission = "system_admin"
	PermissionViewReports   Permission = "view_reports"
	PermissionManageSystem  Permission = "manage_system"
)

// GetUserPermissions returns all permissions for a user based on their role
func (s *PermissionService) GetUserPermissions(userID string) ([]Permission, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.Active {
		return []Permission{}, nil
	}

	return s.GetRolePermissions(user.Role), nil
}

// GetRolePermissions returns all permissions for a specific role
func (s *PermissionService) GetRolePermissions(role domain.Role) []Permission {
	switch role {
	case domain.RoleAdmin:
		return []Permission{
			// All permissions for admin
			PermissionCreateUser, PermissionViewUser, PermissionUpdateUser, PermissionDeleteUser,
			PermissionViewAllUsers, PermissionManageUsers,
			PermissionCreateAgency, PermissionViewAgency, PermissionUpdateAgency, PermissionDeleteAgency,
			PermissionViewAllAgencies, PermissionManageAgencies,
			PermissionCreateProperty, PermissionViewProperty, PermissionUpdateProperty, PermissionDeleteProperty,
			PermissionViewAllProperties, PermissionManageProperties, PermissionAssignProperty, PermissionTransferProperty,
			PermissionViewStats, PermissionSystemAdmin, PermissionViewReports, PermissionManageSystem,
		}
	case domain.RoleAgency:
		return []Permission{
			// Agency can manage their agents and properties
			PermissionCreateUser, PermissionViewUser, PermissionUpdateUser, // For their agents
			PermissionViewAgency, PermissionUpdateAgency, // For their own agency
			PermissionCreateProperty, PermissionViewProperty, PermissionUpdateProperty, PermissionDeleteProperty,
			PermissionAssignProperty, PermissionViewStats,
		}
	case domain.RoleAgent:
		return []Permission{
			// Agent can manage assigned properties
			PermissionViewUser, // Limited to their agency
			PermissionViewAgency, // Their own agency
			PermissionCreateProperty, PermissionViewProperty, PermissionUpdateProperty,
			PermissionAssignProperty, // Limited scope
		}
	case domain.RoleOwner:
		return []Permission{
			// Owner can manage their own properties
			PermissionViewUser, // Limited to themselves
			PermissionUpdateUser, // Limited to themselves
			PermissionCreateProperty, PermissionViewProperty, PermissionUpdateProperty, PermissionDeleteProperty,
		}
	case domain.RoleBuyer:
		return []Permission{
			// Buyer can view properties and update their profile
			PermissionViewUser, // Limited to themselves
			PermissionUpdateUser, // Limited to themselves
			PermissionViewProperty, // Available properties only
		}
	default:
		return []Permission{}
	}
}

// HasPermission checks if a user has a specific permission
func (s *PermissionService) HasPermission(userID string, permission Permission) (bool, error) {
	permissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}

	return false, nil
}

// CanManageUser checks if a user can manage another user
func (s *PermissionService) CanManageUser(managerID, targetUserID string) (bool, error) {
	manager, err := s.userRepo.GetByID(managerID)
	if err != nil {
		return false, fmt.Errorf("failed to get manager: %w", err)
	}

	targetUser, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get target user: %w", err)
	}

	if !manager.Active {
		return false, nil
	}

	switch manager.Role {
	case domain.RoleAdmin:
		return true, nil
	case domain.RoleAgency:
		// Agency can manage their agents and themselves
		return manager.ID == targetUser.ID ||
			(targetUser.Role == domain.RoleAgent && targetUser.AgencyID != nil && *targetUser.AgencyID == manager.ID), nil
	case domain.RoleAgent:
		// Agent can only manage themselves
		return manager.ID == targetUser.ID, nil
	case domain.RoleOwner:
		// Owner can only manage themselves
		return manager.ID == targetUser.ID, nil
	case domain.RoleBuyer:
		// Buyer can only manage themselves
		return manager.ID == targetUser.ID, nil
	default:
		return false, nil
	}
}

// CanManageAgency checks if a user can manage an agency
func (s *PermissionService) CanManageAgency(userID, agencyID string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.Active {
		return false, nil
	}

	switch user.Role {
	case domain.RoleAdmin:
		return true, nil
	case domain.RoleAgency:
		// Agency can manage their own agency
		return user.ID == agencyID, nil
	default:
		return false, nil
	}
}

// CanManageProperty checks if a user can manage a property
func (s *PermissionService) CanManageProperty(userID, propertyID string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	property, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		return false, fmt.Errorf("failed to get property: %w", err)
	}

	if !user.Active {
		return false, nil
	}

	return user.CanManageProperty(property), nil
}

// CanViewProperty checks if a user can view a property
func (s *PermissionService) CanViewProperty(userID, propertyID string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	property, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		return false, fmt.Errorf("failed to get property: %w", err)
	}

	if !user.Active {
		return false, nil
	}

	return user.CanViewProperty(property), nil
}

// CanCreateProperty checks if a user can create properties
func (s *PermissionService) CanCreateProperty(userID string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.Active {
		return false, nil
	}

	switch user.Role {
	case domain.RoleAdmin, domain.RoleAgency, domain.RoleAgent, domain.RoleOwner:
		return true, nil
	default:
		return false, nil
	}
}

// CanAssignProperty checks if a user can assign properties to agents
func (s *PermissionService) CanAssignProperty(userID, propertyID, agentID string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	property, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		return false, fmt.Errorf("failed to get property: %w", err)
	}

	agent, err := s.userRepo.GetByID(agentID)
	if err != nil {
		return false, fmt.Errorf("failed to get agent: %w", err)
	}

	if !user.Active || !agent.Active {
		return false, nil
	}

	if agent.Role != domain.RoleAgent {
		return false, fmt.Errorf("target user is not an agent")
	}

	switch user.Role {
	case domain.RoleAdmin:
		return true, nil
	case domain.RoleAgency:
		// Agency can assign properties to their agents
		return agent.AgencyID != nil && *agent.AgencyID == user.ID, nil
	case domain.RoleAgent:
		// Agent can only assign properties within their agency
		return user.AgencyID != nil && agent.AgencyID != nil && *user.AgencyID == *agent.AgencyID &&
			user.CanManageProperty(property), nil
	case domain.RoleOwner:
		// Owner can assign their properties to any agent
		return user.CanManageProperty(property), nil
	default:
		return false, nil
	}
}

// CanTransferProperty checks if a user can transfer property ownership
func (s *PermissionService) CanTransferProperty(userID, propertyID, newOwnerID string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	property, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		return false, fmt.Errorf("failed to get property: %w", err)
	}

	newOwner, err := s.userRepo.GetByID(newOwnerID)
	if err != nil {
		return false, fmt.Errorf("failed to get new owner: %w", err)
	}

	if !user.Active || !newOwner.Active {
		return false, nil
	}

	if newOwner.Role != domain.RoleOwner && newOwner.Role != domain.RoleBuyer {
		return false, fmt.Errorf("new owner must have owner or buyer role")
	}

	switch user.Role {
	case domain.RoleAdmin:
		return true, nil
	case domain.RoleOwner:
		// Owner can transfer their own properties
		return user.CanManageProperty(property), nil
	default:
		return false, nil
	}
}

// FilterPropertiesForUser filters properties based on user permissions
func (s *PermissionService) FilterPropertiesForUser(userID string, properties []*domain.Property) ([]*domain.Property, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.Active {
		return []*domain.Property{}, nil
	}

	var filteredProperties []*domain.Property

	for _, property := range properties {
		if user.CanViewProperty(property) {
			filteredProperties = append(filteredProperties, property)
		}
	}

	return filteredProperties, nil
}

// GetUserAccessibleProperties returns properties accessible to a user
func (s *PermissionService) GetUserAccessibleProperties(userID string, filters *domain.PropertySearchFilters) ([]*domain.Property, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.Active {
		return []*domain.Property{}, nil
	}

	// Apply user-specific filters based on role
	switch user.Role {
	case domain.RoleAdmin:
		// Admin can access all properties
		break
	case domain.RoleAgency:
		// Agency can access their properties
		filters.AgencyID = &user.ID
	case domain.RoleAgent:
		// Agent can access their assigned properties or agency properties
		if user.AgencyID != nil {
			filters.AgencyID = user.AgencyID
		}
	case domain.RoleOwner:
		// Owner can access their properties
		filters.OwnerID = &user.ID
	case domain.RoleBuyer:
		// Buyer can only access available properties
		filters.Status = []string{domain.StatusAvailable}
	}

	// Get filtered properties (this would need to be implemented in the property repository)
	// For now, we'll return an empty slice
	return []*domain.Property{}, nil
}

// ValidatePropertyCreation validates if a user can create a property with specific attributes
func (s *PermissionService) ValidatePropertyCreation(userID string, property *domain.Property) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.Active {
		return fmt.Errorf("user is not active")
	}

	// Check if user can create properties
	canCreate, err := s.CanCreateProperty(userID)
	if err != nil {
		return fmt.Errorf("failed to check create permission: %w", err)
	}

	if !canCreate {
		return fmt.Errorf("user does not have permission to create properties")
	}

	// Validate property assignments based on user role
	switch user.Role {
	case domain.RoleAdmin:
		// Admin can create any property
		break
	case domain.RoleAgency:
		// Agency must assign itself as the agency
		if property.AgencyID == nil || *property.AgencyID != user.ID {
			return fmt.Errorf("agency must assign itself as the property agency")
		}
	case domain.RoleAgent:
		// Agent must assign their agency
		if user.AgencyID == nil {
			return fmt.Errorf("agent must be associated with an agency")
		}
		if property.AgencyID == nil || *property.AgencyID != *user.AgencyID {
			return fmt.Errorf("agent must assign their agency to the property")
		}
	case domain.RoleOwner:
		// Owner must assign themselves as the owner
		if property.OwnerID == nil || *property.OwnerID != user.ID {
			return fmt.Errorf("owner must assign themselves as the property owner")
		}
	}

	return nil
}

// ValidatePropertyUpdate validates if a user can update a property with specific changes
func (s *PermissionService) ValidatePropertyUpdate(userID string, oldProperty, newProperty *domain.Property) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.Active {
		return fmt.Errorf("user is not active")
	}

	// Check if user can manage the property
	if !user.CanManageProperty(oldProperty) {
		return fmt.Errorf("user does not have permission to update this property")
	}

	// Validate ownership/agency changes
	if oldProperty.OwnerID != newProperty.OwnerID {
		canTransfer, err := s.CanTransferProperty(userID, oldProperty.ID, *newProperty.OwnerID)
		if err != nil {
			return fmt.Errorf("failed to check transfer permission: %w", err)
		}
		if !canTransfer {
			return fmt.Errorf("user does not have permission to transfer property ownership")
		}
	}

	if oldProperty.AgentID != newProperty.AgentID {
		if newProperty.AgentID != nil {
			canAssign, err := s.CanAssignProperty(userID, oldProperty.ID, *newProperty.AgentID)
			if err != nil {
				return fmt.Errorf("failed to check assign permission: %w", err)
			}
			if !canAssign {
				return fmt.Errorf("user does not have permission to assign this agent")
			}
		}
	}

	return nil
}

// GetUserDashboardData returns dashboard data filtered by user permissions
func (s *PermissionService) GetUserDashboardData(userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.Active {
		return nil, fmt.Errorf("user is not active")
	}

	dashboardData := make(map[string]interface{})
	dashboardData["user"] = user
	dashboardData["role"] = user.Role
	dashboardData["permissions"] = s.GetRolePermissions(user.Role)

	switch user.Role {
	case domain.RoleAdmin:
		// Admin gets system-wide statistics
		userStats, _ := s.userRepo.GetStatistics()
		agencyStats, _ := s.agencyRepo.GetStatistics()
		dashboardData["user_stats"] = userStats
		dashboardData["agency_stats"] = agencyStats
	case domain.RoleAgency:
		// Agency gets their performance and agents
		if performance, err := s.agencyRepo.GetPerformance(user.ID); err == nil {
			dashboardData["performance"] = performance
		}
		if agents, err := s.userRepo.GetByAgency(user.ID); err == nil {
			dashboardData["agents"] = agents
		}
	case domain.RoleAgent:
		// Agent gets their assigned properties and agency info
		if user.AgencyID != nil {
			if agency, err := s.agencyRepo.GetByID(*user.AgencyID); err == nil {
				dashboardData["agency"] = agency
			}
		}
	case domain.RoleOwner:
		// Owner gets their properties statistics
		// This would need to be implemented in the property repository
	case domain.RoleBuyer:
		// Buyer gets their preferences and matching properties
		dashboardData["preferences"] = map[string]interface{}{
			"min_budget":            user.MinBudget,
			"max_budget":            user.MaxBudget,
			"interested_provinces":  user.InterestedProvinces,
			"interested_types":      user.InterestedTypes,
		}
	}

	return dashboardData, nil
}

// LogPermissionCheck logs permission checks for auditing
func (s *PermissionService) LogPermissionCheck(userID string, permission Permission, resource string, granted bool) {
	s.logger.Printf("Permission Check - User: %s, Permission: %s, Resource: %s, Granted: %t",
		userID, permission, resource, granted)
}