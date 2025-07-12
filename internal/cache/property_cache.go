package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
)

// PropertyCacheEntry represents a property cache entry
type PropertyCacheEntry struct {
	Data        interface{}
	Size        int64
	CreatedAt   time.Time
	AccessCount int64
	TTL         time.Duration
}

// PropertyCache wraps LRU cache with property-specific functionality
type PropertyCache struct {
	lru            CacheInterface
	mutex          sync.RWMutex
	enabled        bool
	defaultTTL     time.Duration
	searchTTL      time.Duration // Shorter TTL for search results
	statisticsTTL  time.Duration // Longer TTL for statistics
	stats          PropertyCacheStats
	logger         *log.Logger
}

// PropertyCacheStats represents property cache statistics
type PropertyCacheStats struct {
	CacheStats
	SearchHits      int64   `json:"search_hits"`
	SearchMisses    int64   `json:"search_misses"`
	SearchRate      float64 `json:"search_hit_rate"`
	StatisticsHits  int64   `json:"statistics_hits"`
	StatisticsMisses int64   `json:"statistics_misses"`
	StatisticsRate  float64 `json:"statistics_hit_rate"`
	ListHits        int64   `json:"list_hits"`
	ListMisses      int64   `json:"list_misses"`
	ListRate        float64 `json:"list_hit_rate"`
	FilterHits      int64   `json:"filter_hits"`
	FilterMisses    int64   `json:"filter_misses"`
	FilterRate      float64 `json:"filter_hit_rate"`
}

// PropertyCacheConfig defines configuration for property cache
type PropertyCacheConfig struct {
	Enabled        bool
	Capacity       int
	MaxSizeBytes   int64
	DefaultTTL     time.Duration
	SearchTTL      time.Duration
	StatisticsTTL  time.Duration
	Logger         *log.Logger
}

// NewPropertyCache creates a new property cache instance
func NewPropertyCache(config PropertyCacheConfig) *PropertyCache {
	if !config.Enabled {
		return &PropertyCache{enabled: false}
	}

	// Set defaults
	if config.Capacity <= 0 {
		config.Capacity = 1000
	}
	if config.MaxSizeBytes <= 0 {
		config.MaxSizeBytes = 50 * 1024 * 1024 // 50MB
	}
	if config.DefaultTTL <= 0 {
		config.DefaultTTL = 5 * time.Minute
	}
	if config.SearchTTL <= 0 {
		config.SearchTTL = 1 * time.Minute // Búsquedas cambian rápido
	}
	if config.StatisticsTTL <= 0 {
		config.StatisticsTTL = 15 * time.Minute // Estadísticas cambian lento
	}

	lru := NewLRUCache(config.Capacity, config.MaxSizeBytes, config.DefaultTTL)

	return &PropertyCache{
		lru:           lru,
		enabled:       true,
		defaultTTL:    config.DefaultTTL,
		searchTTL:     config.SearchTTL,
		statisticsTTL: config.StatisticsTTL,
		logger:        config.Logger,
	}
}

// GetProperty retrieves a cached property by ID
func (pc *PropertyCache) GetProperty(id string) (*domain.Property, bool) {
	if !pc.enabled {
		return nil, false
	}

	key := fmt.Sprintf("property:%s", id)
	if value, found := pc.lru.Get(key); found {
		if property, ok := value.(*domain.Property); ok {
			pc.incrementHits()
			if pc.logger != nil {
				pc.logger.Printf("Property cache HIT: %s", id)
			}
			return property, true
		}
	}

	pc.incrementMisses()
	if pc.logger != nil {
		pc.logger.Printf("Property cache MISS: %s", id)
	}
	return nil, false
}

// SetProperty stores a property in cache
func (pc *PropertyCache) SetProperty(property *domain.Property) {
	if !pc.enabled || property == nil {
		return
	}

	key := fmt.Sprintf("property:%s", property.ID)
	size := pc.estimatePropertySize(property)
	
	pc.lru.Set(key, property, size)
	
	if pc.logger != nil {
		pc.logger.Printf("Property cached: %s (size: %d bytes)", property.ID, size)
	}
}

// GetSearchResults retrieves cached search results
func (pc *PropertyCache) GetSearchResults(query string, limit int) ([]repository.PropertySearchResult, bool) {
	if !pc.enabled {
		return nil, false
	}

	key := fmt.Sprintf("search:%s:limit:%d", query, limit)
	if value, found := pc.lru.Get(key); found {
		if results, ok := value.([]repository.PropertySearchResult); ok {
			pc.mutex.Lock()
			pc.stats.SearchHits++
			pc.mutex.Unlock()
			
			if pc.logger != nil {
				pc.logger.Printf("Search cache HIT: %s", query)
			}
			return results, true
		}
	}

	pc.mutex.Lock()
	pc.stats.SearchMisses++
	pc.mutex.Unlock()
	
	if pc.logger != nil {
		pc.logger.Printf("Search cache MISS: %s", query)
	}
	return nil, false
}

// SetSearchResults stores search results in cache
func (pc *PropertyCache) SetSearchResults(query string, limit int, results []repository.PropertySearchResult) {
	if !pc.enabled || len(results) == 0 {
		return
	}

	key := fmt.Sprintf("search:%s:limit:%d", query, limit)
	size := pc.estimateSearchResultsSize(results)
	
	// Use shorter TTL for search results
	pc.lru.SetWithTTL(key, results, size, pc.searchTTL)
	
	if pc.logger != nil {
		pc.logger.Printf("Search results cached: %s (count: %d, size: %d bytes)", 
			query, len(results), size)
	}
}

// GetFilterResults retrieves cached filter results
func (pc *PropertyCache) GetFilterResults(province string, minPrice, maxPrice float64) ([]domain.Property, bool) {
	if !pc.enabled {
		return nil, false
	}

	key := fmt.Sprintf("filter:province:%s:price:%.0f-%.0f", province, minPrice, maxPrice)
	if value, found := pc.lru.Get(key); found {
		if properties, ok := value.([]domain.Property); ok {
			pc.mutex.Lock()
			pc.stats.FilterHits++
			pc.mutex.Unlock()
			
			if pc.logger != nil {
				pc.logger.Printf("Filter cache HIT: %s", key)
			}
			return properties, true
		}
	}

	pc.mutex.Lock()
	pc.stats.FilterMisses++
	pc.mutex.Unlock()
	
	if pc.logger != nil {
		pc.logger.Printf("Filter cache MISS: %s", key)
	}
	return nil, false
}

// SetFilterResults stores filter results in cache
func (pc *PropertyCache) SetFilterResults(province string, minPrice, maxPrice float64, properties []domain.Property) {
	if !pc.enabled || len(properties) == 0 {
		return
	}

	key := fmt.Sprintf("filter:province:%s:price:%.0f-%.0f", province, minPrice, maxPrice)
	size := pc.estimatePropertiesSize(properties)
	
	pc.lru.Set(key, properties, size)
	
	if pc.logger != nil {
		pc.logger.Printf("Filter results cached: %s (count: %d, size: %d bytes)", 
			key, len(properties), size)
	}
}

// GetStatistics retrieves cached statistics
func (pc *PropertyCache) GetStatistics(key string) (map[string]interface{}, bool) {
	if !pc.enabled {
		return nil, false
	}

	cacheKey := fmt.Sprintf("stats:%s", key)
	if value, found := pc.lru.Get(cacheKey); found {
		if stats, ok := value.(map[string]interface{}); ok {
			pc.mutex.Lock()
			pc.stats.StatisticsHits++
			pc.mutex.Unlock()
			
			if pc.logger != nil {
				pc.logger.Printf("Statistics cache HIT: %s", key)
			}
			return stats, true
		}
	}

	pc.mutex.Lock()
	pc.stats.StatisticsMisses++
	pc.mutex.Unlock()
	
	if pc.logger != nil {
		pc.logger.Printf("Statistics cache MISS: %s", key)
	}
	return nil, false
}

// SetStatistics stores statistics in cache
func (pc *PropertyCache) SetStatistics(key string, stats map[string]interface{}) {
	if !pc.enabled || len(stats) == 0 {
		return
	}

	cacheKey := fmt.Sprintf("stats:%s", key)
	size := pc.estimateStatsSize(stats)
	
	// Use longer TTL for statistics
	pc.lru.SetWithTTL(cacheKey, stats, size, pc.statisticsTTL)
	
	if pc.logger != nil {
		pc.logger.Printf("Statistics cached: %s (size: %d bytes)", key, size)
	}
}

// InvalidateProperty removes a property from cache
func (pc *PropertyCache) InvalidateProperty(id string) {
	if !pc.enabled {
		return
	}

	key := fmt.Sprintf("property:%s", id)
	pc.lru.Delete(key)
	
	// Also invalidate related caches that might contain this property
	pc.InvalidateSearchResults() // Clear search cache when properties change
	pc.InvalidateStatistics()     // Clear stats cache when properties change
	
	if pc.logger != nil {
		pc.logger.Printf("Property cache invalidated: %s", id)
	}
}

// InvalidateSearchResults clears all search result caches
func (pc *PropertyCache) InvalidateSearchResults() {
	if !pc.enabled {
		return
	}

	// Get all keys and remove search-related ones
	keys := pc.lru.Keys()
	for _, key := range keys {
		if strings.HasPrefix(key, "search:") || strings.HasPrefix(key, "filter:") {
			pc.lru.Delete(key)
		}
	}
	
	if pc.logger != nil {
		pc.logger.Printf("Search results cache invalidated")
	}
}

// InvalidateStatistics clears statistics cache
func (pc *PropertyCache) InvalidateStatistics() {
	if !pc.enabled {
		return
	}

	// Get all keys and remove stats-related ones
	keys := pc.lru.Keys()
	for _, key := range keys {
		if strings.HasPrefix(key, "stats:") {
			pc.lru.Delete(key)
		}
	}
	
	if pc.logger != nil {
		pc.logger.Printf("Statistics cache invalidated")
	}
}

// GetStats returns comprehensive cache statistics
func (pc *PropertyCache) GetStats() PropertyCacheStats {
	if !pc.enabled {
		return PropertyCacheStats{}
	}

	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	// Get base stats from LRU cache
	baseStats := pc.lru.Stats()
	
	// Calculate hit rates
	searchTotal := pc.stats.SearchHits + pc.stats.SearchMisses
	statsTotal := pc.stats.StatisticsHits + pc.stats.StatisticsMisses
	listTotal := pc.stats.ListHits + pc.stats.ListMisses
	filterTotal := pc.stats.FilterHits + pc.stats.FilterMisses

	stats := pc.stats
	stats.CacheStats = baseStats
	
	if searchTotal > 0 {
		stats.SearchRate = float64(stats.SearchHits) / float64(searchTotal) * 100
	}
	if statsTotal > 0 {
		stats.StatisticsRate = float64(stats.StatisticsHits) / float64(statsTotal) * 100
	}
	if listTotal > 0 {
		stats.ListRate = float64(stats.ListHits) / float64(listTotal) * 100
	}
	if filterTotal > 0 {
		stats.FilterRate = float64(stats.FilterHits) / float64(filterTotal) * 100
	}

	return stats
}

// Clear removes all entries from cache
func (pc *PropertyCache) Clear() {
	if !pc.enabled {
		return
	}

	pc.lru.Clear()
	
	// Reset stats
	pc.mutex.Lock()
	pc.stats = PropertyCacheStats{}
	pc.mutex.Unlock()
	
	if pc.logger != nil {
		pc.logger.Printf("Property cache cleared")
	}
}

// IsEnabled returns whether cache is enabled
func (pc *PropertyCache) IsEnabled() bool {
	return pc.enabled
}

// Helper methods for size estimation

func (pc *PropertyCache) estimatePropertySize(property *domain.Property) int64 {
	// Rough estimation of property size in memory
	baseSize := int64(500) // Base struct size
	
	baseSize += int64(len(property.Title))
	baseSize += int64(len(property.Description))
	baseSize += int64(len(property.Province))
	baseSize += int64(len(property.City))
	
	// Estimate JSON fields
	if property.Tags != nil {
		if data, err := json.Marshal(property.Tags); err == nil {
			baseSize += int64(len(data))
		}
	}
	if property.Images != nil {
		if data, err := json.Marshal(property.Images); err == nil {
			baseSize += int64(len(data))
		}
	}
	
	return baseSize
}

func (pc *PropertyCache) estimatePropertiesSize(properties []domain.Property) int64 {
	total := int64(0)
	for _, property := range properties {
		total += pc.estimatePropertySize(&property)
	}
	return total
}

func (pc *PropertyCache) estimateSearchResultsSize(results []repository.PropertySearchResult) int64 {
	total := int64(0)
	for _, result := range results {
		total += pc.estimatePropertySize(&result.Property)
		total += 8 // For the Rank field (float64)
	}
	return total
}

func (pc *PropertyCache) estimateStatsSize(stats map[string]interface{}) int64 {
	if data, err := json.Marshal(stats); err == nil {
		return int64(len(data))
	}
	return 1024 // Default estimate
}

func (pc *PropertyCache) incrementHits() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.stats.Hits++
}

func (pc *PropertyCache) incrementMisses() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.stats.Misses++
}