package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserRole represents the different types of users in the system
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleAgent  UserRole = "agent"
	RoleBuyer  UserRole = "buyer"
	RoleOwner  UserRole = "seller"  // El schema DB tiene "seller" no "owner"
	RoleAgency UserRole = "agency"  // Mantener para compatibilidad
)

// UserStatus represents the user account status
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
	StatusPending   UserStatus = "pending"
)

// User represents a system user with role-based access
type User struct {
	ID                      string     `json:"id" db:"id"`
	FirstName               string     `json:"first_name" db:"first_name"`
	LastName                string     `json:"last_name" db:"last_name"`
	Email                   string     `json:"email" db:"email"`
	Phone                   *string    `json:"phone" db:"phone"`
	Cedula                  *string    `json:"cedula" db:"national_id"`
	DateOfBirth             *time.Time `json:"date_of_birth" db:"date_of_birth"`
	Role                    UserRole   `json:"role" db:"user_type"`
	Active                  bool       `json:"active" db:"active"`
	MinBudget               *float64   `json:"min_budget" db:"min_budget"`
	MaxBudget               *float64   `json:"max_budget" db:"max_budget"`
	PreferredProvinces      []string   `json:"preferred_provinces" db:"preferred_provinces"`
	PreferredPropertyTypes  []string   `json:"preferred_property_types" db:"preferred_property_types"`
	AvatarURL               *string    `json:"avatar_url" db:"avatar_url"`
	Bio                     *string    `json:"bio" db:"bio"`
	RealEstateCompanyID     *string    `json:"real_estate_company_id" db:"real_estate_company_id"`
	ReceiveNotifications    bool       `json:"receive_notifications" db:"receive_notifications"`
	ReceiveNewsletter       bool       `json:"receive_newsletter" db:"receive_newsletter"`
	AgencyID                *string    `json:"agency_id" db:"agency_id"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
	
	// Additional fields for auth functionality (not in DB)
	PasswordHash            string     `json:"-"`
	EmailVerified           bool       `json:"email_verified"`
	EmailVerificationToken  *string    `json:"-"`
	PasswordResetToken      *string    `json:"-"`
	PasswordResetExpires    *time.Time `json:"-"`
	LastLogin               *time.Time `json:"last_login"`
	LastLoginAt             *time.Time `json:"last_login_at"`
	DeletedAt               *time.Time `json:"deleted_at"`
	Status                  UserStatus `json:"status"`
}

// NewUser creates a new user with validation
func NewUser(email, firstName, lastName string, role UserRole) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validateName(firstName); err != nil {
		return nil, err
	}

	if err := validateName(lastName); err != nil {
		return nil, err
	}

	if err := validateRole(role); err != nil {
		return nil, err
	}

	user := &User{
		ID:                   uuid.New().String(),
		Email:                strings.ToLower(strings.TrimSpace(email)),
		FirstName:            strings.TrimSpace(firstName),
		LastName:             strings.TrimSpace(lastName),
		Role:                 role,
		Status:               StatusPending,
		Active:               true,
		ReceiveNotifications: true,
		ReceiveNewsletter:    false,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	return user, nil
}

// IsValid validates the user data
func (u *User) IsValid() error {
	if u.ID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	if err := validateEmail(u.Email); err != nil {
		return err
	}

	if err := validateName(u.FirstName); err != nil {
		return err
	}

	if err := validateName(u.LastName); err != nil {
		return err
	}

	if err := validateRole(u.Role); err != nil {
		return err
	}

	if err := validateStatus(u.Status); err != nil {
		return err
	}

	// Agent must belong to an agency
	if u.Role == RoleAgent && u.AgencyID == nil {
		return fmt.Errorf("agent must belong to an agency")
	}

	// Non-agents cannot have agency association
	if u.Role != RoleAgent && u.AgencyID != nil {
		return fmt.Errorf("only agents can belong to an agency")
	}

	return nil
}

// SetAgency assigns a user to an agency (only for agents)
func (u *User) SetAgency(agencyID string) error {
	if u.Role != RoleAgent {
		return fmt.Errorf("only agents can be assigned to agencies")
	}

	if agencyID == "" {
		return fmt.Errorf("agency ID cannot be empty")
	}

	u.AgencyID = &agencyID
	u.UpdatedAt = time.Now()
	return nil
}

// Activate sets the user status to active
func (u *User) Activate() error {
	if u.Status == StatusSuspended {
		return fmt.Errorf("suspended users cannot be activated directly")
	}

	u.Status = StatusActive
	u.UpdatedAt = time.Now()
	return nil
}

// Deactivate sets the user status to inactive
func (u *User) Deactivate() error {
	u.Status = StatusInactive
	u.UpdatedAt = time.Now()
	return nil
}

// Suspend sets the user status to suspended
func (u *User) Suspend() error {
	u.Status = StatusSuspended
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// CanManageProperty checks if user can manage a specific property
func (u *User) CanManageProperty(propertyOwnerID string, propertyAgencyID *string) bool {
	switch u.Role {
	case RoleAdmin:
		return true
	case RoleOwner:
		return u.ID == propertyOwnerID
	case RoleAgent:
		return u.AgencyID != nil && propertyAgencyID != nil && *u.AgencyID == *propertyAgencyID
	case RoleAgency:
		return propertyAgencyID != nil && u.ID == *propertyAgencyID
	default:
		return false
	}
}

// CanCreateProperty checks if user can create properties
func (u *User) CanCreateProperty() bool {
	return u.Role == RoleAdmin || u.Role == RoleOwner || u.Role == RoleAgent || u.Role == RoleAgency
}

// CanViewProperty checks if user can view a specific property
func (u *User) CanViewProperty(propertyOwnerID string, propertyAgencyID *string) bool {
	// All users can view all properties in a marketplace
	return true
}

// GetRoleLevel returns the hierarchical level of the role
func (u *User) GetRoleLevel() int {
	switch u.Role {
	case RoleAdmin:
		return 5
	case RoleAgency:
		return 4
	case RoleAgent:
		return 3
	case RoleOwner:
		return 2
	case RoleBuyer:
		return 1
	default:
		return 0
	}
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsAgent checks if the user is an agent
func (u *User) IsAgent() bool {
	return u.Role == RoleAgent
}

// IsOwner checks if the user is an owner
func (u *User) IsOwner() bool {
	return u.Role == RoleOwner
}

// IsAgency checks if the user is an agency
func (u *User) IsAgency() bool {
	return u.Role == RoleAgency
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsBuyer checks if the user is a buyer
func (u *User) IsBuyer() bool {
	return u.Role == RoleBuyer
}

// Helper validation functions
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}

	if len(name) > 255 {
		return fmt.Errorf("name cannot exceed 255 characters")
	}

	return nil
}

func validateRole(role UserRole) error {
	validRoles := []UserRole{RoleAdmin, RoleAgency, RoleAgent, RoleOwner, RoleBuyer}
	
	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}

	return fmt.Errorf("invalid role: %s", role)
}

func validateStatus(status UserStatus) error {
	validStatuses := []UserStatus{StatusActive, StatusInactive, StatusSuspended, StatusPending}
	
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status: %s", status)
}

// Role represents the user role enum type for repository use
type Role string

const (
	RoleAdminStr  Role = "admin"
	RoleAgencyStr Role = "agency"
	RoleAgentStr  Role = "agent"
	RoleOwnerStr  Role = "owner"
	RoleBuyerStr  Role = "buyer"
)

// UserSearchParams represents search parameters for users
type UserSearchParams struct {
	Query      string              `json:"query,omitempty"`
	Role       *UserRole           `json:"role,omitempty"`
	Status     *UserStatus         `json:"status,omitempty"`
	Active     *bool               `json:"active,omitempty"`
	Province   string              `json:"province,omitempty"`
	Provinces  []string            `json:"provinces,omitempty"`
	City       string              `json:"city,omitempty"`
	AgencyID   *string             `json:"agency_id,omitempty"`
	MinBudget  *float64            `json:"min_budget,omitempty"`
	MaxBudget  *float64            `json:"max_budget,omitempty"`
	Pagination *PaginationParams   `json:"pagination,omitempty"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	SortBy     string              `json:"sort_by"`
	SortDesc   bool                `json:"sort_desc"`
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers         int                    `json:"total_users"`
	ActiveUsers        int                    `json:"active_users"`
	InactiveUsers      int                    `json:"inactive_users"`
	SuspendedUsers     int                    `json:"suspended_users"`
	PendingUsers       int                    `json:"pending_users"`
	AdminCount         int                    `json:"admin_count"`
	AgencyCount        int                    `json:"agency_count"`
	AgentCount         int                    `json:"agent_count"`
	OwnerCount         int                    `json:"owner_count"`
	BuyerCount         int                    `json:"buyer_count"`
	EmailVerified      int                    `json:"email_verified"`
	WithBudget         int                    `json:"with_budget"`
	AssociatedAgents   int                    `json:"associated_agents"`
	UsersByRole        map[string]int         `json:"users_by_role"`
	UsersByProvince    map[string]int         `json:"users_by_province"`
	NewUsersThisMonth  int                    `json:"new_users_this_month"`
	AverageAge         float64                `json:"average_age"`
	GenderDistribution map[string]int         `json:"gender_distribution"`
}

// Business methods for User

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// Name returns the full name (alias for GetFullName)
func (u *User) Name() string {
	return u.GetFullName()
}