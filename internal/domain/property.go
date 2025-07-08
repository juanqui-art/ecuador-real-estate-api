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
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// NewProperty creates a new property with automatically generated SEO slug
func NewProperty(title, description, province, city, propertyType string, price float64) *Property {
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
		p.Type != ""
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