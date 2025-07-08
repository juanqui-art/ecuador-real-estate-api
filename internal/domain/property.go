package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Property represents a real estate property
type Property struct {
	// Primary identification
	ID   string `json:"id" db:"id"`
	Slug string `json:"slug" db:"slug"` // SEO-friendly URL

	// Basic information
	Title       string  `json:"title" db:"title"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`

	// Location
	Province string `json:"province" db:"province"` // Ecuador province
	City     string `json:"city" db:"city"`
	Sector   string `json:"sector" db:"sector"`   // Neighborhood/sector
	Address  string `json:"address" db:"address"` // Full address

	// Geolocation for maps
	Latitude          float64 `json:"latitude" db:"latitude"`                     // GPS latitude (-4.0 to 2.0 for Ecuador)
	Longitude         float64 `json:"longitude" db:"longitude"`                   // GPS longitude (-92.0 to -75.0 for Ecuador)
	LocationPrecision string  `json:"location_precision" db:"location_precision"` // exact, approximate, sector

	// Property characteristics
	Type      string  `json:"type" db:"type"`           // house, apartment, land, commercial
	Status    string  `json:"status" db:"status"`       // available, sold, rented, reserved
	Bedrooms  int     `json:"bedrooms" db:"bedrooms"`   // Number of bedrooms
	Bathrooms float32 `json:"bathrooms" db:"bathrooms"` // Can be 2.5 bathrooms
	AreaM2    float64 `json:"area_m2" db:"area_m2"`     // Area in square meters

	// Images and media
	MainImage string   `json:"main_image" db:"main_image"`           // URL of main image
	Images    []string `json:"images,omitempty" db:"images"`         // Array of image URLs (JSON in DB)
	VideoTour string   `json:"video_tour,omitempty" db:"video_tour"` // Video tour URL
	Tour360   string   `json:"tour_360,omitempty" db:"tour_360"`     // 360° virtual tour URL

	// Additional pricing
	RentPrice      *float64 `json:"rent_price,omitempty" db:"rent_price"`           // Monthly rent price
	CommonExpenses *float64 `json:"common_expenses,omitempty" db:"common_expenses"` // Common/HOA expenses
	PricePerM2     *float64 `json:"price_per_m2,omitempty" db:"price_per_m2"`       // Price per square meter

	// Detailed characteristics
	YearBuilt      *int   `json:"year_built,omitempty" db:"year_built"` // Construction year
	Floors         *int   `json:"floors,omitempty" db:"floors"`         // Number of floors
	PropertyStatus string `json:"property_status" db:"property_status"` // new, used, renovated
	Furnished      bool   `json:"furnished" db:"furnished"`             // Is furnished

	// Amenities (for frontend filters)
	Garage          bool `json:"garage" db:"garage"`                     // Has garage/parking
	Pool            bool `json:"pool" db:"pool"`                         // Has swimming pool
	Garden          bool `json:"garden" db:"garden"`                     // Has garden
	Terrace         bool `json:"terrace" db:"terrace"`                   // Has terrace
	Balcony         bool `json:"balcony" db:"balcony"`                   // Has balcony
	Security        bool `json:"security" db:"security"`                 // Gated community with security
	Elevator        bool `json:"elevator" db:"elevator"`                 // Has elevator
	AirConditioning bool `json:"air_conditioning" db:"air_conditioning"` // Has AC

	// Marketing and SEO
	Tags      []string `json:"tags,omitempty" db:"tags"`   // Search tags ["luxury", "ocean-view"]
	Featured  bool     `json:"featured" db:"featured"`     // Featured/premium property
	ViewCount int      `json:"view_count" db:"view_count"` // View counter for analytics

	// Relations
	RealEstateCompanyID *string `json:"real_estate_company_id,omitempty" db:"real_estate_company_id"` // FK to real estate companies

	// Audit fields
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Property types constants
const (
	PropertyTypeHouse      = "house"
	PropertyTypeApartment  = "apartment"
	PropertyTypeLand       = "land"
	PropertyTypeCommercial = "commercial"
)

// Property status constants
const (
	PropertyStatusAvailable = "available"
	PropertyStatusSold      = "sold"
	PropertyStatusRented    = "rented"
	PropertyStatusReserved  = "reserved"
)

// Property status constants (condition)
const (
	PropertyConditionNew       = "new"
	PropertyConditionUsed      = "used"
	PropertyConditionRenovated = "renovated"
)

// Location precision constants
const (
	LocationPrecisionExact       = "exact"
	LocationPrecisionApproximate = "approximate"
	LocationPrecisionSector      = "sector"
)

// NewProperty creates a new property with auto-generated ID and slug
func NewProperty(title, description, province, city, propertyType string, price float64) *Property {
	id := uuid.New().String()
	slug := GenerateSlug(title, id)

	return &Property{
		ID:             id,
		Slug:           slug,
		Title:          title,
		Description:    description,
		Price:          price,
		Province:       province,
		City:           city,
		Type:           propertyType,
		Status:         PropertyStatusAvailable,
		PropertyStatus: PropertyConditionUsed, // Default to used
		Bedrooms:       0,
		Bathrooms:      0,
		AreaM2:         0,
		Images:         []string{},
		Tags:           []string{},
		Featured:       false,
		ViewCount:      0,
		Furnished:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// GenerateSlug generates a SEO-friendly slug from title and ID
func GenerateSlug(title, id string) string {
	// Clean title: remove special characters, convert to lowercase
	slug := strings.ToLower(title)

	// Replace accented characters
	replacements := map[string]string{
		"á": "a", "é": "e", "í": "i", "ó": "o", "ú": "u", "ñ": "n",
		"Á": "a", "É": "e", "Í": "i", "Ó": "o", "Ú": "u", "Ñ": "n",
	}
	for old, new := range replacements {
		slug = strings.ReplaceAll(slug, old, new)
	}

	// Keep only letters, numbers, spaces, and hyphens
	reg := regexp.MustCompile(`[^a-z0-9\s\-]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace multiple spaces with single hyphen
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Truncate if too long
	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.Trim(slug, "-")
	}

	// Add short ID for uniqueness
	if len(id) >= 8 {
		slug = slug + "-" + id[:8]
	}

	return slug
}

// UpdateTimestamp updates the modification timestamp
func (p *Property) UpdateTimestamp() {
	p.UpdatedAt = time.Now()
}

// UpdateSlug regenerates the slug based on current title
func (p *Property) UpdateSlug() {
	p.Slug = GenerateSlug(p.Title, p.ID)
}

// Validate validates all property fields
func (p *Property) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(p.Province) == "" {
		return fmt.Errorf("province is required")
	}
	if strings.TrimSpace(p.City) == "" {
		return fmt.Errorf("city is required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if p.AreaM2 < 0 {
		return fmt.Errorf("area cannot be negative")
	}
	if p.Bedrooms < 0 {
		return fmt.Errorf("bedrooms cannot be negative")
	}
	if p.Bathrooms < 0 {
		return fmt.Errorf("bathrooms cannot be negative")
	}

	// Validate type
	if err := p.ValidateType(); err != nil {
		return err
	}

	// Validate status
	if err := p.ValidateStatus(); err != nil {
		return err
	}

	// Validate province
	if err := p.ValidateProvince(); err != nil {
		return err
	}

	// Validate coordinates if provided
	if err := p.ValidateCoordinates(); err != nil {
		return err
	}

	return nil
}

// ValidateType validates the property type
func (p *Property) ValidateType() error {
	validTypes := []string{PropertyTypeHouse, PropertyTypeApartment, PropertyTypeLand, PropertyTypeCommercial}
	for _, validType := range validTypes {
		if p.Type == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid property type: %s", p.Type)
}

// ValidateStatus validates the property status
func (p *Property) ValidateStatus() error {
	validStatuses := []string{PropertyStatusAvailable, PropertyStatusSold, PropertyStatusRented, PropertyStatusReserved}
	for _, validStatus := range validStatuses {
		if p.Status == validStatus {
			return nil
		}
	}
	return fmt.Errorf("invalid property status: %s", p.Status)
}

// ValidateProvince validates Ecuador provinces
func (p *Property) ValidateProvince() error {
	ecuadorProvinces := []string{
		"Azuay", "Bolívar", "Cañar", "Carchi", "Chimborazo", "Cotopaxi",
		"El Oro", "Esmeraldas", "Galápagos", "Guayas", "Imbabura", "Loja",
		"Los Ríos", "Manabí", "Morona Santiago", "Napo", "Orellana", "Pastaza",
		"Pichincha", "Santa Elena", "Santo Domingo", "Sucumbíos", "Tungurahua", "Zamora Chinchipe",
	}

	for _, province := range ecuadorProvinces {
		if p.Province == province {
			return nil
		}
	}
	return fmt.Errorf("invalid Ecuador province: %s", p.Province)
}

// ValidateCoordinates validates GPS coordinates for Ecuador
func (p *Property) ValidateCoordinates() error {
	if p.Latitude != 0 || p.Longitude != 0 {
		// Ecuador bounds: Latitude -4.0 to 2.0, Longitude -92.0 to -75.0
		if p.Latitude < -4.0 || p.Latitude > 2.0 {
			return fmt.Errorf("latitude must be between -4.0 and 2.0 for Ecuador")
		}
		if p.Longitude < -92.0 || p.Longitude > -75.0 {
			return fmt.Errorf("longitude must be between -92.0 and -75.0 for Ecuador")
		}
	}
	return nil
}

// CalculatePricePerM2 calculates and updates price per square meter
func (p *Property) CalculatePricePerM2() {
	if p.AreaM2 > 0 {
		pricePerM2 := p.Price / p.AreaM2
		p.PricePerM2 = &pricePerM2
	}
}

// AddImage adds an image URL to the gallery
func (p *Property) AddImage(imageURL string) error {
	if len(p.Images) >= 20 {
		return fmt.Errorf("maximum 20 images allowed")
	}

	if !isValidImageURL(imageURL) {
		return fmt.Errorf("invalid image URL format")
	}

	p.Images = append(p.Images, imageURL)
	p.UpdateTimestamp()
	return nil
}

// AddTag adds a search tag
func (p *Property) AddTag(tag string) error {
	if len(p.Tags) >= 10 {
		return fmt.Errorf("maximum 10 tags allowed")
	}

	// Check if tag already exists
	for _, existingTag := range p.Tags {
		if existingTag == tag {
			return nil // Already exists, no error
		}
	}

	p.Tags = append(p.Tags, strings.TrimSpace(tag))
	p.UpdateTimestamp()
	return nil
}

// IncrementViewCount increments the view counter
func (p *Property) IncrementViewCount() {
	p.ViewCount++
	// Don't update timestamp for view increments
}

// SetFeatured marks/unmarks property as featured
func (p *Property) SetFeatured(featured bool) {
	p.Featured = featured
	p.UpdateTimestamp()
}

// SetLocation sets the GPS coordinates and precision
func (p *Property) SetLocation(latitude, longitude float64, precision string) error {
	p.Latitude = latitude
	p.Longitude = longitude
	p.LocationPrecision = precision

	if err := p.ValidateCoordinates(); err != nil {
		return err
	}

	p.UpdateTimestamp()
	return nil
}

// GetAllImages returns all images (main + gallery)
func (p *Property) GetAllImages() []string {
	allImages := []string{}

	if p.MainImage != "" {
		allImages = append(allImages, p.MainImage)
	}

	allImages = append(allImages, p.Images...)
	return allImages
}

// IsAvailable checks if property is available for sale/rent
func (p *Property) IsAvailable() bool {
	return p.Status == PropertyStatusAvailable
}

// HasLocation checks if property has GPS coordinates
func (p *Property) HasLocation() bool {
	return p.Latitude != 0 && p.Longitude != 0
}

// GetSummary returns a summary for listings
func (p *Property) GetSummary() map[string]interface{} {
	return map[string]interface{}{
		"id":         p.ID,
		"slug":       p.Slug,
		"title":      p.Title,
		"price":      p.Price,
		"province":   p.Province,
		"city":       p.City,
		"type":       p.Type,
		"status":     p.Status,
		"bedrooms":   p.Bedrooms,
		"bathrooms":  p.Bathrooms,
		"area_m2":    p.AreaM2,
		"main_image": p.MainImage,
		"featured":   p.Featured,
		"view_count": p.ViewCount,
	}
}

// Helper function to validate image URLs
func isValidImageURL(url string) bool {
	if url == "" {
		return false
	}
	// Basic URL validation for images
	pattern := `^https?://.*\.(jpg|jpeg|png|webp)(\?.*)?$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}
