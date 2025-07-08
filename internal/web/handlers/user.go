package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"realty-core/internal/servicio"

	"github.com/gorilla/mux"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service *servicio.UserService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(service *servicio.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUserRequest represents the request structure for creating a user
type CreateUserRequest struct {
	FirstName  string `json:"first_name" validate:"required,min=2,max=100"`
	LastName   string `json:"last_name" validate:"required,min=2,max=100"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required"`
	NationalID string `json:"national_id" validate:"required,len=10"`
	UserType   string `json:"user_type" validate:"required,oneof=buyer seller agent admin"`
}

// CreateBuyerRequest for creating buyers with preferences
type CreateBuyerRequest struct {
	CreateUserRequest
	MinBudget              *float64 `json:"min_budget,omitempty"`
	MaxBudget              *float64 `json:"max_budget,omitempty"`
	PreferredProvinces     []string `json:"preferred_provinces,omitempty"`
	PreferredPropertyTypes []string `json:"preferred_property_types,omitempty"`
}

// CreateAgentRequest for creating agents with real estate company
type CreateAgentRequest struct {
	CreateUserRequest
	RealEstateCompanyID string `json:"real_estate_company_id" validate:"required"`
}

// UpdateUserRequest for updating basic user information
type UpdateUserRequest struct {
	FirstName   string     `json:"first_name" validate:"required,min=2,max=100"`
	LastName    string     `json:"last_name" validate:"required,min=2,max=100"`
	Email       string     `json:"email" validate:"required,email"`
	Phone       string     `json:"phone" validate:"required"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	AvatarURL   string     `json:"avatar_url,omitempty"`
	Bio         string     `json:"bio,omitempty"`
}

// UpdateBuyerPreferencesRequest for updating buyer preferences
type UpdateBuyerPreferencesRequest struct {
	MinBudget              *float64 `json:"min_budget,omitempty"`
	MaxBudget              *float64 `json:"max_budget,omitempty"`
	PreferredProvinces     []string `json:"preferred_provinces,omitempty"`
	PreferredPropertyTypes []string `json:"preferred_property_types,omitempty"`
}

// RegisterUserRoutes registers all user routes
func (h *UserHandler) RegisterUserRoutes(router *mux.Router) {
	// Main routes
	router.HandleFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/users", h.GetUsers).Methods("GET")
	router.HandleFunc("/users/buyers", h.CreateBuyer).Methods("POST")
	router.HandleFunc("/users/agents", h.CreateAgent).Methods("POST")
	router.HandleFunc("/users/search", h.SearchUsers).Methods("GET")
	router.HandleFunc("/users/statistics", h.GetStatistics).Methods("GET")

	// Routes by type
	router.HandleFunc("/users/buyers", h.GetBuyers).Methods("GET")
	router.HandleFunc("/users/sellers", h.GetSellers).Methods("GET")
	router.HandleFunc("/users/agents", h.GetAgents).Methods("GET")

	// Routes with ID
	router.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}/preferences", h.UpdateBuyerPreferences).Methods("PUT")
	router.HandleFunc("/users/{id}/company", h.ChangeRealEstateCompany).Methods("PUT")
	router.HandleFunc("/users/{id}/notifications", h.SetNotificationPreferences).Methods("PUT")
	router.HandleFunc("/users/{id}/activate", h.ActivateUser).Methods("PUT")
	router.HandleFunc("/users/{id}/deactivate", h.DeactivateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

	// Routes by email and national ID
	router.HandleFunc("/users/email/{email}", h.GetUserByEmail).Methods("GET")
	router.HandleFunc("/users/national-id/{national_id}", h.GetUserByNationalID).Methods("GET")

	// Specific search routes
	router.HandleFunc("/users/buyers/property", h.GetBuyersForProperty).Methods("GET")
	router.HandleFunc("/companies/{company_id}/agents", h.GetAgentsByCompany).Methods("GET")

	// Validation routes
	router.HandleFunc("/users/validate/national-id", h.ValidateNationalID).Methods("POST")
	router.HandleFunc("/users/validate/email", h.ValidateEmail).Methods("POST")
	router.HandleFunc("/users/validate/phone", h.ValidatePhone).Methods("POST")
	router.HandleFunc("/users/check/email", h.CheckEmailAvailability).Methods("POST")
	router.HandleFunc("/users/check/national-id", h.CheckNationalIDAvailability).Methods("POST")
}

// CreateUser creates a new basic user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Create(req.FirstName, req.LastName, req.Email, req.Phone, req.NationalID, req.UserType)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// CreateBuyer creates a new buyer user with preferences
func (h *UserHandler) CreateBuyer(w http.ResponseWriter, r *http.Request) {
	var req CreateBuyerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Force user type to buyer
	req.UserType = "buyer"

	user, err := h.service.CreateBuyer(
		req.FirstName, req.LastName, req.Email, req.Phone, req.NationalID,
		req.MinBudget, req.MaxBudget, req.PreferredProvinces, req.PreferredPropertyTypes,
	)
	if err != nil {
		log.Printf("Error creating buyer: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// CreateAgent creates a new agent associated with a real estate company
func (h *UserHandler) CreateAgent(w http.ResponseWriter, r *http.Request) {
	var req CreateAgentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Force user type to agent
	req.UserType = "agent"

	user, err := h.service.CreateAgent(
		req.FirstName, req.LastName, req.Email, req.Phone, req.NationalID, req.RealEstateCompanyID,
	)
	if err != nil {
		log.Printf("Error creating agent: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUsers retrieves all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
		"total": len(users),
	})
}

// GetBuyers retrieves all buyer users
func (h *UserHandler) GetBuyers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetBuyers()
	if err != nil {
		log.Printf("Error getting buyers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"buyers": users,
		"total":  len(users),
	})
}

// GetSellers retrieves all seller users
func (h *UserHandler) GetSellers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetSellers()
	if err != nil {
		log.Printf("Error getting sellers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sellers": users,
		"total":   len(users),
	})
}

// GetAgents retrieves all agent users
func (h *UserHandler) GetAgents(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAgents()
	if err != nil {
		log.Printf("Error getting agents: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"agents": users,
		"total":  len(users),
	})
}

// GetAgentsByCompany retrieves agents for a specific real estate company
func (h *UserHandler) GetAgentsByCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["company_id"]

	users, err := h.service.GetAgentsByCompany(companyID)
	if err != nil {
		log.Printf("Error getting agents for company %s: %v", companyID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"agents":     users,
		"total":      len(users),
		"company_id": companyID,
	})
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.GetByID(id)
	if err != nil {
		log.Printf("Error getting user %s: %v", id, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserByEmail retrieves a user by email
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	user, err := h.service.GetByEmail(email)
	if err != nil {
		log.Printf("Error getting user by email %s: %v", email, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserByNationalID retrieves a user by national ID
func (h *UserHandler) GetUserByNationalID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nationalID := vars["national_id"]

	user, err := h.service.GetByNationalID(nationalID)
	if err != nil {
		log.Printf("Error getting user by national ID %s: %v", nationalID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates basic user information
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Update(
		id, req.FirstName, req.LastName, req.Email, req.Phone,
		req.DateOfBirth, req.AvatarURL, req.Bio,
	)
	if err != nil {
		log.Printf("Error updating user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateBuyerPreferences updates buyer search preferences
func (h *UserHandler) UpdateBuyerPreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateBuyerPreferencesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateBuyerPreferences(
		id, req.MinBudget, req.MaxBudget,
		req.PreferredProvinces, req.PreferredPropertyTypes,
	)
	if err != nil {
		log.Printf("Error updating buyer preferences for user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ChangeRealEstateCompany changes the real estate company for an agent
func (h *UserHandler) ChangeRealEstateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		RealEstateCompanyID string `json:"real_estate_company_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.ChangeRealEstateCompany(id, req.RealEstateCompanyID)
	if err != nil {
		log.Printf("Error changing real estate company for user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// SetNotificationPreferences sets user notification preferences
func (h *UserHandler) SetNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		ReceiveNotifications bool `json:"receive_notifications"`
		ReceiveNewsletter    bool `json:"receive_newsletter"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.SetNotificationPreferences(id, req.ReceiveNotifications, req.ReceiveNewsletter)
	if err != nil {
		log.Printf("Error setting notification preferences for user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ActivateUser activates a user
func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.Activate(id)
	if err != nil {
		log.Printf("Error activating user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeactivateUser deactivates a user
func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.Deactivate(id)
	if err != nil {
		log.Printf("Error deactivating user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deactivated successfully",
		"id":      id,
	})
}

// SearchUsers searches users by name
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		http.Error(w, "Search term 'q' is required", http.StatusBadRequest)
		return
	}

	users, err := h.service.SearchByName(searchTerm)
	if err != nil {
		log.Printf("Error searching users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users":       users,
		"total":       len(users),
		"search_term": searchTerm,
	})
}

// GetBuyersForProperty gets buyers that can afford a specific property price
func (h *UserHandler) GetBuyersForProperty(w http.ResponseWriter, r *http.Request) {
	priceStr := r.URL.Query().Get("price")
	if priceStr == "" {
		http.Error(w, "Property price parameter 'price' is required", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Price must be a valid number", http.StatusBadRequest)
		return
	}

	users, err := h.service.GetBuyersForProperty(price)
	if err != nil {
		log.Printf("Error getting buyers for property: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"buyers": users,
		"total":  len(users),
		"price":  price,
	})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.Delete(id)
	if err != nil {
		log.Printf("Error deleting user %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
		"id":      id,
	})
}

// ValidateNationalID validates a national ID
func (h *UserHandler) ValidateNationalID(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NationalID string `json:"national_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	err := h.service.ValidateNationalID(req.NationalID)
	isValid := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"national_id": req.NationalID,
		"valid":       isValid,
		"message": func() string {
			if isValid {
				return "National ID is valid"
			}
			return err.Error()
		}(),
	})
}

// ValidateEmail validates an email
func (h *UserHandler) ValidateEmail(w http.ResponseWriter, r *http.Request) {
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
func (h *UserHandler) ValidatePhone(w http.ResponseWriter, r *http.Request) {
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

// CheckEmailAvailability checks if an email is available
func (h *UserHandler) CheckEmailAvailability(w http.ResponseWriter, r *http.Request) {
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

// CheckNationalIDAvailability checks if a national ID is available
func (h *UserHandler) CheckNationalIDAvailability(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NationalID string `json:"national_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	available, err := h.service.CheckNationalIDAvailability(req.NationalID)
	if err != nil {
		log.Printf("Error checking national ID availability: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"national_id": req.NationalID,
		"available":   available,
		"message": func() string {
			if available {
				return "National ID is available"
			}
			return "National ID is already in use"
		}(),
	})
}

// GetStatistics returns statistics about users
func (h *UserHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		log.Printf("Error getting user statistics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
