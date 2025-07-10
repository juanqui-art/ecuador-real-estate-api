package cache

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// ImageCacheEntry represents an image cache entry with metadata
type ImageCacheEntry struct {
	Data        []byte
	ContentType string
	Size        int64
	CreatedAt   time.Time
	AccessCount int64
}

// ImageCache wraps LRU cache with image-specific functionality
type ImageCache struct {
	lru        CacheInterface
	mutex      sync.RWMutex
	enabled    bool
	defaultTTL time.Duration
	stats      ImageCacheStats
}

// ImageCacheStats represents image cache statistics
type ImageCacheStats struct {
	CacheStats
	TotalDataSize    int64   `json:"total_data_size"`
	AverageEntrySize int64   `json:"average_entry_size"`
	ThumbnailHits    int64   `json:"thumbnail_hits"`
	VariantHits      int64   `json:"variant_hits"`
	ThumbnailMisses  int64   `json:"thumbnail_misses"`
	VariantMisses    int64   `json:"variant_misses"`
	ThumbnailRate    float64 `json:"thumbnail_hit_rate"`
	VariantRate      float64 `json:"variant_hit_rate"`
}

// ImageCacheConfig represents configuration for image cache
type ImageCacheConfig struct {
	Enabled         bool
	Capacity        int
	MaxSizeBytes    int64
	TTL             time.Duration
	CleanupInterval time.Duration
}

// NewImageCache creates a new image cache
func NewImageCache(config ImageCacheConfig) *ImageCache {
	if !config.Enabled {
		return &ImageCache{
			enabled: false,
		}
	}
	
	// Set defaults
	if config.Capacity <= 0 {
		config.Capacity = 1000
	}
	if config.MaxSizeBytes <= 0 {
		config.MaxSizeBytes = 100 * 1024 * 1024 // 100MB
	}
	if config.TTL <= 0 {
		config.TTL = 1 * time.Hour
	}
	if config.CleanupInterval <= 0 {
		config.CleanupInterval = 10 * time.Minute
	}
	
	cache := &ImageCache{
		lru:        NewLRUCache(config.Capacity, config.MaxSizeBytes, config.TTL),
		enabled:    true,
		defaultTTL: config.TTL,
		stats:      ImageCacheStats{},
	}
	
	// Start cleanup goroutine
	go cache.cleanupRoutine(config.CleanupInterval)
	
	return cache
}

// Get retrieves an image from cache
func (c *ImageCache) Get(key string) ([]byte, string, bool) {
	if !c.enabled {
		return nil, "", false
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Update type-specific stats
	if strings.Contains(key, "thumbnail") {
		c.stats.ThumbnailHits++
	} else {
		c.stats.VariantHits++
	}
	
	value, found := c.lru.Get(key)
	if !found {
		// Update miss stats
		if strings.Contains(key, "thumbnail") {
			c.stats.ThumbnailMisses++
		} else {
			c.stats.VariantMisses++
		}
		return nil, "", false
	}
	
	entry, ok := value.(*ImageCacheEntry)
	if !ok {
		log.Printf("Warning: Invalid cache entry type for key: %s", key)
		c.lru.Delete(key)
		return nil, "", false
	}
	
	// Update access count
	entry.AccessCount++
	
	return entry.Data, entry.ContentType, true
}

// Set stores an image in cache
func (c *ImageCache) Set(key string, data []byte, contentType string) {
	if !c.enabled || len(data) == 0 {
		return
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	entry := &ImageCacheEntry{
		Data:        data,
		ContentType: contentType,
		Size:        int64(len(data)),
		CreatedAt:   time.Now(),
		AccessCount: 1,
	}
	
	c.lru.Set(key, entry, entry.Size)
	
	// Update stats
	c.stats.TotalDataSize += entry.Size
}

// Delete removes an image from cache
func (c *ImageCache) Delete(key string) bool {
	if !c.enabled {
		return false
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Get entry size before deletion for stats
	if value, found := c.lru.Get(key); found {
		if entry, ok := value.(*ImageCacheEntry); ok {
			c.stats.TotalDataSize -= entry.Size
		}
	}
	
	return c.lru.Delete(key)
}

// Clear removes all entries from cache
func (c *ImageCache) Clear() {
	if !c.enabled {
		return
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.lru.Clear()
	c.stats = ImageCacheStats{}
}

// GetThumbnail retrieves a thumbnail from cache
func (c *ImageCache) GetThumbnail(imageID string, size int) ([]byte, string, bool) {
	key := c.generateThumbnailKey(imageID, size)
	return c.Get(key)
}

// SetThumbnail stores a thumbnail in cache
func (c *ImageCache) SetThumbnail(imageID string, size int, data []byte, contentType string) {
	key := c.generateThumbnailKey(imageID, size)
	c.Set(key, data, contentType)
}

// GetVariant retrieves an image variant from cache
func (c *ImageCache) GetVariant(imageID string, width, height, quality int, format string) ([]byte, string, bool) {
	key := c.generateVariantKey(imageID, width, height, quality, format)
	return c.Get(key)
}

// SetVariant stores an image variant in cache
func (c *ImageCache) SetVariant(imageID string, width, height, quality int, format string, data []byte, contentType string) {
	key := c.generateVariantKey(imageID, width, height, quality, format)
	c.Set(key, data, contentType)
}

// InvalidateImage removes all cached variants for an image
func (c *ImageCache) InvalidateImage(imageID string) int {
	if !c.enabled {
		return 0
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	keys := c.lru.Keys()
	removed := 0
	
	for _, key := range keys {
		if strings.HasPrefix(key, imageID+"_") {
			c.lru.Delete(key)
			removed++
		}
	}
	
	return removed
}

// Stats returns cache statistics
func (c *ImageCache) Stats() ImageCacheStats {
	if !c.enabled {
		return ImageCacheStats{}
	}
	
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	baseStats := c.lru.Stats()
	stats := c.stats
	stats.CacheStats = baseStats
	
	// Calculate averages
	if stats.Size > 0 {
		stats.AverageEntrySize = stats.CurrentSize / int64(stats.Size)
	}
	
	// Calculate hit rates
	thumbnailTotal := stats.ThumbnailHits + stats.ThumbnailMisses
	if thumbnailTotal > 0 {
		stats.ThumbnailRate = float64(stats.ThumbnailHits) / float64(thumbnailTotal) * 100
	}
	
	variantTotal := stats.VariantHits + stats.VariantMisses
	if variantTotal > 0 {
		stats.VariantRate = float64(stats.VariantHits) / float64(variantTotal) * 100
	}
	
	return stats
}

// IsEnabled returns whether the cache is enabled
func (c *ImageCache) IsEnabled() bool {
	return c.enabled
}

// Size returns the number of cached items
func (c *ImageCache) Size() int {
	if !c.enabled {
		return 0
	}
	return c.lru.Size()
}

// CurrentSize returns the current cache size in bytes
func (c *ImageCache) CurrentSize() int64 {
	if !c.enabled {
		return 0
	}
	return c.lru.CurrentSize()
}

// CleanupExpired removes expired entries
func (c *ImageCache) CleanupExpired() int {
	if !c.enabled {
		return 0
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	return c.lru.CleanupExpired()
}

// GetPopularImages returns most accessed images
func (c *ImageCache) GetPopularImages(limit int) []PopularImage {
	if !c.enabled {
		return nil
	}
	
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	keys := c.lru.Keys()
	var popular []PopularImage
	
	for _, key := range keys {
		if value, found := c.lru.Get(key); found {
			if entry, ok := value.(*ImageCacheEntry); ok {
				popular = append(popular, PopularImage{
					Key:         key,
					AccessCount: entry.AccessCount,
					Size:        entry.Size,
					CreatedAt:   entry.CreatedAt,
				})
			}
		}
	}
	
	// Sort by access count (simplified - in production might use heap)
	// For now, just return the first 'limit' items
	if len(popular) > limit {
		popular = popular[:limit]
	}
	
	return popular
}

// generateThumbnailKey creates a cache key for thumbnails
func (c *ImageCache) generateThumbnailKey(imageID string, size int) string {
	return fmt.Sprintf("%s_thumbnail_%d", imageID, size)
}

// generateVariantKey creates a cache key for image variants
func (c *ImageCache) generateVariantKey(imageID string, width, height, quality int, format string) string {
	return fmt.Sprintf("%s_variant_%dx%d_q%d_%s", imageID, width, height, quality, format)
}

// cleanupRoutine runs periodic cleanup
func (c *ImageCache) cleanupRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for range ticker.C {
		if !c.enabled {
			return
		}
		
		removed := c.CleanupExpired()
		if removed > 0 {
			log.Printf("Image cache cleanup: removed %d expired entries", removed)
		}
	}
}

// PopularImage represents a popular cached image
type PopularImage struct {
	Key         string    `json:"key"`
	AccessCount int64     `json:"access_count"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}

// ImageCacheInterface defines the interface for image cache operations
type ImageCacheInterface interface {
	Get(key string) ([]byte, string, bool)
	Set(key string, data []byte, contentType string)
	Delete(key string) bool
	Clear()
	GetThumbnail(imageID string, size int) ([]byte, string, bool)
	SetThumbnail(imageID string, size int, data []byte, contentType string)
	GetVariant(imageID string, width, height, quality int, format string) ([]byte, string, bool)
	SetVariant(imageID string, width, height, quality int, format string, data []byte, contentType string)
	InvalidateImage(imageID string) int
	Stats() ImageCacheStats
	IsEnabled() bool
	Size() int
	CurrentSize() int64
	CleanupExpired() int
	GetPopularImages(limit int) []PopularImage
}

// DefaultImageCacheConfig returns default configuration for image cache
func DefaultImageCacheConfig() ImageCacheConfig {
	return ImageCacheConfig{
		Enabled:         true,
		Capacity:        1000,
		MaxSizeBytes:    100 * 1024 * 1024, // 100MB
		TTL:             1 * time.Hour,
		CleanupInterval: 10 * time.Minute,
	}
}

// NewDisabledImageCache returns a disabled image cache
func NewDisabledImageCache() *ImageCache {
	return &ImageCache{
		enabled: false,
	}
}