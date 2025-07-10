package processors

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"strings"
	"time"

	"golang.org/x/image/draw"
	"realty-core/internal/domain"
)

// ImageProcessor handles image processing operations
type ImageProcessor struct {
	maxWidth  int
	maxHeight int
}

// NewImageProcessor creates a new image processor
func NewImageProcessor(maxWidth, maxHeight int) *ImageProcessor {
	if maxWidth <= 0 {
		maxWidth = 3000
	}
	if maxHeight <= 0 {
		maxHeight = 2000
	}
	
	return &ImageProcessor{
		maxWidth:  maxWidth,
		maxHeight: maxHeight,
	}
}

// ProcessImage processes an image with the given options
func (ip *ImageProcessor) ProcessImage(inputData []byte, options domain.ProcessingOptions) ([]byte, *domain.ImageStats, error) {
	start := time.Now()
	originalSize := int64(len(inputData))
	
	// Validate options
	if err := domain.ValidateProcessingOptions(options); err != nil {
		return nil, nil, fmt.Errorf("invalid processing options: %w", err)
	}
	
	// Decode input image
	inputImage, inputFormat, err := image.Decode(bytes.NewReader(inputData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode image: %w", err)
	}
	
	// Process the image
	processedImage, err := ip.processImageWithOptions(inputImage, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process image: %w", err)
	}
	
	// Encode output image
	outputData, err := ip.encodeImage(processedImage, options.Format, options.Quality)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode image: %w", err)
	}
	
	// Calculate statistics
	stats := &domain.ImageStats{
		OriginalSize:     originalSize,
		OptimizedSize:    int64(len(outputData)),
		CompressionRatio: domain.CalculateCompressionRatio(originalSize, int64(len(outputData))),
		ProcessingTime:   time.Since(start).Milliseconds(),
	}
	
	log.Printf("Image processed: %s -> %s, %s compression ratio: %.2f, time: %dms",
		inputFormat, options.Format, formatBytes(originalSize), stats.CompressionRatio, stats.ProcessingTime)
	
	return outputData, stats, nil
}

// processImageWithOptions applies processing options to the image
func (ip *ImageProcessor) processImageWithOptions(inputImage image.Image, options domain.ProcessingOptions) (image.Image, error) {
	bounds := inputImage.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()
	
	// Calculate new dimensions
	newWidth, newHeight := ip.calculateDimensions(originalWidth, originalHeight, options)
	
	// If no resizing needed, return original
	if newWidth == originalWidth && newHeight == originalHeight {
		return inputImage, nil
	}
	
	// Create new image
	newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	
	// Use high-quality scaling
	if options.OptimizeSize {
		draw.BiLinear.Scale(newImage, newImage.Bounds(), inputImage, inputImage.Bounds(), draw.Over, nil)
	} else {
		draw.NearestNeighbor.Scale(newImage, newImage.Bounds(), inputImage, inputImage.Bounds(), draw.Over, nil)
	}
	
	return newImage, nil
}

// calculateDimensions calculates new image dimensions based on constraints
func (ip *ImageProcessor) calculateDimensions(originalWidth, originalHeight int, options domain.ProcessingOptions) (int, int) {
	maxWidth := options.MaxWidth
	maxHeight := options.MaxHeight
	
	// Use processor defaults if not specified
	if maxWidth <= 0 {
		maxWidth = ip.maxWidth
	}
	if maxHeight <= 0 {
		maxHeight = ip.maxHeight
	}
	
	// If original is smaller than max, keep original
	if originalWidth <= maxWidth && originalHeight <= maxHeight {
		return originalWidth, originalHeight
	}
	
	// Calculate scaling ratio
	widthRatio := float64(maxWidth) / float64(originalWidth)
	heightRatio := float64(maxHeight) / float64(originalHeight)
	
	// Use the smaller ratio to maintain aspect ratio
	var ratio float64
	if options.PreserveAspect {
		if widthRatio < heightRatio {
			ratio = widthRatio
		} else {
			ratio = heightRatio
		}
	} else {
		ratio = widthRatio
	}
	
	newWidth := int(float64(originalWidth) * ratio)
	newHeight := int(float64(originalHeight) * ratio)
	
	// Ensure minimum dimensions
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}
	
	return newWidth, newHeight
}

// encodeImage encodes the image with the specified format and quality
func (ip *ImageProcessor) encodeImage(img image.Image, format string, quality int) ([]byte, error) {
	var buffer bytes.Buffer
	
	switch strings.ToLower(format) {
	case "jpg", "jpeg":
		err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, fmt.Errorf("failed to encode as JPEG: %w", err)
		}
	case "png":
		err := png.Encode(&buffer, img)
		if err != nil {
			return nil, fmt.Errorf("failed to encode as PNG: %w", err)
		}
	case "webp":
		return nil, fmt.Errorf("WebP encoding not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	
	return buffer.Bytes(), nil
}

// OptimizeForSize optimizes image with intelligent quality adjustment
func (ip *ImageProcessor) OptimizeForSize(inputData []byte, targetSizeKB int) ([]byte, *domain.ImageStats, error) {
	if targetSizeKB <= 0 {
		targetSizeKB = 1200 // Default 1.2MB
	}
	
	targetSize := int64(targetSizeKB * 1024)
	originalSize := int64(len(inputData))
	
	// If already smaller than target, apply minimal optimization
	if originalSize <= targetSize {
		options := domain.ProcessingOptions{
			MaxWidth:       ip.maxWidth,
			MaxHeight:      ip.maxHeight,
			Quality:        95,
			Format:         "jpg",
			OptimizeSize:   true,
			PreserveAspect: true,
		}
		return ip.ProcessImage(inputData, options)
	}
	
	// Try different quality levels
	qualityLevels := []int{85, 75, 65, 55, 45}
	
	for _, quality := range qualityLevels {
		options := domain.ProcessingOptions{
			MaxWidth:       ip.maxWidth,
			MaxHeight:      ip.maxHeight,
			Quality:        quality,
			Format:         "jpg",
			OptimizeSize:   true,
			PreserveAspect: true,
		}
		
		result, stats, err := ip.ProcessImage(inputData, options)
		if err != nil {
			continue
		}
		
		// If size is acceptable, return this result
		if stats.OptimizedSize <= targetSize {
			return result, stats, nil
		}
	}
	
	// If still too large, try reducing dimensions
	options := domain.ProcessingOptions{
		MaxWidth:       2000,
		MaxHeight:      1500,
		Quality:        45,
		Format:         "jpg",
		OptimizeSize:   true,
		PreserveAspect: true,
	}
	
	return ip.ProcessImage(inputData, options)
}

// GetImageDimensions returns image dimensions without full processing
func (ip *ImageProcessor) GetImageDimensions(inputData []byte) (int, int, string, error) {
	img, format, err := image.Decode(bytes.NewReader(inputData))
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed to decode image: %w", err)
	}
	
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), format, nil
}

// ValidateImageData validates image data and returns basic info
func (ip *ImageProcessor) ValidateImageData(data []byte, maxSize int64) error {
	if len(data) == 0 {
		return fmt.Errorf("empty image data")
	}
	
	if maxSize > 0 && int64(len(data)) > maxSize {
		return fmt.Errorf("image too large: %d bytes, max: %d bytes", len(data), maxSize)
	}
	
	// Try to decode to validate format
	_, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("invalid image format: %w", err)
	}
	
	return nil
}

// GenerateThumbnail generates a thumbnail with fixed dimensions
func (ip *ImageProcessor) GenerateThumbnail(inputData []byte, size int) ([]byte, error) {
	if size <= 0 {
		size = domain.ThumbnailSize
	}
	
	options := domain.ProcessingOptions{
		MaxWidth:       size,
		MaxHeight:      size,
		Quality:        80,
		Format:         "jpg",
		OptimizeSize:   true,
		PreserveAspect: true,
	}
	
	result, _, err := ip.ProcessImage(inputData, options)
	return result, err
}

// GenerateImageVariant generates specific image variant
func (ip *ImageProcessor) GenerateImageVariant(inputData []byte, width, height int, quality int, format string) ([]byte, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: %dx%d", width, height)
	}
	
	if quality <= 0 || quality > 100 {
		quality = domain.DefaultQuality
	}
	
	if format == "" {
		format = "jpg"
	}
	
	options := domain.ProcessingOptions{
		MaxWidth:       width,
		MaxHeight:      height,
		Quality:        quality,
		Format:         format,
		OptimizeSize:   true,
		PreserveAspect: true,
	}
	
	result, _, err := ip.ProcessImage(inputData, options)
	return result, err
}

// formatBytes formats bytes to human readable string
func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
}