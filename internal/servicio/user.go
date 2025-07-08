package servicio

import (
	"fmt"
	"realty-core/internal/domain"
	"realty-core/internal/repositorio"
	"strings"
	"time"
)

// UserService handles user business logic
type UserService struct {
	userRepo              *repositorio.UserRepository
	realEstateCompanyRepo *repositorio.RealEstateCompanyRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo *repositorio.UserRepository, realEstateCompanyRepo *repositorio.RealEstateCompanyRepository) *UserService {
	return &UserService{
		userRepo:              userRepo,
		realEstateCompanyRepo: realEstateCompanyRepo,
	}
}

// Create creates a new user with validation
func (s *UserService) Create(firstName, lastName, email, phone, nationalID, userType string) (*domain.User, error) {
	// Create user
	user := domain.NewUser(firstName, lastName, email, phone, nationalID, userType)

	// Validate the user
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("a user with email %s already exists", email)
	}

	// Check if national ID already exists
	exists, err = s.userRepo.ExistsByNationalID(nationalID)
	if err != nil {
		return nil, fmt.Errorf("error checking national ID existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("a user with national ID %s already exists", nationalID)
	}

	// Validate national ID using database function
	isValid, err := s.userRepo.ValidateNationalID(nationalID)
	if err != nil {
		return nil, fmt.Errorf("error validating national ID: %w", err)
	}
	if !isValid {
		return nil, fmt.Errorf("national ID %s is not valid according to Ecuador algorithm", nationalID)
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

// CreateBuyer creates a new buyer user with preferences
func (s *UserService) CreateBuyer(firstName, lastName, email, phone, nationalID string, minBudget, maxBudget *float64, provinces, propertyTypes []string) (*domain.User, error) {
	// Create buyer
	user := domain.NewBuyer(firstName, lastName, email, phone, nationalID, minBudget, maxBudget, provinces, propertyTypes)

	// Validate the user
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("a user with email %s already exists", email)
	}

	// Check if national ID already exists
	exists, err = s.userRepo.ExistsByNationalID(nationalID)
	if err != nil {
		return nil, fmt.Errorf("error checking national ID existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("a user with national ID %s already exists", nationalID)
	}

	// Validate national ID
	isValid, err := s.userRepo.ValidateNationalID(nationalID)
	if err != nil {
		return nil, fmt.Errorf("error validating national ID: %w", err)
	}
	if !isValid {
		return nil, fmt.Errorf("national ID %s is not valid according to Ecuador algorithm", nationalID)
	}

	// Validate provinces
	for _, province := range provinces {
		if !domain.IsValidEcuadorProvince(province) {
			return nil, fmt.Errorf("invalid Ecuador province: %s", province)
		}
	}

	// Validate property types
	validPropertyTypes := []string{domain.PropertyTypeHouse, domain.PropertyTypeApartment, domain.PropertyTypeLand, domain.PropertyTypeCommercial}
	for _, propertyType := range propertyTypes {
		isValid := false
		for _, validType := range validPropertyTypes {
			if propertyType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, fmt.Errorf("invalid property type: %s", propertyType)
		}
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating buyer: %w", err)
	}

	return user, nil
}

// CreateAgent creates a new agent user associated with a real estate company
func (s *UserService) CreateAgent(firstName, lastName, email, phone, nationalID, realEstateCompanyID string) (*domain.User, error) {
	// Validate that the real estate company exists and is active
	company, err := s.realEstateCompanyRepo.GetByID(realEstateCompanyID)
	if err != nil {
		return nil, fmt.Errorf("real estate company not found: %w", err)
	}
	if !company.Active {
		return nil, fmt.Errorf("cannot create agent for inactive company")
	}

	// Create agent
	user := domain.NewAgent(firstName, lastName, email, phone, nationalID, realEstateCompanyID)

	// Validate the user
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("a user with email %s already exists", email)
	}

	// Check if national ID already exists
	exists, err = s.userRepo.ExistsByNationalID(nationalID)
	if err != nil {
		return nil, fmt.Errorf("error checking national ID existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("a user with national ID %s already exists", nationalID)
	}

	// Validate national ID
	isValid, err := s.userRepo.ValidateNationalID(nationalID)
	if err != nil {
		return nil, fmt.Errorf("error validating national ID: %w", err)
	}
	if !isValid {
		return nil, fmt.Errorf("national ID %s is not valid according to Ecuador algorithm", nationalID)
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating agent: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id string) (*domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("email is required")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByNationalID retrieves a user by national ID
func (s *UserService) GetByNationalID(nationalID string) (*domain.User, error) {
	if strings.TrimSpace(nationalID) == "" {
		return nil, fmt.Errorf("national ID is required")
	}

	user, err := s.userRepo.GetByNationalID(nationalID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAll retrieves all users
func (s *UserService) GetAll() ([]*domain.User, error) {
	return s.userRepo.GetAll()
}

// GetBuyers retrieves all buyer users
func (s *UserService) GetBuyers() ([]*domain.User, error) {
	return s.userRepo.GetBuyers()
}

// GetSellers retrieves all seller users
func (s *UserService) GetSellers() ([]*domain.User, error) {
	return s.userRepo.GetSellers()
}

// GetAgents retrieves all agent users
func (s *UserService) GetAgents() ([]*domain.User, error) {
	return s.userRepo.GetAgents()
}

// GetAgentsByCompany retrieves agents for a specific real estate company
func (s *UserService) GetAgentsByCompany(companyID string) ([]*domain.User, error) {
	// Validate company exists
	_, err := s.realEstateCompanyRepo.GetByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("real estate company not found: %w", err)
	}

	return s.userRepo.GetAgentsByCompany(companyID)
}

// Update updates basic user information
func (s *UserService) Update(id, firstName, lastName, email, phone string, dateOfBirth *time.Time, avatarURL, bio string) (*domain.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Store original email for validation
	originalEmail := user.Email

	// Update fields
	user.FirstName = strings.TrimSpace(firstName)
	user.LastName = strings.TrimSpace(lastName)
	user.Email = strings.ToLower(strings.TrimSpace(email))
	user.Phone = strings.TrimSpace(phone)
	user.DateOfBirth = dateOfBirth
	user.AvatarURL = strings.TrimSpace(avatarURL)
	user.Bio = strings.TrimSpace(bio)

	// Validate updated user
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email changed and if new email already exists
	if user.Email != originalEmail {
		exists, err := s.userRepo.ExistsByEmail(user.Email)
		if err != nil {
			return nil, fmt.Errorf("error checking email existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("a user with email %s already exists", user.Email)
		}
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}

// UpdateBuyerPreferences updates buyer search preferences
func (s *UserService) UpdateBuyerPreferences(id string, minBudget, maxBudget *float64, provinces, propertyTypes []string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if !user.IsBuyer() {
		return nil, fmt.Errorf("user is not a buyer")
	}

	// Set budget
	if err := user.SetBudget(minBudget, maxBudget); err != nil {
		return nil, fmt.Errorf("invalid budget: %w", err)
	}

	// Validate provinces
	for _, province := range provinces {
		if !domain.IsValidEcuadorProvince(province) {
			return nil, fmt.Errorf("invalid Ecuador province: %s", province)
		}
	}

	// Validate property types
	validPropertyTypes := []string{domain.PropertyTypeHouse, domain.PropertyTypeApartment, domain.PropertyTypeLand, domain.PropertyTypeCommercial}
	for _, propertyType := range propertyTypes {
		isValid := false
		for _, validType := range validPropertyTypes {
			if propertyType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, fmt.Errorf("invalid property type: %s", propertyType)
		}
	}

	// Set preferences
	user.SetPreferences(provinces, propertyTypes)

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("error updating buyer preferences: %w", err)
	}

	return user, nil
}

// ChangeRealEstateCompany changes the real estate company for an agent
func (s *UserService) ChangeRealEstateCompany(userID, companyID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsAgent() {
		return nil, fmt.Errorf("user is not an agent")
	}

	// Validate that the real estate company exists and is active
	company, err := s.realEstateCompanyRepo.GetByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("real estate company not found: %w", err)
	}
	if !company.Active {
		return nil, fmt.Errorf("cannot assign agent to inactive company")
	}

	// Set company
	user.SetRealEstateCompany(companyID)

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("error changing real estate company: %w", err)
	}

	return user, nil
}

// Activate activates a user
func (s *UserService) Activate(id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user.Active {
		return user, nil // Already active
	}

	// Activate
	if err := s.userRepo.Activate(id); err != nil {
		return nil, fmt.Errorf("error activating user: %w", err)
	}

	// Get updated user
	return s.userRepo.GetByID(id)
}

// Deactivate deactivates a user
func (s *UserService) Deactivate(id string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if !user.Active {
		return nil // Already inactive
	}

	// Deactivate
	if err := s.userRepo.Deactivate(id); err != nil {
		return fmt.Errorf("error deactivating user: %w", err)
	}

	return nil
}

// SearchByName searches users by name
func (s *UserService) SearchByName(searchTerm string) ([]*domain.User, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return nil, fmt.Errorf("search term is required")
	}

	return s.userRepo.SearchByName(searchTerm)
}

// GetBuyersForProperty gets buyers that can afford a specific property price
func (s *UserService) GetBuyersForProperty(propertyPrice float64) ([]*domain.User, error) {
	if propertyPrice <= 0 {
		return nil, fmt.Errorf("property price must be greater than 0")
	}

	return s.userRepo.GetBuyersForProperty(propertyPrice)
}

// Delete permanently deletes a user
func (s *UserService) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("user ID is required")
	}

	// Verify user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}

// GetStatistics returns statistics about users
func (s *UserService) GetStatistics() (map[string]interface{}, error) {
	return s.userRepo.GetStatistics()
}

// ValidateNationalID validates a national ID format and algorithm
func (s *UserService) ValidateNationalID(nationalID string) error {
	// Basic format validation
	nationalID = strings.TrimSpace(nationalID)
	if len(nationalID) != 10 {
		return fmt.Errorf("national ID must be 10 digits")
	}

	// Check using domain validation
	tempUser := &domain.User{NationalID: nationalID}
	if err := tempUser.ValidateNationalID(); err != nil {
		return err
	}

	// Check using database algorithm
	isValid, err := s.userRepo.ValidateNationalID(nationalID)
	if err != nil {
		return fmt.Errorf("error validating national ID: %w", err)
	}
	if !isValid {
		return fmt.Errorf("national ID %s is not valid according to Ecuador algorithm", nationalID)
	}

	return nil
}

// ValidateEmail validates an email format
func (s *UserService) ValidateEmail(email string) error {
	tempUser := &domain.User{Email: email}
	return tempUser.ValidateEmail()
}

// ValidatePhone validates a phone format for Ecuador
func (s *UserService) ValidatePhone(phone string) error {
	tempUser := &domain.User{Phone: phone}
	return tempUser.ValidatePhone()
}

// CheckEmailAvailability checks if an email is available for use
func (s *UserService) CheckEmailAvailability(email string) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, fmt.Errorf("email is required")
	}

	// First validate email format
	if err := s.ValidateEmail(email); err != nil {
		return false, err
	}

	// Check if it exists
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return false, fmt.Errorf("error checking email availability: %w", err)
	}

	return !exists, nil
}

// CheckNationalIDAvailability checks if a national ID is available for use
func (s *UserService) CheckNationalIDAvailability(nationalID string) (bool, error) {
	if strings.TrimSpace(nationalID) == "" {
		return false, fmt.Errorf("national ID is required")
	}

	// First validate national ID format
	if err := s.ValidateNationalID(nationalID); err != nil {
		return false, err
	}

	// Check if it exists
	exists, err := s.userRepo.ExistsByNationalID(nationalID)
	if err != nil {
		return false, fmt.Errorf("error checking national ID availability: %w", err)
	}

	return !exists, nil
}

// SetNotificationPreferences sets user notification preferences
func (s *UserService) SetNotificationPreferences(id string, notifications, newsletter bool) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.SetNotificationPreferences(notifications, newsletter)

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("error updating notification preferences: %w", err)
	}

	return user, nil
}
