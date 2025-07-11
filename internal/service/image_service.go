package service

import (
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"realty-core/internal/cache"
	"realty-core/internal/domain"
	"realty-core/internal/processors"
	"realty-core/internal/repository"
	"realty-core/internal/storage"
)

// ImageServiceInterface defines the interface for image service operations
type ImageServiceInterface interface {
	// Upload uploads and processes a new image
	Upload(propertyID string, file multipart.File, header *multipart.FileHeader, altText string) (*domain.ImageInfo, error)
	
	// GetImage retrieves image metadata by ID
	GetImage(id string) (*domain.ImageInfo, error)
	
	// GetImagesByProperty retrieves all images for a property
	GetImagesByProperty(propertyID string) ([]domain.ImageInfo, error)
	
	// UpdateImageMetadata updates image metadata
	UpdateImageMetadata(id, altText string, sortOrder int) error
	
	// DeleteImage deletes image and its files
	DeleteImage(id string) error
	
	// ReorderImages reorders images for a property
	ReorderImages(propertyID string, imageIDs []string) error
	
	// SetMainImage sets an image as the main image for a property
	SetMainImage(propertyID, imageID string) error
	
	// GetMainImage gets the main image for a property
	GetMainImage(propertyID string) (*domain.ImageInfo, error)
	
	// GetImageVariant generates and returns an image variant
	GetImageVariant(imageID string, width, height int, format string, quality int) ([]byte, error)
	
	// GetImageStats returns image statistics
	GetImageStats() (map[string]interface{}, error)
	
	// ValidateUpload validates image upload before processing
	ValidateUpload(header *multipart.FileHeader) error
	
	// CleanupTempFiles removes temporary files
	CleanupTempFiles(olderThan time.Duration) error
	
	// GenerateThumbnail generates a thumbnail for an image
	GenerateThumbnail(imageID string, size int) ([]byte, error)
	
	// GetCacheStats returns cache statistics
	GetCacheStats() cache.ImageCacheStats
}

// ImageService implements ImageServiceInterface
type ImageService struct {
	imageRepo     repository.ImageRepository
	propertyRepo  repository.PropertyRepository
	storage       storage.ImageStorage
	processor     *processors.ImageProcessor
	cache         cache.ImageCacheInterface
	maxFileSize   int64
	maxImages     int
	allowedTypes  map[string]string
}

// NewImageService creates a new image service
func NewImageService(
	imageRepo repository.ImageRepository,
	propertyRepo repository.PropertyRepository,
	storage storage.ImageStorage,
	processor *processors.ImageProcessor,
	cache cache.ImageCacheInterface,
) *ImageService {
	return &ImageService{
		imageRepo:    imageRepo,
		propertyRepo: propertyRepo,
		storage:      storage,
		processor:    processor,
		cache:        cache,
		maxFileSize:  domain.MaxUploadSize,
		maxImages:    domain.MaxImagesPerProperty,
		allowedTypes: domain.SupportedMimeTypes,
	}
}

// Upload uploads and processes a new image
func (s *ImageService) Upload(propertyID string, file multipart.File, header *multipart.FileHeader, altText string) (*domain.ImageInfo, error) {
	// Validate property exists
	_, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		return nil, fmt.Errorf("property not found: %w", err)
	}
	
	// Validate upload
	if err := s.ValidateUpload(header); err != nil {
		return nil, fmt.Errorf("upload validation failed: %w", err)
	}
	
	// Check image limit
	count, err := s.imageRepo.GetImageCount(propertyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image count: %w", err)
	}
	
	if count >= s.maxImages {
		return nil, fmt.Errorf("maximum images per property exceeded: %d", s.maxImages)
	}
	
	// Read file data
	fileData := make([]byte, header.Size)
	if _, err := file.Read(fileData); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Validate image data
	if err := s.processor.ValidateImageData(fileData, s.maxFileSize); err != nil {
		return nil, fmt.Errorf("image validation failed: %w", err)
	}
	
	// Get original dimensions
	width, height, format, err := s.processor.GetImageDimensions(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to get image dimensions: %w", err)
	}
	
	// Create image info
	fileName := domain.GenerateImageFileName(propertyID, header.Filename)
	imageInfo := domain.NewImageInfo(propertyID, fileName)
	imageInfo.AltText = altText
	imageInfo.SortOrder = count // Add at the end
	
	// Process image for optimized storage
	optimizedData, stats, err := s.processor.OptimizeForSize(fileData, 1200) // 1.2MB target
	if err != nil {
		return nil, fmt.Errorf("failed to optimize image: %w", err)
	}
	
	// Store optimized image
	storedPath, err := s.storage.Store(optimizedData, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to store image: %w", err)
	}
	
	// Update image info with processing results
	imageInfo.OriginalURL = s.storage.GetURL(storedPath)
	imageInfo.SetProcessingResults(width, height, stats.OptimizedSize, format, 85, true)
	
	// Save to database
	if err := s.imageRepo.Create(imageInfo); err != nil {
		// Clean up stored file on database error
		s.storage.Delete(storedPath)
		return nil, fmt.Errorf("failed to save image metadata: %w", err)
	}
	
	log.Printf("Image uploaded successfully: %s, size: %d -> %d bytes (%.1f%% compression)",
		imageInfo.ID, stats.OriginalSize, stats.OptimizedSize, (1-stats.CompressionRatio)*100)
	
	return imageInfo, nil
}

// GetImage retrieves image metadata by ID
func (s *ImageService) GetImage(id string) (*domain.ImageInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("image ID cannot be empty")
	}
	
	return s.imageRepo.GetByID(id)
}

// GetImagesByProperty retrieves all images for a property
func (s *ImageService) GetImagesByProperty(propertyID string) ([]domain.ImageInfo, error) {
	if propertyID == "" {
		return nil, fmt.Errorf("property ID cannot be empty")
	}
	
	return s.imageRepo.GetByPropertyID(propertyID)
}

// UpdateImageMetadata updates image metadata
func (s *ImageService) UpdateImageMetadata(id, altText string, sortOrder int) error {
	if id == "" {
		return fmt.Errorf("image ID cannot be empty")
	}
	
	// Get existing image
	image, err := s.imageRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}
	
	// Update metadata
	image.UpdateMetadata(altText, sortOrder)
	
	// Save changes
	if err := s.imageRepo.Update(image); err != nil {
		return fmt.Errorf("failed to update image: %w", err)
	}
	
	return nil
}

// DeleteImage deletes image and its files
func (s *ImageService) DeleteImage(id string) error {
	if id == "" {
		return fmt.Errorf("image ID cannot be empty")
	}
	
	// Get image info
	image, err := s.imageRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}
	
	// Delete from storage
	// Extract path from URL
	storedPath := s.extractPathFromURL(image.OriginalURL)
	if storedPath != "" {
		if err := s.storage.Delete(storedPath); err != nil {
			log.Printf("Warning: failed to delete image file %s: %v", storedPath, err)
		}
	}
	
	// Delete variants and thumbnails
	s.deleteImageVariants(image.FileName)
	
	// Invalidate cache
	s.cache.InvalidateImage(id)
	
	// Delete from database
	if err := s.imageRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete image from database: %w", err)
	}
	
	log.Printf("Image deleted successfully: %s", id)
	return nil
}

// ReorderImages reorders images for a property
func (s *ImageService) ReorderImages(propertyID string, imageIDs []string) error {
	if propertyID == "" {
		return fmt.Errorf("property ID cannot be empty")
	}
	
	if len(imageIDs) == 0 {
		return fmt.Errorf("image IDs cannot be empty")
	}
	
	// Validate all images belong to the property
	for _, imageID := range imageIDs {
		image, err := s.imageRepo.GetByID(imageID)
		if err != nil {
			return fmt.Errorf("image not found: %s", imageID)
		}
		
		if image.PropertyID != propertyID {
			return fmt.Errorf("image %s does not belong to property %s", imageID, propertyID)
		}
	}
	
	// Update sort order
	return s.imageRepo.UpdateSortOrder(propertyID, imageIDs)
}

// SetMainImage sets an image as the main image for a property
func (s *ImageService) SetMainImage(propertyID, imageID string) error {
	if propertyID == "" {
		return fmt.Errorf("property ID cannot be empty")
	}
	
	if imageID == "" {
		return fmt.Errorf("image ID cannot be empty")
	}
	
	return s.imageRepo.SetMainImage(propertyID, imageID)
}

// GetMainImage gets the main image for a property
func (s *ImageService) GetMainImage(propertyID string) (*domain.ImageInfo, error) {
	if propertyID == "" {
		return nil, fmt.Errorf("property ID cannot be empty")
	}
	
	return s.imageRepo.GetMainImage(propertyID)
}

// GetImageVariant generates and returns an image variant
func (s *ImageService) GetImageVariant(imageID string, width, height int, format string, quality int) ([]byte, error) {
	if imageID == "" {
		return nil, fmt.Errorf("image ID cannot be empty")
	}
	
	// Check cache first
	if cachedData, _, found := s.cache.GetVariant(imageID, width, height, quality, format); found {
		return cachedData, nil
	}
	
	// Get image info
	image, err := s.imageRepo.GetByID(imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	
	// Generate variant filename
	variantName := fmt.Sprintf("%s_%dx%d_q%d.%s", 
		strings.TrimSuffix(image.FileName, filepath.Ext(image.FileName)),
		width, height, quality, format)
	
	// Check if variant already exists in storage
	variantPath := filepath.Join("variants", variantName)
	if s.storage.Exists(variantPath) {
		data, err := s.storage.Retrieve(variantPath)
		if err == nil {
			// Cache the retrieved data
			contentType := fmt.Sprintf("image/%s", format)
			if format == "jpg" {
				contentType = "image/jpeg"
			}
			s.cache.SetVariant(imageID, width, height, quality, format, data, contentType)
			return data, nil
		}
	}
	
	// Get original image data
	originalPath := s.extractPathFromURL(image.OriginalURL)
	originalData, err := s.storage.Retrieve(originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve original image: %w", err)
	}
	
	// Generate variant
	variantData, err := s.processor.GenerateImageVariant(originalData, width, height, quality, format)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image variant: %w", err)
	}
	
	// Store variant for future use
	if localStorage, ok := s.storage.(*storage.LocalImageStorage); ok {
		localStorage.StoreVariant(variantData, variantName, "variants")
	}
	
	// Cache the generated data
	contentType := fmt.Sprintf("image/%s", format)
	if format == "jpg" {
		contentType = "image/jpeg"
	}
	s.cache.SetVariant(imageID, width, height, quality, format, variantData, contentType)
	
	return variantData, nil
}

// GetImageStats returns image statistics
func (s *ImageService) GetImageStats() (map[string]interface{}, error) {
	stats, err := s.imageRepo.GetImageStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get image stats: %w", err)
	}
	
	// Add storage stats
	storageInfo := s.storage.GetStorageInfo()
	stats["storage_type"] = storageInfo.Type
	stats["storage_total_size"] = storageInfo.TotalSize
	stats["storage_file_count"] = storageInfo.FileCount
	
	// Add cache stats
	cacheStats := s.cache.Stats()
	stats["cache"] = cacheStats
	
	return stats, nil
}

// ValidateUpload validates image upload before processing
func (s *ImageService) ValidateUpload(header *multipart.FileHeader) error {
	if header == nil {
		return fmt.Errorf("file header cannot be nil")
	}
	
	if header.Size == 0 {
		return fmt.Errorf("file is empty")
	}
	
	if header.Size > s.maxFileSize {
		return fmt.Errorf("file too large: %d bytes, max: %d bytes", header.Size, s.maxFileSize)
	}
	
	// Check content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		return fmt.Errorf("content type not specified")
	}
	
	if !domain.IsSupportedMimeType(contentType) {
		return fmt.Errorf("unsupported content type: %s", contentType)
	}
	
	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	format := domain.GetImageFormatFromFilename(header.Filename)
	if format == "" {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
	
	return nil
}

// CleanupTempFiles removes temporary files
func (s *ImageService) CleanupTempFiles(olderThan time.Duration) error {
	if localStorage, ok := s.storage.(*storage.LocalImageStorage); ok {
		return localStorage.CleanupTempFiles(olderThan)
	}
	
	return nil // No cleanup needed for other storage types
}

// extractPathFromURL extracts the storage path from a URL
func (s *ImageService) extractPathFromURL(url string) string {
	if url == "" {
		return ""
	}
	
	// For local storage, URL format is: /path/to/file or baseURL/path/to/file
	// Extract the relative path part
	storageInfo := s.storage.GetStorageInfo()
	if storageInfo.BaseURL != "" {
		// Remove base URL prefix
		url = strings.TrimPrefix(url, storageInfo.BaseURL)
		url = strings.TrimPrefix(url, "/")
	}
	
	return url
}

// deleteImageVariants deletes all variants and thumbnails for an image
func (s *ImageService) deleteImageVariants(fileName string) {
	if fileName == "" {
		return
	}
	
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	
	// Delete thumbnails
	thumbnailPath := filepath.Join("thumbnails", baseName+"_thumb.jpg")
	s.storage.Delete(thumbnailPath)
	
	// Note: For a full implementation, you might want to:
	// 1. Keep track of generated variants in a cache/database
	// 2. Scan the variants directory for files matching the pattern
	// 3. Delete all matching files
	// For now, we'll just log that variants should be cleaned up
	log.Printf("TODO: Clean up variants for image: %s", fileName)
}

// GenerateThumbnail generates a thumbnail for an image
func (s *ImageService) GenerateThumbnail(imageID string, size int) ([]byte, error) {
	if imageID == "" {
		return nil, fmt.Errorf("image ID cannot be empty")
	}
	
	// Check cache first
	if cachedData, _, found := s.cache.GetThumbnail(imageID, size); found {
		return cachedData, nil
	}
	
	// Get image info
	image, err := s.imageRepo.GetByID(imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	
	// Generate thumbnail filename
	thumbnailName := fmt.Sprintf("%s_thumb_%d.jpg", 
		strings.TrimSuffix(image.FileName, filepath.Ext(image.FileName)), size)
	
	// Check if thumbnail already exists in storage
	thumbnailPath := filepath.Join("thumbnails", thumbnailName)
	if s.storage.Exists(thumbnailPath) {
		data, err := s.storage.Retrieve(thumbnailPath)
		if err == nil {
			// Cache the retrieved data
			s.cache.SetThumbnail(imageID, size, data, "image/jpeg")
			return data, nil
		}
	}
	
	// Get original image data
	originalPath := s.extractPathFromURL(image.OriginalURL)
	originalData, err := s.storage.Retrieve(originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve original image: %w", err)
	}
	
	// Generate thumbnail
	thumbnailData, err := s.processor.GenerateThumbnail(originalData, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate thumbnail: %w", err)
	}
	
	// Store thumbnail for future use
	if localStorage, ok := s.storage.(*storage.LocalImageStorage); ok {
		localStorage.StoreVariant(thumbnailData, thumbnailName, "thumbnails")
	}
	
	// Cache the generated data
	s.cache.SetThumbnail(imageID, size, thumbnailData, "image/jpeg")
	
	return thumbnailData, nil
}

// GetCacheStats returns cache statistics
func (s *ImageService) GetCacheStats() cache.ImageCacheStats {
	return s.cache.Stats()
}

// GetPaginatedImages gets paginated images (simplified version)
func (s *ImageService) GetPaginatedImages(pagination *domain.PaginationParams) ([]domain.ImageInfo, error) {
	if pagination == nil {
		pagination = domain.NewPaginationParams()
	}
	
	if err := pagination.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pagination parameters: %w", err)
	}
	
	// For simplicity, get all images and paginate manually
	// In a real implementation, this would be done at the database level
	// Get images using GetImagesByFormat with empty format to get all
	allImages, err := s.imageRepo.GetImagesByFormat("")
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}
	
	// Manual pagination
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	
	if offset >= len(allImages) {
		return []domain.ImageInfo{}, nil
	}
	
	end := offset + limit
	if end > len(allImages) {
		end = len(allImages)
	}
	
	return allImages[offset:end], nil
}

// CountImages returns the total count of images
func (s *ImageService) CountImages() (int, error) {
	// Use existing GetImagesByFormat and count the results
	allImages, err := s.imageRepo.GetImagesByFormat("")
	if err != nil {
		return 0, fmt.Errorf("failed to count images: %w", err)
	}
	return len(allImages), nil
}