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

// AgencyHandlerSimple handles HTTP requests for agencies (simplified)
type AgencyHandlerSimple struct {
	agencyService     *service.AgencyService
	permissionService *service.PermissionService
	logger            *log.Logger
}

// NewAgencyHandlerSimple creates a new simplified agency handler
func NewAgencyHandlerSimple(
	agencyService *service.AgencyService,
	permissionService *service.PermissionService,
	logger *log.Logger,
) *AgencyHandlerSimple {
	return &AgencyHandlerSimple{
		agencyService:     agencyService,
		permissionService: permissionService,
		logger:            logger,
	}
}

// CreateAgencyRequest represents the request to create an agency
type CreateAgencyRequest struct {
	Name          string  `json:"name"`
	RUC           string  `json:"ruc"`
	Address       string  `json:"address"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	Website       string  `json:"website,omitempty"`
	Description   string  `json:"description,omitempty"`
	LicenseNumber string  `json:"license_number"`
	Commission    float64 `json:"commission,omitempty"`
}

// UpdateAgencyRequest represents the request to update an agency
type UpdateAgencyRequest struct {
	Name        string  `json:"name,omitempty"`
	Address     string  `json:"address,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	Email       string  `json:"email,omitempty"`
	Website     string  `json:"website,omitempty"`
	Description string  `json:"description,omitempty"`
	Commission  float64 `json:"commission,omitempty"`
}

// CreateAgency handles agency creation
func (h *AgencyHandlerSimple) CreateAgency(w http.ResponseWriter, r *http.Request) {
	var req CreateAgencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	agency, err := h.agencyService.CreateAgency(req.Name, req.RUC, req.Address, req.Phone, req.Email, req.LicenseNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set optional fields
	if req.Website != "" {
		agency.Website = &req.Website
	}
	if req.Description != "" {
		agency.Description = &req.Description
	}
	if req.Commission > 0 {
		agency.Commission = req.Commission
	}

	h.sendJSONResponse(w, agency, http.StatusCreated)
}

// GetAgency handles getting an agency by ID
func (h *AgencyHandlerSimple) GetAgency(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Agency ID required", http.StatusBadRequest)
		return
	}

	agency, err := h.agencyService.GetAgency(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.sendJSONResponse(w, agency, http.StatusOK)
}

// UpdateAgency handles updating an agency
func (h *AgencyHandlerSimple) UpdateAgency(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Agency ID required", http.StatusBadRequest)
		return
	}

	var req UpdateAgencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	agency, err := h.agencyService.GetAgency(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Update fields
	if req.Name != "" {
		agency.Name = req.Name
	}
	if req.Address != "" {
		agency.Address = req.Address
	}
	if req.Phone != "" {
		agency.Phone = req.Phone
	}
	if req.Email != "" {
		agency.Email = req.Email
	}
	if req.Website != "" {
		agency.Website = &req.Website
	}
	if req.Description != "" {
		agency.Description = &req.Description
	}
	if req.Commission > 0 {
		agency.Commission = req.Commission
	}

	if err := h.agencyService.UpdateAgency(agency); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, agency, http.StatusOK)
}

// DeleteAgency handles agency deletion
func (h *AgencyHandlerSimple) DeleteAgency(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Agency ID required", http.StatusBadRequest)
		return
	}

	if err := h.agencyService.DeleteAgency(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchAgencies handles agency search
func (h *AgencyHandlerSimple) SearchAgencies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	name := query.Get("name")
	province := query.Get("province")
	activeStr := query.Get("active")
	
	var active *bool
	if activeStr != "" {
		if activeStr == "true" {
			active = &[]bool{true}[0]
		} else if activeStr == "false" {
			active = &[]bool{false}[0]
		}
	}
	
	limit := 10
	offset := 0
	
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	
	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	params := &domain.AgencySearchParams{
		Query:    name,
		Province: province,
		Active:   active,
		Pagination: &domain.PaginationParams{
			Page:     offset/limit + 1,
			PageSize: limit,
		},
	}

	agencies, pagination, err := h.agencyService.SearchAgencies(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total := 0
	if pagination != nil {
		total = pagination.TotalRecords
	}

	response := map[string]interface{}{
		"agencies": agencies,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetActiveAgencies handles getting active agencies
func (h *AgencyHandlerSimple) GetActiveAgencies(w http.ResponseWriter, r *http.Request) {
	agencies, err := h.agencyService.GetActiveAgencies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, map[string]interface{}{
		"agencies": agencies,
		"count":    len(agencies),
	}, http.StatusOK)
}

// GetAgenciesByServiceArea handles getting agencies by service area
func (h *AgencyHandlerSimple) GetAgenciesByServiceArea(w http.ResponseWriter, r *http.Request) {
	area := h.extractIDFromPath(r.URL.Path)
	if area == "" {
		http.Error(w, "Service area required", http.StatusBadRequest)
		return
	}

	agencies, err := h.agencyService.GetAgenciesByServiceArea(area)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, map[string]interface{}{
		"agencies":     agencies,
		"service_area": area,
		"count":        len(agencies),
	}, http.StatusOK)
}

// GetAgenciesBySpecialty handles getting agencies by specialty
func (h *AgencyHandlerSimple) GetAgenciesBySpecialty(w http.ResponseWriter, r *http.Request) {
	specialty := h.extractIDFromPath(r.URL.Path)
	if specialty == "" {
		http.Error(w, "Specialty required", http.StatusBadRequest)
		return
	}

	agencies, err := h.agencyService.GetAgenciesBySpecialty(specialty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, map[string]interface{}{
		"agencies":  agencies,
		"specialty": specialty,
		"count":     len(agencies),
	}, http.StatusOK)
}

// GetAgencyAgents handles getting agents of an agency
func (h *AgencyHandlerSimple) GetAgencyAgents(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Agency ID required", http.StatusBadRequest)
		return
	}

	agents, err := h.agencyService.GetAgencyAgents(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, map[string]interface{}{
		"agents":    agents,
		"agency_id": id,
		"count":     len(agents),
	}, http.StatusOK)
}

// SetAgencyLicense handles setting agency license
func (h *AgencyHandlerSimple) SetAgencyLicense(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Agency ID required", http.StatusBadRequest)
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	licenseNumber := req["license_number"]
	if licenseNumber == "" {
		http.Error(w, "License number required", http.StatusBadRequest)
		return
	}

	if err := h.agencyService.SetLicenseNumber(id, licenseNumber); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "License set successfully"}, http.StatusOK)
}

// GetAgencyStatistics handles getting agency statistics
func (h *AgencyHandlerSimple) GetAgencyStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.agencyService.GetAgencyStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, stats, http.StatusOK)
}

// GetAgencyPerformance handles getting agency performance metrics
func (h *AgencyHandlerSimple) GetAgencyPerformance(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "Agency ID required", http.StatusBadRequest)
		return
	}

	performance, err := h.agencyService.GetAgencyPerformance(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, performance, http.StatusOK)
}

// Helper functions

func (h *AgencyHandlerSimple) extractIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3] // /api/agencies/{id}
	}
	return ""
}

func (h *AgencyHandlerSimple) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}