package servicio

import (
	"fmt"
	"realty-core/internal/domain"
	"realty-core/internal/repositorio"
	"strings"
)

// RealEstateCompanyService handles real estate company business logic
type RealEstateCompanyService struct {
	companyRepo *repositorio.RealEstateCompanyRepository
}

// NewRealEstateCompanyService creates a new real estate company service instance
func NewRealEstateCompanyService(companyRepo *repositorio.RealEstateCompanyRepository) *RealEstateCompanyService {
	return &RealEstateCompanyService{
		companyRepo: companyRepo,
	}
}

// Create creates a new real estate company with validation
func (s *RealEstateCompanyService) Create(name, ruc, address, phone, email string) (*domain.RealEstateCompany, error) {
	// Create company
	company := domain.NewRealEstateCompany(name, ruc, address, phone, email)

	// Validate the company
	if err := company.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if RUC already exists
	exists, err := s.companyRepo.ExistsByRUC(ruc)
	if err != nil {
		return nil, fmt.Errorf("error checking RUC existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("a company with RUC %s already exists", ruc)
	}

	// Check if email already exists
	exists, err = s.companyRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("a company with email %s already exists", email)
	}

	// Validate RUC using the database function
	isValid, err := s.companyRepo.ValidateRUC(ruc)
	if err != nil {
		return nil, fmt.Errorf("error validating RUC: %w", err)
	}
	if !isValid {
		return nil, fmt.Errorf("RUC %s is not valid according to Ecuador algorithm", ruc)
	}

	// Save to database
	if err := s.companyRepo.Create(company); err != nil {
		return nil, fmt.Errorf("error creating real estate company: %w", err)
	}

	return company, nil
}

// GetByID retrieves a real estate company by ID
func (s *RealEstateCompanyService) GetByID(id string) (*domain.RealEstateCompany, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("company ID is required")
	}

	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return company, nil
}

// GetByRUC retrieves a real estate company by RUC
func (s *RealEstateCompanyService) GetByRUC(ruc string) (*domain.RealEstateCompany, error) {
	if strings.TrimSpace(ruc) == "" {
		return nil, fmt.Errorf("RUC is required")
	}

	company, err := s.companyRepo.GetByRUC(ruc)
	if err != nil {
		return nil, err
	}

	return company, nil
}

// GetAll retrieves all real estate companies
func (s *RealEstateCompanyService) GetAll() ([]*domain.RealEstateCompany, error) {
	return s.companyRepo.GetAll()
}

// GetActive retrieves all active real estate companies
func (s *RealEstateCompanyService) GetActive() ([]*domain.RealEstateCompany, error) {
	return s.companyRepo.GetActive()
}

// Update updates an existing real estate company
func (s *RealEstateCompanyService) Update(id, name, ruc, address, phone, email, website, description, logoURL string) (*domain.RealEstateCompany, error) {
	// Get existing company
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Store original values for validation
	originalRUC := company.RUC
	originalEmail := company.Email

	// Update fields
	company.Name = strings.TrimSpace(name)
	company.RUC = strings.TrimSpace(ruc)
	company.Address = strings.TrimSpace(address)
	company.Phone = strings.TrimSpace(phone)
	company.Email = strings.ToLower(strings.TrimSpace(email))
	company.Website = strings.TrimSpace(website)
	company.Description = strings.TrimSpace(description)
	company.LogoURL = strings.TrimSpace(logoURL)

	// Validate updated company
	if err := company.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if RUC changed and if new RUC already exists
	if company.RUC != originalRUC {
		exists, err := s.companyRepo.ExistsByRUC(company.RUC)
		if err != nil {
			return nil, fmt.Errorf("error checking RUC existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("a company with RUC %s already exists", company.RUC)
		}

		// Validate new RUC
		isValid, err := s.companyRepo.ValidateRUC(company.RUC)
		if err != nil {
			return nil, fmt.Errorf("error validating RUC: %w", err)
		}
		if !isValid {
			return nil, fmt.Errorf("RUC %s is not valid according to Ecuador algorithm", company.RUC)
		}
	}

	// Check if email changed and if new email already exists
	if company.Email != originalEmail {
		exists, err := s.companyRepo.ExistsByEmail(company.Email)
		if err != nil {
			return nil, fmt.Errorf("error checking email existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("a company with email %s already exists", company.Email)
		}
	}

	// Validate optional fields
	if company.Website != "" {
		if err := company.ValidateWebsite(); err != nil {
			return nil, fmt.Errorf("invalid website: %w", err)
		}
	}

	if company.LogoURL != "" {
		if err := company.ValidateLogoURL(); err != nil {
			return nil, fmt.Errorf("invalid logo URL: %w", err)
		}
	}

	// Save changes
	if err := s.companyRepo.Update(company); err != nil {
		return nil, fmt.Errorf("error updating real estate company: %w", err)
	}

	return company, nil
}

// UpdateContactInfo updates contact information
func (s *RealEstateCompanyService) UpdateContactInfo(id, phone, email string) (*domain.RealEstateCompany, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Store original email for validation
	originalEmail := company.Email

	// Update contact info with validation
	if err := company.SetContactInfo(phone, email); err != nil {
		return nil, fmt.Errorf("invalid contact information: %w", err)
	}

	// Check if email changed and if new email already exists
	if company.Email != originalEmail {
		exists, err := s.companyRepo.ExistsByEmail(company.Email)
		if err != nil {
			return nil, fmt.Errorf("error checking email existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("a company with email %s already exists", company.Email)
		}
	}

	// Save changes
	if err := s.companyRepo.Update(company); err != nil {
		return nil, fmt.Errorf("error updating contact information: %w", err)
	}

	return company, nil
}

// UpdateWebsite updates the company website
func (s *RealEstateCompanyService) UpdateWebsite(id, website string) (*domain.RealEstateCompany, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Set and validate website
	if err := company.SetWebsite(website); err != nil {
		return nil, fmt.Errorf("invalid website: %w", err)
	}

	// Save changes
	if err := s.companyRepo.Update(company); err != nil {
		return nil, fmt.Errorf("error updating website: %w", err)
	}

	return company, nil
}

// UpdateLogoURL updates the company logo URL
func (s *RealEstateCompanyService) UpdateLogoURL(id, logoURL string) (*domain.RealEstateCompany, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Set and validate logo URL
	if err := company.SetLogoURL(logoURL); err != nil {
		return nil, fmt.Errorf("invalid logo URL: %w", err)
	}

	// Save changes
	if err := s.companyRepo.Update(company); err != nil {
		return nil, fmt.Errorf("error updating logo URL: %w", err)
	}

	return company, nil
}

// Activate activates a real estate company
func (s *RealEstateCompanyService) Activate(id string) (*domain.RealEstateCompany, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if company.Active {
		return company, nil // Already active
	}

	// Activate
	if err := s.companyRepo.Activate(id); err != nil {
		return nil, fmt.Errorf("error activating company: %w", err)
	}

	// Get updated company
	return s.companyRepo.GetByID(id)
}

// Deactivate deactivates a real estate company
func (s *RealEstateCompanyService) Deactivate(id string) error {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return err
	}

	if !company.Active {
		return nil // Already inactive
	}

	// Deactivate
	if err := s.companyRepo.Deactivate(id); err != nil {
		return fmt.Errorf("error deactivating company: %w", err)
	}

	return nil
}

// SearchByName searches real estate companies by name
func (s *RealEstateCompanyService) SearchByName(searchTerm string) ([]*domain.RealEstateCompany, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return nil, fmt.Errorf("search term is required")
	}

	return s.companyRepo.SearchByName(searchTerm)
}

// Delete permanently deletes a real estate company
func (s *RealEstateCompanyService) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("company ID is required")
	}

	// Verify company exists
	_, err := s.companyRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.companyRepo.Delete(id)
}

// GetStatistics returns statistics about real estate companies
func (s *RealEstateCompanyService) GetStatistics() (map[string]interface{}, error) {
	return s.companyRepo.GetStatistics()
}

// GetCompaniesWithPropertyCount returns companies with their property counts
func (s *RealEstateCompanyService) GetCompaniesWithPropertyCount() ([]map[string]interface{}, error) {
	return s.companyRepo.GetCompaniesWithPropertyCount()
}

// ValidateRUC validates a RUC format and algorithm
func (s *RealEstateCompanyService) ValidateRUC(ruc string) error {
	// Basic format validation
	ruc = strings.TrimSpace(ruc)
	if len(ruc) != 13 {
		return fmt.Errorf("RUC must be 13 digits")
	}

	// Check using domain validation
	tempCompany := &domain.RealEstateCompany{RUC: ruc}
	if err := tempCompany.ValidateRUC(); err != nil {
		return err
	}

	// Check using database algorithm
	isValid, err := s.companyRepo.ValidateRUC(ruc)
	if err != nil {
		return fmt.Errorf("error validating RUC: %w", err)
	}
	if !isValid {
		return fmt.Errorf("RUC %s is not valid according to Ecuador algorithm", ruc)
	}

	return nil
}

// ValidateEmail validates an email format
func (s *RealEstateCompanyService) ValidateEmail(email string) error {
	tempCompany := &domain.RealEstateCompany{Email: email}
	return tempCompany.ValidateEmail()
}

// ValidatePhone validates a phone format for Ecuador
func (s *RealEstateCompanyService) ValidatePhone(phone string) error {
	tempCompany := &domain.RealEstateCompany{Phone: phone}
	return tempCompany.ValidatePhone()
}

// CheckRUCAvailability checks if a RUC is available for use
func (s *RealEstateCompanyService) CheckRUCAvailability(ruc string) (bool, error) {
	if strings.TrimSpace(ruc) == "" {
		return false, fmt.Errorf("RUC is required")
	}

	// First validate RUC format
	if err := s.ValidateRUC(ruc); err != nil {
		return false, err
	}

	// Check if it exists
	exists, err := s.companyRepo.ExistsByRUC(ruc)
	if err != nil {
		return false, fmt.Errorf("error checking RUC availability: %w", err)
	}

	return !exists, nil
}

// CheckEmailAvailability checks if an email is available for use
func (s *RealEstateCompanyService) CheckEmailAvailability(email string) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, fmt.Errorf("email is required")
	}

	// First validate email format
	if err := s.ValidateEmail(email); err != nil {
		return false, err
	}

	// Check if it exists
	exists, err := s.companyRepo.ExistsByEmail(email)
	if err != nil {
		return false, fmt.Errorf("error checking email availability: %w", err)
	}

	return !exists, nil
}
