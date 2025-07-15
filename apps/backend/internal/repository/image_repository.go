package repository

import (
	"database/sql"
	"fmt"
	"time"

	"realty-core/internal/domain"
)

// ImageRepository defines the interface for image metadata operations
type ImageRepository interface {
	// Create saves a new image record
	Create(image *domain.ImageInfo) error
	
	// GetByID retrieves image by ID
	GetByID(id string) (*domain.ImageInfo, error)
	
	// GetByPropertyID retrieves all images for a property
	GetByPropertyID(propertyID string) ([]domain.ImageInfo, error)
	
	// Update updates image metadata
	Update(image *domain.ImageInfo) error
	
	// Delete removes image record
	Delete(id string) error
	
	// UpdateSortOrder updates the sort order of images for a property
	UpdateSortOrder(propertyID string, imageIDs []string) error
	
	// GetMainImage gets the main image for a property
	GetMainImage(propertyID string) (*domain.ImageInfo, error)
	
	// SetMainImage sets an image as the main image for a property
	SetMainImage(propertyID, imageID string) error
	
	// GetImageCount returns the total number of images for a property
	GetImageCount(propertyID string) (int, error)
	
	// GetImagesByFormat retrieves images by format
	GetImagesByFormat(format string) ([]domain.ImageInfo, error)
	
	// GetImageStats returns image statistics
	GetImageStats() (map[string]interface{}, error)
}

// PostgreSQLImageRepository implements ImageRepository using PostgreSQL
type PostgreSQLImageRepository struct {
	db *sql.DB
}

// NewPostgreSQLImageRepository creates a new PostgreSQL image repository
func NewPostgreSQLImageRepository(db *sql.DB) *PostgreSQLImageRepository {
	return &PostgreSQLImageRepository{db: db}
}

// Create saves a new image record
func (r *PostgreSQLImageRepository) Create(image *domain.ImageInfo) error {
	if image == nil {
		return fmt.Errorf("image cannot be nil")
	}
	
	if err := image.Validate(); err != nil {
		return fmt.Errorf("image validation failed: %w", err)
	}
	
	query := `
		INSERT INTO images (
			id, property_id, file_name, original_url, alt_text, sort_order,
			size, width, height, format, quality, is_optimized, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)`
	
	_, err := r.db.Exec(query,
		image.ID, image.PropertyID, image.FileName, image.OriginalURL, image.AltText,
		image.SortOrder, image.Size, image.Width, image.Height, image.Format,
		image.Quality, image.IsOptimized, image.CreatedAt, image.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create image: %w", err)
	}
	
	return nil
}

// GetByID retrieves image by ID
func (r *PostgreSQLImageRepository) GetByID(id string) (*domain.ImageInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("image ID cannot be empty")
	}
	
	query := `
		SELECT id, property_id, file_name, original_url, alt_text, sort_order,
			   size, width, height, format, quality, is_optimized, created_at, updated_at
		FROM images
		WHERE id = $1`
	
	image := &domain.ImageInfo{}
	
	err := r.db.QueryRow(query, id).Scan(
		&image.ID, &image.PropertyID, &image.FileName, &image.OriginalURL,
		&image.AltText, &image.SortOrder, &image.Size, &image.Width,
		&image.Height, &image.Format, &image.Quality, &image.IsOptimized,
		&image.CreatedAt, &image.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("image not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	
	return image, nil
}

// GetByPropertyID retrieves all images for a property
func (r *PostgreSQLImageRepository) GetByPropertyID(propertyID string) ([]domain.ImageInfo, error) {
	if propertyID == "" {
		return nil, fmt.Errorf("property ID cannot be empty")
	}
	
	query := `
		SELECT id, property_id, file_name, original_url, alt_text, sort_order,
			   size, width, height, format, quality, is_optimized, created_at, updated_at
		FROM images
		WHERE property_id = $1
		ORDER BY sort_order ASC, created_at ASC`
	
	rows, err := r.db.Query(query, propertyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query images: %w", err)
	}
	defer rows.Close()
	
	var images []domain.ImageInfo
	
	for rows.Next() {
		var image domain.ImageInfo
		err := rows.Scan(
			&image.ID, &image.PropertyID, &image.FileName, &image.OriginalURL,
			&image.AltText, &image.SortOrder, &image.Size, &image.Width,
			&image.Height, &image.Format, &image.Quality, &image.IsOptimized,
			&image.CreatedAt, &image.UpdatedAt)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		
		images = append(images, image)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	
	return images, nil
}

// Update updates image metadata
func (r *PostgreSQLImageRepository) Update(image *domain.ImageInfo) error {
	if image == nil {
		return fmt.Errorf("image cannot be nil")
	}
	
	if err := image.Validate(); err != nil {
		return fmt.Errorf("image validation failed: %w", err)
	}
	
	image.UpdatedAt = time.Now()
	
	query := `
		UPDATE images SET
			file_name = $2, original_url = $3, alt_text = $4, sort_order = $5,
			size = $6, width = $7, height = $8, format = $9, quality = $10,
			is_optimized = $11, updated_at = $12
		WHERE id = $1`
	
	result, err := r.db.Exec(query,
		image.ID, image.FileName, image.OriginalURL, image.AltText,
		image.SortOrder, image.Size, image.Width, image.Height,
		image.Format, image.Quality, image.IsOptimized, image.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to update image: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("image not found: %s", image.ID)
	}
	
	return nil
}

// Delete removes image record
func (r *PostgreSQLImageRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("image ID cannot be empty")
	}
	
	query := `DELETE FROM images WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("image not found: %s", id)
	}
	
	return nil
}

// UpdateSortOrder updates the sort order of images for a property
func (r *PostgreSQLImageRepository) UpdateSortOrder(propertyID string, imageIDs []string) error {
	if propertyID == "" {
		return fmt.Errorf("property ID cannot be empty")
	}
	
	if len(imageIDs) == 0 {
		return fmt.Errorf("image IDs cannot be empty")
	}
	
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Update sort order for each image
	for i, imageID := range imageIDs {
		query := `UPDATE images SET sort_order = $1, updated_at = $2 WHERE id = $3 AND property_id = $4`
		
		result, err := tx.Exec(query, i, time.Now(), imageID, propertyID)
		if err != nil {
			return fmt.Errorf("failed to update sort order for image %s: %w", imageID, err)
		}
		
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected for image %s: %w", imageID, err)
		}
		
		if rowsAffected == 0 {
			return fmt.Errorf("image not found or belongs to different property: %s", imageID)
		}
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// GetMainImage gets the main image for a property
func (r *PostgreSQLImageRepository) GetMainImage(propertyID string) (*domain.ImageInfo, error) {
	if propertyID == "" {
		return nil, fmt.Errorf("property ID cannot be empty")
	}
	
	query := `
		SELECT id, property_id, file_name, original_url, alt_text, sort_order,
			   size, width, height, format, quality, is_optimized, created_at, updated_at
		FROM images
		WHERE property_id = $1
		ORDER BY sort_order ASC, created_at ASC
		LIMIT 1`
	
	image := &domain.ImageInfo{}
	
	err := r.db.QueryRow(query, propertyID).Scan(
		&image.ID, &image.PropertyID, &image.FileName, &image.OriginalURL,
		&image.AltText, &image.SortOrder, &image.Size, &image.Width,
		&image.Height, &image.Format, &image.Quality, &image.IsOptimized,
		&image.CreatedAt, &image.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no images found for property: %s", propertyID)
		}
		return nil, fmt.Errorf("failed to get main image: %w", err)
	}
	
	return image, nil
}

// SetMainImage sets an image as the main image for a property
func (r *PostgreSQLImageRepository) SetMainImage(propertyID, imageID string) error {
	if propertyID == "" {
		return fmt.Errorf("property ID cannot be empty")
	}
	
	if imageID == "" {
		return fmt.Errorf("image ID cannot be empty")
	}
	
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// First, verify the image exists and belongs to the property
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM images WHERE id = $1 AND property_id = $2)`
	err = tx.QueryRow(checkQuery, imageID, propertyID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check image existence: %w", err)
	}
	
	if !exists {
		return fmt.Errorf("image not found or belongs to different property: %s", imageID)
	}
	
	// Set the selected image as sort_order = 0
	updateQuery := `UPDATE images SET sort_order = 0, updated_at = $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, time.Now(), imageID)
	if err != nil {
		return fmt.Errorf("failed to set main image: %w", err)
	}
	
	// Update other images' sort_order to be > 0
	incrementQuery := `
		UPDATE images SET sort_order = sort_order + 1, updated_at = $1
		WHERE property_id = $2 AND id != $3`
	
	_, err = tx.Exec(incrementQuery, time.Now(), propertyID, imageID)
	if err != nil {
		return fmt.Errorf("failed to update other images sort order: %w", err)
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// GetImageCount returns the total number of images for a property
func (r *PostgreSQLImageRepository) GetImageCount(propertyID string) (int, error) {
	if propertyID == "" {
		return 0, fmt.Errorf("property ID cannot be empty")
	}
	
	query := `SELECT COUNT(*) FROM images WHERE property_id = $1`
	
	var count int
	err := r.db.QueryRow(query, propertyID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get image count: %w", err)
	}
	
	return count, nil
}

// GetImagesByFormat retrieves images by format
func (r *PostgreSQLImageRepository) GetImagesByFormat(format string) ([]domain.ImageInfo, error) {
	if format == "" {
		return nil, fmt.Errorf("format cannot be empty")
	}
	
	query := `
		SELECT id, property_id, file_name, original_url, alt_text, sort_order,
			   size, width, height, format, quality, is_optimized, created_at, updated_at
		FROM images
		WHERE format = $1
		ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, format)
	if err != nil {
		return nil, fmt.Errorf("failed to query images by format: %w", err)
	}
	defer rows.Close()
	
	var images []domain.ImageInfo
	
	for rows.Next() {
		var image domain.ImageInfo
		err := rows.Scan(
			&image.ID, &image.PropertyID, &image.FileName, &image.OriginalURL,
			&image.AltText, &image.SortOrder, &image.Size, &image.Width,
			&image.Height, &image.Format, &image.Quality, &image.IsOptimized,
			&image.CreatedAt, &image.UpdatedAt)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		
		images = append(images, image)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	
	return images, nil
}

// GetImageStats returns image statistics
func (r *PostgreSQLImageRepository) GetImageStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total images
	var totalImages int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM images`).Scan(&totalImages)
	if err != nil {
		return nil, fmt.Errorf("failed to get total images: %w", err)
	}
	stats["total_images"] = totalImages
	
	// Total size
	var totalSize int64
	err = r.db.QueryRow(`SELECT COALESCE(SUM(size), 0) FROM images`).Scan(&totalSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get total size: %w", err)
	}
	stats["total_size"] = totalSize
	
	// Average size
	var avgSize float64
	err = r.db.QueryRow(`SELECT COALESCE(AVG(size), 0) FROM images`).Scan(&avgSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get average size: %w", err)
	}
	stats["average_size"] = avgSize
	
	// Format distribution
	formatQuery := `SELECT format, COUNT(*) FROM images GROUP BY format ORDER BY COUNT(*) DESC`
	rows, err := r.db.Query(formatQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query format distribution: %w", err)
	}
	defer rows.Close()
	
	formats := make(map[string]int)
	for rows.Next() {
		var format string
		var count int
		if err := rows.Scan(&format, &count); err != nil {
			return nil, fmt.Errorf("failed to scan format distribution: %w", err)
		}
		formats[format] = count
	}
	stats["formats"] = formats
	
	// Properties with images
	var propertiesWithImages int
	err = r.db.QueryRow(`SELECT COUNT(DISTINCT property_id) FROM images`).Scan(&propertiesWithImages)
	if err != nil {
		return nil, fmt.Errorf("failed to get properties with images: %w", err)
	}
	stats["properties_with_images"] = propertiesWithImages
	
	// Average images per property
	var avgImagesPerProperty float64
	if propertiesWithImages > 0 {
		avgImagesPerProperty = float64(totalImages) / float64(propertiesWithImages)
	}
	stats["average_images_per_property"] = avgImagesPerProperty
	
	// Optimization rate
	var optimizedImages int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM images WHERE is_optimized = true`).Scan(&optimizedImages)
	if err != nil {
		return nil, fmt.Errorf("failed to get optimized images: %w", err)
	}
	stats["optimized_images"] = optimizedImages
	
	var optimizationRate float64
	if totalImages > 0 {
		optimizationRate = float64(optimizedImages) / float64(totalImages) * 100
	}
	stats["optimization_rate"] = optimizationRate
	
	return stats, nil
}

// CreateImageTable creates the images table if it doesn't exist
func (r *PostgreSQLImageRepository) CreateImageTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS images (
			id VARCHAR(36) PRIMARY KEY,
			property_id VARCHAR(36) NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			original_url TEXT NOT NULL,
			alt_text TEXT DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			size BIGINT DEFAULT 0,
			width INTEGER DEFAULT 0,
			height INTEGER DEFAULT 0,
			format VARCHAR(10) DEFAULT '',
			quality INTEGER DEFAULT 85,
			is_optimized BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
		);
		
		CREATE INDEX IF NOT EXISTS idx_images_property_id ON images(property_id);
		CREATE INDEX IF NOT EXISTS idx_images_sort_order ON images(property_id, sort_order);
		CREATE INDEX IF NOT EXISTS idx_images_format ON images(format);
		CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at);
	`
	
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create images table: %w", err)
	}
	
	return nil
}