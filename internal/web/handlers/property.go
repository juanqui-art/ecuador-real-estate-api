package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"realty-core/internal/servicio"

	"github.com/gorilla/mux"
)

// PropertyHandler handles HTTP requests for properties
type PropertyHandler struct {
	service *servicio.PropertyService
}

// NewPropertyHandler creates a new property handler instance
func NewPropertyHandler(service *servicio.PropertyService) *PropertyHandler {
	return &PropertyHandler{service: service}
}

// CreatePropertyRequest represents the request structure for creating a property
type CreatePropertyRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Province    string  `json:"province" validate:"required"`
	City        string  `json:"city" validate:"required"`
	Sector      string  `json:"sector"`
	Address     string  `json:"address"`
	Type        string  `json:"type" validate:"required,oneof=house apartment land commercial"`
}

// UpdatePropertyRequest represents the request structure for updating a property
type UpdatePropertyRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Province    string  `json:"province" validate:"required"`
	City        string  `json:"city" validate:"required"`
	Sector      string  `json:"sector"`
	Address     string  `json:"address"`
	Type        string  `json:"type" validate:"required,oneof=house apartment land commercial"`
}

// UpdatePropertyDetailsRequest for updating detailed property information
type UpdatePropertyDetailsRequest struct {
	Bedrooms  int     `json:"bedrooms" validate:"min=0"`
	Bathrooms float32 `json:"bathrooms" validate:"min=0"`
	AreaM2    float64 `json:"area_m2" validate:"min=0"`
	YearBuilt *int    `json:"year_built,omitempty"`
	Floors    *int    `json:"floors,omitempty"`
}

// SetLocationRequest for setting property coordinates
type SetLocationRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-4,max=2"`
	Longitude float64 `json:"longitude" validate:"required,min=-92,max=-75"`
	Precision string  `json:"precision" validate:"required,oneof=exact approximate sector"`
}

// RegisterPropertyRoutes registers all property routes
func (h *PropertyHandler) RegisterPropertyRoutes(router *mux.Router) {
	// Main CRUD routes
	router.HandleFunc("/properties", h.CreateProperty).Methods("POST")
	router.HandleFunc("/properties", h.GetProperties).Methods("GET")
	router.HandleFunc("/properties/available", h.GetAvailableProperties).Methods("GET")
	router.HandleFunc("/properties/featured", h.GetFeaturedProperties).Methods("GET")
	router.HandleFunc("/properties/search", h.SearchProperties).Methods("GET")
	router.HandleFunc("/properties/statistics", h.GetStatistics).Methods("GET")

	// Routes with ID
	router.HandleFunc("/properties/{id}", h.GetProperty).Methods("GET")
	router.HandleFunc("/properties/{id}", h.UpdateProperty).Methods("PUT")
	router.HandleFunc("/properties/{id}", h.DeleteProperty).Methods("DELETE")
	router.HandleFunc("/properties/{id}/details", h.UpdatePropertyDetails).Methods("PUT")
	router.HandleFunc("/properties/{id}/location", h.SetPropertyLocation).Methods("PUT")
	router.HandleFunc("/properties/{id}/main-image", h.SetMainImage).Methods("PUT")
	router.HandleFunc("/properties/{id}/images", h.AddImage).Methods("POST")
	router.HandleFunc("/properties/{id}/tags", h.AddTag).Methods("POST")
	router.HandleFunc("/properties/{id}/featured", h.SetFeatured).Methods("PUT")
	router.HandleFunc("/properties/{id}/status", h.ChangeStatus).Methods("PUT")
	router.HandleFunc("/properties/{id}/company", h.AssignToCompany).Methods("PUT")
	router.HandleFunc("/properties/{id}/company", h.UnassignFromCompany).Methods("DELETE")

	// Routes by slug
	router.HandleFunc("/properties/slug/{slug}", h.GetPropertyBySlug).Methods("GET")

	// Company-specific routes
	router.HandleFunc("/companies/{company_id}/properties", h.GetPropertiesByCompany).Methods("GET")

	// Budget-based search
	router.HandleFunc("/properties/budget", h.GetPropertiesInBudget).Methods("GET")
}

// CreateProperty creates a new property
func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var req CreatePropertyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.Create(req.Title, req.Description, req.Province, req.City, req.Type, req.Price)
	if err != nil {
		log.Printf("Error creating property: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(property)
}

// GetProperties retrieves all properties with optional filters
func (h *PropertyHandler) GetProperties(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filters
	filters := make(map[string]interface{})

	if province := r.URL.Query().Get("province"); province != "" {
		filters["province"] = province
	}
	if city := r.URL.Query().Get("city"); city != "" {
		filters["city"] = city
	}
	if propertyType := r.URL.Query().Get("type"); propertyType != "" {
		filters["type"] = propertyType
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters["min_price"] = minPrice
		}
	}
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters["max_price"] = maxPrice
		}
	}
	if bedroomsStr := r.URL.Query().Get("bedrooms"); bedroomsStr != "" {
		if bedrooms, err := strconv.Atoi(bedroomsStr); err == nil {
			filters["bedrooms"] = bedrooms
		}
	}
	if featured := r.URL.Query().Get("featured"); featured == "true" {
		filters["featured"] = true
	}

	properties, err := h.service.GetAll(filters)
	if err != nil {
		log.Printf("Error getting properties: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": properties,
		"total":      len(properties),
		"filters":    filters,
	})
}

// GetAvailableProperties retrieves all available properties
func (h *PropertyHandler) GetAvailableProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.service.GetAvailable()
	if err != nil {
		log.Printf("Error getting available properties: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": properties,
		"total":      len(properties),
	})
}

// GetFeaturedProperties retrieves all featured properties
func (h *PropertyHandler) GetFeaturedProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.service.GetFeatured()
	if err != nil {
		log.Printf("Error getting featured properties: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": properties,
		"total":      len(properties),
	})
}

// GetProperty retrieves a property by ID
func (h *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	property, err := h.service.GetByID(id)
	if err != nil {
		log.Printf("Error getting property %s: %v", id, err)
		http.Error(w, "Property not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// GetPropertyBySlug retrieves a property by slug (and increments view count)
func (h *PropertyHandler) GetPropertyBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	property, err := h.service.GetBySlug(slug)
	if err != nil {
		log.Printf("Error getting property by slug %s: %v", slug, err)
		http.Error(w, "Property not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// UpdateProperty updates basic property information
func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdatePropertyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.Update(id, req.Title, req.Description, req.Province, req.City, req.Type, req.Price)
	if err != nil {
		log.Printf("Error updating property %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// UpdatePropertyDetails updates detailed property information
func (h *PropertyHandler) UpdatePropertyDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdatePropertyDetailsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.UpdateDetails(id, req.Bedrooms, req.Bathrooms, req.AreaM2, req.YearBuilt, req.Floors)
	if err != nil {
		log.Printf("Error updating property details %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// SetPropertyLocation sets GPS coordinates for a property
func (h *PropertyHandler) SetPropertyLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req SetLocationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.SetLocation(id, req.Latitude, req.Longitude, req.Precision)
	if err != nil {
		log.Printf("Error setting property location %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// SetMainImage sets the main image for a property
func (h *PropertyHandler) SetMainImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		ImageURL string `json:"image_url" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.SetMainImage(id, req.ImageURL)
	if err != nil {
		log.Printf("Error setting main image for property %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// AddImage adds an image to a property's gallery
func (h *PropertyHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		ImageURL string `json:"image_url" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.AddImage(id, req.ImageURL)
	if err != nil {
		log.Printf("Error adding image to property %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// AddTag adds a search tag to a property
func (h *PropertyHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Tag string `json:"tag" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.AddTag(id, req.Tag)
	if err != nil {
		log.Printf("Error adding tag to property %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// SetFeatured marks/unmarks a property as featured
func (h *PropertyHandler) SetFeatured(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Featured bool `json:"featured"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.SetFeatured(id, req.Featured)
	if err != nil {
		log.Printf("Error setting featured status for property %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// ChangeStatus changes the status of a property
func (h *PropertyHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Status string `json:"status" validate:"required,oneof=available sold rented reserved"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.ChangeStatus(id, req.Status)
	if err != nil {
		log.Printf("Error changing status for property %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// AssignToCompany assigns a property to a real estate company
func (h *PropertyHandler) AssignToCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		CompanyID string `json:"company_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	property, err := h.service.AssignToCompany(id, req.CompanyID)
	if err != nil {
		log.Printf("Error assigning property %s to company: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// UnassignFromCompany removes a property from a real estate company
func (h *PropertyHandler) UnassignFromCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	property, err := h.service.UnassignFromCompany(id)
	if err != nil {
		log.Printf("Error unassigning property %s from company: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// GetPropertiesByCompany retrieves all properties for a specific company
func (h *PropertyHandler) GetPropertiesByCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["company_id"]

	properties, err := h.service.GetByCompany(companyID)
	if err != nil {
		log.Printf("Error getting properties for company %s: %v", companyID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": properties,
		"total":      len(properties),
		"company_id": companyID,
	})
}

// SearchProperties performs full-text search on properties
func (h *PropertyHandler) SearchProperties(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		http.Error(w, "Search term 'q' is required", http.StatusBadRequest)
		return
	}

	properties, err := h.service.Search(searchTerm)
	if err != nil {
		log.Printf("Error searching properties: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties":  properties,
		"total":       len(properties),
		"search_term": searchTerm,
	})
}

// GetPropertiesInBudget retrieves properties within a budget range
func (h *PropertyHandler) GetPropertiesInBudget(w http.ResponseWriter, r *http.Request) {
	var minPrice, maxPrice *float64

	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if price, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = &price
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if price, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = &price
		}
	}

	if minPrice == nil && maxPrice == nil {
		http.Error(w, "At least one of min_price or max_price is required", http.StatusBadRequest)
		return
	}

	properties, err := h.service.GetPropertiesInBudgetRange(minPrice, maxPrice)
	if err != nil {
		log.Printf("Error getting properties in budget: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": properties,
		"total":      len(properties),
		"min_price":  minPrice,
		"max_price":  maxPrice,
	})
}

// DeleteProperty deletes a property
func (h *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.Delete(id)
	if err != nil {
		log.Printf("Error deleting property %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Property deleted successfully",
		"id":      id,
	})
}

// GetStatistics returns property statistics
func (h *PropertyHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		log.Printf("Error getting property statistics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
