package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"realty-core/internal/domain"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (
			id, first_name, last_name, email, phone, national_id, date_of_birth, 
			user_type, active, min_budget, max_budget, preferred_provinces, 
			preferred_property_types, avatar_url, bio, real_estate_company_id,
			receive_notifications, receive_newsletter, agency_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)`

	_, err := r.db.Exec(query,
		user.ID, user.FirstName, user.LastName, user.Email, user.Phone, user.Cedula,
		user.DateOfBirth, user.Role, user.Active, user.MinBudget, user.MaxBudget,
		pq.Array(user.PreferredProvinces), pq.Array(user.PreferredPropertyTypes),
		user.AvatarURL, user.Bio, user.RealEstateCompanyID,
		user.ReceiveNotifications, user.ReceiveNewsletter, user.AgencyID,
		user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth, 
			   user_type, active, min_budget, max_budget, preferred_provinces, 
			   preferred_property_types, avatar_url, bio, real_estate_company_id,
			   receive_notifications, receive_newsletter, agency_id, created_at, updated_at
		FROM users 
		WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone,
		&user.Cedula, &user.DateOfBirth, &user.Role, &user.Active,
		&user.MinBudget, &user.MaxBudget, pq.Array(&user.PreferredProvinces),
		pq.Array(&user.PreferredPropertyTypes), &user.AvatarURL, &user.Bio,
		&user.RealEstateCompanyID, &user.ReceiveNotifications, &user.ReceiveNewsletter,
		&user.AgencyID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with id: %s", id)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth, 
			   user_type, active, min_budget, max_budget, preferred_provinces, 
			   preferred_property_types, avatar_url, bio, real_estate_company_id,
			   receive_notifications, receive_newsletter, agency_id, password_hash, created_at, updated_at
		FROM users 
		WHERE email = $1`

	user := &domain.User{}
	var provincesJSON, propertyTypesJSON []byte
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone,
		&user.Cedula, &user.DateOfBirth, &user.Role, &user.Active,
		&user.MinBudget, &user.MaxBudget, &provincesJSON,
		&propertyTypesJSON, &user.AvatarURL, &user.Bio,
		&user.RealEstateCompanyID, &user.ReceiveNotifications, &user.ReceiveNewsletter,
		&user.AgencyID, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == nil {
		// Parse JSONB arrays after successful scan
		if len(provincesJSON) > 0 {
			json.Unmarshal(provincesJSON, &user.PreferredProvinces)
		}
		if len(propertyTypesJSON) > 0 {
			json.Unmarshal(propertyTypesJSON, &user.PreferredPropertyTypes)
		}
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetByNationalID retrieves a user by national ID
func (r *UserRepository) GetByNationalID(national_id string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth, 
			   user_type, active, min_budget, max_budget, preferred_provinces, 
			   preferred_property_types, avatar_url, bio, real_estate_company_id,
			   receive_notifications, receive_newsletter, agency_id, created_at, updated_at
		FROM users 
		WHERE national_id = $1`

	user := &domain.User{}
	err := r.db.QueryRow(query, national_id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone,
		&user.Cedula, &user.DateOfBirth, &user.Role, &user.Active,
		&user.MinBudget, &user.MaxBudget, pq.Array(&user.PreferredProvinces),
		pq.Array(&user.PreferredPropertyTypes), &user.AvatarURL, &user.Bio,
		&user.RealEstateCompanyID, &user.ReceiveNotifications, &user.ReceiveNewsletter,
		&user.AgencyID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with national ID: %s", national_id)
		}
		return nil, fmt.Errorf("failed to get user by national ID: %w", err)
	}

	return user, nil
}

// Update updates a user in the database
func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users SET 
			first_name = $2, last_name = $3, email = $4, phone = $5, 
			national_id = $6, date_of_birth = $7, user_type = $8, active = $9, 
			min_budget = $10, max_budget = $11, preferred_provinces = $12, 
			preferred_property_types = $13, avatar_url = $14, bio = $15, 
			real_estate_company_id = $16, receive_notifications = $17, 
			receive_newsletter = $18, agency_id = $19, updated_at = $20
		WHERE id = $1`

	_, err := r.db.Exec(query,
		user.ID, user.FirstName, user.LastName, user.Email, user.Phone,
		user.Cedula, user.DateOfBirth, user.Role, user.Active,
		user.MinBudget, user.MaxBudget, pq.Array(user.PreferredProvinces),
		pq.Array(user.PreferredPropertyTypes), user.AvatarURL, user.Bio,
		user.RealEstateCompanyID, user.ReceiveNotifications, user.ReceiveNewsletter,
		user.AgencyID, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user from the database
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetByUserType retrieves users by role
func (r *UserRepository) GetByUserType(role domain.UserRole) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth, 
			   user_type, active, min_budget, max_budget, preferred_provinces, 
			   preferred_property_types, avatar_url, bio, real_estate_company_id,
			   receive_notifications, receive_newsletter, agency_id, created_at, updated_at
		FROM users 
		WHERE user_type = $1 AND active = TRUE
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone,
			&user.Cedula, &user.DateOfBirth, &user.Role, &user.Active,
			&user.MinBudget, &user.MaxBudget, pq.Array(&user.PreferredProvinces),
			pq.Array(&user.PreferredPropertyTypes), &user.AvatarURL, &user.Bio,
			&user.RealEstateCompanyID, &user.ReceiveNotifications, &user.ReceiveNewsletter,
			&user.AgencyID, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// GetByAgency retrieves users by agency ID
func (r *UserRepository) GetByAgency(agencyID string) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth, 
			   user_type, active, min_budget, max_budget, preferred_provinces, 
			   preferred_property_types, avatar_url, bio, real_estate_company_id,
			   receive_notifications, receive_newsletter, agency_id, created_at, updated_at
		FROM users 
		WHERE agency_id = $1 AND active = TRUE
		ORDER BY first_name, last_name`

	rows, err := r.db.Query(query, agencyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by agency: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone,
			&user.Cedula, &user.DateOfBirth, &user.Role, &user.Active,
			&user.MinBudget, &user.MaxBudget, pq.Array(&user.PreferredProvinces),
			pq.Array(&user.PreferredPropertyTypes), &user.AvatarURL, &user.Bio,
			&user.RealEstateCompanyID, &user.ReceiveNotifications, &user.ReceiveNewsletter,
			&user.AgencyID, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Search searches users with filters
func (r *UserRepository) Search(params *domain.UserSearchParams) ([]*domain.User, int, error) {
	// Build base query
	baseQuery := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth, 
			   user_type, active, min_budget, max_budget, preferred_provinces, 
			   preferred_property_types, avatar_url, bio, real_estate_company_id,
			   receive_notifications, receive_newsletter, agency_id, created_at, updated_at
		FROM users WHERE 1=1`

	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`

	var args []interface{}
	var conditions []string
	argIndex := 1

	// Add search conditions
	if params.Query != "" {
		conditions = append(conditions, fmt.Sprintf(`to_tsvector('spanish', first_name || ' ' || last_name) @@ plainto_tsquery('spanish', $%d)`, argIndex))
		args = append(args, params.Query)
		argIndex++
	}

	if params.Role != nil {
		conditions = append(conditions, fmt.Sprintf(`user_type = $%d`, argIndex))
		args = append(args, *params.Role)
		argIndex++
	}

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf(`active = $%d`, argIndex))
		args = append(args, *params.Active)
		argIndex++
	}

	if params.AgencyID != nil {
		conditions = append(conditions, fmt.Sprintf(`agency_id = $%d`, argIndex))
		args = append(args, *params.AgencyID)
		argIndex++
	}

	if len(params.Provinces) > 0 {
		conditions = append(conditions, fmt.Sprintf(`preferred_provinces && $%d`, argIndex))
		args = append(args, pq.Array(params.Provinces))
		argIndex++
	}

	if params.MinBudget != nil {
		conditions = append(conditions, fmt.Sprintf(`min_budget >= $%d`, argIndex))
		args = append(args, *params.MinBudget)
		argIndex++
	}

	if params.MaxBudget != nil {
		conditions = append(conditions, fmt.Sprintf(`max_budget <= $%d`, argIndex))
		args = append(args, *params.MaxBudget)
		argIndex++
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
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
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
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone,
			&user.Cedula, &user.DateOfBirth, &user.Role, &user.Active,
			&user.MinBudget, &user.MaxBudget, pq.Array(&user.PreferredProvinces),
			pq.Array(&user.PreferredPropertyTypes), &user.AvatarURL, &user.Bio,
			&user.RealEstateCompanyID, &user.ReceiveNotifications, &user.ReceiveNewsletter,
			&user.AgencyID, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, totalCount, nil
}

// GetBuyersByBudget finds buyers who can afford a specific property price
func (r *UserRepository) GetBuyersByBudget(price float64) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, national_id, date_of_birth, 
			   user_type, active, min_budget, max_budget, preferred_provinces, 
			   preferred_property_types, avatar_url, bio, real_estate_company_id,
			   receive_notifications, receive_newsletter, agency_id, created_at, updated_at
		FROM users 
		WHERE user_type = 'buyer' AND active = TRUE 
		  AND min_budget IS NOT NULL AND max_budget IS NOT NULL
		  AND $1 >= min_budget AND $1 <= max_budget
		ORDER BY max_budget DESC`

	rows, err := r.db.Query(query, price)
	if err != nil {
		return nil, fmt.Errorf("failed to get buyers by budget: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone,
			&user.Cedula, &user.DateOfBirth, &user.Role, &user.Active,
			&user.MinBudget, &user.MaxBudget, pq.Array(&user.PreferredProvinces),
			pq.Array(&user.PreferredPropertyTypes), &user.AvatarURL, &user.Bio,
			&user.RealEstateCompanyID, &user.ReceiveNotifications, &user.ReceiveNewsletter,
			&user.AgencyID, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// GetStatistics returns user statistics
func (r *UserRepository) GetStatistics() (*domain.UserStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_users,
			COUNT(*) FILTER (WHERE active = TRUE) as active_users,
			COUNT(*) FILTER (WHERE user_type = 'admin') as admin_count,
			COUNT(*) FILTER (WHERE user_type = 'agency') as agency_count,
			COUNT(*) FILTER (WHERE user_type = 'agent') as agent_count,
			COUNT(*) FILTER (WHERE user_type = 'seller') as owner_count,
			COUNT(*) FILTER (WHERE user_type = 'buyer') as buyer_count,
			COUNT(*) FILTER (WHERE min_budget IS NOT NULL AND max_budget IS NOT NULL) as with_budget,
			COUNT(*) FILTER (WHERE agency_id IS NOT NULL) as associated_agents
		FROM users`

	stats := &domain.UserStats{}
	err := r.db.QueryRow(query).Scan(
		&stats.TotalUsers, &stats.ActiveUsers, &stats.AdminCount,
		&stats.AgencyCount, &stats.AgentCount, &stats.OwnerCount,
		&stats.BuyerCount, &stats.WithBudget, &stats.AssociatedAgents,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %w", err)
	}

	return stats, nil
}

// SetEmailVerified sets the email verification status (placeholder - auth fields not in schema)
func (r *UserRepository) SetEmailVerified(userID string, verified bool) error {
	// Note: email_verified field does not exist in the actual database schema
	// This is a placeholder for future authentication implementation
	query := `
		UPDATE users 
		SET updated_at = $2
		WHERE id = $1`

	_, err := r.db.Exec(query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// SetPasswordResetToken sets the password reset token and expiry (placeholder - auth fields not in schema)
func (r *UserRepository) SetPasswordResetToken(userID, token string, expiry time.Time) error {
	// Note: password_reset_token fields do not exist in the actual database schema
	// This is a placeholder for future authentication implementation
	query := `
		UPDATE users 
		SET updated_at = $2
		WHERE id = $1`

	_, err := r.db.Exec(query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// ClearPasswordResetToken clears the password reset token (placeholder - auth fields not in schema)
func (r *UserRepository) ClearPasswordResetToken(userID string) error {
	// Note: password_reset_token fields do not exist in the actual database schema
	// This is a placeholder for future authentication implementation
	query := `
		UPDATE users 
		SET updated_at = $2
		WHERE id = $1`

	_, err := r.db.Exec(query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(userID string) error {
	now := time.Now()
	query := `
		UPDATE users 
		SET updated_at = $2
		WHERE id = $1`

	_, err := r.db.Exec(query, userID, now)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}