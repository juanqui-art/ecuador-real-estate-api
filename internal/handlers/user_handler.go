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

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService       *service.UserService
	permissionService *service.PermissionService
	logger            *log.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	userService *service.UserService,
	permissionService *service.PermissionService,
	logger *log.Logger,
) *UserHandler {
	return &UserHandler{
		userService:       userService,
		permissionService: permissionService,
		logger:            logger,
	}
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Email            string   `json:"email"`
	Phone            string   `json:"phone"`
	Cedula           string   `json:"cedula"`
	Password         string   `json:"password"`
	Role             string   `json:"role"`
	AgencyID         *string  `json:"agency_id,omitempty"`
	MinBudget        *float64 `json:"min_budget,omitempty"`
	MaxBudget        *float64 `json:"max_budget,omitempty"`
	InterestedProvinces []string `json:"interested_provinces,omitempty"`
	InterestedTypes  []string `json:"interested_types,omitempty"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName           string   `json:"first_name,omitempty"`
	LastName            string   `json:"last_name,omitempty"`
	Phone               string   `json:"phone,omitempty"`
	Bio                 string   `json:"bio,omitempty"`
	AvatarURL           string   `json:"avatar_url,omitempty"`
	MinBudget           *float64 `json:"min_budget,omitempty"`
	MaxBudget           *float64 `json:"max_budget,omitempty"`
	InterestedProvinces []string `json:"interested_provinces,omitempty"`
	InterestedTypes     []string `json:"interested_types,omitempty"`
}

// ChangePasswordRequest represents the request to change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// LoginRequest represents the request to authenticate
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get current user from context (would be set by auth middleware)
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check permissions
	hasPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionCreateUser)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !hasPermission {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	// Validate role
	role := domain.Role(req.Role)
	if !role.IsValid() {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid role", "")
		return
	}

	// Create user based on role
	var user *domain.User
	switch role {
	case domain.RoleAgent:
		if req.AgencyID == nil {
			WriteErrorResponse(w, http.StatusBadRequest, "Agency ID required for agents", "")
			return
		}
		user, err = h.userService.CreateAgent(req.FirstName, req.LastName, req.Email, req.Phone, req.Cedula, req.Password, *req.AgencyID)
	case domain.RoleBuyer:
		user, err = h.userService.CreateBuyer(req.FirstName, req.LastName, req.Email, req.Phone, req.Cedula, req.Password, req.MinBudget, req.MaxBudget, req.InterestedProvinces, req.InterestedTypes)
	default:
		user, err = h.userService.CreateUser(req.FirstName, req.LastName, req.Email, req.Phone, req.Cedula, req.Password, role)
	}

	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to create user", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusCreated, user)
}

// GetUser handles GET /api/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := GetPathParam(r, "id")
	if userID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "User ID is required", "")
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user can view this user
	canManage, err := h.permissionService.CanManageUser(currentUserID, userID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, user)
}

// UpdateUser handles PUT /api/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := GetPathParam(r, "id")
	if userID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "User ID is required", "")
		return
	}

	var req UpdateUserRequest
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

	// Check if user can manage this user
	canManage, err := h.permissionService.CanManageUser(currentUserID, userID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !canManage {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	// Get existing user
	user, err := h.userService.GetUser(userID)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	// Update preferences for buyers
	if user.Role == domain.RoleBuyer {
		if req.MinBudget != nil && req.MaxBudget != nil {
			if err := user.SetBudget(req.MinBudget, req.MaxBudget); err != nil {
				WriteErrorResponse(w, http.StatusBadRequest, "Invalid budget", err.Error())
				return
			}
		}

		if len(req.InterestedProvinces) > 0 {
			user.InterestedProvinces = req.InterestedProvinces
		}

		if len(req.InterestedTypes) > 0 {
			user.InterestedTypes = req.InterestedTypes
		}
	}

	if err := h.userService.UpdateUser(user); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, user)
}

// DeleteUser handles DELETE /api/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := GetPathParam(r, "id")
	if userID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "User ID is required", "")
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check permissions
	hasPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionDeleteUser)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !hasPermission {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to delete user", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// SearchUsers handles GET /api/users
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check permissions
	hasPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionViewAllUsers)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !hasPermission {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	// Parse query parameters
	params := &domain.UserSearchParams{
		Query:      r.URL.Query().Get("query"),
		Pagination: ParsePaginationParams(r),
	}

	if roleStr := r.URL.Query().Get("role"); roleStr != "" {
		role := domain.Role(roleStr)
		if role.IsValid() {
			params.Role = &role
		}
	}

	if activeStr := r.URL.Query().Get("active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			params.Active = &active
		}
	}

	if agencyID := r.URL.Query().Get("agency_id"); agencyID != "" {
		params.AgencyID = &agencyID
	}

	if provinces := r.URL.Query().Get("provinces"); provinces != "" {
		params.Provinces = strings.Split(provinces, ",")
	}

	if minBudgetStr := r.URL.Query().Get("min_budget"); minBudgetStr != "" {
		if minBudget, err := strconv.ParseFloat(minBudgetStr, 64); err == nil {
			params.MinBudget = &minBudget
		}
	}

	if maxBudgetStr := r.URL.Query().Get("max_budget"); maxBudgetStr != "" {
		if maxBudget, err := strconv.ParseFloat(maxBudgetStr, 64); err == nil {
			params.MaxBudget = &maxBudget
		}
	}

	users, pagination, err := h.userService.SearchUsers(params)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search users", err.Error())
		return
	}

	response := domain.PaginatedResponse{
		Data:       users,
		Pagination: pagination,
	}

	WriteJSONResponse(w, http.StatusOK, response)
}

// GetUsersByRole handles GET /api/users/role/{role}
func (h *UserHandler) GetUsersByRole(w http.ResponseWriter, r *http.Request) {
	roleStr := GetPathParam(r, "role")
	if roleStr == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Role is required", "")
		return
	}

	role := domain.Role(roleStr)
	if !role.IsValid() {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid role", "")
		return
	}

	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check permissions
	hasPermission, err := h.permissionService.HasPermission(currentUserID, service.PermissionViewAllUsers)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Permission check failed", err.Error())
		return
	}

	if !hasPermission {
		WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "")
		return
	}

	users, err := h.userService.GetUsersByRole(role)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get users by role", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, users)
}

// Login handles POST /api/auth/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication failed", err.Error())
		return
	}

	// Here you would typically generate a JWT token or session
	// For now, we'll just return the user
	response := map[string]interface{}{
		"user":         user,
		"permissions":  h.permissionService.GetRolePermissions(user.Role),
		"message":      "Login successful",
	}

	WriteJSONResponse(w, http.StatusOK, response)
}

// ChangePassword handles POST /api/users/{id}/password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := GetPathParam(r, "id")
	if userID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "User ID is required", "")
		return
	}

	var req ChangePasswordRequest
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

	// Users can only change their own password
	if currentUserID != userID {
		WriteErrorResponse(w, http.StatusForbidden, "Can only change your own password", "")
		return
	}

	if err := h.userService.ChangePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to change password", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

// GetUserStatistics handles GET /api/users/statistics
func (h *UserHandler) GetUserStatistics(w http.ResponseWriter, r *http.Request) {
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

	stats, err := h.userService.GetUserStatistics()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get user statistics", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, stats)
}

// GetDashboard handles GET /api/users/dashboard
func (h *UserHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	currentUserID := GetUserIDFromContext(r.Context())
	if currentUserID == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	dashboardData, err := h.permissionService.GetUserDashboardData(currentUserID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get dashboard data", err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusOK, dashboardData)
}

// RegisterUserRoutes registers all user-related routes
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	// User CRUD operations
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", h.GetUser)
	mux.HandleFunc("PUT /api/users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.DeleteUser)
	mux.HandleFunc("GET /api/users", h.SearchUsers)
	mux.HandleFunc("GET /api/users/role/{role}", h.GetUsersByRole)
	
	// Authentication
	mux.HandleFunc("POST /api/auth/login", h.Login)
	
	// User management
	mux.HandleFunc("POST /api/users/{id}/password", h.ChangePassword)
	
	// Statistics and dashboard
	mux.HandleFunc("GET /api/users/statistics", h.GetUserStatistics)
	mux.HandleFunc("GET /api/users/dashboard", h.GetDashboard)
}