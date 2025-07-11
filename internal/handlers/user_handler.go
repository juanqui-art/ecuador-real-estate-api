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

// UserHandlerSimple handles HTTP requests for users (simplified)
type UserHandlerSimple struct {
	userService       *service.UserServiceSimple
	permissionService *service.PermissionService
	logger            *log.Logger
}

// NewUserHandlerSimple creates a new simplified user handler
func NewUserHandlerSimple(
	userService *service.UserServiceSimple,
	permissionService *service.PermissionService,
	logger *log.Logger,
) *UserHandlerSimple {
	return &UserHandlerSimple{
		userService:       userService,
		permissionService: permissionService,
		logger:            logger,
	}
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Cedula    string `json:"cedula"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Bio       string `json:"bio"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// CreateUser handles user creation
func (h *UserHandlerSimple) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate role
	role := domain.UserRole(req.Role)
	validRoles := []domain.UserRole{domain.RoleAdmin, domain.RoleAgent, domain.RoleBuyer, domain.RoleOwner}
	valid := false
	for _, r := range validRoles {
		if role == r {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.FirstName, req.LastName, req.Email, req.Phone, req.Cedula, req.Password, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendJSONResponse(w, user, http.StatusCreated)
}

// GetUser handles getting a user by ID
func (h *UserHandlerSimple) GetUser(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.sendJSONResponse(w, user, http.StatusOK)
}

// UpdateUser handles updating a user
func (h *UserHandlerSimple) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Update fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = &req.Phone
	user.Bio = &req.Bio

	if err := h.userService.UpdateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, user, http.StatusOK)
}

// DeleteUser handles user deletion
func (h *UserHandlerSimple) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Login handles user authentication
func (h *UserHandlerSimple) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Remove sensitive data before sending response
	user.PasswordHash = ""
	user.EmailVerificationToken = nil
	user.PasswordResetToken = nil

	h.sendJSONResponse(w, map[string]interface{}{
		"user": user,
		"message": "Login successful",
	}, http.StatusOK)
}

// ChangePassword handles password changes
func (h *UserHandlerSimple) ChangePassword(w http.ResponseWriter, r *http.Request) {
	id := h.extractIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.userService.ChangePassword(id, req.OldPassword, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "Password changed successfully"}, http.StatusOK)
}

// SearchUsers handles user search
func (h *UserHandlerSimple) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	name := query.Get("name")
	roleStr := query.Get("role")
	activeStr := query.Get("active")
	
	var role domain.UserRole
	if roleStr != "" {
		role = domain.UserRole(roleStr)
	}
	
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

	users, total, err := h.userService.SearchUsers("", name, role, active, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"users": users,
		"total": total,
		"limit": limit,
		"offset": offset,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetUsersByRole handles getting users by role
func (h *UserHandlerSimple) GetUsersByRole(w http.ResponseWriter, r *http.Request) {
	roleStr := h.extractIDFromPath(r.URL.Path)
	if roleStr == "" {
		http.Error(w, "Role required", http.StatusBadRequest)
		return
	}

	role := domain.UserRole(roleStr)
	users, err := h.userService.GetUsersByRole(role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, map[string]interface{}{
		"users": users,
		"role": roleStr,
		"count": len(users),
	}, http.StatusOK)
}

// GetUserStatistics handles getting user statistics
func (h *UserHandlerSimple) GetUserStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.userService.GetUserStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, stats, http.StatusOK)
}

// GetUserDashboard handles getting user dashboard data
func (h *UserHandlerSimple) GetUserDashboard(w http.ResponseWriter, r *http.Request) {
	// Simple dashboard with basic stats
	stats, err := h.userService.GetUserStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dashboard := map[string]interface{}{
		"statistics": stats,
		"message": "User dashboard",
		"timestamp": "2025-01-10",
	}

	h.sendJSONResponse(w, dashboard, http.StatusOK)
}

// Helper functions

func (h *UserHandlerSimple) extractIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3] // /api/users/{id}
	}
	return ""
}

func (h *UserHandlerSimple) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}