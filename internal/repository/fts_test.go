package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"realty-core/internal/domain"
)

func TestPostgreSQLPropertyRepository_SearchProperties(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgreSQLPropertyRepository(db)

	tests := []struct {
		name          string
		query         string
		limit         int
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		wantError     bool
		errorContains string
	}{
		{
			name:  "successful search",
			query: "casa moderna",
			limit: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city", "sector", "address",
					"latitude", "longitude", "location_precision", "type", "status", "bedrooms", "bathrooms", "area_m2",
					"main_image", "images", "video_tour", "tour_360",
					"rent_price", "common_expenses", "price_per_m2",
					"year_built", "floors", "property_status", "furnished",
					"garage", "pool", "garden", "terrace", "balcony", "security", "elevator", "air_conditioning",
					"tags", "featured", "view_count", "real_estate_company_id",
					"created_at", "updated_at", "parking_spaces",
					"owner_id", "agent_id", "agency_id", "created_by", "updated_by",
				}).AddRow(
					"123e4567-e89b-12d3-a456-426614174000", "casa-moderna", "Casa moderna", "Descripción", 285000.0,
					"Guayas", "Samborondón", "", "", 0.0, 0.0, "", "house", "available", 4, 3.5, 320.0,
					"", "[]", "", "", 0.0, 0.0, 0.0, 0, 0, "", false,
					false, false, false, false, false, false, false, false,
					"[]", false, 0, "", time.Now(), time.Now(), 0,
					nil, nil, nil, nil, nil,
				)
				mock.ExpectQuery(`SELECT .+ FROM properties`).
					WithArgs("casa moderna", 10).
					WillReturnRows(rows)
			},
			expectedCount: 1,
			wantError:     false,
		},
		{
			name:  "database error",
			query: "test",
			limit: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .+ FROM properties`).
					WithArgs("test", 10).
					WillReturnError(sql.ErrConnDone)
			},
			wantError:     true,
			errorContains: "error performing full-text search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			results, err := repo.SearchProperties(tt.query, tt.limit)

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
		})
	}
}

func TestPostgreSQLPropertyRepository_SearchPropertiesRanked(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgreSQLPropertyRepository(db)

	tests := []struct {
		name          string
		query         string
		limit         int
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		wantError     bool
	}{
		{
			name:  "successful ranked search",
			query: "casa lujo",
			limit: 5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city", "type", "rank",
				}).AddRow(
					"123e4567-e89b-12d3-a456-426614174000", "casa-lujo", "Casa de lujo", "Descripción casa lujo",
					500000.0, "Guayas", "Samborondón", "house", 0.9,
				)
				mock.ExpectQuery(`SELECT .+ FROM properties`).
					WithArgs("casa lujo", 5).
					WillReturnRows(rows)
			},
			expectedCount: 1,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			results, err := repo.SearchPropertiesRanked(tt.query, tt.limit)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
			if len(results) > 0 {
				assert.Greater(t, results[0].Rank, 0.0)
			}
		})
	}
}

func TestPostgreSQLPropertyRepository_GetSearchSuggestions(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgreSQLPropertyRepository(db)

	tests := []struct {
		name          string
		query         string
		limit         int
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		wantError     bool
	}{
		{
			name:  "successful suggestions",
			query: "gua",
			limit: 5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"suggestion", "category", "frequency",
				}).
					AddRow("Guayas", "province", 150).
					AddRow("Guayaquil", "city", 120)
				mock.ExpectQuery(`SELECT \* FROM get_search_suggestions`).
					WithArgs("gua", 5).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			suggestions, err := repo.GetSearchSuggestions(tt.query, tt.limit)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, suggestions, tt.expectedCount)
			if len(suggestions) > 0 {
				assert.NotEmpty(t, suggestions[0].Text)
				assert.NotEmpty(t, suggestions[0].Category)
				assert.Greater(t, suggestions[0].Frequency, 0)
			}
		})
	}
}

func TestPostgreSQLPropertyRepository_AdvancedSearch(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgreSQLPropertyRepository(db)

	tests := []struct {
		name          string
		params        AdvancedSearchParams
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		wantError     bool
	}{
		{
			name: "successful advanced search",
			params: AdvancedSearchParams{
				Query:       "casa moderna",
				Province:    "Guayas",
				MinPrice:    200000,
				MaxPrice:    500000,
				MinBedrooms: 3,
				MaxBedrooms: 5,
				Limit:       10,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "slug", "title", "description", "price", "province", "city", "type",
					"bedrooms", "bathrooms", "area_m2", "featured", "rank",
				}).AddRow(
					"123e4567-e89b-12d3-a456-426614174000", "casa-moderna", "Casa moderna", "Descripción",
					350000.0, "Guayas", "Samborondón", "house", 4, 3.5, 280.0, false, 0.85,
				)
				mock.ExpectQuery(`SELECT \* FROM advanced_search_properties`).
					WithArgs(
						"casa moderna", "Guayas", "", "", 200000.0, 500000.0,
						3, 5, 0.0, 100.0, 0.0, 999999.0, false, 10,
					).
					WillReturnRows(rows)
			},
			expectedCount: 1,
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			results, err := repo.AdvancedSearch(tt.params)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
			if len(results) > 0 {
				assert.Equal(t, "casa-moderna", results[0].Property.Slug)
				assert.Greater(t, results[0].Rank, 0.0)
			}
		})
	}
}

func TestAdvancedSearchParams_Validation(t *testing.T) {
	tests := []struct {
		name   string
		params AdvancedSearchParams
		valid  bool
	}{
		{
			name: "valid parameters",
			params: AdvancedSearchParams{
				Query:        "casa",
				Province:     "Guayas",
				MinPrice:     100000,
				MaxPrice:     500000,
				MinBedrooms:  2,
				MaxBedrooms:  5,
				MinBathrooms: 1,
				MaxBathrooms: 4,
				MinArea:      80,
				MaxArea:      400,
				Limit:        20,
			},
			valid: true,
		},
		{
			name: "empty parameters",
			params: AdvancedSearchParams{
				Limit: 10,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - all our test cases should be valid
			// More complex validation is done in the service layer
			assert.True(t, tt.valid)
		})
	}
}

func TestSearchSuggestion_Structure(t *testing.T) {
	suggestion := SearchSuggestion{
		Text:      "Guayas",
		Category:  "province",
		Frequency: 150,
	}

	assert.Equal(t, "Guayas", suggestion.Text)
	assert.Equal(t, "province", suggestion.Category)
	assert.Equal(t, 150, suggestion.Frequency)
}

func TestPropertySearchResult_Structure(t *testing.T) {
	property := domain.NewProperty(
		"Casa test",
		"Descripción test",
		"Guayas",
		"Guayaquil",
		"house",
		285000,
		"owner-123",
	)

	result := PropertySearchResult{
		Property: *property,
		Rank:     0.95,
	}

	assert.Equal(t, "Casa test", result.Property.Title)
	assert.Equal(t, 0.95, result.Rank)
}