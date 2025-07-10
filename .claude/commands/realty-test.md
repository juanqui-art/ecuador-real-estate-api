# Realty Testing Framework

Create comprehensive tests for: $ARGUMENTS

## Context - Current Test Coverage
We have 157 tests with 90%+ coverage across:
- **Domain tests:** Property validation, business rules
- **Service tests:** Business logic with mocks
- **Repository tests:** Database operations with SQL mocks
- **Handler tests:** HTTP integration tests
- **Cache tests:** LRU cache performance and functionality
- **Storage tests:** File operations and image processing

## Testing Strategy:
1. **Unit Tests:** Test individual functions in isolation
2. **Integration Tests:** Test component interactions
3. **Performance Tests:** Benchmarks and load testing
4. **End-to-End Tests:** Complete workflows

## Test Patterns:

### **Table-driven tests:**
```go
func TestPropertyValidation(t *testing.T) {
    tests := []struct {
        name        string
        property    Property
        wantErr     bool
        expectedErr string
    }{
        {
            name: "valid property",
            property: Property{
                Title: "Casa en Quito",
                Price: 150000,
                Province: "Pichincha",
                City: "Quito",
                Type: "casa",
            },
            wantErr: false,
        },
        {
            name: "invalid province",
            property: Property{
                Title: "Casa en Ciudad Inexistente",
                Price: 150000,
                Province: "Inexistente",
                City: "Quito",
                Type: "casa",
            },
            wantErr: true,
            expectedErr: "invalid province",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.property.Validate()
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedErr)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### **Mock usage with testify:**
```go
func TestPropertyService_CreateProperty(t *testing.T) {
    mockRepo := new(MockPropertyRepository)
    mockImageRepo := new(MockImageRepository)
    service := NewPropertyService(mockRepo, mockImageRepo)

    // Setup mock expectations
    mockRepo.On("Create", mock.AnythingOfType("*domain.Property")).Return(nil)

    // Test execution
    property, err := service.CreateProperty("Test Property", "Description", "Pichincha", "Quito", "casa", 150000)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, property)
    assert.Equal(t, "Test Property", property.Title)
    mockRepo.AssertExpectations(t)
}
```

### **HTTP Handler testing:**
```go
func TestPropertyHandler_CreateProperty(t *testing.T) {
    mockService := new(MockPropertyService)
    handler := NewPropertyHandler(mockService)

    // Test data
    requestBody := `{
        "title": "Test Property",
        "description": "Test Description",
        "price": 150000,
        "province": "Pichincha",
        "city": "Quito",
        "type": "casa"
    }`

    // Setup mock
    mockService.On("CreateProperty", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
        Return(&domain.Property{ID: "test-id"}, nil)

    // Create request
    req := httptest.NewRequest(http.MethodPost, "/api/properties", strings.NewReader(requestBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    // Execute
    handler.CreateProperty(w, req)

    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    mockService.AssertExpectations(t)
}
```

### **Database integration tests:**
```go
func TestPropertyRepository_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()

    repo := NewPostgreSQLPropertyRepository(db)

    // Test data
    property := &domain.Property{
        Title:    "Integration Test Property",
        Price:    200000,
        Province: "Guayas",
        City:     "Guayaquil",
        Type:     "apartamento",
    }

    // Test Create
    err := repo.Create(property)
    assert.NoError(t, err)
    assert.NotEmpty(t, property.ID)

    // Test Get
    retrieved, err := repo.GetByID(property.ID)
    assert.NoError(t, err)
    assert.Equal(t, property.Title, retrieved.Title)

    // Test Update
    retrieved.Price = 250000
    err = repo.Update(retrieved.ID, retrieved)
    assert.NoError(t, err)

    // Test Delete
    err = repo.Delete(retrieved.ID)
    assert.NoError(t, err)
}
```

### **Cache performance tests:**
```go
func BenchmarkPropertyCache(b *testing.B) {
    cache := NewImageCache(DefaultImageCacheConfig())
    
    // Setup test data
    properties := make([]domain.Property, 1000)
    for i := range properties {
        properties[i] = domain.Property{
            ID:    fmt.Sprintf("prop-%d", i),
            Title: fmt.Sprintf("Property %d", i),
        }
    }

    b.ResetTimer()
    
    b.Run("Cache_Set", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            prop := properties[i%len(properties)]
            cache.Set(prop.ID, prop, "application/json")
        }
    })

    b.Run("Cache_Get", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            prop := properties[i%len(properties)]
            cache.Get(prop.ID)
        }
    })
}
```

## Test utilities and helpers:
```go
// Database setup for integration tests
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", "postgres://test:test@localhost/test_db?sslmode=disable")
    require.NoError(t, err)
    
    // Run migrations
    runMigrations(t, db)
    
    return db
}

// Property factory for tests
func CreateTestProperty(overrides ...func(*domain.Property)) *domain.Property {
    prop := &domain.Property{
        Title:    "Test Property",
        Price:    150000,
        Province: "Pichincha",
        City:     "Quito",
        Type:     "casa",
        Status:   "available",
    }
    
    for _, override := range overrides {
        override(prop)
    }
    
    return prop
}
```

## Common test scenarios:
- **Property validation:** Test all validation rules
- **Search functionality:** Test FTS queries and filters
- **Image processing:** Test upload, resize, and caching
- **Cache performance:** Test hit rates and memory usage
- **API endpoints:** Test all HTTP methods and error cases
- **Database operations:** Test CRUD and complex queries

## Test organization:
- **Unit tests:** `*_test.go` files alongside source code
- **Integration tests:** `integration_test.go` files
- **Performance tests:** `benchmark_test.go` files
- **Test utilities:** `testutils/` package for shared helpers

## Testing commands (use with Make):
```bash
make test-cache      # Test cache functionality
make test-images     # Test image processing
make test-properties # Test property CRUD
make test-handlers   # Test HTTP endpoints
make test-coverage   # Generate coverage report
make test-bench      # Run benchmarks
```

## Common use cases:
- "test property search with multiple filters"
- "test image upload with invalid formats"
- "test cache eviction under memory pressure"
- "test property validation for Ecuador provinces"
- "test API pagination with large datasets"
- "benchmark property search performance"

## Testing best practices:
- Use descriptive test names
- Test both happy path and error cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Test edge cases and boundary conditions
- Include performance benchmarks for critical paths

## Mock interfaces:
```go
type MockPropertyRepository struct {
    mock.Mock
}

func (m *MockPropertyRepository) Create(property *domain.Property) error {
    args := m.Called(property)
    return args.Error(0)
}

func (m *MockPropertyRepository) GetByID(id string) (*domain.Property, error) {
    args := m.Called(id)
    return args.Get(0).(*domain.Property), args.Error(1)
}
```

## Coverage goals:
- **Domain:** 90%+ coverage for business logic
- **Service:** 90%+ coverage for application logic
- **Repository:** 85%+ coverage for data access
- **Handlers:** 90%+ coverage for HTTP endpoints
- **Cache:** 95%+ coverage for caching logic

## Output format:
- Complete test files with proper structure
- Table-driven tests for multiple scenarios
- Mock implementations for dependencies
- Integration tests for complex workflows
- Performance benchmarks for critical operations
- Test utilities and helpers for common operations