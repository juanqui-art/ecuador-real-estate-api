package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents a system user (buyer, seller, agent, admin)
type User struct {
	// Primary identification
	ID string `json:"id" db:"id"`

	// Personal information
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	Email       string     `json:"email" db:"email"`
	Phone       string     `json:"phone" db:"phone"`
	NationalID  string     `json:"national_id" db:"national_id"` // Ecuador cedula (10 digits)
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`

	// User type and status
	UserType string `json:"user_type" db:"user_type"` // buyer, seller, agent, admin
	Active   bool   `json:"active" db:"active"`       // Is user active

	// Search preferences (for buyers)
	MinBudget              *float64 `json:"min_budget,omitempty" db:"min_budget"`                             // Minimum budget
	MaxBudget              *float64 `json:"max_budget,omitempty" db:"max_budget"`                             // Maximum budget
	PreferredProvinces     []string `json:"preferred_provinces,omitempty" db:"preferred_provinces"`           // Provinces of interest (JSON in DB)
	PreferredPropertyTypes []string `json:"preferred_property_types,omitempty" db:"preferred_property_types"` // Property types of interest (JSON in DB)

	// Profile information
	AvatarURL string `json:"avatar_url,omitempty" db:"avatar_url"` // Avatar image URL
	Bio       string `json:"bio,omitempty" db:"bio"`               // User biography/description

	// Relationship with RealEstateCompany (for agents)
	RealEstateCompanyID *string `json:"real_estate_company_id,omitempty" db:"real_estate_company_id"` // FK to real estate companies

	// Notification preferences
	ReceiveNotifications bool `json:"receive_notifications" db:"receive_notifications"` // Receive app notifications
	ReceiveNewsletter    bool `json:"receive_newsletter" db:"receive_newsletter"`       // Receive newsletter emails

	// Audit fields
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// User type constants
const (
	UserTypeBuyer  = "buyer"  // User looking to buy/rent
	UserTypeSeller = "seller" // User selling/renting properties
	UserTypeAgent  = "agent"  // Real estate agent
	UserTypeAdmin  = "admin"  // System administrator
)

// NewUser creates a new basic user
func NewUser(firstName, lastName, email, phone, nationalID, userType string) *User {
	return &User{
		ID:                     uuid.New().String(),
		FirstName:              strings.TrimSpace(firstName),
		LastName:               strings.TrimSpace(lastName),
		Email:                  strings.ToLower(strings.TrimSpace(email)),
		Phone:                  strings.TrimSpace(phone),
		NationalID:             strings.TrimSpace(nationalID),
		UserType:               userType,
		Active:                 true, // Active by default
		PreferredProvinces:     []string{},
		PreferredPropertyTypes: []string{},
		ReceiveNotifications:   true,  // Default to receiving notifications
		ReceiveNewsletter:      false, // Default to no newsletter
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
}

// NewBuyer creates a new buyer user with preferences
func NewBuyer(firstName, lastName, email, phone, nationalID string, minBudget, maxBudget *float64, provinces, propertyTypes []string) *User {
	user := NewUser(firstName, lastName, email, phone, nationalID, UserTypeBuyer)
	user.MinBudget = minBudget
	user.MaxBudget = maxBudget
	user.PreferredProvinces = provinces
	user.PreferredPropertyTypes = propertyTypes
	return user
}

// NewAgent creates a new agent user associated with a real estate company
func NewAgent(firstName, lastName, email, phone, nationalID, realEstateCompanyID string) *User {
	user := NewUser(firstName, lastName, email, phone, nationalID, UserTypeAgent)
	user.RealEstateCompanyID = &realEstateCompanyID
	return user
}

// UpdateTimestamp updates the modification timestamp
func (u *User) UpdateTimestamp() {
	u.UpdatedAt = time.Now()
}

// Validate validates all user fields
func (u *User) Validate() error {
	if strings.TrimSpace(u.FirstName) == "" {
		return fmt.Errorf("first name is required")
	}
	if strings.TrimSpace(u.LastName) == "" {
		return fmt.Errorf("last name is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(u.Phone) == "" {
		return fmt.Errorf("phone is required")
	}
	if strings.TrimSpace(u.NationalID) == "" {
		return fmt.Errorf("national ID is required")
	}

	// Format validations
	if err := u.ValidateEmail(); err != nil {
		return err
	}
	if err := u.ValidatePhone(); err != nil {
		return err
	}
	if err := u.ValidateNationalID(); err != nil {
		return err
	}
	if err := u.ValidateUserType(); err != nil {
		return err
	}
	if err := u.ValidateBudget(); err != nil {
		return err
	}

	return nil
}

// ValidateEmail validates email format
func (u *User) ValidateEmail() error {
	email := strings.TrimSpace(u.Email)
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)

	if !matched {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}

// ValidatePhone validates Ecuador phone format
func (u *User) ValidatePhone() error {
	phone := strings.TrimSpace(u.Phone)

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

// ValidateNationalID validates Ecuador national ID (cedula)
func (u *User) ValidateNationalID() error {
	nationalID := strings.TrimSpace(u.NationalID)

	if len(nationalID) != 10 {
		return fmt.Errorf("national ID must be 10 digits")
	}

	// Must contain only numbers
	if !regexp.MustCompile(`^\d{10}$`).MatchString(nationalID) {
		return fmt.Errorf("national ID must contain only numbers")
	}

	// Basic validation: first two digits should be valid province (01-24)
	provinceCode := nationalID[:2]
	if provinceCode < "01" || provinceCode > "24" {
		return fmt.Errorf("invalid province code in national ID")
	}

	// Third digit must be less than 6 for natural persons
	if nationalID[2] >= '6' {
		return fmt.Errorf("invalid national ID format")
	}

	// Validate using Ecuador algorithm
	if !IsValidEcuadorNationalID(nationalID) {
		return fmt.Errorf("invalid Ecuador national ID: %s", nationalID)
	}

	return nil
}

// ValidateUserType validates the user type
func (u *User) ValidateUserType() error {
	validTypes := []string{UserTypeBuyer, UserTypeSeller, UserTypeAgent, UserTypeAdmin}
	for _, validType := range validTypes {
		if u.UserType == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid user type: %s", u.UserType)
}

// ValidateBudget validates budget for buyers
func (u *User) ValidateBudget() error {
	if u.MinBudget != nil && *u.MinBudget < 0 {
		return fmt.Errorf("minimum budget cannot be negative")
	}
	if u.MaxBudget != nil && *u.MaxBudget < 0 {
		return fmt.Errorf("maximum budget cannot be negative")
	}
	if u.MinBudget != nil && u.MaxBudget != nil && *u.MinBudget > *u.MaxBudget {
		return fmt.Errorf("minimum budget cannot be greater than maximum")
	}
	return nil
}

// Activate activates the user
func (u *User) Activate() {
	u.Active = true
	u.UpdateTimestamp()
}

// Deactivate deactivates the user
func (u *User) Deactivate() {
	u.Active = false
	u.UpdateTimestamp()
}

// SetProfile updates profile information
func (u *User) SetProfile(avatarURL, bio string) {
	u.AvatarURL = avatarURL
	u.Bio = bio
	u.UpdateTimestamp()
}

// SetRealEstateCompany assigns a real estate company to an agent
func (u *User) SetRealEstateCompany(companyID string) {
	u.RealEstateCompanyID = &companyID
	u.UpdateTimestamp()
}

// SetBudget sets the budget range for buyers
func (u *User) SetBudget(minBudget, maxBudget *float64) error {
	u.MinBudget = minBudget
	u.MaxBudget = maxBudget

	if err := u.ValidateBudget(); err != nil {
		return err
	}

	u.UpdateTimestamp()
	return nil
}

// SetPreferences sets search preferences for buyers
func (u *User) SetPreferences(provinces, propertyTypes []string) {
	u.PreferredProvinces = provinces
	u.PreferredPropertyTypes = propertyTypes
	u.UpdateTimestamp()
}

// AddPreferredPropertyType adds a property type to preferences
func (u *User) AddPreferredPropertyType(propertyType string) error {
	// Validate property type
	validTypes := []string{PropertyTypeHouse, PropertyTypeApartment, PropertyTypeLand, PropertyTypeCommercial}
	isValid := false
	for _, validType := range validTypes {
		if validType == propertyType {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid property type: %s", propertyType)
	}

	// Check if already exists
	for _, existing := range u.PreferredPropertyTypes {
		if existing == propertyType {
			return nil // Already exists, no error
		}
	}

	u.PreferredPropertyTypes = append(u.PreferredPropertyTypes, propertyType)
	u.UpdateTimestamp()
	return nil
}

// AddPreferredProvince adds a province to preferences
func (u *User) AddPreferredProvince(province string) error {
	if !IsValidEcuadorProvince(province) {
		return fmt.Errorf("invalid Ecuador province: %s", province)
	}

	// Check if already exists
	for _, existing := range u.PreferredProvinces {
		if existing == province {
			return nil // Already exists, no error
		}
	}

	u.PreferredProvinces = append(u.PreferredProvinces, province)
	u.UpdateTimestamp()
	return nil
}

// SetNotificationPreferences sets notification preferences
func (u *User) SetNotificationPreferences(notifications, newsletter bool) {
	u.ReceiveNotifications = notifications
	u.ReceiveNewsletter = newsletter
	u.UpdateTimestamp()
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// GetDisplayName returns display name with status if inactive
func (u *User) GetDisplayName() string {
	if u.Active {
		return u.GetFullName()
	}
	return u.GetFullName() + " (Inactive)"
}

// GetSummary returns a summary for listings
func (u *User) GetSummary() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.GetDisplayName(),
		"email":      u.Email,
		"phone":      u.Phone,
		"user_type":  u.UserType,
		"active":     u.Active,
		"created_at": u.CreatedAt,
	}
}

// IsBuyer checks if user is a buyer
func (u *User) IsBuyer() bool {
	return u.UserType == UserTypeBuyer
}

// IsSeller checks if user is a seller
func (u *User) IsSeller() bool {
	return u.UserType == UserTypeSeller
}

// IsAgent checks if user is an agent
func (u *User) IsAgent() bool {
	return u.UserType == UserTypeAgent
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.UserType == UserTypeAdmin
}

// HasBudgetConfigured checks if user has budget configured
func (u *User) HasBudgetConfigured() bool {
	return u.MinBudget != nil || u.MaxBudget != nil
}

// CanAffordProperty checks if a property is within budget range
func (u *User) CanAffordProperty(propertyPrice float64) bool {
	if !u.HasBudgetConfigured() {
		return true // If no budget set, all properties are valid
	}

	if u.MinBudget != nil && propertyPrice < *u.MinBudget {
		return false
	}

	if u.MaxBudget != nil && propertyPrice > *u.MaxBudget {
		return false
	}

	return true
}

// CanReceiveNotifications checks if user can receive notifications
func (u *User) CanReceiveNotifications() bool {
	return u.Active && u.ReceiveNotifications
}

// FormatPhone formats the phone number for display
func (u *User) FormatPhone() string {
	phone := u.Phone
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

// IsValidEcuadorNationalID validates Ecuador national ID using the algorithm
func IsValidEcuadorNationalID(nationalID string) bool {
	cleanID := strings.TrimSpace(nationalID)

	if len(cleanID) != 10 || !regexp.MustCompile(`^\d{10}$`).MatchString(cleanID) {
		return false
	}

	// Convert to digits
	digits := make([]int, 10)
	for i, r := range cleanID {
		digits[i] = int(r - '0')
	}

	// Third digit must be less than 6 for natural persons
	if digits[2] >= 6 {
		return false
	}

	// Apply Ecuador national ID algorithm
	sum := 0
	for i := 0; i < 9; i++ {
		if i%2 == 0 {
			// Even positions (0,2,4,6,8)
			product := digits[i] * 2
			if product > 9 {
				product -= 9
			}
			sum += product
		} else {
			// Odd positions (1,3,5,7)
			sum += digits[i]
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit == digits[9]
}

// IsValidEcuadorProvince validates if a province is valid in Ecuador
func IsValidEcuadorProvince(province string) bool {
	ecuadorProvinces := []string{
		"Azuay", "Bolívar", "Cañar", "Carchi", "Chimborazo", "Cotopaxi",
		"El Oro", "Esmeraldas", "Galápagos", "Guayas", "Imbabura", "Loja",
		"Los Ríos", "Manabí", "Morona Santiago", "Napo", "Orellana", "Pastaza",
		"Pichincha", "Santa Elena", "Santo Domingo", "Sucumbíos", "Tungurahua", "Zamora Chinchipe",
	}

	for _, validProvince := range ecuadorProvinces {
		if validProvince == province {
			return true
		}
	}
	return false
}
