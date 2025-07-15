package domain

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ImageInfo represents image metadata and URLs
type ImageInfo struct {
	ID           string    `json:"id"`
	PropertyID   string    `json:"property_id"`
	FileName     string    `json:"file_name"`
	OriginalURL  string    `json:"original_url"`
	AltText      string    `json:"alt_text"`
	SortOrder    int       `json:"sort_order"`
	Size         int64     `json:"size"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	Format       string    `json:"format"`
	Quality      int       `json:"quality"`
	IsOptimized  bool      `json:"is_optimized"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ImageUploadRequest represents request for uploading images
type ImageUploadRequest struct {
	PropertyID string `json:"property_id"`
	AltText    string `json:"alt_text"`
	SortOrder  int    `json:"sort_order"`
}

// ImageReorderRequest represents request for reordering images
type ImageReorderRequest struct {
	ImageIDs []string `json:"image_ids"`
}

// ProcessingOptions represents image processing configuration
type ProcessingOptions struct {
	MaxWidth       int     `json:"max_width"`
	MaxHeight      int     `json:"max_height"`
	Quality        int     `json:"quality"`
	Format         string  `json:"format"`
	OptimizeSize   bool    `json:"optimize_size"`
	PreserveAspect bool    `json:"preserve_aspect"`
}

// ImageStats represents processing statistics
type ImageStats struct {
	OriginalSize     int64   `json:"original_size"`
	OptimizedSize    int64   `json:"optimized_size"`
	CompressionRatio float64 `json:"compression_ratio"`
	ProcessingTime   int64   `json:"processing_time_ms"`
}

// NewImageInfo creates a new ImageInfo with generated ID
func NewImageInfo(propertyID, fileName string) *ImageInfo {
	return &ImageInfo{
		ID:          uuid.New().String(),
		PropertyID:  propertyID,
		FileName:    fileName,
		AltText:     "",
		SortOrder:   0,
		Size:        0,
		Width:       0,
		Height:      0,
		Format:      "",
		Quality:     85,
		IsOptimized: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateMetadata updates image metadata
func (img *ImageInfo) UpdateMetadata(altText string, sortOrder int) {
	img.AltText = altText
	img.SortOrder = sortOrder
	img.UpdatedAt = time.Now()
}

// SetProcessingResults sets results from image processing
func (img *ImageInfo) SetProcessingResults(width, height int, size int64, format string, quality int, isOptimized bool) {
	img.Width = width
	img.Height = height
	img.Size = size
	img.Format = format
	img.Quality = quality
	img.IsOptimized = isOptimized
	img.UpdatedAt = time.Now()
}

// GetOptimizedURL returns URL for optimized image with parameters
func (img *ImageInfo) GetOptimizedURL(baseURL string, width, height int, format string, quality int) string {
	if width == 0 && height == 0 && format == "" && quality == 0 {
		return img.OriginalURL
	}
	
	params := []string{}
	if width > 0 {
		params = append(params, fmt.Sprintf("w=%d", width))
	}
	if height > 0 {
		params = append(params, fmt.Sprintf("h=%d", height))
	}
	if format != "" {
		params = append(params, fmt.Sprintf("f=%s", format))
	}
	if quality > 0 {
		params = append(params, fmt.Sprintf("q=%d", quality))
	}
	
	separator := "?"
	if strings.Contains(img.OriginalURL, "?") {
		separator = "&"
	}
	
	if len(params) > 0 {
		return fmt.Sprintf("%s%s%s", img.OriginalURL, separator, strings.Join(params, "&"))
	}
	
	return img.OriginalURL
}

// Validate validates image information
func (img *ImageInfo) Validate() error {
	if img.PropertyID == "" {
		return fmt.Errorf("property_id is required")
	}
	
	if img.FileName == "" {
		return fmt.Errorf("file_name is required")
	}
	
	if img.Size < 0 {
		return fmt.Errorf("size must be non-negative")
	}
	
	if img.Width < 0 || img.Height < 0 {
		return fmt.Errorf("width and height must be non-negative")
	}
	
	if img.Quality < 1 || img.Quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100")
	}
	
	if img.SortOrder < 0 {
		return fmt.Errorf("sort_order must be non-negative")
	}
	
	return nil
}

// IsValidImageFormat checks if format is supported
func IsValidImageFormat(format string) bool {
	validFormats := []string{"jpg", "jpeg", "png", "webp", "avif"}
	format = strings.ToLower(format)
	
	for _, valid := range validFormats {
		if format == valid {
			return true
		}
	}
	
	return false
}

// GetImageFormatFromFilename extracts format from filename
func GetImageFormatFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".jpg", ".jpeg":
		return "jpg"
	case ".png":
		return "png"
	case ".webp":
		return "webp"
	case ".avif":
		return "avif"
	default:
		return ""
	}
}

// GenerateImageFileName generates unique filename for image
func GenerateImageFileName(propertyID string, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".jpg"
	}
	
	// Generate unique filename
	uniqueID := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s%s", propertyID, uniqueID, ext)
}

// ValidateProcessingOptions validates processing options
func ValidateProcessingOptions(opts ProcessingOptions) error {
	if opts.MaxWidth < 0 || opts.MaxHeight < 0 {
		return fmt.Errorf("max dimensions must be non-negative")
	}
	
	if opts.Quality < 1 || opts.Quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100")
	}
	
	if opts.Format != "" && !IsValidImageFormat(opts.Format) {
		return fmt.Errorf("invalid image format: %s", opts.Format)
	}
	
	return nil
}

// DefaultProcessingOptions returns default processing options
func DefaultProcessingOptions() ProcessingOptions {
	return ProcessingOptions{
		MaxWidth:       3000,
		MaxHeight:      2000,
		Quality:        85,
		Format:         "jpg",
		OptimizeSize:   true,
		PreserveAspect: true,
	}
}

// CalculateCompressionRatio calculates compression ratio
func CalculateCompressionRatio(originalSize, compressedSize int64) float64 {
	if originalSize == 0 {
		return 0.0
	}
	
	return float64(compressedSize) / float64(originalSize)
}

// GetImageSizeCategory returns size category for image
func GetImageSizeCategory(size int64) string {
	const (
		Small  = 100 * 1024    // 100KB
		Medium = 500 * 1024    // 500KB
		Large  = 2 * 1024 * 1024 // 2MB
	)
	
	switch {
	case size <= Small:
		return "small"
	case size <= Medium:
		return "medium"
	case size <= Large:
		return "large"
	default:
		return "xlarge"
	}
}

// Common image processing constants
const (
	MaxUploadSize     = int64(10 * 1024 * 1024) // 10MB
	MaxImageWidth     = 4000
	MaxImageHeight    = 4000
	DefaultQuality    = 85
	ThumbnailSize     = 150
	MediumSize        = 800
	LargeSize         = 1200
	MaxImagesPerProperty = 50
)

// Supported MIME types
var SupportedMimeTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/jpg":  "jpg",
	"image/png":  "png",
	"image/webp": "webp",
	"image/avif": "avif",
}

// GetFormatFromMimeType returns format from MIME type
func GetFormatFromMimeType(mimeType string) string {
	if format, exists := SupportedMimeTypes[strings.ToLower(mimeType)]; exists {
		return format
	}
	return ""
}

// IsSupportedMimeType checks if MIME type is supported
func IsSupportedMimeType(mimeType string) bool {
	_, exists := SupportedMimeTypes[strings.ToLower(mimeType)]
	return exists
}