package service

import (
	"fmt"
	"log"

	"realty-core/internal/domain"
)

// enrichPropertyWithImages enriches a single property with image data
func (s *PropertyService) enrichPropertyWithImages(property *domain.Property) {
	if property == nil || s.imageRepo == nil {
		return
	}

	// Get images for the property
	images, err := s.imageRepo.GetByPropertyID(property.ID)
	if err != nil {
		log.Printf("Warning: failed to get images for property %s: %v", property.ID, err)
		return
	}

	// Convert images to URLs array for backward compatibility
	var imageURLs []string
	for _, img := range images {
		imageURLs = append(imageURLs, img.OriginalURL)
	}

	// Set images array
	property.Images = imageURLs

	// Set main image (first image in sort order)
	if len(images) > 0 {
		property.MainImage = &images[0].OriginalURL
	}
}

// enrichPropertiesWithImages enriches multiple properties with image data
func (s *PropertyService) enrichPropertiesWithImages(properties []domain.Property) {
	if s.imageRepo == nil {
		return
	}

	for i := range properties {
		s.enrichPropertyWithImages(&properties[i])
	}
}

// GetPropertyImages returns images for a property
func (s *PropertyService) GetPropertyImages(propertyID string) ([]domain.ImageInfo, error) {
	if propertyID == "" {
		return nil, fmt.Errorf("property ID required")
	}

	// Verify property exists
	_, err := s.repo.GetByID(propertyID)
	if err != nil {
		return nil, fmt.Errorf("property not found: %w", err)
	}

	// Get images
	images, err := s.imageRepo.GetByPropertyID(propertyID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving property images: %w", err)
	}

	return images, nil
}

// GetPropertyMainImage returns the main image for a property
func (s *PropertyService) GetPropertyMainImage(propertyID string) (*domain.ImageInfo, error) {
	if propertyID == "" {
		return nil, fmt.Errorf("property ID required")
	}

	// Verify property exists
	_, err := s.repo.GetByID(propertyID)
	if err != nil {
		return nil, fmt.Errorf("property not found: %w", err)
	}

	// Get main image
	mainImage, err := s.imageRepo.GetMainImage(propertyID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving main image: %w", err)
	}

	return mainImage, nil
}

// GetPropertyImageCount returns the number of images for a property
func (s *PropertyService) GetPropertyImageCount(propertyID string) (int, error) {
	if propertyID == "" {
		return 0, fmt.Errorf("property ID required")
	}

	// Verify property exists
	_, err := s.repo.GetByID(propertyID)
	if err != nil {
		return 0, fmt.Errorf("property not found: %w", err)
	}

	// Get image count
	count, err := s.imageRepo.GetImageCount(propertyID)
	if err != nil {
		return 0, fmt.Errorf("error retrieving image count: %w", err)
	}

	return count, nil
}

// HasImages checks if a property has images
func (s *PropertyService) HasImages(propertyID string) (bool, error) {
	count, err := s.GetPropertyImageCount(propertyID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPropertiesWithImages returns properties that have images
func (s *PropertyService) GetPropertiesWithImages() ([]domain.Property, error) {
	if s.imageRepo == nil {
		return nil, fmt.Errorf("image repository not available")
	}

	// Get all properties
	properties, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving properties: %w", err)
	}

	// Filter properties with images
	var propertiesWithImages []domain.Property
	for _, property := range properties {
		count, err := s.imageRepo.GetImageCount(property.ID)
		if err != nil {
			log.Printf("Warning: failed to get image count for property %s: %v", property.ID, err)
			continue
		}

		if count > 0 {
			// Enrich with image data
			s.enrichPropertyWithImages(&property)
			propertiesWithImages = append(propertiesWithImages, property)
		}
	}

	return propertiesWithImages, nil
}

// GetPropertiesWithoutImages returns properties that don't have images
func (s *PropertyService) GetPropertiesWithoutImages() ([]domain.Property, error) {
	if s.imageRepo == nil {
		return nil, fmt.Errorf("image repository not available")
	}

	// Get all properties
	properties, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving properties: %w", err)
	}

	// Filter properties without images
	var propertiesWithoutImages []domain.Property
	for _, property := range properties {
		count, err := s.imageRepo.GetImageCount(property.ID)
		if err != nil {
			log.Printf("Warning: failed to get image count for property %s: %v", property.ID, err)
			continue
		}

		if count == 0 {
			propertiesWithoutImages = append(propertiesWithoutImages, property)
		}
	}

	return propertiesWithoutImages, nil
}

// GetExtendedStatistics returns extended property statistics including image data
func (s *PropertyService) GetExtendedStatistics() (map[string]interface{}, error) {
	// Get basic statistics
	stats, err := s.GetStatistics()
	if err != nil {
		return nil, err
	}

	// Add image statistics if image repo is available
	if s.imageRepo != nil {
		imageStats, err := s.imageRepo.GetImageStats()
		if err != nil {
			log.Printf("Warning: failed to get image statistics: %v", err)
		} else {
			stats["images"] = imageStats
		}

		// Count properties with/without images
		propertiesWithImages, err := s.GetPropertiesWithImages()
		if err != nil {
			log.Printf("Warning: failed to get properties with images: %v", err)
		} else {
			stats["properties_with_images"] = len(propertiesWithImages)
		}

		propertiesWithoutImages, err := s.GetPropertiesWithoutImages()
		if err != nil {
			log.Printf("Warning: failed to get properties without images: %v", err)
		} else {
			stats["properties_without_images"] = len(propertiesWithoutImages)
		}
	}

	return stats, nil
}