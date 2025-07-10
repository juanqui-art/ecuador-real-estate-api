package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// UserService handles business logic for users
type UserService struct {
	userRepo   *repository.UserRepository
	agencyRepo *repository.AgencyRepository
	logger     *log.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository, agencyRepo *repository.AgencyRepository, logger *log.Logger) *UserService {
	return &UserService{
		userRepo:   userRepo,
		agencyRepo: agencyRepo,
		logger:     logger,
	}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(firstName, lastName, email, phone, cedula, password string, role domain.Role) (*domain.User, error) {
	// Validate basic data
	if firstName == "" || lastName == "" || email == "" || phone == "" || cedula == "" || password == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	// Check if user already exists
	if existing, _ := s.userRepo.GetByEmail(email); existing != nil {
		return nil, fmt.Errorf("user with email already exists")
	}

	if existing, _ := s.userRepo.GetByCedula(cedula); existing != nil {
		return nil, fmt.Errorf("user with cedula already exists")
	}

	// Create user
	user := domain.NewUser(firstName, lastName, email, phone, cedula, role)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	// Validate user
	if !user.IsValid() {
		return nil, fmt.Errorf("invalid user data")
	}

	// Validate business rules
	if err := user.ValidateBusinessRules(); err != nil {
		return nil, fmt.Errorf("business rule validation failed: %w", err)
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Printf("User created successfully: %s (%s)", user.GetFullName(), user.Email)
	return user, nil
}

// CreateAgent creates a new agent associated with an agency
func (s *UserService) CreateAgent(firstName, lastName, email, phone, cedula, password, agencyID string) (*domain.User, error) {
	// Validate agency exists and is active
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return nil, fmt.Errorf("agency not found: %w", err)
	}

	if !agency.Active {
		return nil, fmt.Errorf("cannot create agent for inactive agency")
	}

	if !agency.IsLicenseValid() {
		return nil, fmt.Errorf("cannot create agent for agency with invalid license")
	}

	// Create agent
	user, err := s.CreateUser(firstName, lastName, email, phone, cedula, password, domain.RoleAgent)
	if err != nil {
		return nil, err
	}

	// Associate with agency
	if err := user.SetAgency(&agencyID); err != nil {
		return nil, fmt.Errorf("failed to associate agent with agency: %w", err)
	}

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update agent with agency: %w", err)
	}

	s.logger.Printf("Agent created successfully: %s for agency %s", user.GetFullName(), agency.Name)
	return user, nil
}

// CreateBuyer creates a new buyer with budget preferences
func (s *UserService) CreateBuyer(firstName, lastName, email, phone, cedula, password string, minBudget, maxBudget *float64, provinces, propertyTypes []string) (*domain.User, error) {
	// Create buyer
	user, err := s.CreateUser(firstName, lastName, email, phone, cedula, password, domain.RoleBuyer)
	if err != nil {
		return nil, err
	}

	// Set budget preferences
	if minBudget != nil && maxBudget != nil {
		if err := user.SetBudget(minBudget, maxBudget); err != nil {
			return nil, fmt.Errorf("failed to set budget: %w", err)
		}
	}

	// Set interested provinces
	for _, province := range provinces {
		if err := user.AddInterestedProvince(province); err != nil {
			return nil, fmt.Errorf("failed to add interested province: %w", err)
		}
	}

	// Set interested property types
	for _, propertyType := range propertyTypes {
		if err := user.AddInterestedType(propertyType); err != nil {
			return nil, fmt.Errorf("failed to add interested property type: %w", err)
		}
	}

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update buyer preferences: %w", err)
	}

	s.logger.Printf("Buyer created successfully: %s", user.GetFullName())
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(user *domain.User) error {
	// Validate user
	if !user.IsValid() {
		return fmt.Errorf("invalid user data")
	}

	// Validate business rules
	if err := user.ValidateBusinessRules(); err != nil {
		return fmt.Errorf("business rule validation failed: %w", err)
	}

	// Update timestamp
	user.UpdateTimestamp()

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Printf("User updated successfully: %s", user.GetFullName())
	return nil
}

// DeleteUser soft deletes a user (deactivates)
func (s *UserService) DeleteUser(id string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Deactivate user
	user.Deactivate()

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	s.logger.Printf("User deactivated successfully: %s", user.GetFullName())
	return nil
}

// AuthenticateUser authenticates a user with email and password
func (s *UserService) AuthenticateUser(email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if !user.Active {
		return nil, fmt.Errorf("user account is deactivated")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		s.logger.Printf("Failed to update last login for user %s: %v", user.ID, err)
	}

	s.logger.Printf("User authenticated successfully: %s", user.Email)
	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdateTimestamp()

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.Printf("Password changed successfully for user: %s", user.Email)
	return nil
}

// RequestPasswordReset generates a password reset token
func (s *UserService) RequestPasswordReset(email string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Set token expiry (24 hours)
	expiry := time.Now().Add(24 * time.Hour)

	// Save token to database
	if err := s.userRepo.SetPasswordResetToken(user.ID, token, expiry); err != nil {
		return "", fmt.Errorf("failed to save reset token: %w", err)
	}

	s.logger.Printf("Password reset token generated for user: %s", user.Email)
	return token, nil
}

// ResetPassword resets a user's password using a reset token
func (s *UserService) ResetPassword(token, newPassword string) error {
	// Find user by reset token
	user, err := s.userRepo.GetByEmail("") // We need to modify this to search by token
	if err != nil {
		return fmt.Errorf("invalid reset token")
	}

	// Verify token and expiry
	if user.PasswordResetToken == nil || *user.PasswordResetToken != token {
		return fmt.Errorf("invalid reset token")
	}

	if user.PasswordResetExpires == nil || time.Now().After(*user.PasswordResetExpires) {
		return fmt.Errorf("reset token has expired")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdateTimestamp()

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Clear reset token
	if err := s.userRepo.ClearPasswordResetToken(user.ID); err != nil {
		s.logger.Printf("Failed to clear reset token for user %s: %v", user.ID, err)
	}

	s.logger.Printf("Password reset successfully for user: %s", user.Email)
	return nil
}

// SendEmailVerification generates and sends an email verification token
func (s *UserService) SendEmailVerification(userID string) (string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if user.EmailVerified {
		return "", fmt.Errorf("email already verified")
	}

	// Generate verification token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate verification token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	user.EmailVerificationToken = &token
	user.UpdateTimestamp()

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return "", fmt.Errorf("failed to save verification token: %w", err)
	}

	s.logger.Printf("Email verification token generated for user: %s", user.Email)
	return token, nil
}

// VerifyEmail verifies a user's email using a verification token
func (s *UserService) VerifyEmail(token string) error {
	// Find user by verification token (we need to modify the repository for this)
	// For now, we'll implement a simple approach
	user, err := s.userRepo.GetByEmail("") // We need to modify this to search by token
	if err != nil {
		return fmt.Errorf("invalid verification token")
	}

	// Verify token
	if user.EmailVerificationToken == nil || *user.EmailVerificationToken != token {
		return fmt.Errorf("invalid verification token")
	}

	// Set email as verified
	if err := s.userRepo.SetEmailVerified(user.ID, true); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	s.logger.Printf("Email verified successfully for user: %s", user.Email)
	return nil
}

// SearchUsers searches users with filters
func (s *UserService) SearchUsers(params *domain.UserSearchParams) ([]*domain.User, *domain.Pagination, error) {
	if params.Pagination == nil {
		params.Pagination = domain.NewPaginationParams()
	}

	users, totalCount, err := s.userRepo.Search(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search users: %w", err)
	}

	pagination := domain.NewPagination(
		params.Pagination.Page,
		params.Pagination.PageSize,
		totalCount,
	)

	return users, pagination, nil
}

// GetUsersByRole retrieves users by role
func (s *UserService) GetUsersByRole(role domain.Role) ([]*domain.User, error) {
	users, err := s.userRepo.GetByRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}

	return users, nil
}

// GetAgentsByAgency retrieves agents for a specific agency
func (s *UserService) GetAgentsByAgency(agencyID string) ([]*domain.User, error) {
	// Verify agency exists
	_, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return nil, fmt.Errorf("agency not found: %w", err)
	}

	users, err := s.userRepo.GetByAgency(agencyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents by agency: %w", err)
	}

	return users, nil
}

// GetBuyersByBudget finds buyers who can afford a specific property price
func (s *UserService) GetBuyersByBudget(price float64) ([]*domain.User, error) {
	users, err := s.userRepo.GetBuyersByBudget(price)
	if err != nil {
		return nil, fmt.Errorf("failed to get buyers by budget: %w", err)
	}

	return users, nil
}

// GetUserStatistics returns user statistics
func (s *UserService) GetUserStatistics() (*domain.UserStats, error) {
	stats, err := s.userRepo.GetStatistics()
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %w", err)
	}

	return stats, nil
}

// AssignAgentToAgency assigns an agent to a different agency
func (s *UserService) AssignAgentToAgency(agentID, agencyID string) error {
	// Get agent
	agent, err := s.userRepo.GetByID(agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent: %w", err)
	}

	if agent.Role != domain.RoleAgent {
		return fmt.Errorf("user is not an agent")
	}

	// Verify agency exists and is active
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("agency not found: %w", err)
	}

	if !agency.Active {
		return fmt.Errorf("cannot assign agent to inactive agency")
	}

	if !agency.IsLicenseValid() {
		return fmt.Errorf("cannot assign agent to agency with invalid license")
	}

	// Assign agent to agency
	if err := agent.SetAgency(&agencyID); err != nil {
		return fmt.Errorf("failed to assign agent to agency: %w", err)
	}

	// Update in database
	if err := s.userRepo.Update(agent); err != nil {
		return fmt.Errorf("failed to update agent assignment: %w", err)
	}

	s.logger.Printf("Agent %s assigned to agency %s", agent.GetFullName(), agency.Name)
	return nil
}

// UpdateUserPreferences updates buyer preferences
func (s *UserService) UpdateUserPreferences(userID string, minBudget, maxBudget *float64, provinces, propertyTypes []string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Set budget preferences
	if minBudget != nil && maxBudget != nil {
		if err := user.SetBudget(minBudget, maxBudget); err != nil {
			return fmt.Errorf("failed to set budget: %w", err)
		}
	}

	// Clear and set interested provinces
	user.InterestedProvinces = []string{}
	for _, province := range provinces {
		if err := user.AddInterestedProvince(province); err != nil {
			return fmt.Errorf("failed to add interested province: %w", err)
		}
	}

	// Clear and set interested property types
	user.InterestedTypes = []string{}
	for _, propertyType := range propertyTypes {
		if err := user.AddInterestedType(propertyType); err != nil {
			return fmt.Errorf("failed to add interested property type: %w", err)
		}
	}

	// Update in database
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user preferences: %w", err)
	}

	s.logger.Printf("User preferences updated successfully: %s", user.GetFullName())
	return nil
}