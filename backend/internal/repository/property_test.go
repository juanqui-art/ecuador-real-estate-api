package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realty-core/internal/domain"
)

// Helper function to create a test property
func createTestProperty() *domain.Property {
	return domain.NewProperty(
		"Beautiful house in Samborondón",
		"Modern house with pool",
		"Guayas",
		"Samborondón",
		"house",
		285000,
		"owner-123",
	)
}

// Helper function to setup mock database
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func TestNewPostgreSQLPropertyRepository(t *testing.T) {
	db, _ := setupMockDB(t)
	defer db.Close()

	repo := NewPostgreSQLPropertyRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestPostgreSQLPropertyRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		property      *domain.Property
		mockSetup     func(sqlmock.Sqlmock)
		wantError     bool
		errorContains string
	}{
		{
			name:     "successful creation",
			property: createTestProperty(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO properties`).
					WithArgs(
						sqlmock.AnyArg(), // id
						sqlmock.AnyArg(), // slug
						sqlmock.AnyArg(), // title
						sqlmock.AnyArg(), // description
						sqlmock.AnyArg(), // price
						sqlmock.AnyArg(), // province
						sqlmock.AnyArg(), // city
						sqlmock.AnyArg(), // sector
						sqlmock.AnyArg(), // address
						sqlmock.AnyArg(), // latitude
						sqlmock.AnyArg(), // longitude
						sqlmock.AnyArg(), // location_precision
						sqlmock.AnyArg(), // type
						sqlmock.AnyArg(), // status
						sqlmock.AnyArg(), // bedrooms
						sqlmock.AnyArg(), // bathrooms
						sqlmock.AnyArg(), // area_m2
						sqlmock.AnyArg(), // main_image
						sqlmock.AnyArg(), // images
						sqlmock.AnyArg(), // video_tour
						sqlmock.AnyArg(), // tour_360
						sqlmock.AnyArg(), // rent_price
						sqlmock.AnyArg(), // common_expenses
						sqlmock.AnyArg(), // price_per_m2
						sqlmock.AnyArg(), // year_built
						sqlmock.AnyArg(), // floors
						sqlmock.AnyArg(), // property_status
						sqlmock.AnyArg(), // furnished
						sqlmock.AnyArg(), // garage
						sqlmock.AnyArg(), // pool
						sqlmock.AnyArg(), // garden
						sqlmock.AnyArg(), // terrace
						sqlmock.AnyArg(), // balcony
						sqlmock.AnyArg(), // security
						sqlmock.AnyArg(), // elevator
						sqlmock.AnyArg(), // air_conditioning
						sqlmock.AnyArg(), // tags
						sqlmock.AnyArg(), // featured
						sqlmock.AnyArg(), // view_count
						sqlmock.AnyArg(), // real_estate_company_id
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
						sqlmock.AnyArg(), // parking_spaces
						sqlmock.AnyArg(), // owner_id
						sqlmock.AnyArg(), // agent_id
						sqlmock.AnyArg(), // agency_id
						sqlmock.AnyArg(), // created_by
						sqlmock.AnyArg(), // updated_by
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantError: false,
		},
		{
			name:     "database error",
			property: createTestProperty(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO properties`).
					WillReturnError(errors.New("database connection failed"))
			},
			wantError:     true,
			errorContains: "error creating property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewPostgreSQLPropertyRepository(db)

			err := repo.Create(tt.property)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLPropertyRepository_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockSetup     func(sqlmock.Sqlmock)
		wantError     bool
		errorContains string
		validateResult func(*domain.Property) bool
	}{
		{
			name: "successful retrieval",
			id:   "test-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city",
					"sector", "address", "latitude", "longitude", "location_precision",
					"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
					"images", "video_tour", "tour_360", "rent_price", "common_expenses",
					"price_per_m2", "year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
					"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				}).AddRow(
					"test-id", "test-slug", "Test Title", "Test Description", 100000.0, "Guayas", "Samborondón",
					nil, nil, nil, nil, "approximate", "house", "available", 3, 2.5, 150.0, nil,
					`[]`, nil, nil, nil, nil, nil, nil, nil, "used", false, false, false, false,
					false, false, false, false, false, `[]`, false, 0, nil, time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil,
				)
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE id = \$1`).
					WithArgs("test-id").
					WillReturnRows(rows)
			},
			wantError: false,
			validateResult: func(p *domain.Property) bool {
				return p.ID == "test-id" && p.Title == "Test Title"
			},
		},
		{
			name: "property not found",
			id:   "nonexistent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE id = \$1`).
					WithArgs("nonexistent-id").
					WillReturnError(sql.ErrNoRows)
			},
			wantError:     true,
			errorContains: "property not found",
		},
		{
			name: "database error",
			id:   "test-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE id = \$1`).
					WithArgs("test-id").
					WillReturnError(errors.New("database connection failed"))
			},
			wantError:     true,
			errorContains: "error retrieving property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewPostgreSQLPropertyRepository(db)

			property, err := repo.GetByID(tt.id)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, property)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, property)
				if tt.validateResult != nil {
					assert.True(t, tt.validateResult(property))
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLPropertyRepository_GetBySlug(t *testing.T) {
	tests := []struct {
		name          string
		slug          string
		mockSetup     func(sqlmock.Sqlmock)
		wantError     bool
		errorContains string
		validateResult func(*domain.Property) bool
	}{
		{
			name: "successful retrieval",
			slug: "test-slug-12345678",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city",
					"sector", "address", "latitude", "longitude", "location_precision",
					"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
					"images", "video_tour", "tour_360", "rent_price", "common_expenses",
					"price_per_m2", "year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
					"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				}).AddRow(
					"test-id", "test-slug-12345678", "Test Title", "Test Description", 100000.0, "Guayas", "Samborondón",
					nil, nil, nil, nil, "approximate", "house", "available", 3, 2.5, 150.0, nil,
					`[]`, nil, nil, nil, nil, nil, nil, nil, "used", false, false, false, false,
					false, false, false, false, false, `[]`, false, 0, nil, time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil, 0,
					nil, nil, nil, nil, nil,
				)
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE slug = \$1`).
					WithArgs("test-slug-12345678").
					WillReturnRows(rows)
			},
			wantError: false,
			validateResult: func(p *domain.Property) bool {
				return p.Slug == "test-slug-12345678" && p.Title == "Test Title"
			},
		},
		{
			name: "property not found",
			slug: "nonexistent-slug",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE slug = \$1`).
					WithArgs("nonexistent-slug").
					WillReturnError(sql.ErrNoRows)
			},
			wantError:     true,
			errorContains: "property not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewPostgreSQLPropertyRepository(db)

			property, err := repo.GetBySlug(tt.slug)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, property)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, property)
				if tt.validateResult != nil {
					assert.True(t, tt.validateResult(property))
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLPropertyRepository_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		wantError     bool
		errorContains string
		expectedCount int
	}{
		{
			name: "successful retrieval with multiple properties",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city",
					"sector", "address", "latitude", "longitude", "location_precision",
					"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
					"images", "video_tour", "tour_360", "rent_price", "common_expenses",
					"price_per_m2", "year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
					"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				}).AddRow(
					"id1", "slug1", "Title 1", "Description 1", 100000.0, "Guayas", "Samborondón",
					nil, nil, nil, nil, "approximate", "house", "available", 3, 2.5, 150.0, nil,
					`[]`, nil, nil, nil, nil, nil, nil, nil, "used", false, false, false, false,
					false, false, false, false, false, `[]`, false, 0, nil, time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil,
				).AddRow(
					"id2", "slug2", "Title 2", "Description 2", 200000.0, "Pichincha", "Quito",
					nil, nil, nil, nil, "approximate", "apartment", "available", 2, 2.0, 80.0, nil,
					`[]`, nil, nil, nil, nil, nil, nil, nil, "used", false, false, false, false,
					false, false, false, false, false, `[]`, false, 0, nil, time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil,
				)
				mock.ExpectQuery(`SELECT .+ FROM properties ORDER BY featured DESC, created_at DESC`).
					WillReturnRows(rows)
			},
			wantError:     false,
			expectedCount: 2,
		},
		{
			name: "empty result set",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city",
					"sector", "address", "latitude", "longitude", "location_precision",
					"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
					"images", "video_tour", "tour_360", "rent_price", "common_expenses",
					"price_per_m2", "year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
					"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				})
				mock.ExpectQuery(`SELECT .+ FROM properties ORDER BY featured DESC, created_at DESC`).
					WillReturnRows(rows)
			},
			wantError:     false,
			expectedCount: 0,
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM properties ORDER BY featured DESC, created_at DESC`).
					WillReturnError(errors.New("database connection failed"))
			},
			wantError:     true,
			errorContains: "error querying properties",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewPostgreSQLPropertyRepository(db)

			properties, err := repo.GetAll()

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, properties)
			} else {
				assert.NoError(t, err)
				assert.Len(t, properties, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLPropertyRepository_Update(t *testing.T) {
	tests := []struct {
		name          string
		property      *domain.Property
		mockSetup     func(sqlmock.Sqlmock)
		wantError     bool
		errorContains string
	}{
		{
			name:     "successful update",
			property: createTestProperty(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE properties SET`).
					WithArgs(
						sqlmock.AnyArg(), // slug
						sqlmock.AnyArg(), // title
						sqlmock.AnyArg(), // description
						sqlmock.AnyArg(), // price
						sqlmock.AnyArg(), // province
						sqlmock.AnyArg(), // city
						sqlmock.AnyArg(), // sector
						sqlmock.AnyArg(), // address
						sqlmock.AnyArg(), // latitude
						sqlmock.AnyArg(), // longitude
						sqlmock.AnyArg(), // location_precision
						sqlmock.AnyArg(), // type
						sqlmock.AnyArg(), // status
						sqlmock.AnyArg(), // bedrooms
						sqlmock.AnyArg(), // bathrooms
						sqlmock.AnyArg(), // area_m2
						sqlmock.AnyArg(), // main_image
						sqlmock.AnyArg(), // images
						sqlmock.AnyArg(), // video_tour
						sqlmock.AnyArg(), // tour_360
						sqlmock.AnyArg(), // rent_price
						sqlmock.AnyArg(), // common_expenses
						sqlmock.AnyArg(), // price_per_m2
						sqlmock.AnyArg(), // year_built
						sqlmock.AnyArg(), // floors
						sqlmock.AnyArg(), // property_status
						sqlmock.AnyArg(), // furnished
						sqlmock.AnyArg(), // garage
						sqlmock.AnyArg(), // pool
						sqlmock.AnyArg(), // garden
						sqlmock.AnyArg(), // terrace
						sqlmock.AnyArg(), // balcony
						sqlmock.AnyArg(), // security
						sqlmock.AnyArg(), // elevator
						sqlmock.AnyArg(), // air_conditioning
						sqlmock.AnyArg(), // tags
						sqlmock.AnyArg(), // featured
						sqlmock.AnyArg(), // view_count
						sqlmock.AnyArg(), // real_estate_company_id
						sqlmock.AnyArg(), // updated_at
						sqlmock.AnyArg(), // id (WHERE clause)
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			name:     "property not found",
			property: createTestProperty(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE properties SET`).
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
			wantError:     true,
			errorContains: "property not found",
		},
		{
			name:     "database error",
			property: createTestProperty(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE properties SET`).
					WillReturnError(errors.New("database connection failed"))
			},
			wantError:     true,
			errorContains: "error updating property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewPostgreSQLPropertyRepository(db)

			err := repo.Update(tt.property)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLPropertyRepository_Delete(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockSetup     func(sqlmock.Sqlmock)
		wantError     bool
		errorContains string
	}{
		{
			name: "successful deletion",
			id:   "test-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM properties WHERE id = \$1`).
					WithArgs("test-id").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			name: "property not found",
			id:   "nonexistent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM properties WHERE id = \$1`).
					WithArgs("nonexistent-id").
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
			wantError:     true,
			errorContains: "property not found",
		},
		{
			name: "database error",
			id:   "test-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM properties WHERE id = \$1`).
					WithArgs("test-id").
					WillReturnError(errors.New("database connection failed"))
			},
			wantError:     true,
			errorContains: "error deleting property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewPostgreSQLPropertyRepository(db)

			err := repo.Delete(tt.id)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLPropertyRepository_GetByProvince(t *testing.T) {
	tests := []struct {
		name          string
		province      string
		mockSetup     func(sqlmock.Sqlmock)
		wantError     bool
		errorContains string
		expectedCount int
	}{
		{
			name:     "successful filtering",
			province: "Guayas",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city",
					"sector", "address", "latitude", "longitude", "location_precision",
					"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
					"images", "video_tour", "tour_360", "rent_price", "common_expenses",
					"price_per_m2", "year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
					"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				}).AddRow(
					"id1", "slug1", "Title 1", "Description 1", 100000.0, "Guayas", "Samborondón",
					nil, nil, nil, nil, "approximate", "house", "available", 3, 2.5, 150.0, nil,
					`[]`, nil, nil, nil, nil, nil, nil, nil, "used", false, false, false, false,
					false, false, false, false, false, `[]`, false, 0, nil, time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil,
				)
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE province = \$1 ORDER BY featured DESC, created_at DESC`).
					WithArgs("Guayas").
					WillReturnRows(rows)
			},
			wantError:     false,
			expectedCount: 1,
		},
		{
			name:     "no properties found",
			province: "Loja",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city",
					"sector", "address", "latitude", "longitude", "location_precision",
					"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
					"images", "video_tour", "tour_360", "rent_price", "common_expenses",
					"price_per_m2", "year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
					"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				})
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE province = \$1 ORDER BY featured DESC, created_at DESC`).
					WithArgs("Loja").
					WillReturnRows(rows)
			},
			wantError:     false,
			expectedCount: 0,
		},
		{
			name:     "database error",
			province: "Guayas",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE province = \$1 ORDER BY featured DESC, created_at DESC`).
					WithArgs("Guayas").
					WillReturnError(errors.New("database connection failed"))
			},
			wantError:     true,
			errorContains: "error querying properties by province",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewPostgreSQLPropertyRepository(db)

			properties, err := repo.GetByProvince(tt.province)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, properties)
			} else {
				assert.NoError(t, err)
				assert.Len(t, properties, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLPropertyRepository_GetByPriceRange(t *testing.T) {
	tests := []struct {
		name          string
		minPrice      float64
		maxPrice      float64
		mockSetup     func(sqlmock.Sqlmock)
		wantError     bool
		errorContains string
		expectedCount int
	}{
		{
			name:     "successful filtering",
			minPrice: 100000,
			maxPrice: 300000,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city",
					"sector", "address", "latitude", "longitude", "location_precision",
					"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
					"images", "video_tour", "tour_360", "rent_price", "common_expenses",
					"price_per_m2", "year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
					"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				}).AddRow(
					"id1", "slug1", "Title 1", "Description 1", 200000.0, "Guayas", "Samborondón",
					nil, nil, nil, nil, "approximate", "house", "available", 3, 2.5, 150.0, nil,
					`[]`, nil, nil, nil, nil, nil, nil, nil, "used", false, false, false, false,
					false, false, false, false, false, `[]`, false, 0, nil, time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil,
				)
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE price >= \$1 AND price <= \$2 ORDER BY featured DESC, created_at DESC`).
					WithArgs(100000.0, 300000.0).
					WillReturnRows(rows)
			},
			wantError:     false,
			expectedCount: 1,
		},
		{
			name:     "no properties in range",
			minPrice: 500000,
			maxPrice: 1000000,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city",
					"sector", "address", "latitude", "longitude", "location_precision",
					"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
					"images", "video_tour", "tour_360", "rent_price", "common_expenses",
					"price_per_m2", "year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
					"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				})
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE price >= \$1 AND price <= \$2 ORDER BY featured DESC, created_at DESC`).
					WithArgs(500000.0, 1000000.0).
					WillReturnRows(rows)
			},
			wantError:     false,
			expectedCount: 0,
		},
		{
			name:     "database error",
			minPrice: 100000,
			maxPrice: 300000,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM properties WHERE price >= \$1 AND price <= \$2 ORDER BY featured DESC, created_at DESC`).
					WithArgs(100000.0, 300000.0).
					WillReturnError(errors.New("database connection failed"))
			},
			wantError:     true,
			errorContains: "error querying properties by price range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewPostgreSQLPropertyRepository(db)

			properties, err := repo.GetByPriceRange(tt.minPrice, tt.maxPrice)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, properties)
			} else {
				assert.NoError(t, err)
				assert.Len(t, properties, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Test helper function for scanning properties
func TestScanProperty(t *testing.T) {
	// This tests the scanProperty helper function indirectly through GetByID
	db, mock := setupMockDB(t)
	defer db.Close()

	// Test with valid JSON data
	rows := sqlmock.NewRows([]string{
		"id", "slug", "title", "description", "price", "province", "city",
		"sector", "address", "latitude", "longitude", "location_precision",
		"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
		"images", "video_tour", "tour_360", "rent_price", "common_expenses",
		"price_per_m2", "year_built", "floors", "property_status", "furnished",
		"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
		"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
		"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
	}).AddRow(
		"test-id", "test-slug", "Test Title", "Test Description", 100000.0, "Guayas", "Samborondón",
		"Test Sector", "Test Address", -2.1667, -79.9, "exact", "house", "available", 3, 2.5, 150.0, "main.jpg",
		`["image1.jpg","image2.jpg"]`, "video.mp4", "tour360.html", 1200.0, 150.0, 666.67, 2020, 2, "new", true,
		true, true, true, true, true, true, true, true, `["luxury","pool","garden"]`, true, 25, "company-id", time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil,
	)

	mock.ExpectQuery(`SELECT .+ FROM properties WHERE id = \$1`).
		WithArgs("test-id").
		WillReturnRows(rows)

	repo := NewPostgreSQLPropertyRepository(db)
	property, err := repo.GetByID("test-id")

	assert.NoError(t, err)
	assert.NotNil(t, property)
	assert.Equal(t, "test-id", property.ID)
	assert.Equal(t, "Test Title", property.Title)
	assert.Len(t, property.Images, 2)
	assert.Equal(t, "image1.jpg", property.Images[0])
	assert.Len(t, property.Tags, 3)
	assert.Contains(t, property.Tags, "luxury")
	assert.Equal(t, -2.1667, *property.Latitude)
	assert.Equal(t, -79.9, *property.Longitude)
	assert.Equal(t, "Test Sector", *property.Sector)
	assert.Equal(t, true, property.Furnished)
	assert.Equal(t, true, property.Pool)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test JSON marshaling/unmarshaling edge cases
func TestJSONHandling(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	// Test with invalid JSON for images
	rows := sqlmock.NewRows([]string{
		"id", "slug", "title", "description", "price", "province", "city",
		"sector", "address", "latitude", "longitude", "location_precision",
		"type", "status", "bedrooms", "bathrooms", "area_m2", "main_image",
		"images", "video_tour", "tour_360", "rent_price", "common_expenses",
		"price_per_m2", "year_built", "floors", "property_status", "furnished",
		"garage", "pool", "garden", "terrace", "balcony", "security", "elevator",
		"air_conditioning", "tags", "featured", "view_count", "real_estate_company_id",
		"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
	}).AddRow(
		"test-id", "test-slug", "Test Title", "Test Description", 100000.0, "Guayas", "Samborondón",
		nil, nil, nil, nil, "approximate", "house", "available", 3, 2.5, 150.0, nil,
		`invalid json`, nil, nil, nil, nil, nil, nil, nil, "used", false, false, false, false,
		false, false, false, false, false, `[]`, false, 0, nil, time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil,
	)

	mock.ExpectQuery(`SELECT .+ FROM properties WHERE id = \$1`).
		WithArgs("test-id").
		WillReturnRows(rows)

	repo := NewPostgreSQLPropertyRepository(db)
	property, err := repo.GetByID("test-id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error converting images from JSON")
	assert.Nil(t, property)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test property creation with complex data
func TestCreatePropertyWithComplexData(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	property := createTestProperty()
	// Add complex data
	property.Images = []string{"image1.jpg", "image2.jpg", "image3.jpg"}
	property.Tags = []string{"luxury", "pool", "garden", "modern"}
	latitude := -2.1667
	longitude := -79.9
	property.Latitude = &latitude
	property.Longitude = &longitude
	property.SetFeatured(true)
	property.Pool = true
	property.Garden = true

	mock.ExpectExec(`INSERT INTO properties`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewPostgreSQLPropertyRepository(db)
	err := repo.Create(property)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}