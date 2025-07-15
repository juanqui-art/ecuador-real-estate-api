package service

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// UserServiceSimple handles basic user operations
type UserServiceSimple struct {
	userRepo   *repository.UserRepository
	agencyRepo *repository.AgencyRepository
	logger     *log.Logger
}

// NewUserService creates a new simplified user service
func NewUserService(userRepo *repository.UserRepository, agencyRepo *repository.AgencyRepository, logger *log.Logger) *UserServiceSimple {
	return &UserServiceSimple{
		userRepo:   userRepo,
		agencyRepo: agencyRepo,
		logger:     logger,
	}
}

// CreateUser creates a new user with validation
func (s *UserServiceSimple) CreateUser(firstName, lastName, email, phone, cedula, password string, role domain.UserRole) (*domain.User, error) {
	// Validate basic data
	if firstName == "" || lastName == "" || email == "" || phone == "" || cedula == "" || password == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	// Check if user already exists
	if existing, _ := s.userRepo.GetByEmail(email); existing != nil {
		return nil, fmt.Errorf("user with email already exists")
	}

	if existing, _ := s.userRepo.GetByNationalID(cedula); existing != nil {
		return nil, fmt.Errorf("user with national ID already exists")
	}

	// Create user with basic data
	user, err := domain.NewUser(email, firstName, lastName, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	user.Cedula = &cedula
	user.Phone = &phone
	user.Active = true
	user.Status = domain.StatusActive

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	// Validate user
	if err := user.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Printf("User created successfully: %s (%s)", user.Name(), user.Email)
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserServiceSimple) GetUser(id string) (*domain.User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (s *UserServiceSimple) UpdateUser(user *domain.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	// Validate user data
	if err := user.IsValid(); err != nil {
		return fmt.Errorf("invalid user data: %w", err)
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Printf("User updated successfully: %s", user.Name())
	return nil
}

// DeleteUser soft deletes a user
func (s *UserServiceSimple) DeleteUser(id string) error {
	if id == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Soft delete
	now := time.Now()
	user.DeletedAt = &now
	user.Active = false
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Printf("User deleted successfully: %s", user.Name())
	return nil
}

// AuthenticateUser validates user credentials
func (s *UserServiceSimple) AuthenticateUser(email, password string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.Active {
		return nil, fmt.Errorf("account is inactive")
	}

	if user.DeletedAt != nil {
		return nil, fmt.Errorf("account not found")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = time.Now()
	s.userRepo.Update(user)

	s.logger.Printf("User authenticated successfully: %s", user.Email)
	return user, nil
}

// ChangePassword changes user password
func (s *UserServiceSimple) ChangePassword(userID, oldPassword, newPassword string) error {
	if userID == "" || oldPassword == "" || newPassword == "" {
		return fmt.Errorf("all fields are required")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.Printf("Password changed successfully for user: %s", user.Email)
	return nil
}

// SearchUsers searches users with filters
func (s *UserServiceSimple) SearchUsers(email, name string, role domain.UserRole, active *bool, limit, offset int) ([]*domain.User, int, error) {
	params := &domain.UserSearchParams{
		Query:  name,
		Active: active,
		Pagination: &domain.PaginationParams{
			Page:     offset/limit + 1,
			PageSize: limit,
		},
	}
	
	if role != "" {
		params.Role = &role
	}

	users, total, err := s.userRepo.Search(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	return users, total, nil
}

// GetUsersByRole gets all users with a specific role
func (s *UserServiceSimple) GetUsersByRole(role domain.UserRole) ([]*domain.User, error) {
	params := &domain.UserSearchParams{
		Role: &role,
		Pagination: &domain.PaginationParams{
			Page:     1,
			PageSize: 1000,
		},
	}

	users, _, err := s.userRepo.Search(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}

	return users, nil
}

// GetUserStatistics returns user statistics
func (s *UserServiceSimple) GetUserStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count users by role
	for _, role := range []domain.UserRole{domain.RoleAdmin, domain.RoleAgent, domain.RoleBuyer, domain.RoleOwner} {
		users, err := s.GetUsersByRole(role)
		if err != nil {
			continue
		}
		stats[string(role)+"_count"] = len(users)
	}

	// Total active users
	activeTrue := true
	activeParams := &domain.UserSearchParams{
		Active: &activeTrue,
		Pagination: &domain.PaginationParams{
			Page:     1,
			PageSize: 10000,
		},
	}
	activeUsers, _, err := s.userRepo.Search(activeParams)
	if err == nil {
		stats["total_active"] = len(activeUsers)
	}

	stats["last_updated"] = time.Now()
	return stats, nil
}