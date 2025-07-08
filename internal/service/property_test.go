package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// MockPropertyRepository is a mock implementation of PropertyRepository
type MockPropertyRepository struct {
	mock.Mock
}

func (m *MockPropertyRepository) Create(property *domain.Property) error {
	args := m.Called(property)
	return args.Error(0)
}

func (m *MockPropertyRepository) GetByID(id string) (*domain.Property, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockPropertyRepository) GetBySlug(slug string) (*domain.Property, error) {
	args := m.Called(slug)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockPropertyRepository) GetAll() ([]domain.Property, error) {
	args := m.Called()
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyRepository) Update(property *domain.Property) error {
	args := m.Called(property)
	return args.Error(0)
}

func (m *MockPropertyRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPropertyRepository) GetByProvince(province string) ([]domain.Property, error) {
	args := m.Called(province)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyRepository) GetByPriceRange(minPrice, maxPrice float64) ([]domain.Property, error) {
	args := m.Called(minPrice, maxPrice)
	return args.Get(0).([]domain.Property), args.Error(1)
}

// FTS methods for MockPropertyRepository
func (m *MockPropertyRepository) SearchProperties(query string, limit int) ([]domain.Property, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyRepository) SearchPropertiesRanked(query string, limit int) ([]repository.PropertySearchResult, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]repository.PropertySearchResult), args.Error(1)
}

func (m *MockPropertyRepository) GetSearchSuggestions(query string, limit int) ([]repository.SearchSuggestion, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]repository.SearchSuggestion), args.Error(1)
}

func (m *MockPropertyRepository) AdvancedSearch(params repository.AdvancedSearchParams) ([]repository.PropertySearchResult, error) {
	args := m.Called(params)
	return args.Get(0).([]repository.PropertySearchResult), args.Error(1)
}

// Helper function to create a test property
func createTestProperty() *domain.Property {
	return domain.NewProperty(
		"Beautiful house in Samborondón",
		"Modern house with pool",
		"Guayas",
		"Samborondón",
		"house",
		285000,
	)
}

func TestNewPropertyService(t *testing.T) {
	mockRepo := &MockPropertyRepository{}
	service := NewPropertyService(mockRepo)
	
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestPropertyService_CreateProperty(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		description   string
		province      string
		city          string
		propertyType  string
		price         float64
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
	}{
		{
			name:         "valid property creation",
			title:        "Beautiful house in Samborondón",
			description:  "Modern house with pool",
			province:     "Guayas",
			city:         "Samborondón",
			propertyType: "house",
			price:        285000,
			mockSetup: func(m *MockPropertyRepository) {
				m.On("Create", mock.AnythingOfType("*domain.Property")).Return(nil)
			},
			wantError: false,
		},
		{
			name:          "empty title",
			title:         "",
			description:   "Modern house with pool",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "house",
			price:         285000,
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "title is required",
		},
		{
			name:          "short title",
			title:         "Short",
			description:   "Modern house with pool",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "house",
			price:         285000,
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "title must be at least 10 characters",
		},
		{
			name:          "invalid province",
			title:         "Beautiful house in Samborondón",
			description:   "Modern house with pool",
			province:      "InvalidProvince",
			city:          "Samborondón",
			propertyType:  "house",
			price:         285000,
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "invalid province",
		},
		{
			name:          "invalid property type",
			title:         "Beautiful house in Samborondón",
			description:   "Modern house with pool",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "invalid",
			price:         285000,
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "invalid property type",
		},
		{
			name:          "zero price",
			title:         "Beautiful house in Samborondón",
			description:   "Modern house with pool",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "house",
			price:         0,
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "price must be greater than 0",
		},
		{
			name:         "repository error",
			title:        "Beautiful house in Samborondón",
			description:  "Modern house with pool",
			province:     "Guayas",
			city:         "Samborondón",
			propertyType: "house",
			price:        285000,
			mockSetup: func(m *MockPropertyRepository) {
				m.On("Create", mock.AnythingOfType("*domain.Property")).Return(errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error creating property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			property, err := service.CreateProperty(
				tt.title,
				tt.description,
				tt.province,
				tt.city,
				tt.propertyType,
				tt.price,
			)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, property)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, property)
				assert.Equal(t, tt.title, property.Title)
				assert.Equal(t, tt.description, property.Description)
				assert.Equal(t, tt.province, property.Province)
				assert.Equal(t, tt.city, property.City)
				assert.Equal(t, tt.propertyType, property.Type)
				assert.Equal(t, tt.price, property.Price)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_GetProperty(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
	}{
		{
			name: "successful retrieval",
			id:   "test-id",
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Update", mock.AnythingOfType("*domain.Property")).Return(nil)
			},
			wantError: false,
		},
		{
			name:          "empty id",
			id:            "",
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "property ID required",
		},
		{
			name: "property not found",
			id:   "nonexistent-id",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetByID", "nonexistent-id").Return((*domain.Property)(nil), errors.New("property not found"))
			},
			wantError:     true,
			errorContains: "error retrieving property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			property, err := service.GetProperty(tt.id)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, property)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, property)
				// Verify that view count was incremented
				assert.Equal(t, 1, property.ViewCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_GetPropertyBySlug(t *testing.T) {
	tests := []struct {
		name          string
		slug          string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
	}{
		{
			name: "successful retrieval",
			slug: "beautiful-house-12345678",
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetBySlug", "beautiful-house-12345678").Return(property, nil)
				m.On("Update", mock.AnythingOfType("*domain.Property")).Return(nil)
			},
			wantError: false,
		},
		{
			name:          "empty slug",
			slug:          "",
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "property slug required",
		},
		{
			name:          "invalid slug format",
			slug:          "INVALID-SLUG",
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "invalid slug format",
		},
		{
			name: "property not found",
			slug: "nonexistent-slug-12345678",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetBySlug", "nonexistent-slug-12345678").Return((*domain.Property)(nil), errors.New("property not found"))
			},
			wantError:     true,
			errorContains: "error retrieving property by slug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			property, err := service.GetPropertyBySlug(tt.slug)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, property)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, property)
				// Verify that view count was incremented
				assert.Equal(t, 1, property.ViewCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_ListProperties(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
		expectedCount int
	}{
		{
			name: "successful listing",
			mockSetup: func(m *MockPropertyRepository) {
				properties := []domain.Property{
					*createTestProperty(),
					*createTestProperty(),
				}
				m.On("GetAll").Return(properties, nil)
			},
			wantError:     false,
			expectedCount: 2,
		},
		{
			name: "empty list",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetAll").Return([]domain.Property{}, nil)
			},
			wantError:     false,
			expectedCount: 0,
		},
		{
			name: "repository error",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetAll").Return([]domain.Property{}, errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error listing properties",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			properties, err := service.ListProperties()

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, properties)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, properties)
				assert.Len(t, properties, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_UpdateProperty(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		title         string
		description   string
		province      string
		city          string
		propertyType  string
		price         float64
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
	}{
		{
			name:         "successful update",
			id:           "test-id",
			title:        "Updated Beautiful house",
			description:  "Updated description",
			province:     "Guayas",
			city:         "Samborondón",
			propertyType: "house",
			price:        300000,
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Update", mock.AnythingOfType("*domain.Property")).Return(nil)
			},
			wantError: false,
		},
		{
			name:         "property not found",
			id:           "nonexistent-id",
			title:        "Updated Beautiful house",
			description:  "Updated description",
			province:     "Guayas",
			city:         "Samborondón",
			propertyType: "house",
			price:        300000,
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetByID", "nonexistent-id").Return((*domain.Property)(nil), errors.New("property not found"))
			},
			wantError:     true,
			errorContains: "property not found",
		},
		{
			name:          "invalid updated data",
			id:            "test-id",
			title:         "", // Invalid title
			description:   "Updated description",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "house",
			price:         300000,
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
			},
			wantError:     true,
			errorContains: "title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			property, err := service.UpdateProperty(
				tt.id,
				tt.title,
				tt.description,
				tt.province,
				tt.city,
				tt.propertyType,
				tt.price,
			)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, property)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, property)
				assert.Equal(t, tt.title, property.Title)
				assert.Equal(t, tt.description, property.Description)
				assert.Equal(t, tt.province, property.Province)
				assert.Equal(t, tt.city, property.City)
				assert.Equal(t, tt.propertyType, property.Type)
				assert.Equal(t, tt.price, property.Price)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_DeleteProperty(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
	}{
		{
			name: "successful deletion",
			id:   "test-id",
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Delete", "test-id").Return(nil)
			},
			wantError: false,
		},
		{
			name:          "empty id",
			id:            "",
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "property ID required",
		},
		{
			name: "property not found",
			id:   "nonexistent-id",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetByID", "nonexistent-id").Return((*domain.Property)(nil), errors.New("property not found"))
			},
			wantError:     true,
			errorContains: "property not found",
		},
		{
			name: "repository delete error",
			id:   "test-id",
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Delete", "test-id").Return(errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error deleting property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			err := service.DeleteProperty(tt.id)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_FilterByProvince(t *testing.T) {
	tests := []struct {
		name          string
		province      string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
		expectedCount int
	}{
		{
			name:     "successful filtering",
			province: "Guayas",
			mockSetup: func(m *MockPropertyRepository) {
				properties := []domain.Property{
					*createTestProperty(),
				}
				m.On("GetByProvince", "Guayas").Return(properties, nil)
			},
			wantError:     false,
			expectedCount: 1,
		},
		{
			name:          "empty province",
			province:      "",
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "province required",
		},
		{
			name:          "invalid province",
			province:      "InvalidProvince",
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "invalid province",
		},
		{
			name:     "repository error",
			province: "Guayas",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetByProvince", "Guayas").Return([]domain.Property{}, errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error filtering properties by province",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			properties, err := service.FilterByProvince(tt.province)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, properties)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, properties)
				assert.Len(t, properties, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_FilterByPriceRange(t *testing.T) {
	tests := []struct {
		name          string
		minPrice      float64
		maxPrice      float64
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
		expectedCount int
	}{
		{
			name:     "successful filtering",
			minPrice: 100000,
			maxPrice: 500000,
			mockSetup: func(m *MockPropertyRepository) {
				properties := []domain.Property{
					*createTestProperty(),
				}
				m.On("GetByPriceRange", 100000.0, 500000.0).Return(properties, nil)
			},
			wantError:     false,
			expectedCount: 1,
		},
		{
			name:          "negative min price",
			minPrice:      -1000,
			maxPrice:      500000,
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "prices must be positive",
		},
		{
			name:          "negative max price",
			minPrice:      100000,
			maxPrice:      -1000,
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "prices must be positive",
		},
		{
			name:          "min price greater than max price",
			minPrice:      500000,
			maxPrice:      100000,
			mockSetup:     func(m *MockPropertyRepository) {},
			wantError:     true,
			errorContains: "minimum price cannot be greater than maximum price",
		},
		{
			name:     "repository error",
			minPrice: 100000,
			maxPrice: 500000,
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetByPriceRange", 100000.0, 500000.0).Return([]domain.Property{}, errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error filtering properties by price range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			properties, err := service.FilterByPriceRange(tt.minPrice, tt.maxPrice)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, properties)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, properties)
				assert.Len(t, properties, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_GetStatistics(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
		validateStats func(map[string]interface{}) bool
	}{
		{
			name: "successful statistics calculation",
			mockSetup: func(m *MockPropertyRepository) {
				property1 := createTestProperty()
				property1.Type = "house"
				property1.Status = "available"
				property1.Province = "Guayas"
				property1.Price = 200000

				property2 := createTestProperty()
				property2.Type = "apartment"
				property2.Status = "sold"
				property2.Province = "Pichincha"
				property2.Price = 150000

				properties := []domain.Property{*property1, *property2}
				m.On("GetAll").Return(properties, nil)
			},
			wantError: false,
			validateStats: func(stats map[string]interface{}) bool {
				totalProps := stats["total_properties"].(int)
				avgPrice := stats["average_price"].(float64)
				byType := stats["by_type"].(map[string]int)
				byStatus := stats["by_status"].(map[string]int)
				byProvince := stats["by_province"].(map[string]int)

				return totalProps == 2 &&
					avgPrice == 175000.0 &&
					byType["house"] == 1 &&
					byType["apartment"] == 1 &&
					byStatus["available"] == 1 &&
					byStatus["sold"] == 1 &&
					byProvince["Guayas"] == 1 &&
					byProvince["Pichincha"] == 1
			},
		},
		{
			name: "empty properties list",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetAll").Return([]domain.Property{}, nil)
			},
			wantError: false,
			validateStats: func(stats map[string]interface{}) bool {
				totalProps := stats["total_properties"].(int)
				avgPrice := stats["average_price"].(float64)
				return totalProps == 0 && avgPrice == 0.0
			},
		},
		{
			name: "repository error",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetAll").Return([]domain.Property{}, errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error retrieving properties",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			stats, err := service.GetStatistics()

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, stats)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, stats)
				if tt.validateStats != nil {
					assert.True(t, tt.validateStats(stats))
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_SetPropertyLocation(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		latitude      float64
		longitude     float64
		precision     string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
	}{
		{
			name:      "successful location setting",
			id:        "test-id",
			latitude:  -2.1667,
			longitude: -79.9,
			precision: "exact",
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Update", mock.AnythingOfType("*domain.Property")).Return(nil)
			},
			wantError: false,
		},
		{
			name:      "property not found",
			id:        "nonexistent-id",
			latitude:  -2.1667,
			longitude: -79.9,
			precision: "exact",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetByID", "nonexistent-id").Return((*domain.Property)(nil), errors.New("property not found"))
			},
			wantError:     true,
			errorContains: "property not found",
		},
		{
			name:      "invalid coordinates",
			id:        "test-id",
			latitude:  40.7128, // New York coordinates
			longitude: -74.0060,
			precision: "exact",
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
			},
			wantError:     true,
			errorContains: "error setting location",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			err := service.SetPropertyLocation(tt.id, tt.latitude, tt.longitude, tt.precision)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_SetPropertyFeatured(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		featured      bool
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
	}{
		{
			name:     "successful featured setting",
			id:       "test-id",
			featured: true,
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Update", mock.AnythingOfType("*domain.Property")).Return(nil)
			},
			wantError: false,
		},
		{
			name:     "property not found",
			id:       "nonexistent-id",
			featured: true,
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetByID", "nonexistent-id").Return((*domain.Property)(nil), errors.New("property not found"))
			},
			wantError:     true,
			errorContains: "property not found",
		},
		{
			name:     "repository update error",
			id:       "test-id",
			featured: true,
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Update", mock.AnythingOfType("*domain.Property")).Return(errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error updating property featured status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			err := service.SetPropertyFeatured(tt.id, tt.featured)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_AddPropertyTag(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		tag           string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
	}{
		{
			name: "successful tag addition",
			id:   "test-id",
			tag:  "luxury",
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Update", mock.AnythingOfType("*domain.Property")).Return(nil)
			},
			wantError: false,
		},
		{
			name: "property not found",
			id:   "nonexistent-id",
			tag:  "luxury",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("GetByID", "nonexistent-id").Return((*domain.Property)(nil), errors.New("property not found"))
			},
			wantError:     true,
			errorContains: "property not found",
		},
		{
			name: "repository update error",
			id:   "test-id",
			tag:  "luxury",
			mockSetup: func(m *MockPropertyRepository) {
				property := createTestProperty()
				m.On("GetByID", "test-id").Return(property, nil)
				m.On("Update", mock.AnythingOfType("*domain.Property")).Return(errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error adding tag to property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			err := service.AddPropertyTag(tt.id, tt.tag)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_SearchProperties_Legacy(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		mockSetup     func(*MockPropertyRepository)
		wantError     bool
		errorContains string
		expectedCount int
	}{
		{
			name:  "successful search with results using FTS",
			query: "beautiful",
			mockSetup: func(m *MockPropertyRepository) {
				properties := []domain.Property{*createTestProperty()}
				m.On("SearchProperties", "beautiful", 50).Return(properties, nil)
			},
			wantError:     false,
			expectedCount: 1,
		},
		{
			name:  "empty query returns all",
			query: "",
			mockSetup: func(m *MockPropertyRepository) {
				properties := []domain.Property{
					*createTestProperty(),
					*createTestProperty(),
				}
				m.On("GetAll").Return(properties, nil)
			},
			wantError:     false,
			expectedCount: 2,
		},
		{
			name:  "query too short",
			query: "a",
			mockSetup: func(m *MockPropertyRepository) {
				// No mock setup needed as validation fails before repo call
			},
			wantError:     true,
			errorContains: "search query must be at least 2 characters",
		},
		{
			name:  "repository error with FTS",
			query: "test",
			mockSetup: func(m *MockPropertyRepository) {
				m.On("SearchProperties", "test", 50).Return([]domain.Property{}, errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error performing search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockPropertyRepository{}
			tt.mockSetup(mockRepo)
			service := NewPropertyService(mockRepo)

			properties, err := service.SearchProperties(tt.query)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, properties)
			} else {
				assert.NoError(t, err)
				assert.Len(t, properties, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_validatePropertyData(t *testing.T) {
	service := NewPropertyService(&MockPropertyRepository{})

	tests := []struct {
		name          string
		title         string
		province      string
		city          string
		propertyType  string
		price         float64
		wantError     bool
		errorContains string
	}{
		{
			name:         "valid data",
			title:        "Beautiful house in Samborondón",
			province:     "Guayas",
			city:         "Samborondón",
			propertyType: "house",
			price:        285000,
			wantError:    false,
		},
		{
			name:          "empty title",
			title:         "",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "house",
			price:         285000,
			wantError:     true,
			errorContains: "title is required",
		},
		{
			name:          "title too long",
			title:         string(make([]byte, 260)), // 260 characters
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "house",
			price:         285000,
			wantError:     true,
			errorContains: "title cannot exceed 255 characters",
		},
		{
			name:          "empty province",
			title:         "Beautiful house in Samborondón",
			province:      "",
			city:          "Samborondón",
			propertyType:  "house",
			price:         285000,
			wantError:     true,
			errorContains: "province is required",
		},
		{
			name:          "empty city",
			title:         "Beautiful house in Samborondón",
			province:      "Guayas",
			city:          "",
			propertyType:  "house",
			price:         285000,
			wantError:     true,
			errorContains: "city is required",
		},
		{
			name:          "empty property type",
			title:         "Beautiful house in Samborondón",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "",
			price:         285000,
			wantError:     true,
			errorContains: "property type is required",
		},
		{
			name:          "zero price",
			title:         "Beautiful house in Samborondón",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "house",
			price:         0,
			wantError:     true,
			errorContains: "price must be greater than 0",
		},
		{
			name:          "negative price",
			title:         "Beautiful house in Samborondón",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "house",
			price:         -1000,
			wantError:     true,
			errorContains: "price must be greater than 0",
		},
		{
			name:          "invalid province",
			title:         "Beautiful house in Samborondón",
			province:      "InvalidProvince",
			city:          "Samborondón",
			propertyType:  "house",
			price:         285000,
			wantError:     true,
			errorContains: "invalid province",
		},
		{
			name:          "invalid property type",
			title:         "Beautiful house in Samborondón",
			province:      "Guayas",
			city:          "Samborondón",
			propertyType:  "invalid",
			price:         285000,
			wantError:     true,
			errorContains: "invalid property type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validatePropertyData(tt.title, tt.province, tt.city, tt.propertyType, tt.price)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}