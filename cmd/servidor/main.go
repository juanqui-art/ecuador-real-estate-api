package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"realty-core/internal/repositorio"
	"realty-core/internal/servicio"
	"realty-core/internal/web"
	"realty-core/internal/web/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables del sistema")
	}

	// Obtener configuración de variables de entorno
	databaseURL := obtenerVariable("DATABASE_URL", "postgresql://admin:password@localhost:5432/inmobiliaria_db")
	puerto := obtenerVariable("PORT", "8080")
	logLevel := obtenerVariable("LOG_LEVEL", "info")

	log.Printf("Iniciando servidor de propiedades inmobiliarias...")
	log.Printf("Log Level: %s", logLevel)
	log.Printf("Puerto: %s", puerto)

	// Conectar a la base de datos
	log.Println("Conectando a PostgreSQL...")
	db, err := repositorio.ConectarBaseDatos(databaseURL)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	// Crear repositorio
	propiedadRepo := repositorio.NuevoPropiedadRepositoryPostgres(db)

	// Crear servicio
	propiedadService := servicio.NuevoPropiedadService(propiedadRepo)

	// Crear handler
	propiedadHandler := handlers.NuevoPropiedadHandler(propiedadService)

	// Configurar rutas
	router := web.ConfigurarRutas(propiedadHandler)

	// Configurar servidor HTTP
	servidor := &http.Server{
		Addr:         ":" + puerto,
		Handler:      agregarMiddleware(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Servidor iniciado en http://localhost:%s", puerto)
	log.Printf("📚 Documentación disponible en http://localhost:%s/", puerto)
	log.Printf("🏥 Health check en http://localhost:%s/api/salud", puerto)

	// Iniciar servidor
	if err := servidor.ListenAndServe(); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}

// obtenerVariable obtiene una variable de entorno o devuelve un valor por defecto
func obtenerVariable(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// agregarMiddleware agrega middleware común a todas las rutas
func agregarMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Middleware de logging
		start := time.Now()

		// Configurar headers CORS para desarrollo
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Manejar preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Ejecutar el handler siguiente
		next.ServeHTTP(w, r)

		// Log de la petición
		duration := time.Since(start)
		log.Printf("%s %s - %s", r.Method, r.URL.Path, duration)
	})
}
