package servicio

import (
	"fmt"
	"realty-core/internal/domain"
	"realty-core/internal/repositorio"
	"strings"
)

// PropertyService handles property business logic
type PropertyService struct {
	propertyRepo          *repositorio.PropertyRepository
	realEstateCompanyRepo *repositorio.RealEstateCompanyRepository
}

// NewPropertyService creates a new property service instance
func NewPropertyService(propertyRepo *repositorio.PropertyRepository, realEstateCompanyRepo *repositorio.RealEstateCompanyRepository) *PropertyService {
	return &PropertyService{
		propertyRepo:          propertyRepo,
		realEstateCompanyRepo: realEstateCompanyRepo,
	}
}

// Create creates a new property with validation
func (s *PropertyService) Create(title, description, province, city, propertyType string, price float64) (*domain.Property, error) {
	// Create property with auto-generated ID and slug
	property := domain.NewProperty(title, description, province, city, propertyType, price)

	// Validate the property
	if err := property.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if slug already exists and make it unique if needed
	if err := s.ensureUniqueSlug(property); err != nil {
		return nil, fmt.Errorf("error ensuring unique slug: %w", err)
	}

	// Save to database
	if err := s.propertyRepo.Create(property); err != nil {
		return nil, fmt.Errorf("error creating property: %w", err)
	}

	return property, nil
}

// GetByID retrieves a property by ID
func (s *PropertyService) GetByID(id string) (*domain.Property, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("property ID is required")
	}

	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return property, nil
}

// GetBySlug retrieves a property by slug and increments view count
func (s *PropertyService) GetBySlug(slug string) (*domain.Property, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, fmt.Errorf("property slug is required")
	}

	property, err := s.propertyRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Increment view count
	if err := s.propertyRepo.IncrementViewCount(property.ID); err != nil {
		// Log error but don't fail the request for view count
		fmt.Printf("Warning: failed to increment view count for property %s: %v\n", property.ID, err)
	} else {
		property.IncrementViewCount()
	}

	return property, nil
}

// GetAll retrieves all properties with optional filters
func (s *PropertyService) GetAll(filters map[string]interface{}) ([]*domain.Property, error) {
	return s.propertyRepo.GetAll(filters)
}

// GetAvailable retrieves all available properties
func (s *PropertyService) GetAvailable() ([]*domain.Property, error) {
	return s.propertyRepo.GetAvailable()
}

// GetFeatured retrieves all featured properties
func (s *PropertyService) GetFeatured() ([]*domain.Property, error) {
	return s.propertyRepo.GetFeatured()
}

// Update updates an existing property
func (s *PropertyService) Update(id, title, description, province, city, propertyType string, price float64) (*domain.Property, error) {
	// Get existing property
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	property.Title = strings.TrimSpace(title)
	property.Description = strings.TrimSpace(description)
	property.Province = strings.TrimSpace(province)
	property.City = strings.TrimSpace(city)
	property.Type = strings.TrimSpace(propertyType)
	property.Price = price

	// Update slug if title changed
	if property.Title != title {
		property.UpdateSlug()
		if err := s.ensureUniqueSlug(property); err != nil {
			return nil, fmt.Errorf("error ensuring unique slug: %w", err)
		}
	}

	// Validate updated property
	if err := property.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Calculate price per m2 if area is set
	property.CalculatePricePerM2()

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating property: %w", err)
	}

	return property, nil
}

// UpdateDetails updates detailed property information
func (s *PropertyService) UpdateDetails(id string, bedrooms int, bathrooms float32, areaM2 float64, yearBuilt *int, floors *int) (*domain.Property, error) {
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	property.Bedrooms = bedrooms
	property.Bathrooms = bathrooms
	property.AreaM2 = areaM2
	property.YearBuilt = yearBuilt
	property.Floors = floors

	// Validate
	if err := property.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Calculate price per m2
	property.CalculatePricePerM2()

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating property details: %w", err)
	}

	return property, nil
}

// SetLocation sets GPS coordinates for a property
func (s *PropertyService) SetLocation(id string, latitude, longitude float64, precision string) (*domain.Property, error) {
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Set location
	if err := property.SetLocation(latitude, longitude, precision); err != nil {
		return nil, fmt.Errorf("invalid location: %w", err)
	}

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating property location: %w", err)
	}

	return property, nil
}

// AddImage adds an image to a property's gallery
func (s *PropertyService) AddImage(id, imageURL string) (*domain.Property, error) {
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Add image
	if err := property.AddImage(imageURL); err != nil {
		return nil, fmt.Errorf("error adding image: %w", err)
	}

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating property images: %w", err)
	}

	return property, nil
}

// SetMainImage sets the main image for a property
func (s *PropertyService) SetMainImage(id, imageURL string) (*domain.Property, error) {
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate image URL format
	if !isValidImageURL(imageURL) {
		return nil, fmt.Errorf("invalid image URL format")
	}

	property.MainImage = imageURL
	property.UpdateTimestamp()

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating main image: %w", err)
	}

	return property, nil
}

// AddTag adds a search tag to a property
func (s *PropertyService) AddTag(id, tag string) (*domain.Property, error) {
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Add tag
	if err := property.AddTag(tag); err != nil {
		return nil, fmt.Errorf("error adding tag: %w", err)
	}

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating property tags: %w", err)
	}

	return property, nil
}

// SetFeatured marks/unmarks a property as featured
func (s *PropertyService) SetFeatured(id string, featured bool) (*domain.Property, error) {
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	property.SetFeatured(featured)

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating featured status: %w", err)
	}

	return property, nil
}

// ChangeStatus changes the status of a property (available, sold, rented, reserved)
func (s *PropertyService) ChangeStatus(id, status string) (*domain.Property, error) {
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	property.Status = status

	// Validate status
	if err := property.ValidateStatus(); err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	property.UpdateTimestamp()

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating property status: %w", err)
	}

	return property, nil
}

// AssignToCompany assigns a property to a real estate company
func (s *PropertyService) AssignToCompany(propertyID, companyID string) (*domain.Property, error) {
	// Validate company exists and is active
	company, err := s.realEstateCompanyRepo.GetByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("real estate company not found: %w", err)
	}

	if !company.Active {
		return nil, fmt.Errorf("cannot assign property to inactive company")
	}

	// Get property
	property, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		return nil, err
	}

	// Assign company
	property.RealEstateCompanyID = &companyID
	property.UpdateTimestamp()

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error assigning property to company: %w", err)
	}

	return property, nil
}

// UnassignFromCompany removes a property from a real estate company
func (s *PropertyService) UnassignFromCompany(propertyID string) (*domain.Property, error) {
	property, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		return nil, err
	}

	// Remove company assignment
	property.RealEstateCompanyID = nil
	property.UpdateTimestamp()

	// Save changes
	if err := s.propertyRepo.Update(property); err != nil {
		return nil, fmt.Errorf("error unassigning property from company: %w", err)
	}

	return property, nil
}

// GetByCompany retrieves all properties for a specific company
func (s *PropertyService) GetByCompany(companyID string) ([]*domain.Property, error) {
	// Validate company exists
	_, err := s.realEstateCompanyRepo.GetByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("real estate company not found: %w", err)
	}

	return s.propertyRepo.GetByCompany(companyID)
}

// Search performs full-text search on properties
func (s *PropertyService) Search(searchTerm string) ([]*domain.Property, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return nil, fmt.Errorf("search term is required")
	}

	return s.propertyRepo.Search(searchTerm)
}

// Delete deletes a property
func (s *PropertyService) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("property ID is required")
	}

	return s.propertyRepo.Delete(id)
}

// GetPropertiesInBudgetRange gets properties within a price range
func (s *PropertyService) GetPropertiesInBudgetRange(minPrice, maxPrice *float64) ([]*domain.Property, error) {
	filters := make(map[string]interface{})
	filters["status"] = domain.PropertyStatusAvailable

	if minPrice != nil && *minPrice > 0 {
		filters["min_price"] = *minPrice
	}
	if maxPrice != nil && *maxPrice > 0 {
		filters["max_price"] = *maxPrice
	}

	return s.propertyRepo.GetAll(filters)
}

// GetStatistics returns property statistics
func (s *PropertyService) GetStatistics() (map[string]interface{}, error) {
	allProperties, err := s.propertyRepo.GetAll(nil)
	if err != nil {
		return nil, fmt.Errorf("error getting properties for statistics: %w", err)
	}

	stats := map[string]interface{}{
		"total_properties": len(allProperties),
		"available":        0,
		"sold":             0,
		"rented":           0,
		"reserved":         0,
		"featured":         0,
		"avg_price":        0.0,
		"total_value":      0.0,
	}

	if len(allProperties) == 0 {
		return stats, nil
	}

	var totalPrice float64
	for _, property := range allProperties {
		totalPrice += property.Price

		switch property.Status {
		case domain.PropertyStatusAvailable:
			stats["available"] = stats["available"].(int) + 1
		case domain.PropertyStatusSold:
			stats["sold"] = stats["sold"].(int) + 1
		case domain.PropertyStatusRented:
			stats["rented"] = stats["rented"].(int) + 1
		case domain.PropertyStatusReserved:
			stats["reserved"] = stats["reserved"].(int) + 1
		}

		if property.Featured {
			stats["featured"] = stats["featured"].(int) + 1
		}
	}

	stats["avg_price"] = totalPrice / float64(len(allProperties))
	stats["total_value"] = totalPrice

	return stats, nil
}

// ValidateTitle validates a property title
func (s *PropertyService) ValidateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title is required")
	}
	if len(title) > 255 {
		return fmt.Errorf("title cannot exceed 255 characters")
	}
	return nil
}

// ValidatePrice validates a property price
func (s *PropertyService) ValidatePrice(price float64) error {
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if price > 10000000 { // 10 million USD max
		return fmt.Errorf("price cannot exceed $10,000,000")
	}
	return nil
}

// ensureUniqueSlug ensures the property has a unique slug
func (s *PropertyService) ensureUniqueSlug(property *domain.Property) error {
	originalSlug := property.Slug
	counter := 1

	for {
		// Check if slug exists
		_, err := s.propertyRepo.GetBySlug(property.Slug)
		if err != nil {
			// Slug doesn't exist, we can use it
			break
		}

		// Slug exists, try with counter
		property.Slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++

		if counter > 100 {
			return fmt.Errorf("unable to generate unique slug after 100 attempts")
		}
	}

	return nil
}

// Helper function to validate image URLs
func isValidImageURL(url string) bool {
	if url == "" {
		return false
	}
	// Basic URL validation for images
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
