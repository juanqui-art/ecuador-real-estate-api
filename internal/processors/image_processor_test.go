package processors

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"realty-core/internal/domain"
)

// createTestImage creates a test image for testing
func createTestImage(width, height int, format string) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Fill with a simple pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, image.NewRGBA(image.Rectangle{}).ColorModel().Convert(
				image.NewUniform(image.NewRGBA(image.Rectangle{}).ColorModel().Convert(
					image.NewRGBA(image.Rectangle{}).At(x%256, y%256)))))
		}
	}
	
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	case "png":
		png.Encode(&buf, img)
	default:
		jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	}
	
	return buf.Bytes()
}

func TestNewImageProcessor(t *testing.T) {
	tests := []struct {
		name      string
		maxWidth  int
		maxHeight int
		wantWidth int
		wantHeight int
	}{
		{
			name:       "valid dimensions",
			maxWidth:   1920,
			maxHeight:  1080,
			wantWidth:  1920,
			wantHeight: 1080,
		},
		{
			name:       "zero dimensions should use defaults",
			maxWidth:   0,
			maxHeight:  0,
			wantWidth:  3000,
			wantHeight: 2000,
		},
		{
			name:       "negative dimensions should use defaults",
			maxWidth:   -100,
			maxHeight:  -50,
			wantWidth:  3000,
			wantHeight: 2000,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewImageProcessor(tt.maxWidth, tt.maxHeight)
			
			assert.NotNil(t, processor)
			assert.Equal(t, tt.wantWidth, processor.maxWidth)
			assert.Equal(t, tt.wantHeight, processor.maxHeight)
		})
	}
}

func TestImageProcessor_ValidateImageData(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	tests := []struct {
		name      string
		data      []byte
		maxSize   int64
		wantError bool
	}{
		{
			name:      "empty data",
			data:      []byte{},
			maxSize:   1024,
			wantError: true,
		},
		{
			name:      "valid jpeg data",
			data:      createTestImage(100, 100, "jpeg"),
			maxSize:   100000,
			wantError: false,
		},
		{
			name:      "data too large",
			data:      createTestImage(100, 100, "jpeg"),
			maxSize:   100,
			wantError: true,
		},
		{
			name:      "invalid image data",
			data:      []byte("not an image"),
			maxSize:   1024,
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processor.ValidateImageData(tt.data, tt.maxSize)
			
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestImageProcessor_GetImageDimensions(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	tests := []struct {
		name       string
		data       []byte
		wantWidth  int
		wantHeight int
		wantFormat string
		wantError  bool
	}{
		{
			name:       "valid jpeg",
			data:       createTestImage(200, 150, "jpeg"),
			wantWidth:  200,
			wantHeight: 150,
			wantFormat: "jpeg",
			wantError:  false,
		},
		{
			name:       "valid png",
			data:       createTestImage(300, 200, "png"),
			wantWidth:  300,
			wantHeight: 200,
			wantFormat: "png",
			wantError:  false,
		},
		{
			name:      "invalid data",
			data:      []byte("not an image"),
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height, format, err := processor.GetImageDimensions(tt.data)
			
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantWidth, width)
				assert.Equal(t, tt.wantHeight, height)
				assert.Equal(t, tt.wantFormat, format)
			}
		})
	}
}

func TestImageProcessor_ProcessImage(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	tests := []struct {
		name      string
		inputData []byte
		options   domain.ProcessingOptions
		wantError bool
	}{
		{
			name:      "valid jpeg processing",
			inputData: createTestImage(500, 300, "jpeg"),
			options: domain.ProcessingOptions{
				MaxWidth:       400,
				MaxHeight:      300,
				Quality:        80,
				Format:         "jpg",
				OptimizeSize:   true,
				PreserveAspect: true,
			},
			wantError: false,
		},
		{
			name:      "invalid options",
			inputData: createTestImage(500, 300, "jpeg"),
			options: domain.ProcessingOptions{
				MaxWidth:       -100,
				MaxHeight:      -50,
				Quality:        80,
				Format:         "jpg",
				OptimizeSize:   true,
				PreserveAspect: true,
			},
			wantError: true,
		},
		{
			name:      "unsupported format",
			inputData: createTestImage(500, 300, "jpeg"),
			options: domain.ProcessingOptions{
				MaxWidth:       400,
				MaxHeight:      300,
				Quality:        80,
				Format:         "webp",
				OptimizeSize:   true,
				PreserveAspect: true,
			},
			wantError: true,
		},
		{
			name:      "invalid input data",
			inputData: []byte("not an image"),
			options: domain.ProcessingOptions{
				MaxWidth:       400,
				MaxHeight:      300,
				Quality:        80,
				Format:         "jpg",
				OptimizeSize:   true,
				PreserveAspect: true,
			},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, stats, err := processor.ProcessImage(tt.inputData, tt.options)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Nil(t, stats)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, stats)
				assert.Greater(t, len(result), 0)
				assert.Greater(t, stats.OriginalSize, int64(0))
				assert.Greater(t, stats.OptimizedSize, int64(0))
				assert.Greater(t, stats.ProcessingTime, int64(0))
			}
		})
	}
}

func TestImageProcessor_OptimizeForSize(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	tests := []struct {
		name         string
		inputData    []byte
		targetSizeKB int
		wantError    bool
	}{
		{
			name:         "valid optimization",
			inputData:    createTestImage(800, 600, "jpeg"),
			targetSizeKB: 500,
			wantError:    false,
		},
		{
			name:         "zero target size should use default",
			inputData:    createTestImage(800, 600, "jpeg"),
			targetSizeKB: 0,
			wantError:    false,
		},
		{
			name:         "invalid input data",
			inputData:    []byte("not an image"),
			targetSizeKB: 500,
			wantError:    true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, stats, err := processor.OptimizeForSize(tt.inputData, tt.targetSizeKB)
			
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, stats)
				assert.Greater(t, len(result), 0)
				
				// Check that optimization actually reduced size (in most cases)
				if stats.OriginalSize > int64(tt.targetSizeKB*1024) {
					assert.LessOrEqual(t, stats.OptimizedSize, stats.OriginalSize)
				}
			}
		})
	}
}

func TestImageProcessor_GenerateThumbnail(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	tests := []struct {
		name      string
		inputData []byte
		size      int
		wantError bool
	}{
		{
			name:      "valid thumbnail generation",
			inputData: createTestImage(800, 600, "jpeg"),
			size:      150,
			wantError: false,
		},
		{
			name:      "zero size should use default",
			inputData: createTestImage(800, 600, "jpeg"),
			size:      0,
			wantError: false,
		},
		{
			name:      "invalid input data",
			inputData: []byte("not an image"),
			size:      150,
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.GenerateThumbnail(tt.inputData, tt.size)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result), 0)
			}
		})
	}
}

func TestImageProcessor_GenerateImageVariant(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	tests := []struct {
		name      string
		inputData []byte
		width     int
		height    int
		quality   int
		format    string
		wantError bool
	}{
		{
			name:      "valid variant generation",
			inputData: createTestImage(800, 600, "jpeg"),
			width:     400,
			height:    300,
			quality:   80,
			format:    "jpg",
			wantError: false,
		},
		{
			name:      "zero quality should use default",
			inputData: createTestImage(800, 600, "jpeg"),
			width:     400,
			height:    300,
			quality:   0,
			format:    "jpg",
			wantError: false,
		},
		{
			name:      "empty format should use default",
			inputData: createTestImage(800, 600, "jpeg"),
			width:     400,
			height:    300,
			quality:   80,
			format:    "",
			wantError: false,
		},
		{
			name:      "invalid dimensions",
			inputData: createTestImage(800, 600, "jpeg"),
			width:     -100,
			height:    -50,
			quality:   80,
			format:    "jpg",
			wantError: true,
		},
		{
			name:      "invalid input data",
			inputData: []byte("not an image"),
			width:     400,
			height:    300,
			quality:   80,
			format:    "jpg",
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.GenerateImageVariant(tt.inputData, tt.width, tt.height, tt.quality, tt.format)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result), 0)
			}
		})
	}
}

func TestImageProcessor_calculateDimensions(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	tests := []struct {
		name           string
		originalWidth  int
		originalHeight int
		options        domain.ProcessingOptions
		wantWidth      int
		wantHeight     int
	}{
		{
			name:           "no scaling needed",
			originalWidth:  800,
			originalHeight: 600,
			options: domain.ProcessingOptions{
				MaxWidth:       1920,
				MaxHeight:      1080,
				PreserveAspect: true,
			},
			wantWidth:  800,
			wantHeight: 600,
		},
		{
			name:           "scale down width",
			originalWidth:  2000,
			originalHeight: 1000,
			options: domain.ProcessingOptions{
				MaxWidth:       1000,
				MaxHeight:      1000,
				PreserveAspect: true,
			},
			wantWidth:  1000,
			wantHeight: 500,
		},
		{
			name:           "scale down height",
			originalWidth:  1000,
			originalHeight: 2000,
			options: domain.ProcessingOptions{
				MaxWidth:       1000,
				MaxHeight:      1000,
				PreserveAspect: true,
			},
			wantWidth:  500,
			wantHeight: 1000,
		},
		{
			name:           "use processor defaults when options are zero",
			originalWidth:  3000,
			originalHeight: 2500,
			options: domain.ProcessingOptions{
				MaxWidth:       0,
				MaxHeight:      0,
				PreserveAspect: true,
			},
			wantWidth:  1920,
			wantHeight: 1600,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight := processor.calculateDimensions(tt.originalWidth, tt.originalHeight, tt.options)
			
			assert.Equal(t, tt.wantWidth, gotWidth)
			assert.Equal(t, tt.wantHeight, gotHeight)
		})
	}
}

func TestImageProcessor_encodeImage(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	tests := []struct {
		name      string
		img       image.Image
		format    string
		quality   int
		wantError bool
	}{
		{
			name:      "valid jpeg encoding",
			img:       img,
			format:    "jpg",
			quality:   80,
			wantError: false,
		},
		{
			name:      "valid png encoding",
			img:       img,
			format:    "png",
			quality:   80,
			wantError: false,
		},
		{
			name:      "unsupported format",
			img:       img,
			format:    "webp",
			quality:   80,
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.encodeImage(tt.img, tt.format, tt.quality)
			
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Greater(t, len(result), 0)
			}
		})
	}
}

func TestImageProcessor_Integration(t *testing.T) {
	processor := NewImageProcessor(1920, 1080)
	
	// Create a large test image
	largeImageData := createTestImage(2000, 1500, "jpeg")
	
	t.Run("full processing pipeline", func(t *testing.T) {
		options := domain.ProcessingOptions{
			MaxWidth:       1000,
			MaxHeight:      800,
			Quality:        85,
			Format:         "jpg",
			OptimizeSize:   true,
			PreserveAspect: true,
		}
		
		result, stats, err := processor.ProcessImage(largeImageData, options)
		
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, stats)
		
		// Verify the result is a valid image
		_, _, _, err = processor.GetImageDimensions(result)
		assert.NoError(t, err)
		
		// Verify processing reduced file size
		assert.Less(t, stats.OptimizedSize, stats.OriginalSize)
		
		// Verify compression ratio is reasonable
		assert.Greater(t, stats.CompressionRatio, 0.1)
		assert.Less(t, stats.CompressionRatio, 1.0)
	})
	
	t.Run("optimization for size", func(t *testing.T) {
		targetSizeKB := 200
		
		result, stats, err := processor.OptimizeForSize(largeImageData, targetSizeKB)
		
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, stats)
		
		// Verify the result is a valid image
		_, _, _, err = processor.GetImageDimensions(result)
		assert.NoError(t, err)
		
		// Verify size is within reasonable bounds
		targetSize := int64(targetSizeKB * 1024)
		if stats.OriginalSize > targetSize {
			assert.LessOrEqual(t, stats.OptimizedSize, targetSize*2) // Allow some margin
		}
	})
}