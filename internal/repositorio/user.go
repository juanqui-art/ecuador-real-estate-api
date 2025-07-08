package repositorio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"realty-core/internal/domain"
	"time"
)

// UserRepository handles user data access operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *domain.User) error {
	// Serialize JSON fields
	provincesJSON, err := json.Marshal(user.PreferredProvinces)
	if err != nil {
		return fmt.Errorf("error marshaling preferred provinces: %w", err)
	}

	propertyTypesJSON, err := json.Marshal(user.PreferredPropertyTypes)
	if err != nil {
		return fmt.Errorf("error marshaling preferred property types: %w", err)
	}

	query := `
		INSERT INTO users (
			id, first_name, last_name, email, phone, national_id, date_of_birth,
			user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)`

	_, err = r.db.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.NationalID,
		nullTimePtr(user.DateOfBirth),
		user.UserType,
		user.Active,
		nullFloat64Ptr(user.MinBudget),
		nullFloat64Ptr(user.MaxBudget),
		string(provincesJSON),
		string(propertyTypesJSON),
		user.AvatarURL,
		user.Bio,
		nullStringPtr(user.RealEstateCompanyID),
		user.ReceiveNotifications,
		user.ReceiveNewsletter,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth,
			   user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			   avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			   created_at, updated_at
		FROM users WHERE id = $1`

	user := &domain.User{}
	var dateOfBirth sql.NullTime
	var minBudget, maxBudget sql.NullFloat64
	var provincesJSON, propertyTypesJSON string
	var realEstateCompanyID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.NationalID,
		&dateOfBirth,
		&user.UserType,
		&user.Active,
		&minBudget,
		&maxBudget,
		&provincesJSON,
		&propertyTypesJSON,
		&user.AvatarURL,
		&user.Bio,
		&realEstateCompanyID,
		&user.ReceiveNotifications,
		&user.ReceiveNewsletter,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with ID: %s", id)
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Handle nullable fields
	if dateOfBirth.Valid {
		user.DateOfBirth = &dateOfBirth.Time
	}
	if minBudget.Valid {
		user.MinBudget = &minBudget.Float64
	}
	if maxBudget.Valid {
		user.MaxBudget = &maxBudget.Float64
	}
	if realEstateCompanyID.Valid {
		user.RealEstateCompanyID = &realEstateCompanyID.String
	}

	// Deserialize JSON fields
	if err := json.Unmarshal([]byte(provincesJSON), &user.PreferredProvinces); err != nil {
		return nil, fmt.Errorf("error unmarshaling preferred provinces: %w", err)
	}
	if err := json.Unmarshal([]byte(propertyTypesJSON), &user.PreferredPropertyTypes); err != nil {
		return nil, fmt.Errorf("error unmarshaling preferred property types: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth,
			   user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			   avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			   created_at, updated_at
		FROM users WHERE email = $1`

	user := &domain.User{}
	var dateOfBirth sql.NullTime
	var minBudget, maxBudget sql.NullFloat64
	var provincesJSON, propertyTypesJSON string
	var realEstateCompanyID sql.NullString

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.NationalID,
		&dateOfBirth,
		&user.UserType,
		&user.Active,
		&minBudget,
		&maxBudget,
		&provincesJSON,
		&propertyTypesJSON,
		&user.AvatarURL,
		&user.Bio,
		&realEstateCompanyID,
		&user.ReceiveNotifications,
		&user.ReceiveNewsletter,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Handle nullable fields (same logic as GetByID)
	if dateOfBirth.Valid {
		user.DateOfBirth = &dateOfBirth.Time
	}
	if minBudget.Valid {
		user.MinBudget = &minBudget.Float64
	}
	if maxBudget.Valid {
		user.MaxBudget = &maxBudget.Float64
	}
	if realEstateCompanyID.Valid {
		user.RealEstateCompanyID = &realEstateCompanyID.String
	}

	// Deserialize JSON fields
	if err := json.Unmarshal([]byte(provincesJSON), &user.PreferredProvinces); err != nil {
		return nil, fmt.Errorf("error unmarshaling preferred provinces: %w", err)
	}
	if err := json.Unmarshal([]byte(propertyTypesJSON), &user.PreferredPropertyTypes); err != nil {
		return nil, fmt.Errorf("error unmarshaling preferred property types: %w", err)
	}

	return user, nil
}

// GetByNationalID retrieves a user by their national ID
func (r *UserRepository) GetByNationalID(nationalID string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth,
			   user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			   avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			   created_at, updated_at
		FROM users WHERE national_id = $1`

	user := &domain.User{}
	var dateOfBirth sql.NullTime
	var minBudget, maxBudget sql.NullFloat64
	var provincesJSON, propertyTypesJSON string
	var realEstateCompanyID sql.NullString

	err := r.db.QueryRow(query, nationalID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.NationalID,
		&dateOfBirth,
		&user.UserType,
		&user.Active,
		&minBudget,
		&maxBudget,
		&provincesJSON,
		&propertyTypesJSON,
		&user.AvatarURL,
		&user.Bio,
		&realEstateCompanyID,
		&user.ReceiveNotifications,
		&user.ReceiveNewsletter,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with national ID: %s", nationalID)
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Handle nullable fields (same logic as GetByID)
	if dateOfBirth.Valid {
		user.DateOfBirth = &dateOfBirth.Time
	}
	if minBudget.Valid {
		user.MinBudget = &minBudget.Float64
	}
	if maxBudget.Valid {
		user.MaxBudget = &maxBudget.Float64
	}
	if realEstateCompanyID.Valid {
		user.RealEstateCompanyID = &realEstateCompanyID.String
	}

	// Deserialize JSON fields
	if err := json.Unmarshal([]byte(provincesJSON), &user.PreferredProvinces); err != nil {
		return nil, fmt.Errorf("error unmarshaling preferred provinces: %w", err)
	}
	if err := json.Unmarshal([]byte(propertyTypesJSON), &user.PreferredPropertyTypes); err != nil {
		return nil, fmt.Errorf("error unmarshaling preferred property types: %w", err)
	}

	return user, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth,
			   user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			   avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			   created_at, updated_at
		FROM users 
		ORDER BY first_name, last_name`

	return r.queryUsers(query)
}

// GetByType retrieves users filtered by type
func (r *UserRepository) GetByType(userType string) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth,
			   user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			   avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			   created_at, updated_at
		FROM users 
		WHERE user_type = $1 AND active = true
		ORDER BY first_name, last_name`

	return r.queryUsersWithArgs(query, userType)
}

// GetBuyers retrieves all buyer users
func (r *UserRepository) GetBuyers() ([]*domain.User, error) {
	return r.GetByType(domain.UserTypeBuyer)
}

// GetSellers retrieves all seller users
func (r *UserRepository) GetSellers() ([]*domain.User, error) {
	return r.GetByType(domain.UserTypeSeller)
}

// GetAgents retrieves all agent users
func (r *UserRepository) GetAgents() ([]*domain.User, error) {
	return r.GetByType(domain.UserTypeAgent)
}

// GetAgentsByCompany retrieves agents for a specific real estate company
func (r *UserRepository) GetAgentsByCompany(companyID string) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth,
			   user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			   avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			   created_at, updated_at
		FROM users 
		WHERE user_type = 'agent' AND real_estate_company_id = $1 AND active = true
		ORDER BY first_name, last_name`

	return r.queryUsersWithArgs(query, companyID)
}

// Update updates an existing user
func (r *UserRepository) Update(user *domain.User) error {
	// Update timestamp
	user.UpdateTimestamp()

	// Serialize JSON fields
	provincesJSON, err := json.Marshal(user.PreferredProvinces)
	if err != nil {
		return fmt.Errorf("error marshaling preferred provinces: %w", err)
	}

	propertyTypesJSON, err := json.Marshal(user.PreferredPropertyTypes)
	if err != nil {
		return fmt.Errorf("error marshaling preferred property types: %w", err)
	}

	query := `
		UPDATE users SET
			first_name = $2, last_name = $3, email = $4, phone = $5, national_id = $6,
			date_of_birth = $7, user_type = $8, active = $9, min_budget = $10, max_budget = $11,
			preferred_provinces = $12, preferred_property_types = $13, avatar_url = $14, bio = $15,
			real_estate_company_id = $16, receive_notifications = $17, receive_newsletter = $18,
			updated_at = $19
		WHERE id = $1`

	result, err := r.db.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.NationalID,
		nullTimePtr(user.DateOfBirth),
		user.UserType,
		user.Active,
		nullFloat64Ptr(user.MinBudget),
		nullFloat64Ptr(user.MaxBudget),
		string(provincesJSON),
		string(propertyTypesJSON),
		user.AvatarURL,
		user.Bio,
		nullStringPtr(user.RealEstateCompanyID),
		user.ReceiveNotifications,
		user.ReceiveNewsletter,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID: %s", user.ID)
	}

	return nil
}

// Deactivate marks a user as inactive
func (r *UserRepository) Deactivate(id string) error {
	query := `UPDATE users SET active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deactivating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID: %s", id)
	}

	return nil
}

// Activate marks a user as active
func (r *UserRepository) Activate(id string) error {
	query := `UPDATE users SET active = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error activating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID: %s", id)
	}

	return nil
}

// SearchByName searches users by name using full-text search
func (r *UserRepository) SearchByName(searchTerm string) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth,
			   user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			   avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			   created_at, updated_at
		FROM users 
		WHERE active = true
		  AND to_tsvector('english', first_name || ' ' || last_name || ' ' || bio) @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(to_tsvector('english', first_name || ' ' || last_name || ' ' || bio), plainto_tsquery('english', $1)) DESC`

	return r.queryUsersWithArgs(query, searchTerm)
}

// GetBuyersForProperty gets buyers that can afford a specific property price
func (r *UserRepository) GetBuyersForProperty(propertyPrice float64) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth,
			   user_type, active, min_budget, max_budget, preferred_provinces, preferred_property_types,
			   avatar_url, bio, real_estate_company_id, receive_notifications, receive_newsletter,
			   created_at, updated_at
		FROM users 
		WHERE user_type = 'buyer'
		  AND active = true
		  AND (min_budget IS NULL OR $1 >= min_budget)
		  AND (max_budget IS NULL OR $1 <= max_budget)
		ORDER BY max_budget DESC NULLS LAST, created_at DESC`

	return r.queryUsersWithArgs(query, propertyPrice)
}

// ExistsByEmail checks if a user with the given email already exists
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}

	return count > 0, nil
}

// ExistsByNationalID checks if a user with the given national ID already exists
func (r *UserRepository) ExistsByNationalID(nationalID string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE national_id = $1`

	var count int
	err := r.db.QueryRow(query, nationalID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking national ID existence: %w", err)
	}

	return count > 0, nil
}

// ValidateNationalID validates if a national ID is valid using the database function
func (r *UserRepository) ValidateNationalID(nationalID string) (bool, error) {
	query := `SELECT validate_ecuador_national_id($1)`

	var isValid bool
	err := r.db.QueryRow(query, nationalID).Scan(&isValid)
	if err != nil {
		return false, fmt.Errorf("error validating national ID: %w", err)
	}

	return isValid, nil
}

// GetStatistics returns statistics about users
func (r *UserRepository) GetStatistics() (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_users,
			COUNT(*) FILTER (WHERE active = true) as active_users,
			COUNT(*) FILTER (WHERE user_type = 'buyer' AND active = true) as buyers,
			COUNT(*) FILTER (WHERE user_type = 'seller' AND active = true) as sellers,
			COUNT(*) FILTER (WHERE user_type = 'agent' AND active = true) as agents,
			COUNT(*) FILTER (WHERE user_type = 'admin' AND active = true) as admins,
			COUNT(*) FILTER (WHERE user_type = 'buyer' AND active = true AND (min_budget IS NOT NULL OR max_budget IS NOT NULL)) as users_with_budget
		FROM users`

	var totalUsers, activeUsers, buyers, sellers, agents, admins, usersWithBudget int64

	err := r.db.QueryRow(query).Scan(&totalUsers, &activeUsers, &buyers, &sellers, &agents, &admins, &usersWithBudget)
	if err != nil {
		return nil, fmt.Errorf("error getting user statistics: %w", err)
	}

	stats := map[string]interface{}{
		"total_users":       totalUsers,
		"active_users":      activeUsers,
		"buyers":            buyers,
		"sellers":           sellers,
		"agents":            agents,
		"admins":            admins,
		"users_with_budget": usersWithBudget,
	}

	return stats, nil
}

// Delete permanently deletes a user (use with caution)
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID: %s", id)
	}

	return nil
}

// Helper methods for querying users
func (r *UserRepository) queryUsers(query string) ([]*domain.User, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

func (r *UserRepository) queryUsersWithArgs(query string, args ...interface{}) ([]*domain.User, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

func (r *UserRepository) scanUsers(rows *sql.Rows) ([]*domain.User, error) {
	var users []*domain.User

	for rows.Next() {
		user := &domain.User{}
		var dateOfBirth sql.NullTime
		var minBudget, maxBudget sql.NullFloat64
		var provincesJSON, propertyTypesJSON string
		var realEstateCompanyID sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.NationalID,
			&dateOfBirth,
			&user.UserType,
			&user.Active,
			&minBudget,
			&maxBudget,
			&provincesJSON,
			&propertyTypesJSON,
			&user.AvatarURL,
			&user.Bio,
			&realEstateCompanyID,
			&user.ReceiveNotifications,
			&user.ReceiveNewsletter,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}

		// Handle nullable fields
		if dateOfBirth.Valid {
			user.DateOfBirth = &dateOfBirth.Time
		}
		if minBudget.Valid {
			user.MinBudget = &minBudget.Float64
		}
		if maxBudget.Valid {
			user.MaxBudget = &maxBudget.Float64
		}
		if realEstateCompanyID.Valid {
			user.RealEstateCompanyID = &realEstateCompanyID.String
		}

		// Deserialize JSON fields
		if err := json.Unmarshal([]byte(provincesJSON), &user.PreferredProvinces); err != nil {
			return nil, fmt.Errorf("error unmarshaling preferred provinces: %w", err)
		}
		if err := json.Unmarshal([]byte(propertyTypesJSON), &user.PreferredPropertyTypes); err != nil {
			return nil, fmt.Errorf("error unmarshaling preferred property types: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

// Helper function for handling null time values
func nullTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
