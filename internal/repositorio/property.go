package repositorio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"realty-core/internal/domain"
	"strings"
)

// PropertyRepository handles property data access operations
type PropertyRepository struct {
	db *sql.DB
}

// NewPropertyRepository creates a new property repository instance
func NewPropertyRepository(db *sql.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

// Create creates a new property in the database
func (r *PropertyRepository) Create(property *domain.Property) error {
	// Serialize JSON fields
	imagesJSON, err := json.Marshal(property.Images)
	if err != nil {
		return fmt.Errorf("error marshaling images: %w", err)
	}

	tagsJSON, err := json.Marshal(property.Tags)
	if err != nil {
		return fmt.Errorf("error marshaling tags: %w", err)
	}

	query := `
		INSERT INTO properties (
			id, slug, title, description, price, province, city, sector, address,
			latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			main_image, images, video_tour, tour_360, rent_price, common_expenses, price_per_m2,
			year_built, floors, property_status, furnished, garage, pool, garden, terrace,
			balcony, security, elevator, air_conditioning, tags, featured, view_count,
			real_estate_company_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32,
			$33, $34, $35, $36, $37, $38, $39, $40, $41, $42
		)`

	_, err = r.db.Exec(query,
		property.ID, property.Slug, property.Title, property.Description, property.Price,
		property.Province, property.City, property.Sector, property.Address,
		nullFloat64(property.Latitude), nullFloat64(property.Longitude), property.LocationPrecision,
		property.Type, property.Status, property.Bedrooms, property.Bathrooms, property.AreaM2,
		nullString(property.MainImage), string(imagesJSON), nullString(property.VideoTour), nullString(property.Tour360),
		nullFloat64Ptr(property.RentPrice), nullFloat64Ptr(property.CommonExpenses), nullFloat64Ptr(property.PricePerM2),
		nullIntPtr(property.YearBuilt), nullIntPtr(property.Floors), property.PropertyStatus, property.Furnished,
		property.Garage, property.Pool, property.Garden, property.Terrace, property.Balcony,
		property.Security, property.Elevator, property.AirConditioning, string(tagsJSON),
		property.Featured, property.ViewCount, nullStringPtr(property.RealEstateCompanyID),
		property.CreatedAt, property.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating property: %w", err)
	}

	return nil
}

// GetByID retrieves a property by its ID
func (r *PropertyRepository) GetByID(id string) (*domain.Property, error) {
	query := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   main_image, images, video_tour, tour_360, rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished, garage, pool, garden, terrace,
			   balcony, security, elevator, air_conditioning, tags, featured, view_count,
			   real_estate_company_id, created_at, updated_at
		FROM properties WHERE id = $1`

	property := &domain.Property{}
	var imagesJSON, tagsJSON string
	var latitude, longitude sql.NullFloat64
	var mainImage, videoTour, tour360 sql.NullString
	var rentPrice, commonExpenses, pricePerM2 sql.NullFloat64
	var yearBuilt, floors sql.NullInt32
	var realEstateCompanyID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
		&property.Province, &property.City, &property.Sector, &property.Address,
		&latitude, &longitude, &property.LocationPrecision,
		&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
		&mainImage, &imagesJSON, &videoTour, &tour360, &rentPrice, &commonExpenses, &pricePerM2,
		&yearBuilt, &floors, &property.PropertyStatus, &property.Furnished,
		&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
		&property.Security, &property.Elevator, &property.AirConditioning, &tagsJSON,
		&property.Featured, &property.ViewCount, &realEstateCompanyID,
		&property.CreatedAt, &property.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("property not found with ID: %s", id)
		}
		return nil, fmt.Errorf("error retrieving property: %w", err)
	}

	// Handle nullable fields
	property.Latitude = latitude.Float64
	property.Longitude = longitude.Float64
	property.MainImage = mainImage.String
	property.VideoTour = videoTour.String
	property.Tour360 = tour360.String

	if rentPrice.Valid {
		property.RentPrice = &rentPrice.Float64
	}
	if commonExpenses.Valid {
		property.CommonExpenses = &commonExpenses.Float64
	}
	if pricePerM2.Valid {
		property.PricePerM2 = &pricePerM2.Float64
	}
	if yearBuilt.Valid {
		value := int(yearBuilt.Int32)
		property.YearBuilt = &value
	}
	if floors.Valid {
		value := int(floors.Int32)
		property.Floors = &value
	}
	if realEstateCompanyID.Valid {
		property.RealEstateCompanyID = &realEstateCompanyID.String
	}

	// Deserialize JSON fields
	if err := json.Unmarshal([]byte(imagesJSON), &property.Images); err != nil {
		return nil, fmt.Errorf("error unmarshaling images: %w", err)
	}
	if err := json.Unmarshal([]byte(tagsJSON), &property.Tags); err != nil {
		return nil, fmt.Errorf("error unmarshaling tags: %w", err)
	}

	return property, nil
}

// GetBySlug retrieves a property by its slug
func (r *PropertyRepository) GetBySlug(slug string) (*domain.Property, error) {
	query := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   main_image, images, video_tour, tour_360, rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished, garage, pool, garden, terrace,
			   balcony, security, elevator, air_conditioning, tags, featured, view_count,
			   real_estate_company_id, created_at, updated_at
		FROM properties WHERE slug = $1`

	property := &domain.Property{}
	var imagesJSON, tagsJSON string
	var latitude, longitude sql.NullFloat64
	var mainImage, videoTour, tour360 sql.NullString
	var rentPrice, commonExpenses, pricePerM2 sql.NullFloat64
	var yearBuilt, floors sql.NullInt32
	var realEstateCompanyID sql.NullString

	err := r.db.QueryRow(query, slug).Scan(
		&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
		&property.Province, &property.City, &property.Sector, &property.Address,
		&latitude, &longitude, &property.LocationPrecision,
		&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
		&mainImage, &imagesJSON, &videoTour, &tour360, &rentPrice, &commonExpenses, &pricePerM2,
		&yearBuilt, &floors, &property.PropertyStatus, &property.Furnished,
		&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
		&property.Security, &property.Elevator, &property.AirConditioning, &tagsJSON,
		&property.Featured, &property.ViewCount, &realEstateCompanyID,
		&property.CreatedAt, &property.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("property not found with slug: %s", slug)
		}
		return nil, fmt.Errorf("error retrieving property: %w", err)
	}

	// Handle nullable fields (same logic as GetByID)
	property.Latitude = latitude.Float64
	property.Longitude = longitude.Float64
	property.MainImage = mainImage.String
	property.VideoTour = videoTour.String
	property.Tour360 = tour360.String

	if rentPrice.Valid {
		property.RentPrice = &rentPrice.Float64
	}
	if commonExpenses.Valid {
		property.CommonExpenses = &commonExpenses.Float64
	}
	if pricePerM2.Valid {
		property.PricePerM2 = &pricePerM2.Float64
	}
	if yearBuilt.Valid {
		value := int(yearBuilt.Int32)
		property.YearBuilt = &value
	}
	if floors.Valid {
		value := int(floors.Int32)
		property.Floors = &value
	}
	if realEstateCompanyID.Valid {
		property.RealEstateCompanyID = &realEstateCompanyID.String
	}

	// Deserialize JSON fields
	if err := json.Unmarshal([]byte(imagesJSON), &property.Images); err != nil {
		return nil, fmt.Errorf("error unmarshaling images: %w", err)
	}
	if err := json.Unmarshal([]byte(tagsJSON), &property.Tags); err != nil {
		return nil, fmt.Errorf("error unmarshaling tags: %w", err)
	}

	return property, nil
}

// GetAll retrieves all properties with optional filters
func (r *PropertyRepository) GetAll(filters map[string]interface{}) ([]*domain.Property, error) {
	baseQuery := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   main_image, images, video_tour, tour_360, rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished, garage, pool, garden, terrace,
			   balcony, security, elevator, air_conditioning, tags, featured, view_count,
			   real_estate_company_id, created_at, updated_at
		FROM properties`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Apply filters
	if province, ok := filters["province"].(string); ok && province != "" {
		conditions = append(conditions, fmt.Sprintf("province = $%d", argIndex))
		args = append(args, province)
		argIndex++
	}

	if city, ok := filters["city"].(string); ok && city != "" {
		conditions = append(conditions, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, city)
		argIndex++
	}

	if propertyType, ok := filters["type"].(string); ok && propertyType != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, propertyType)
		argIndex++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	if minPrice, ok := filters["min_price"].(float64); ok && minPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, minPrice)
		argIndex++
	}

	if maxPrice, ok := filters["max_price"].(float64); ok && maxPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, maxPrice)
		argIndex++
	}

	if bedrooms, ok := filters["bedrooms"].(int); ok && bedrooms > 0 {
		conditions = append(conditions, fmt.Sprintf("bedrooms >= $%d", argIndex))
		args = append(args, bedrooms)
		argIndex++
	}

	if featured, ok := filters["featured"].(bool); ok && featured {
		conditions = append(conditions, "featured = true")
	}

	// Build final query
	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY featured DESC, created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying properties: %w", err)
	}
	defer rows.Close()

	var properties []*domain.Property
	for rows.Next() {
		property := &domain.Property{}
		var imagesJSON, tagsJSON string
		var latitude, longitude sql.NullFloat64
		var mainImage, videoTour, tour360 sql.NullString
		var rentPrice, commonExpenses, pricePerM2 sql.NullFloat64
		var yearBuilt, floors sql.NullInt32
		var realEstateCompanyID sql.NullString

		err := rows.Scan(
			&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
			&property.Province, &property.City, &property.Sector, &property.Address,
			&latitude, &longitude, &property.LocationPrecision,
			&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
			&mainImage, &imagesJSON, &videoTour, &tour360, &rentPrice, &commonExpenses, &pricePerM2,
			&yearBuilt, &floors, &property.PropertyStatus, &property.Furnished,
			&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
			&property.Security, &property.Elevator, &property.AirConditioning, &tagsJSON,
			&property.Featured, &property.ViewCount, &realEstateCompanyID,
			&property.CreatedAt, &property.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning property: %w", err)
		}

		// Handle nullable fields
		property.Latitude = latitude.Float64
		property.Longitude = longitude.Float64
		property.MainImage = mainImage.String
		property.VideoTour = videoTour.String
		property.Tour360 = tour360.String

		if rentPrice.Valid {
			property.RentPrice = &rentPrice.Float64
		}
		if commonExpenses.Valid {
			property.CommonExpenses = &commonExpenses.Float64
		}
		if pricePerM2.Valid {
			property.PricePerM2 = &pricePerM2.Float64
		}
		if yearBuilt.Valid {
			value := int(yearBuilt.Int32)
			property.YearBuilt = &value
		}
		if floors.Valid {
			value := int(floors.Int32)
			property.Floors = &value
		}
		if realEstateCompanyID.Valid {
			property.RealEstateCompanyID = &realEstateCompanyID.String
		}

		// Deserialize JSON fields
		if err := json.Unmarshal([]byte(imagesJSON), &property.Images); err != nil {
			return nil, fmt.Errorf("error unmarshaling images: %w", err)
		}
		if err := json.Unmarshal([]byte(tagsJSON), &property.Tags); err != nil {
			return nil, fmt.Errorf("error unmarshaling tags: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}

// GetAvailable retrieves all available properties
func (r *PropertyRepository) GetAvailable() ([]*domain.Property, error) {
	filters := map[string]interface{}{
		"status": domain.PropertyStatusAvailable,
	}
	return r.GetAll(filters)
}

// GetFeatured retrieves all featured properties
func (r *PropertyRepository) GetFeatured() ([]*domain.Property, error) {
	filters := map[string]interface{}{
		"status":   domain.PropertyStatusAvailable,
		"featured": true,
	}
	return r.GetAll(filters)
}

// Update updates an existing property
func (r *PropertyRepository) Update(property *domain.Property) error {
	// Update timestamp
	property.UpdateTimestamp()

	// Serialize JSON fields
	imagesJSON, err := json.Marshal(property.Images)
	if err != nil {
		return fmt.Errorf("error marshaling images: %w", err)
	}

	tagsJSON, err := json.Marshal(property.Tags)
	if err != nil {
		return fmt.Errorf("error marshaling tags: %w", err)
	}

	query := `
		UPDATE properties SET
			slug = $2, title = $3, description = $4, price = $5, province = $6, city = $7,
			sector = $8, address = $9, latitude = $10, longitude = $11, location_precision = $12,
			type = $13, status = $14, bedrooms = $15, bathrooms = $16, area_m2 = $17,
			main_image = $18, images = $19, video_tour = $20, tour_360 = $21,
			rent_price = $22, common_expenses = $23, price_per_m2 = $24, year_built = $25,
			floors = $26, property_status = $27, furnished = $28, garage = $29, pool = $30,
			garden = $31, terrace = $32, balcony = $33, security = $34, elevator = $35,
			air_conditioning = $36, tags = $37, featured = $38, view_count = $39,
			real_estate_company_id = $40, updated_at = $41
		WHERE id = $1`

	_, err = r.db.Exec(query,
		property.ID, property.Slug, property.Title, property.Description, property.Price,
		property.Province, property.City, property.Sector, property.Address,
		nullFloat64(property.Latitude), nullFloat64(property.Longitude), property.LocationPrecision,
		property.Type, property.Status, property.Bedrooms, property.Bathrooms, property.AreaM2,
		nullString(property.MainImage), string(imagesJSON), nullString(property.VideoTour), nullString(property.Tour360),
		nullFloat64Ptr(property.RentPrice), nullFloat64Ptr(property.CommonExpenses), nullFloat64Ptr(property.PricePerM2),
		nullIntPtr(property.YearBuilt), nullIntPtr(property.Floors), property.PropertyStatus, property.Furnished,
		property.Garage, property.Pool, property.Garden, property.Terrace, property.Balcony,
		property.Security, property.Elevator, property.AirConditioning, string(tagsJSON),
		property.Featured, property.ViewCount, nullStringPtr(property.RealEstateCompanyID),
		property.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating property: %w", err)
	}

	return nil
}

// Delete deletes a property by ID
func (r *PropertyRepository) Delete(id string) error {
	query := `DELETE FROM properties WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting property: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("property not found with ID: %s", id)
	}

	return nil
}

// Search performs full-text search on properties
func (r *PropertyRepository) Search(searchTerm string) ([]*domain.Property, error) {
	query := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   main_image, images, video_tour, tour_360, rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished, garage, pool, garden, terrace,
			   balcony, security, elevator, air_conditioning, tags, featured, view_count,
			   real_estate_company_id, created_at, updated_at
		FROM properties 
		WHERE status = 'available'
		  AND to_tsvector('english', title || ' ' || COALESCE(description, '')) @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(to_tsvector('english', title || ' ' || COALESCE(description, '')), plainto_tsquery('english', $1)) DESC,
				 featured DESC, created_at DESC`

	rows, err := r.db.Query(query, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("error searching properties: %w", err)
	}
	defer rows.Close()

	var properties []*domain.Property
	for rows.Next() {
		property := &domain.Property{}
		var imagesJSON, tagsJSON string
		var latitude, longitude sql.NullFloat64
		var mainImage, videoTour, tour360 sql.NullString
		var rentPrice, commonExpenses, pricePerM2 sql.NullFloat64
		var yearBuilt, floors sql.NullInt32
		var realEstateCompanyID sql.NullString

		err := rows.Scan(
			&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
			&property.Province, &property.City, &property.Sector, &property.Address,
			&latitude, &longitude, &property.LocationPrecision,
			&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
			&mainImage, &imagesJSON, &videoTour, &tour360, &rentPrice, &commonExpenses, &pricePerM2,
			&yearBuilt, &floors, &property.PropertyStatus, &property.Furnished,
			&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
			&property.Security, &property.Elevator, &property.AirConditioning, &tagsJSON,
			&property.Featured, &property.ViewCount, &realEstateCompanyID,
			&property.CreatedAt, &property.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning property: %w", err)
		}

		// Handle nullable fields (same logic as GetByID)
		property.Latitude = latitude.Float64
		property.Longitude = longitude.Float64
		property.MainImage = mainImage.String
		property.VideoTour = videoTour.String
		property.Tour360 = tour360.String

		if rentPrice.Valid {
			property.RentPrice = &rentPrice.Float64
		}
		if commonExpenses.Valid {
			property.CommonExpenses = &commonExpenses.Float64
		}
		if pricePerM2.Valid {
			property.PricePerM2 = &pricePerM2.Float64
		}
		if yearBuilt.Valid {
			value := int(yearBuilt.Int32)
			property.YearBuilt = &value
		}
		if floors.Valid {
			value := int(floors.Int32)
			property.Floors = &value
		}
		if realEstateCompanyID.Valid {
			property.RealEstateCompanyID = &realEstateCompanyID.String
		}

		// Deserialize JSON fields
		if err := json.Unmarshal([]byte(imagesJSON), &property.Images); err != nil {
			return nil, fmt.Errorf("error unmarshaling images: %w", err)
		}
		if err := json.Unmarshal([]byte(tagsJSON), &property.Tags); err != nil {
			return nil, fmt.Errorf("error unmarshaling tags: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}

// IncrementViewCount increments the view count for a property
func (r *PropertyRepository) IncrementViewCount(id string) error {
	query := `UPDATE properties SET view_count = view_count + 1 WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error incrementing view count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("property not found with ID: %s", id)
	}

	return nil
}

// GetByCompany retrieves all properties for a specific real estate company
func (r *PropertyRepository) GetByCompany(companyID string) ([]*domain.Property, error) {
	filters := map[string]interface{}{
		"real_estate_company_id": companyID,
	}
	return r.GetAll(filters)
}

// Helper functions for handling null values
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullStringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func nullFloat64(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

func nullFloat64Ptr(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

func nullIntPtr(i *int) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*i), Valid: true}
}
