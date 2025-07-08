package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// MockFTSPropertyRepository is a mock for FTS-specific functionality
type MockFTSPropertyRepository struct {
	mock.Mock
}

func (m *MockFTSPropertyRepository) Create(property *domain.Property) error {
	args := m.Called(property)
	return args.Error(0)
}

func (m *MockFTSPropertyRepository) GetByID(id string) (*domain.Property, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockFTSPropertyRepository) GetBySlug(slug string) (*domain.Property, error) {
	args := m.Called(slug)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockFTSPropertyRepository) GetAll() ([]domain.Property, error) {
	args := m.Called()
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockFTSPropertyRepository) Update(property *domain.Property) error {
	args := m.Called(property)
	return args.Error(0)
}

func (m *MockFTSPropertyRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockFTSPropertyRepository) GetByProvince(province string) ([]domain.Property, error) {
	args := m.Called(province)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockFTSPropertyRepository) GetByPriceRange(minPrice, maxPrice float64) ([]domain.Property, error) {
	args := m.Called(minPrice, maxPrice)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockFTSPropertyRepository) SearchProperties(query string, limit int) ([]domain.Property, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockFTSPropertyRepository) SearchPropertiesRanked(query string, limit int) ([]repository.PropertySearchResult, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]repository.PropertySearchResult), args.Error(1)
}

func (m *MockFTSPropertyRepository) GetSearchSuggestions(query string, limit int) ([]repository.SearchSuggestion, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]repository.SearchSuggestion), args.Error(1)
}

func (m *MockFTSPropertyRepository) AdvancedSearch(params repository.AdvancedSearchParams) ([]repository.PropertySearchResult, error) {
	args := m.Called(params)
	return args.Get(0).([]repository.PropertySearchResult), args.Error(1)
}

func TestPropertyService_SearchProperties_FTS(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		mockSetup     func(*MockFTSPropertyRepository)
		expectedCount int
		wantError     bool
		errorContains string
	}{
		{
			name:  "successful FTS search",
			query: "casa moderna piscina",
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				properties := []domain.Property{
					*createTestProperty(),
				}
				mockRepo.On("SearchProperties", "casa moderna piscina", 50).Return(properties, nil)
			},
			expectedCount: 1,
			wantError:     false,
		},
		{
			name:  "empty query returns all properties",
			query: "",
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				properties := []domain.Property{
					*createTestProperty(),
					*createTestProperty(),
				}
				mockRepo.On("GetAll").Return(properties, nil)
			},
			expectedCount: 2,
			wantError:     false,
		},
		{
			name:  "query too short",
			query: "a",
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed as validation fails before repo call
			},
			wantError:     true,
			errorContains: "search query must be at least 2 characters",
		},
		{
			name:  "repository error",
			query: "casa error",
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				mockRepo.On("SearchProperties", "casa error", 50).Return([]domain.Property{}, errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error performing search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockFTSPropertyRepository{}
			tt.mockSetup(mockRepo)

			service := NewPropertyService(mockRepo)
			results, err := service.SearchProperties(tt.query)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_SearchPropertiesRanked(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		limit         int
		mockSetup     func(*MockFTSPropertyRepository)
		expectedCount int
		wantError     bool
		errorContains string
	}{
		{
			name:  "successful ranked search",
			query: "casa lujo",
			limit: 10,
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				results := []repository.PropertySearchResult{
					{
						Property: *createTestProperty(),
						Rank:     0.95,
					},
				}
				mockRepo.On("SearchPropertiesRanked", "casa lujo", 10).Return(results, nil)
			},
			expectedCount: 1,
			wantError:     false,
		},
		{
			name:  "empty query error",
			query: "",
			limit: 10,
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed
			},
			wantError:     true,
			errorContains: "search query required",
		},
		{
			name:  "query too short",
			query: "a",
			limit: 10,
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed
			},
			wantError:     true,
			errorContains: "search query must be at least 2 characters",
		},
		{
			name:  "invalid limit gets normalized",
			query: "casa test",
			limit: 0,
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				results := []repository.PropertySearchResult{}
				mockRepo.On("SearchPropertiesRanked", "casa test", 50).Return(results, nil) // Limit normalized to 50
			},
			expectedCount: 0,
			wantError:     false,
		},
		{
			name:  "repository error",
			query: "casa error",
			limit: 10,
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				mockRepo.On("SearchPropertiesRanked", "casa error", 10).Return([]repository.PropertySearchResult{}, errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error performing ranked search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockFTSPropertyRepository{}
			tt.mockSetup(mockRepo)

			service := NewPropertyService(mockRepo)
			results, err := service.SearchPropertiesRanked(tt.query, tt.limit)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
			if len(results) > 0 {
				assert.GreaterOrEqual(t, results[0].Rank, 0.0)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_GetSearchSuggestions(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		limit         int
		mockSetup     func(*MockFTSPropertyRepository)
		expectedCount int
		wantError     bool
	}{
		{
			name:  "successful suggestions",
			query: "gua",
			limit: 5,
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				suggestions := []repository.SearchSuggestion{
					{Text: "Guayas", Category: "province", Frequency: 150},
					{Text: "Guayaquil", Category: "city", Frequency: 120},
				}
				mockRepo.On("GetSearchSuggestions", "gua", 5).Return(suggestions, nil)
			},
			expectedCount: 2,
			wantError:     false,
		},
		{
			name:  "empty query returns empty suggestions",
			query: "",
			limit: 5,
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed
			},
			expectedCount: 0,
			wantError:     false,
		},
		{
			name:  "invalid limit gets normalized",
			query: "test",
			limit: 0,
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				suggestions := []repository.SearchSuggestion{}
				mockRepo.On("GetSearchSuggestions", "test", 10).Return(suggestions, nil) // Limit normalized to 10
			},
			expectedCount: 0,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockFTSPropertyRepository{}
			tt.mockSetup(mockRepo)

			service := NewPropertyService(mockRepo)
			suggestions, err := service.GetSearchSuggestions(tt.query, tt.limit)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, suggestions, tt.expectedCount)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_AdvancedSearch(t *testing.T) {
	tests := []struct {
		name          string
		params        repository.AdvancedSearchParams
		mockSetup     func(*MockFTSPropertyRepository)
		expectedCount int
		wantError     bool
		errorContains string
	}{
		{
			name: "successful advanced search",
			params: repository.AdvancedSearchParams{
				Query:       "casa moderna",
				Province:    "Guayas",
				MinPrice:    200000,
				MaxPrice:    500000,
				MinBedrooms: 3,
				MaxBedrooms: 5,
				Limit:       10,
			},
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				results := []repository.PropertySearchResult{
					{
						Property: *createTestProperty(),
						Rank:     0.88,
					},
				}
				mockRepo.On("AdvancedSearch", mock.AnythingOfType("repository.AdvancedSearchParams")).Return(results, nil)
			},
			expectedCount: 1,
			wantError:     false,
		},
		{
			name: "invalid price range",
			params: repository.AdvancedSearchParams{
				MinPrice: 500000,
				MaxPrice: 200000,
			},
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed
			},
			wantError:     true,
			errorContains: "minimum price cannot be greater than maximum price",
		},
		{
			name: "negative prices",
			params: repository.AdvancedSearchParams{
				MinPrice: -100000,
				MaxPrice: 500000,
			},
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed
			},
			wantError:     true,
			errorContains: "prices must be positive",
		},
		{
			name: "invalid province",
			params: repository.AdvancedSearchParams{
				Province: "InvalidProvince",
				Limit:    10,
			},
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed
			},
			wantError:     true,
			errorContains: "invalid province",
		},
		{
			name: "invalid property type",
			params: repository.AdvancedSearchParams{
				Type:  "invalid_type",
				Limit: 10,
			},
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed
			},
			wantError:     true,
			errorContains: "invalid property type",
		},
		{
			name: "search query too short",
			params: repository.AdvancedSearchParams{
				Query: "a",
				Limit: 10,
			},
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				// No mock setup needed
			},
			wantError:     true,
			errorContains: "search query must be at least 2 characters",
		},
		{
			name: "repository error",
			params: repository.AdvancedSearchParams{
				Query: "casa",
				Limit: 10,
			},
			mockSetup: func(mockRepo *MockFTSPropertyRepository) {
				mockRepo.On("AdvancedSearch", mock.AnythingOfType("repository.AdvancedSearchParams")).Return([]repository.PropertySearchResult{}, errors.New("database error"))
			},
			wantError:     true,
			errorContains: "error performing advanced search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockFTSPropertyRepository{}
			tt.mockSetup(mockRepo)

			service := NewPropertyService(mockRepo)
			results, err := service.AdvancedSearch(tt.params)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPropertyService_AdvancedSearch_ParameterNormalization(t *testing.T) {
	mockRepo := &MockFTSPropertyRepository{}
	
	// Test parameter normalization
	params := repository.AdvancedSearchParams{
		Query: "casa",
		Limit: 0, // Should be normalized to 50
	}

	expectedParams := repository.AdvancedSearchParams{
		Query: "casa",
		Limit: 50, // Normalized
	}

	mockRepo.On("AdvancedSearch", expectedParams).Return([]repository.PropertySearchResult{}, nil)

	service := NewPropertyService(mockRepo)
	_, err := service.AdvancedSearch(params)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPropertyService_SearchPropertiesRanked_LimitNormalization(t *testing.T) {
	tests := []struct {
		name          string
		inputLimit    int
		expectedLimit int
	}{
		{
			name:          "zero limit normalized to 50",
			inputLimit:    0,
			expectedLimit: 50,
		},
		{
			name:          "negative limit normalized to 50",
			inputLimit:    -5,
			expectedLimit: 50,
		},
		{
			name:          "excessive limit normalized to 50",
			inputLimit:    150,
			expectedLimit: 50,
		},
		{
			name:          "valid limit preserved",
			inputLimit:    25,
			expectedLimit: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockFTSPropertyRepository{}
			mockRepo.On("SearchPropertiesRanked", "test", tt.expectedLimit).Return([]repository.PropertySearchResult{}, nil)

			service := NewPropertyService(mockRepo)
			_, err := service.SearchPropertiesRanked("test", tt.inputLimit)

			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}