package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"realty-core/internal/servicio"

	"github.com/gorilla/mux"
)

// InmobiliariaHandler maneja las peticiones HTTP relacionadas con inmobiliarias
type InmobiliariaHandler struct {
	servicio *servicio.InmobiliariaService
}

// NuevoInmobiliariaHandler crea una nueva instancia del handler
func NuevoInmobiliariaHandler(servicio *servicio.InmobiliariaService) *InmobiliariaHandler {
	return &InmobiliariaHandler{servicio: servicio}
}

// CrearInmobiliariaRequest representa la estructura de datos para crear una inmobiliaria
type CrearInmobiliariaRequest struct {
	Nombre      string `json:"nombre" validate:"required,min=2,max=255"`
	RUC         string `json:"ruc" validate:"required,len=13"`
	Direccion   string `json:"direccion" validate:"required,min=10,max=255"`
	Telefono    string `json:"telefono" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	SitioWeb    string `json:"sitio_web,omitempty"`
	Descripcion string `json:"descripcion,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
}

// ActualizarInmobiliariaRequest representa la estructura para actualizar una inmobiliaria
type ActualizarInmobiliariaRequest struct {
	Nombre      string `json:"nombre" validate:"required,min=2,max=255"`
	Direccion   string `json:"direccion" validate:"required,min=10,max=255"`
	Telefono    string `json:"telefono" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	SitioWeb    string `json:"sitio_web,omitempty"`
	Descripcion string `json:"descripcion,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
	Activa      *bool  `json:"activa,omitempty"`
}

// RegistrarRutasInmobiliaria registra todas las rutas para inmobiliarias
func (h *InmobiliariaHandler) RegistrarRutasInmobiliaria(router *mux.Router) {
	// Rutas principales
	router.HandleFunc("/inmobiliarias", h.CrearInmobiliaria).Methods("POST")
	router.HandleFunc("/inmobiliarias", h.ObtenerInmobiliarias).Methods("GET")
	router.HandleFunc("/inmobiliarias/activas", h.ObtenerInmobiliariasActivas).Methods("GET")
	router.HandleFunc("/inmobiliarias/buscar", h.BuscarInmobiliarias).Methods("GET")
	router.HandleFunc("/inmobiliarias/estadisticas", h.ObtenerEstadisticas).Methods("GET")

	// Rutas con ID
	router.HandleFunc("/inmobiliarias/{id}", h.ObtenerInmobiliaria).Methods("GET")
	router.HandleFunc("/inmobiliarias/{id}", h.ActualizarInmobiliaria).Methods("PUT")
	router.HandleFunc("/inmobiliarias/{id}/perfil", h.ActualizarPerfil).Methods("PUT")
	router.HandleFunc("/inmobiliarias/{id}/desactivar", h.DesactivarInmobiliaria).Methods("PUT")
	router.HandleFunc("/inmobiliarias/{id}/reactivar", h.ReactivarInmobiliaria).Methods("PUT")

	// Rutas por RUC
	router.HandleFunc("/inmobiliarias/ruc/{ruc}", h.ObtenerInmobiliariaPorRUC).Methods("GET")

	// Rutas de validación
	router.HandleFunc("/inmobiliarias/validar/ruc", h.ValidarRUC).Methods("POST")
	router.HandleFunc("/inmobiliarias/validar/email", h.ValidarEmail).Methods("POST")
	router.HandleFunc("/inmobiliarias/validar/telefono", h.ValidarTelefono).Methods("POST")
}

// CrearInmobiliaria crea una nueva inmobiliaria
func (h *InmobiliariaHandler) CrearInmobiliaria(w http.ResponseWriter, r *http.Request) {
	var req CrearInmobiliariaRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Crear inmobiliaria
	inmobiliaria, err := h.servicio.Crear(req.Nombre, req.RUC, req.Direccion, req.Telefono, req.Email)
	if err != nil {
		log.Printf("Error al crear inmobiliaria: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Actualizar campos opcionales si se proporcionaron
	if req.SitioWeb != "" || req.Descripcion != "" || req.LogoURL != "" {
		inmobiliaria, err = h.servicio.ActualizarPerfil(
			inmobiliaria.ID, inmobiliaria.Nombre, inmobiliaria.Direccion,
			inmobiliaria.Telefono, inmobiliaria.Email, req.SitioWeb, req.Descripcion, req.LogoURL,
		)
		if err != nil {
			log.Printf("Error al actualizar perfil de inmobiliaria: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inmobiliaria)
}

// ObtenerInmobiliarias obtiene todas las inmobiliarias
func (h *InmobiliariaHandler) ObtenerInmobiliarias(w http.ResponseWriter, r *http.Request) {
	inmobiliarias, err := h.servicio.ObtenerTodas()
	if err != nil {
		log.Printf("Error al obtener inmobiliarias: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"inmobiliarias": inmobiliarias,
		"total":         len(inmobiliarias),
	})
}

// ObtenerInmobiliariasActivas obtiene solo las inmobiliarias activas
func (h *InmobiliariaHandler) ObtenerInmobiliariasActivas(w http.ResponseWriter, r *http.Request) {
	inmobiliarias, err := h.servicio.ObtenerActivas()
	if err != nil {
		log.Printf("Error al obtener inmobiliarias activas: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"inmobiliarias": inmobiliarias,
		"total":         len(inmobiliarias),
	})
}

// ObtenerInmobiliaria obtiene una inmobiliaria por ID
func (h *InmobiliariaHandler) ObtenerInmobiliaria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	inmobiliaria, err := h.servicio.ObtenerPorID(id)
	if err != nil {
		log.Printf("Error al obtener inmobiliaria %s: %v", id, err)
		http.Error(w, "Inmobiliaria no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inmobiliaria)
}

// ObtenerInmobiliariaPorRUC obtiene una inmobiliaria por RUC
func (h *InmobiliariaHandler) ObtenerInmobiliariaPorRUC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ruc := vars["ruc"]

	inmobiliaria, err := h.servicio.ObtenerPorRUC(ruc)
	if err != nil {
		log.Printf("Error al obtener inmobiliaria por RUC %s: %v", ruc, err)
		http.Error(w, "Inmobiliaria no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inmobiliaria)
}

// ActualizarInmobiliaria actualiza una inmobiliaria completa
func (h *InmobiliariaHandler) ActualizarInmobiliaria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req ActualizarInmobiliariaRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Determinar si se debe cambiar el estado activo
	activa := true // Por defecto mantener activa
	if req.Activa != nil {
		activa = *req.Activa
	}

	inmobiliaria, err := h.servicio.Actualizar(
		id, req.Nombre, req.Direccion, req.Telefono, req.Email,
		req.SitioWeb, req.Descripcion, req.LogoURL, activa,
	)
	if err != nil {
		log.Printf("Error al actualizar inmobiliaria %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inmobiliaria)
}

// ActualizarPerfil actualiza solo el perfil de la inmobiliaria
func (h *InmobiliariaHandler) ActualizarPerfil(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req ActualizarInmobiliariaRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	inmobiliaria, err := h.servicio.ActualizarPerfil(
		id, req.Nombre, req.Direccion, req.Telefono, req.Email,
		req.SitioWeb, req.Descripcion, req.LogoURL,
	)
	if err != nil {
		log.Printf("Error al actualizar perfil de inmobiliaria %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inmobiliaria)
}

// DesactivarInmobiliaria desactiva una inmobiliaria
func (h *InmobiliariaHandler) DesactivarInmobiliaria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.servicio.Desactivar(id)
	if err != nil {
		log.Printf("Error al desactivar inmobiliaria %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Inmobiliaria desactivada exitosamente",
		"id":      id,
	})
}

// ReactivarInmobiliaria reactiva una inmobiliaria
func (h *InmobiliariaHandler) ReactivarInmobiliaria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	inmobiliaria, err := h.servicio.Reactivar(id)
	if err != nil {
		log.Printf("Error al reactivar inmobiliaria %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inmobiliaria)
}

// BuscarInmobiliarias busca inmobiliarias por nombre
func (h *InmobiliariaHandler) BuscarInmobiliarias(w http.ResponseWriter, r *http.Request) {
	nombre := r.URL.Query().Get("nombre")
	if nombre == "" {
		http.Error(w, "Parámetro 'nombre' requerido", http.StatusBadRequest)
		return
	}

	inmobiliarias, err := h.servicio.BuscarPorNombre(nombre)
	if err != nil {
		log.Printf("Error al buscar inmobiliarias: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"inmobiliarias": inmobiliarias,
		"total":         len(inmobiliarias),
		"termino":       nombre,
	})
}

// ValidarRUC valida un RUC ecuatoriano
func (h *InmobiliariaHandler) ValidarRUC(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RUC string `json:"ruc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	err := h.servicio.ValidarRUC(req.RUC)
	esValido := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ruc":    req.RUC,
		"valido": esValido,
		"mensaje": func() string {
			if esValido {
				return "RUC válido"
			}
			return err.Error()
		}(),
	})
}

// ValidarEmail valida un email
func (h *InmobiliariaHandler) ValidarEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	err := h.servicio.ValidarEmail(req.Email)
	esValido := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"email":  req.Email,
		"valido": esValido,
		"mensaje": func() string {
			if esValido {
				return "Email válido"
			}
			return err.Error()
		}(),
	})
}

// ValidarTelefono valida un teléfono ecuatoriano
func (h *InmobiliariaHandler) ValidarTelefono(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Telefono string `json:"telefono"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	err := h.servicio.ValidarTelefono(req.Telefono)
	esValido := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"telefono": req.Telefono,
		"valido":   esValido,
		"mensaje": func() string {
			if esValido {
				return "Teléfono válido"
			}
			return err.Error()
		}(),
	})
}

// ObtenerEstadisticas obtiene estadísticas de inmobiliarias
func (h *InmobiliariaHandler) ObtenerEstadisticas(w http.ResponseWriter, r *http.Request) {
	estadisticas, err := h.servicio.ObtenerEstadisticas()
	if err != nil {
		log.Printf("Error al obtener estadísticas de inmobiliarias: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estadisticas)
}

// Middleware para logging de requests (opcional)
func (h *InmobiliariaHandler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Inmobiliaria API: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// PaginationParams extrae parámetros de paginación de la URL
func extraerParametrosPaginacion(r *http.Request) (int, int) {
	pagina := 1
	limite := 20

	if p := r.URL.Query().Get("pagina"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			pagina = parsed
		}
	}

	if l := r.URL.Query().Get("limite"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limite = parsed
		}
	}

	return pagina, limite
}
