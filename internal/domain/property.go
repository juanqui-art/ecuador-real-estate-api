package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Property represents a real estate property matching the PostgreSQL 'properties' table
type Property struct {
	ID                    string    `json:"id" db:"id"`
	Slug                  string    `json:"slug" db:"slug"`
	Title                 string    `json:"title" db:"title"`
	Description           string    `json:"description" db:"description"`
	Price                 float64   `json:"price" db:"price"`
	Province              string    `json:"province" db:"province"`
	City                  string    `json:"city" db:"city"`
	Sector                *string   `json:"sector" db:"sector"`
	Address               *string   `json:"address" db:"address"`
	Latitude              *float64  `json:"latitude" db:"latitude"`
	Longitude             *float64  `json:"longitude" db:"longitude"`
	LocationPrecision     string    `json:"location_precision" db:"location_precision"`
	Type                  string    `json:"type" db:"type"`
	Status                string    `json:"status" db:"status"`
	Bedrooms              int       `json:"bedrooms" db:"bedrooms"`
	Bathrooms             float32   `json:"bathrooms" db:"bathrooms"`
	AreaM2                float64   `json:"area_m2" db:"area_m2"`
	ParkingSpaces         int       `json:"parking_spaces" db:"parking_spaces"`
	MainImage             *string   `json:"main_image" db:"main_image"`
	Images                []string  `json:"images" db:"images"`
	VideoTour             *string   `json:"video_tour" db:"video_tour"`
	Tour360               *string   `json:"tour_360" db:"tour_360"`
	RentPrice             *float64  `json:"rent_price" db:"rent_price"`
	CommonExpenses        *float64  `json:"common_expenses" db:"common_expenses"`
	PricePerM2            *float64  `json:"price_per_m2" db:"price_per_m2"`
	YearBuilt             *int      `json:"year_built" db:"year_built"`
	Floors                *int      `json:"floors" db:"floors"`
	PropertyStatus        string    `json:"property_status" db:"property_status"`
	Furnished             bool      `json:"furnished" db:"furnished"`
	Garage                bool      `json:"garage" db:"garage"`
	Pool                  bool      `json:"pool" db:"pool"`
	Garden                bool      `json:"garden" db:"garden"`
	Terrace               bool      `json:"terrace" db:"terrace"`
	Balcony               bool      `json:"balcony" db:"balcony"`
	Security              bool      `json:"security" db:"security"`
	Elevator              bool      `json:"elevator" db:"elevator"`
	AirConditioning       bool      `json:"air_conditioning" db:"air_conditioning"`
	Tags                  []string  `json:"tags" db:"tags"`
	Featured              bool      `json:"featured" db:"featured"`
	ViewCount             int       `json:"view_count" db:"view_count"`
	RealEstateCompanyID   *string   `json:"real_estate_company_id" db:"real_estate_company_id"`
	// User relationships for role-based system
	OwnerID               *string   `json:"owner_id" db:"owner_id"`
	AgentID               *string   `json:"agent_id" db:"agent_id"`
	AgencyID              *string   `json:"agency_id" db:"agency_id"`
	CreatedBy             *string   `json:"created_by" db:"created_by"`
	UpdatedBy             *string   `json:"updated_by" db:"updated_by"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// NewProperty creates a new property with automatically generated SEO slug
func NewProperty(title, description, province, city, propertyType string, price float64, ownerID string) *Property {
	id := uuid.New().String()
	slug := GenerateSlug(title, id)

	return &Property{
		ID:                id,
		Slug:              slug,
		Title:             title,
		Description:       description,
		Price:             price,
		Province:          province,
		City:              city,
		Type:              propertyType,
		Status:            StatusAvailable,
		LocationPrecision: PrecisionApproximate,
		PropertyStatus:    PropertyStatusUsed,
		Bedrooms:          0,
		Bathrooms:         0,
		AreaM2:            0,
		ParkingSpaces:     0,
		Images:            []string{},
		Tags:              []string{},
		Featured:          false,
		ViewCount:         0,
		Furnished:         false,
		Garage:            false,
		Pool:              false,
		Garden:            false,
		Terrace:           false,
		Balcony:           false,
		Security:          false,
		Elevator:          false,
		AirConditioning:   false,
		Sector:            nil,
		Address:           nil,
		OwnerID:           &ownerID,
		CreatedBy:         &ownerID,
		MainImage:         nil,
		VideoTour:         nil,
		Tour360:           nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// IsValid validates the required fields of the property
func (p *Property) IsValid() bool {
	return p.Title != "" &&
		p.Price > 0 &&
		p.Province != "" &&
		p.City != "" &&
		p.Type != "" &&
		p.ParkingSpaces >= 0
}

// UpdateTimestamp updates the modification date
func (p *Property) UpdateTimestamp() {
	p.UpdatedAt = time.Now()
}

// SetLocation sets the GPS coordinates of the property
func (p *Property) SetLocation(latitude, longitude float64, precision string) error {
	if !IsValidEcuadorCoordinates(latitude, longitude) {
		return fmt.Errorf("coordinates outside Ecuadorian territory: lat=%.6f, lng=%.6f", latitude, longitude)
	}

	if !IsValidLocationPrecision(precision) {
		return fmt.Errorf("invalid location precision: %s", precision)
	}

	p.Latitude = &latitude
	p.Longitude = &longitude
	p.LocationPrecision = precision
	p.UpdateTimestamp()

	return nil
}

// UpdateSlug regenerates the slug when title changes
func (p *Property) UpdateSlug() {
	p.Slug = GenerateSlug(p.Title, p.ID)
	p.UpdateTimestamp()
}

// IncrementViews increments the view counter
func (p *Property) IncrementViews() {
	p.ViewCount++
}

// AddTag adds a search tag
func (p *Property) AddTag(tag string) {
	if tag != "" && !p.HasTag(tag) {
		p.Tags = append(p.Tags, strings.ToLower(strings.TrimSpace(tag)))
		p.UpdateTimestamp()
	}
}

// HasTag checks if property has a specific tag
func (p *Property) HasTag(tag string) bool {
	tagLower := strings.ToLower(strings.TrimSpace(tag))
	for _, t := range p.Tags {
		if t == tagLower {
			return true
		}
	}
	return false
}

// SetFeatured marks the property as featured
func (p *Property) SetFeatured(featured bool) {
	p.Featured = featured
	p.UpdateTimestamp()
}

// SetParkingSpaces sets the number of parking spaces for the property
func (p *Property) SetParkingSpaces(parkingSpaces int) error {
	if parkingSpaces < 0 {
		return fmt.Errorf("parking spaces must be non-negative")
	}
	p.ParkingSpaces = parkingSpaces
	p.UpdateTimestamp()
	return nil
}

// Constants for property types
const (
	TypeHouse      = "house"
	TypeApartment  = "apartment"
	TypeLand       = "land"
	TypeCommercial = "commercial"
)

// Constants for property status
const (
	StatusAvailable = "available"
	StatusSold      = "sold"
	StatusRented    = "rented"
	StatusReserved  = "reserved"
)

// Constants for location precision
const (
	PrecisionExact       = "exact"
	PrecisionApproximate = "approximate"
	PrecisionSector      = "sector"
)

// Constants for property condition
const (
	PropertyStatusNew       = "new"
	PropertyStatusUsed      = "used"
	PropertyStatusRenovated = "renovated"
)

// EcuadorProvinces lists the valid provinces of Ecuador
var EcuadorProvinces = []string{
	"Azuay", "Bolívar", "Cañar", "Carchi", "Chimborazo",
	"Cotopaxi", "El Oro", "Esmeraldas", "Galápagos",
	"Guayas", "Imbabura", "Loja", "Los Ríos", "Manabí",
	"Morona Santiago", "Napo", "Orellana", "Pastaza",
	"Pichincha", "Santa Elena", "Santo Domingo",
	"Sucumbíos", "Tungurahua", "Zamora Chinchipe",
}

// IsValidProvince verifies if a province is valid in Ecuador
func IsValidProvince(province string) bool {
	for _, p := range EcuadorProvinces {
		if p == province {
			return true
		}
	}
	return false
}

// IsValidEcuadorCoordinates verifies if coordinates are within Ecuador
func IsValidEcuadorCoordinates(latitude, longitude float64) bool {
	// Approximate geographic limits of Ecuador
	// Latitude: -5.0 (south) to 2.0 (north)
	// Longitude: -92.0 (west) to -75.0 (east)
	validLatitude := latitude >= -5.0 && latitude <= 2.0
	validLongitude := longitude >= -92.0 && longitude <= -75.0
	return validLatitude && validLongitude
}

// IsValidLocationPrecision verifies if the location precision is valid
func IsValidLocationPrecision(precision string) bool {
	validPrecisions := []string{PrecisionExact, PrecisionApproximate, PrecisionSector}
	for _, p := range validPrecisions {
		if p == precision {
			return true
		}
	}
	return false
}

// IsValidPropertyType verifies if the property type is valid
func IsValidPropertyType(propertyType string) bool {
	validTypes := []string{TypeHouse, TypeApartment, TypeLand, TypeCommercial}
	for _, t := range validTypes {
		if t == propertyType {
			return true
		}
	}
	return false
}

// IsValidPropertyStatus verifies if the property status is valid
func IsValidPropertyStatus(status string) bool {
	validStatuses := []string{PropertyStatusNew, PropertyStatusUsed, PropertyStatusRenovated}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// GenerateSlug creates an SEO-friendly slug from the property title
func GenerateSlug(title, id string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace special characters and spaces with hyphens
	// Keep only letters, numbers and spaces
	slug = regexp.MustCompile(`[^a-z0-9\s]+`).ReplaceAllString(slug, "")

	// Replace multiple spaces with single space
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, " ")

	// Convert spaces to hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove hyphens at beginning and end
	slug = strings.Trim(slug, "-")

	// Truncate if too long (maximum 50 characters before ID)
	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.Trim(slug, "-")
	}

	// Add short ID at the end to avoid duplicates
	shortID := id
	if len(id) > 8 {
		shortID = id[:8]
	}

	return slug + "-" + shortID
}

// IsValidSlug verifies if a string can be a valid slug
func IsValidSlug(slug string) bool {
	if slug == "" {
		return false
	}
	// A valid slug only contains lowercase letters, numbers and hyphens
	// Cannot start or end with hyphen
	matched, _ := regexp.MatchString(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`, slug)
	return matched
}

// AssignToAgency assigns the property to an agency and optionally an agent
func (p *Property) AssignToAgency(agencyID string, agentID *string, userID string) error {
	if agencyID == "" {
		return fmt.Errorf("agency ID cannot be empty")
	}

	p.AgencyID = &agencyID
	p.AgentID = agentID
	p.UpdatedBy = &userID
	p.UpdatedAt = time.Now()
	
	return nil
}

// RemoveFromAgency removes the property from agency management
func (p *Property) RemoveFromAgency(userID string) error {
	p.AgencyID = nil
	p.AgentID = nil
	p.UpdatedBy = &userID
	p.UpdatedAt = time.Now()
	
	return nil
}

// AssignToAgent assigns the property to a specific agent
func (p *Property) AssignToAgent(agentID string, userID string) error {
	if agentID == "" {
		return fmt.Errorf("agent ID cannot be empty")
	}

	p.AgentID = &agentID
	p.UpdatedBy = &userID
	p.UpdatedAt = time.Now()
	
	return nil
}

// RemoveFromAgent removes the property from agent assignment
func (p *Property) RemoveFromAgent(userID string) error {
	p.AgentID = nil
	p.UpdatedBy = &userID
	p.UpdatedAt = time.Now()
	
	return nil
}

// TransferOwnership transfers the property to a new owner
func (p *Property) TransferOwnership(newOwnerID string, userID string) error {
	if newOwnerID == "" {
		return fmt.Errorf("new owner ID cannot be empty")
	}

	p.OwnerID = &newOwnerID
	p.UpdatedBy = &userID
	p.UpdatedAt = time.Now()
	
	return nil
}

// IsOwnedBy checks if the property is owned by a specific user
func (p *Property) IsOwnedBy(userID string) bool {
	return p.OwnerID != nil && *p.OwnerID == userID
}

// IsManagedByAgency checks if the property is managed by a specific agency
func (p *Property) IsManagedByAgency(agencyID string) bool {
	return p.AgencyID != nil && *p.AgencyID == agencyID
}

// IsAssignedToAgent checks if the property is assigned to a specific agent
func (p *Property) IsAssignedToAgent(agentID string) bool {
	return p.AgentID != nil && *p.AgentID == agentID
}

// GetOwnerID returns the owner ID if available
func (p *Property) GetOwnerID() *string {
	return p.OwnerID
}

// GetAgencyID returns the agency ID if available
func (p *Property) GetAgencyID() *string {
	return p.AgencyID
}

// GetAgentID returns the agent ID if available
func (p *Property) GetAgentID() *string {
	return p.AgentID
}

// CanBeModifiedBy checks if a user can modify this property based on role relationships
func (p *Property) CanBeModifiedBy(userID string, userRole UserRole, userAgencyID *string) bool {
	// Admin can modify any property
	if userRole == RoleAdmin {
		return true
	}

	// Owner can modify their own property
	if userRole == RoleOwner && p.IsOwnedBy(userID) {
		return true
	}

	// Agency can modify properties they manage
	if userRole == RoleAgency && userAgencyID != nil && p.IsManagedByAgency(*userAgencyID) {
		return true
	}

	// Agent can modify properties assigned to them within their agency
	if userRole == RoleAgent && userAgencyID != nil && 
		p.IsManagedByAgency(*userAgencyID) && p.IsAssignedToAgent(userID) {
		return true
	}

	return false
}

// CanBeViewedBy checks if a user can view this property (all users can view, but some have additional access)
func (p *Property) CanBeViewedBy(userID string, userRole UserRole, userAgencyID *string) bool {
	// All users can view properties
	return true
}

// UpdateBy updates the UpdatedBy field and timestamp
func (p *Property) UpdateBy(userID string) {
	p.UpdatedBy = &userID
	p.UpdatedAt = time.Now()
}

// PaginationParams represents pagination parameters for queries
type PaginationParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
}

// NewPaginationParams creates default pagination parameters
func NewPaginationParams() *PaginationParams {
	return &PaginationParams{
		Page:     1,
		PageSize: 20,
		SortBy:   "created_at",
		SortDesc: true,
	}
}

// GetOffset calculates the SQL OFFSET value based on page and page_size
func (p *PaginationParams) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the SQL LIMIT value
func (p *PaginationParams) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// GetOrderBy returns the SQL ORDER BY clause
func (p *PaginationParams) GetOrderBy() string {
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"title":      true,
		"price":      true,
		"area_m2":    true,
		"bedrooms":   true,
		"bathrooms":  true,
		"view_count": true,
	}

	if !validSortFields[p.SortBy] {
		p.SortBy = "created_at"
	}

	order := "ASC"
	if p.SortDesc {
		order = "DESC"
	}

	return fmt.Sprintf("%s %s", p.SortBy, order)
}

// Validate validates pagination parameters
func (p *PaginationParams) Validate() error {
	if p.Page <= 0 {
		return fmt.Errorf("page must be greater than 0")
	}
	if p.PageSize <= 0 {
		return fmt.Errorf("page_size must be greater than 0")
	}
	if p.PageSize > 100 {
		return fmt.Errorf("page_size cannot exceed 100")
	}
	return nil
}

// PaginatedResponse represents a paginated response with metadata
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	CurrentPage  int  `json:"current_page"`
	PageSize     int  `json:"page_size"`
	TotalPages   int  `json:"total_pages"`
	TotalRecords int  `json:"total_records"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

// NewPagination creates pagination metadata
func NewPagination(currentPage, pageSize, totalRecords int) *Pagination {
	totalPages := (totalRecords + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return &Pagination{
		CurrentPage:  currentPage,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasNext:      currentPage < totalPages,
		HasPrev:      currentPage > 1,
	}
}

// Role-based property management methods

// SetOwner sets the property owner
func (p *Property) SetOwner(ownerID string) {
	p.OwnerID = &ownerID
	p.UpdateTimestamp()
}

// SetAgent sets the property agent
func (p *Property) SetAgent(agentID string) {
	p.AgentID = &agentID
	p.UpdateTimestamp()
}

// SetAgency sets the property agency
func (p *Property) SetAgency(agencyID string) {
	p.AgencyID = &agencyID
	p.UpdateTimestamp()
}

// SetCreatedBy sets who created the property
func (p *Property) SetCreatedBy(userID string) {
	p.CreatedBy = &userID
}

// SetUpdatedBy sets who last updated the property
func (p *Property) SetUpdatedBy(userID string) {
	p.UpdatedBy = &userID
	p.UpdateTimestamp()
}


// IsAssignedToAgency checks if the property is assigned to a specific agency
func (p *Property) IsAssignedToAgency(agencyID string) bool {
	return p.AgencyID != nil && *p.AgencyID == agencyID
}

// ValidateRoleBasedRules validates role-based business rules
func (p *Property) ValidateRoleBasedRules() error {
	// If property has an agent, it must have an agency
	if p.AgentID != nil && p.AgencyID == nil {
		return fmt.Errorf("property with assigned agent must have an agency")
	}

	// Property must have either an owner or an agency
	if p.OwnerID == nil && p.AgencyID == nil {
		return fmt.Errorf("property must have either an owner or an agency")
	}

	return nil
}

// GetManagers returns all users who can manage this property
func (p *Property) GetManagers() []string {
	var managers []string
	
	if p.OwnerID != nil {
		managers = append(managers, *p.OwnerID)
	}
	if p.AgentID != nil {
		managers = append(managers, *p.AgentID)
	}
	if p.AgencyID != nil {
		managers = append(managers, *p.AgencyID)
	}
	
	return managers
}

// PropertyWithRelations represents a property with its related entities
type PropertyWithRelations struct {
	Property *Property `json:"property"`
	Owner    *User     `json:"owner,omitempty"`
	Agent    *User     `json:"agent,omitempty"`
	Agency   *Agency   `json:"agency,omitempty"`
}

// PropertySearchFilters represents enhanced search filters with role-based filtering
type PropertySearchFilters struct {
	// Basic filters
	Query         string   `json:"query"`
	MinPrice      *float64 `json:"min_price"`
	MaxPrice      *float64 `json:"max_price"`
	PropertyTypes []string `json:"property_types"`
	Provinces     []string `json:"provinces"`
	Cities        []string `json:"cities"`
	Sectors       []string `json:"sectors"`
	MinBedrooms   *int     `json:"min_bedrooms"`
	MaxBedrooms   *int     `json:"max_bedrooms"`
	MinBathrooms  *float32 `json:"min_bathrooms"`
	MaxBathrooms  *float32 `json:"max_bathrooms"`
	MinArea       *float64 `json:"min_area"`
	MaxArea       *float64 `json:"max_area"`
	Status        []string `json:"status"`
	Featured      *bool    `json:"featured"`
	
	// Role-based filters
	OwnerID   *string `json:"owner_id"`
	AgentID   *string `json:"agent_id"`
	AgencyID  *string `json:"agency_id"`
	CreatedBy *string `json:"created_by"`
	
	// Additional filters
	HasPool           *bool    `json:"has_pool"`
	HasGarden         *bool    `json:"has_garden"`
	HasTerrace        *bool    `json:"has_terrace"`
	HasBalcony        *bool    `json:"has_balcony"`
	HasSecurity       *bool    `json:"has_security"`
	HasElevator       *bool    `json:"has_elevator"`
	HasAirCondition   *bool    `json:"has_air_condition"`
	HasParking        *bool    `json:"has_parking"`
	MinParkingSpaces  *int     `json:"min_parking_spaces"`
	Furnished         *bool    `json:"furnished"`
	Tags              []string `json:"tags"`
	
	// Pagination
	Pagination *PaginationParams `json:"pagination"`
}

// NewPropertySearchFilters creates default search filters
func NewPropertySearchFilters() *PropertySearchFilters {
	return &PropertySearchFilters{
		Pagination: NewPaginationParams(),
	}
}