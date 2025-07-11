package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"realty-core/internal/domain"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PropertyRepository defines the data access operations for properties
type PropertyRepository interface {
	Create(property *domain.Property) error
	GetByID(id string) (*domain.Property, error)
	GetBySlug(slug string) (*domain.Property, error)
	GetAll() ([]domain.Property, error)
	Update(property *domain.Property) error
	Delete(id string) error
	GetByProvince(province string) ([]domain.Property, error)
	GetByPriceRange(minPrice, maxPrice float64) ([]domain.Property, error)
	// Full-text search methods
	SearchProperties(query string, limit int) ([]domain.Property, error)
	SearchPropertiesRanked(query string, limit int) ([]PropertySearchResult, error)
	GetSearchSuggestions(query string, limit int) ([]SearchSuggestion, error)
	AdvancedSearch(params AdvancedSearchParams) ([]PropertySearchResult, error)
	// Pagination methods
	GetAllPaginated(pagination *domain.PaginationParams) ([]domain.Property, int, error)
	GetByProvincePaginated(province string, pagination *domain.PaginationParams) ([]domain.Property, int, error)
	GetByPriceRangePaginated(minPrice, maxPrice float64, pagination *domain.PaginationParams) ([]domain.Property, int, error)
	SearchPropertiesPaginated(query string, pagination *domain.PaginationParams) ([]domain.Property, int, error)
	SearchPropertiesRankedPaginated(query string, pagination *domain.PaginationParams) ([]PropertySearchResult, int, error)
	AdvancedSearchPaginated(params AdvancedSearchParams, pagination *domain.PaginationParams) ([]PropertySearchResult, int, error)
}

// PropertySearchResult represents a search result with ranking
type PropertySearchResult struct {
	Property domain.Property
	Rank     float64
}

// SearchSuggestion represents a search suggestion
type SearchSuggestion struct {
	Text      string
	Category  string
	Frequency int
}

// AdvancedSearchParams holds parameters for advanced search
type AdvancedSearchParams struct {
	Query        string
	Province     string
	City         string
	Type         string
	MinPrice     float64
	MaxPrice     float64
	MinBedrooms  int
	MaxBedrooms  int
	MinBathrooms float64
	MaxBathrooms float64
	MinArea      float64
	MaxArea      float64
	FeaturedOnly bool
	Limit        int
}

// PostgreSQLPropertyRepository implements PropertyRepository using PostgreSQL
type PostgreSQLPropertyRepository struct {
	db *sql.DB
}

// NewPostgreSQLPropertyRepository creates a new instance of the repository
func NewPostgreSQLPropertyRepository(db *sql.DB) *PostgreSQLPropertyRepository {
	return &PostgreSQLPropertyRepository{db: db}
}

// Create inserts a new property into the database
func (r *PostgreSQLPropertyRepository) Create(property *domain.Property) error {
	// Convert slices to JSON for storage in JSONB
	imagesJSON, err := json.Marshal(property.Images)
	if err != nil {
		return fmt.Errorf("error converting images to JSON: %w", err)
	}

	tagsJSON, err := json.Marshal(property.Tags)
	if err != nil {
		return fmt.Errorf("error converting tags to JSON: %w", err)
	}

	query := `
		INSERT INTO properties (
			id, slug, title, description, price, province, city, sector, address,
			latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			parking_spaces, main_image, images, video_tour, tour_360,
			rent_price, common_expenses, price_per_m2,
			year_built, floors, property_status, furnished,
			garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			tags, featured, view_count, real_estate_company_id,
			owner_id, agent_id, agency_id, created_by, updated_by,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32,
			$33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48
		)
	`

	_, err = r.db.Exec(
		query,
		property.ID, property.Slug, property.Title, property.Description, property.Price,
		property.Province, property.City, property.Sector, property.Address,
		property.Latitude, property.Longitude, property.LocationPrecision,
		property.Type, property.Status, property.Bedrooms, property.Bathrooms, property.AreaM2,
		property.ParkingSpaces, property.MainImage, string(imagesJSON), property.VideoTour, property.Tour360,
		property.RentPrice, property.CommonExpenses, property.PricePerM2,
		property.YearBuilt, property.Floors, property.PropertyStatus, property.Furnished,
		property.Garage, property.Pool, property.Garden, property.Terrace, property.Balcony,
		property.Security, property.Elevator, property.AirConditioning,
		string(tagsJSON), property.Featured, property.ViewCount, property.RealEstateCompanyID,
		property.OwnerID, property.AgentID, property.AgencyID, property.CreatedBy, property.UpdatedBy,
		property.CreatedAt, property.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating property: %w", err)
	}

	log.Printf("Property created successfully: %s", property.ID)
	return nil
}

// GetByID retrieves a property by its ID
func (r *PostgreSQLPropertyRepository) GetByID(id string) (*domain.Property, error) {
	query := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		WHERE id = $1
	`

	var property domain.Property
	var imagesJSON, tagsJSON string

	err := r.db.QueryRow(query, id).Scan(
		&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
		&property.Province, &property.City, &property.Sector, &property.Address,
		&property.Latitude, &property.Longitude, &property.LocationPrecision,
		&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
		&property.ParkingSpaces, &property.MainImage, &imagesJSON, &property.VideoTour, &property.Tour360,
		&property.RentPrice, &property.CommonExpenses, &property.PricePerM2,
		&property.YearBuilt, &property.Floors, &property.PropertyStatus, &property.Furnished,
		&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
		&property.Security, &property.Elevator, &property.AirConditioning,
		&tagsJSON, &property.Featured, &property.ViewCount, &property.RealEstateCompanyID,
		&property.OwnerID, &property.AgentID, &property.AgencyID, &property.CreatedBy, &property.UpdatedBy,
		&property.CreatedAt, &property.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("property not found: %s", id)
		}
		return nil, fmt.Errorf("error retrieving property: %w", err)
	}

	// Convert JSON back to slices
	if imagesJSON != "" {
		err = json.Unmarshal([]byte(imagesJSON), &property.Images)
		if err != nil {
			return nil, fmt.Errorf("error converting images from JSON: %w", err)
		}
	}

	if tagsJSON != "" {
		err = json.Unmarshal([]byte(tagsJSON), &property.Tags)
		if err != nil {
			return nil, fmt.Errorf("error converting tags from JSON: %w", err)
		}
	}

	return &property, nil
}

// GetBySlug retrieves a property by its SEO slug
func (r *PostgreSQLPropertyRepository) GetBySlug(slug string) (*domain.Property, error) {
	query := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		WHERE slug = $1
	`

	var property domain.Property
	var imagesJSON, tagsJSON string

	err := r.db.QueryRow(query, slug).Scan(
		&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
		&property.Province, &property.City, &property.Sector, &property.Address,
		&property.Latitude, &property.Longitude, &property.LocationPrecision,
		&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
		&property.ParkingSpaces, &property.MainImage, &imagesJSON, &property.VideoTour, &property.Tour360,
		&property.RentPrice, &property.CommonExpenses, &property.PricePerM2,
		&property.YearBuilt, &property.Floors, &property.PropertyStatus, &property.Furnished,
		&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
		&property.Security, &property.Elevator, &property.AirConditioning,
		&tagsJSON, &property.Featured, &property.ViewCount, &property.RealEstateCompanyID,
		&property.OwnerID, &property.AgentID, &property.AgencyID, &property.CreatedBy, &property.UpdatedBy,
		&property.CreatedAt, &property.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("property not found with slug: %s", slug)
		}
		return nil, fmt.Errorf("error retrieving property by slug: %w", err)
	}

	// Convert JSON back to slices
	if imagesJSON != "" {
		err = json.Unmarshal([]byte(imagesJSON), &property.Images)
		if err != nil {
			return nil, fmt.Errorf("error converting images from JSON: %w", err)
		}
	}

	if tagsJSON != "" {
		err = json.Unmarshal([]byte(tagsJSON), &property.Tags)
		if err != nil {
			return nil, fmt.Errorf("error converting tags from JSON: %w", err)
		}
	}

	return &property, nil
}

// GetAll returns all properties (with pagination in a real implementation)
func (r *PostgreSQLPropertyRepository) GetAll() ([]domain.Property, error) {
	query := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		ORDER BY featured DESC, created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying properties: %w", err)
	}
	defer rows.Close()

	var properties []domain.Property

	for rows.Next() {
		var property domain.Property
		var imagesJSON, tagsJSON string

		err := rows.Scan(
			&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
			&property.Province, &property.City, &property.Sector, &property.Address,
			&property.Latitude, &property.Longitude, &property.LocationPrecision,
			&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
			&property.ParkingSpaces, &property.MainImage, &imagesJSON, &property.VideoTour, &property.Tour360,
			&property.RentPrice, &property.CommonExpenses, &property.PricePerM2,
			&property.YearBuilt, &property.Floors, &property.PropertyStatus, &property.Furnished,
			&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
			&property.Security, &property.Elevator, &property.AirConditioning,
			&tagsJSON, &property.Featured, &property.ViewCount, &property.RealEstateCompanyID,
			&property.OwnerID, &property.AgentID, &property.AgencyID, &property.CreatedBy, &property.UpdatedBy,
			&property.CreatedAt, &property.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning property: %w", err)
		}

		// Convert JSON back to slices
		if imagesJSON != "" {
			err = json.Unmarshal([]byte(imagesJSON), &property.Images)
			if err != nil {
				property.Images = []string{} // Continue with empty slice if JSON is invalid
			}
		}

		if tagsJSON != "" {
			err = json.Unmarshal([]byte(tagsJSON), &property.Tags)
			if err != nil {
				property.Tags = []string{} // Continue with empty slice if JSON is invalid
			}
		}

		properties = append(properties, property)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating properties: %w", err)
	}

	return properties, nil
}

// Update modifies an existing property
func (r *PostgreSQLPropertyRepository) Update(property *domain.Property) error {
	property.UpdateTimestamp()
	property.UpdateSlug()

	// Convert slices to JSON
	imagesJSON, err := json.Marshal(property.Images)
	if err != nil {
		return fmt.Errorf("error converting images to JSON: %w", err)
	}

	tagsJSON, err := json.Marshal(property.Tags)
	if err != nil {
		return fmt.Errorf("error converting tags to JSON: %w", err)
	}

	query := `
		UPDATE properties SET 
			slug = $2, title = $3, description = $4, price = $5, province = $6, city = $7,
			sector = $8, address = $9, latitude = $10, longitude = $11, location_precision = $12,
			type = $13, status = $14, bedrooms = $15, bathrooms = $16, area_m2 = $17,
			parking_spaces = $18, main_image = $19, images = $20, video_tour = $21, tour_360 = $22,
			rent_price = $23, common_expenses = $24, price_per_m2 = $25,
			year_built = $26, floors = $27, property_status = $28, furnished = $29,
			garage = $30, pool = $31, garden = $32, terrace = $33, balcony = $34,
			security = $35, elevator = $36, air_conditioning = $37,
			tags = $38, featured = $39, view_count = $40, real_estate_company_id = $41,
			updated_at = $42
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		property.ID, property.Slug, property.Title, property.Description, property.Price,
		property.Province, property.City, property.Sector, property.Address,
		property.Latitude, property.Longitude, property.LocationPrecision,
		property.Type, property.Status, property.Bedrooms, property.Bathrooms, property.AreaM2,
		property.ParkingSpaces, property.MainImage, string(imagesJSON), property.VideoTour, property.Tour360,
		property.RentPrice, property.CommonExpenses, property.PricePerM2,
		property.YearBuilt, property.Floors, property.PropertyStatus, property.Furnished,
		property.Garage, property.Pool, property.Garden, property.Terrace, property.Balcony,
		property.Security, property.Elevator, property.AirConditioning,
		string(tagsJSON), property.Featured, property.ViewCount, property.RealEstateCompanyID,
		property.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating property: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("property not found: %s", property.ID)
	}

	log.Printf("Property updated successfully: %s", property.ID)
	return nil
}

// Delete removes a property from the database
func (r *PostgreSQLPropertyRepository) Delete(id string) error {
	query := `DELETE FROM properties WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting property: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking delete result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("property not found: %s", id)
	}

	log.Printf("Property deleted successfully: %s", id)
	return nil
}

// GetByProvince filters properties by province
func (r *PostgreSQLPropertyRepository) GetByProvince(province string) ([]domain.Property, error) {
	query := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		WHERE province = $1
		ORDER BY featured DESC, created_at DESC
	`

	rows, err := r.db.Query(query, province)
	if err != nil {
		return nil, fmt.Errorf("error querying properties by province: %w", err)
	}
	defer rows.Close()

	var properties []domain.Property

	for rows.Next() {
		var property domain.Property
		var imagesJSON, tagsJSON string

		err := rows.Scan(
			&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
			&property.Province, &property.City, &property.Sector, &property.Address,
			&property.Latitude, &property.Longitude, &property.LocationPrecision,
			&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
			&property.ParkingSpaces, &property.MainImage, &imagesJSON, &property.VideoTour, &property.Tour360,
			&property.RentPrice, &property.CommonExpenses, &property.PricePerM2,
			&property.YearBuilt, &property.Floors, &property.PropertyStatus, &property.Furnished,
			&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
			&property.Security, &property.Elevator, &property.AirConditioning,
			&tagsJSON, &property.Featured, &property.ViewCount, &property.RealEstateCompanyID,
			&property.OwnerID, &property.AgentID, &property.AgencyID, &property.CreatedBy, &property.UpdatedBy,
			&property.CreatedAt, &property.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning property: %w", err)
		}

		// Convert JSON back to slices
		if imagesJSON != "" {
			json.Unmarshal([]byte(imagesJSON), &property.Images)
		}
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &property.Tags)
		}

		properties = append(properties, property)
	}

	return properties, nil
}

// GetByPriceRange filters properties by price range
func (r *PostgreSQLPropertyRepository) GetByPriceRange(minPrice, maxPrice float64) ([]domain.Property, error) {
	query := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		WHERE price >= $1 AND price <= $2
		ORDER BY featured DESC, created_at DESC
	`

	rows, err := r.db.Query(query, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("error querying properties by price range: %w", err)
	}
	defer rows.Close()

	var properties []domain.Property

	for rows.Next() {
		var property domain.Property
		var imagesJSON, tagsJSON string

		err := rows.Scan(
			&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
			&property.Province, &property.City, &property.Sector, &property.Address,
			&property.Latitude, &property.Longitude, &property.LocationPrecision,
			&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
			&property.ParkingSpaces, &property.MainImage, &imagesJSON, &property.VideoTour, &property.Tour360,
			&property.RentPrice, &property.CommonExpenses, &property.PricePerM2,
			&property.YearBuilt, &property.Floors, &property.PropertyStatus, &property.Furnished,
			&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
			&property.Security, &property.Elevator, &property.AirConditioning,
			&tagsJSON, &property.Featured, &property.ViewCount, &property.RealEstateCompanyID,
			&property.OwnerID, &property.AgentID, &property.AgencyID, &property.CreatedBy, &property.UpdatedBy,
			&property.CreatedAt, &property.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning property: %w", err)
		}

		// Convert JSON back to slices
		if imagesJSON != "" {
			json.Unmarshal([]byte(imagesJSON), &property.Images)
		}
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &property.Tags)
		}

		properties = append(properties, property)
	}

	return properties, nil
}

// SearchProperties performs basic full-text search
func (r *PostgreSQLPropertyRepository) SearchProperties(query string, limit int) ([]domain.Property, error) {
	if limit <= 0 {
		limit = 50
	}
	
	sqlQuery := `
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		WHERE search_vector @@ plainto_tsquery('spanish', $1)
		ORDER BY 
			ts_rank_cd(search_vector, plainto_tsquery('spanish', $1)) DESC,
			featured DESC,
			created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(sqlQuery, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error performing full-text search: %w", err)
	}
	defer rows.Close()

	var properties []domain.Property
	for rows.Next() {
		var property domain.Property
		var imagesJSON, tagsJSON string

		err := rows.Scan(
			&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
			&property.Province, &property.City, &property.Sector, &property.Address,
			&property.Latitude, &property.Longitude, &property.LocationPrecision,
			&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
			&property.ParkingSpaces, &property.MainImage, &imagesJSON, &property.VideoTour, &property.Tour360,
			&property.RentPrice, &property.CommonExpenses, &property.PricePerM2,
			&property.YearBuilt, &property.Floors, &property.PropertyStatus, &property.Furnished,
			&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
			&property.Security, &property.Elevator, &property.AirConditioning,
			&tagsJSON, &property.Featured, &property.ViewCount, &property.RealEstateCompanyID,
			&property.OwnerID, &property.AgentID, &property.AgencyID, &property.CreatedBy, &property.UpdatedBy,
			&property.CreatedAt, &property.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning search result: %w", err)
		}

		// Convert JSON back to slices
		if imagesJSON != "" {
			json.Unmarshal([]byte(imagesJSON), &property.Images)
		}
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &property.Tags)
		}

		properties = append(properties, property)
	}

	return properties, nil
}

// SearchPropertiesRanked performs full-text search with ranking scores
func (r *PostgreSQLPropertyRepository) SearchPropertiesRanked(query string, limit int) ([]PropertySearchResult, error) {
	if limit <= 0 {
		limit = 50
	}

	sqlQuery := `
		SELECT id, slug, title, description, price, province, city, type,
			   ts_rank_cd(search_vector, plainto_tsquery('spanish', $1)) as rank
		FROM properties 
		WHERE search_vector @@ plainto_tsquery('spanish', $1)
		ORDER BY 
			ts_rank_cd(search_vector, plainto_tsquery('spanish', $1)) DESC,
			featured DESC,
			created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(sqlQuery, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error performing ranked search: %w", err)
	}
	defer rows.Close()

	var results []PropertySearchResult
	for rows.Next() {
		var result PropertySearchResult
		var rank sql.NullFloat64

		err := rows.Scan(
			&result.Property.ID, &result.Property.Slug, &result.Property.Title, 
			&result.Property.Description, &result.Property.Price,
			&result.Property.Province, &result.Property.City, &result.Property.Type,
			&rank,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning ranked search result: %w", err)
		}

		if rank.Valid {
			result.Rank = rank.Float64
		}

		results = append(results, result)
	}

	return results, nil
}

// GetSearchSuggestions returns search suggestions based on existing data
func (r *PostgreSQLPropertyRepository) GetSearchSuggestions(query string, limit int) ([]SearchSuggestion, error) {
	if limit <= 0 {
		limit = 10
	}

	sqlQuery := `
		SELECT * FROM get_search_suggestions($1, $2)
	`

	rows, err := r.db.Query(sqlQuery, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting search suggestions: %w", err)
	}
	defer rows.Close()

	var suggestions []SearchSuggestion
	for rows.Next() {
		var suggestion SearchSuggestion
		err := rows.Scan(&suggestion.Text, &suggestion.Category, &suggestion.Frequency)
		if err != nil {
			return nil, fmt.Errorf("error scanning search suggestion: %w", err)
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// AdvancedSearch performs advanced search with multiple filters
func (r *PostgreSQLPropertyRepository) AdvancedSearch(params AdvancedSearchParams) ([]PropertySearchResult, error) {
	if params.Limit <= 0 {
		params.Limit = 50
	}
	if params.MaxPrice == 0 {
		params.MaxPrice = 999999999
	}
	if params.MaxBedrooms == 0 {
		params.MaxBedrooms = 100
	}
	if params.MaxBathrooms == 0 {
		params.MaxBathrooms = 100
	}
	if params.MaxArea == 0 {
		params.MaxArea = 999999
	}

	sqlQuery := `
		SELECT * FROM advanced_search_properties($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	rows, err := r.db.Query(
		sqlQuery,
		params.Query, params.Province, params.City, params.Type,
		params.MinPrice, params.MaxPrice,
		params.MinBedrooms, params.MaxBedrooms,
		params.MinBathrooms, params.MaxBathrooms,
		params.MinArea, params.MaxArea,
		params.FeaturedOnly, params.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("error performing advanced search: %w", err)
	}
	defer rows.Close()

	var results []PropertySearchResult
	for rows.Next() {
		var result PropertySearchResult
		var rank sql.NullFloat64

		err := rows.Scan(
			&result.Property.ID, &result.Property.Slug, &result.Property.Title,
			&result.Property.Description, &result.Property.Price,
			&result.Property.Province, &result.Property.City, &result.Property.Type,
			&result.Property.Bedrooms, &result.Property.Bathrooms, &result.Property.AreaM2,
			&result.Property.Featured, &rank,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning advanced search result: %w", err)
		}

		if rank.Valid {
			result.Rank = rank.Float64
		}

		results = append(results, result)
	}

	return results, nil
}

// Pagination methods

// GetAllPaginated returns paginated properties with total count
func (r *PostgreSQLPropertyRepository) GetAllPaginated(pagination *domain.PaginationParams) ([]domain.Property, int, error) {
	// Get total count
	countQuery := "SELECT COUNT(*) FROM properties"
	var totalCount int
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting properties: %w", err)
	}

	// Get paginated data
	query := fmt.Sprintf(`
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, pagination.GetOrderBy())

	rows, err := r.db.Query(query, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, 0, fmt.Errorf("error querying paginated properties: %w", err)
	}
	defer rows.Close()

	properties, err := r.scanProperties(rows)
	if err != nil {
		return nil, 0, err
	}

	return properties, totalCount, nil
}

// GetByProvincePaginated returns paginated properties filtered by province
func (r *PostgreSQLPropertyRepository) GetByProvincePaginated(province string, pagination *domain.PaginationParams) ([]domain.Property, int, error) {
	// Get total count
	countQuery := "SELECT COUNT(*) FROM properties WHERE province = $1"
	var totalCount int
	err := r.db.QueryRow(countQuery, province).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting properties by province: %w", err)
	}

	// Get paginated data
	query := fmt.Sprintf(`
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		WHERE province = $1
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, pagination.GetOrderBy())

	rows, err := r.db.Query(query, province, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, 0, fmt.Errorf("error querying paginated properties by province: %w", err)
	}
	defer rows.Close()

	properties, err := r.scanProperties(rows)
	if err != nil {
		return nil, 0, err
	}

	return properties, totalCount, nil
}

// GetByPriceRangePaginated returns paginated properties filtered by price range
func (r *PostgreSQLPropertyRepository) GetByPriceRangePaginated(minPrice, maxPrice float64, pagination *domain.PaginationParams) ([]domain.Property, int, error) {
	// Get total count
	countQuery := "SELECT COUNT(*) FROM properties WHERE price >= $1 AND price <= $2"
	var totalCount int
	err := r.db.QueryRow(countQuery, minPrice, maxPrice).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting properties by price range: %w", err)
	}

	// Get paginated data
	query := fmt.Sprintf(`
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		WHERE price >= $1 AND price <= $2
		ORDER BY %s
		LIMIT $3 OFFSET $4
	`, pagination.GetOrderBy())

	rows, err := r.db.Query(query, minPrice, maxPrice, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, 0, fmt.Errorf("error querying paginated properties by price range: %w", err)
	}
	defer rows.Close()

	properties, err := r.scanProperties(rows)
	if err != nil {
		return nil, 0, err
	}

	return properties, totalCount, nil
}

// SearchPropertiesPaginated performs paginated full-text search
func (r *PostgreSQLPropertyRepository) SearchPropertiesPaginated(query string, pagination *domain.PaginationParams) ([]domain.Property, int, error) {
	// Get total count
	countQuery := "SELECT COUNT(*) FROM properties WHERE search_vector @@ plainto_tsquery('spanish', $1)"
	var totalCount int
	err := r.db.QueryRow(countQuery, query).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting search results: %w", err)
	}

	// Get paginated data with FTS ranking
	sqlQuery := fmt.Sprintf(`
		SELECT id, slug, title, description, price, province, city, sector, address,
			   latitude, longitude, location_precision, type, status, bedrooms, bathrooms, area_m2,
			   parking_spaces, main_image, images, video_tour, tour_360,
			   rent_price, common_expenses, price_per_m2,
			   year_built, floors, property_status, furnished,
			   garage, pool, garden, terrace, balcony, security, elevator, air_conditioning,
			   tags, featured, view_count, real_estate_company_id,
			   owner_id, agent_id, agency_id, created_by, updated_by,
			   created_at, updated_at
		FROM properties 
		WHERE search_vector @@ plainto_tsquery('spanish', $1)
		ORDER BY 
			ts_rank_cd(search_vector, plainto_tsquery('spanish', $1)) DESC,
			%s
		LIMIT $2 OFFSET $3
	`, pagination.GetOrderBy())

	rows, err := r.db.Query(sqlQuery, query, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, 0, fmt.Errorf("error performing paginated search: %w", err)
	}
	defer rows.Close()

	properties, err := r.scanProperties(rows)
	if err != nil {
		return nil, 0, err
	}

	return properties, totalCount, nil
}

// SearchPropertiesRankedPaginated performs paginated full-text search with ranking
func (r *PostgreSQLPropertyRepository) SearchPropertiesRankedPaginated(query string, pagination *domain.PaginationParams) ([]PropertySearchResult, int, error) {
	// Get total count
	countQuery := "SELECT COUNT(*) FROM properties WHERE search_vector @@ plainto_tsquery('spanish', $1)"
	var totalCount int
	err := r.db.QueryRow(countQuery, query).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting ranked search results: %w", err)
	}

	// Get paginated data with ranking
	sqlQuery := `
		SELECT id, slug, title, description, price, province, city, type,
			   ts_rank_cd(search_vector, plainto_tsquery('spanish', $1)) as rank
		FROM properties 
		WHERE search_vector @@ plainto_tsquery('spanish', $1)
		ORDER BY 
			ts_rank_cd(search_vector, plainto_tsquery('spanish', $1)) DESC,
			featured DESC,
			created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(sqlQuery, query, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, 0, fmt.Errorf("error performing paginated ranked search: %w", err)
	}
	defer rows.Close()

	var results []PropertySearchResult
	for rows.Next() {
		var result PropertySearchResult
		var rank sql.NullFloat64

		err := rows.Scan(
			&result.Property.ID, &result.Property.Slug, &result.Property.Title,
			&result.Property.Description, &result.Property.Price,
			&result.Property.Province, &result.Property.City, &result.Property.Type,
			&rank,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning ranked search result: %w", err)
		}

		if rank.Valid {
			result.Rank = rank.Float64
		}

		results = append(results, result)
	}

	return results, totalCount, nil
}

// AdvancedSearchPaginated performs paginated advanced search
func (r *PostgreSQLPropertyRepository) AdvancedSearchPaginated(params AdvancedSearchParams, pagination *domain.PaginationParams) ([]PropertySearchResult, int, error) {
	// For now, use existing advanced search function and apply pagination manually
	// In a real implementation, you'd modify the stored procedure to support pagination
	if params.Limit <= 0 {
		params.Limit = pagination.GetLimit()
	}
	
	// Get total count using a simpler query
	countQuery := `
		SELECT COUNT(*) FROM properties 
		WHERE ($1 = '' OR search_vector @@ plainto_tsquery('spanish', $1))
		AND ($2 = '' OR province = $2)
		AND ($3 = '' OR city = $3)
		AND ($4 = '' OR type = $4)
		AND price >= $5 AND price <= $6
		AND bedrooms >= $7 AND bedrooms <= $8
		AND bathrooms >= $9 AND bathrooms <= $10
		AND area_m2 >= $11 AND area_m2 <= $12
		AND ($13 = false OR featured = true)
	`
	
	maxPrice := params.MaxPrice
	if maxPrice == 0 {
		maxPrice = 999999999
	}
	maxBedrooms := params.MaxBedrooms
	if maxBedrooms == 0 {
		maxBedrooms = 100
	}
	maxBathrooms := params.MaxBathrooms
	if maxBathrooms == 0 {
		maxBathrooms = 100
	}
	maxArea := params.MaxArea
	if maxArea == 0 {
		maxArea = 999999
	}
	
	var totalCount int
	err := r.db.QueryRow(countQuery, 
		params.Query, params.Province, params.City, params.Type,
		params.MinPrice, maxPrice,
		params.MinBedrooms, maxBedrooms,
		params.MinBathrooms, maxBathrooms,
		params.MinArea, maxArea,
		params.FeaturedOnly).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting advanced search results: %w", err)
	}

	// Use existing advanced search with offset simulation
	offset := pagination.GetOffset()
	adjustedParams := params
	adjustedParams.Limit = pagination.GetLimit() + offset
	
	allResults, err := r.AdvancedSearch(adjustedParams)
	if err != nil {
		return nil, 0, err
	}
	
	// Apply offset manually
	if offset >= len(allResults) {
		return []PropertySearchResult{}, totalCount, nil
	}
	
	end := offset + pagination.GetLimit()
	if end > len(allResults) {
		end = len(allResults)
	}
	
	results := allResults[offset:end]
	return results, totalCount, nil
}

// scanProperties is a helper function to scan properties from rows
func (r *PostgreSQLPropertyRepository) scanProperties(rows *sql.Rows) ([]domain.Property, error) {
	var properties []domain.Property

	for rows.Next() {
		var property domain.Property
		var imagesJSON, tagsJSON string

		err := rows.Scan(
			&property.ID, &property.Slug, &property.Title, &property.Description, &property.Price,
			&property.Province, &property.City, &property.Sector, &property.Address,
			&property.Latitude, &property.Longitude, &property.LocationPrecision,
			&property.Type, &property.Status, &property.Bedrooms, &property.Bathrooms, &property.AreaM2,
			&property.ParkingSpaces, &property.MainImage, &imagesJSON, &property.VideoTour, &property.Tour360,
			&property.RentPrice, &property.CommonExpenses, &property.PricePerM2,
			&property.YearBuilt, &property.Floors, &property.PropertyStatus, &property.Furnished,
			&property.Garage, &property.Pool, &property.Garden, &property.Terrace, &property.Balcony,
			&property.Security, &property.Elevator, &property.AirConditioning,
			&tagsJSON, &property.Featured, &property.ViewCount, &property.RealEstateCompanyID,
			&property.OwnerID, &property.AgentID, &property.AgencyID, &property.CreatedBy, &property.UpdatedBy,
			&property.CreatedAt, &property.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning property: %w", err)
		}

		// Convert JSON back to slices
		if imagesJSON != "" {
			err = json.Unmarshal([]byte(imagesJSON), &property.Images)
			if err != nil {
				property.Images = []string{} // Continue with empty slice if JSON is invalid
			}
		}

		if tagsJSON != "" {
			err = json.Unmarshal([]byte(tagsJSON), &property.Tags)
			if err != nil {
				property.Tags = []string{} // Continue with empty slice if JSON is invalid
			}
		}

		properties = append(properties, property)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating properties: %w", err)
	}

	return properties, nil
}

// ConnectDatabase establishes connection to PostgreSQL
func ConnectDatabase(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("PostgreSQL connection established successfully")
	return db, nil
}