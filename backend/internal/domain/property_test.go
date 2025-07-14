package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProperty(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		province    string
		city        string
		propType    string
		price       float64
		wantValid   bool
	}{
		{
			name:        "valid property",
			title:       "Beautiful house in Samborondon",
			description: "Modern house with pool",
			province:    "Guayas",
			city:        "Samborondón",
			propType:    "house",
			price:       285000,
			wantValid:   true,
		},
		{
			name:        "empty title",
			title:       "",
			description: "Modern house with pool",
			province:    "Guayas",
			city:        "Samborondón",
			propType:    "house",
			price:       285000,
			wantValid:   false,
		},
		{
			name:        "zero price",
			title:       "Beautiful house in Samborondon",
			description: "Modern house with pool",
			province:    "Guayas",
			city:        "Samborondón",
			propType:    "house",
			price:       0,
			wantValid:   false,
		},
		{
			name:        "empty province",
			title:       "Beautiful house in Samborondon",
			description: "Modern house with pool",
			province:    "",
			city:        "Samborondón",
			propType:    "house",
			price:       285000,
			wantValid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ownerID := uuid.New().String()
			property := NewProperty(tt.title, tt.description, tt.province, tt.city, tt.propType, tt.price, ownerID)
			
			// Basic assertions
			require.NotNil(t, property)
			assert.NotEmpty(t, property.ID)
			assert.NotEmpty(t, property.Slug)
			assert.Equal(t, tt.title, property.Title)
			assert.Equal(t, tt.description, property.Description)
			assert.Equal(t, tt.province, property.Province)
			assert.Equal(t, tt.city, property.City)
			assert.Equal(t, tt.propType, property.Type)
			assert.Equal(t, tt.price, property.Price)
			
			// Default values
			assert.Equal(t, StatusAvailable, property.Status)
			assert.Equal(t, PrecisionApproximate, property.LocationPrecision)
			assert.Equal(t, PropertyStatusUsed, property.PropertyStatus)
			assert.False(t, property.Featured)
			assert.Zero(t, property.ViewCount)
			assert.Empty(t, property.Images)
			assert.Empty(t, property.Tags)
			
			// Timestamps
			assert.WithinDuration(t, time.Now(), property.CreatedAt, time.Second)
			assert.WithinDuration(t, time.Now(), property.UpdatedAt, time.Second)
			
			// Validation
			assert.Equal(t, tt.wantValid, property.IsValid())
		})
	}
}

func TestProperty_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		property func() *Property
		want     bool
	}{
		{
			name: "valid property",
			property: func() *Property {
				ownerID := uuid.New().String()
				return NewProperty("Valid Title", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
			},
			want: true,
		},
		{
			name: "empty title",
			property: func() *Property {
				ownerID := uuid.New().String()
				p := NewProperty("Valid Title", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
				p.Title = ""
				return p
			},
			want: false,
		},
		{
			name: "zero price",
			property: func() *Property {
				ownerID := uuid.New().String()
				p := NewProperty("Valid Title", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
				p.Price = 0
				return p
			},
			want: false,
		},
		{
			name: "negative price",
			property: func() *Property {
				ownerID := uuid.New().String()
				p := NewProperty("Valid Title", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
				p.Price = -1000
				return p
			},
			want: false,
		},
		{
			name: "empty province",
			property: func() *Property {
				ownerID := uuid.New().String()
				p := NewProperty("Valid Title", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
				p.Province = ""
				return p
			},
			want: false,
		},
		{
			name: "empty city",
			property: func() *Property {
				ownerID := uuid.New().String()
				p := NewProperty("Valid Title", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
				p.City = ""
				return p
			},
			want: false,
		},
		{
			name: "empty type",
			property: func() *Property {
				ownerID := uuid.New().String()
				p := NewProperty("Valid Title", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
				p.Type = ""
				return p
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property := tt.property()
			assert.Equal(t, tt.want, property.IsValid())
		})
	}
}

func TestProperty_UpdateTimestamp(t *testing.T) {
	ownerID := uuid.New().String()
	property := NewProperty("Test", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
	originalTime := property.UpdatedAt
	
	// Wait a small amount to ensure time difference
	time.Sleep(time.Millisecond)
	
	property.UpdateTimestamp()
	
	assert.True(t, property.UpdatedAt.After(originalTime))
}

func TestProperty_SetLocation(t *testing.T) {
	ownerID := uuid.New().String()
	property := NewProperty("Test", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
	
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		precision string
		wantError bool
	}{
		{
			name:      "valid Ecuador coordinates",
			latitude:  -2.1667,    // Guayaquil
			longitude: -79.9,      // Guayaquil
			precision: PrecisionExact,
			wantError: false,
		},
		{
			name:      "coordinates outside Ecuador",
			latitude:  40.7128,    // New York
			longitude: -74.0060,   // New York
			precision: PrecisionExact,
			wantError: true,
		},
		{
			name:      "invalid precision",
			latitude:  -2.1667,
			longitude: -79.9,
			precision: "invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := property.SetLocation(tt.latitude, tt.longitude, tt.precision)
			
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.latitude, *property.Latitude)
				assert.Equal(t, tt.longitude, *property.Longitude)
				assert.Equal(t, tt.precision, property.LocationPrecision)
			}
		})
	}
}

func TestProperty_AddTag(t *testing.T) {
	ownerID := uuid.New().String()
	property := NewProperty("Test", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
	
	// Test adding valid tag
	property.AddTag("luxury")
	assert.Contains(t, property.Tags, "luxury")
	assert.Len(t, property.Tags, 1)
	
	// Test adding duplicate tag (should not add)
	property.AddTag("luxury")
	assert.Len(t, property.Tags, 1)
	
	// Test adding tag with different case
	property.AddTag("LUXURY")
	assert.Len(t, property.Tags, 1) // Should not add duplicate
	
	// Test adding empty tag
	property.AddTag("")
	assert.Len(t, property.Tags, 1) // Should not add empty
	
	// Test adding another tag
	property.AddTag("pool")
	assert.Contains(t, property.Tags, "pool")
	assert.Len(t, property.Tags, 2)
}

func TestProperty_HasTag(t *testing.T) {
	ownerID := uuid.New().String()
	property := NewProperty("Test", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
	property.AddTag("luxury")
	property.AddTag("pool")
	
	assert.True(t, property.HasTag("luxury"))
	assert.True(t, property.HasTag("LUXURY")) // Case insensitive
	assert.True(t, property.HasTag("pool"))
	assert.False(t, property.HasTag("garden"))
	assert.False(t, property.HasTag(""))
}

func TestProperty_SetFeatured(t *testing.T) {
	ownerID := uuid.New().String()
	property := NewProperty("Test", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
	assert.False(t, property.Featured)
	
	property.SetFeatured(true)
	assert.True(t, property.Featured)
	
	property.SetFeatured(false)
	assert.False(t, property.Featured)
}

func TestProperty_IncrementViews(t *testing.T) {
	ownerID := uuid.New().String()
	property := NewProperty("Test", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
	assert.Zero(t, property.ViewCount)
	
	property.IncrementViews()
	assert.Equal(t, 1, property.ViewCount)
	
	property.IncrementViews()
	assert.Equal(t, 2, property.ViewCount)
}

func TestIsValidProvince(t *testing.T) {
	tests := []struct {
		province string
		want     bool
	}{
		{"Guayas", true},
		{"Pichincha", true},
		{"Azuay", true},
		{"Invalid Province", false},
		{"", false},
		{"guayas", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.province, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidProvince(tt.province))
		})
	}
}

func TestIsValidEcuadorCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		want      bool
	}{
		{"Quito", -0.2295, -78.5243, true},
		{"Guayaquil", -2.1667, -79.9, true},
		{"Cuenca", -2.9001, -79.0059, true},
		{"North of Ecuador", 3.0, -78.0, false},
		{"South of Ecuador", -6.0, -78.0, false},
		{"East of Ecuador", -2.0, -70.0, false},
		{"West of Ecuador", -2.0, -95.0, false},
		{"New York", 40.7128, -74.0060, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidEcuadorCoordinates(tt.latitude, tt.longitude))
		})
	}
}

func TestIsValidLocationPrecision(t *testing.T) {
	tests := []struct {
		precision string
		want      bool
	}{
		{PrecisionExact, true},
		{PrecisionApproximate, true},
		{PrecisionSector, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.precision, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidLocationPrecision(tt.precision))
		})
	}
}

func TestIsValidPropertyType(t *testing.T) {
	tests := []struct {
		propType string
		want     bool
	}{
		{TypeHouse, true},
		{TypeApartment, true},
		{TypeLand, true},
		{TypeCommercial, true},
		{"invalid", false},
		{"", false},
		{"HOUSE", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.propType, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidPropertyType(tt.propType))
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name  string
		title string
		id    string
		want  string
	}{
		{
			name:  "basic title",
			title: "Beautiful House",
			id:    "12345678",
			want:  "beautiful-house-12345678",
		},
		{
			name:  "title with special characters",
			title: "Casa Hermosa con Piscina!!!",
			id:    "12345678",
			want:  "casa-hermosa-con-piscina-12345678",
		},
		{
			name:  "title with multiple spaces",
			title: "Beautiful    House    with    Pool",
			id:    "12345678",
			want:  "beautiful-house-with-pool-12345678",
		},
		{
			name:  "very long title",
			title: "This is a very long title that should be truncated because it exceeds the maximum length allowed",
			id:    "12345678",
			want:  "this-is-a-very-long-title-that-should-be-truncated-12345678",
		},
		{
			name:  "long UUID",
			title: "Test House",
			id:    "550e8400-e29b-41d4-a716-446655440000",
			want:  "test-house-550e8400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.title, tt.id)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestProperty_UpdateSlug(t *testing.T) {
	ownerID := uuid.New().String()
	property := NewProperty("Original Title", "Description", "Guayas", "Samborondón", "house", 100000, ownerID)
	originalSlug := property.Slug
	
	property.Title = "New Title"
	property.UpdateSlug()
	
	assert.NotEqual(t, originalSlug, property.Slug)
	assert.Contains(t, property.Slug, "new-title")
}

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		slug string
		want bool
	}{
		{"beautiful-house-12345", true},
		{"casa-moderna-abc123", true},
		{"test", true},
		{"", false},
		{"-invalid-start", false},
		{"invalid-end-", false},
		{"invalid--double-dash", true}, // This should actually be valid
		{"UPPERCASE", false},          // Uppercase not allowed in slugs
		{"with_underscore", false},    // Underscores not allowed
		{"with spaces", false},        // Spaces not allowed
	}

	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidSlug(tt.slug))
		})
	}
}

// Pagination tests

func TestNewPaginationParams(t *testing.T) {
	params := NewPaginationParams()
	
	assert.Equal(t, 1, params.Page)
	assert.Equal(t, 20, params.PageSize)
	assert.Equal(t, "created_at", params.SortBy)
	assert.True(t, params.SortDesc)
}

func TestPaginationParams_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		want     int
	}{
		{
			name:     "page 1",
			page:     1,
			pageSize: 20,
			want:     0,
		},
		{
			name:     "page 2",
			page:     2,
			pageSize: 20,
			want:     20,
		},
		{
			name:     "page 3 with different page size",
			page:     3,
			pageSize: 10,
			want:     20,
		},
		{
			name:     "page 0 should default to 1",
			page:     0,
			pageSize: 20,
			want:     0,
		},
		{
			name:     "negative page should default to 1",
			page:     -1,
			pageSize: 20,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &PaginationParams{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}
			assert.Equal(t, tt.want, params.GetOffset())
		})
	}
}

func TestPaginationParams_GetLimit(t *testing.T) {
	tests := []struct {
		name     string
		pageSize int
		want     int
	}{
		{
			name:     "valid page size",
			pageSize: 20,
			want:     20,
		},
		{
			name:     "zero page size should default to 20",
			pageSize: 0,
			want:     20,
		},
		{
			name:     "negative page size should default to 20",
			pageSize: -1,
			want:     20,
		},
		{
			name:     "page size over limit should be capped at 100",
			pageSize: 150,
			want:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &PaginationParams{
				PageSize: tt.pageSize,
			}
			assert.Equal(t, tt.want, params.GetLimit())
		})
	}
}

func TestPaginationParams_GetOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		sortBy   string
		sortDesc bool
		want     string
	}{
		{
			name:     "valid sort field ascending",
			sortBy:   "price",
			sortDesc: false,
			want:     "price ASC",
		},
		{
			name:     "valid sort field descending",
			sortBy:   "price",
			sortDesc: true,
			want:     "price DESC",
		},
		{
			name:     "invalid sort field should default to created_at",
			sortBy:   "invalid_field",
			sortDesc: false,
			want:     "created_at ASC",
		},
		{
			name:     "empty sort field should default to created_at",
			sortBy:   "",
			sortDesc: true,
			want:     "created_at DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &PaginationParams{
				SortBy:   tt.sortBy,
				SortDesc: tt.sortDesc,
			}
			assert.Equal(t, tt.want, params.GetOrderBy())
		})
	}
}

func TestPaginationParams_Validate(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		wantErr  bool
	}{
		{
			name:     "valid parameters",
			page:     1,
			pageSize: 20,
			wantErr:  false,
		},
		{
			name:     "page zero should error",
			page:     0,
			pageSize: 20,
			wantErr:  true,
		},
		{
			name:     "negative page should error",
			page:     -1,
			pageSize: 20,
			wantErr:  true,
		},
		{
			name:     "page size zero should error",
			page:     1,
			pageSize: 0,
			wantErr:  true,
		},
		{
			name:     "negative page size should error",
			page:     1,
			pageSize: -1,
			wantErr:  true,
		},
		{
			name:     "page size over limit should error",
			page:     1,
			pageSize: 150,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &PaginationParams{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}
			err := params.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name         string
		currentPage  int
		pageSize     int
		totalRecords int
		wantPages    int
		wantHasNext  bool
		wantHasPrev  bool
	}{
		{
			name:         "first page with more records",
			currentPage:  1,
			pageSize:     10,
			totalRecords: 25,
			wantPages:    3,
			wantHasNext:  true,
			wantHasPrev:  false,
		},
		{
			name:         "middle page",
			currentPage:  2,
			pageSize:     10,
			totalRecords: 25,
			wantPages:    3,
			wantHasNext:  true,
			wantHasPrev:  true,
		},
		{
			name:         "last page",
			currentPage:  3,
			pageSize:     10,
			totalRecords: 25,
			wantPages:    3,
			wantHasNext:  false,
			wantHasPrev:  true,
		},
		{
			name:         "exact page fit",
			currentPage:  2,
			pageSize:     10,
			totalRecords: 20,
			wantPages:    2,
			wantHasNext:  false,
			wantHasPrev:  true,
		},
		{
			name:         "empty result set",
			currentPage:  1,
			pageSize:     10,
			totalRecords: 0,
			wantPages:    1,
			wantHasNext:  false,
			wantHasPrev:  false,
		},
		{
			name:         "single page",
			currentPage:  1,
			pageSize:     10,
			totalRecords: 5,
			wantPages:    1,
			wantHasNext:  false,
			wantHasPrev:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagination := NewPagination(tt.currentPage, tt.pageSize, tt.totalRecords)
			
			assert.Equal(t, tt.currentPage, pagination.CurrentPage)
			assert.Equal(t, tt.pageSize, pagination.PageSize)
			assert.Equal(t, tt.totalRecords, pagination.TotalRecords)
			assert.Equal(t, tt.wantPages, pagination.TotalPages)
			assert.Equal(t, tt.wantHasNext, pagination.HasNext)
			assert.Equal(t, tt.wantHasPrev, pagination.HasPrev)
		})
	}
}

func TestPaginatedResponse(t *testing.T) {
	// Create test data
	properties := []Property{
		func() Property { ownerID := uuid.New().String(); return *NewProperty("House 1", "Description 1", "Guayas", "Samborondón", "house", 100000, ownerID) }(),
		func() Property { ownerID := uuid.New().String(); return *NewProperty("House 2", "Description 2", "Guayas", "Samborondón", "house", 200000, ownerID) }(),
	}
	
	pagination := NewPagination(1, 10, 25)
	
	response := &PaginatedResponse{
		Data:       properties,
		Pagination: pagination,
	}
	
	assert.NotNil(t, response.Data)
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, properties, response.Data)
	assert.Equal(t, pagination, response.Pagination)
}