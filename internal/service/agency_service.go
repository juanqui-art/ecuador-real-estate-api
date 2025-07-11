package service

import (
	"fmt"
	"log"
	"time"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// AgencyService handles business logic for agencies
type AgencyService struct {
	agencyRepo *repository.AgencyRepository
	userRepo   *repository.UserRepository
	logger     *log.Logger
}

// NewAgencyService creates a new agency service
func NewAgencyService(agencyRepo *repository.AgencyRepository, userRepo *repository.UserRepository, logger *log.Logger) *AgencyService {
	return &AgencyService{
		agencyRepo: agencyRepo,
		userRepo:   userRepo,
		logger:     logger,
	}
}

// CreateAgency creates a new agency with validation
func (s *AgencyService) CreateAgency(name, ruc, address, phone, email, licenseNumber string) (*domain.Agency, error) {
	// Validate basic data
	if name == "" || ruc == "" || address == "" || phone == "" || email == "" || licenseNumber == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	// Check if agency already exists
	if existing, _ := s.agencyRepo.GetByRUC(ruc); existing != nil {
		return nil, fmt.Errorf("agency with RUC already exists")
	}

	// Create agency
	agency, err := domain.NewAgency(name, ruc, address, phone, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create agency: %w", err)
	}
	agency.LicenseNumber = licenseNumber

	// Validate agency
	if err := agency.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid agency data: %w", err)
	}

	// Validate business rules
	if err := agency.ValidateBusinessRules(); err != nil {
		return nil, fmt.Errorf("business rule validation failed: %w", err)
	}

	// Save to database
	if err := s.agencyRepo.Create(agency); err != nil {
		return nil, fmt.Errorf("failed to create agency: %w", err)
	}

	s.logger.Printf("Agency created successfully: %s (%s)", agency.Name, agency.RUC)
	return agency, nil
}

// GetAgency retrieves an agency by ID
func (s *AgencyService) GetAgency(id string) (*domain.Agency, error) {
	agency, err := s.agencyRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get agency: %w", err)
	}

	return agency, nil
}

// GetAgencyByRUC retrieves an agency by RUC
func (s *AgencyService) GetAgencyByRUC(ruc string) (*domain.Agency, error) {
	agency, err := s.agencyRepo.GetByRUC(ruc)
	if err != nil {
		return nil, fmt.Errorf("failed to get agency by RUC: %w", err)
	}

	return agency, nil
}

// UpdateAgency updates agency information
func (s *AgencyService) UpdateAgency(agency *domain.Agency) error {
	// Validate agency
	if err := agency.IsValid(); err != nil {
		return fmt.Errorf("invalid agency data: %w", err)
	}

	// Validate business rules
	if err := agency.ValidateBusinessRules(); err != nil {
		return fmt.Errorf("business rule validation failed: %w", err)
	}

	// Update timestamp
	agency.UpdateTimestamp()

	// Update in database
	if err := s.agencyRepo.Update(agency); err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	s.logger.Printf("Agency updated successfully: %s", agency.Name)
	return nil
}

// DeleteAgency soft deletes an agency (deactivates)
func (s *AgencyService) DeleteAgency(id string) error {
	agency, err := s.agencyRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get agency: %w", err)
	}

	// Check if agency has active agents
	agents, err := s.userRepo.GetByAgency(id)
	if err != nil {
		return fmt.Errorf("failed to check agency agents: %w", err)
	}

	if len(agents) > 0 {
		return fmt.Errorf("cannot delete agency with active agents")
	}

	// Deactivate agency
	agency.Deactivate()

	// Update in database
	if err := s.agencyRepo.Update(agency); err != nil {
		return fmt.Errorf("failed to deactivate agency: %w", err)
	}

	s.logger.Printf("Agency deactivated successfully: %s", agency.Name)
	return nil
}

// GetActiveAgencies retrieves all active agencies
func (s *AgencyService) GetActiveAgencies() ([]*domain.Agency, error) {
	agencies, err := s.agencyRepo.GetActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get active agencies: %w", err)
	}

	return agencies, nil
}

// SearchAgencies searches agencies with filters
func (s *AgencyService) SearchAgencies(params *domain.AgencySearchParams) ([]*domain.Agency, *domain.Pagination, error) {
	if params.Pagination == nil {
		params.Pagination = domain.NewPaginationParams()
	}

	agencies, totalCount, err := s.agencyRepo.Search(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search agencies: %w", err)
	}

	pagination := domain.NewPagination(
		params.Pagination.Page,
		params.Pagination.PageSize,
		totalCount,
	)

	return agencies, pagination, nil
}

// GetAgenciesByServiceArea retrieves agencies that serve a specific area
func (s *AgencyService) GetAgenciesByServiceArea(province string) ([]*domain.Agency, error) {
	// Validate province
	if !domain.IsValidProvince(province) {
		return nil, fmt.Errorf("invalid province: %s", province)
	}

	agencies, err := s.agencyRepo.GetByServiceArea(province)
	if err != nil {
		return nil, fmt.Errorf("failed to get agencies by service area: %w", err)
	}

	return agencies, nil
}

// GetAgenciesBySpecialty retrieves agencies with a specific specialty
func (s *AgencyService) GetAgenciesBySpecialty(specialty string) ([]*domain.Agency, error) {
	agencies, err := s.agencyRepo.GetBySpecialty(specialty)
	if err != nil {
		return nil, fmt.Errorf("failed to get agencies by specialty: %w", err)
	}

	return agencies, nil
}

// GetAgencyStatistics returns agency statistics
func (s *AgencyService) GetAgencyStatistics() (*domain.AgencyStats, error) {
	stats, err := s.agencyRepo.GetStatistics()
	if err != nil {
		return nil, fmt.Errorf("failed to get agency statistics: %w", err)
	}

	return stats, nil
}

// GetAgencyPerformance returns performance metrics for an agency
func (s *AgencyService) GetAgencyPerformance(agencyID string) (*domain.AgencyPerformance, error) {
	performance, err := s.agencyRepo.GetPerformance(agencyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agency performance: %w", err)
	}

	return performance, nil
}

// GetAgencyWithAgents returns agency with its associated agents
func (s *AgencyService) GetAgencyWithAgents(agencyID string) (*domain.AgencyWithAgents, error) {
	agencyWithAgents, err := s.agencyRepo.GetWithAgents(agencyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agency with agents: %w", err)
	}

	return agencyWithAgents, nil
}

// AddSpecialtyToAgency adds a specialty to an agency
func (s *AgencyService) AddSpecialtyToAgency(agencyID, specialty string) error {
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("failed to get agency: %w", err)
	}

	if err := agency.AddSpecialty(specialty); err != nil {
		return fmt.Errorf("failed to add specialty: %w", err)
	}

	if err := s.agencyRepo.Update(agency); err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	s.logger.Printf("Specialty '%s' added to agency: %s", specialty, agency.Name)
	return nil
}

// AddServiceAreaToAgency adds a service area to an agency
func (s *AgencyService) AddServiceAreaToAgency(agencyID, province string) error {
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("failed to get agency: %w", err)
	}

	if err := agency.AddServiceArea(province); err != nil {
		return fmt.Errorf("failed to add service area: %w", err)
	}

	if err := s.agencyRepo.Update(agency); err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	s.logger.Printf("Service area '%s' added to agency: %s", province, agency.Name)
	return nil
}

// SetAgencyCommission sets the commission percentage for an agency
func (s *AgencyService) SetAgencyCommission(agencyID string, commission float64) error {
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("failed to get agency: %w", err)
	}

	if err := agency.SetCommission(commission); err != nil {
		return fmt.Errorf("failed to set commission: %w", err)
	}

	if err := s.agencyRepo.Update(agency); err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	s.logger.Printf("Commission set to %.2f%% for agency: %s", commission, agency.Name)
	return nil
}

// SetAgencyLicense sets the license information for an agency
func (s *AgencyService) SetAgencyLicense(agencyID, licenseNumber string, expiry *time.Time) error {
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("failed to get agency: %w", err)
	}

	if err := agency.SetLicense(licenseNumber, expiry); err != nil {
		return fmt.Errorf("failed to set license: %w", err)
	}

	if err := s.agencyRepo.Update(agency); err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	s.logger.Printf("License updated for agency: %s", agency.Name)
	return nil
}

// SetAgencySocialMedia sets social media links for an agency
func (s *AgencyService) SetAgencySocialMedia(agencyID, platform, url string) error {
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("failed to get agency: %w", err)
	}

	if err := agency.SetSocialMedia(platform, url); err != nil {
		return fmt.Errorf("failed to set social media: %w", err)
	}

	if err := s.agencyRepo.Update(agency); err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	s.logger.Printf("Social media '%s' updated for agency: %s", platform, agency.Name)
	return nil
}

// ActivateAgency activates an agency
func (s *AgencyService) ActivateAgency(agencyID string) error {
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("failed to get agency: %w", err)
	}

	// Validate license before activation
	if !agency.IsLicenseValid() {
		return fmt.Errorf("cannot activate agency with invalid license")
	}

	if err := s.agencyRepo.Activate(agencyID); err != nil {
		return fmt.Errorf("failed to activate agency: %w", err)
	}

	s.logger.Printf("Agency activated successfully: %s", agency.Name)
	return nil
}

// DeactivateAgency deactivates an agency
func (s *AgencyService) DeactivateAgency(agencyID string) error {
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("failed to get agency: %w", err)
	}

	// Check if agency has active agents
	agents, err := s.userRepo.GetByAgency(agencyID)
	if err != nil {
		return fmt.Errorf("failed to check agency agents: %w", err)
	}

	if len(agents) > 0 {
		return fmt.Errorf("cannot deactivate agency with active agents")
	}

	if err := s.agencyRepo.Deactivate(agencyID); err != nil {
		return fmt.Errorf("failed to deactivate agency: %w", err)
	}

	s.logger.Printf("Agency deactivated successfully: %s", agency.Name)
	return nil
}

// ValidateAgencyForPropertyManagement validates if an agency can manage properties
func (s *AgencyService) ValidateAgencyForPropertyManagement(agencyID string) error {
	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("agency not found: %w", err)
	}

	if !agency.Active {
		return fmt.Errorf("agency is not active")
	}

	if !agency.IsLicenseValid() {
		return fmt.Errorf("agency license is invalid or expired")
	}

	return nil
}

// GetAgenciesForProperty recommends agencies for a property based on location and type
func (s *AgencyService) GetAgenciesForProperty(province, propertyType string) ([]*domain.Agency, error) {
	// Validate province
	if !domain.IsValidProvince(province) {
		return nil, fmt.Errorf("invalid province: %s", province)
	}

	// Validate property type
	if !domain.IsValidPropertyType(propertyType) {
		return nil, fmt.Errorf("invalid property type: %s", propertyType)
	}

	// Get agencies that serve the area
	agencies, err := s.agencyRepo.GetByServiceArea(province)
	if err != nil {
		return nil, fmt.Errorf("failed to get agencies by service area: %w", err)
	}

	// Filter by specialty if applicable
	var matchingAgencies []*domain.Agency
	for _, agency := range agencies {
		// Check if agency specializes in the property type
		hasMatchingSpecialty := false
		for _, specialty := range agency.Specialties {
			if specialty == "residencial" && (propertyType == "house" || propertyType == "apartment") {
				hasMatchingSpecialty = true
				break
			}
			if specialty == "comercial" && propertyType == "commercial" {
				hasMatchingSpecialty = true
				break
			}
			if specialty == "terrenos" && propertyType == "land" {
				hasMatchingSpecialty = true
				break
			}
		}

		// Include agency if it has matching specialty or general capabilities
		if hasMatchingSpecialty || len(agency.Specialties) == 0 {
			matchingAgencies = append(matchingAgencies, agency)
		}
	}

	return matchingAgencies, nil
}

// TransferAgentToAgency transfers an agent from one agency to another
func (s *AgencyService) TransferAgentToAgency(agentID, fromAgencyID, toAgencyID string) error {
	// Validate source agency
	fromAgency, err := s.agencyRepo.GetByID(fromAgencyID)
	if err != nil {
		return fmt.Errorf("source agency not found: %w", err)
	}

	// Validate target agency
	toAgency, err := s.agencyRepo.GetByID(toAgencyID)
	if err != nil {
		return fmt.Errorf("target agency not found: %w", err)
	}

	if !toAgency.Active {
		return fmt.Errorf("cannot transfer agent to inactive agency")
	}

	if !toAgency.IsLicenseValid() {
		return fmt.Errorf("cannot transfer agent to agency with invalid license")
	}

	// Get and validate agent
	agent, err := s.userRepo.GetByID(agentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	if agent.Role != domain.RoleAgent {
		return fmt.Errorf("user is not an agent")
	}

	if agent.AgencyID == nil || *agent.AgencyID != fromAgencyID {
		return fmt.Errorf("agent does not belong to source agency")
	}

	// Transfer agent
	if err := agent.SetAgency(toAgencyID); err != nil {
		return fmt.Errorf("failed to transfer agent: %w", err)
	}

	// Update in database
	if err := s.userRepo.Update(agent); err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	s.logger.Printf("Agent %s transferred from %s to %s", agent.Name, fromAgency.Name, toAgency.Name)
	return nil
}

// GetAgencyAgents gets all agents belonging to an agency
func (s *AgencyService) GetAgencyAgents(agencyID string) ([]*domain.User, error) {
	if agencyID == "" {
		return nil, fmt.Errorf("agency ID cannot be empty")
	}

	// Verify agency exists
	if _, err := s.agencyRepo.GetByID(agencyID); err != nil {
		return nil, fmt.Errorf("agency not found: %w", err)
	}

	// Get agents by agency
	agents, err := s.userRepo.GetByAgency(agencyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agency agents: %w", err)
	}

	return agents, nil
}

// SetLicenseNumber sets the license number for an agency (simplified version)
func (s *AgencyService) SetLicenseNumber(agencyID, licenseNumber string) error {
	if agencyID == "" || licenseNumber == "" {
		return fmt.Errorf("agency ID and license number are required")
	}

	agency, err := s.agencyRepo.GetByID(agencyID)
	if err != nil {
		return fmt.Errorf("agency not found: %w", err)
	}

	// Update license number
	agency.LicenseNumber = licenseNumber
	agency.UpdatedAt = time.Now()

	if err := s.agencyRepo.Update(agency); err != nil {
		return fmt.Errorf("failed to update agency license: %w", err)
	}

	s.logger.Printf("License number updated for agency: %s", agency.Name)
	return nil
}