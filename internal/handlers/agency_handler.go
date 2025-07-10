package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"realty-core/internal/domain"
	"realty-core/internal/service"
)

// AgencyHandler handles HTTP requests for agencies
type AgencyHandler struct {
	agencyService     *service.AgencyService
	permissionService *service.PermissionService
	logger            *log.Logger
}

// NewAgencyHandler creates a new agency handler
func NewAgencyHandler(
	agencyService *service.AgencyService,
	permissionService *service.PermissionService,
	logger *log.Logger,
) *AgencyHandler {
	return &AgencyHandler{
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
	Name           string            `json:"name,omitempty"`
	Address        string            `json:"address,omitempty"`
	Phone          string            `json:"phone,omitempty"`
	Email          string            `json:"email,omitempty"`
	Website        string            `json:"website,omitempty"`
	Description    string            `json:"description,omitempty"`
	LogoURL        string            `json:"logo_url,omitempty"`
	Commission     *float64          `json:"commission,omitempty"`
	BusinessHours  string            `json:"business_hours,omitempty"`
	SocialMedia    map[string]string `json:"social_media,omitempty"`
	Specialties    []string          `json:"specialties,omitempty"`
	ServiceAreas   []string          `json:"service_areas,omitempty"`
}

// SetLicenseRequest represents the request to set license information
type SetLicenseRequest struct {
	LicenseNumber string     `json:"license_number"`
	LicenseExpiry *time.Time `json:"license_expiry,omitempty"`
}

// AddSpecialtyRequest represents the request to add a specialty
type AddSpecialtyRequest struct {
	Specialty string `json:"specialty"`
}

// AddServiceAreaRequest represents the request to add a service area
type AddServiceAreaRequest struct {
	Province string `json:"province"`
}

// SetCommissionRequest represents the request to set commission
type SetCommissionRequest struct {
	Commission float64 `json:"commission"`
}

// CreateAgency handles POST /api/agencies
func (h *AgencyHandler) CreateAgency(w http.ResponseWriter, r *http.Request) {
	var req CreateAgencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check permissions
	hasPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionCreateAgency)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !hasPermission {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	agency, err := h.agencyService.CreateAgency(req.Name, req.RUC, req.Address, req.Phone, req.Email, req.LicenseNumber)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to create agency", err.Error())
		return
	}

	// Set optional fields
	if req.Website != "" {
		agency.Website = req.Website
	}
	if req.Description != "" {
		agency.Description = req.Description
	}
	if req.Commission > 0 {
		if err := agency.SetCommission(req.Commission); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "Invalid commission", err.Error())
			return
		}
	}

	if err := h.agencyService.UpdateAgency(agency); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to update agency", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusCreated, agency)
}

// GetAgency handles GET /api/agencies/{id}
func (h *AgencyHandler) GetAgency(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can view this agency
	canManage, err := h.permissionService.CanManageAgency(currentUserID, agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	// Allow viewing active agencies for non-managers
	if !canManage {
		hasViewPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionViewAgency)
		if err != nil || !hasViewPermission {
			WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
			return
		}
	}

	agency, err := h.agencyService.GetAgency(agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Agency not found", err.Error())
		return
	}

	// If user can't manage, only show active agencies
	if !canManage && !agency.Active {
		WriteErrorResponse(w, http.StatusNotFound, "Agency not found", "")
		return
	}

	WriteJSONResponse(w, http.StatusOK, agency)
}

// UpdateAgency handles PUT /api/agencies/{id}
func (h *AgencyHandler) UpdateAgency(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	var req UpdateAgencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can manage this agency
	canManage, err := h.permissionService.CanManageAgency(currentUserID, agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	// Get existing agency
	agency, err := h.agencyService.GetAgency(agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Agency not found", err.Error())
		return
	}

	// Update fields if provided
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
		agency.Website = req.Website
	}
	if req.Description != "" {
		agency.Description = req.Description
	}
	if req.LogoURL != "" {
		agency.LogoURL = req.LogoURL
	}
	if req.BusinessHours != "" {
		agency.BusinessHours = req.BusinessHours
	}
	if req.Commission != nil {
		if err := agency.SetCommission(*req.Commission); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "Invalid commission", err.Error())
			return
		}
	}
	if req.SocialMedia != nil {
		for platform, url := range req.SocialMedia {
			if err := agency.SetSocialMedia(platform, url); err != nil {
				WriteErrorResponse(w, http.StatusBadRequest, "Invalid social media", err.Error())
				return
			}
		}
	}
	if len(req.Specialties) > 0 {
		agency.Specialties = req.Specialties
	}
	if len(req.ServiceAreas) > 0 {
		agency.ServiceAreas = req.ServiceAreas
	}

	if err := h.agencyService.UpdateAgency(agency); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to update agency", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, agency)
}

// DeleteAgency handles DELETE /api/agencies/{id}
func (h *AgencyHandler) DeleteAgency(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check permissions
	hasPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionDeleteAgency)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !hasPermission {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	if err := h.agencyService.DeleteAgency(agencyID); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to delete agency", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Agency deleted successfully"})
}

// SearchAgencies handles GET /api/agencies
func (h *AgencyHandler) SearchAgencies(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check permissions
	hasPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionViewAllAgencies)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !hasPermission {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	// Parse query parameters
	params := &domain.AgencySearchParams{
		Query:      r.URL.Query().Get("query"),
		Pagination: ParsePaginationParams(r),
	}

	if activeStr := r.URL.Query().Get("active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			params.Active = &active
		}
	}

	if serviceAreas := r.URL.Query().Get("service_areas"); serviceAreas != "" {
		params.ServiceAreas = strings.Split(serviceAreas, ",")
	}

	if specialties := r.URL.Query().Get("specialties"); specialties != "" {
		params.Specialties = strings.Split(specialties, ",")
	}

	if minCommissionStr := r.URL.Query().Get("min_commission"); minCommissionStr != "" {
		if minCommission, err := strconv.ParseFloat(minCommissionStr, 64); err == nil {
			params.MinCommission = &minCommission
		}
	}

	if maxCommissionStr := r.URL.Query().Get("max_commission"); maxCommissionStr != "" {
		if maxCommission, err := strconv.ParseFloat(maxCommissionStr, 64); err == nil {
			params.MaxCommission = &maxCommission
		}
	}

	if licenseValidStr := r.URL.Query().Get("license_valid"); licenseValidStr != "" {
		if licenseValid, err := strconv.ParseBool(licenseValidStr); err == nil {
			params.LicenseValid = &licenseValid
		}
	}

	agencies, pagination, err := h.agencyService.SearchAgencies(params)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search agencies", err.Error())
		return
	}

	response := domain.PaginatedResponse{
		Data:       agencies,
		Pagination: pagination,
	}

	WriteJSONResponse(w, http.StatusOK, response)
}

// GetActiveAgencies handles GET /api/agencies/active
func (h *AgencyHandler) GetActiveAgencies(w http.ResponseWriter, r *http.Request) {
	agencies, err := h.agencyService.GetActiveAgencies()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get active agencies", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, agencies)
}

// GetAgenciesByServiceArea handles GET /api/agencies/service-area/{province}
func (h *AgencyHandler) GetAgenciesByServiceArea(w http.ResponseWriter, r *http.Request) {
	province := GetPathParam(r, "province")
	if province == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Province is required", "")
		return
	}

	agencies, err := h.agencyService.GetAgenciesByServiceArea(province)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get agencies by service area", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, agencies)
}

// GetAgenciesBySpecialty handles GET /api/agencies/specialty/{specialty}
func (h *AgencyHandler) GetAgenciesBySpecialty(w http.ResponseWriter, r *http.Request) {
	specialty := GetPathParam(r, "specialty")
	if specialty == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Specialty is required", "")
		return
	}

	agencies, err := h.agencyService.GetAgenciesBySpecialty(specialty)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get agencies by specialty", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, agencies)
}

// GetAgencyWithAgents handles GET /api/agencies/{id}/agents
func (h *AgencyHandler) GetAgencyWithAgents(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can view this agency's agents
	canManage, err := h.permissionService.CanManageAgency(currentUserID, agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	agencyWithAgents, err := h.agencyService.GetAgencyWithAgents(agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Agency not found", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, agencyWithAgents)
}

// SetAgencyLicense handles POST /api/agencies/{id}/license
func (h *AgencyHandler) SetAgencyLicense(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	var req SetLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can manage this agency
	canManage, err := h.permissionService.CanManageAgency(currentUserID, agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	if err := h.agencyService.SetAgencyLicense(agencyID, req.LicenseNumber, req.LicenseExpiry); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to set license", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "License updated successfully"})
}

// AddSpecialty handles POST /api/agencies/{id}/specialties
func (h *AgencyHandler) AddSpecialty(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	var req AddSpecialtyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can manage this agency
	canManage, err := h.permissionService.CanManageAgency(currentUserID, agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	if err := h.agencyService.AddSpecialtyToAgency(agencyID, req.Specialty); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to add specialty", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Specialty added successfully"})
}

// AddServiceArea handles POST /api/agencies/{id}/service-areas
func (h *AgencyHandler) AddServiceArea(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	var req AddServiceAreaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can manage this agency
	canManage, err := h.permissionService.CanManageAgency(currentUserID, agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	if err := h.agencyService.AddServiceAreaToAgency(agencyID, req.Province); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to add service area", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Service area added successfully"})
}

// SetCommission handles POST /api/agencies/{id}/commission
func (h *AgencyHandler) SetCommission(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	var req SetCommissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can manage this agency
	canManage, err := h.permissionService.CanManageAgency(currentUserID, agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	if err := h.agencyService.SetAgencyCommission(agencyID, req.Commission); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to set commission", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Commission updated successfully"})
}

// GetAgencyStatistics handles GET /api/agencies/statistics
func (h *AgencyHandler) GetAgencyStatistics(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check permissions
	hasPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionViewStats)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !hasPermission {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	stats, err := h.agencyService.GetAgencyStatistics()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get agency statistics", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, stats)
}

// GetAgencyPerformance handles GET /api/agencies/{id}/performance
func (h *AgencyHandler) GetAgencyPerformance(w http.ResponseWriter, r *http.Request) {
	agencyID := GetPathParam(r, "id")
	if agencyID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Agency ID is required", "")
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can view this agency's performance
	canManage, err := h.permissionService.CanManageAgency(currentUserID, agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	performance, err := h.agencyService.GetAgencyPerformance(agencyID)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "Agency not found", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, performance)
}

// RegisterAgencyRoutes registers all agency-related routes
func (h *AgencyHandler) RegisterRoutes(mux *http.ServeMux) {
	// Agency CRUD operations
	mux.HandleFunc("POST /api/agencies", h.CreateAgency)
	mux.HandleFunc("GET /api/agencies/{id}", h.GetAgency)
	mux.HandleFunc("PUT /api/agencies/{id}", h.UpdateAgency)
	mux.HandleFunc("DELETE /api/agencies/{id}", h.DeleteAgency)
	mux.HandleFunc("GET /api/agencies", h.SearchAgencies)
	
	// Public agency queries
	mux.HandleFunc("GET /api/agencies/active", h.GetActiveAgencies)
	mux.HandleFunc("GET /api/agencies/service-area/{province}", h.GetAgenciesByServiceArea)
	mux.HandleFunc("GET /api/agencies/specialty/{specialty}", h.GetAgenciesBySpecialty)
	
	// Agency management
	mux.HandleFunc("GET /api/agencies/{id}/agents", h.GetAgencyWithAgents)
	mux.HandleFunc("POST /api/agencies/{id}/license", h.SetAgencyLicense)
	mux.HandleFunc("POST /api/agencies/{id}/specialties", h.AddSpecialty)
	mux.HandleFunc("POST /api/agencies/{id}/service-areas", h.AddServiceArea)
	mux.HandleFunc("POST /api/agencies/{id}/commission", h.SetCommission)
	
	// Statistics and performance
	mux.HandleFunc("GET /api/agencies/statistics", h.GetAgencyStatistics)
	mux.HandleFunc("GET /api/agencies/{id}/performance", h.GetAgencyPerformance)
}