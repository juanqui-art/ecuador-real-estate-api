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
	RoleAgency UserRole = "agency"
	RoleAgent  UserRole = "agent"
	RoleOwner  UserRole = "owner"
	RoleBuyer  UserRole = "buyer"
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
	ID          string     `json:"id" db:"id"`
	Email       string     `json:"email" db:"email"`
	Name        string     `json:"name" db:"name"`
	Phone       *string    `json:"phone" db:"phone"`
	Role        UserRole   `json:"role" db:"role"`
	Status      UserStatus `json:"status" db:"status"`
	AgencyID    *string    `json:"agency_id" db:"agency_id"`
	Avatar      *string    `json:"avatar" db:"avatar"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}

// NewUser creates a new user with validation
func NewUser(email, name string, role UserRole) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validateName(name); err != nil {
		return nil, err
	}

	if err := validateRole(role); err != nil {
		return nil, err
	}

	user := &User{
		ID:        uuid.New().String(),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Name:      strings.TrimSpace(name),
		Role:      role,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

	if err := validateName(u.Name); err != nil {
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