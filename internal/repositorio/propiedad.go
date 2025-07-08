package repositorio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"realty-core/internal/dominio"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// PropiedadRepository define las operaciones de acceso a datos
// En Go, usamos interfaces para definir contratos
type PropiedadRepository interface {
	Crear(propiedad *dominio.Propiedad) error
	ObtenerPorID(id string) (*dominio.Propiedad, error)
	ObtenerPorSlug(slug string) (*dominio.Propiedad, error)
	ObtenerTodas() ([]dominio.Propiedad, error)
	Actualizar(propiedad *dominio.Propiedad) error
	Eliminar(id string) error
}

// PropiedadRepositoryPostgres implementa PropiedadRepository usando PostgreSQL
type PropiedadRepositoryPostgres struct {
	db *sql.DB
}

// NuevoPropiedadRepositoryPostgres crea una nueva instancia del repositorio
func NuevoPropiedadRepositoryPostgres(db *sql.DB) *PropiedadRepositoryPostgres {
	return &PropiedadRepositoryPostgres{db: db}
}

// Crear inserta una nueva propiedad en la base de datos
func (r *PropiedadRepositoryPostgres) Crear(propiedad *dominio.Propiedad) error {
	// Convertir slices a JSON para almacenar en JSONB
	imagenesJSON, err := json.Marshal(propiedad.Imagenes)
	if err != nil {
		return fmt.Errorf("error al convertir imágenes a JSON: %w", err)
	}

	tagsJSON, err := json.Marshal(propiedad.Tags)
	if err != nil {
		return fmt.Errorf("error al convertir tags a JSON: %w", err)
	}

	// Query SQL para insertar una propiedad con todos los campos
	query := `
		INSERT INTO propiedades (
			id, slug, titulo, descripcion, precio, provincia, ciudad, 
			sector, direccion, tipo, estado, dormitorios, banos, 
			area_m2, fecha_creacion, fecha_actualizacion,
			latitud, longitud, precision_ubicacion,
			imagen_principal, imagenes, video_tour, tour_360,
			precio_alquiler, gastos_comunes, precio_m2,
			ano_construccion, pisos, estado_propiedad, amoblada,
			garage, piscina, jardin, terraza, balcon, seguridad, ascensor, aire_acondicionado,
			tags, destacada, visitas_contador, inmobiliaria_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42
		)
	`

	// Ejecutar la query con todos los campos
	_, err = r.db.Exec(
		query,
		propiedad.ID,
		propiedad.Slug,
		propiedad.Titulo,
		propiedad.Descripcion,
		propiedad.Precio,
		propiedad.Provincia,
		propiedad.Ciudad,
		propiedad.Sector,
		propiedad.Direccion,
		propiedad.Tipo,
		propiedad.Estado,
		propiedad.Dormitorios,
		propiedad.Banos,
		propiedad.AreaM2,
		propiedad.FechaCreacion,
		propiedad.FechaActualizacion,
		// Geolocalización
		propiedad.Latitud,
		propiedad.Longitud,
		propiedad.PrecisionUbicacion,
		// Imágenes y media
		propiedad.ImagenPrincipal,
		string(imagenesJSON),
		propiedad.VideoTour,
		propiedad.Tour360,
		// Precios adicionales
		propiedad.PrecioAlquiler,
		propiedad.GastosComunes,
		propiedad.PrecioM2,
		// Características
		propiedad.AnoConstruccion,
		propiedad.Pisos,
		propiedad.EstadoPropiedad,
		propiedad.Amoblada,
		// Amenidades
		propiedad.Garage,
		propiedad.Piscina,
		propiedad.Jardin,
		propiedad.Terraza,
		propiedad.Balcon,
		propiedad.Seguridad,
		propiedad.Ascensor,
		propiedad.AireAcondicionado,
		// Marketing y SEO
		string(tagsJSON),
		propiedad.Destacada,
		propiedad.VisitasContador,
		propiedad.InmobiliariaID,
	)

	if err != nil {
		return fmt.Errorf("error al crear propiedad: %w", err)
	}

	log.Printf("Propiedad creada exitosamente: %s", propiedad.ID)
	return nil
}

// ObtenerPorID busca una propiedad por su ID
func (r *PropiedadRepositoryPostgres) ObtenerPorID(id string) (*dominio.Propiedad, error) {
	query := `
		SELECT id, slug, titulo, descripcion, precio, provincia, ciudad, 
			   sector, direccion, tipo, estado, dormitorios, banos, 
			   area_m2, fecha_creacion, fecha_actualizacion,
			   latitud, longitud, precision_ubicacion,
			   imagen_principal, imagenes, video_tour, tour_360,
			   precio_alquiler, gastos_comunes, precio_m2,
			   ano_construccion, pisos, estado_propiedad, amoblada,
			   garage, piscina, jardin, terraza, balcon, seguridad, ascensor, aire_acondicionado,
			   tags, destacada, visitas_contador, inmobiliaria_id
		FROM propiedades 
		WHERE id = $1
	`

	var propiedad dominio.Propiedad
	var imagenesJSON, tagsJSON string

	// QueryRow ejecuta la query y devuelve una fila
	err := r.db.QueryRow(query, id).Scan(
		&propiedad.ID,
		&propiedad.Slug,
		&propiedad.Titulo,
		&propiedad.Descripcion,
		&propiedad.Precio,
		&propiedad.Provincia,
		&propiedad.Ciudad,
		&propiedad.Sector,
		&propiedad.Direccion,
		&propiedad.Tipo,
		&propiedad.Estado,
		&propiedad.Dormitorios,
		&propiedad.Banos,
		&propiedad.AreaM2,
		&propiedad.FechaCreacion,
		&propiedad.FechaActualizacion,
		// Geolocalización
		&propiedad.Latitud,
		&propiedad.Longitud,
		&propiedad.PrecisionUbicacion,
		// Imágenes y media
		&propiedad.ImagenPrincipal,
		&imagenesJSON,
		&propiedad.VideoTour,
		&propiedad.Tour360,
		// Precios adicionales
		&propiedad.PrecioAlquiler,
		&propiedad.GastosComunes,
		&propiedad.PrecioM2,
		// Características
		&propiedad.AnoConstruccion,
		&propiedad.Pisos,
		&propiedad.EstadoPropiedad,
		&propiedad.Amoblada,
		// Amenidades
		&propiedad.Garage,
		&propiedad.Piscina,
		&propiedad.Jardin,
		&propiedad.Terraza,
		&propiedad.Balcon,
		&propiedad.Seguridad,
		&propiedad.Ascensor,
		&propiedad.AireAcondicionado,
		// Marketing y SEO
		&tagsJSON,
		&propiedad.Destacada,
		&propiedad.VisitasContador,
		&propiedad.InmobiliariaID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("propiedad no encontrada: %s", id)
		}
		return nil, fmt.Errorf("error al obtener propiedad: %w", err)
	}

	// Convertir JSON de vuelta a slices
	if imagenesJSON != "" {
		err = json.Unmarshal([]byte(imagenesJSON), &propiedad.Imagenes)
		if err != nil {
			return nil, fmt.Errorf("error al convertir imágenes desde JSON: %w", err)
		}
	}

	if tagsJSON != "" {
		err = json.Unmarshal([]byte(tagsJSON), &propiedad.Tags)
		if err != nil {
			return nil, fmt.Errorf("error al convertir tags desde JSON: %w", err)
		}
	}

	return &propiedad, nil
}

// ObtenerPorSlug busca una propiedad por su slug SEO
func (r *PropiedadRepositoryPostgres) ObtenerPorSlug(slug string) (*dominio.Propiedad, error) {
	query := `
		SELECT id, slug, titulo, descripcion, precio, provincia, ciudad, 
			   sector, direccion, tipo, estado, dormitorios, banos, 
			   area_m2, fecha_creacion, fecha_actualizacion,
			   latitud, longitud, precision_ubicacion,
			   imagen_principal, imagenes, video_tour, tour_360,
			   precio_alquiler, gastos_comunes, precio_m2,
			   ano_construccion, pisos, estado_propiedad, amoblada,
			   garage, piscina, jardin, terraza, balcon, seguridad, ascensor, aire_acondicionado,
			   tags, destacada, visitas_contador, inmobiliaria_id
		FROM propiedades 
		WHERE slug = $1
	`

	var propiedad dominio.Propiedad
	var imagenesJSON, tagsJSON string

	// QueryRow ejecuta la query y devuelve una fila
	err := r.db.QueryRow(query, slug).Scan(
		&propiedad.ID,
		&propiedad.Slug,
		&propiedad.Titulo,
		&propiedad.Descripcion,
		&propiedad.Precio,
		&propiedad.Provincia,
		&propiedad.Ciudad,
		&propiedad.Sector,
		&propiedad.Direccion,
		&propiedad.Tipo,
		&propiedad.Estado,
		&propiedad.Dormitorios,
		&propiedad.Banos,
		&propiedad.AreaM2,
		&propiedad.FechaCreacion,
		&propiedad.FechaActualizacion,
		// Geolocalización
		&propiedad.Latitud,
		&propiedad.Longitud,
		&propiedad.PrecisionUbicacion,
		// Imágenes y media
		&propiedad.ImagenPrincipal,
		&imagenesJSON,
		&propiedad.VideoTour,
		&propiedad.Tour360,
		// Precios adicionales
		&propiedad.PrecioAlquiler,
		&propiedad.GastosComunes,
		&propiedad.PrecioM2,
		// Características
		&propiedad.AnoConstruccion,
		&propiedad.Pisos,
		&propiedad.EstadoPropiedad,
		&propiedad.Amoblada,
		// Amenidades
		&propiedad.Garage,
		&propiedad.Piscina,
		&propiedad.Jardin,
		&propiedad.Terraza,
		&propiedad.Balcon,
		&propiedad.Seguridad,
		&propiedad.Ascensor,
		&propiedad.AireAcondicionado,
		// Marketing y SEO
		&tagsJSON,
		&propiedad.Destacada,
		&propiedad.VisitasContador,
		&propiedad.InmobiliariaID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("propiedad no encontrada con slug: %s", slug)
		}
		return nil, fmt.Errorf("error al obtener propiedad por slug: %w", err)
	}

	// Convertir JSON de vuelta a slices
	if imagenesJSON != "" {
		err = json.Unmarshal([]byte(imagenesJSON), &propiedad.Imagenes)
		if err != nil {
			return nil, fmt.Errorf("error al convertir imágenes desde JSON: %w", err)
		}
	}

	if tagsJSON != "" {
		err = json.Unmarshal([]byte(tagsJSON), &propiedad.Tags)
		if err != nil {
			return nil, fmt.Errorf("error al convertir tags desde JSON: %w", err)
		}
	}

	return &propiedad, nil
}

// ObtenerTodas devuelve todas las propiedades (versión ligera para listados)
func (r *PropiedadRepositoryPostgres) ObtenerTodas() ([]dominio.Propiedad, error) {
	query := `
		SELECT id, slug, titulo, descripcion, precio, provincia, ciudad, 
			   sector, direccion, tipo, estado, dormitorios, banos, 
			   area_m2, fecha_creacion, fecha_actualizacion,
			   imagen_principal, tags, destacada, visitas_contador, inmobiliaria_id
		FROM propiedades 
		ORDER BY destacada DESC, fecha_creacion DESC
	`

	// Query ejecuta la consulta y devuelve múltiples filas
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener propiedades: %w", err)
	}
	defer rows.Close() // Importante: siempre cerrar rows

	var propiedades []dominio.Propiedad

	// Iterar sobre todas las filas
	for rows.Next() {
		var propiedad dominio.Propiedad
		var tagsJSON string

		err := rows.Scan(
			&propiedad.ID,
			&propiedad.Slug,
			&propiedad.Titulo,
			&propiedad.Descripcion,
			&propiedad.Precio,
			&propiedad.Provincia,
			&propiedad.Ciudad,
			&propiedad.Sector,
			&propiedad.Direccion,
			&propiedad.Tipo,
			&propiedad.Estado,
			&propiedad.Dormitorios,
			&propiedad.Banos,
			&propiedad.AreaM2,
			&propiedad.FechaCreacion,
			&propiedad.FechaActualizacion,
			&propiedad.ImagenPrincipal,
			&tagsJSON,
			&propiedad.Destacada,
			&propiedad.VisitasContador,
			&propiedad.InmobiliariaID,
		)
		if err != nil {
			return nil, fmt.Errorf("error al leer propiedad: %w", err)
		}

		// Convertir tags JSON
		if tagsJSON != "" {
			err = json.Unmarshal([]byte(tagsJSON), &propiedad.Tags)
			if err != nil {
				// Si falla la conversión, continuar sin tags
				propiedad.Tags = []string{}
			}
		}

		propiedades = append(propiedades, propiedad)
	}

	// Verificar si hubo errores durante la iteración
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar propiedades: %w", err)
	}

	return propiedades, nil
}

// Actualizar modifica una propiedad existente
func (r *PropiedadRepositoryPostgres) Actualizar(propiedad *dominio.Propiedad) error {
	// Actualizar la fecha de modificación y regenerar slug
	propiedad.ActualizarFecha()
	propiedad.ActualizarSlug()

	query := `
		UPDATE propiedades SET 
			slug = $2, titulo = $3, descripcion = $4, precio = $5, provincia = $6, 
			ciudad = $7, sector = $8, direccion = $9, tipo = $10, 
			estado = $11, dormitorios = $12, banos = $13, area_m2 = $14, 
			fecha_actualizacion = $15
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		propiedad.ID,
		propiedad.Slug,
		propiedad.Titulo,
		propiedad.Descripcion,
		propiedad.Precio,
		propiedad.Provincia,
		propiedad.Ciudad,
		propiedad.Sector,
		propiedad.Direccion,
		propiedad.Tipo,
		propiedad.Estado,
		propiedad.Dormitorios,
		propiedad.Banos,
		propiedad.AreaM2,
		propiedad.FechaActualizacion,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar propiedad: %w", err)
	}

	// Verificar si la propiedad fue actualizada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("propiedad no encontrada: %s", propiedad.ID)
	}

	log.Printf("Propiedad actualizada exitosamente: %s", propiedad.ID)
	return nil
}

// Eliminar borra una propiedad de la base de datos
func (r *PropiedadRepositoryPostgres) Eliminar(id string) error {
	query := `DELETE FROM propiedades WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar propiedad: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar eliminación: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("propiedad no encontrada: %s", id)
	}

	log.Printf("Propiedad eliminada exitosamente: %s", id)
	return nil
}

// ConectarBaseDatos establece la conexión a PostgreSQL
func ConectarBaseDatos(databaseURL string) (*sql.DB, error) {
	// Abrir conexión a PostgreSQL
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %w", err)
	}

	// Verificar que la conexión funciona
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)                 // Máximo 25 conexiones abiertas
	db.SetMaxIdleConns(25)                 // Máximo 25 conexiones inactivas
	db.SetConnMaxLifetime(5 * time.Minute) // Tiempo de vida de conexión

	log.Println("Conexión a PostgreSQL establecida exitosamente")
	return db, nil
}
