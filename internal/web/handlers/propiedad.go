package handlers

import (
	"encoding/json"

	"log"
	"net/http"
	"strconv"
	"strings"

	"realty-core/internal/servicio"
)

// PropiedadHandler maneja las peticiones HTTP para propiedades
type PropiedadHandler struct {
	servicio *servicio.PropiedadService
}

// NuevoPropiedadHandler crea una nueva instancia del handler
func NuevoPropiedadHandler(servicio *servicio.PropiedadService) *PropiedadHandler {
	return &PropiedadHandler{servicio: servicio}
}

// CrearPropiedadRequest representa la estructura de la petición para crear una propiedad
type CrearPropiedadRequest struct {
	Titulo      string  `json:"titulo"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Provincia   string  `json:"provincia"`
	Ciudad      string  `json:"ciudad"`
	Tipo        string  `json:"tipo"`
}

// FiltroRequest representa los filtros para buscar propiedades
type FiltroRequest struct {
	Provincia string  `json:"provincia"`
	PrecioMin float64 `json:"precio_min"`
	PrecioMax float64 `json:"precio_max"`
}

// ErrorResponse representa una respuesta de error estandarizada
type ErrorResponse struct {
	Error   string `json:"error"`
	Mensaje string `json:"mensaje"`
}

// SuccessResponse representa una respuesta exitosa estandarizada
type SuccessResponse struct {
	Datos   interface{} `json:"datos"`
	Mensaje string      `json:"mensaje"`
}

// CrearPropiedad maneja POST /api/propiedades
func (h *PropiedadHandler) CrearPropiedad(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método POST
	if r.Method != http.MethodPost {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Decodificar JSON del cuerpo de la petición
	var req CrearPropiedadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responderError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	// Llamar al servicio para crear la propiedad
	propiedad, err := h.servicio.CrearPropiedad(
		req.Titulo,
		req.Descripcion,
		req.Provincia,
		req.Ciudad,
		req.Tipo,
		req.Precio,
	)

	if err != nil {
		h.responderError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Responder con la propiedad creada
	h.responderExito(w, http.StatusCreated, propiedad, "Propiedad creada exitosamente")
}

// ObtenerPropiedad maneja GET /api/propiedades/{id}
func (h *PropiedadHandler) ObtenerPropiedad(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método GET
	if r.Method != http.MethodGet {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Extraer ID de la URL
	id := h.extraerIDDeURL(r.URL.Path)
	if id == "" {
		h.responderError(w, http.StatusBadRequest, "ID de propiedad requerido")
		return
	}

	// Obtener la propiedad
	propiedad, err := h.servicio.ObtenerPropiedad(id)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrada") {
			h.responderError(w, http.StatusNotFound, err.Error())
		} else {
			h.responderError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Responder con la propiedad
	h.responderExito(w, http.StatusOK, propiedad, "Propiedad obtenida exitosamente")
}

// ObtenerPropiedadPorSlug maneja GET /api/propiedades/slug/{slug}
func (h *PropiedadHandler) ObtenerPropiedadPorSlug(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método GET
	if r.Method != http.MethodGet {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Extraer slug de la URL
	slug := h.extraerSlugDeURL(r.URL.Path)
	if slug == "" {
		h.responderError(w, http.StatusBadRequest, "Slug de propiedad requerido")
		return
	}

	// Obtener la propiedad por slug
	propiedad, err := h.servicio.ObtenerPropiedadPorSlug(slug)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrada") {
			h.responderError(w, http.StatusNotFound, err.Error())
		} else {
			h.responderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// Responder con la propiedad
	h.responderExito(w, http.StatusOK, propiedad, "Propiedad obtenida por slug exitosamente")
}

// ListarPropiedades maneja GET /api/propiedades
func (h *PropiedadHandler) ListarPropiedades(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método GET
	if r.Method != http.MethodGet {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Obtener todas las propiedades
	propiedades, err := h.servicio.ListarPropiedades()
	if err != nil {
		h.responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Responder con las propiedades
	h.responderExito(w, http.StatusOK, propiedades, "Propiedades obtenidas exitosamente")
}

// ActualizarPropiedad maneja PUT /api/propiedades/{id}
func (h *PropiedadHandler) ActualizarPropiedad(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método PUT
	if r.Method != http.MethodPut {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Extraer ID de la URL
	id := h.extraerIDDeURL(r.URL.Path)
	if id == "" {
		h.responderError(w, http.StatusBadRequest, "ID de propiedad requerido")
		return
	}

	// Decodificar JSON del cuerpo de la petición
	var req CrearPropiedadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responderError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	// Llamar al servicio para actualizar la propiedad
	propiedad, err := h.servicio.ActualizarPropiedad(
		id,
		req.Titulo,
		req.Descripcion,
		req.Provincia,
		req.Ciudad,
		req.Tipo,
		req.Precio,
	)

	if err != nil {
		if strings.Contains(err.Error(), "no encontrada") {
			h.responderError(w, http.StatusNotFound, err.Error())
		} else {
			h.responderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// Responder con la propiedad actualizada
	h.responderExito(w, http.StatusOK, propiedad, "Propiedad actualizada exitosamente")
}

// EliminarPropiedad maneja DELETE /api/propiedades/{id}
func (h *PropiedadHandler) EliminarPropiedad(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método DELETE
	if r.Method != http.MethodDelete {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Extraer ID de la URL
	id := h.extraerIDDeURL(r.URL.Path)
	if id == "" {
		h.responderError(w, http.StatusBadRequest, "ID de propiedad requerido")
		return
	}

	// Eliminar la propiedad
	err := h.servicio.EliminarPropiedad(id)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrada") {
			h.responderError(w, http.StatusNotFound, err.Error())
		} else {
			h.responderError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Responder con confirmación
	h.responderExito(w, http.StatusOK, nil, "Propiedad eliminada exitosamente")
}

// FiltrarPropiedades maneja GET /api/propiedades/filtrar
func (h *PropiedadHandler) FiltrarPropiedades(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método GET
	if r.Method != http.MethodGet {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Obtener parámetros de query
	query := r.URL.Query()
	provincia := query.Get("provincia")
	precioMinStr := query.Get("precio_min")
	precioMaxStr := query.Get("precio_max")

	// Filtrar por provincia si se proporciona
	if provincia != "" {
		propiedades, err := h.servicio.FiltrarPorProvincia(provincia)
		if err != nil {
			h.responderError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.responderExito(w, http.StatusOK, propiedades, "Propiedades filtradas por provincia")
		return
	}

	// Filtrar por rango de precio si se proporciona
	if precioMinStr != "" && precioMaxStr != "" {
		precioMin, err := strconv.ParseFloat(precioMinStr, 64)
		if err != nil {
			h.responderError(w, http.StatusBadRequest, "Precio mínimo inválido")
			return
		}

		precioMax, err := strconv.ParseFloat(precioMaxStr, 64)
		if err != nil {
			h.responderError(w, http.StatusBadRequest, "Precio máximo inválido")
			return
		}

		propiedades, err := h.servicio.FiltrarPorRangoPrecio(precioMin, precioMax)
		if err != nil {
			h.responderError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.responderExito(w, http.StatusOK, propiedades, "Propiedades filtradas por precio")
		return
	}

	// Si no hay filtros, devolver todas las propiedades
	propiedades, err := h.servicio.ListarPropiedades()
	if err != nil {
		h.responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.responderExito(w, http.StatusOK, propiedades, "Todas las propiedades")
}

// ObtenerEstadisticas maneja GET /api/propiedades/estadisticas
func (h *PropiedadHandler) ObtenerEstadisticas(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método GET
	if r.Method != http.MethodGet {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Obtener estadísticas
	stats, err := h.servicio.ObtenerEstadisticas()
	if err != nil {
		h.responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Responder con las estadísticas
	h.responderExito(w, http.StatusOK, stats, "Estadísticas obtenidas exitosamente")
}

// HealthCheck maneja GET /api/salud
func (h *PropiedadHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Verificar que sea método GET
	if r.Method != http.MethodGet {
		h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	salud := map[string]string{
		"estado":   "saludable",
		"servicio": "api-propiedades",
		"version":  "1.0.0",
	}

	h.responderExito(w, http.StatusOK, salud, "Servicio funcionando correctamente")
}

// Métodos auxiliares privados

// extraerIDDeURL extrae el ID de la URL (último segmento)
func (h *PropiedadHandler) extraerIDDeURL(path string) string {
	// Eliminar slash final si existe
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	// Dividir por slash y obtener último segmento
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

// extraerSlugDeURL extrae el slug de la URL para rutas /api/propiedades/slug/{slug}
func (h *PropiedadHandler) extraerSlugDeURL(path string) string {
	// Eliminar slash final si existe
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	// Dividir por slash
	parts := strings.Split(path, "/")

	// Buscar el patrón /api/propiedades/slug/{slug}
	// parts debería ser: ["", "api", "propiedades", "slug", "{slug}"]
	if len(parts) >= 5 && parts[1] == "api" && parts[2] == "propiedades" && parts[3] == "slug" {
		return parts[4]
	}

	return ""
}

// responderError envía una respuesta de error en formato JSON
func (h *PropiedadHandler) responderError(w http.ResponseWriter, status int, mensaje string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResp := ErrorResponse{
		Error:   http.StatusText(status),
		Mensaje: mensaje,
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("Error al codificar respuesta de error: %v", err)
	}
}

// responderExito envía una respuesta exitosa en formato JSON
func (h *PropiedadHandler) responderExito(w http.ResponseWriter, status int, datos interface{}, mensaje string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	successResp := SuccessResponse{
		Datos:   datos,
		Mensaje: mensaje,
	}

	if err := json.NewEncoder(w).Encode(successResp); err != nil {
		log.Printf("Error al codificar respuesta exitosa: %v", err)
	}
}
