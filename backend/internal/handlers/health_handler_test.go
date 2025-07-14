package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"realty-core/internal/cache"
	"realty-core/internal/domain"
	"realty-core/internal/repository"
	"realty-core/internal/service"
)

// Mock implementations for testing

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Ping() error {
	args := m.Called()
	return args.Error(0)
}

type MockPropertyRepo struct {
	mock.Mock
}

func (m *MockPropertyRepo) Create(property *domain.Property) error {
	args := m.Called(property)
	return args.Error(0)
}

func (m *MockPropertyRepo) GetByID(id string) (*domain.Property, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockPropertyRepo) GetBySlug(slug string) (*domain.Property, error) {
	args := m.Called(slug)
	return args.Get(0).(*domain.Property), args.Error(1)
}

func (m *MockPropertyRepo) GetAll() ([]domain.Property, error) {
	args := m.Called()
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyRepo) Update(property *domain.Property) error {
	args := m.Called(property)
	return args.Error(0)
}

func (m *MockPropertyRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPropertyRepo) FilterByProvince(province string) ([]domain.Property, error) {
	args := m.Called(province)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyRepo) FilterByPriceRange(minPrice, maxPrice float64) ([]domain.Property, error) {
	args := m.Called(minPrice, maxPrice)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyRepo) SearchProperties(query string) ([]domain.Property, error) {
	args := m.Called(query)
	return args.Get(0).([]domain.Property), args.Error(1)
}

func (m *MockPropertyRepo) SearchPropertiesRanked(query string, limit int) ([]repository.PropertySearchResult, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]repository.PropertySearchResult), args.Error(1)
}

func (m *MockPropertyRepo) GetSearchSuggestions(query string, limit int) ([]repository.SearchSuggestion, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]repository.SearchSuggestion), args.Error(1)
}

func (m *MockPropertyRepo) AdvancedSearch(params repository.AdvancedSearchParams) ([]repository.PropertySearchResult, error) {
	args := m.Called(params)
	return args.Get(0).([]repository.PropertySearchResult), args.Error(1)
}

func (m *MockPropertyRepo) ListPropertiesPaginated(pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyRepo) FilterByProvincePaginated(province string, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(province, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyRepo) FilterByPriceRangePaginated(minPrice, maxPrice float64, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(minPrice, maxPrice, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyRepo) SearchPropertiesPaginated(query string, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(query, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyRepo) SearchPropertiesRankedPaginated(query string, pagination *domain.PaginationParams) (*domain.PaginatedResponse, error) {
	args := m.Called(query, pagination)
	return args.Get(0).(*domain.PaginatedResponse), args.Error(1)
}

func (m *MockPropertyRepo) AdvancedSearchPaginated(params repository.AdvancedSearchParams, pagination *domain.PaginationParams) ([]repository.PropertySearchResult, int, error) {
	args := m.Called(params, pagination)
	return args.Get(0).([]repository.PropertySearchResult), args.Int(1), args.Error(2)
}

// Simple mock for other repos
type MockImageRepo struct {
	mock.Mock
}

func (m *MockImageRepo) GetByPropertyID(propertyID string) ([]domain.ImageInfo, error) {
	args := m.Called(propertyID)
	return args.Get(0).([]domain.ImageInfo), args.Error(1)
}

type MockImageCache struct {
	mock.Mock
}

func (m *MockImageCache) IsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockImageCache) Stats() cache.ImageCacheStats {
	args := m.Called()
	return args.Get(0).(cache.ImageCacheStats)
}

func (m *MockImageCache) Get(key string) ([]byte, string, bool) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.String(1), args.Bool(2)
}

func (m *MockImageCache) Set(key string, data []byte, contentType string) {
	m.Called(key, data, contentType)
}

func (m *MockImageCache) Delete(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockImageCache) Clear() {
	m.Called()
}

func (m *MockImageCache) GetThumbnail(imageID string, size int) ([]byte, string, bool) {
	args := m.Called(imageID, size)
	return args.Get(0).([]byte), args.String(1), args.Bool(2)
}

func (m *MockImageCache) SetThumbnail(imageID string, size int, data []byte, contentType string) {
	m.Called(imageID, size, data, contentType)
}

func (m *MockImageCache) GetVariant(imageID string, width, height, quality int, format string) ([]byte, string, bool) {
	args := m.Called(imageID, width, height, quality, format)
	return args.Get(0).([]byte), args.String(1), args.Bool(2)
}

func (m *MockImageCache) SetVariant(imageID string, width, height, quality int, format string, data []byte, contentType string) {
	m.Called(imageID, width, height, quality, format, data, contentType)
}

func (m *MockImageCache) InvalidateImage(imageID string) int {
	args := m.Called(imageID)
	return args.Int(0)
}

func (m *MockImageCache) Size() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockImageCache) CurrentSize() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockImageCache) CleanupExpired() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockImageCache) GetPopularImages(limit int) []cache.PopularImage {
	args := m.Called(limit)
	return args.Get(0).([]cache.PopularImage)
}

// Test functions

func TestHealthHandler_BasicHealthCheck(t *testing.T) {
	// Setup
	handler := &HealthHandler{}
	
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	
	// Execute
	handler.BasicHealthCheck(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "realty-core", response["service"])
	assert.Equal(t, "1.9.0", response["version"])
	assert.NotNil(t, response["timestamp"])
}

func TestHealthHandler_ReadinessCheck_Healthy(t *testing.T) {
	// Setup
	mockDB := &sql.DB{} // This is a simplified mock
	handler := &HealthHandler{db: mockDB}
	
	req := httptest.NewRequest("GET", "/api/health/ready", nil)
	w := httptest.NewRecorder()
	
	// Execute
	handler.ReadinessCheck(w, req)
	
	// Assert
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.NotNil(t, response["ready"])
	assert.NotNil(t, response["checks"])
	assert.NotNil(t, response["timestamp"])
}

func TestHealthHandler_LivenessCheck(t *testing.T) {
	// Setup
	handler := &HealthHandler{}
	
	req := httptest.NewRequest("GET", "/api/health/live", nil)
	w := httptest.NewRecorder()
	
	// Execute
	handler.LivenessCheck(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, true, response["alive"])
	assert.NotNil(t, response["timestamp"])
	assert.NotNil(t, response["uptime"])
}

func TestHealthHandler_CacheHealth(t *testing.T) {
	// Setup mocks
	mockImageCache := &MockImageCache{}
	mockImageCache.On("IsEnabled").Return(true)
	mockImageCache.On("Stats").Return(cache.ImageCacheStats{
		CacheStats: cache.CacheStats{
			Hits:     100,
			Misses:   20,
			HitRate:  83.3,
			Size:     50,
			Capacity: 100,
		},
	})
	
	// Create a property service with cache (simplified for testing)
	// cacheConfig := cache.PropertyCacheConfig{
	// 	Enabled:  true,
	// 	Capacity: 100,
	// }
	// propertyCache := cache.NewPropertyCache(cacheConfig)
	
	// Create handler
	handler := &HealthHandler{
		imageCache:      mockImageCache,
		propertyService: &service.PropertyService{},
	}
	
	// Since we can't easily inject the property service with cache,
	// we'll use a simplified version
	handler.propertyService = nil // This will test the nil check
	
	req := httptest.NewRequest("GET", "/api/health/cache", nil)
	w := httptest.NewRecorder()
	
	// Execute
	handler.CacheHealth(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var response CacheHealthStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.True(t, response.ImageCache.Enabled)
	assert.Equal(t, 83.3, response.ImageCache.HitRate)
	assert.False(t, response.PropertyCache.Enabled) // Should be false when service is nil
	
	mockImageCache.AssertExpectations(t)
}

func TestHealthHandler_MetricsEndpoint(t *testing.T) {
	// Setup
	mockImageCache := &MockImageCache{}
	mockImageCache.On("IsEnabled").Return(true)
	mockImageCache.On("Stats").Return(cache.ImageCacheStats{
		CacheStats: cache.CacheStats{
			Hits:   100,
			Misses: 20,
		},
	})
	
	handler := &HealthHandler{
		imageCache:      mockImageCache,
		propertyService: nil, // Simplified for testing
	}
	
	req := httptest.NewRequest("GET", "/api/metrics", nil)
	w := httptest.NewRecorder()
	
	// Execute
	handler.MetricsEndpoint(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	
	body := w.Body.String()
	assert.Contains(t, body, "realty_core_uptime_seconds")
	assert.Contains(t, body, "realty_core_memory_alloc_bytes")
	assert.Contains(t, body, "realty_core_goroutines")
	assert.Contains(t, body, "realty_core_image_cache_hits_total 100")
	assert.Contains(t, body, "realty_core_image_cache_misses_total 20")
	
	mockImageCache.AssertExpectations(t)
}

func TestHealthHandler_DetailedHealthCheck_Healthy(t *testing.T) {
	// Setup simplified handler for testing basic functionality
	mockImageCache := &MockImageCache{}
	mockImageCache.On("IsEnabled").Return(true)
	mockImageCache.On("Stats").Return(cache.ImageCacheStats{
		CacheStats: cache.CacheStats{
			HitRate: 80.0,
		},
	})
	
	handler := &HealthHandler{
		db:           &sql.DB{}, // Simplified
		propertyRepo: nil,       // Simplified for basic test
		imageCache:   mockImageCache,
	}
	
	req := httptest.NewRequest("GET", "/api/health/detailed", nil)
	w := httptest.NewRecorder()
	
	// Execute
	handler.DetailedHealthCheck(w, req)
	
	// Assert
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	
	var response HealthStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.NotEmpty(t, response.Status)
	assert.Equal(t, "1.9.0", response.Version)
	assert.NotZero(t, response.Uptime)
	assert.NotNil(t, response.Services)
	assert.NotNil(t, response.System)
	
	// Check that system health is populated
	assert.NotZero(t, response.System.Memory.Alloc)
	assert.Greater(t, response.System.Goroutines, 0)
	
	mockImageCache.AssertExpectations(t)
}

func TestHealthHandler_GetSystemHealth(t *testing.T) {
	// Setup
	handler := &HealthHandler{}
	
	// Execute
	systemHealth := handler.getSystemHealth()
	
	// Assert
	assert.NotZero(t, systemHealth.Memory.Alloc)
	assert.NotZero(t, systemHealth.Runtime.NumCPU)
	assert.NotEmpty(t, systemHealth.Runtime.Version)
	assert.NotEmpty(t, systemHealth.Runtime.OS)
	assert.NotEmpty(t, systemHealth.Runtime.Arch)
	assert.Greater(t, systemHealth.Goroutines, 0)
}

func TestHealthHandler_CheckCacheHealth(t *testing.T) {
	// Setup
	mockImageCache := &MockImageCache{}
	mockImageCache.On("IsEnabled").Return(true)
	mockImageCache.On("Stats").Return(cache.ImageCacheStats{
		CacheStats: cache.CacheStats{
			HitRate: 85.5,
		},
	})
	
	handler := &HealthHandler{
		imageCache: mockImageCache,
	}
	
	// Execute
	cacheHealth := handler.checkCacheHealth()
	
	// Assert
	assert.Equal(t, "healthy", cacheHealth.Status)
	assert.Greater(t, cacheHealth.ResponseTime, time.Duration(0))
	assert.NotZero(t, cacheHealth.LastChecked)
	
	mockImageCache.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkHealthHandler_BasicHealthCheck(b *testing.B) {
	handler := &HealthHandler{}
	req := httptest.NewRequest("GET", "/api/health", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.BasicHealthCheck(w, req)
	}
}

func BenchmarkHealthHandler_MetricsEndpoint(b *testing.B) {
	mockImageCache := &MockImageCache{}
	mockImageCache.On("IsEnabled").Return(true)
	mockImageCache.On("Stats").Return(cache.ImageCacheStats{
		CacheStats: cache.CacheStats{
			Hits:   1000,
			Misses: 100,
		},
	})
	
	handler := &HealthHandler{
		imageCache: mockImageCache,
	}
	req := httptest.NewRequest("GET", "/api/metrics", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.MetricsEndpoint(w, req)
	}
}