package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewImageCache(t *testing.T) {
	tests := []struct {
		name   string
		config ImageCacheConfig
		want   bool // enabled
	}{
		{
			name: "enabled cache with default config",
			config: ImageCacheConfig{
				Enabled: true,
			},
			want: true,
		},
		{
			name: "disabled cache",
			config: ImageCacheConfig{
				Enabled: false,
			},
			want: false,
		},
		{
			name: "enabled cache with custom config",
			config: ImageCacheConfig{
				Enabled:         true,
				Capacity:        500,
				MaxSizeBytes:    50 * 1024 * 1024,
				TTL:             30 * time.Minute,
				CleanupInterval: 5 * time.Minute,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewImageCache(tt.config)
			
			assert.NotNil(t, cache)
			assert.Equal(t, tt.want, cache.IsEnabled())
			
			if tt.want {
				assert.NotNil(t, cache.lru)
			}
		})
	}
}

func TestImageCache_SetAndGet(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	// Test setting and getting data
	data := []byte("test image data")
	contentType := "image/jpeg"
	
	cache.Set("test-key", data, contentType)
	
	retrievedData, retrievedContentType, found := cache.Get("test-key")
	assert.True(t, found)
	assert.Equal(t, data, retrievedData)
	assert.Equal(t, contentType, retrievedContentType)
	
	// Test getting non-existent key
	_, _, found = cache.Get("non-existent")
	assert.False(t, found)
}

func TestImageCache_SetAndGet_Disabled(t *testing.T) {
	config := ImageCacheConfig{
		Enabled: false,
	}
	cache := NewImageCache(config)

	// Test with disabled cache
	data := []byte("test image data")
	contentType := "image/jpeg"
	
	cache.Set("test-key", data, contentType)
	
	retrievedData, retrievedContentType, found := cache.Get("test-key")
	assert.False(t, found)
	assert.Nil(t, retrievedData)
	assert.Empty(t, retrievedContentType)
}

func TestImageCache_ThumbnailOperations(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	imageID := "test-image-123"
	size := 150
	data := []byte("thumbnail data")
	contentType := "image/jpeg"
	
	// Test setting thumbnail
	cache.SetThumbnail(imageID, size, data, contentType)
	
	// Test getting thumbnail
	retrievedData, retrievedContentType, found := cache.GetThumbnail(imageID, size)
	assert.True(t, found)
	assert.Equal(t, data, retrievedData)
	assert.Equal(t, contentType, retrievedContentType)
	
	// Test getting non-existent thumbnail
	_, _, found = cache.GetThumbnail("non-existent", size)
	assert.False(t, found)
}

func TestImageCache_VariantOperations(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	imageID := "test-image-123"
	width := 800
	height := 600
	quality := 85
	format := "jpg"
	data := []byte("variant data")
	contentType := "image/jpeg"
	
	// Test setting variant
	cache.SetVariant(imageID, width, height, quality, format, data, contentType)
	
	// Test getting variant
	retrievedData, retrievedContentType, found := cache.GetVariant(imageID, width, height, quality, format)
	assert.True(t, found)
	assert.Equal(t, data, retrievedData)
	assert.Equal(t, contentType, retrievedContentType)
	
	// Test getting non-existent variant
	_, _, found = cache.GetVariant("non-existent", width, height, quality, format)
	assert.False(t, found)
}

func TestImageCache_Delete(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	data := []byte("test image data")
	contentType := "image/jpeg"
	
	// Set data
	cache.Set("test-key", data, contentType)
	
	// Verify it exists
	_, _, found := cache.Get("test-key")
	assert.True(t, found)
	
	// Delete it
	deleted := cache.Delete("test-key")
	assert.True(t, deleted)
	
	// Verify it's gone
	_, _, found = cache.Get("test-key")
	assert.False(t, found)
	
	// Delete non-existent key
	deleted = cache.Delete("non-existent")
	assert.False(t, deleted)
}

func TestImageCache_Clear(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	// Add some data
	cache.Set("key1", []byte("data1"), "image/jpeg")
	cache.Set("key2", []byte("data2"), "image/png")
	
	assert.Equal(t, 2, cache.Size())
	
	// Clear cache
	cache.Clear()
	
	assert.Equal(t, 0, cache.Size())
	
	// Verify data is gone
	_, _, found1 := cache.Get("key1")
	assert.False(t, found1)
	
	_, _, found2 := cache.Get("key2")
	assert.False(t, found2)
}

func TestImageCache_InvalidateImage(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	imageID := "test-image-123"
	
	// Add thumbnail and variant for the image
	cache.SetThumbnail(imageID, 150, []byte("thumb"), "image/jpeg")
	cache.SetVariant(imageID, 800, 600, 85, "jpg", []byte("variant"), "image/jpeg")
	
	// Add data for another image
	cache.SetThumbnail("other-image", 150, []byte("other"), "image/jpeg")
	
	assert.Equal(t, 3, cache.Size())
	
	// Invalidate the first image
	removed := cache.InvalidateImage(imageID)
	assert.Equal(t, 2, removed)
	assert.Equal(t, 1, cache.Size())
	
	// Verify first image data is gone
	_, _, found1 := cache.GetThumbnail(imageID, 150)
	assert.False(t, found1)
	
	_, _, found2 := cache.GetVariant(imageID, 800, 600, 85, "jpg")
	assert.False(t, found2)
	
	// Verify other image data still exists
	_, _, found3 := cache.GetThumbnail("other-image", 150)
	assert.True(t, found3)
}

func TestImageCache_Stats(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	// Initial stats
	stats := cache.Stats()
	assert.Equal(t, 0, stats.Size)
	assert.Equal(t, int64(0), stats.CurrentSize)
	assert.Equal(t, int64(0), stats.ThumbnailHits)
	assert.Equal(t, int64(0), stats.VariantHits)
	assert.Equal(t, int64(0), stats.ThumbnailMisses)
	assert.Equal(t, int64(0), stats.VariantMisses)
	assert.Equal(t, float64(0), stats.ThumbnailRate)
	assert.Equal(t, float64(0), stats.VariantRate)
	
	// Add some data and test stats
	cache.SetThumbnail("image1", 150, []byte("thumb1"), "image/jpeg")
	cache.SetVariant("image1", 800, 600, 85, "jpg", []byte("variant1"), "image/jpeg")
	
	// Test hits and misses
	cache.GetThumbnail("image1", 150)     // Hit
	cache.GetThumbnail("image1", 200)     // Miss
	cache.GetVariant("image1", 800, 600, 85, "jpg") // Hit
	cache.GetVariant("image1", 400, 300, 80, "jpg") // Miss
	
	stats = cache.Stats()
	assert.Equal(t, 2, stats.Size)
	assert.Greater(t, stats.CurrentSize, int64(0))
	assert.Equal(t, int64(2), stats.ThumbnailHits)
	assert.Equal(t, int64(2), stats.VariantHits)
	assert.Equal(t, int64(1), stats.ThumbnailMisses)
	assert.Equal(t, int64(1), stats.VariantMisses)
	assert.Equal(t, float64(66.66666666666666), stats.ThumbnailRate)
	assert.Equal(t, float64(66.66666666666666), stats.VariantRate)
}

func TestImageCache_GenerateKeys(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	// Test thumbnail key generation
	thumbnailKey := cache.generateThumbnailKey("image123", 150)
	assert.Equal(t, "image123_thumbnail_150", thumbnailKey)
	
	// Test variant key generation
	variantKey := cache.generateVariantKey("image123", 800, 600, 85, "jpg")
	assert.Equal(t, "image123_variant_800x600_q85_jpg", variantKey)
}

func TestImageCache_DisabledOperations(t *testing.T) {
	config := ImageCacheConfig{
		Enabled: false,
	}
	cache := NewImageCache(config)

	// All operations should be no-ops or return defaults
	assert.False(t, cache.IsEnabled())
	assert.Equal(t, 0, cache.Size())
	assert.Equal(t, int64(0), cache.CurrentSize())
	
	// Set operations should not crash
	cache.Set("key", []byte("data"), "image/jpeg")
	cache.SetThumbnail("image", 150, []byte("thumb"), "image/jpeg")
	cache.SetVariant("image", 800, 600, 85, "jpg", []byte("variant"), "image/jpeg")
	
	// Get operations should return false
	_, _, found1 := cache.Get("key")
	assert.False(t, found1)
	
	_, _, found2 := cache.GetThumbnail("image", 150)
	assert.False(t, found2)
	
	_, _, found3 := cache.GetVariant("image", 800, 600, 85, "jpg")
	assert.False(t, found3)
	
	// Delete operations should return false
	assert.False(t, cache.Delete("key"))
	
	// Invalidate should return 0
	assert.Equal(t, 0, cache.InvalidateImage("image"))
	
	// Cleanup should return 0
	assert.Equal(t, 0, cache.CleanupExpired())
	
	// Stats should return empty
	stats := cache.Stats()
	assert.Equal(t, ImageCacheStats{}, stats)
}

func TestImageCache_GetPopularImages(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	// Add some images
	cache.SetThumbnail("image1", 150, []byte("thumb1"), "image/jpeg")
	cache.SetThumbnail("image2", 150, []byte("thumb2"), "image/jpeg")
	cache.SetThumbnail("image3", 150, []byte("thumb3"), "image/jpeg")
	
	// Access them to increase access count
	cache.GetThumbnail("image1", 150)
	cache.GetThumbnail("image1", 150)
	cache.GetThumbnail("image2", 150)
	
	// Get popular images
	popular := cache.GetPopularImages(5)
	assert.Len(t, popular, 3)
	
	// Verify structure
	for _, img := range popular {
		assert.NotEmpty(t, img.Key)
		assert.Greater(t, img.AccessCount, int64(0))
		assert.Greater(t, img.Size, int64(0))
		assert.False(t, img.CreatedAt.IsZero())
	}
}

func TestImageCache_DefaultConfig(t *testing.T) {
	config := DefaultImageCacheConfig()
	
	assert.True(t, config.Enabled)
	assert.Equal(t, 1000, config.Capacity)
	assert.Equal(t, int64(100*1024*1024), config.MaxSizeBytes)
	assert.Equal(t, 1*time.Hour, config.TTL)
	assert.Equal(t, 10*time.Minute, config.CleanupInterval)
}

func TestNewDisabledImageCache(t *testing.T) {
	cache := NewDisabledImageCache()
	
	assert.NotNil(t, cache)
	assert.False(t, cache.IsEnabled())
}

func TestImageCache_EmptyDataHandling(t *testing.T) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     10,
		MaxSizeBytes: 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	// Test setting empty data (should not be stored)
	cache.Set("empty-key", []byte{}, "image/jpeg")
	
	_, _, found := cache.Get("empty-key")
	assert.False(t, found)
	assert.Equal(t, 0, cache.Size())
}

func BenchmarkImageCache_Set(b *testing.B) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     1000,
		MaxSizeBytes: 10 * 1024 * 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	data := make([]byte, 1024) // 1KB image data
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("image-%d", i%1000)
		cache.Set(key, data, "image/jpeg")
	}
}

func BenchmarkImageCache_Get(b *testing.B) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     1000,
		MaxSizeBytes: 10 * 1024 * 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	data := make([]byte, 1024) // 1KB image data
	
	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("image-%d", i)
		cache.Set(key, data, "image/jpeg")
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("image-%d", i%1000)
		cache.Get(key)
	}
}

func BenchmarkImageCache_Thumbnail(b *testing.B) {
	config := ImageCacheConfig{
		Enabled:      true,
		Capacity:     1000,
		MaxSizeBytes: 10 * 1024 * 1024,
		TTL:          1 * time.Hour,
	}
	cache := NewImageCache(config)

	data := make([]byte, 512) // 512B thumbnail data
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		imageID := fmt.Sprintf("image-%d", i%1000)
		if i%2 == 0 {
			cache.SetThumbnail(imageID, 150, data, "image/jpeg")
		} else {
			cache.GetThumbnail(imageID, 150)
		}
	}
}