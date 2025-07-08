package repositorio

import (
	"database/sql"
	"fmt"
	"realty-core/internal/domain"
)

// RealEstateCompanyRepository handles real estate company data access operations
type RealEstateCompanyRepository struct {
	db *sql.DB
}

// NewRealEstateCompanyRepository creates a new real estate company repository instance
func NewRealEstateCompanyRepository(db *sql.DB) *RealEstateCompanyRepository {
	return &RealEstateCompanyRepository{db: db}
}

// Create creates a new real estate company in the database
func (r *RealEstateCompanyRepository) Create(company *domain.RealEstateCompany) error {
	query := `
		INSERT INTO real_estate_companies (
			id, name, ruc, address, description, phone, email, website, logo_url, active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Exec(query,
		company.ID,
		company.Name,
		company.RUC,
		company.Address,
		company.Description,
		company.Phone,
		company.Email,
		company.Website,
		company.LogoURL,
		company.Active,
		company.CreatedAt,
		company.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating real estate company: %w", err)
	}

	return nil
}

// GetByID retrieves a real estate company by its ID
func (r *RealEstateCompanyRepository) GetByID(id string) (*domain.RealEstateCompany, error) {
	query := `
		SELECT id, name, ruc, address, description, phone, email, website, logo_url, active, created_at, updated_at
		FROM real_estate_companies 
		WHERE id = $1`

	company := &domain.RealEstateCompany{}

	err := r.db.QueryRow(query, id).Scan(
		&company.ID,
		&company.Name,
		&company.RUC,
		&company.Address,
		&company.Description,
		&company.Phone,
		&company.Email,
		&company.Website,
		&company.LogoURL,
		&company.Active,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("real estate company not found with ID: %s", id)
		}
		return nil, fmt.Errorf("error retrieving real estate company: %w", err)
	}

	return company, nil
}

// GetByRUC retrieves a real estate company by its RUC
func (r *RealEstateCompanyRepository) GetByRUC(ruc string) (*domain.RealEstateCompany, error) {
	query := `
		SELECT id, name, ruc, address, description, phone, email, website, logo_url, active, created_at, updated_at
		FROM real_estate_companies 
		WHERE ruc = $1`

	company := &domain.RealEstateCompany{}

	err := r.db.QueryRow(query, ruc).Scan(
		&company.ID,
		&company.Name,
		&company.RUC,
		&company.Address,
		&company.Description,
		&company.Phone,
		&company.Email,
		&company.Website,
		&company.LogoURL,
		&company.Active,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("real estate company not found with RUC: %s", ruc)
		}
		return nil, fmt.Errorf("error retrieving real estate company: %w", err)
	}

	return company, nil
}

// GetAll retrieves all real estate companies
func (r *RealEstateCompanyRepository) GetAll() ([]*domain.RealEstateCompany, error) {
	query := `
		SELECT id, name, ruc, address, description, phone, email, website, logo_url, active, created_at, updated_at
		FROM real_estate_companies 
		ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying real estate companies: %w", err)
	}
	defer rows.Close()

	var companies []*domain.RealEstateCompany
	for rows.Next() {
		company := &domain.RealEstateCompany{}

		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.RUC,
			&company.Address,
			&company.Description,
			&company.Phone,
			&company.Email,
			&company.Website,
			&company.LogoURL,
			&company.Active,
			&company.CreatedAt,
			&company.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning real estate company: %w", err)
		}

		companies = append(companies, company)
	}

	return companies, nil
}

// GetActive retrieves all active real estate companies
func (r *RealEstateCompanyRepository) GetActive() ([]*domain.RealEstateCompany, error) {
	query := `
		SELECT id, name, ruc, address, description, phone, email, website, logo_url, active, created_at, updated_at
		FROM real_estate_companies 
		WHERE active = true
		ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying active real estate companies: %w", err)
	}
	defer rows.Close()

	var companies []*domain.RealEstateCompany
	for rows.Next() {
		company := &domain.RealEstateCompany{}

		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.RUC,
			&company.Address,
			&company.Description,
			&company.Phone,
			&company.Email,
			&company.Website,
			&company.LogoURL,
			&company.Active,
			&company.CreatedAt,
			&company.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning real estate company: %w", err)
		}

		companies = append(companies, company)
	}

	return companies, nil
}

// Update updates an existing real estate company
func (r *RealEstateCompanyRepository) Update(company *domain.RealEstateCompany) error {
	// Update timestamp
	company.UpdateTimestamp()

	query := `
		UPDATE real_estate_companies SET
			name = $2, ruc = $3, address = $4, description = $5, phone = $6, 
			email = $7, website = $8, logo_url = $9, active = $10, updated_at = $11
		WHERE id = $1`

	result, err := r.db.Exec(query,
		company.ID,
		company.Name,
		company.RUC,
		company.Address,
		company.Description,
		company.Phone,
		company.Email,
		company.Website,
		company.LogoURL,
		company.Active,
		company.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating real estate company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("real estate company not found with ID: %s", company.ID)
	}

	return nil
}

// Deactivate marks a real estate company as inactive
func (r *RealEstateCompanyRepository) Deactivate(id string) error {
	query := `UPDATE real_estate_companies SET active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deactivating real estate company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("real estate company not found with ID: %s", id)
	}

	return nil
}

// Activate marks a real estate company as active
func (r *RealEstateCompanyRepository) Activate(id string) error {
	query := `UPDATE real_estate_companies SET active = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error activating real estate company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("real estate company not found with ID: %s", id)
	}

	return nil
}

// SearchByName searches real estate companies by name using full-text search
func (r *RealEstateCompanyRepository) SearchByName(searchTerm string) ([]*domain.RealEstateCompany, error) {
	query := `
		SELECT id, name, ruc, address, description, phone, email, website, logo_url, active, created_at, updated_at
		FROM real_estate_companies 
		WHERE active = true
		  AND to_tsvector('english', name || ' ' || description) @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(to_tsvector('english', name || ' ' || description), plainto_tsquery('english', $1)) DESC`

	rows, err := r.db.Query(query, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("error searching real estate companies: %w", err)
	}
	defer rows.Close()

	var companies []*domain.RealEstateCompany
	for rows.Next() {
		company := &domain.RealEstateCompany{}

		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.RUC,
			&company.Address,
			&company.Description,
			&company.Phone,
			&company.Email,
			&company.Website,
			&company.LogoURL,
			&company.Active,
			&company.CreatedAt,
			&company.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning real estate company: %w", err)
		}

		companies = append(companies, company)
	}

	return companies, nil
}

// ExistsByRUC checks if a real estate company with the given RUC already exists
func (r *RealEstateCompanyRepository) ExistsByRUC(ruc string) (bool, error) {
	query := `SELECT COUNT(*) FROM real_estate_companies WHERE ruc = $1`

	var count int
	err := r.db.QueryRow(query, ruc).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking RUC existence: %w", err)
	}

	return count > 0, nil
}

// ExistsByEmail checks if a real estate company with the given email already exists
func (r *RealEstateCompanyRepository) ExistsByEmail(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM real_estate_companies WHERE email = $1`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}

	return count > 0, nil
}

// ValidateRUC validates if a RUC is valid using the database function
func (r *RealEstateCompanyRepository) ValidateRUC(ruc string) (bool, error) {
	query := `SELECT validate_ecuador_ruc($1)`

	var isValid bool
	err := r.db.QueryRow(query, ruc).Scan(&isValid)
	if err != nil {
		return false, fmt.Errorf("error validating RUC: %w", err)
	}

	return isValid, nil
}

// GetStatistics returns statistics about real estate companies
func (r *RealEstateCompanyRepository) GetStatistics() (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_companies,
			COUNT(*) FILTER (WHERE active = true) as active_companies,
			COUNT(*) FILTER (WHERE active = false) as inactive_companies,
			ROUND(
				(COUNT(*) FILTER (WHERE active = true) * 100.0 / NULLIF(COUNT(*), 0)), 
				2
			) as active_percentage
		FROM real_estate_companies`

	var totalCompanies, activeCompanies, inactiveCompanies int64
	var activePercentage float64

	err := r.db.QueryRow(query).Scan(&totalCompanies, &activeCompanies, &inactiveCompanies, &activePercentage)
	if err != nil {
		return nil, fmt.Errorf("error getting real estate company statistics: %w", err)
	}

	stats := map[string]interface{}{
		"total_companies":    totalCompanies,
		"active_companies":   activeCompanies,
		"inactive_companies": inactiveCompanies,
		"active_percentage":  activePercentage,
	}

	return stats, nil
}

// GetCompaniesWithPropertyCount returns companies along with their property counts
func (r *RealEstateCompanyRepository) GetCompaniesWithPropertyCount() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			rec.id,
			rec.name,
			rec.active,
			rec.created_at,
			COUNT(p.id) as total_properties,
			COUNT(p.id) FILTER (WHERE p.status = 'available') as available_properties,
			COUNT(p.id) FILTER (WHERE p.status = 'sold') as sold_properties,
			COUNT(p.id) FILTER (WHERE p.status = 'rented') as rented_properties
		FROM real_estate_companies rec
		LEFT JOIN properties p ON rec.id = p.real_estate_company_id
		GROUP BY rec.id, rec.name, rec.active, rec.created_at
		ORDER BY rec.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying companies with property count: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, name string
		var active bool
		var createdAt string
		var totalProperties, availableProperties, soldProperties, rentedProperties int

		err := rows.Scan(
			&id, &name, &active, &createdAt,
			&totalProperties, &availableProperties, &soldProperties, &rentedProperties,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning company with property count: %w", err)
		}

		result := map[string]interface{}{
			"id":                   id,
			"name":                 name,
			"active":               active,
			"created_at":           createdAt,
			"total_properties":     totalProperties,
			"available_properties": availableProperties,
			"sold_properties":      soldProperties,
			"rented_properties":    rentedProperties,
		}

		results = append(results, result)
	}

	return results, nil
}

// Delete permanently deletes a real estate company (use with caution)
func (r *RealEstateCompanyRepository) Delete(id string) error {
	// First, check if the company has any properties or agents
	var propertyCount, agentCount int

	// Check properties
	err := r.db.QueryRow("SELECT COUNT(*) FROM properties WHERE real_estate_company_id = $1", id).Scan(&propertyCount)
	if err != nil {
		return fmt.Errorf("error checking property count: %w", err)
	}

	// Check agents
	err = r.db.QueryRow("SELECT COUNT(*) FROM users WHERE real_estate_company_id = $1 AND user_type = 'agent'", id).Scan(&agentCount)
	if err != nil {
		return fmt.Errorf("error checking agent count: %w", err)
	}

	if propertyCount > 0 || agentCount > 0 {
		return fmt.Errorf("cannot delete company: has %d properties and %d agents. Consider deactivating instead", propertyCount, agentCount)
	}

	// Delete the company
	query := `DELETE FROM real_estate_companies WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting real estate company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("real estate company not found with ID: %s", id)
	}

	return nil
}
