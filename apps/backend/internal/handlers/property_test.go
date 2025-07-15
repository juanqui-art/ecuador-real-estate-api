package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// MockPropertyService is a mock implementation of PropertyServiceInterface
type MockPropertyService struct {
	mock.Mock
}

func (m *MockPropertyService) CreateProperty(title, description, province, city, propertyType string, price float64, parkingSpaces int) (*domain.Property, error) {
	args := m.Called(title, description, province, city, propertyType, price, parkingSpaces)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockPropertyService) GetProperty(id string) (*domain.Property, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockPropertyService) GetPropertyBySlug(slug string) (*domain.Property, error) {
	args := m.Called(slug)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockPropertyService) ListProperties() ([]domain.Property, error) {
	args := m.Called()
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyService) UpdateProperty(id, title, description, province, city, propertyType string, price float64) (*domain.Property, error) {
	args := m.Called(id, title, description, province, city, propertyType, price)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockPropertyService) DeleteProperty(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPropertyService) FilterByProvince(province string) ([]domain.Property, error) {
	args := m.Called(province)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyService) FilterByPriceRange(minPrice, maxPrice float64) ([]domain.Property, error) {
	args := m.Called(minPrice, maxPrice)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyService) GetStatistics() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockPropertyService) SetPropertyLocation(id string, latitude, longitude float64, precision string) error {
	args := m.Called(id, latitude, longitude, precision)
	return args.Error(0)
}

func (m *MockPropertyService) SetPropertyFeatured(id string, featured bool) error {
	args := m.Called(id, featured)
	return args.Error(0)
}

func (m *MockPropertyService) AddPropertyTag(id, tag string) error {
	args := m.Called(id, tag)
	return args.Error(0)
}

func (m *MockPropertyService) SetPropertyParkingSpaces(id string, parkingSpaces int) error {
	args := m.Called(id, parkingSpaces)
	return args.Error(0)
}

func (m *MockPropertyService) SearchProperties(query string) ([]domain.Property, error) {
	args := m.Called(query)
	return args.Get(0).([]domain.Property), args.Error(1)
}

// Enhanced search methods for FTS functionality
func (m *MockPropertyService) SearchPropertiesRanked(query string, limit int) ([]repository.PropertySearchResult, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]repository.PropertySearchResult), args.Error(1)
}

func (m *MockPropertyService) GetSearchSuggestions(query string, limit int) ([]repository.SearchSuggestion, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]repository.SearchSuggestion), args.Error(1)
}

func (m *MockPropertyService) AdvancedSearch(params repository.AdvancedSearchParams) ([]repository.PropertySearchResult, error) {
	args := m.Called(params)
	return args.Get(0).([]repository.PropertySearchResult), args.Error(1)
}

// Pagination methods for MockPropertyService
func (m *MockPropertyService) ListPropertiesPaginated(pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyService) FilterByProvincePaginated(province string, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(province, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyService) FilterByPriceRangePaginated(minPrice, maxPrice float64, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(minPrice, maxPrice, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyService) SearchPropertiesPaginated(query string, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(query, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyService) SearchPropertiesRankedPaginated(query string, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(query, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyService) AdvancedSearchPaginated(params repository.AdvancedSearchParams, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(params, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
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
		"owner-123",
	)
}

func TestNewPropertyHandler(t *testing.T) {
	mockService := &MockPropertyService{}
	handler := NewPropertyHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestPropertyHandler_CreateProperty(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful creation",
			method: http.MethodPost,
			requestBody: CreatePropertyRequest{
				Title:       "Beautiful house in Samborondón",
				Description: "Modern house with pool",
				Province:    "Guayas",
				City:        "Samborondón",
				Type:        "house",
				Price:       285000,
			},
			mockSetup: func(m *MockPropertyService) {
				property := createTestProperty()
				m.On("CreateProperty", "Beautiful house in Samborondón", "Modern house with pool", "Guayas", "Samborondón", "house", 285000.0, 0).
					Return(property, nil)
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Property created successfully", response.Message)
				
				// Verify property data
				propertyData, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "Beautiful house in Samborondón", propertyData["title"])
				assert.Equal(t, "Guayas", propertyData["province"])
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			requestBody:    nil,
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:        "invalid JSON",
			method:      http.MethodPost,
			requestBody: "invalid json",
			mockSetup:   func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},
		{
			name:   "service error",
			method: http.MethodPost,
			requestBody: CreatePropertyRequest{
				Title:       "Invalid Property",
				Description: "Description",
				Province:    "InvalidProvince",
				City:        "City",
				Type:        "house",
				Price:       100000,
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("CreateProperty", "Invalid Property", "Description", "InvalidProvince", "City", "house", 100000.0, 0).
					Return((*domain.Property)(nil), errors.New("invalid province: InvalidProvince"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid province: InvalidProvince",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			var body []byte
			var err error
			if tt.requestBody != nil {
				if str, ok := tt.requestBody.(string); ok {
					body = []byte(str)
				} else {
					body, err = json.Marshal(tt.requestBody)
					assert.NoError(t, err)
				}
			}

			req := httptest.NewRequest(tt.method, "/api/properties", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateProperty(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_GetProperty(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful retrieval",
			method: http.MethodGet,
			url:    "/api/properties/test-id",
			mockSetup: func(m *MockPropertyService) {
				property := createTestProperty()
				property.ID = "test-id"
				m.On("GetProperty", "test-id").Return(property, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Property retrieved successfully", response.Message)

				propertyData, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "test-id", propertyData["id"])
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			url:            "/api/properties/test-id",
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "route without ID",
			method:         http.MethodGet,
			url:            "/api/properties/",
			mockSetup:      func(m *MockPropertyService) {
				m.On("GetProperty", "properties").Return((*domain.Property)(nil), errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:   "property not found",
			method: http.MethodGet,
			url:    "/api/properties/nonexistent-id",
			mockSetup: func(m *MockPropertyService) {
				m.On("GetProperty", "nonexistent-id").
					Return((*domain.Property)(nil), errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:   "service error",
			method: http.MethodGet,
			url:    "/api/properties/test-id",
			mockSetup: func(m *MockPropertyService) {
				m.On("GetProperty", "test-id").
					Return((*domain.Property)(nil), errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.GetProperty(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_GetPropertyBySlug(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful retrieval",
			method: http.MethodGet,
			url:    "/api/properties/slug/beautiful-house-12345678",
			mockSetup: func(m *MockPropertyService) {
				property := createTestProperty()
				property.Slug = "beautiful-house-12345678"
				m.On("GetPropertyBySlug", "beautiful-house-12345678").Return(property, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Property retrieved by slug successfully", response.Message)

				propertyData, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "beautiful-house-12345678", propertyData["slug"])
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			url:            "/api/properties/slug/test-slug",
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "empty slug",
			method:         http.MethodGet,
			url:            "/api/properties/slug/",
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Property slug required",
		},
		{
			name:   "property not found",
			method: http.MethodGet,
			url:    "/api/properties/slug/nonexistent-slug",
			mockSetup: func(m *MockPropertyService) {
				m.On("GetPropertyBySlug", "nonexistent-slug").
					Return((*domain.Property)(nil), errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:   "invalid slug format",
			method: http.MethodGet,
			url:    "/api/properties/slug/invalid-slug",
			mockSetup: func(m *MockPropertyService) {
				m.On("GetPropertyBySlug", "invalid-slug").
					Return((*domain.Property)(nil), errors.New("invalid slug format"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid slug format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.GetPropertyBySlug(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_ListProperties(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful listing",
			method: http.MethodGet,
			mockSetup: func(m *MockPropertyService) {
				properties := []domain.Property{
					*createTestProperty(),
					*createTestProperty(),
				}
				m.On("ListProperties").Return(properties, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Properties retrieved successfully", response.Message)

				properties, ok := response.Data.([]interface{})
				assert.True(t, ok)
				assert.Len(t, properties, 2)
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:   "empty list",
			method: http.MethodGet,
			mockSetup: func(m *MockPropertyService) {
				m.On("ListProperties").Return([]domain.Property{}, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Properties retrieved successfully", response.Message)

				properties, ok := response.Data.([]interface{})
				assert.True(t, ok)
				assert.Len(t, properties, 0)
			},
		},
		{
			name:   "service error",
			method: http.MethodGet,
			mockSetup: func(m *MockPropertyService) {
				m.On("ListProperties").Return([]domain.Property{}, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			req := httptest.NewRequest(tt.method, "/api/properties", nil)
			rec := httptest.NewRecorder()

			handler.ListProperties(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_UpdateProperty(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		requestBody    interface{}
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful update",
			method: http.MethodPut,
			url:    "/api/properties/test-id",
			requestBody: CreatePropertyRequest{
				Title:       "Updated Beautiful house",
				Description: "Updated description",
				Province:    "Guayas",
				City:        "Samborondón",
				Type:        "house",
				Price:       300000,
			},
			mockSetup: func(m *MockPropertyService) {
				property := createTestProperty()
				property.Title = "Updated Beautiful house"
				property.Price = 300000
				m.On("UpdateProperty", "test-id", "Updated Beautiful house", "Updated description", "Guayas", "Samborondón", "house", 300000.0).
					Return(property, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Property updated successfully", response.Message)

				propertyData, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "Updated Beautiful house", propertyData["title"])
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			url:            "/api/properties/test-id",
			requestBody:    nil,
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:        "route without ID",
			method:      http.MethodPut,
			url:         "/api/properties/",
			requestBody: CreatePropertyRequest{},
			mockSetup: func(m *MockPropertyService) {
				m.On("UpdateProperty", "properties", "", "", "", "", "", 0.0).
					Return((*domain.Property)(nil), errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:        "invalid JSON",
			method:      http.MethodPut,
			url:         "/api/properties/test-id",
			requestBody: "invalid json",
			mockSetup:   func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},
		{
			name:   "property not found",
			method: http.MethodPut,
			url:    "/api/properties/nonexistent-id",
			requestBody: CreatePropertyRequest{
				Title:       "Updated title",
				Description: "Updated description",
				Province:    "Guayas",
				City:        "Samborondón",
				Type:        "house",
				Price:       300000,
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("UpdateProperty", "nonexistent-id", "Updated title", "Updated description", "Guayas", "Samborondón", "house", 300000.0).
					Return((*domain.Property)(nil), errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:   "validation error",
			method: http.MethodPut,
			url:    "/api/properties/test-id",
			requestBody: CreatePropertyRequest{
				Title:       "",
				Description: "Updated description",
				Province:    "Guayas",
				City:        "Samborondón",
				Type:        "house",
				Price:       300000,
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("UpdateProperty", "test-id", "", "Updated description", "Guayas", "Samborondón", "house", 300000.0).
					Return((*domain.Property)(nil), errors.New("title is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			var body []byte
			var err error
			if tt.requestBody != nil {
				if str, ok := tt.requestBody.(string); ok {
					body = []byte(str)
				} else {
					body, err = json.Marshal(tt.requestBody)
					assert.NoError(t, err)
				}
			}

			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.UpdateProperty(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_DeleteProperty(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful deletion",
			method: http.MethodDelete,
			url:    "/api/properties/test-id",
			mockSetup: func(m *MockPropertyService) {
				m.On("DeleteProperty", "test-id").Return(nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Property deleted successfully", response.Message)
				assert.Nil(t, response.Data)
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			url:            "/api/properties/test-id",
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:   "route without ID",
			method: http.MethodDelete,
			url:    "/api/properties/",
			mockSetup: func(m *MockPropertyService) {
				m.On("DeleteProperty", "properties").
					Return(errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:   "property not found",
			method: http.MethodDelete,
			url:    "/api/properties/nonexistent-id",
			mockSetup: func(m *MockPropertyService) {
				m.On("DeleteProperty", "nonexistent-id").
					Return(errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:   "service error",
			method: http.MethodDelete,
			url:    "/api/properties/test-id",
			mockSetup: func(m *MockPropertyService) {
				m.On("DeleteProperty", "test-id").
					Return(errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.DeleteProperty(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_FilterProperties(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "search by query",
			method: http.MethodGet,
			url:    "/api/properties/filter?q=beautiful",
			mockSetup: func(m *MockPropertyService) {
				properties := []domain.Property{*createTestProperty()}
				m.On("SearchProperties", "beautiful").Return(properties, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Properties filtered by search query", response.Message)

				properties, ok := response.Data.([]interface{})
				assert.True(t, ok)
				assert.Len(t, properties, 1)
			},
		},
		{
			name:   "filter by province",
			method: http.MethodGet,
			url:    "/api/properties/filter?province=Guayas",
			mockSetup: func(m *MockPropertyService) {
				properties := []domain.Property{*createTestProperty()}
				m.On("FilterByProvince", "Guayas").Return(properties, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Properties filtered by province", response.Message)
			},
		},
		{
			name:   "filter by price range",
			method: http.MethodGet,
			url:    "/api/properties/filter?min_price=100000&max_price=500000",
			mockSetup: func(m *MockPropertyService) {
				properties := []domain.Property{*createTestProperty()}
				m.On("FilterByPriceRange", 100000.0, 500000.0).Return(properties, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Properties filtered by price range", response.Message)
			},
		},
		{
			name:   "no filters - return all",
			method: http.MethodGet,
			url:    "/api/properties/filter",
			mockSetup: func(m *MockPropertyService) {
				properties := []domain.Property{*createTestProperty(), *createTestProperty()}
				m.On("ListProperties").Return(properties, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "All properties", response.Message)

				properties, ok := response.Data.([]interface{})
				assert.True(t, ok)
				assert.Len(t, properties, 2)
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			url:            "/api/properties/filter",
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:   "invalid min price",
			method: http.MethodGet,
			url:    "/api/properties/filter?min_price=invalid&max_price=500000",
			mockSetup: func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid minimum price",
		},
		{
			name:   "invalid max price",
			method: http.MethodGet,
			url:    "/api/properties/filter?min_price=100000&max_price=invalid",
			mockSetup: func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid maximum price",
		},
		{
			name:   "search service error",
			method: http.MethodGet,
			url:    "/api/properties/filter?q=test",
			mockSetup: func(m *MockPropertyService) {
				m.On("SearchProperties", "test").Return([]domain.Property{}, errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "database error",
		},
		{
			name:   "province filter error",
			method: http.MethodGet,
			url:    "/api/properties/filter?province=InvalidProvince",
			mockSetup: func(m *MockPropertyService) {
				m.On("FilterByProvince", "InvalidProvince").Return([]domain.Property{}, errors.New("invalid province"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid province",
		},
		{
			name:   "price range filter error",
			method: http.MethodGet,
			url:    "/api/properties/filter?min_price=500000&max_price=100000",
			mockSetup: func(m *MockPropertyService) {
				m.On("FilterByPriceRange", 500000.0, 100000.0).Return([]domain.Property{}, errors.New("minimum price cannot be greater"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "minimum price cannot be greater",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.FilterProperties(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_GetStatistics(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful statistics retrieval",
			method: http.MethodGet,
			mockSetup: func(m *MockPropertyService) {
				stats := map[string]interface{}{
					"total_properties": 10,
					"average_price":    250000.0,
					"by_type": map[string]int{
						"house":     6,
						"apartment": 4,
					},
					"by_province": map[string]int{
						"Guayas":    7,
						"Pichincha": 3,
					},
				}
				m.On("GetStatistics").Return(stats, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Statistics retrieved successfully", response.Message)

				stats, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(10), stats["total_properties"])
				assert.Equal(t, 250000.0, stats["average_price"])
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:   "service error",
			method: http.MethodGet,
			mockSetup: func(m *MockPropertyService) {
				m.On("GetStatistics").Return(map[string]interface{}{}, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			req := httptest.NewRequest(tt.method, "/api/properties/statistics", nil)
			rec := httptest.NewRecorder()

			handler.GetStatistics(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_SetPropertyLocation(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		requestBody    interface{}
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful location setting",
			method: http.MethodPost,
			url:    "/api/properties/test-id/location",
			requestBody: map[string]interface{}{
				"latitude":  -2.1667,
				"longitude": -79.9,
				"precision": "exact",
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("SetPropertyLocation", "test-id", -2.1667, -79.9, "exact").Return(nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Property location updated successfully", response.Message)
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			url:            "/api/properties/test-id/location",
			requestBody:    nil,
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:        "route without ID",
			method:      http.MethodPost,
			url:         "/api/properties//location",
			requestBody: map[string]interface{}{},
			mockSetup:   func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Property ID required",
		},
		{
			name:        "invalid JSON",
			method:      http.MethodPost,
			url:         "/api/properties/test-id/location",
			requestBody: "invalid json",
			mockSetup:   func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},
		{
			name:   "property not found",
			method: http.MethodPost,
			url:    "/api/properties/nonexistent-id/location",
			requestBody: map[string]interface{}{
				"latitude":  -2.1667,
				"longitude": -79.9,
				"precision": "exact",
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("SetPropertyLocation", "nonexistent-id", -2.1667, -79.9, "exact").
					Return(errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:   "invalid coordinates",
			method: http.MethodPost,
			url:    "/api/properties/test-id/location",
			requestBody: map[string]interface{}{
				"latitude":  40.7128, // New York coordinates
				"longitude": -74.0060,
				"precision": "exact",
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("SetPropertyLocation", "test-id", 40.7128, -74.0060, "exact").
					Return(errors.New("coordinates outside Ecuadorian territory"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "coordinates outside Ecuadorian territory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			var body []byte
			var err error
			if tt.requestBody != nil {
				if str, ok := tt.requestBody.(string); ok {
					body = []byte(str)
				} else {
					body, err = json.Marshal(tt.requestBody)
					assert.NoError(t, err)
				}
			}

			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.SetPropertyLocation(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_SetPropertyFeatured(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		requestBody    interface{}
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful featured setting",
			method: http.MethodPost,
			url:    "/api/properties/test-id/featured",
			requestBody: map[string]interface{}{
				"featured": true,
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("SetPropertyFeatured", "test-id", true).Return(nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Property featured status updated successfully", response.Message)
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			url:            "/api/properties/test-id/featured",
			requestBody:    nil,
			mockSetup:      func(m *MockPropertyService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:        "route without ID",
			method:      http.MethodPost,
			url:         "/api/properties//featured",
			requestBody: map[string]interface{}{},
			mockSetup:   func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Property ID required",
		},
		{
			name:        "invalid JSON",
			method:      http.MethodPost,
			url:         "/api/properties/test-id/featured",
			requestBody: "invalid json",
			mockSetup:   func(m *MockPropertyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},
		{
			name:   "property not found",
			method: http.MethodPost,
			url:    "/api/properties/nonexistent-id/featured",
			requestBody: map[string]interface{}{
				"featured": true,
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("SetPropertyFeatured", "nonexistent-id", true).
					Return(errors.New("property not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "property not found",
		},
		{
			name:   "service error",
			method: http.MethodPost,
			url:    "/api/properties/test-id/featured",
			requestBody: map[string]interface{}{
				"featured": true,
			},
			mockSetup: func(m *MockPropertyService) {
				m.On("SetPropertyFeatured", "test-id", true).
					Return(errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)
			handler := NewPropertyHandler(mockService)

			var body []byte
			var err error
			if tt.requestBody != nil {
				if str, ok := tt.requestBody.(string); ok {
					body = []byte(str)
				} else {
					body, err = json.Marshal(tt.requestBody)
					assert.NoError(t, err)
				}
			}

			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.SetPropertyFeatured(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedError  string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "successful health check",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Service is running correctly", response.Message)

				health, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "healthy", health["status"])
				assert.Equal(t, "real-estate-api", health["service"])
				assert.Equal(t, "1.0.0", health["version"])
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			handler := NewPropertyHandler(mockService)

			req := httptest.NewRequest(tt.method, "/api/health", nil)
			rec := httptest.NewRecorder()

			handler.HealthCheck(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_ExtractIDFromURL(t *testing.T) {
	handler := NewPropertyHandler(&MockPropertyService{})

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "valid path with ID",
			path:     "/api/properties/test-id-123",
			expected: "test-id-123",
		},
		{
			name:     "path with trailing slash",
			path:     "/api/properties/test-id-123/",
			expected: "test-id-123",
		},
		{
			name:     "nested path",
			path:     "/api/properties/test-id/location",
			expected: "location",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "",
		},
		{
			name:     "root path",
			path:     "/",
			expected: "",
		},
		{
			name:     "UUID format",
			path:     "/api/properties/550e8400-e29b-41d4-a716-446655440000",
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.extractIDFromURL(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPropertyHandler_ExtractSlugFromURL(t *testing.T) {
	handler := NewPropertyHandler(&MockPropertyService{})

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "valid slug path",
			path:     "/api/properties/slug/beautiful-house-12345678",
			expected: "beautiful-house-12345678",
		},
		{
			name:     "slug path with trailing slash",
			path:     "/api/properties/slug/beautiful-house-12345678/",
			expected: "beautiful-house-12345678",
		},
		{
			name:     "invalid path format",
			path:     "/api/properties/test-id",
			expected: "",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "",
		},
		{
			name:     "incomplete slug path",
			path:     "/api/properties/slug/",
			expected: "",
		},
		{
			name:     "complex slug",
			path:     "/api/properties/slug/casa-moderna-con-piscina-y-jardin-12345678",
			expected: "casa-moderna-con-piscina-y-jardin-12345678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.extractSlugFromURL(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPropertyHandler_ErrorResponse(t *testing.T) {
	mockService := &MockPropertyService{}
	mockService.On("GetProperty", "nonexistent").Return((*domain.Property)(nil), errors.New("property not found"))
	handler := NewPropertyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/api/properties/nonexistent", nil)
	rec := httptest.NewRecorder()

	// Call respondError directly by triggering a scenario that uses it
	handler.GetProperty(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var errorResp ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.False(t, errorResp.Success)
	assert.Contains(t, errorResp.Message, "property not found")
	
	mockService.AssertExpectations(t)
}

func TestPropertyHandler_SuccessResponse(t *testing.T) {
	mockService := &MockPropertyService{}
	handler := NewPropertyHandler(mockService)

	// Setup mock for successful response
	properties := []domain.Property{*createTestProperty()}
	mockService.On("ListProperties").Return(properties, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/properties", nil)
	rec := httptest.NewRecorder()

	handler.ListProperties(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var successResp SuccessResponse
	err := json.Unmarshal(rec.Body.Bytes(), &successResp)
	assert.NoError(t, err)
	assert.NotNil(t, successResp.Data)
	assert.Equal(t, "Properties retrieved successfully", successResp.Message)

	mockService.AssertExpectations(t)
}