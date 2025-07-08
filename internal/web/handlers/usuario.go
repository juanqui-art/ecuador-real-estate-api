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

// UsuarioHandler maneja las peticiones HTTP relacionadas con usuarios
type UsuarioHandler struct {
	servicio *servicio.UsuarioService
}

// NuevoUsuarioHandler crea una nueva instancia del handler
func NuevoUsuarioHandler(servicio *servicio.UsuarioService) *UsuarioHandler {
	return &UsuarioHandler{servicio: servicio}
}

// CrearUsuarioRequest representa la estructura de datos para crear un usuario
type CrearUsuarioRequest struct {
	Nombre      string `json:"nombre" validate:"required,min=2,max=100"`
	Apellido    string `json:"apellido" validate:"required,min=2,max=100"`
	Email       string `json:"email" validate:"required,email"`
	Telefono    string `json:"telefono" validate:"required"`
	Cedula      string `json:"cedula" validate:"required,len=10"`
	TipoUsuario string `json:"tipo_usuario" validate:"required,oneof=comprador vendedor agente admin"`
}

// CrearCompradorRequest para crear compradores con preferencias
type CrearCompradorRequest struct {
	CrearUsuarioRequest
	PresupuestoMin        *float64 `json:"presupuesto_min,omitempty"`
	PresupuestoMax        *float64 `json:"presupuesto_max,omitempty"`
	ProvinciasInteres     []string `json:"provincias_interes,omitempty"`
	TiposPropiedadInteres []string `json:"tipos_propiedad_interes,omitempty"`
}

// CrearAgenteRequest para crear agentes con inmobiliaria
type CrearAgenteRequest struct {
	CrearUsuarioRequest
	InmobiliariaID string `json:"inmobiliaria_id" validate:"required"`
}

// ActualizarUsuarioRequest para actualizar datos básicos del usuario
type ActualizarUsuarioRequest struct {
	Nombre          string     `json:"nombre" validate:"required,min=2,max=100"`
	Apellido        string     `json:"apellido" validate:"required,min=2,max=100"`
	Email           string     `json:"email" validate:"required,email"`
	Telefono        string     `json:"telefono" validate:"required"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	AvatarURL       string     `json:"avatar_url,omitempty"`
	Bio             string     `json:"bio,omitempty"`
}

// ActualizarPreferenciasRequest para actualizar preferencias de comprador
type ActualizarPreferenciasRequest struct {
	PresupuestoMin        *float64 `json:"presupuesto_min,omitempty"`
	PresupuestoMax        *float64 `json:"presupuesto_max,omitempty"`
	ProvinciasInteres     []string `json:"provincias_interes,omitempty"`
	TiposPropiedadInteres []string `json:"tipos_propiedad_interes,omitempty"`
}

// RegistrarRutasUsuario registra todas las rutas para usuarios
func (h *UsuarioHandler) RegistrarRutasUsuario(router *mux.Router) {
	// Rutas principales
	router.HandleFunc("/usuarios", h.CrearUsuario).Methods("POST")
	router.HandleFunc("/usuarios", h.ObtenerUsuarios).Methods("GET")
	router.HandleFunc("/usuarios/compradores", h.CrearComprador).Methods("POST")
	router.HandleFunc("/usuarios/agentes", h.CrearAgente).Methods("POST")
	router.HandleFunc("/usuarios/buscar", h.BuscarUsuarios).Methods("GET")
	router.HandleFunc("/usuarios/estadisticas", h.ObtenerEstadisticas).Methods("GET")

	// Rutas por tipo
	router.HandleFunc("/usuarios/compradores", h.ObtenerCompradores).Methods("GET")
	router.HandleFunc("/usuarios/vendedores", h.ObtenerVendedores).Methods("GET")
	router.HandleFunc("/usuarios/agentes", h.ObtenerAgentes).Methods("GET")

	// Rutas con ID
	router.HandleFunc("/usuarios/{id}", h.ObtenerUsuario).Methods("GET")
	router.HandleFunc("/usuarios/{id}", h.ActualizarUsuario).Methods("PUT")
	router.HandleFunc("/usuarios/{id}/preferencias", h.ActualizarPreferencias).Methods("PUT")
	router.HandleFunc("/usuarios/{id}/inmobiliaria", h.CambiarInmobiliaria).Methods("PUT")
	router.HandleFunc("/usuarios/{id}/desactivar", h.DesactivarUsuario).Methods("PUT")
	router.HandleFunc("/usuarios/{id}/reactivar", h.ReactivarUsuario).Methods("PUT")

	// Rutas por email y cédula
	router.HandleFunc("/usuarios/email/{email}", h.ObtenerUsuarioPorEmail).Methods("GET")
	router.HandleFunc("/usuarios/cedula/{cedula}", h.ObtenerUsuarioPorCedula).Methods("GET")

	// Rutas de búsqueda específica
	router.HandleFunc("/usuarios/compradores/propiedad", h.BuscarCompradoresParaPropiedad).Methods("GET")
	router.HandleFunc("/inmobiliarias/{inmobiliaria_id}/agentes", h.ObtenerAgentesPorInmobiliaria).Methods("GET")

	// Rutas de validación
	router.HandleFunc("/usuarios/validar/cedula", h.ValidarCedula).Methods("POST")
	router.HandleFunc("/usuarios/validar/email", h.ValidarEmail).Methods("POST")
	router.HandleFunc("/usuarios/validar/telefono", h.ValidarTelefono).Methods("POST")
}

// CrearUsuario crea un nuevo usuario básico
func (h *UsuarioHandler) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	var req CrearUsuarioRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	usuario, err := h.servicio.Crear(req.Nombre, req.Apellido, req.Email, req.Telefono, req.Cedula, req.TipoUsuario)
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

// CrearComprador crea un nuevo usuario comprador con preferencias
func (h *UsuarioHandler) CrearComprador(w http.ResponseWriter, r *http.Request) {
	var req CrearCompradorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Forzar tipo comprador
	req.TipoUsuario = "comprador"

	usuario, err := h.servicio.CrearComprador(
		req.Nombre, req.Apellido, req.Email, req.Telefono, req.Cedula,
		req.PresupuestoMin, req.PresupuestoMax, req.ProvinciasInteres, req.TiposPropiedadInteres,
	)
	if err != nil {
		log.Printf("Error al crear comprador: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

// CrearAgente crea un nuevo agente asociado a una inmobiliaria
func (h *UsuarioHandler) CrearAgente(w http.ResponseWriter, r *http.Request) {
	var req CrearAgenteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Forzar tipo agente
	req.TipoUsuario = "agente"

	usuario, err := h.servicio.CrearAgente(
		req.Nombre, req.Apellido, req.Email, req.Telefono, req.Cedula, req.InmobiliariaID,
	)
	if err != nil {
		log.Printf("Error al crear agente: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

// ObtenerUsuarios obtiene todos los usuarios
func (h *UsuarioHandler) ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.servicio.ObtenerTodos()
	if err != nil {
		log.Printf("Error al obtener usuarios: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"usuarios": usuarios,
		"total":    len(usuarios),
	})
}

// ObtenerCompradores obtiene todos los usuarios compradores
func (h *UsuarioHandler) ObtenerCompradores(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.servicio.ObtenerCompradores()
	if err != nil {
		log.Printf("Error al obtener compradores: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"compradores": usuarios,
		"total":       len(usuarios),
	})
}

// ObtenerVendedores obtiene todos los usuarios vendedores
func (h *UsuarioHandler) ObtenerVendedores(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.servicio.ObtenerVendedores()
	if err != nil {
		log.Printf("Error al obtener vendedores: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"vendedores": usuarios,
		"total":      len(usuarios),
	})
}

// ObtenerAgentes obtiene todos los usuarios agentes
func (h *UsuarioHandler) ObtenerAgentes(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.servicio.ObtenerAgentes()
	if err != nil {
		log.Printf("Error al obtener agentes: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"agentes": usuarios,
		"total":   len(usuarios),
	})
}

// ObtenerAgentesPorInmobiliaria obtiene agentes de una inmobiliaria específica
func (h *UsuarioHandler) ObtenerAgentesPorInmobiliaria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inmobiliariaID := vars["inmobiliaria_id"]

	usuarios, err := h.servicio.ObtenerAgentesPorInmobiliaria(inmobiliariaID)
	if err != nil {
		log.Printf("Error al obtener agentes de inmobiliaria %s: %v", inmobiliariaID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"agentes":         usuarios,
		"total":           len(usuarios),
		"inmobiliaria_id": inmobiliariaID,
	})
}

// ObtenerUsuario obtiene un usuario por ID
func (h *UsuarioHandler) ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	usuario, err := h.servicio.ObtenerPorID(id)
	if err != nil {
		log.Printf("Error al obtener usuario %s: %v", id, err)
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// ObtenerUsuarioPorEmail obtiene un usuario por email
func (h *UsuarioHandler) ObtenerUsuarioPorEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	usuario, err := h.servicio.ObtenerPorEmail(email)
	if err != nil {
		log.Printf("Error al obtener usuario por email %s: %v", email, err)
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// ObtenerUsuarioPorCedula obtiene un usuario por cédula
func (h *UsuarioHandler) ObtenerUsuarioPorCedula(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cedula := vars["cedula"]

	usuario, err := h.servicio.ObtenerPorCedula(cedula)
	if err != nil {
		log.Printf("Error al obtener usuario por cédula %s: %v", cedula, err)
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// ActualizarUsuario actualiza los datos básicos de un usuario
func (h *UsuarioHandler) ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req ActualizarUsuarioRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	usuario, err := h.servicio.Actualizar(
		id, req.Nombre, req.Apellido, req.Email, req.Telefono,
		req.FechaNacimiento, req.AvatarURL, req.Bio,
	)
	if err != nil {
		log.Printf("Error al actualizar usuario %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// ActualizarPreferencias actualiza las preferencias de un comprador
func (h *UsuarioHandler) ActualizarPreferencias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req ActualizarPreferenciasRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	usuario, err := h.servicio.ActualizarPreferenciasComprador(
		id, req.PresupuestoMin, req.PresupuestoMax,
		req.ProvinciasInteres, req.TiposPropiedadInteres,
	)
	if err != nil {
		log.Printf("Error al actualizar preferencias de usuario %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// CambiarInmobiliaria cambia la inmobiliaria de un agente
func (h *UsuarioHandler) CambiarInmobiliaria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		InmobiliariaID string `json:"inmobiliaria_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	usuario, err := h.servicio.CambiarInmobiliaria(id, req.InmobiliariaID)
	if err != nil {
		log.Printf("Error al cambiar inmobiliaria de usuario %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// DesactivarUsuario desactiva un usuario
func (h *UsuarioHandler) DesactivarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.servicio.Desactivar(id)
	if err != nil {
		log.Printf("Error al desactivar usuario %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Usuario desactivado exitosamente",
		"id":      id,
	})
}

// ReactivarUsuario reactiva un usuario
func (h *UsuarioHandler) ReactivarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	usuario, err := h.servicio.Reactivar(id)
	if err != nil {
		log.Printf("Error al reactivar usuario %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// BuscarUsuarios busca usuarios por nombre
func (h *UsuarioHandler) BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nombre := r.URL.Query().Get("nombre")
	if nombre == "" {
		http.Error(w, "Parámetro 'nombre' requerido", http.StatusBadRequest)
		return
	}

	usuarios, err := h.servicio.BuscarPorNombre(nombre)
	if err != nil {
		log.Printf("Error al buscar usuarios: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"usuarios": usuarios,
		"total":    len(usuarios),
		"termino":  nombre,
	})
}

// BuscarCompradoresParaPropiedad busca compradores que puedan pagar una propiedad
func (h *UsuarioHandler) BuscarCompradoresParaPropiedad(w http.ResponseWriter, r *http.Request) {
	precioStr := r.URL.Query().Get("precio")
	if precioStr == "" {
		http.Error(w, "Parámetro 'precio' requerido", http.StatusBadRequest)
		return
	}

	precio, err := strconv.ParseFloat(precioStr, 64)
	if err != nil {
		http.Error(w, "Precio debe ser un número válido", http.StatusBadRequest)
		return
	}

	usuarios, err := h.servicio.BuscarCompradoresParaPropiedad(precio)
	if err != nil {
		log.Printf("Error al buscar compradores para propiedad: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"compradores": usuarios,
		"total":       len(usuarios),
		"precio":      precio,
	})
}

// ValidarCedula valida una cédula ecuatoriana
func (h *UsuarioHandler) ValidarCedula(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Cedula string `json:"cedula"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	err := h.servicio.ValidarCedula(req.Cedula)
	esValida := err == nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cedula": req.Cedula,
		"valida": esValida,
		"mensaje": func() string {
			if esValida {
				return "Cédula válida"
			}
			return err.Error()
		}(),
	})
}

// ValidarEmail valida un email
func (h *UsuarioHandler) ValidarEmail(w http.ResponseWriter, r *http.Request) {
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
func (h *UsuarioHandler) ValidarTelefono(w http.ResponseWriter, r *http.Request) {
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

// ObtenerEstadisticas obtiene estadísticas de usuarios
func (h *UsuarioHandler) ObtenerEstadisticas(w http.ResponseWriter, r *http.Request) {
	estadisticas, err := h.servicio.ObtenerEstadisticas()
	if err != nil {
		log.Printf("Error al obtener estadísticas de usuarios: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estadisticas)
}
