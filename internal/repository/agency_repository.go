package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"realty-core/internal/domain"
)

// AgencyRepository handles database operations for agencies
type AgencyRepository struct {
	db *sql.DB
}

// NewAgencyRepository creates a new agency repository
func NewAgencyRepository(db *sql.DB) *AgencyRepository {
	return &AgencyRepository{db: db}
}

// Create creates a new agency in the database
func (r *AgencyRepository) Create(agency *domain.Agency) error {
	query := `
		INSERT INTO agencies (
			id, name, ruc, address, phone, email, website, description, 
			logo_url, active, license_number, license_expiry, commission, 
			business_hours, social_media, specialties, service_areas, 
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)`

	// Convert maps and slices to JSON
	socialMediaJSON, err := json.Marshal(agency.SocialMedia)
	if err != nil {
		return fmt.Errorf("failed to marshal social media: %w", err)
	}

	specialtiesJSON, err := json.Marshal(agency.Specialties)
	if err != nil {
		return fmt.Errorf("failed to marshal specialties: %w", err)
	}

	serviceAreasJSON, err := json.Marshal(agency.ServiceAreas)
	if err != nil {
		return fmt.Errorf("failed to marshal service areas: %w", err)
	}

	_, err = r.db.Exec(query,
		agency.ID, agency.Name, agency.RUC, agency.Address, agency.Phone,
		agency.Email, agency.Website, agency.Description, agency.LogoURL,
		agency.Active, agency.LicenseNumber, agency.LicenseExpiry,
		agency.Commission, agency.BusinessHours, socialMediaJSON,
		specialtiesJSON, serviceAreasJSON, agency.CreatedAt, agency.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create agency: %w", err)
	}

	return nil
}

// GetByID retrieves an agency by ID
func (r *AgencyRepository) GetByID(id string) (*domain.Agency, error) {
	query := `
		SELECT id, name, ruc, address, phone, email, website, description, 
			   logo_url, active, license_number, license_expiry, commission, 
			   business_hours, social_media, specialties, service_areas, 
			   created_at, updated_at
		FROM agencies 
		WHERE id = $1`

	agency := &domain.Agency{}
	var socialMediaJSON, specialtiesJSON, serviceAreasJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&agency.ID, &agency.Name, &agency.RUC, &agency.Address, &agency.Phone,
		&agency.Email, &agency.Website, &agency.Description, &agency.LogoURL,
		&agency.Active, &agency.LicenseNumber, &agency.LicenseExpiry,
		&agency.Commission, &agency.BusinessHours, &socialMediaJSON,
		&specialtiesJSON, &serviceAreasJSON, &agency.CreatedAt, &agency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agency not found with id: %s", id)
		}
		return nil, fmt.Errorf("failed to get agency by id: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(socialMediaJSON, &agency.SocialMedia); err != nil {
		return nil, fmt.Errorf("failed to unmarshal social media: %w", err)
	}

	if err := json.Unmarshal(specialtiesJSON, &agency.Specialties); err != nil {
		return nil, fmt.Errorf("failed to unmarshal specialties: %w", err)
	}

	if err := json.Unmarshal(serviceAreasJSON, &agency.ServiceAreas); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service areas: %w", err)
	}

	return agency, nil
}

// GetByRUC retrieves an agency by RUC
func (r *AgencyRepository) GetByRUC(ruc string) (*domain.Agency, error) {
	query := `
		SELECT id, name, ruc, address, phone, email, website, description, 
			   logo_url, active, license_number, license_expiry, commission, 
			   business_hours, social_media, specialties, service_areas, 
			   created_at, updated_at
		FROM agencies 
		WHERE ruc = $1`

	agency := &domain.Agency{}
	var socialMediaJSON, specialtiesJSON, serviceAreasJSON []byte

	err := r.db.QueryRow(query, ruc).Scan(
		&agency.ID, &agency.Name, &agency.RUC, &agency.Address, &agency.Phone,
		&agency.Email, &agency.Website, &agency.Description, &agency.LogoURL,
		&agency.Active, &agency.LicenseNumber, &agency.LicenseExpiry,
		&agency.Commission, &agency.BusinessHours, &socialMediaJSON,
		&specialtiesJSON, &serviceAreasJSON, &agency.CreatedAt, &agency.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agency not found with ruc: %s", ruc)
		}
		return nil, fmt.Errorf("failed to get agency by ruc: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(socialMediaJSON, &agency.SocialMedia); err != nil {
		return nil, fmt.Errorf("failed to unmarshal social media: %w", err)
	}

	if err := json.Unmarshal(specialtiesJSON, &agency.Specialties); err != nil {
		return nil, fmt.Errorf("failed to unmarshal specialties: %w", err)
	}

	if err := json.Unmarshal(serviceAreasJSON, &agency.ServiceAreas); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service areas: %w", err)
	}

	return agency, nil
}

// Update updates an agency in the database
func (r *AgencyRepository) Update(agency *domain.Agency) error {
	query := `
		UPDATE agencies SET 
			name = $2, ruc = $3, address = $4, phone = $5, email = $6, 
			website = $7, description = $8, logo_url = $9, active = $10, 
			license_number = $11, license_expiry = $12, commission = $13, 
			business_hours = $14, social_media = $15, specialties = $16, 
			service_areas = $17, updated_at = $18
		WHERE id = $1`

	// Convert maps and slices to JSON
	socialMediaJSON, err := json.Marshal(agency.SocialMedia)
	if err != nil {
		return fmt.Errorf("failed to marshal social media: %w", err)
	}

	specialtiesJSON, err := json.Marshal(agency.Specialties)
	if err != nil {
		return fmt.Errorf("failed to marshal specialties: %w", err)
	}

	serviceAreasJSON, err := json.Marshal(agency.ServiceAreas)
	if err != nil {
		return fmt.Errorf("failed to marshal service areas: %w", err)
	}

	_, err = r.db.Exec(query,
		agency.ID, agency.Name, agency.RUC, agency.Address, agency.Phone,
		agency.Email, agency.Website, agency.Description, agency.LogoURL,
		agency.Active, agency.LicenseNumber, agency.LicenseExpiry,
		agency.Commission, agency.BusinessHours, socialMediaJSON,
		specialtiesJSON, serviceAreasJSON, agency.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	return nil
}

// Delete deletes an agency from the database
func (r *AgencyRepository) Delete(id string) error {
	query := `DELETE FROM agencies WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete agency: %w", err)
	}
	return nil
}

// GetActive retrieves all active agencies
func (r *AgencyRepository) GetActive() ([]*domain.Agency, error) {
	query := `
		SELECT id, name, ruc, address, phone, email, website, description, 
			   logo_url, active, license_number, license_expiry, commission, 
			   business_hours, social_media, specialties, service_areas, 
			   created_at, updated_at
		FROM agencies 
		WHERE active = TRUE 
		  AND (license_expiry IS NULL OR license_expiry > CURRENT_TIMESTAMP)
		ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active agencies: %w", err)
	}
	defer rows.Close()

	var agencies []*domain.Agency
	for rows.Next() {
		agency := &domain.Agency{}
		var socialMediaJSON, specialtiesJSON, serviceAreasJSON []byte

		err := rows.Scan(
			&agency.ID, &agency.Name, &agency.RUC, &agency.Address, &agency.Phone,
			&agency.Email, &agency.Website, &agency.Description, &agency.LogoURL,
			&agency.Active, &agency.LicenseNumber, &agency.LicenseExpiry,
			&agency.Commission, &agency.BusinessHours, &socialMediaJSON,
			&specialtiesJSON, &serviceAreasJSON, &agency.CreatedAt, &agency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agency: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(socialMediaJSON, &agency.SocialMedia); err != nil {
			return nil, fmt.Errorf("failed to unmarshal social media: %w", err)
		}

		if err := json.Unmarshal(specialtiesJSON, &agency.Specialties); err != nil {
			return nil, fmt.Errorf("failed to unmarshal specialties: %w", err)
		}

		if err := json.Unmarshal(serviceAreasJSON, &agency.ServiceAreas); err != nil {
			return nil, fmt.Errorf("failed to unmarshal service areas: %w", err)
		}

		agencies = append(agencies, agency)
	}

	return agencies, nil
}

// Search searches agencies with filters
func (r *AgencyRepository) Search(params *domain.AgencySearchParams) ([]*domain.Agency, int, error) {
	// Build base query
	baseQuery := `
		SELECT id, name, ruc, address, phone, email, website, description, 
			   logo_url, active, license_number, license_expiry, commission, 
			   business_hours, social_media, specialties, service_areas, 
			   created_at, updated_at
		FROM agencies WHERE 1=1`

	countQuery := `SELECT COUNT(*) FROM agencies WHERE 1=1`

	var args []interface{}
	var conditions []string
	argIndex := 1

	// Add search conditions
	if params.Query != "" {
		conditions = append(conditions, fmt.Sprintf(`to_tsvector('spanish', name || ' ' || description) @@ plainto_tsquery('spanish', $%d)`, argIndex))
		args = append(args, params.Query)
		argIndex++
	}

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf(`active = $%d`, argIndex))
		args = append(args, *params.Active)
		argIndex++
	}

	if len(params.ServiceAreas) > 0 {
		conditions = append(conditions, fmt.Sprintf(`service_areas @> $%d`, argIndex))
		serviceAreasJSON, _ := json.Marshal(params.ServiceAreas)
		args = append(args, string(serviceAreasJSON))
		argIndex++
	}

	if len(params.Specialties) > 0 {
		conditions = append(conditions, fmt.Sprintf(`specialties @> $%d`, argIndex))
		specialtiesJSON, _ := json.Marshal(params.Specialties)
		args = append(args, string(specialtiesJSON))
		argIndex++
	}

	if params.MinCommission != nil {
		conditions = append(conditions, fmt.Sprintf(`commission >= $%d`, argIndex))
		args = append(args, *params.MinCommission)
		argIndex++
	}

	if params.MaxCommission != nil {
		conditions = append(conditions, fmt.Sprintf(`commission <= $%d`, argIndex))
		args = append(args, *params.MaxCommission)
		argIndex++
	}

	if params.LicenseValid != nil && *params.LicenseValid {
		conditions = append(conditions, `(license_expiry IS NULL OR license_expiry > CURRENT_TIMESTAMP)`)
	}

	// Add conditions to queries
	if len(conditions) > 0 {
		conditionStr := " AND " + strings.Join(conditions, " AND ")
		baseQuery += conditionStr
		countQuery += conditionStr
	}

	// Get total count
	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count agencies: %w", err)
	}

	// Add pagination and sorting
	if params.Pagination != nil {
		orderBy := params.Pagination.GetOrderBy()
		baseQuery += fmt.Sprintf(` ORDER BY %s LIMIT %d OFFSET %d`,
			orderBy, params.Pagination.GetLimit(), params.Pagination.GetOffset())
	}

	// Execute main query
	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search agencies: %w", err)
	}
	defer rows.Close()

	var agencies []*domain.Agency
	for rows.Next() {
		agency := &domain.Agency{}
		var socialMediaJSON, specialtiesJSON, serviceAreasJSON []byte

		err := rows.Scan(
			&agency.ID, &agency.Name, &agency.RUC, &agency.Address, &agency.Phone,
			&agency.Email, &agency.Website, &agency.Description, &agency.LogoURL,
			&agency.Active, &agency.LicenseNumber, &agency.LicenseExpiry,
			&agency.Commission, &agency.BusinessHours, &socialMediaJSON,
			&specialtiesJSON, &serviceAreasJSON, &agency.CreatedAt, &agency.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan agency: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(socialMediaJSON, &agency.SocialMedia); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal social media: %w", err)
		}

		if err := json.Unmarshal(specialtiesJSON, &agency.Specialties); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal specialties: %w", err)
		}

		if err := json.Unmarshal(serviceAreasJSON, &agency.ServiceAreas); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal service areas: %w", err)
		}

		agencies = append(agencies, agency)
	}

	return agencies, totalCount, nil
}

// GetByServiceArea retrieves agencies that serve a specific area
func (r *AgencyRepository) GetByServiceArea(province string) ([]*domain.Agency, error) {
	query := `
		SELECT id, name, ruc, address, phone, email, website, description, 
			   logo_url, active, license_number, license_expiry, commission, 
			   business_hours, social_media, specialties, service_areas, 
			   created_at, updated_at
		FROM agencies 
		WHERE active = TRUE 
		  AND (license_expiry IS NULL OR license_expiry > CURRENT_TIMESTAMP)
		  AND service_areas @> $1
		ORDER BY name`

	provinceJSON, _ := json.Marshal([]string{province})
	rows, err := r.db.Query(query, string(provinceJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to get agencies by service area: %w", err)
	}
	defer rows.Close()

	var agencies []*domain.Agency
	for rows.Next() {
		agency := &domain.Agency{}
		var socialMediaJSON, specialtiesJSON, serviceAreasJSON []byte

		err := rows.Scan(
			&agency.ID, &agency.Name, &agency.RUC, &agency.Address, &agency.Phone,
			&agency.Email, &agency.Website, &agency.Description, &agency.LogoURL,
			&agency.Active, &agency.LicenseNumber, &agency.LicenseExpiry,
			&agency.Commission, &agency.BusinessHours, &socialMediaJSON,
			&specialtiesJSON, &serviceAreasJSON, &agency.CreatedAt, &agency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agency: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(socialMediaJSON, &agency.SocialMedia); err != nil {
			return nil, fmt.Errorf("failed to unmarshal social media: %w", err)
		}

		if err := json.Unmarshal(specialtiesJSON, &agency.Specialties); err != nil {
			return nil, fmt.Errorf("failed to unmarshal specialties: %w", err)
		}

		if err := json.Unmarshal(serviceAreasJSON, &agency.ServiceAreas); err != nil {
			return nil, fmt.Errorf("failed to unmarshal service areas: %w", err)
		}

		agencies = append(agencies, agency)
	}

	return agencies, nil
}

// GetBySpecialty retrieves agencies with a specific specialty
func (r *AgencyRepository) GetBySpecialty(specialty string) ([]*domain.Agency, error) {
	query := `
		SELECT id, name, ruc, address, phone, email, website, description, 
			   logo_url, active, license_number, license_expiry, commission, 
			   business_hours, social_media, specialties, service_areas, 
			   created_at, updated_at
		FROM agencies 
		WHERE active = TRUE 
		  AND (license_expiry IS NULL OR license_expiry > CURRENT_TIMESTAMP)
		  AND specialties @> $1
		ORDER BY name`

	specialtyJSON, _ := json.Marshal([]string{specialty})
	rows, err := r.db.Query(query, string(specialtyJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to get agencies by specialty: %w", err)
	}
	defer rows.Close()

	var agencies []*domain.Agency
	for rows.Next() {
		agency := &domain.Agency{}
		var socialMediaJSON, specialtiesJSON, serviceAreasJSON []byte

		err := rows.Scan(
			&agency.ID, &agency.Name, &agency.RUC, &agency.Address, &agency.Phone,
			&agency.Email, &agency.Website, &agency.Description, &agency.LogoURL,
			&agency.Active, &agency.LicenseNumber, &agency.LicenseExpiry,
			&agency.Commission, &agency.BusinessHours, &socialMediaJSON,
			&specialtiesJSON, &serviceAreasJSON, &agency.CreatedAt, &agency.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agency: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(socialMediaJSON, &agency.SocialMedia); err != nil {
			return nil, fmt.Errorf("failed to unmarshal social media: %w", err)
		}

		if err := json.Unmarshal(specialtiesJSON, &agency.Specialties); err != nil {
			return nil, fmt.Errorf("failed to unmarshal specialties: %w", err)
		}

		if err := json.Unmarshal(serviceAreasJSON, &agency.ServiceAreas); err != nil {
			return nil, fmt.Errorf("failed to unmarshal service areas: %w", err)
		}

		agencies = append(agencies, agency)
	}

	return agencies, nil
}

// GetStatistics returns agency statistics
func (r *AgencyRepository) GetStatistics() (*domain.AgencyStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_agencies,
			COUNT(*) FILTER (WHERE active = TRUE) as active_agencies,
			COUNT(*) FILTER (WHERE license_number IS NOT NULL AND license_number != '') as licensed_agencies,
			COUNT(*) FILTER (WHERE license_expiry IS NOT NULL AND license_expiry <= CURRENT_TIMESTAMP) as expired_licenses,
			COALESCE(AVG(commission), 0) as average_commission,
			(SELECT COUNT(*) FROM users WHERE role = 'agent' AND agency_id IS NOT NULL) as total_agents,
			(SELECT COUNT(*) FROM properties WHERE agency_id IS NOT NULL) as total_properties
		FROM agencies`

	stats := &domain.AgencyStats{}
	err := r.db.QueryRow(query).Scan(
		&stats.TotalAgencies, &stats.ActiveAgencies, &stats.LicensedAgencies,
		&stats.ExpiredLicenses, &stats.AverageCommission, &stats.TotalAgents,
		&stats.TotalProperties,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get agency statistics: %w", err)
	}

	return stats, nil
}

// GetPerformance returns performance metrics for an agency
func (r *AgencyRepository) GetPerformance(agencyID string) (*domain.AgencyPerformance, error) {
	query := `
		SELECT 
			a.id as agency_id,
			a.name as agency_name,
			COUNT(p.id) as total_properties,
			COUNT(p.id) FILTER (WHERE p.status = 'sold') as sold_properties,
			COUNT(p.id) FILTER (WHERE p.status = 'rented') as rented_properties,
			COALESCE(SUM(CASE WHEN p.status = 'sold' THEN p.price END), 0) as total_sales_value,
			COALESCE(SUM(CASE WHEN p.status = 'rented' THEN p.rent_price END), 0) as total_rent_value,
			COALESCE(AVG(p.price), 0) as average_property_value,
			(SELECT COUNT(*) FROM users WHERE agency_id = a.id AND role = 'agent') as total_agents,
			(SELECT COUNT(*) FROM users WHERE agency_id = a.id AND role = 'agent' AND active = TRUE) as active_agents,
			CASE 
				WHEN COUNT(p.id) > 0 THEN 
					ROUND((COUNT(p.id) FILTER (WHERE p.status IN ('sold', 'rented')) * 100.0 / COUNT(p.id)), 2)
				ELSE 0
			END as conversion_rate
		FROM agencies a
		LEFT JOIN properties p ON a.id = p.agency_id
		WHERE a.id = $1
		GROUP BY a.id, a.name`

	performance := &domain.AgencyPerformance{}
	err := r.db.QueryRow(query, agencyID).Scan(
		&performance.AgencyID, &performance.AgencyName,
		&performance.TotalProperties, &performance.SoldProperties,
		&performance.RentedProperties, &performance.TotalSalesValue,
		&performance.TotalRentValue, &performance.AveragePropertyValue,
		&performance.TotalAgents, &performance.ActiveAgents,
		&performance.ConversionRate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agency not found with id: %s", agencyID)
		}
		return nil, fmt.Errorf("failed to get agency performance: %w", err)
	}

	return performance, nil
}

// GetWithAgents returns agency with its associated agents
func (r *AgencyRepository) GetWithAgents(agencyID string) (*domain.AgencyWithAgents, error) {
	// Get the agency first
	agency, err := r.GetByID(agencyID)
	if err != nil {
		return nil, err
	}

	// Get associated agents
	userRepo := NewUserRepository(r.db)
	agents, err := userRepo.GetByAgency(agencyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents for agency: %w", err)
	}

	return &domain.AgencyWithAgents{
		Agency: agency,
		Agents: agents,
	}, nil
}

// Activate activates an agency
func (r *AgencyRepository) Activate(id string) error {
	query := `
		UPDATE agencies 
		SET active = TRUE, updated_at = $2
		WHERE id = $1`

	_, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to activate agency: %w", err)
	}

	return nil
}

// Deactivate deactivates an agency
func (r *AgencyRepository) Deactivate(id string) error {
	query := `
		UPDATE agencies 
		SET active = FALSE, updated_at = $2
		WHERE id = $1`

	_, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to deactivate agency: %w", err)
	}

	return nil
}