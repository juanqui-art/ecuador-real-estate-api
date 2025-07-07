package web

import (
	"net/http"

	"realty-core/internal/web/handlers"
)

// ConfigurarRutas configura todas las rutas de la aplicación
func ConfigurarRutas(propiedadHandler *handlers.PropiedadHandler) *http.ServeMux {
	// Crear un nuevo ServeMux (router)
	mux := http.NewServeMux()

	// Rutas para propiedades
	mux.HandleFunc("/api/propiedades", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			propiedadHandler.ListarPropiedades(w, r)
		case http.MethodPost:
			propiedadHandler.CrearPropiedad(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Ruta para obtener/actualizar/eliminar propiedad específica
	mux.HandleFunc("/api/propiedades/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			propiedadHandler.ObtenerPropiedad(w, r)
		case http.MethodPut:
			propiedadHandler.ActualizarPropiedad(w, r)
		case http.MethodDelete:
			propiedadHandler.EliminarPropiedad(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Ruta para obtener propiedad por slug (debe ir antes de rutas genéricas)
	mux.HandleFunc("/api/propiedades/slug/", propiedadHandler.ObtenerPropiedadPorSlug)

	// Ruta para filtrar propiedades
	mux.HandleFunc("/api/propiedades/filtrar", propiedadHandler.FiltrarPropiedades)

	// Ruta para estadísticas
	mux.HandleFunc("/api/propiedades/estadisticas", propiedadHandler.ObtenerEstadisticas)

	// Ruta de health check
	mux.HandleFunc("/api/salud", propiedadHandler.HealthCheck)

	// Ruta raíz con información del API
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"mensaje": "API de Propiedades Inmobiliarias",
			"version": "1.0.0",
			"endpoints": {
				"GET /api/salud": "Estado del servicio",
				"GET /api/propiedades": "Listar todas las propiedades",
				"POST /api/propiedades": "Crear nueva propiedad",
				"GET /api/propiedades/{id}": "Obtener propiedad por ID",
				"GET /api/propiedades/slug/{slug}": "Obtener propiedad por slug SEO",
				"PUT /api/propiedades/{id}": "Actualizar propiedad",
				"DELETE /api/propiedades/{id}": "Eliminar propiedad",
				"GET /api/propiedades/filtrar": "Filtrar propiedades (query params: provincia, precio_min, precio_max)",
				"GET /api/propiedades/estadisticas": "Obtener estadísticas"
			},
			"seo": {
				"descripcion": "Los slugs SEO se generan automáticamente desde el título",
				"formato": "titulo-normalizado-{id}",
				"ejemplo": "/api/propiedades/slug/casa-moderna-samborondon-abcd1234"
			}
		}`))
	})

	return mux
}