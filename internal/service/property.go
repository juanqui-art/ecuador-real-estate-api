package service

import (
	"fmt"
	"strings"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// PropertyServiceInterface defines the business logic operations for properties
type PropertyServiceInterface interface {
	CreateProperty(title, description, province, city, propertyType string, price float64) (*domain.Property, error)
	GetProperty(id string) (*domain.Property, error)
	GetPropertyBySlug(slug string) (*domain.Property, error)
	ListProperties() ([]domain.Property, error)
	UpdateProperty(id, title, description, province, city, propertyType string, price float64) (*domain.Property, error)
	DeleteProperty(id string) error
	FilterByProvince(province string) ([]domain.Property, error)
	FilterByPriceRange(minPrice, maxPrice float64) ([]domain.Property, error)
	GetStatistics() (map[string]interface{}, error)
	SetPropertyLocation(id string, latitude, longitude float64, precision string) error
	SetPropertyFeatured(id string, featured bool) error
	AddPropertyTag(id, tag string) error
	SearchProperties(query string) ([]domain.Property, error)
	// Enhanced search methods
	SearchPropertiesRanked(query string, limit int) ([]repository.PropertySearchResult, error)
	GetSearchSuggestions(query string, limit int) ([]repository.SearchSuggestion, error)
	AdvancedSearch(params repository.AdvancedSearchParams) ([]repository.PropertySearchResult, error)
}

// PropertyService handles business logic for properties
type PropertyService struct {
	repo repository.PropertyRepository
}

// NewPropertyService creates a new instance of the service
func NewPropertyService(repo repository.PropertyRepository) *PropertyService {
	return &PropertyService{repo: repo}
}

// CreateProperty creates a new property with validations
func (s *PropertyService) CreateProperty(title, description, province, city, propertyType string, price float64) (*domain.Property, error) {
	// Validate input data
	if err := s.validatePropertyData(title, province, city, propertyType, price); err != nil {
		return nil, err
	}

	// Clean and normalize data
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	province = strings.TrimSpace(province)
	city = strings.TrimSpace(city)
	propertyType = strings.ToLower(strings.TrimSpace(propertyType))

	// Create the property
	property := domain.NewProperty(title, description, province, city, propertyType, price)

	// Validate the complete property
	if !property.IsValid() {
		return nil, fmt.Errorf("invalid property data")
	}

	// Save to database
	if err := s.repo.Create(property); err != nil {
		return nil, fmt.Errorf("error creating property: %w", err)
	}

	return property, nil
}

// GetProperty retrieves a property by ID
func (s *PropertyService) GetProperty(id string) (*domain.Property, error) {
	if id == "" {
		return nil, fmt.Errorf("property ID required")
	}

	property, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving property: %w", err)
	}

	// Increment view count
	property.IncrementViews()
	s.repo.Update(property)

	return property, nil
}

// GetPropertyBySlug retrieves a property by SEO slug
func (s *PropertyService) GetPropertyBySlug(slug string) (*domain.Property, error) {
	if slug == "" {
		return nil, fmt.Errorf("property slug required")
	}

	// Validate slug format
	if !domain.IsValidSlug(slug) {
		return nil, fmt.Errorf("invalid slug format: %s", slug)
	}

	property, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("error retrieving property by slug: %w", err)
	}

	// Increment view count
	property.IncrementViews()
	s.repo.Update(property)

	return property, nil
}

// ListProperties retrieves all properties
func (s *PropertyService) ListProperties() ([]domain.Property, error) {
	properties, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error listing properties: %w", err)
	}

	return properties, nil
}

// UpdateProperty modifies an existing property
func (s *PropertyService) UpdateProperty(id, title, description, province, city, propertyType string, price float64) (*domain.Property, error) {
	// Check if property exists
	property, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("property not found: %w", err)
	}

	// Validate new data
	if err := s.validatePropertyData(title, province, city, propertyType, price); err != nil {
		return nil, err
	}

	// Update fields
	property.Title = strings.TrimSpace(title)
	property.Description = strings.TrimSpace(description)
	property.Province = strings.TrimSpace(province)
	property.City = strings.TrimSpace(city)
	property.Type = strings.ToLower(strings.TrimSpace(propertyType))
	property.Price = price

	// Update slug if title changed
	property.UpdateSlug()

	// Validate updated property
	if !property.IsValid() {
		return nil, fmt.Errorf("invalid updated property data")
	}

	// Save changes
	if err := s.repo.Update(property); err != nil {
		return nil, fmt.Errorf("error updating property: %w", err)
	}

	return property, nil
}

// DeleteProperty removes a property by ID
func (s *PropertyService) DeleteProperty(id string) error {
	if id == "" {
		return fmt.Errorf("property ID required")
	}

	// Verify property exists before deleting
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("property not found: %w", err)
	}

	// Delete the property
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("error deleting property: %w", err)
	}

	return nil
}

// FilterByProvince filters properties by province
func (s *PropertyService) FilterByProvince(province string) ([]domain.Property, error) {
	if province == "" {
		return nil, fmt.Errorf("province required")
	}

	// Validate province
	if !domain.IsValidProvince(province) {
		return nil, fmt.Errorf("invalid province: %s", province)
	}

	properties, err := s.repo.GetByProvince(province)
	if err != nil {
		return nil, fmt.Errorf("error filtering properties by province: %w", err)
	}

	return properties, nil
}

// FilterByPriceRange filters properties by price range
func (s *PropertyService) FilterByPriceRange(minPrice, maxPrice float64) ([]domain.Property, error) {
	if minPrice < 0 || maxPrice < 0 {
		return nil, fmt.Errorf("prices must be positive")
	}

	if minPrice > maxPrice {
		return nil, fmt.Errorf("minimum price cannot be greater than maximum price")
	}

	properties, err := s.repo.GetByPriceRange(minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("error filtering properties by price range: %w", err)
	}

	return properties, nil
}

// GetStatistics returns basic property statistics
func (s *PropertyService) GetStatistics() (map[string]interface{}, error) {
	properties, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving properties: %w", err)
	}

	stats := make(map[string]interface{})
	stats["total_properties"] = len(properties)

	// Count by type
	typeCount := make(map[string]int)
	// Count by status
	statusCount := make(map[string]int)
	// Count by province
	provinceCount := make(map[string]int)
	// Calculate average price
	var totalPrice float64

	for _, property := range properties {
		typeCount[property.Type]++
		statusCount[property.Status]++
		provinceCount[property.Province]++
		totalPrice += property.Price
	}

	stats["by_type"] = typeCount
	stats["by_status"] = statusCount
	stats["by_province"] = provinceCount

	if len(properties) > 0 {
		stats["average_price"] = totalPrice / float64(len(properties))
	} else {
		stats["average_price"] = float64(0)
	}

	return stats, nil
}

// SetPropertyLocation sets GPS coordinates for a property
func (s *PropertyService) SetPropertyLocation(id string, latitude, longitude float64, precision string) error {
	property, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("property not found: %w", err)
	}

	if err := property.SetLocation(latitude, longitude, precision); err != nil {
		return fmt.Errorf("error setting location: %w", err)
	}

	if err := s.repo.Update(property); err != nil {
		return fmt.Errorf("error updating property location: %w", err)
	}

	return nil
}

// SetPropertyFeatured marks or unmarks a property as featured
func (s *PropertyService) SetPropertyFeatured(id string, featured bool) error {
	property, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("property not found: %w", err)
	}

	property.SetFeatured(featured)

	if err := s.repo.Update(property); err != nil {
		return fmt.Errorf("error updating property featured status: %w", err)
	}

	return nil
}

// AddPropertyTag adds a search tag to a property
func (s *PropertyService) AddPropertyTag(id, tag string) error {
	property, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("property not found: %w", err)
	}

	property.AddTag(tag)

	if err := s.repo.Update(property); err != nil {
		return fmt.Errorf("error adding tag to property: %w", err)
	}

	return nil
}

// SearchProperties performs PostgreSQL full-text search
func (s *PropertyService) SearchProperties(query string) ([]domain.Property, error) {
	if query == "" {
		return s.repo.GetAll()
	}

	// Clean and validate search query
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return nil, fmt.Errorf("search query must be at least 2 characters")
	}

	// Use PostgreSQL FTS for efficient search
	properties, err := s.repo.SearchProperties(query, 50)
	if err != nil {
		return nil, fmt.Errorf("error performing search: %w", err)
	}

	return properties, nil
}

// SearchPropertiesRanked performs ranked full-text search with relevance scores
func (s *PropertyService) SearchPropertiesRanked(query string, limit int) ([]repository.PropertySearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query required")
	}

	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return nil, fmt.Errorf("search query must be at least 2 characters")
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	results, err := s.repo.SearchPropertiesRanked(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error performing ranked search: %w", err)
	}

	return results, nil
}

// GetSearchSuggestions returns autocomplete suggestions for search
func (s *PropertyService) GetSearchSuggestions(query string, limit int) ([]repository.SearchSuggestion, error) {
	if query == "" {
		return []repository.SearchSuggestion{}, nil
	}

	query = strings.TrimSpace(query)
	if len(query) < 1 {
		return []repository.SearchSuggestion{}, nil
	}

	if limit <= 0 || limit > 20 {
		limit = 10
	}

	suggestions, err := s.repo.GetSearchSuggestions(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting search suggestions: %w", err)
	}

	return suggestions, nil
}

// AdvancedSearch performs advanced search with multiple filters
func (s *PropertyService) AdvancedSearch(params repository.AdvancedSearchParams) ([]repository.PropertySearchResult, error) {
	// Validate parameters
	if params.MinPrice < 0 || params.MaxPrice < 0 {
		return nil, fmt.Errorf("prices must be positive")
	}

	if params.MinPrice > 0 && params.MaxPrice > 0 && params.MinPrice > params.MaxPrice {
		return nil, fmt.Errorf("minimum price cannot be greater than maximum price")
	}

	if params.MinBedrooms < 0 || params.MaxBedrooms < 0 {
		return nil, fmt.Errorf("bedroom counts must be positive")
	}

	if params.MinBathrooms < 0 || params.MaxBathrooms < 0 {
		return nil, fmt.Errorf("bathroom counts must be positive")
	}

	if params.MinArea < 0 || params.MaxArea < 0 {
		return nil, fmt.Errorf("area values must be positive")
	}

	// Validate province if provided
	if params.Province != "" && !domain.IsValidProvince(params.Province) {
		return nil, fmt.Errorf("invalid province: %s", params.Province)
	}

	// Validate property type if provided
	if params.Type != "" && !domain.IsValidPropertyType(params.Type) {
		return nil, fmt.Errorf("invalid property type: %s", params.Type)
	}

	// Clean search query
	if params.Query != "" {
		params.Query = strings.TrimSpace(params.Query)
		if len(params.Query) < 2 {
			return nil, fmt.Errorf("search query must be at least 2 characters")
		}
	}

	// Set reasonable limits
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 50
	}

	results, err := s.repo.AdvancedSearch(params)
	if err != nil {
		return nil, fmt.Errorf("error performing advanced search: %w", err)
	}

	return results, nil
}

// validatePropertyData validates basic property creation/update data
func (s *PropertyService) validatePropertyData(title, province, city, propertyType string, price float64) error {
	// Validate required fields
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title is required")
	}

	if len(strings.TrimSpace(title)) < 10 {
		return fmt.Errorf("title must be at least 10 characters")
	}

	if len(strings.TrimSpace(title)) > 255 {
		return fmt.Errorf("title cannot exceed 255 characters")
	}

	if strings.TrimSpace(province) == "" {
		return fmt.Errorf("province is required")
	}

	if strings.TrimSpace(city) == "" {
		return fmt.Errorf("city is required")
	}

	if strings.TrimSpace(propertyType) == "" {
		return fmt.Errorf("property type is required")
	}

	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	// Validate Ecuadorian province
	if !domain.IsValidProvince(province) {
		return fmt.Errorf("invalid province: %s", province)
	}

	// Validate property type
	if !domain.IsValidPropertyType(strings.ToLower(strings.TrimSpace(propertyType))) {
		return fmt.Errorf("invalid property type: %s. Valid types: house, apartment, land, commercial", propertyType)
	}

	return nil
}