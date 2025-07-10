package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AgencyStatus represents the agency status
type AgencyStatus string

const (
	AgencyStatusActive    AgencyStatus = "active"
	AgencyStatusInactive  AgencyStatus = "inactive"
	AgencyStatusSuspended AgencyStatus = "suspended"
	AgencyStatusPending   AgencyStatus = "pending"
)

// Agency represents a real estate agency
type Agency struct {
	ID          string       `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Email       string       `json:"email" db:"email"`
	Phone       string       `json:"phone" db:"phone"`
	Address     string       `json:"address" db:"address"`
	City        string       `json:"city" db:"city"`
	Province    string       `json:"province" db:"province"`
	License     string       `json:"license" db:"license"`     // RUC or business license
	Website     *string      `json:"website" db:"website"`
	Description *string      `json:"description" db:"description"`
	Logo        *string      `json:"logo" db:"logo"`
	Status      AgencyStatus `json:"status" db:"status"`
	OwnerID     string       `json:"owner_id" db:"owner_id"`   // User who owns/manages the agency
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at" db:"deleted_at"`
}

// NewAgency creates a new agency with validation
func NewAgency(name, email, phone, address, city, province, license, ownerID string) (*Agency, error) {
	if err := validateAgencyName(name); err != nil {
		return nil, err
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validatePhone(phone); err != nil {
		return nil, err
	}

	if err := validateAddress(address); err != nil {
		return nil, err
	}

	if err := validateEcuadorProvince(province); err != nil {
		return nil, err
	}

	if err := validateLicense(license); err != nil {
		return nil, err
	}

	if ownerID == "" {
		return nil, fmt.Errorf("owner ID cannot be empty")
	}

	agency := &Agency{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Phone:     strings.TrimSpace(phone),
		Address:   strings.TrimSpace(address),
		City:      strings.TrimSpace(city),
		Province:  strings.TrimSpace(province),
		License:   strings.TrimSpace(license),
		Status:    AgencyStatusPending,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return agency, nil
}

// IsValid validates the agency data
func (a *Agency) IsValid() error {
	if a.ID == "" {
		return fmt.Errorf("agency ID cannot be empty")
	}

	if err := validateAgencyName(a.Name); err != nil {
		return err
	}

	if err := validateEmail(a.Email); err != nil {
		return err
	}

	if err := validatePhone(a.Phone); err != nil {
		return err
	}

	if err := validateAddress(a.Address); err != nil {
		return err
	}

	if err := validateEcuadorProvince(a.Province); err != nil {
		return err
	}

	if err := validateLicense(a.License); err != nil {
		return err
	}

	if a.OwnerID == "" {
		return fmt.Errorf("owner ID cannot be empty")
	}

	if err := validateAgencyStatus(a.Status); err != nil {
		return err
	}

	// Optional fields validation
	if a.Website != nil && *a.Website != "" {
		if err := validateWebsite(*a.Website); err != nil {
			return err
		}
	}

	return nil
}

// Activate sets the agency status to active
func (a *Agency) Activate() error {
	if a.Status == AgencyStatusSuspended {
		return fmt.Errorf("suspended agencies cannot be activated directly")
	}

	a.Status = AgencyStatusActive
	a.UpdatedAt = time.Now()
	return nil
}

// Deactivate sets the agency status to inactive
func (a *Agency) Deactivate() error {
	a.Status = AgencyStatusInactive
	a.UpdatedAt = time.Now()
	return nil
}

// Suspend sets the agency status to suspended
func (a *Agency) Suspend() error {
	a.Status = AgencyStatusSuspended
	a.UpdatedAt = time.Now()
	return nil
}

// UpdateInfo updates the agency information
func (a *Agency) UpdateInfo(name, phone, address, city, province string) error {
	if err := validateAgencyName(name); err != nil {
		return err
	}

	if err := validatePhone(phone); err != nil {
		return err
	}

	if err := validateAddress(address); err != nil {
		return err
	}

	if err := validateEcuadorProvince(province); err != nil {
		return err
	}

	a.Name = strings.TrimSpace(name)
	a.Phone = strings.TrimSpace(phone)
	a.Address = strings.TrimSpace(address)
	a.City = strings.TrimSpace(city)
	a.Province = strings.TrimSpace(province)
	a.UpdatedAt = time.Now()

	return nil
}

// SetWebsite updates the agency website
func (a *Agency) SetWebsite(website string) error {
	if website == "" {
		a.Website = nil
		a.UpdatedAt = time.Now()
		return nil
	}

	if err := validateWebsite(website); err != nil {
		return err
	}

	a.Website = &website
	a.UpdatedAt = time.Now()
	return nil
}

// SetDescription updates the agency description
func (a *Agency) SetDescription(description string) error {
	if description == "" {
		a.Description = nil
		a.UpdatedAt = time.Now()
		return nil
	}

	if len(description) > 1000 {
		return fmt.Errorf("description cannot exceed 1000 characters")
	}

	a.Description = &description
	a.UpdatedAt = time.Now()
	return nil
}

// SetLogo updates the agency logo
func (a *Agency) SetLogo(logoURL string) error {
	if logoURL == "" {
		a.Logo = nil
		a.UpdatedAt = time.Now()
		return nil
	}

	a.Logo = &logoURL
	a.UpdatedAt = time.Now()
	return nil
}

// IsActive checks if the agency is active
func (a *Agency) IsActive() bool {
	return a.Status == AgencyStatusActive
}

// CanManageProperty checks if agency can manage a property
func (a *Agency) CanManageProperty(propertyAgencyID *string) bool {
	if !a.IsActive() {
		return false
	}

	return propertyAgencyID != nil && *propertyAgencyID == a.ID
}

// GetDisplayName returns the agency display name
func (a *Agency) GetDisplayName() string {
	return a.Name
}

// GetContactInfo returns formatted contact information
func (a *Agency) GetContactInfo() string {
	return fmt.Sprintf("%s - %s - %s", a.Name, a.Phone, a.Email)
}

// GetFullAddress returns the complete address
func (a *Agency) GetFullAddress() string {
	return fmt.Sprintf("%s, %s, %s", a.Address, a.City, a.Province)
}

// Agency validation functions
func validateAgencyName(name string) error {
	if name == "" {
		return fmt.Errorf("agency name cannot be empty")
	}

	if len(name) < 2 {
		return fmt.Errorf("agency name must be at least 2 characters long")
	}

	if len(name) > 255 {
		return fmt.Errorf("agency name cannot exceed 255 characters")
	}

	return nil
}

func validatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone cannot be empty")
	}

	// Ecuador phone number validation (basic)
	// Format: +593xxxxxxxxx or 0xxxxxxxxx
	phoneRegex := regexp.MustCompile(`^(\+593|0)[0-9]{9}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid Ecuador phone number format")
	}

	return nil
}

func validateAddress(address string) error {
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if len(address) < 5 {
		return fmt.Errorf("address must be at least 5 characters long")
	}

	if len(address) > 500 {
		return fmt.Errorf("address cannot exceed 500 characters")
	}

	return nil
}

func validateLicense(license string) error {
	if license == "" {
		return fmt.Errorf("license cannot be empty")
	}

	// Ecuador RUC validation (basic)
	// Format: 13 digits
	rucRegex := regexp.MustCompile(`^[0-9]{13}$`)
	if !rucRegex.MatchString(license) {
		return fmt.Errorf("invalid Ecuador RUC format (must be 13 digits)")
	}

	return nil
}

func validateWebsite(website string) error {
	if website == "" {
		return nil
	}

	// Basic URL validation
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(website) {
		return fmt.Errorf("invalid website URL format")
	}

	return nil
}

func validateEcuadorProvince(province string) error {
	if province == "" {
		return fmt.Errorf("province cannot be empty")
	}

	validProvinces := []string{
		"Azuay", "Bolívar", "Cañar", "Carchi", "Chimborazo", "Cotopaxi",
		"El Oro", "Esmeraldas", "Galápagos", "Guayas", "Imbabura", "Loja",
		"Los Ríos", "Manabí", "Morona Santiago", "Napo", "Orellana", "Pastaza",
		"Pichincha", "Santa Elena", "Santo Domingo", "Sucumbíos", "Tungurahua", "Zamora Chinchipe",
	}

	for _, validProvince := range validProvinces {
		if province == validProvince {
			return nil
		}
	}

	return fmt.Errorf("invalid Ecuador province: %s", province)
}

func validateAgencyStatus(status AgencyStatus) error {
	validStatuses := []AgencyStatus{
		AgencyStatusActive, AgencyStatusInactive, 
		AgencyStatusSuspended, AgencyStatusPending,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid agency status: %s", status)
}