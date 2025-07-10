# Realty API Endpoints

Create REST API endpoints for: $ARGUMENTS

## Context - Existing API Structure
Current endpoints (26+ available):
- **Properties CRUD:** GET/POST/PUT/DELETE /api/properties
- **Search:** /api/properties/filter, /api/properties/search/ranked
- **Images:** 13 endpoints for image management
- **Cache:** /api/images/cache/stats

## API Design Patterns:
1. **RESTful conventions:**
   - Use HTTP methods correctly (GET, POST, PUT, DELETE)
   - Resource-based URLs: /api/properties/{id}
   - Proper HTTP status codes (200, 201, 400, 404, 500)
   - Consistent response format

2. **Request/Response structure:**
   ```go
   type SuccessResponse struct {
       Success bool        `json:"success"`
       Message string      `json:"message"`
       Data    interface{} `json:"data,omitempty"`
   }
   
   type ErrorResponse struct {
       Success bool   `json:"success"`
       Message string `json:"message"`
   }
   ```

3. **Handler patterns:**
   - Method validation first
   - Input parsing and validation
   - Business logic delegation to service layer
   - Proper error handling and logging
   - Response formatting

## Common endpoint patterns:

### **CRUD Operations:**
- `GET /api/properties` - List with pagination
- `POST /api/properties` - Create new property
- `GET /api/properties/{id}` - Get by ID
- `PUT /api/properties/{id}` - Update property
- `DELETE /api/properties/{id}` - Soft delete

### **Search and Filtering:**
- `GET /api/properties/filter?city=Quito&price_min=100000`
- `GET /api/properties/search/ranked?q=casa+piscina`
- `POST /api/properties/search/advanced` - Complex filters

### **Specialized endpoints:**
- `POST /api/properties/{id}/favorite` - User favorites
- `GET /api/properties/statistics` - Market stats
- `POST /api/properties/{id}/contact` - Contact owner

## Request validation patterns:
```go
// Use struct tags for validation
type CreatePropertyRequest struct {
    Title       string  `json:"title" validate:"required,min=10,max=200"`
    Description string  `json:"description" validate:"required,min=50"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Province    string  `json:"province" validate:"required,oneof=Azuay Bolívar Cañar..."`
    City        string  `json:"city" validate:"required,min=2"`
    Type        string  `json:"type" validate:"required,oneof=casa apartamento terreno comercial"`
}
```

## Pagination implementation:
```go
type PaginationParams struct {
    Page     int    `json:"page" form:"page"`
    PageSize int    `json:"page_size" form:"page_size"`
    OrderBy  string `json:"order_by" form:"order_by"`
    OrderDir string `json:"order_dir" form:"order_dir"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}
```

## Error handling:
- Validate HTTP method first
- Parse and validate input parameters
- Handle business logic errors from service layer
- Return appropriate HTTP status codes
- Log errors for debugging

## Security considerations:
- Input validation and sanitization
- Rate limiting for search endpoints
- CORS headers if needed
- Request size limits
- SQL injection prevention (using parameterized queries)

## Performance optimizations:
- Use pagination for large datasets
- Implement caching for frequent requests
- Optimize database queries
- Use connection pooling
- Add request timeouts

## Testing requirements:
- Unit tests for handler functions using httptest
- Integration tests for complete request/response flow
- Test error scenarios and edge cases
- Performance tests for search endpoints

## Common use cases:
- "create property favorites endpoint"
- "add property comparison API"
- "implement property statistics endpoint"
- "create property contact form API"
- "add property view tracking"
- "implement property recommendations"

## Integration with existing systems:
- **Property Service:** Delegate business logic
- **Image Service:** Handle property images
- **Cache Service:** Cache search results
- **FTS Service:** Full-text search integration

## Middleware integration:
- Request logging middleware
- CORS middleware
- Authentication middleware (future)
- Rate limiting middleware
- Error recovery middleware

## Documentation:
- API documentation with examples
- Request/response schemas
- Error codes and messages
- Rate limits and pagination info

## Example implementation:
```go
func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CreatePropertyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate input
    if err := h.validator.Struct(req); err != nil {
        h.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Business logic
    property, err := h.service.CreateProperty(req.Title, req.Description, ...)
    if err != nil {
        h.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.sendSuccessResponse(w, "Property created successfully", property)
}
```

## Output format:
- Complete handler function with proper error handling
- Request/response struct definitions
- Input validation logic
- Integration with service layer
- Unit tests for the endpoint
- API documentation for the endpoint