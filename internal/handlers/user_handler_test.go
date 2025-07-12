package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests básicos para UserHandler enfocados en validación de entrada
func TestUserHandler_InputValidation(t *testing.T) {
	// Estos tests verifican las validaciones de entrada y manejo de errores básicos

	tests := []struct {
		name           string
		method         string  
		path           string
		body           string
		handlerFunc    func(*UserHandlerSimple, http.ResponseWriter, *http.Request)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "GetUser with empty ID",
			method: http.MethodGet,
			path:   "/api/users/",
			body:   "",
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.GetUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "User ID required",
		},
		{
			name:   "UpdateUser with empty ID",
			method: http.MethodPut,
			path:   "/api/users/",
			body:   "",
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.UpdateUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "User ID required",
		},
		{
			name:   "DeleteUser with empty ID",
			method: http.MethodDelete,
			path:   "/api/users/",
			body:   "",
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.DeleteUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "User ID required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandlerSimple{}
			
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}
			rr := httptest.NewRecorder()
			
			tt.handlerFunc(handler, rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

func TestUserHandler_InvalidJSON(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		handlerFunc    func(*UserHandlerSimple, http.ResponseWriter, *http.Request)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "CreateUser with invalid JSON",
			method: http.MethodPost,
			path:   "/api/users",
			body:   `{"invalid": json}`,
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.CreateUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON",
		},
		{
			name:   "UpdateUser with invalid JSON",
			method: http.MethodPut,
			path:   "/api/users/test-id",
			body:   `{"invalid": json}`,
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.UpdateUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON",
		},
		{
			name:   "Login with invalid JSON",
			method: http.MethodPost,
			path:   "/api/auth/login",
			body:   `{"invalid": json}`,
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.Login(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON",
		},
		{
			name:   "ChangePassword with invalid JSON",
			method: http.MethodPost,
			path:   "/api/users/test-id/password",
			body:   `{"invalid": json}`,
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.ChangePassword(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandlerSimple{}
			
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			
			tt.handlerFunc(handler, rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

func TestUserHandler_EmptyPaths(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		handlerFunc    func(*UserHandlerSimple, http.ResponseWriter, *http.Request)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "GetUser with empty ID",
			method: http.MethodGet,
			path:   "/api/users/",
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.GetUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "User ID required",
		},
		{
			name:   "UpdateUser with empty ID",
			method: http.MethodPut,
			path:   "/api/users/",
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.UpdateUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "User ID required",
		},
		{
			name:   "DeleteUser with empty ID",
			method: http.MethodDelete,
			path:   "/api/users/",
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.DeleteUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "User ID required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandlerSimple{}
			
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			
			tt.handlerFunc(handler, rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

func TestUserHandler_MissingFields(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		handlerFunc    func(*UserHandlerSimple, http.ResponseWriter, *http.Request)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "CreateUser with missing required fields",
			method: http.MethodPost,
			path:   "/api/users",
			body:   `{"first_name": "Juan"}`, // Solo first_name, faltan otros campos
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.CreateUser(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid role",
		},
		{
			name:   "Login with missing password",
			method: http.MethodPost,
			path:   "/api/auth/login",
			body:   `{"email": "test@example.com"}`, // Solo email, falta password
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.Login(w, r)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials",
		},
		{
			name:   "ChangePassword with missing fields",
			method: http.MethodPost,
			path:   "/api/users/test-id/password",
			body:   `{"old_password": "old123"}`, // Solo old_password, falta new_password
			handlerFunc: func(h *UserHandlerSimple, w http.ResponseWriter, r *http.Request) {
				h.ChangePassword(w, r)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "all fields are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandlerSimple{}
			
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			
			tt.handlerFunc(handler, rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

// Tests básicos completados para UserHandler
// Los tests que requieren servicios reales se omiten por limitaciones de setup