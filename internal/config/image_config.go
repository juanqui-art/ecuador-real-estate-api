package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"realty-core/internal/cache"
)

// ImageConfig represents configuration for the image system
type ImageConfig struct {
	// Storage configuration
	StorageBasePath string
	StorageBaseURL  string
	StorageMaxSize  int64
	
	// Processing configuration
	ProcessorMaxWidth  int
	ProcessorMaxHeight int
	DefaultQuality     int
	
	// Cache configuration
	CacheEnabled         bool
	CacheCapacity        int
	CacheMaxSizeBytes    int64
	CacheTTL             time.Duration
	CacheCleanupInterval time.Duration
	
	// Upload limits
	MaxUploadSize        int64
	MaxImagesPerProperty int
	
	// Feature flags
	EnableImageOptimization bool
	EnableVariantGeneration bool
	EnableThumbnailCache    bool
}

// LoadImageConfig loads image configuration from environment variables
func LoadImageConfig() *ImageConfig {
	config := &ImageConfig{
		// Storage defaults
		StorageBasePath: getEnv("IMAGE_STORAGE_BASE_PATH", "./storage/images"),
		StorageBaseURL:  getEnv("IMAGE_STORAGE_BASE_URL", "/images"),
		StorageMaxSize:  getEnvInt64("IMAGE_STORAGE_MAX_SIZE", 1024*1024*1024), // 1GB
		
		// Processing defaults
		ProcessorMaxWidth:  getEnvInt("IMAGE_PROCESSOR_MAX_WIDTH", 3000),
		ProcessorMaxHeight: getEnvInt("IMAGE_PROCESSOR_MAX_HEIGHT", 2000),
		DefaultQuality:     getEnvInt("IMAGE_DEFAULT_QUALITY", 85),
		
		// Cache defaults
		CacheEnabled:         getEnvBool("IMAGE_CACHE_ENABLED", true),
		CacheCapacity:        getEnvInt("IMAGE_CACHE_CAPACITY", 1000),
		CacheMaxSizeBytes:    getEnvInt64("IMAGE_CACHE_MAX_SIZE", 100*1024*1024), // 100MB
		CacheTTL:             getEnvDuration("IMAGE_CACHE_TTL", 1*time.Hour),
		CacheCleanupInterval: getEnvDuration("IMAGE_CACHE_CLEANUP_INTERVAL", 10*time.Minute),
		
		// Upload limits
		MaxUploadSize:        getEnvInt64("IMAGE_MAX_UPLOAD_SIZE", 10*1024*1024), // 10MB
		MaxImagesPerProperty: getEnvInt("IMAGE_MAX_IMAGES_PER_PROPERTY", 50),
		
		// Feature flags
		EnableImageOptimization: getEnvBool("IMAGE_ENABLE_OPTIMIZATION", true),
		EnableVariantGeneration: getEnvBool("IMAGE_ENABLE_VARIANT_GENERATION", true),
		EnableThumbnailCache:    getEnvBool("IMAGE_ENABLE_THUMBNAIL_CACHE", true),
	}
	
	return config
}

// ToCacheConfig converts ImageConfig to cache.ImageCacheConfig
func (c *ImageConfig) ToCacheConfig() cache.ImageCacheConfig {
	return cache.ImageCacheConfig{
		Enabled:         c.CacheEnabled,
		Capacity:        c.CacheCapacity,
		MaxSizeBytes:    c.CacheMaxSizeBytes,
		TTL:             c.CacheTTL,
		CleanupInterval: c.CacheCleanupInterval,
	}
}

// Validate validates the configuration
func (c *ImageConfig) Validate() error {
	if c.StorageBasePath == "" {
		return fmt.Errorf("storage base path cannot be empty")
	}
	
	if c.StorageMaxSize <= 0 {
		return fmt.Errorf("storage max size must be positive")
	}
	
	if c.ProcessorMaxWidth <= 0 || c.ProcessorMaxHeight <= 0 {
		return fmt.Errorf("processor max dimensions must be positive")
	}
	
	if c.DefaultQuality < 1 || c.DefaultQuality > 100 {
		return fmt.Errorf("default quality must be between 1 and 100")
	}
	
	if c.CacheEnabled {
		if c.CacheCapacity <= 0 {
			return fmt.Errorf("cache capacity must be positive when cache is enabled")
		}
		
		if c.CacheMaxSizeBytes <= 0 {
			return fmt.Errorf("cache max size must be positive when cache is enabled")
		}
		
		if c.CacheTTL <= 0 {
			return fmt.Errorf("cache TTL must be positive when cache is enabled")
		}
	}
	
	if c.MaxUploadSize <= 0 {
		return fmt.Errorf("max upload size must be positive")
	}
	
	if c.MaxImagesPerProperty <= 0 {
		return fmt.Errorf("max images per property must be positive")
	}
	
	return nil
}

// GetSummary returns a summary of the configuration for logging
func (c *ImageConfig) GetSummary() map[string]interface{} {
	return map[string]interface{}{
		"storage_base_path":          c.StorageBasePath,
		"storage_base_url":           c.StorageBaseURL,
		"storage_max_size_mb":        c.StorageMaxSize / (1024 * 1024),
		"processor_max_width":        c.ProcessorMaxWidth,
		"processor_max_height":       c.ProcessorMaxHeight,
		"default_quality":            c.DefaultQuality,
		"cache_enabled":              c.CacheEnabled,
		"cache_capacity":             c.CacheCapacity,
		"cache_max_size_mb":          c.CacheMaxSizeBytes / (1024 * 1024),
		"cache_ttl_minutes":          c.CacheTTL.Minutes(),
		"cache_cleanup_interval_min": c.CacheCleanupInterval.Minutes(),
		"max_upload_size_mb":         c.MaxUploadSize / (1024 * 1024),
		"max_images_per_property":    c.MaxImagesPerProperty,
		"enable_optimization":        c.EnableImageOptimization,
		"enable_variant_generation":  c.EnableVariantGeneration,
		"enable_thumbnail_cache":     c.EnableThumbnailCache,
	}
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// Development configuration
func DevelopmentImageConfig() *ImageConfig {
	return &ImageConfig{
		StorageBasePath:         "./storage/images",
		StorageBaseURL:          "http://localhost:8080/images",
		StorageMaxSize:          512 * 1024 * 1024, // 512MB
		ProcessorMaxWidth:       2000,
		ProcessorMaxHeight:      1500,
		DefaultQuality:          85,
		CacheEnabled:            true,
		CacheCapacity:           500,
		CacheMaxSizeBytes:       50 * 1024 * 1024, // 50MB
		CacheTTL:                30 * time.Minute,
		CacheCleanupInterval:    5 * time.Minute,
		MaxUploadSize:           5 * 1024 * 1024, // 5MB
		MaxImagesPerProperty:    20,
		EnableImageOptimization: true,
		EnableVariantGeneration: true,
		EnableThumbnailCache:    true,
	}
}

// Production configuration
func ProductionImageConfig() *ImageConfig {
	return &ImageConfig{
		StorageBasePath:         "/var/lib/inmobiliaria/images",
		StorageBaseURL:          "https://example.com/images",
		StorageMaxSize:          10 * 1024 * 1024 * 1024, // 10GB
		ProcessorMaxWidth:       4000,
		ProcessorMaxHeight:      3000,
		DefaultQuality:          90,
		CacheEnabled:            true,
		CacheCapacity:           2000,
		CacheMaxSizeBytes:       500 * 1024 * 1024, // 500MB
		CacheTTL:                2 * time.Hour,
		CacheCleanupInterval:    30 * time.Minute,
		MaxUploadSize:           20 * 1024 * 1024, // 20MB
		MaxImagesPerProperty:    100,
		EnableImageOptimization: true,
		EnableVariantGeneration: true,
		EnableThumbnailCache:    true,
	}
}

// Test configuration
func TestImageConfig() *ImageConfig {
	return &ImageConfig{
		StorageBasePath:         "./test_storage/images",
		StorageBaseURL:          "http://localhost:8080/test/images",
		StorageMaxSize:          100 * 1024 * 1024, // 100MB
		ProcessorMaxWidth:       1000,
		ProcessorMaxHeight:      1000,
		DefaultQuality:          80,
		CacheEnabled:            false, // Disabled for testing
		CacheCapacity:           100,
		CacheMaxSizeBytes:       10 * 1024 * 1024, // 10MB
		CacheTTL:                5 * time.Minute,
		CacheCleanupInterval:    1 * time.Minute,
		MaxUploadSize:           2 * 1024 * 1024, // 2MB
		MaxImagesPerProperty:    10,
		EnableImageOptimization: true,
		EnableVariantGeneration: true,
		EnableThumbnailCache:    false,
	}
}