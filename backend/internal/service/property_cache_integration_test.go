package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"realty-core/internal/cache"
	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// Test the integration of PropertyCache with PropertyService
func TestPropertyService_CacheIntegration(t *testing.T) {
	t.Run("GetProperty caches on first call", func(t *testing.T) {
		// Setup
		mockRepo := &MockPropertyRepository{}
		mockImageRepo := &MockImageRepository{}
		
		// Create service with cache enabled
		cacheConfig := cache.PropertyCacheConfig{
			Enabled:    true,
			Capacity:   100,
			DefaultTTL: 5 * time.Minute,
		}
		propertyCache := cache.NewPropertyCache(cacheConfig)
		service := NewPropertyServiceWithCache(mockRepo, mockImageRepo, propertyCache)

		testProperty := &domain.Property{
			ID:       "test-1",
			Title:    "Test Property",
			Province: "Pichincha",
			City:     "Quito",
			Type:     "house",
			Price:    100000,
			Status:   "available",
		}

		// Setup mocks - repo should be called only once
		mockRepo.On("GetByID", "test-1").Return(testProperty, nil).Once()
		mockRepo.On("Update", testProperty).Return(nil).Once()
		mockImageRepo.On("GetByPropertyID", "test-1").Return([]domain.ImageInfo{}, nil).Times(2)

		// First call should go to repo
		result1, err := service.GetProperty("test-1")
		assert.NoError(t, err)
		assert.Equal(t, "test-1", result1.ID)

		// Second call should come from cache (repo not called again)
		result2, err := service.GetProperty("test-1")
		assert.NoError(t, err)
		assert.Equal(t, "test-1", result2.ID)

		// Verify cache stats
		stats := service.GetCacheStats()
		assert.Equal(t, int64(1), stats.Hits)
		assert.Equal(t, int64(1), stats.Misses)

		mockRepo.AssertExpectations(t)
		mockImageRepo.AssertExpectations(t)
	})

	t.Run("SearchPropertiesRanked caches results", func(t *testing.T) {
		// Setup
		mockRepo := &MockPropertyRepository{}
		mockImageRepo := &MockImageRepository{}
		
		cacheConfig := cache.PropertyCacheConfig{
			Enabled:   true,
			Capacity:  100,
			SearchTTL: 1 * time.Minute,
		}
		propertyCache := cache.NewPropertyCache(cacheConfig)
		service := NewPropertyServiceWithCache(mockRepo, mockImageRepo, propertyCache)

		searchResults := []repository.PropertySearchResult{
			{
				Property: domain.Property{
					ID:    "search-1",
					Title: "Search Result 1",
					Price: 150000,
				},
				Rank: 0.95,
			},
		}

		// Setup mock - repo should be called only once
		mockRepo.On("SearchPropertiesRanked", "test query", 50).Return(searchResults, nil).Once()

		// First call should go to repo
		result1, err := service.SearchPropertiesRanked("test query", 50)
		assert.NoError(t, err)
		assert.Len(t, result1, 1)
		assert.Equal(t, "search-1", result1[0].Property.ID)

		// Second call should come from cache
		result2, err := service.SearchPropertiesRanked("test query", 50)
		assert.NoError(t, err)
		assert.Len(t, result2, 1)
		assert.Equal(t, "search-1", result2[0].Property.ID)

		// Verify cache stats
		stats := service.GetCacheStats()
		assert.Equal(t, int64(1), stats.SearchHits)
		assert.Equal(t, int64(1), stats.SearchMisses)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetStatistics caches results", func(t *testing.T) {
		// Setup
		mockRepo := &MockPropertyRepository{}
		mockImageRepo := &MockImageRepository{}
		
		cacheConfig := cache.PropertyCacheConfig{
			Enabled:       true,
			Capacity:      100,
			StatisticsTTL: 15 * time.Minute,
		}
		propertyCache := cache.NewPropertyCache(cacheConfig)
		service := NewPropertyServiceWithCache(mockRepo, mockImageRepo, propertyCache)

		testProperties := []domain.Property{
			{ID: "1", Type: "house", Status: "available", Province: "Pichincha", Price: 100000},
			{ID: "2", Type: "apartment", Status: "sold", Province: "Guayas", Price: 200000},
		}

		// Setup mock - repo should be called only once
		mockRepo.On("GetAll").Return(testProperties, nil).Once()

		// First call should go to repo
		result1, err := service.GetStatistics()
		assert.NoError(t, err)
		assert.Equal(t, 2, result1["total_properties"])

		// Second call should come from cache
		result2, err := service.GetStatistics()
		assert.NoError(t, err)
		assert.Equal(t, 2, result2["total_properties"])

		// Verify cache stats
		stats := service.GetCacheStats()
		assert.Equal(t, int64(1), stats.StatisticsHits)
		assert.Equal(t, int64(1), stats.StatisticsMisses)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Cache invalidation on property changes", func(t *testing.T) {
		// Setup
		mockRepo := &MockPropertyRepository{}
		mockImageRepo := &MockImageRepository{}
		
		cacheConfig := cache.PropertyCacheConfig{
			Enabled:    true,
			Capacity:   100,
			DefaultTTL: 5 * time.Minute,
		}
		propertyCache := cache.NewPropertyCache(cacheConfig)
		service := NewPropertyServiceWithCache(mockRepo, mockImageRepo, propertyCache)

		testProperty := &domain.Property{
			ID:       "test-invalidate",
			Title:    "Test Property",
			Province: "Pichincha",
			City:     "Quito",
			Type:     "house",
			Price:    100000,
			Status:   "available",
		}

		// Cache some statistics first
		testProperties := []domain.Property{*testProperty}
		mockRepo.On("GetAll").Return(testProperties, nil).Once()
		_, err := service.GetStatistics()
		assert.NoError(t, err)

		// Check that cache has statistics cached
		initialStats := service.GetCacheStats()
		assert.Greater(t, initialStats.Size, 0)

		// Update the property (should invalidate caches)
		mockRepo.On("GetByID", "test-invalidate").Return(testProperty, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*domain.Property")).Return(nil).Once()
		_, err = service.UpdateProperty("test-invalidate", "Updated Title", "Updated Description", "Pichincha", "Quito", "house", 150000)
		assert.NoError(t, err)

		// Verify that statistics cache was invalidated by checking it's empty
		// (note: individual property cache invalidation is harder to test due to immediate re-caching on get)
		clearedStats := service.GetCacheStats()
		assert.Equal(t, 0, clearedStats.Size) // Statistics cache should be cleared

		mockRepo.AssertExpectations(t)
	})

	t.Run("Cache with disabled configuration", func(t *testing.T) {
		// Setup with disabled cache
		mockRepo := &MockPropertyRepository{}
		mockImageRepo := &MockImageRepository{}
		
		cacheConfig := cache.PropertyCacheConfig{
			Enabled: false,
		}
		propertyCache := cache.NewPropertyCache(cacheConfig)
		service := NewPropertyServiceWithCache(mockRepo, mockImageRepo, propertyCache)

		testProperty := &domain.Property{
			ID:       "test-disabled",
			Title:    "Test Property",
			Province: "Pichincha",
			City:     "Quito",
			Type:     "house",
			Price:    100000,
			Status:   "available",
		}

		// Setup mocks - repo should be called every time since cache is disabled
		mockRepo.On("GetByID", "test-disabled").Return(testProperty, nil).Times(2)
		mockRepo.On("Update", testProperty).Return(nil).Times(2)
		mockImageRepo.On("GetByPropertyID", "test-disabled").Return([]domain.ImageInfo{}, nil).Times(2)

		// Both calls should go to repo
		_, err := service.GetProperty("test-disabled")
		assert.NoError(t, err)

		_, err = service.GetProperty("test-disabled")
		assert.NoError(t, err)

		// Cache stats should be empty (disabled)
		stats := service.GetCacheStats()
		assert.Equal(t, int64(0), stats.Hits)
		assert.Equal(t, int64(0), stats.Misses)

		mockRepo.AssertExpectations(t)
		mockImageRepo.AssertExpectations(t)
	})

	t.Run("Cache clear functionality", func(t *testing.T) {
		// Setup
		mockRepo := &MockPropertyRepository{}
		mockImageRepo := &MockImageRepository{}
		
		cacheConfig := cache.PropertyCacheConfig{
			Enabled:    true,
			Capacity:   100,
			DefaultTTL: 5 * time.Minute,
		}
		propertyCache := cache.NewPropertyCache(cacheConfig)
		service := NewPropertyServiceWithCache(mockRepo, mockImageRepo, propertyCache)

		// Add some data to cache
		testProperties := []domain.Property{
			{ID: "1", Type: "house", Status: "available", Province: "Pichincha", Price: 100000},
		}
		mockRepo.On("GetAll").Return(testProperties, nil).Once()

		// Cache some statistics
		_, err := service.GetStatistics()
		assert.NoError(t, err)

		// Verify cache has data
		stats := service.GetCacheStats()
		assert.Greater(t, stats.Size, 0)

		// Clear cache
		service.ClearCache()

		// Verify cache is empty
		clearedStats := service.GetCacheStats()
		assert.Equal(t, 0, clearedStats.Size)
		assert.Equal(t, int64(0), clearedStats.StatisticsHits)
		assert.Equal(t, int64(0), clearedStats.StatisticsMisses)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository error handling with cache", func(t *testing.T) {
		// Setup
		mockRepo := &MockPropertyRepository{}
		mockImageRepo := &MockImageRepository{}
		
		cacheConfig := cache.PropertyCacheConfig{
			Enabled:    true,
			Capacity:   100,
			DefaultTTL: 5 * time.Minute,
		}
		propertyCache := cache.NewPropertyCache(cacheConfig)
		service := NewPropertyServiceWithCache(mockRepo, mockImageRepo, propertyCache)

		// Setup mock to return error
		mockRepo.On("GetByID", "error-test").Return((*domain.Property)(nil), errors.New("database error"))

		// Call should return error and not cache anything
		_, err := service.GetProperty("error-test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Cache should remain empty
		stats := service.GetCacheStats()
		assert.Equal(t, int64(0), stats.Hits)
		assert.Equal(t, int64(1), stats.Misses)

		mockRepo.AssertExpectations(t)
	})
}