package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"realty-core/internal/domain"
	"realty-core/internal/repository"
	"realty-core/internal/service"
)

// PropertyHandler handles HTTP requests for properties
type PropertyHandler struct {
	service service.PropertyServiceInterface
}

// NewPropertyHandler creates a new instance of the handler
func NewPropertyHandler(service service.PropertyServiceInterface) *PropertyHandler {
	return &PropertyHandler{service: service}
}

// CreatePropertyRequest represents the request structure for creating a property
// Updated to match complete domain Property struct - ALL 50+ fields supported (2025)
type CreatePropertyRequest struct {
	// Basic Information
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	Type          string  `json:"type"`
	Status        string  `json:"status"`
	
	// Location (expanded with all domain fields)
	Province          string  `json:"province"`
	City              string  `json:"city"`
	Sector            string  `json:"sector,omitempty"`
	Address           string  `json:"address,omitempty"`
	Latitude          float64 `json:"latitude,omitempty"`
	Longitude         float64 `json:"longitude,omitempty"`
	LocationPrecision string  `json:"location_precision,omitempty"`
	
	// Property Characteristics (expanded)
	Bedrooms      int     `json:"bedrooms"`
	Bathrooms     float32 `json:"bathrooms"`
	AreaM2        float64 `json:"area_m2"`
	ParkingSpaces int     `json:"parking_spaces"`
	YearBuilt     *int    `json:"year_built,omitempty"`
	Floors        *int    `json:"floors,omitempty"`
	
	// Additional Pricing
	RentPrice      *float64 `json:"rent_price,omitempty"`
	CommonExpenses *float64 `json:"common_expenses,omitempty"`
	PricePerM2     *float64 `json:"price_per_m2,omitempty"`
	
	// Multimedia
	MainImage *string  `json:"main_image,omitempty"`
	Images    []string `json:"images,omitempty"`
	VideoTour *string  `json:"video_tour,omitempty"`
	Tour360   *string  `json:"tour_360,omitempty"`
	
	// State and Classification
	PropertyStatus string   `json:"property_status,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Featured       bool     `json:"featured"`
	
	// Amenities (boolean fields) - complete set
	Garden            bool `json:"garden"`
	Pool              bool `json:"pool"`
	Elevator          bool `json:"elevator"`
	Balcony           bool `json:"balcony"`
	Terrace           bool `json:"terrace"`
	Garage            bool `json:"garage"`
	Furnished         bool `json:"furnished"`
	AirConditioning   bool `json:"air_conditioning"`
	Security          bool `json:"security"`
	
	// Ownership System (optional for forms, handled by backend)
	RealEstateCompanyID *string `json:"real_estate_company_id,omitempty"`
	OwnerID             *string `json:"owner_id,omitempty"`
	AgentID             *string `json:"agent_id,omitempty"`
	AgencyID            *string `json:"agency_id,omitempty"`
	
	// Contact Information (temporary until user system)
	ContactPhone  string `json:"contact_phone"`
	ContactEmail  string `json:"contact_email"`
	Notes         string `json:"notes,omitempty"`
}


// CreateProperty handles POST /api/properties
func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Convert CreatePropertyRequest to service CreatePropertyFullRequest
	// Updated to map ALL 50+ fields from expanded structs (2025)
	serviceReq := service.CreatePropertyFullRequest{
		// Basic Information
		Title:         req.Title,
		Description:   req.Description,
		Price:         req.Price,
		Type:          req.Type,
		Status:        req.Status,
		
		// Location (expanded)
		Province:          req.Province,
		City:              req.City,
		Sector:            req.Sector,
		Address:           req.Address,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		LocationPrecision: req.LocationPrecision,
		
		// Property Characteristics (expanded)
		Bedrooms:      req.Bedrooms,
		Bathrooms:     req.Bathrooms,
		AreaM2:        req.AreaM2,
		ParkingSpaces: req.ParkingSpaces,
		YearBuilt:     req.YearBuilt,
		Floors:        req.Floors,
		
		// Additional Pricing
		RentPrice:      req.RentPrice,
		CommonExpenses: req.CommonExpenses,
		PricePerM2:     req.PricePerM2,
		
		// Multimedia
		MainImage: req.MainImage,
		Images:    req.Images,
		VideoTour: req.VideoTour,
		Tour360:   req.Tour360,
		
		// State and Classification
		PropertyStatus: req.PropertyStatus,
		Tags:           req.Tags,
		Featured:       req.Featured,
		
		// Amenities (complete set)
		Garden:            req.Garden,
		Pool:              req.Pool,
		Elevator:          req.Elevator,
		Balcony:           req.Balcony,
		Terrace:           req.Terrace,
		Garage:            req.Garage,
		Furnished:         req.Furnished,
		AirConditioning:   req.AirConditioning,
		Security:          req.Security,
		
		// Ownership System
		RealEstateCompanyID: req.RealEstateCompanyID,
		OwnerID:             req.OwnerID,
		AgentID:             req.AgentID,
		AgencyID:            req.AgencyID,
		
		// Contact Information
		ContactPhone:  req.ContactPhone,
		ContactEmail:  req.ContactEmail,
		Notes:         req.Notes,
	}

	property, err := h.service.CreatePropertyComplete(serviceReq)

	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, property, "Property created successfully")
}

// GetProperty handles GET /api/properties/{id}
func (h *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := h.extractIDFromURL(r.URL.Path)
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Property ID required")
		return
	}

	property, err := h.service.GetProperty(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.respondSuccess(w, http.StatusOK, property, "Property retrieved successfully")
}

// GetPropertyBySlug handles GET /api/properties/slug/{slug}
func (h *PropertyHandler) GetPropertyBySlug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	slug := h.extractSlugFromURL(r.URL.Path)
	if slug == "" {
		h.respondError(w, http.StatusBadRequest, "Property slug required")
		return
	}

	property, err := h.service.GetPropertyBySlug(slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.respondSuccess(w, http.StatusOK, property, "Property retrieved by slug successfully")
}

// ListProperties handles GET /api/properties
func (h *PropertyHandler) ListProperties(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	properties, err := h.service.ListProperties()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, properties, "Properties retrieved successfully")
}

// UpdateProperty handles PUT /api/properties/{id}
func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := h.extractIDFromURL(r.URL.Path)
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Property ID required")
		return
	}

	var req CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	property, err := h.service.UpdateProperty(
		id,
		req.Title,
		req.Description,
		req.Province,
		req.City,
		req.Type,
		req.Price,
	)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.respondSuccess(w, http.StatusOK, property, "Property updated successfully")
}

// DeleteProperty handles DELETE /api/properties/{id}
func (h *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := h.extractIDFromURL(r.URL.Path)
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Property ID required")
		return
	}

	err := h.service.DeleteProperty(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.respondSuccess(w, http.StatusOK, nil, "Property deleted successfully")
}

// FilterProperties handles GET /api/properties/filter (basic filtering)
func (h *PropertyHandler) FilterProperties(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	province := query.Get("province")
	minPriceStr := query.Get("min_price")
	maxPriceStr := query.Get("max_price")
	searchQuery := query.Get("q")

	// Search by query if provided
	if searchQuery != "" {
		properties, err := h.service.SearchProperties(searchQuery)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondSuccess(w, http.StatusOK, properties, "Properties filtered by search query")
		return
	}

	// Filter by province if provided
	if province != "" {
		properties, err := h.service.FilterByProvince(province)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondSuccess(w, http.StatusOK, properties, "Properties filtered by province")
		return
	}

	// Filter by price range if provided
	if minPriceStr != "" && maxPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, "Invalid minimum price")
			return
		}

		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, "Invalid maximum price")
			return
		}

		properties, err := h.service.FilterByPriceRange(minPrice, maxPrice)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondSuccess(w, http.StatusOK, properties, "Properties filtered by price range")
		return
	}

	// If no filters, return all properties
	properties, err := h.service.ListProperties()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, properties, "All properties")
}

// SearchRanked handles GET /api/properties/search/ranked
func (h *PropertyHandler) SearchRanked(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	searchQuery := query.Get("q")
	limitStr := query.Get("limit")

	if searchQuery == "" {
		h.respondError(w, http.StatusBadRequest, "Search query required")
		return
	}

	limit := 50
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			h.respondError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
		limit = parsedLimit
	}

	results, err := h.service.SearchPropertiesRanked(searchQuery, limit)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, results, "Ranked search results retrieved successfully")
}

// SearchSuggestions handles GET /api/properties/search/suggestions
func (h *PropertyHandler) SearchSuggestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	searchQuery := query.Get("q")
	limitStr := query.Get("limit")

	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			h.respondError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
		limit = parsedLimit
	}

	suggestions, err := h.service.GetSearchSuggestions(searchQuery, limit)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, suggestions, "Search suggestions retrieved successfully")
}

// AdvancedSearch handles POST /api/properties/search/advanced
func (h *PropertyHandler) AdvancedSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Query        string  `json:"query"`
		Province     string  `json:"province"`
		City         string  `json:"city"`
		Type         string  `json:"type"`
		MinPrice     float64 `json:"min_price"`
		MaxPrice     float64 `json:"max_price"`
		MinBedrooms  int     `json:"min_bedrooms"`
		MaxBedrooms  int     `json:"max_bedrooms"`
		MinBathrooms float64 `json:"min_bathrooms"`
		MaxBathrooms float64 `json:"max_bathrooms"`
		MinArea      float64 `json:"min_area"`
		MaxArea      float64 `json:"max_area"`
		FeaturedOnly bool    `json:"featured_only"`
		Limit        int     `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Build search parameters
	params := repository.AdvancedSearchParams{
		Query:        req.Query,
		Province:     req.Province,
		City:         req.City,
		Type:         req.Type,
		MinPrice:     req.MinPrice,
		MaxPrice:     req.MaxPrice,
		MinBedrooms:  req.MinBedrooms,
		MaxBedrooms:  req.MaxBedrooms,
		MinBathrooms: req.MinBathrooms,
		MaxBathrooms: req.MaxBathrooms,
		MinArea:      req.MinArea,
		MaxArea:      req.MaxArea,
		FeaturedOnly: req.FeaturedOnly,
		Limit:        req.Limit,
	}

	results, err := h.service.AdvancedSearch(params)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, results, "Advanced search results retrieved successfully")
}

// GetStatistics handles GET /api/properties/statistics
func (h *PropertyHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	stats, err := h.service.GetStatistics()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, stats, "Statistics retrieved successfully")
}

// SetPropertyLocation handles POST /api/properties/{id}/location
func (h *PropertyHandler) SetPropertyLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := h.extractIDFromNestedURL(r.URL.Path)
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Property ID required")
		return
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Precision string  `json:"precision"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	err := h.service.SetPropertyLocation(id, req.Latitude, req.Longitude, req.Precision)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.respondSuccess(w, http.StatusOK, nil, "Property location updated successfully")
}

// SetPropertyFeatured handles POST /api/properties/{id}/featured
func (h *PropertyHandler) SetPropertyFeatured(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := h.extractIDFromNestedURL(r.URL.Path)
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Property ID required")
		return
	}

	var req struct {
		Featured bool `json:"featured"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	err := h.service.SetPropertyFeatured(id, req.Featured)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.respondSuccess(w, http.StatusOK, nil, "Property featured status updated successfully")
}

// SetPropertyParkingSpaces handles POST /api/properties/{id}/parking-spaces
func (h *PropertyHandler) SetPropertyParkingSpaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := h.extractIDFromNestedURL(r.URL.Path)
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Property ID required")
		return
	}

	var req struct {
		ParkingSpaces int `json:"parking_spaces"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	err := h.service.SetPropertyParkingSpaces(id, req.ParkingSpaces)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.respondSuccess(w, http.StatusOK, nil, "Property parking spaces updated successfully")
}

// HealthCheck handles GET /api/health
func (h *PropertyHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	health := map[string]string{
		"status":  "healthy",
		"service": "real-estate-api",
		"version": "1.0.0",
	}

	h.respondSuccess(w, http.StatusOK, health, "Service is running correctly")
}

// Helper methods

// extractIDFromURL extracts the ID from the URL (last segment)
func (h *PropertyHandler) extractIDFromURL(path string) string {
	// Remove trailing slash if exists
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	// Split by slash and get last segment
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

// extractIDFromNestedURL extracts the ID from nested URLs like /api/properties/{id}/location
func (h *PropertyHandler) extractIDFromNestedURL(path string) string {
	// Remove trailing slash if exists
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	// Split by slash
	parts := strings.Split(path, "/")

	// Look for pattern /api/properties/{id}/action
	// parts should be: ["", "api", "properties", "{id}", "action"]
	if len(parts) >= 4 && parts[1] == "api" && parts[2] == "properties" {
		return parts[3] // Return the ID part
	}

	return ""
}

// extractSlugFromURL extracts the slug from URL for routes /api/properties/slug/{slug}
func (h *PropertyHandler) extractSlugFromURL(path string) string {
	// Remove trailing slash if exists
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	// Split by slash
	parts := strings.Split(path, "/")

	// Look for pattern /api/properties/slug/{slug}
	// parts should be: ["", "api", "properties", "slug", "{slug}"]
	if len(parts) >= 5 && parts[1] == "api" && parts[2] == "properties" && parts[3] == "slug" {
		return parts[4]
	}

	return ""
}

// respondError sends an error response in JSON format
func (h *PropertyHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResp := ErrorResponse{
		Success: false,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// respondSuccess sends a successful response in JSON format
func (h *PropertyHandler) respondSuccess(w http.ResponseWriter, status int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	successResp := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(successResp); err != nil {
		log.Printf("Error encoding success response: %v", err)
	}
}

// Pagination handlers

// ListPropertiesPaginated handles GET /api/properties/paginated
func (h *PropertyHandler) ListPropertiesPaginated(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pagination, err := h.parsePaginationParams(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.ListPropertiesPaginated(pagination)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, result, "Paginated properties retrieved successfully")
}

// FilterPropertiesPaginated handles GET /api/properties/filter/paginated
func (h *PropertyHandler) FilterPropertiesPaginated(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	province := query.Get("province")
	minPriceStr := query.Get("min_price")
	maxPriceStr := query.Get("max_price")
	searchQuery := query.Get("q")

	pagination, err := h.parsePaginationParams(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var result *domain.PaginatedResponse

	// Search by query if provided
	if searchQuery != "" {
		result, err = h.service.SearchPropertiesPaginated(searchQuery, pagination)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondSuccess(w, http.StatusOK, result, "Paginated properties filtered by search query")
		return
	}

	// Filter by province if provided
	if province != "" {
		result, err = h.service.FilterByProvincePaginated(province, pagination)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondSuccess(w, http.StatusOK, result, "Paginated properties filtered by province")
		return
	}

	// Filter by price range if provided
	if minPriceStr != "" && maxPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, "Invalid minimum price")
			return
		}

		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, "Invalid maximum price")
			return
		}

		result, err = h.service.FilterByPriceRangePaginated(minPrice, maxPrice, pagination)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondSuccess(w, http.StatusOK, result, "Paginated properties filtered by price range")
		return
	}

	// If no filters, return all properties paginated
	result, err = h.service.ListPropertiesPaginated(pagination)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, result, "All paginated properties")
}

// SearchRankedPaginated handles GET /api/properties/search/ranked/paginated
func (h *PropertyHandler) SearchRankedPaginated(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	searchQuery := query.Get("q")

	if searchQuery == "" {
		h.respondError(w, http.StatusBadRequest, "Search query required")
		return
	}

	pagination, err := h.parsePaginationParams(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.SearchPropertiesRankedPaginated(searchQuery, pagination)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, result, "Paginated ranked search results retrieved successfully")
}

// AdvancedSearchPaginated handles POST /api/properties/search/advanced/paginated
func (h *PropertyHandler) AdvancedSearchPaginated(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Query        string                   `json:"query"`
		Province     string                   `json:"province"`
		City         string                   `json:"city"`
		Type         string                   `json:"type"`
		MinPrice     float64                  `json:"min_price"`
		MaxPrice     float64                  `json:"max_price"`
		MinBedrooms  int                      `json:"min_bedrooms"`
		MaxBedrooms  int                      `json:"max_bedrooms"`
		MinBathrooms float64                  `json:"min_bathrooms"`
		MaxBathrooms float64                  `json:"max_bathrooms"`
		MinArea      float64                  `json:"min_area"`
		MaxArea      float64                  `json:"max_area"`
		FeaturedOnly bool                     `json:"featured_only"`
		Pagination   *domain.PaginationParams `json:"pagination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Build search parameters
	params := repository.AdvancedSearchParams{
		Query:        req.Query,
		Province:     req.Province,
		City:         req.City,
		Type:         req.Type,
		MinPrice:     req.MinPrice,
		MaxPrice:     req.MaxPrice,
		MinBedrooms:  req.MinBedrooms,
		MaxBedrooms:  req.MaxBedrooms,
		MinBathrooms: req.MinBathrooms,
		MaxBathrooms: req.MaxBathrooms,
		MinArea:      req.MinArea,
		MaxArea:      req.MaxArea,
		FeaturedOnly: req.FeaturedOnly,
	}

	pagination := req.Pagination
	if pagination == nil {
		pagination = domain.NewPaginationParams()
	}

	result, err := h.service.AdvancedSearchPaginated(params, pagination)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, result, "Paginated advanced search results retrieved successfully")
}

// parsePaginationParams parses pagination parameters from URL query string
func (h *PropertyHandler) parsePaginationParams(r *http.Request) (*domain.PaginationParams, error) {
	query := r.URL.Query()
	
	pagination := domain.NewPaginationParams()
	
	// Parse page
	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("invalid page parameter: %s", pageStr)
		}
		pagination.Page = page
	}
	
	// Parse page_size
	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid page_size parameter: %s", pageSizeStr)
		}
		pagination.PageSize = pageSize
	}
	
	// Parse sort_by
	if sortBy := query.Get("sort_by"); sortBy != "" {
		pagination.SortBy = sortBy
	}
	
	// Parse sort_desc
	if sortDescStr := query.Get("sort_desc"); sortDescStr != "" {
		sortDesc, err := strconv.ParseBool(sortDescStr)
		if err != nil {
			return nil, fmt.Errorf("invalid sort_desc parameter: %s", sortDescStr)
		}
		pagination.SortDesc = sortDesc
	}
	
	return pagination, nil
}