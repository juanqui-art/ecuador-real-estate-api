package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"realty-core/internal/domain"
	"realty-core/internal/service"
)

// PaginationHandlerSimple handles HTTP requests for pagination operations
type PaginationHandlerSimple struct {
	propertyService *service.PropertyService
	imageService    *service.ImageService
	userService     *service.UserServiceSimple
	agencyService   *service.AgencyService
	logger          *log.Logger
}

// NewPaginationHandlerSimple creates a new pagination handler
func NewPaginationHandlerSimple(
	propertyService *service.PropertyService,
	imageService *service.ImageService,
	userService *service.UserServiceSimple,
	agencyService *service.AgencyService,
	logger *log.Logger,
) *PaginationHandlerSimple {
	return &PaginationHandlerSimple{
		propertyService: propertyService,
		imageService:    imageService,
		userService:     userService,
		agencyService:   agencyService,
		logger:          logger,
	}
}

// GetPaginatedProperties handles paginated property retrieval
func (h *PaginationHandlerSimple) GetPaginatedProperties(w http.ResponseWriter, r *http.Request) {
	params := h.extractPaginationParams(r)
	
	// Use existing property service with pagination
	properties, err := h.propertyService.GetPaginatedProperties(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Count total properties for pagination metadata
	totalCount, err := h.propertyService.CountProperties()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pagination := domain.NewPagination(params.Page, params.PageSize, totalCount)
	
	response := domain.PaginatedResponse{
		Data:       properties,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetPaginatedImages handles paginated image retrieval
func (h *PaginationHandlerSimple) GetPaginatedImages(w http.ResponseWriter, r *http.Request) {
	params := h.extractPaginationParams(r)
	
	images, err := h.imageService.GetPaginatedImages(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := h.imageService.CountImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pagination := domain.NewPagination(params.Page, params.PageSize, totalCount)
	
	response := domain.PaginatedResponse{
		Data:       images,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetPaginatedUsers handles paginated user retrieval
func (h *PaginationHandlerSimple) GetPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	params := h.extractPaginationParams(r)
	
	// Use the correct method signature
	users, totalCount, err := h.userService.SearchUsers("", "", domain.UserRole(""), nil, params.PageSize, params.GetOffset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pagination := domain.NewPagination(params.Page, params.PageSize, totalCount)
	
	response := domain.PaginatedResponse{
		Data:       users,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetPaginatedAgencies handles paginated agency retrieval
func (h *PaginationHandlerSimple) GetPaginatedAgencies(w http.ResponseWriter, r *http.Request) {
	params := h.extractPaginationParams(r)
	
	searchParams := &domain.AgencySearchParams{
		Pagination: params,
	}

	agencies, pagination, err := h.agencyService.SearchAgencies(searchParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := domain.PaginatedResponse{
		Data:       agencies,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetPaginatedSearch handles paginated search across multiple entities
func (h *PaginationHandlerSimple) GetPaginatedSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	params := h.extractPaginationParams(r)
	
	// Search across properties, users, and agencies
	searchResults := make(map[string]interface{})
	
	// Search properties
	if properties, err := h.propertyService.SearchPropertiesSimple(query, params); err == nil {
		searchResults["properties"] = properties
	}

	// Search users
	if users, _, err := h.userService.SearchUsers("", query, domain.UserRole(""), nil, params.PageSize, params.GetOffset()); err == nil {
		searchResults["users"] = users
	}

	// Search agencies
	agencySearchParams := &domain.AgencySearchParams{
		Query:      query,
		Pagination: params,
	}
	if agencies, _, err := h.agencyService.SearchAgencies(agencySearchParams); err == nil {
		searchResults["agencies"] = agencies
	}

	pagination := domain.NewPagination(params.Page, params.PageSize, len(searchResults))
	
	response := domain.PaginatedResponse{
		Data:       searchResults,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetPaginationStats handles pagination statistics
func (h *PaginationHandlerSimple) GetPaginationStats(w http.ResponseWriter, r *http.Request) {
	stats := make(map[string]interface{})

	// Get counts for different entities
	if propertyCount, err := h.propertyService.CountProperties(); err == nil {
		stats["total_properties"] = propertyCount
	}

	if imageCount, err := h.imageService.CountImages(); err == nil {
		stats["total_images"] = imageCount
	}

	// Default pagination settings
	defaultParams := domain.NewPaginationParams()
	stats["default_page_size"] = defaultParams.PageSize
	stats["max_page_size"] = 100
	stats["default_sort_by"] = defaultParams.SortBy
	stats["default_sort_desc"] = defaultParams.SortDesc

	// Available sort fields
	stats["available_sort_fields"] = []string{
		"created_at", "updated_at", "title", "price", 
		"area_m2", "bedrooms", "bathrooms", "view_count",
	}

	h.sendJSONResponse(w, stats, http.StatusOK)
}

// AdvancedPaginationRequest represents advanced pagination request
type AdvancedPaginationRequest struct {
	Entity     string                 `json:"entity"`     // "properties", "users", "agencies", "images"
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	SortBy     string                 `json:"sort_by"`
	SortDesc   bool                   `json:"sort_desc"`
	Filters    map[string]interface{} `json:"filters"`
	SearchTerm string                 `json:"search_term"`
}

// HandleAdvancedPagination handles advanced pagination requests
func (h *PaginationHandlerSimple) HandleAdvancedPagination(w http.ResponseWriter, r *http.Request) {
	var req AdvancedPaginationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	params := &domain.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		SortDesc: req.SortDesc,
	}

	// Validate pagination parameters
	if err := params.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result interface{}
	var pagination *domain.Pagination
	var err error

	switch req.Entity {
	case "properties":
		result, err = h.propertyService.GetPaginatedProperties(params)
		if err == nil {
			if count, countErr := h.propertyService.CountProperties(); countErr == nil {
				pagination = domain.NewPagination(params.Page, params.PageSize, count)
			}
		}
	case "users":
		users, totalCount, searchErr := h.userService.SearchUsers("", req.SearchTerm, domain.UserRole(""), nil, params.PageSize, params.GetOffset())
		if searchErr == nil {
			result = users
			pagination = domain.NewPagination(params.Page, params.PageSize, totalCount)
		}
		err = searchErr
	case "agencies":
		searchParams := &domain.AgencySearchParams{
			Query:      req.SearchTerm,
			Pagination: params,
		}
		result, pagination, err = h.agencyService.SearchAgencies(searchParams)
	case "images":
		result, err = h.imageService.GetPaginatedImages(params)
		if err == nil {
			if count, countErr := h.imageService.CountImages(); countErr == nil {
				pagination = domain.NewPagination(params.Page, params.PageSize, count)
			}
		}
	default:
		http.Error(w, "Invalid entity type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := domain.PaginatedResponse{
		Data:       result,
		Pagination: pagination,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// Helper functions

func (h *PaginationHandlerSimple) extractPaginationParams(r *http.Request) *domain.PaginationParams {
	params := domain.NewPaginationParams()
	
	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	
	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 {
			params.PageSize = ps
		}
	}
	
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}
	
	if sortDesc := r.URL.Query().Get("sort_desc"); sortDesc != "" {
		params.SortDesc = strings.ToLower(sortDesc) == "true"
	}
	
	return params
}

func (h *PaginationHandlerSimple) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}