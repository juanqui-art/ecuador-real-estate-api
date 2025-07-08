package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"realty-core/internal/servicio"

	"github.com/gorilla/mux"
)

// RealEstateCompanyHandler handles HTTP requests for real estate companies
type RealEstateCompanyHandler struct {
	service *servicio.RealEstateCompanyService
}

// NewRealEstateCompanyHandler creates a new real estate company handler instance
func NewRealEstateCompanyHandler(service *servicio.RealEstateCompanyService) *RealEstateCompanyHandler {
	return &RealEstateCompanyHandler{service: service}
}

// CreateRealEstateCompanyRequest represents the request structure for creating a real estate company
type CreateRealEstateCompanyRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	RUC         string `json:"ruc" validate:"required,len=13"`
	Address     string `json:"address" validate:"required,min=5,max=500"`
	Phone       string `json:"phone" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Website     string `json:"website,omitempty"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
}

// UpdateRealEstateCompanyRequest represents the request structure for updating a real estate company
type UpdateRealEstateCompanyRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	RUC         string `json:"ruc" validate:"required,len=13"`
	Address     string `json:"address" validate:"required,min=5,max=500"`
	Phone       string `json:"phone" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Website     string `json:"website,omitempty"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
}

// UpdateContactInfoRequest for updating contact information
type UpdateContactInfoRequest struct {
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// RegisterRealEstateCompanyRoutes registers all real estate company routes
func (h *RealEstateCompanyHandler) RegisterRealEstateCompanyRoutes(router *mux.Router) {
	// Main CRUD routes
	router.HandleFunc("/companies", h.CreateRealEstateCompany).Methods("POST")
	router.HandleFunc("/companies", h.GetRealEstateCompanies).Methods("GET")
	router.HandleFunc("/companies/active", h.GetActiveRealEstateCompanies).Methods("GET")
	router.HandleFunc("/companies/search", h.SearchRealEstateCompanies).Methods("GET")
	router.HandleFunc("/companies/statistics", h.GetStatistics).Methods("GET")
	router.HandleFunc("/companies/with-properties", h.GetCompaniesWithPropertyCount).Methods("GET")

	// Routes with ID
	router.HandleFunc("/companies/{id}", h.GetRealEstateCompany).Methods("GET")
	router.HandleFunc("/companies/{id}", h.UpdateRealEstateCompany).Methods("PUT")
	router.HandleFunc("/companies/{id}", h.DeleteRealEstateCompany).Methods("DELETE")
	router.HandleFunc("/companies/{id}/contact", h.UpdateContactInfo).Methods("PUT")
	router.HandleFunc("/companies/{id}/website", h.UpdateWebsite).Methods("PUT")
	router.HandleFunc("/companies/{id}/logo", h.UpdateLogoURL).Methods("PUT")
	router.HandleFunc("/companies/{id}/activate", h.ActivateRealEstateCompany).Methods("PUT")
	router.HandleFunc("/companies/{id}/deactivate", h.DeactivateRealEstateCompany).Methods("PUT")

	// Routes by RUC
	router.HandleFunc("/companies/ruc/{ruc}", h.GetRealEstateCompanyByRUC).Methods("GET")

	// Validation routes
	router.HandleFunc("/companies/validate/ruc", h.ValidateRUC).Methods("POST")
	router.HandleFunc("/companies/validate/email", h.ValidateEmail).Methods("POST")
	router.HandleFunc("/companies/validate/phone", h.ValidatePhone).Methods("POST")
	router.HandleFunc("/companies/check/ruc", h.CheckRUCAvailability).Methods("POST")
	router.HandleFunc("/companies/check/email", h.CheckEmailAvailability).Methods("POST")
}

// CreateRealEstateCompany creates a new real estate company
func (h *RealEstateCompanyHandler) CreateRealEstateCompany(w http.ResponseWriter, r *http.Request) {
	var req CreateRealEstateCompanyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	company, err := h.service.Create(req.Name, req.RUC, req.Address, req.Phone, req.Email)
	if err != nil {
		log.Printf("Error creating real estate company: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update optional fields if provided
	if req.Website != "" || req.Description != "" || req.LogoURL != "" {
		updatedCompany, err := h.service.Update(
			company.ID, company.Name, company.RUC, company.Address,
			company.Phone, company.Email, req.Website, req.Description, req.LogoURL,
		)
		if err != nil {
			log.Printf("Warning: Error updating optional fields for company %s: %v", company.ID, err)
			// Continue with the basic company creation
		} else {
			company = updatedCompany
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}

// GetRealEstateCompanies retrieves all real estate companies
func (h *RealEstateCompanyHandler) GetRealEstateCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.GetAll()
	if err != nil {
		log.Printf("Error getting real estate companies: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"companies": companies,
		"total":     len(companies),
	})
}

// GetActiveRealEstateCompanies retrieves all active real estate companies
func (h *RealEstateCompanyHandler) GetActiveRealEstateCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.GetActive()
	if err != nil {
		log.Printf("Error getting active real estate companies: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"companies": companies,
		"total":     len(companies),
	})
}

// GetRealEstateCompany retrieves a real estate company by ID
func (h *RealEstateCompanyHandler) GetRealEstateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	company, err := h.service.GetByID(id)
	if err != nil {
		log.Printf("Error getting real estate company %s: %v", id, err)
		http.Error(w, "Real estate company not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// GetRealEstateCompanyByRUC retrieves a real estate company by RUC
func (h *RealEstateCompanyHandler) GetRealEstateCompanyByRUC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ruc := vars["ruc"]

	company, err := h.service.GetByRUC(ruc)
	if err != nil {
		log.Printf("Error getting real estate company by RUC %s: %v", ruc, err)
		http.Error(w, "Real estate company not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// UpdateRealEstateCompany updates a real estate company
func (h *RealEstateCompanyHandler) UpdateRealEstateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateRealEstateCompanyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	company, err := h.service.Update(
		id, req.Name, req.RUC, req.Address, req.Phone, req.Email,
		req.Website, req.Description, req.LogoURL,
	)
	if err != nil {
		log.Printf("Error updating real estate company %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// UpdateContactInfo updates contact information for a real estate company
func (h *RealEstateCompanyHandler) UpdateContactInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateContactInfoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	company, err := h.service.UpdateContactInfo(id, req.Phone, req.Email)
	if err != nil {
		log.Printf("Error updating contact info for company %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// UpdateWebsite updates the website for a real estate company
func (h *RealEstateCompanyHandler) UpdateWebsite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Website string `json:"website"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	company, err := h.service.UpdateWebsite(id, req.Website)
	if err != nil {
		log.Printf("Error updating website for company %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// UpdateLogoURL updates the logo URL for a real estate company
func (h *RealEstateCompanyHandler) UpdateLogoURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		LogoURL string `json:"logo_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	company, err := h.service.UpdateLogoURL(id, req.LogoURL)
	if err != nil {
		log.Printf("Error updating logo URL for company %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// ActivateRealEstateCompany activates a real estate company
func (h *RealEstateCompanyHandler) ActivateRealEstateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	company, err := h.service.Activate(id)
	if err != nil {
		log.Printf("Error activating real estate company %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// DeactivateRealEstateCompany deactivates a real estate company
func (h *RealEstateCompanyHandler) DeactivateRealEstateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.Deactivate(id)
	if err != nil {
		log.Printf("Error deactivating real estate company %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Real estate company deactivated successfully",
		"id":      id,
	})
}

// SearchRealEstateCompanies searches real estate companies by name
func (h *RealEstateCompanyHandler) SearchRealEstateCompanies(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		http.Error(w, "Search term 'q' is required", http.StatusBadRequest)
		return
	}

	companies, err := h.service.SearchByName(searchTerm)
	if err != nil {
		log.Printf("Error searching real estate companies: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"companies":   companies,
		"total":       len(companies),
		"search_term": searchTerm,
	})
}

// DeleteRealEstateCompany deletes a real estate company
func (h *RealEstateCompanyHandler) DeleteRealEstateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.Delete(id)
	if err != nil {
		log.Printf("Error deleting real estate company %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Real estate company deleted successfully",
		"id":      id,
	})
}

// ValidateRUC validates a RUC
func (h *RealEstateCompanyHandler) ValidateRUC(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RUC string `json:"ruc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	err := h.service.ValidateRUC(req.RUC)
	isValid := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ruc":   req.RUC,
		"valid": isValid,
		"message": func() string {
			if isValid {
				return "RUC is valid"
			}
			return err.Error()
		}(),
	})
}

// ValidateEmail validates an email
func (h *RealEstateCompanyHandler) ValidateEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	err := h.service.ValidateEmail(req.Email)
	isValid := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"email": req.Email,
		"valid": isValid,
		"message": func() string {
			if isValid {
				return "Email is valid"
			}
			return err.Error()
		}(),
	})
}

// ValidatePhone validates a phone number
func (h *RealEstateCompanyHandler) ValidatePhone(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	err := h.service.ValidatePhone(req.Phone)
	isValid := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"phone": req.Phone,
		"valid": isValid,
		"message": func() string {
			if isValid {
				return "Phone is valid"
			}
			return err.Error()
		}(),
	})
}

// CheckRUCAvailability checks if a RUC is available
func (h *RealEstateCompanyHandler) CheckRUCAvailability(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RUC string `json:"ruc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	available, err := h.service.CheckRUCAvailability(req.RUC)
	if err != nil {
		log.Printf("Error checking RUC availability: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ruc":       req.RUC,
		"available": available,
		"message": func() string {
			if available {
				return "RUC is available"
			}
			return "RUC is already in use"
		}(),
	})
}

// CheckEmailAvailability checks if an email is available
func (h *RealEstateCompanyHandler) CheckEmailAvailability(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	available, err := h.service.CheckEmailAvailability(req.Email)
	if err != nil {
		log.Printf("Error checking email availability: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"email":     req.Email,
		"available": available,
		"message": func() string {
			if available {
				return "Email is available"
			}
			return "Email is already in use"
		}(),
	})
}

// GetStatistics returns statistics about real estate companies
func (h *RealEstateCompanyHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		log.Printf("Error getting real estate company statistics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetCompaniesWithPropertyCount returns companies with their property counts
func (h *RealEstateCompanyHandler) GetCompaniesWithPropertyCount(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.GetCompaniesWithPropertyCount()
	if err != nil {
		log.Printf("Error getting companies with property count: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"companies": companies,
		"total":     len(companies),
	})
}
