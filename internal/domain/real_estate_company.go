package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RealEstateCompany represents a real estate company/agency
type RealEstateCompany struct {
	// Primary identification
	ID string `json:"id" db:"id"`

	// Basic company information
	Name        string `json:"name" db:"name"`
	RUC         string `json:"ruc" db:"ruc"`                 // Ecuador tax ID (13 digits)
	Address     string `json:"address" db:"address"`         // Physical address
	Description string `json:"description" db:"description"` // Company description

	// Contact information
	Phone   string `json:"phone" db:"phone"`
	Email   string `json:"email" db:"email"`
	Website string `json:"website" db:"website"`   // Company website
	LogoURL string `json:"logo_url" db:"logo_url"` // Company logo

	// Status
	Active bool `json:"active" db:"active"` // Is company active

	// Audit fields
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewRealEstateCompany creates a new real estate company
func NewRealEstateCompany(name, ruc, address, phone, email string) *RealEstateCompany {
	return &RealEstateCompany{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		RUC:       strings.TrimSpace(ruc),
		Address:   strings.TrimSpace(address),
		Phone:     strings.TrimSpace(phone),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Active:    true, // Active by default
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdateTimestamp updates the modification timestamp
func (r *RealEstateCompany) UpdateTimestamp() {
	r.UpdatedAt = time.Now()
}

// Validate validates all company fields
func (r *RealEstateCompany) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("company name is required")
	}
	if strings.TrimSpace(r.RUC) == "" {
		return fmt.Errorf("RUC is required")
	}
	if strings.TrimSpace(r.Address) == "" {
		return fmt.Errorf("address is required")
	}
	if strings.TrimSpace(r.Phone) == "" {
		return fmt.Errorf("phone is required")
	}
	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("email is required")
	}

	// Validate formats
	if err := r.ValidateRUC(); err != nil {
		return err
	}
	if err := r.ValidateEmail(); err != nil {
		return err
	}
	if err := r.ValidatePhone(); err != nil {
		return err
	}

	return nil
}

// ValidateRUC validates Ecuador RUC format (13 digits ending in 001 for companies)
func (r *RealEstateCompany) ValidateRUC() error {
	ruc := strings.TrimSpace(r.RUC)

	// Must be exactly 13 digits
	if len(ruc) != 13 {
		return fmt.Errorf("RUC must be 13 digits")
	}

	// Must contain only numbers
	if !regexp.MustCompile(`^\d{13}$`).MatchString(ruc) {
		return fmt.Errorf("RUC must contain only numbers")
	}

	// Must end in 001 for companies
	if !strings.HasSuffix(ruc, "001") {
		return fmt.Errorf("company RUC must end in 001")
	}

	// Basic validation: first two digits should be valid province (01-24)
	provinceCode := ruc[:2]
	if provinceCode < "01" || provinceCode > "24" {
		return fmt.Errorf("invalid province code in RUC")
	}

	return nil
}

// ValidateEmail validates email format
func (r *RealEstateCompany) ValidateEmail() error {
	email := strings.TrimSpace(r.Email)

	// Basic email regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)

	if !matched {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}

// ValidatePhone validates Ecuador phone format
func (r *RealEstateCompany) ValidatePhone() error {
	phone := strings.TrimSpace(r.Phone)

	// Ecuador phone patterns:
	// Landline: 02-XXXXXXX, 03-XXXXXXX, etc. (9 digits with area code)
	// Mobile: 09-XXXXXXXX (10 digits)
	// With country code: +593-X-XXXXXXX

	// Remove spaces, dashes, parentheses
	cleanPhone := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Validate different formats
	if len(cleanPhone) == 9 {
		// Landline without country code (0X-XXXXXXX)
		if !regexp.MustCompile(`^0[2-7]\d{7}$`).MatchString(cleanPhone) {
			return fmt.Errorf("invalid landline phone format: %s", phone)
		}
	} else if len(cleanPhone) == 10 {
		// Mobile without country code (09-XXXXXXXX)
		if !regexp.MustCompile(`^09\d{8}$`).MatchString(cleanPhone) {
			return fmt.Errorf("invalid mobile phone format: %s", phone)
		}
	} else if len(cleanPhone) == 12 && strings.HasPrefix(cleanPhone, "593") {
		// With country code (+593-X-XXXXXXX)
		if !regexp.MustCompile(`^593[0-9]\d{7,8}$`).MatchString(cleanPhone) {
			return fmt.Errorf("invalid phone format with country code: %s", phone)
		}
	} else {
		return fmt.Errorf("invalid phone format: %s", phone)
	}

	return nil
}

// ValidateWebsite validates website URL format
func (r *RealEstateCompany) ValidateWebsite() error {
	if r.Website == "" {
		return nil // Website is optional
	}

	website := strings.TrimSpace(r.Website)
	urlRegex := `^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`
	matched, _ := regexp.MatchString(urlRegex, website)

	if !matched {
		return fmt.Errorf("invalid website URL format: %s", website)
	}

	return nil
}

// ValidateLogoURL validates logo URL format
func (r *RealEstateCompany) ValidateLogoURL() error {
	if r.LogoURL == "" {
		return nil // Logo is optional
	}

	logoURL := strings.TrimSpace(r.LogoURL)
	// Logo should be an image URL
	imageURLRegex := `^https?://.*\.(jpg|jpeg|png|webp|svg)(\?.*)?$`
	matched, _ := regexp.MatchString(imageURLRegex, logoURL)

	if !matched {
		return fmt.Errorf("invalid logo URL format (must be jpg, png, webp, or svg): %s", logoURL)
	}

	return nil
}

// SetRUC sets and validates the RUC
func (r *RealEstateCompany) SetRUC(ruc string) error {
	r.RUC = strings.TrimSpace(ruc)
	if err := r.ValidateRUC(); err != nil {
		return err
	}
	r.UpdateTimestamp()
	return nil
}

// SetContactInfo updates contact information
func (r *RealEstateCompany) SetContactInfo(phone, email string) error {
	r.Phone = strings.TrimSpace(phone)
	r.Email = strings.ToLower(strings.TrimSpace(email))

	if err := r.ValidatePhone(); err != nil {
		return err
	}
	if err := r.ValidateEmail(); err != nil {
		return err
	}

	r.UpdateTimestamp()
	return nil
}

// SetWebsite sets and validates the website URL
func (r *RealEstateCompany) SetWebsite(website string) error {
	r.Website = strings.TrimSpace(website)
	if err := r.ValidateWebsite(); err != nil {
		return err
	}
	r.UpdateTimestamp()
	return nil
}

// SetLogoURL sets and validates the logo URL
func (r *RealEstateCompany) SetLogoURL(logoURL string) error {
	r.LogoURL = strings.TrimSpace(logoURL)
	if err := r.ValidateLogoURL(); err != nil {
		return err
	}
	r.UpdateTimestamp()
	return nil
}

// Activate activates the company
func (r *RealEstateCompany) Activate() {
	r.Active = true
	r.UpdateTimestamp()
}

// Deactivate deactivates the company
func (r *RealEstateCompany) Deactivate() {
	r.Active = false
	r.UpdateTimestamp()
}

// GetDisplayName returns the display name (with status indicator if inactive)
func (r *RealEstateCompany) GetDisplayName() string {
	if r.Active {
		return r.Name
	}
	return r.Name + " (Inactive)"
}

// GetSummary returns a summary for listings
func (r *RealEstateCompany) GetSummary() map[string]interface{} {
	return map[string]interface{}{
		"id":      r.ID,
		"name":    r.GetDisplayName(),
		"ruc":     r.RUC,
		"phone":   r.Phone,
		"email":   r.Email,
		"website": r.Website,
		"active":  r.Active,
	}
}

// FormatPhone formats the phone number for display
func (r *RealEstateCompany) FormatPhone() string {
	phone := r.Phone
	cleanPhone := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Add formatting based on length
	if len(cleanPhone) == 9 && cleanPhone[0] == '0' {
		// Landline: 0X-XXX-XXXX
		return fmt.Sprintf("%s-%s-%s", cleanPhone[:2], cleanPhone[2:5], cleanPhone[5:])
	} else if len(cleanPhone) == 10 && strings.HasPrefix(cleanPhone, "09") {
		// Mobile: 09-XXXX-XXXX
		return fmt.Sprintf("%s-%s-%s", cleanPhone[:2], cleanPhone[2:6], cleanPhone[6:])
	}

	return phone // Return original if no pattern matches
}

// IsValidRUC validates a RUC using the Ecuador algorithm (mod 11)
func IsValidRUC(ruc string) bool {
	cleanRUC := strings.TrimSpace(ruc)

	if len(cleanRUC) != 13 || !regexp.MustCompile(`^\d{13}$`).MatchString(cleanRUC) {
		return false
	}

	if !strings.HasSuffix(cleanRUC, "001") {
		return false
	}

	// Convert to digits
	digits := make([]int, 10)
	for i := 0; i < 10; i++ {
		digits[i] = int(cleanRUC[i] - '0')
	}

	// Third digit must be less than 6 for companies
	if digits[2] >= 6 {
		return false
	}

	// Apply mod 11 algorithm
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}

	checkDigit := 11 - (sum % 11)
	if checkDigit == 11 {
		checkDigit = 0
	} else if checkDigit == 10 {
		checkDigit = 1
	}

	return checkDigit == digits[9]
}
