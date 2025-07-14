package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewImageInfo(t *testing.T) {
	propertyID := "test-property-id"
	fileName := "test.jpg"
	
	image := NewImageInfo(propertyID, fileName)
	
	assert.NotNil(t, image)
	assert.NotEmpty(t, image.ID)
	assert.Equal(t, propertyID, image.PropertyID)
	assert.Equal(t, fileName, image.FileName)
	assert.Equal(t, "", image.AltText)
	assert.Equal(t, 0, image.SortOrder)
	assert.Equal(t, int64(0), image.Size)
	assert.Equal(t, 0, image.Width)
	assert.Equal(t, 0, image.Height)
	assert.Equal(t, "", image.Format)
	assert.Equal(t, 85, image.Quality)
	assert.False(t, image.IsOptimized)
	assert.False(t, image.CreatedAt.IsZero())
	assert.False(t, image.UpdatedAt.IsZero())
}

func TestImageInfo_UpdateMetadata(t *testing.T) {
	image := NewImageInfo("test-property", "test.jpg")
	originalUpdatedAt := image.UpdatedAt
	
	// Sleep briefly to ensure timestamp changes
	time.Sleep(time.Millisecond)
	
	altText := "Test alt text"
	sortOrder := 5
	
	image.UpdateMetadata(altText, sortOrder)
	
	assert.Equal(t, altText, image.AltText)
	assert.Equal(t, sortOrder, image.SortOrder)
	assert.True(t, image.UpdatedAt.After(originalUpdatedAt))
}

func TestImageInfo_SetProcessingResults(t *testing.T) {
	image := NewImageInfo("test-property", "test.jpg")
	originalUpdatedAt := image.UpdatedAt
	
	// Sleep briefly to ensure timestamp changes
	time.Sleep(time.Millisecond)
	
	width := 1920
	height := 1080
	size := int64(1024 * 1024)
	format := "jpg"
	quality := 85
	isOptimized := true
	
	image.SetProcessingResults(width, height, size, format, quality, isOptimized)
	
	assert.Equal(t, width, image.Width)
	assert.Equal(t, height, image.Height)
	assert.Equal(t, size, image.Size)
	assert.Equal(t, format, image.Format)
	assert.Equal(t, quality, image.Quality)
	assert.Equal(t, isOptimized, image.IsOptimized)
	assert.True(t, image.UpdatedAt.After(originalUpdatedAt))
}

func TestImageInfo_GetOptimizedURL(t *testing.T) {
	image := NewImageInfo("test-property", "test.jpg")
	image.OriginalURL = "http://example.com/image.jpg"
	baseURL := "http://example.com"
	
	tests := []struct {
		name     string
		width    int
		height   int
		format   string
		quality  int
		expected string
	}{
		{
			name:     "no parameters should return original URL",
			width:    0,
			height:   0,
			format:   "",
			quality:  0,
			expected: image.OriginalURL,
		},
		{
			name:     "width only",
			width:    800,
			height:   0,
			format:   "",
			quality:  0,
			expected: image.OriginalURL + "?w=800",
		},
		{
			name:     "height only",
			width:    0,
			height:   600,
			format:   "",
			quality:  0,
			expected: image.OriginalURL + "?h=600",
		},
		{
			name:     "width and height",
			width:    800,
			height:   600,
			format:   "",
			quality:  0,
			expected: image.OriginalURL + "?w=800&h=600",
		},
		{
			name:     "all parameters",
			width:    800,
			height:   600,
			format:   "webp",
			quality:  80,
			expected: image.OriginalURL + "?w=800&h=600&f=webp&q=80",
		},
		{
			name:     "format and quality only",
			width:    0,
			height:   0,
			format:   "webp",
			quality:  80,
			expected: image.OriginalURL + "?f=webp&q=80",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := image.GetOptimizedURL(baseURL, tt.width, tt.height, tt.format, tt.quality)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestImageInfo_GetOptimizedURL_ExistingQueryParams(t *testing.T) {
	image := NewImageInfo("test-property", "test.jpg")
	image.OriginalURL = "http://example.com/image.jpg?existing=param"
	baseURL := "http://example.com"
	
	result := image.GetOptimizedURL(baseURL, 800, 600, "webp", 80)
	expected := image.OriginalURL + "&w=800&h=600&f=webp&q=80"
	
	assert.Equal(t, expected, result)
}

func TestImageInfo_Validate(t *testing.T) {
	tests := []struct {
		name    string
		image   *ImageInfo
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid image",
			image: &ImageInfo{
				PropertyID: "test-property",
				FileName:   "test.jpg",
				Size:       1024,
				Width:      800,
				Height:     600,
				Quality:    85,
				SortOrder:  0,
			},
			wantErr: false,
		},
		{
			name: "missing property ID",
			image: &ImageInfo{
				PropertyID: "",
				FileName:   "test.jpg",
			},
			wantErr: true,
			errMsg:  "property_id is required",
		},
		{
			name: "missing file name",
			image: &ImageInfo{
				PropertyID: "test-property",
				FileName:   "",
			},
			wantErr: true,
			errMsg:  "file_name is required",
		},
		{
			name: "negative size",
			image: &ImageInfo{
				PropertyID: "test-property",
				FileName:   "test.jpg",
				Size:       -1,
			},
			wantErr: true,
			errMsg:  "size must be non-negative",
		},
		{
			name: "negative width",
			image: &ImageInfo{
				PropertyID: "test-property",
				FileName:   "test.jpg",
				Width:      -1,
			},
			wantErr: true,
			errMsg:  "width and height must be non-negative",
		},
		{
			name: "negative height",
			image: &ImageInfo{
				PropertyID: "test-property",
				FileName:   "test.jpg",
				Height:     -1,
			},
			wantErr: true,
			errMsg:  "width and height must be non-negative",
		},
		{
			name: "quality too low",
			image: &ImageInfo{
				PropertyID: "test-property",
				FileName:   "test.jpg",
				Quality:    0,
			},
			wantErr: true,
			errMsg:  "quality must be between 1 and 100",
		},
		{
			name: "quality too high",
			image: &ImageInfo{
				PropertyID: "test-property",
				FileName:   "test.jpg",
				Quality:    101,
			},
			wantErr: true,
			errMsg:  "quality must be between 1 and 100",
		},
		{
			name: "negative sort order",
			image: &ImageInfo{
				PropertyID: "test-property",
				FileName:   "test.jpg",
				Quality:    50, // Set valid quality to avoid triggering quality error first
				SortOrder:  -1,
			},
			wantErr: true,
			errMsg:  "sort_order must be non-negative",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.image.Validate()
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidImageFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   bool
	}{
		{"jpg", "jpg", true},
		{"jpeg", "jpeg", true},
		{"png", "png", true},
		{"webp", "webp", true},
		{"avif", "avif", true},
		{"uppercase JPG", "JPG", true},
		{"uppercase JPEG", "JPEG", true},
		{"mixed case WebP", "WebP", true},
		{"invalid format", "gif", false},
		{"empty format", "", false},
		{"random string", "random", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidImageFormat(tt.format)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGetImageFormatFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{"jpg extension", "image.jpg", "jpg"},
		{"jpeg extension", "image.jpeg", "jpg"},
		{"png extension", "image.png", "png"},
		{"webp extension", "image.webp", "webp"},
		{"avif extension", "image.avif", "avif"},
		{"uppercase extension", "image.JPG", "jpg"},
		{"mixed case extension", "image.Png", "png"},
		{"no extension", "image", ""},
		{"unknown extension", "image.gif", ""},
		{"multiple dots", "image.backup.jpg", "jpg"},
		{"empty filename", "", ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetImageFormatFromFilename(tt.filename)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenerateImageFileName(t *testing.T) {
	propertyID := "test-property-123"
	originalFilename := "original.jpg"
	
	result := GenerateImageFileName(propertyID, originalFilename)
	
	assert.NotEmpty(t, result)
	assert.True(t, strings.HasPrefix(result, propertyID+"_"))
	assert.True(t, strings.HasSuffix(result, ".jpg"))
	
	// Should contain an 8-character UUID segment
	parts := strings.Split(result, "_")
	assert.Len(t, parts, 2)
	assert.Equal(t, propertyID, parts[0])
	
	filenameWithExt := parts[1]
	assert.True(t, strings.HasSuffix(filenameWithExt, ".jpg"))
	
	uuidPart := strings.TrimSuffix(filenameWithExt, ".jpg")
	assert.Len(t, uuidPart, 8)
}

func TestGenerateImageFileName_NoExtension(t *testing.T) {
	propertyID := "test-property-123"
	originalFilename := "original"
	
	result := GenerateImageFileName(propertyID, originalFilename)
	
	assert.NotEmpty(t, result)
	assert.True(t, strings.HasPrefix(result, propertyID+"_"))
	assert.True(t, strings.HasSuffix(result, ".jpg")) // Should default to .jpg
}

func TestValidateProcessingOptions(t *testing.T) {
	tests := []struct {
		name    string
		options ProcessingOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid options",
			options: ProcessingOptions{
				MaxWidth:       800,
				MaxHeight:      600,
				Quality:        85,
				Format:         "jpg",
				OptimizeSize:   true,
				PreserveAspect: true,
			},
			wantErr: false,
		},
		{
			name: "negative max width",
			options: ProcessingOptions{
				MaxWidth:  -100,
				MaxHeight: 600,
				Quality:   85,
				Format:    "jpg",
			},
			wantErr: true,
			errMsg:  "max dimensions must be non-negative",
		},
		{
			name: "negative max height",
			options: ProcessingOptions{
				MaxWidth:  800,
				MaxHeight: -100,
				Quality:   85,
				Format:    "jpg",
			},
			wantErr: true,
			errMsg:  "max dimensions must be non-negative",
		},
		{
			name: "quality too low",
			options: ProcessingOptions{
				MaxWidth:  800,
				MaxHeight: 600,
				Quality:   0,
				Format:    "jpg",
			},
			wantErr: true,
			errMsg:  "quality must be between 1 and 100",
		},
		{
			name: "quality too high",
			options: ProcessingOptions{
				MaxWidth:  800,
				MaxHeight: 600,
				Quality:   101,
				Format:    "jpg",
			},
			wantErr: true,
			errMsg:  "quality must be between 1 and 100",
		},
		{
			name: "invalid format",
			options: ProcessingOptions{
				MaxWidth:  800,
				MaxHeight: 600,
				Quality:   85,
				Format:    "gif",
			},
			wantErr: true,
			errMsg:  "invalid image format: gif",
		},
		{
			name: "empty format is valid",
			options: ProcessingOptions{
				MaxWidth:  800,
				MaxHeight: 600,
				Quality:   85,
				Format:    "",
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProcessingOptions(tt.options)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDefaultProcessingOptions(t *testing.T) {
	options := DefaultProcessingOptions()
	
	assert.Equal(t, 3000, options.MaxWidth)
	assert.Equal(t, 2000, options.MaxHeight)
	assert.Equal(t, 85, options.Quality)
	assert.Equal(t, "jpg", options.Format)
	assert.True(t, options.OptimizeSize)
	assert.True(t, options.PreserveAspect)
}

func TestCalculateCompressionRatio(t *testing.T) {
	tests := []struct {
		name           string
		originalSize   int64
		compressedSize int64
		want           float64
	}{
		{
			name:           "50% compression",
			originalSize:   1000,
			compressedSize: 500,
			want:           0.5,
		},
		{
			name:           "75% compression",
			originalSize:   1000,
			compressedSize: 750,
			want:           0.75,
		},
		{
			name:           "no compression",
			originalSize:   1000,
			compressedSize: 1000,
			want:           1.0,
		},
		{
			name:           "zero original size",
			originalSize:   0,
			compressedSize: 500,
			want:           0.0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateCompressionRatio(tt.originalSize, tt.compressedSize)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGetImageSizeCategory(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{
			name: "small image",
			size: 50 * 1024, // 50KB
			want: "small",
		},
		{
			name: "medium image",
			size: 300 * 1024, // 300KB
			want: "medium",
		},
		{
			name: "large image",
			size: 1024 * 1024, // 1MB
			want: "large",
		},
		{
			name: "xlarge image",
			size: 5 * 1024 * 1024, // 5MB
			want: "xlarge",
		},
		{
			name: "boundary small",
			size: 100 * 1024, // exactly 100KB
			want: "small",
		},
		{
			name: "boundary medium",
			size: 500 * 1024, // exactly 500KB
			want: "medium",
		},
		{
			name: "boundary large",
			size: 2 * 1024 * 1024, // exactly 2MB
			want: "large",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetImageSizeCategory(tt.size)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGetFormatFromMimeType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		want     string
	}{
		{"image/jpeg", "image/jpeg", "jpg"},
		{"image/jpg", "image/jpg", "jpg"},
		{"image/png", "image/png", "png"},
		{"image/webp", "image/webp", "webp"},
		{"image/avif", "image/avif", "avif"},
		{"uppercase JPEG", "IMAGE/JPEG", "jpg"},
		{"mixed case PNG", "Image/PNG", "png"},
		{"unsupported type", "image/gif", ""},
		{"non-image type", "text/plain", ""},
		{"empty type", "", ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFormatFromMimeType(tt.mimeType)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestIsSupportedMimeType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		want     bool
	}{
		{"image/jpeg", "image/jpeg", true},
		{"image/jpg", "image/jpg", true},
		{"image/png", "image/png", true},
		{"image/webp", "image/webp", true},
		{"image/avif", "image/avif", true},
		{"uppercase JPEG", "IMAGE/JPEG", true},
		{"mixed case PNG", "Image/PNG", true},
		{"unsupported type", "image/gif", false},
		{"non-image type", "text/plain", false},
		{"empty type", "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSupportedMimeType(tt.mimeType)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have reasonable values
	assert.Equal(t, int64(10*1024*1024), MaxUploadSize)
	assert.Equal(t, 4000, MaxImageWidth)
	assert.Equal(t, 4000, MaxImageHeight)
	assert.Equal(t, 85, DefaultQuality)
	assert.Equal(t, 150, ThumbnailSize)
	assert.Equal(t, 800, MediumSize)
	assert.Equal(t, 1200, LargeSize)
	assert.Equal(t, 50, MaxImagesPerProperty)
	
	// Test that supported MIME types map is not empty
	assert.NotEmpty(t, SupportedMimeTypes)
	assert.Contains(t, SupportedMimeTypes, "image/jpeg")
	assert.Contains(t, SupportedMimeTypes, "image/png")
}