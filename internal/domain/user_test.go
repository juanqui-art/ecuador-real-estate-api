package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		firstName   string
		lastName    string
		role        UserRole
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid owner user",
			email:       "juan@example.com",
			firstName:   "Juan",
			lastName:    "Pérez",
			role:        RoleOwner,
			expectError: false,
		},
		{
			name:        "valid agent user",
			email:       "agent@realty.com",
			firstName:   "María",
			lastName:    "López",
			role:        RoleAgent,
			expectError: false,
		},
		{
			name:        "valid buyer user",
			email:       "buyer@gmail.com",
			firstName:   "Carlos",
			lastName:    "Ruiz",
			role:        RoleBuyer,
			expectError: false,
		},
		{
			name:        "invalid email format",
			email:       "invalid-email",
			firstName:   "Test",
			lastName:    "User",
			role:        RoleOwner,
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "empty email",
			email:       "",
			firstName:   "Test",
			lastName:    "User",
			role:        RoleOwner,
			expectError: true,
			errorMsg:    "email cannot be empty",
		},
		{
			name:        "empty first name",
			email:       "test@example.com",
			firstName:   "",
			lastName:    "User",
			role:        RoleOwner,
			expectError: true,
			errorMsg:    "name cannot be empty",
		},
		{
			name:        "empty last name",
			email:       "test@example.com",
			firstName:   "Test",
			lastName:    "",
			role:        RoleOwner,
			expectError: true,
			errorMsg:    "name cannot be empty",
		},
		{
			name:        "short first name",
			email:       "test@example.com",
			firstName:   "A",
			lastName:    "User",
			role:        RoleOwner,
			expectError: true,
			errorMsg:    "name must be at least 2 characters long",
		},
		{
			name:        "invalid role",
			email:       "test@example.com",
			firstName:   "Test",
			lastName:    "User",
			role:        UserRole("invalid"),
			expectError: true,
			errorMsg:    "invalid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.firstName, tt.lastName, tt.role)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.NotEmpty(t, user.ID)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, tt.role, user.Role)
				assert.Equal(t, StatusPending, user.Status)
				assert.NotZero(t, user.CreatedAt)
				assert.NotZero(t, user.UpdatedAt)
				assert.Nil(t, user.DeletedAt)
			}
		})
	}
}

func TestUserIsValid(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid owner user",
			user: &User{
				ID:        uuid.New().String(),
				Email:     "owner@example.com",
				FirstName: "Property",
				LastName:  "Owner",
				Role:      RoleOwner,
				Status:    StatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
		},
		{
			name: "valid agent with agency",
			user: &User{
				ID:        uuid.New().String(),
				Email:     "agent@realty.com",
				FirstName: "Real Estate",
				LastName:  "Agent",
				Role:      RoleAgent,
				Status:    StatusActive,
				AgencyID:  stringPtr(uuid.New().String()),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
		},
		{
			name: "agent without agency",
			user: &User{
				ID:        uuid.New().String(),
				Email:     "agent@realty.com",
				FirstName: "Real Estate",
				LastName:  "Agent",
				Role:      RoleAgent,
				Status:    StatusActive,
				AgencyID:  nil,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: true,
			errorMsg:    "agent must belong to an agency",
		},
		{
			name: "owner with agency",
			user: &User{
				ID:        uuid.New().String(),
				Email:     "owner@example.com",
				FirstName: "Property",
				LastName:  "Owner",
				Role:      RoleOwner,
				Status:    StatusActive,
				AgencyID:  stringPtr(uuid.New().String()),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: true,
			errorMsg:    "only agents can belong to an agency",
		},
		{
			name: "user with empty ID",
			user: &User{
				ID:        "",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				Role:      RoleOwner,
				Status:    StatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: true,
			errorMsg:    "user ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.IsValid()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserSetAgency(t *testing.T) {
	agencyID := uuid.New().String()

	tests := []struct {
		name        string
		user        *User
		agencyID    string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid agent assignment",
			user: &User{
				ID:   uuid.New().String(),
				Role: RoleAgent,
			},
			agencyID:    agencyID,
			expectError: false,
		},
		{
			name: "owner cannot be assigned to agency",
			user: &User{
				ID:   uuid.New().String(),
				Role: RoleOwner,
			},
			agencyID:    agencyID,
			expectError: true,
			errorMsg:    "only agents can be assigned to agencies",
		},
		{
			name: "empty agency ID",
			user: &User{
				ID:   uuid.New().String(),
				Role: RoleAgent,
			},
			agencyID:    "",
			expectError: true,
			errorMsg:    "agency ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.SetAgency(tt.agencyID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.agencyID, *tt.user.AgencyID)
			}
		})
	}
}

func TestUserStatusManagement(t *testing.T) {
	user := &User{
		ID:     uuid.New().String(),
		Status: StatusPending,
	}

	// Test activation
	err := user.Activate()
	assert.NoError(t, err)
	assert.Equal(t, StatusActive, user.Status)

	// Test deactivation
	err = user.Deactivate()
	assert.NoError(t, err)
	assert.Equal(t, StatusInactive, user.Status)

	// Test suspension
	err = user.Suspend()
	assert.NoError(t, err)
	assert.Equal(t, StatusSuspended, user.Status)

	// Test suspended user cannot be activated directly
	err = user.Activate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "suspended users cannot be activated directly")
}

func TestUserPermissions(t *testing.T) {
	ownerID := uuid.New().String()
	agencyID := uuid.New().String()
	agentID := uuid.New().String()

	tests := []struct {
		name              string
		user              *User
		propertyOwnerID   string
		propertyAgencyID  *string
		canManage         bool
		canCreate         bool
		canView           bool
	}{
		{
			name: "admin can manage any property",
			user: &User{
				ID:   uuid.New().String(),
				Role: RoleAdmin,
			},
			propertyOwnerID:   ownerID,
			propertyAgencyID:  nil,
			canManage:         true,
			canCreate:         true,
			canView:           true,
		},
		{
			name: "owner can manage their own property",
			user: &User{
				ID:   ownerID,
				Role: RoleOwner,
			},
			propertyOwnerID:   ownerID,
			propertyAgencyID:  nil,
			canManage:         true,
			canCreate:         true,
			canView:           true,
		},
		{
			name: "owner cannot manage other's property",
			user: &User{
				ID:   uuid.New().String(),
				Role: RoleOwner,
			},
			propertyOwnerID:   ownerID,
			propertyAgencyID:  nil,
			canManage:         false,
			canCreate:         true,
			canView:           true,
		},
		{
			name: "agent can manage property assigned to their agency",
			user: &User{
				ID:       agentID,
				Role:     RoleAgent,
				AgencyID: &agencyID,
			},
			propertyOwnerID:   ownerID,
			propertyAgencyID:  &agencyID,
			canManage:         true,
			canCreate:         true,
			canView:           true,
		},
		{
			name: "agent cannot manage property from different agency",
			user: &User{
				ID:       agentID,
				Role:     RoleAgent,
				AgencyID: stringPtr(uuid.New().String()),
			},
			propertyOwnerID:   ownerID,
			propertyAgencyID:  &agencyID,
			canManage:         false,
			canCreate:         true,
			canView:           true,
		},
		{
			name: "buyer can only view properties",
			user: &User{
				ID:   uuid.New().String(),
				Role: RoleBuyer,
			},
			propertyOwnerID:   ownerID,
			propertyAgencyID:  nil,
			canManage:         false,
			canCreate:         false,
			canView:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canManage := tt.user.CanManageProperty(tt.propertyOwnerID, tt.propertyAgencyID)
			canCreate := tt.user.CanCreateProperty()
			canView := tt.user.CanViewProperty(tt.propertyOwnerID, tt.propertyAgencyID)

			assert.Equal(t, tt.canManage, canManage, "CanManageProperty mismatch")
			assert.Equal(t, tt.canCreate, canCreate, "CanCreateProperty mismatch")
			assert.Equal(t, tt.canView, canView, "CanViewProperty mismatch")
		})
	}
}

func TestUserRoleLevel(t *testing.T) {
	tests := []struct {
		role     UserRole
		expected int
	}{
		{RoleAdmin, 5},
		{RoleAgency, 4},
		{RoleAgent, 3},
		{RoleOwner, 2},
		{RoleBuyer, 1},
		{UserRole("invalid"), 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			user := &User{Role: tt.role}
			level := user.GetRoleLevel()
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestUserRoleCheckers(t *testing.T) {
	tests := []struct {
		role     UserRole
		isAdmin  bool
		isAgency bool
		isAgent  bool
		isOwner  bool
		isBuyer  bool
	}{
		{RoleAdmin, true, false, false, false, false},
		{RoleAgency, false, true, false, false, false},
		{RoleAgent, false, false, true, false, false},
		{RoleOwner, false, false, false, true, false},
		{RoleBuyer, false, false, false, false, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			user := &User{Role: tt.role}
			
			assert.Equal(t, tt.isAdmin, user.IsAdmin())
			assert.Equal(t, tt.isAgency, user.IsAgency())
			assert.Equal(t, tt.isAgent, user.IsAgent())
			assert.Equal(t, tt.isOwner, user.IsOwner())
			assert.Equal(t, tt.isBuyer, user.IsBuyer())
		})
	}
}

func TestUserLastLogin(t *testing.T) {
	user := &User{
		ID:   uuid.New().String(),
		Role: RoleOwner,
	}

	assert.Nil(t, user.LastLoginAt)
	
	user.UpdateLastLogin()
	
	assert.NotNil(t, user.LastLoginAt)
	assert.True(t, user.LastLoginAt.After(time.Now().Add(-time.Second)))
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}