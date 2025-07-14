package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

func TestPropertyHandler_SearchRanked(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		wantError      bool
	}{
		{
			name: "successful ranked search",
			url:  "/api/properties/search/ranked?q=casa+moderna&limit=10",
			mockSetup: func(mockService *MockPropertyService) {
				results := []repository.PropertySearchResult{
					{
						Property: *createTestProperty(),
						Rank:     0.95,
					},
				}
				mockService.On("SearchPropertiesRanked", "casa moderna", 10).Return(results, nil)
			},
			expectedStatus: http.StatusOK,
			wantError:      false,
		},
		{
			name: "missing query parameter",
			url:  "/api/properties/search/ranked",
			mockSetup: func(mockService *MockPropertyService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			wantError:      true,
		},
		{
			name: "invalid limit parameter",
			url:  "/api/properties/search/ranked?q=casa&limit=invalid",
			mockSetup: func(mockService *MockPropertyService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			wantError:      true,
		},
		{
			name: "service error",
			url:  "/api/properties/search/ranked?q=casa&limit=5",
			mockSetup: func(mockService *MockPropertyService) {
				mockService.On("SearchPropertiesRanked", "casa", 5).Return([]repository.PropertySearchResult{}, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			wantError:      true,
		},
		{
			name: "default limit when not specified",
			url:  "/api/properties/search/ranked?q=casa",
			mockSetup: func(mockService *MockPropertyService) {
				results := []repository.PropertySearchResult{}
				mockService.On("SearchPropertiesRanked", "casa", 50).Return(results, nil) // Default limit 50
			},
			expectedStatus: http.StatusOK,
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)

			handler := NewPropertyHandler(mockService)
			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			handler.SearchRanked(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.wantError {
				var response SuccessResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Ranked search results retrieved successfully", response.Message)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_SearchSuggestions(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		wantError      bool
	}{
		{
			name: "successful suggestions",
			url:  "/api/properties/search/suggestions?q=gua&limit=5",
			mockSetup: func(mockService *MockPropertyService) {
				suggestions := []repository.SearchSuggestion{
					{Text: "Guayas", Category: "province", Frequency: 150},
					{Text: "Guayaquil", Category: "city", Frequency: 120},
				}
				mockService.On("GetSearchSuggestions", "gua", 5).Return(suggestions, nil)
			},
			expectedStatus: http.StatusOK,
			wantError:      false,
		},
		{
			name: "empty query",
			url:  "/api/properties/search/suggestions?q=",
			mockSetup: func(mockService *MockPropertyService) {
				suggestions := []repository.SearchSuggestion{}
				mockService.On("GetSearchSuggestions", "", 10).Return(suggestions, nil) // Default limit 10
			},
			expectedStatus: http.StatusOK,
			wantError:      false,
		},
		{
			name: "invalid limit parameter",
			url:  "/api/properties/search/suggestions?q=test&limit=invalid",
			mockSetup: func(mockService *MockPropertyService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			wantError:      true,
		},
		{
			name: "service error",
			url:  "/api/properties/search/suggestions?q=test",
			mockSetup: func(mockService *MockPropertyService) {
				mockService.On("GetSearchSuggestions", "test", 10).Return([]repository.SearchSuggestion{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)

			handler := NewPropertyHandler(mockService)
			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			handler.SearchSuggestions(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.wantError {
				var response SuccessResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Search suggestions retrieved successfully", response.Message)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_AdvancedSearch(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		wantError      bool
	}{
		{
			name: "successful advanced search",
			requestBody: map[string]interface{}{
				"query":         "casa moderna",
				"province":      "Guayas",
				"min_price":     200000,
				"max_price":     500000,
				"min_bedrooms":  3,
				"max_bedrooms":  5,
				"featured_only": false,
				"limit":         10,
			},
			mockSetup: func(mockService *MockPropertyService) {
				results := []repository.PropertySearchResult{
					{
						Property: *createTestProperty(),
						Rank:     0.88,
					},
				}
				mockService.On("AdvancedSearch", mock.AnythingOfType("repository.AdvancedSearchParams")).Return(results, nil)
			},
			expectedStatus: http.StatusOK,
			wantError:      false,
		},
		{
			name: "minimal search parameters",
			requestBody: map[string]interface{}{
				"query": "casa",
				"limit": 20,
			},
			mockSetup: func(mockService *MockPropertyService) {
				results := []repository.PropertySearchResult{}
				mockService.On("AdvancedSearch", mock.AnythingOfType("repository.AdvancedSearchParams")).Return(results, nil)
			},
			expectedStatus: http.StatusOK,
			wantError:      false,
		},
		{
			name:        "invalid JSON",
			requestBody: "invalid json",
			mockSetup: func(mockService *MockPropertyService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			wantError:      true,
		},
		{
			name: "service validation error",
			requestBody: map[string]interface{}{
				"query":     "a", // Too short
				"min_price": -100000,
				"limit":     10,
			},
			mockSetup: func(mockService *MockPropertyService) {
				mockService.On("AdvancedSearch", mock.AnythingOfType("repository.AdvancedSearchParams")).Return([]repository.PropertySearchResult{}, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			wantError:      true,
		},
		{
			name: "empty request body",
			requestBody: map[string]interface{}{
				"limit": 10,
			},
			mockSetup: func(mockService *MockPropertyService) {
				results := []repository.PropertySearchResult{}
				mockService.On("AdvancedSearch", mock.AnythingOfType("repository.AdvancedSearchParams")).Return(results, nil)
			},
			expectedStatus: http.StatusOK,
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.requestBody)
			}

			handler := NewPropertyHandler(mockService)
			req := httptest.NewRequest("POST", "/api/properties/search/advanced", &body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.AdvancedSearch(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.wantError {
				var response SuccessResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Advanced search results retrieved successfully", response.Message)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPropertyHandler_SearchRanked_MethodNotAllowed(t *testing.T) {
	mockService := &MockPropertyService{}
	handler := NewPropertyHandler(mockService)

	req := httptest.NewRequest("POST", "/api/properties/search/ranked?q=test", nil)
	rr := httptest.NewRecorder()

	handler.SearchRanked(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPropertyHandler_SearchSuggestions_MethodNotAllowed(t *testing.T) {
	mockService := &MockPropertyService{}
	handler := NewPropertyHandler(mockService)

	req := httptest.NewRequest("POST", "/api/properties/search/suggestions?q=test", nil)
	rr := httptest.NewRecorder()

	handler.SearchSuggestions(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPropertyHandler_AdvancedSearch_MethodNotAllowed(t *testing.T) {
	mockService := &MockPropertyService{}
	handler := NewPropertyHandler(mockService)

	req := httptest.NewRequest("GET", "/api/properties/search/advanced", nil)
	rr := httptest.NewRecorder()

	handler.AdvancedSearch(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPropertyHandler_FilterProperties_EnhancedSearch(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*MockPropertyService)
		expectedStatus int
		wantError      bool
	}{
		{
			name: "filter with enhanced search query",
			url:  "/api/properties/filter?q=casa+moderna+terraza",
			mockSetup: func(mockService *MockPropertyService) {
				properties := []domain.Property{*createTestProperty()}
				mockService.On("SearchProperties", "casa moderna terraza").Return(properties, nil)
			},
			expectedStatus: http.StatusOK,
			wantError:      false,
		},
		{
			name: "search query too short",
			url:  "/api/properties/filter?q=a",
			mockSetup: func(mockService *MockPropertyService) {
				mockService.On("SearchProperties", "a").Return([]domain.Property{}, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockPropertyService{}
			tt.mockSetup(mockService)

			handler := NewPropertyHandler(mockService)
			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			handler.FilterProperties(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.wantError {
				var response SuccessResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Properties filtered by search query", response.Message)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAdvancedSearchRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		requestBody map[string]interface{}
		shouldParse bool
	}{
		{
			name: "all fields provided",
			requestBody: map[string]interface{}{
				"query":         "casa",
				"province":      "Guayas",
				"city":          "Guayaquil",
				"type":          "house",
				"min_price":     100000,
				"max_price":     500000,
				"min_bedrooms":  2,
				"max_bedrooms":  5,
				"min_bathrooms": 1.5,
				"max_bathrooms": 4.0,
				"min_area":      80,
				"max_area":      400,
				"featured_only": true,
				"limit":         25,
			},
			shouldParse: true,
		},
		{
			name: "minimal fields",
			requestBody: map[string]interface{}{
				"query": "casa",
			},
			shouldParse: true,
		},
		{
			name: "numeric fields as strings should fail",
			requestBody: map[string]interface{}{
				"query":     "casa",
				"min_price": "100000", // String instead of number
			},
			shouldParse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			var req struct {
				Query        string  `json:"query"`
				Province     string  `json:"province"`
				City         string  `json:"city"`
				Type         string  `json:"type"`
				MinPrice     float64 `json:"min_price"`
				MaxPrice     float64 `json:"max_price"`
				MinBedrooms  int     `json:"min_bedrooms"`
				MaxBedrooms  int     `json:"max_bedrooms"`
				MinBathrooms float64 `json:"min_bathrooms"`
				MaxBathrooms float64 `json:"max_bathrooms"`
				MinArea      float64 `json:"min_area"`
				MaxArea      float64 `json:"max_area"`
				FeaturedOnly bool    `json:"featured_only"`
				Limit        int     `json:"limit"`
			}

			err := json.NewDecoder(&body).Decode(&req)

			if tt.shouldParse {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}