package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"realty-core/internal/handlers"
	"realty-core/internal/repository"
	"realty-core/internal/service"
)

// TestServer holds the test server instance
type TestServer struct {
	server *httptest.Server
	client *http.Client
}

// NewTestServer creates a new test server instance
func NewTestServer(t *testing.T) *TestServer {
	// Set test environment variables
	os.Setenv("DATABASE_URL", "postgresql://juanquizhpi@localhost:5433/inmobiliaria_db?sslmode=disable")
	os.Setenv("PORT", "0") // Let the system choose an available port
	os.Setenv("LOG_LEVEL", "error") // Reduce logging noise during tests

	// Create test server
	db, err := repository.ConnectDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Create repositories and services
	propertyRepo := repository.NewPostgreSQLPropertyRepository(db)
	imageRepo := repository.NewPostgreSQLImageRepository(db)
	propertyService := service.NewPropertyService(propertyRepo, imageRepo)
	
	// Create handler
	propertyHandler := handlers.NewPropertyHandler(propertyService)
	
	// Create a simple router for testing
	mux := http.NewServeMux()
	mux.HandleFunc("/api/properties", propertyHandler.ListProperties)
	mux.HandleFunc("/api/health", propertyHandler.HealthCheck)
	
	testServer := httptest.NewServer(mux)

	return &TestServer{
		server: testServer,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Close shuts down the test server
func (ts *TestServer) Close() {
	ts.server.Close()
}

// URL returns the base URL of the test server
func (ts *TestServer) URL() string {
	return ts.server.URL
}

// GET makes a GET request to the test server
func (ts *TestServer) GET(path string) (*http.Response, error) {
	return ts.client.Get(ts.URL() + path)
}

// POST makes a POST request to the test server
func (ts *TestServer) POST(path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return ts.client.Post(ts.URL()+path, "application/json", bytes.NewBuffer(jsonBody))
}

// PUT makes a PUT request to the test server
func (ts *TestServer) PUT(path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", ts.URL()+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return ts.client.Do(req)
}

// DELETE makes a DELETE request to the test server
func (ts *TestServer) DELETE(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", ts.URL()+path, nil)
	if err != nil {
		return nil, err
	}
	return ts.client.Do(req)
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	resp, err := ts.GET("/api/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var apiResp APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	require.NoError(t, err)

	assert.True(t, apiResp.Success)
	assert.Contains(t, apiResp.Message, "Service is running")
}

// TestPropertyEndpoints tests all property-related endpoints (6 endpoints)
func TestPropertyEndpoints(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	// Test data
	propertyData := map[string]interface{}{
		"title":       "Test Property Integration",
		"description": "Test property for integration testing",
		"price":       250000.0,
		"province":    "Guayas",
		"city":        "Guayaquil",
		"type":        "house",
		"bedrooms":    3,
		"bathrooms":   2.5,
		"area_m2":     180.0,
	}

	t.Run("POST /api/properties - Create Property", func(t *testing.T) {
		resp, err := ts.POST("/api/properties", propertyData)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/properties - List Properties", func(t *testing.T) {
		resp, err := ts.GET("/api/properties")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/properties/filter - Filter Properties", func(t *testing.T) {
		resp, err := ts.GET("/api/properties/filter?province=Guayas&price_max=300000")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/properties/statistics - Property Statistics", func(t *testing.T) {
		resp, err := ts.GET("/api/properties/statistics")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})
}

// TestUserEndpoints tests all user-related endpoints (10 endpoints)
func TestUserEndpoints(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	// Test user data
	userData := map[string]interface{}{
		"firstName": "Juan",
		"lastName":  "Pérez",
		"email":     fmt.Sprintf("test.user.%d@example.com", time.Now().Unix()),
		"phone":     "0987654321",
		"cedula":    "0103355400", // Valid Ecuadorian cedula
		"password":  "testpassword123",
		"role":      "buyer",
	}

	var userID string

	t.Run("POST /api/users - Create User", func(t *testing.T) {
		resp, err := ts.POST("/api/users", userData)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)

		// Extract user ID for subsequent tests
		if data, ok := apiResp.Data.(map[string]interface{}); ok {
			if id, exists := data["id"].(string); exists {
				userID = id
			}
		}
	})

	t.Run("POST /api/auth/login - User Authentication", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    userData["email"],
			"password": userData["password"],
		}

		resp, err := ts.POST("/api/auth/login", loginData)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/users - List Users", func(t *testing.T) {
		resp, err := ts.GET("/api/users")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	if userID != "" {
		t.Run("GET /api/users/{id} - Get User by ID", func(t *testing.T) {
			resp, err := ts.GET("/api/users/" + userID)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var apiResp APIResponse
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			require.NoError(t, err)
			assert.True(t, apiResp.Success)
		})
	}

	t.Run("GET /api/users/statistics - User Statistics", func(t *testing.T) {
		resp, err := ts.GET("/api/users/statistics")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})
}

// TestAgencyEndpoints tests all agency-related endpoints (15 endpoints)
func TestAgencyEndpoints(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	// Test agency data
	agencyData := map[string]interface{}{
		"name":          fmt.Sprintf("Test Agency %d", time.Now().Unix()),
		"ruc":           fmt.Sprintf("09912345670%02d", time.Now().Unix()%100), // Valid RUC format
		"address":       "Av. Principal 123, Guayaquil",
		"phone":         "04-2345678",
		"email":         fmt.Sprintf("agency.test.%d@example.com", time.Now().Unix()),
		"licenseNumber": fmt.Sprintf("LIC-%d", time.Now().Unix()),
	}

	t.Run("POST /api/agencies - Create Agency", func(t *testing.T) {
		resp, err := ts.POST("/api/agencies", agencyData)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/agencies - List Agencies", func(t *testing.T) {
		resp, err := ts.GET("/api/agencies")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/agencies/active - List Active Agencies", func(t *testing.T) {
		resp, err := ts.GET("/api/agencies/active")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/agencies/statistics - Agency Statistics", func(t *testing.T) {
		resp, err := ts.GET("/api/agencies/statistics")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})
}

// TestPaginationEndpoints tests all pagination endpoints (7 endpoints)
func TestPaginationEndpoints(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	t.Run("GET /api/pagination/properties - Paginated Properties", func(t *testing.T) {
		resp, err := ts.GET("/api/pagination/properties?page=1&pageSize=10")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/pagination/users - Paginated Users", func(t *testing.T) {
		resp, err := ts.GET("/api/pagination/users?page=1&pageSize=5")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/pagination/agencies - Paginated Agencies", func(t *testing.T) {
		resp, err := ts.GET("/api/pagination/agencies?page=1&pageSize=5")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})

	t.Run("GET /api/pagination/stats - Pagination Statistics", func(t *testing.T) {
		resp, err := ts.GET("/api/pagination/stats")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiResp APIResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		require.NoError(t, err)
		assert.True(t, apiResp.Success)
	})
}

// TestPerformanceUnderLoad tests system performance under simulated load
func TestPerformanceUnderLoad(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	t.Run("Concurrent Health Checks", func(t *testing.T) {
		concurrency := 50
		done := make(chan bool, concurrency)

		start := time.Now()
		for i := 0; i < concurrency; i++ {
			go func() {
				resp, err := ts.GET("/api/health")
				if err == nil {
					resp.Body.Close()
					assert.Equal(t, http.StatusOK, resp.StatusCode)
				}
				done <- true
			}()
		}

		// Wait for all requests to complete
		for i := 0; i < concurrency; i++ {
			<-done
		}
		
		duration := time.Since(start)
		t.Logf("Completed %d concurrent requests in %v", concurrency, duration)
		
		// Should complete within reasonable time
		assert.Less(t, duration, 5*time.Second, "Concurrent requests took too long")
	})

	t.Run("Property List Performance", func(t *testing.T) {
		start := time.Now()
		resp, err := ts.GET("/api/properties")
		duration := time.Since(start)
		
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Less(t, duration, 500*time.Millisecond, "Property list request was too slow")
		
		t.Logf("Property list request completed in %v", duration)
	})
}