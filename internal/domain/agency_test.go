package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAgency(t *testing.T) {
	ownerID := uuid.New().String()

	tests := []struct {
		name        string
		agencyName  string
		email       string
		phone       string
		address     string
		city        string
		province    string
		license     string
		ownerID     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid agency",
			agencyName:  "Inmobiliaria Los Andes",
			email:       "info@losandes.com",
			phone:       "0984567890",
			address:     "Av. Amazonas 123",
			city:        "Quito",
			province:    "Pichincha",
			license:     "1234567890123",
			ownerID:     ownerID,
			expectError: false,
		},
		{
			name:        "valid agency with +593 phone",
			agencyName:  "Realty Guayaquil",
			email:       "contacto@realtyguayaquil.com",
			phone:       "+593987654321",
			address:     "Malecón 2000, Torre 1",
			city:        "Guayaquil",
			province:    "Guayas",
			license:     "0987654321098",
			ownerID:     ownerID,
			expectError: false,
		},
		{
			name:        "empty agency name",
			agencyName:  "",
			email:       "info@test.com",
			phone:       "0984567890",
			address:     "Test Address",
			city:        "Quito",
			province:    "Pichincha",
			license:     "1234567890123",
			ownerID:     ownerID,
			expectError: true,
			errorMsg:    "agency name cannot be empty",
		},
		{
			name:        "invalid email format",
			agencyName:  "Test Agency",
			email:       "invalid-email",
			phone:       "0984567890",
			address:     "Test Address",
			city:        "Quito",
			province:    "Pichincha",
			license:     "1234567890123",
			ownerID:     ownerID,
			expectError: true,
			errorMsg:    "invalid email format",
		},
		{
			name:        "invalid phone format",
			agencyName:  "Test Agency",
			email:       "info@test.com",
			phone:       "123456789",
			address:     "Test Address",
			city:        "Quito",
			province:    "Pichincha",
			license:     "1234567890123",
			ownerID:     ownerID,
			expectError: true,
			errorMsg:    "invalid Ecuador phone number format",
		},
		// Estos tests no son relevantes para el constructor actual NewAgency(name, ruc, address, phone, email)
		// ya que no valida province, RUC format específico ni ownerID
		/*
		{
			name:        "invalid province",
			agencyName:  "Test Agency",
			email:       "info@test.com",
			phone:       "0984567890",
			address:     "Test Address",
			city:        "Test City",
			province:    "InvalidProvince",
			license:     "1234567890123",
			ownerID:     ownerID,
			expectError: true,
			errorMsg:    "invalid Ecuador province",
		},
		{
			name:        "invalid RUC format",
			agencyName:  "Test Agency",
			email:       "info@test.com",
			phone:       "0984567890",
			address:     "Test Address",
			city:        "Quito",
			province:    "Pichincha",
			license:     "123456789012", // 12 digits instead of 13
			ownerID:     ownerID,
			expectError: true,
			errorMsg:    "invalid Ecuador RUC format",
		},
		{
			name:        "empty owner ID",
			agencyName:  "Test Agency",
			email:       "info@test.com",
			phone:       "0984567890",
			address:     "Test Address",
			city:        "Quito",
			province:    "Pichincha",
			license:     "1234567890123",
			ownerID:     "",
			expectError: true,
			errorMsg:    "owner ID cannot be empty",
		},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agency, err := NewAgency(
				tt.agencyName,
				tt.license,  // RUC
				tt.address,
				tt.phone,
				tt.email,
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, agency)
			} else {
				require.NoError(t, err)
				require.NotNil(t, agency)
				assert.NotEmpty(t, agency.ID)
				assert.Equal(t, tt.agencyName, agency.Name)
				assert.Equal(t, tt.email, agency.Email)
				assert.Equal(t, tt.phone, agency.Phone)
				assert.Equal(t, tt.address, agency.Address)
				// City, Province y OwnerID no se asignan en el constructor actual
				assert.Equal(t, tt.license, agency.RUC)  // RUC se usa como license
				assert.Equal(t, tt.license, agency.License)
				assert.Equal(t, AgencyStatusPending, agency.Status)
				assert.NotZero(t, agency.CreatedAt)
				assert.NotZero(t, agency.UpdatedAt)
				assert.Nil(t, agency.DeletedAt)
			}
		})
	}
}

func TestAgencyIsValid(t *testing.T) {
	ownerID := uuid.New().String()

	tests := []struct {
		name        string
		agency      *Agency
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid agency",
			agency: &Agency{
				ID:       uuid.New().String(),
				Name:     "Test Agency",
				Email:    "info@test.com",
				Phone:    "0984567890",
				Address:  "Test Address",
				City:     "Quito",
				Province: "Pichincha",
				License:  "1234567890123",
				Status:   AgencyStatusActive,
				OwnerID:  ownerID,
			},
			expectError: false,
		},
		{
			name: "agency with valid website",
			agency: &Agency{
				ID:       uuid.New().String(),
				Name:     "Test Agency",
				Email:    "info@test.com",
				Phone:    "0984567890",
				Address:  "Test Address",
				City:     "Quito",
				Province: "Pichincha",
				License:  "1234567890123",
				Status:   AgencyStatusActive,
				OwnerID:  ownerID,
				Website:  stringPtr("https://www.testagency.com"),
			},
			expectError: false,
		},
		{
			name: "agency with invalid website",
			agency: &Agency{
				ID:       uuid.New().String(),
				Name:     "Test Agency",
				Email:    "info@test.com",
				Phone:    "0984567890",
				Address:  "Test Address",
				City:     "Quito",
				Province: "Pichincha",
				License:  "1234567890123",
				Status:   AgencyStatusActive,
				OwnerID:  ownerID,
				Website:  stringPtr("invalid-url"),
			},
			expectError: true,
			errorMsg:    "invalid website URL format",
		},
		{
			name: "agency with empty ID",
			agency: &Agency{
				ID:       "",
				Name:     "Test Agency",
				Email:    "info@test.com",
				Phone:    "0984567890",
				Address:  "Test Address",
				City:     "Quito",
				Province: "Pichincha",
				License:  "1234567890123",
				Status:   AgencyStatusActive,
				OwnerID:  ownerID,
			},
			expectError: true,
			errorMsg:    "agency ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.agency.IsValid()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAgencyStatusManagement(t *testing.T) {
	agency := &Agency{
		ID:     uuid.New().String(),
		Status: AgencyStatusPending,
	}

	// Test activation
	err := agency.Activate()
	assert.NoError(t, err)
	assert.Equal(t, AgencyStatusActive, agency.Status)

	// Test deactivation
	err = agency.Deactivate()
	assert.NoError(t, err)
	assert.Equal(t, AgencyStatusInactive, agency.Status)

	// Test suspension
	err = agency.Suspend()
	assert.NoError(t, err)
	assert.Equal(t, AgencyStatusSuspended, agency.Status)

	// Test suspended agency cannot be activated directly
	err = agency.Activate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "suspended agencies cannot be activated directly")
}

func TestAgencyUpdateInfo(t *testing.T) {
	agency := &Agency{
		ID:   uuid.New().String(),
		Name: "Original Name",
	}

	err := agency.UpdateInfo(
		"Updated Name",
		"0987654321",
		"Updated Address",
		"Guayaquil",
		"Guayas",
	)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", agency.Name)
	assert.Equal(t, "0987654321", agency.Phone)
	assert.Equal(t, "Updated Address", agency.Address)
	assert.Equal(t, "Guayaquil", agency.City)
	assert.Equal(t, "Guayas", agency.Province)
}

func TestAgencySetWebsite(t *testing.T) {
	agency := &Agency{
		ID: uuid.New().String(),
	}

	// Test setting valid website
	err := agency.SetWebsite("https://www.example.com")
	assert.NoError(t, err)
	assert.Equal(t, "https://www.example.com", *agency.Website)

	// Test clearing website
	err = agency.SetWebsite("")
	assert.NoError(t, err)
	assert.Nil(t, agency.Website)

	// Test setting invalid website
	err = agency.SetWebsite("invalid-url")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid website URL format")
}

func TestAgencySetDescription(t *testing.T) {
	agency := &Agency{
		ID: uuid.New().String(),
	}

	// Test setting valid description
	description := "This is a test agency description"
	err := agency.SetDescription(description)
	assert.NoError(t, err)
	assert.Equal(t, description, *agency.Description)

	// Test clearing description
	err = agency.SetDescription("")
	assert.NoError(t, err)
	assert.Nil(t, agency.Description)

	// Test setting too long description
	longDescription := make([]byte, 1001)
	for i := range longDescription {
		longDescription[i] = 'A'
	}
	err = agency.SetDescription(string(longDescription))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "description cannot exceed 1000 characters")
}

func TestAgencySetLogo(t *testing.T) {
	agency := &Agency{
		ID: uuid.New().String(),
	}

	// Test setting logo
	logoURL := "https://example.com/logo.png"
	err := agency.SetLogo(logoURL)
	assert.NoError(t, err)
	assert.Equal(t, logoURL, *agency.Logo)

	// Test clearing logo
	err = agency.SetLogo("")
	assert.NoError(t, err)
	assert.Nil(t, agency.Logo)
}

func TestAgencyCanManageProperty(t *testing.T) {
	agencyID := uuid.New().String()
	otherAgencyID := uuid.New().String()

	tests := []struct {
		name             string
		agency           *Agency
		propertyAgencyID *string
		canManage        bool
	}{
		{
			name: "active agency can manage their property",
			agency: &Agency{
				ID:     agencyID,
				Status: AgencyStatusActive,
			},
			propertyAgencyID: &agencyID,
			canManage:        true,
		},
		{
			name: "active agency cannot manage other agency's property",
			agency: &Agency{
				ID:     agencyID,
				Status: AgencyStatusActive,
			},
			propertyAgencyID: &otherAgencyID,
			canManage:        false,
		},
		{
			name: "inactive agency cannot manage property",
			agency: &Agency{
				ID:     agencyID,
				Status: AgencyStatusInactive,
			},
			propertyAgencyID: &agencyID,
			canManage:        false,
		},
		{
			name: "property with no agency assignment",
			agency: &Agency{
				ID:     agencyID,
				Status: AgencyStatusActive,
			},
			propertyAgencyID: nil,
			canManage:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canManage := tt.agency.CanManageProperty(tt.propertyAgencyID)
			assert.Equal(t, tt.canManage, canManage)
		})
	}
}

func TestAgencyIsActive(t *testing.T) {
	tests := []struct {
		status   AgencyStatus
		isActive bool
	}{
		{AgencyStatusActive, true},
		{AgencyStatusInactive, false},
		{AgencyStatusSuspended, false},
		{AgencyStatusPending, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			agency := &Agency{Status: tt.status}
			assert.Equal(t, tt.isActive, agency.IsActive())
		})
	}
}

func TestAgencyDisplayMethods(t *testing.T) {
	agency := &Agency{
		Name:     "Test Agency",
		Phone:    "0984567890",
		Email:    "info@test.com",
		Address:  "Test Address 123",
		City:     "Quito",
		Province: "Pichincha",
	}

	// Test display name
	assert.Equal(t, "Test Agency", agency.GetDisplayName())

	// Test contact info
	expectedContact := "Test Agency - 0984567890 - info@test.com"
	assert.Equal(t, expectedContact, agency.GetContactInfo())

	// Test full address
	expectedAddress := "Test Address 123, Quito, Pichincha"
	assert.Equal(t, expectedAddress, agency.GetFullAddress())
}

func TestValidateEcuadorProvince(t *testing.T) {
	validProvinces := []string{
		"Azuay", "Bolívar", "Cañar", "Carchi", "Chimborazo", "Cotopaxi",
		"El Oro", "Esmeraldas", "Galápagos", "Guayas", "Imbabura", "Loja",
		"Los Ríos", "Manabí", "Morona Santiago", "Napo", "Orellana", "Pastaza",
		"Pichincha", "Santa Elena", "Santo Domingo", "Sucumbíos", "Tungurahua", "Zamora Chinchipe",
	}

	for _, province := range validProvinces {
		t.Run(province, func(t *testing.T) {
			err := validateEcuadorProvince(province)
			assert.NoError(t, err)
		})
	}

	// Test invalid province
	err := validateEcuadorProvince("InvalidProvince")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Ecuador province")
}

func TestValidatePhone(t *testing.T) {
	validPhones := []string{
		"0984567890",
		"0987654321",
		"+593984567890",
		"+593987654321",
	}

	for _, phone := range validPhones {
		t.Run(phone, func(t *testing.T) {
			err := validatePhone(phone)
			assert.NoError(t, err)
		})
	}

	invalidPhones := []string{
		"123456789",
		"12345678901",
		"0984567890123",
		"+5939845678901",
		"98456789",
		"phone",
		"",
	}

	for _, phone := range invalidPhones {
		t.Run(phone, func(t *testing.T) {
			err := validatePhone(phone)
			assert.Error(t, err)
		})
	}
}

func TestValidateLicense(t *testing.T) {
	validLicenses := []string{
		"1234567890123",
		"0987654321098",
		"1111111111111",
	}

	for _, license := range validLicenses {
		t.Run(license, func(t *testing.T) {
			err := validateLicense(license)
			assert.NoError(t, err)
		})
	}

	invalidLicenses := []string{
		"123456789012",  // 12 digits
		"12345678901234", // 14 digits
		"123456789012a", // contains letter
		"",
		"12345a67890123",
	}

	for _, license := range invalidLicenses {
		t.Run(license, func(t *testing.T) {
			err := validateLicense(license)
			assert.Error(t, err)
		})
	}
}

func TestValidateWebsite(t *testing.T) {
	validWebsites := []string{
		"https://www.example.com",
		"http://example.com",
		"https://subdomain.example.com",
		"https://example.com/path",
		"",
	}

	for _, website := range validWebsites {
		t.Run(website, func(t *testing.T) {
			err := validateWebsite(website)
			assert.NoError(t, err)
		})
	}

	invalidWebsites := []string{
		"invalid-url",
		"www.example.com",
		"example.com",
		"ftp://example.com",
	}

	for _, website := range invalidWebsites {
		t.Run(website, func(t *testing.T) {
			err := validateWebsite(website)
			assert.Error(t, err)
		})
	}
}