package repositorio

import (
	"database/sql"
	"fmt"
	"log"

	"realty-core/internal/dominio"
)

// InmobiliariaRepository define las operaciones de acceso a datos para inmobiliarias
type InmobiliariaRepository interface {
	Crear(inmobiliaria *dominio.Inmobiliaria) error
	ObtenerPorID(id string) (*dominio.Inmobiliaria, error)
	ObtenerPorRUC(ruc string) (*dominio.Inmobiliaria, error)
	ObtenerTodas() ([]dominio.Inmobiliaria, error)
	ObtenerActivas() ([]dominio.Inmobiliaria, error)
	Actualizar(inmobiliaria *dominio.Inmobiliaria) error
	Desactivar(id string) error
	BuscarPorNombre(nombre string) ([]dominio.Inmobiliaria, error)
}

// InmobiliariaRepositoryPostgres implementa InmobiliariaRepository usando PostgreSQL
type InmobiliariaRepositoryPostgres struct {
	db *sql.DB
}

// NuevoInmobiliariaRepositoryPostgres crea una nueva instancia del repositorio
func NuevoInmobiliariaRepositoryPostgres(db *sql.DB) *InmobiliariaRepositoryPostgres {
	return &InmobiliariaRepositoryPostgres{db: db}
}

// Crear inserta una nueva inmobiliaria en la base de datos
func (r *InmobiliariaRepositoryPostgres) Crear(inmobiliaria *dominio.Inmobiliaria) error {
	query := `
		INSERT INTO inmobiliarias (
			id, nombre, ruc, direccion, telefono, email, 
			sitio_web, descripcion, logo_url, activa, 
			fecha_creacion, fecha_actualizacion
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.Exec(
		query,
		inmobiliaria.ID,
		inmobiliaria.Nombre,
		inmobiliaria.RUC,
		inmobiliaria.Direccion,
		inmobiliaria.Telefono,
		inmobiliaria.Email,
		inmobiliaria.SitioWeb,
		inmobiliaria.Descripcion,
		inmobiliaria.LogoURL,
		inmobiliaria.Activa,
		inmobiliaria.FechaCreacion,
		inmobiliaria.FechaActualizacion,
	)

	if err != nil {
		return fmt.Errorf("error al crear inmobiliaria: %w", err)
	}

	log.Printf("Inmobiliaria creada exitosamente: %s", inmobiliaria.ID)
	return nil
}

// ObtenerPorID busca una inmobiliaria por su ID
func (r *InmobiliariaRepositoryPostgres) ObtenerPorID(id string) (*dominio.Inmobiliaria, error) {
	query := `
		SELECT id, nombre, ruc, direccion, telefono, email, 
			   sitio_web, descripcion, logo_url, activa, 
			   fecha_creacion, fecha_actualizacion
		FROM inmobiliarias 
		WHERE id = $1
	`

	var inmobiliaria dominio.Inmobiliaria

	err := r.db.QueryRow(query, id).Scan(
		&inmobiliaria.ID,
		&inmobiliaria.Nombre,
		&inmobiliaria.RUC,
		&inmobiliaria.Direccion,
		&inmobiliaria.Telefono,
		&inmobiliaria.Email,
		&inmobiliaria.SitioWeb,
		&inmobiliaria.Descripcion,
		&inmobiliaria.LogoURL,
		&inmobiliaria.Activa,
		&inmobiliaria.FechaCreacion,
		&inmobiliaria.FechaActualizacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inmobiliaria no encontrada: %s", id)
		}
		return nil, fmt.Errorf("error al obtener inmobiliaria: %w", err)
	}

	return &inmobiliaria, nil
}

// ObtenerPorRUC busca una inmobiliaria por su RUC
func (r *InmobiliariaRepositoryPostgres) ObtenerPorRUC(ruc string) (*dominio.Inmobiliaria, error) {
	query := `
		SELECT id, nombre, ruc, direccion, telefono, email, 
			   sitio_web, descripcion, logo_url, activa, 
			   fecha_creacion, fecha_actualizacion
		FROM inmobiliarias 
		WHERE ruc = $1
	`

	var inmobiliaria dominio.Inmobiliaria

	err := r.db.QueryRow(query, ruc).Scan(
		&inmobiliaria.ID,
		&inmobiliaria.Nombre,
		&inmobiliaria.RUC,
		&inmobiliaria.Direccion,
		&inmobiliaria.Telefono,
		&inmobiliaria.Email,
		&inmobiliaria.SitioWeb,
		&inmobiliaria.Descripcion,
		&inmobiliaria.LogoURL,
		&inmobiliaria.Activa,
		&inmobiliaria.FechaCreacion,
		&inmobiliaria.FechaActualizacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inmobiliaria no encontrada con RUC: %s", ruc)
		}
		return nil, fmt.Errorf("error al obtener inmobiliaria por RUC: %w", err)
	}

	return &inmobiliaria, nil
}

// ObtenerTodas devuelve todas las inmobiliarias
func (r *InmobiliariaRepositoryPostgres) ObtenerTodas() ([]dominio.Inmobiliaria, error) {
	query := `
		SELECT id, nombre, ruc, direccion, telefono, email, 
			   sitio_web, descripcion, logo_url, activa, 
			   fecha_creacion, fecha_actualizacion
		FROM inmobiliarias 
		ORDER BY nombre ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener inmobiliarias: %w", err)
	}
	defer rows.Close()

	var inmobiliarias []dominio.Inmobiliaria

	for rows.Next() {
		var inmobiliaria dominio.Inmobiliaria
		err := rows.Scan(
			&inmobiliaria.ID,
			&inmobiliaria.Nombre,
			&inmobiliaria.RUC,
			&inmobiliaria.Direccion,
			&inmobiliaria.Telefono,
			&inmobiliaria.Email,
			&inmobiliaria.SitioWeb,
			&inmobiliaria.Descripcion,
			&inmobiliaria.LogoURL,
			&inmobiliaria.Activa,
			&inmobiliaria.FechaCreacion,
			&inmobiliaria.FechaActualizacion,
		)
		if err != nil {
			return nil, fmt.Errorf("error al leer inmobiliaria: %w", err)
		}
		inmobiliarias = append(inmobiliarias, inmobiliaria)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar inmobiliarias: %w", err)
	}

	return inmobiliarias, nil
}

// ObtenerActivas devuelve solo las inmobiliarias activas
func (r *InmobiliariaRepositoryPostgres) ObtenerActivas() ([]dominio.Inmobiliaria, error) {
	query := `
		SELECT id, nombre, ruc, direccion, telefono, email, 
			   sitio_web, descripcion, logo_url, activa, 
			   fecha_creacion, fecha_actualizacion
		FROM inmobiliarias 
		WHERE activa = TRUE
		ORDER BY nombre ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener inmobiliarias activas: %w", err)
	}
	defer rows.Close()

	var inmobiliarias []dominio.Inmobiliaria

	for rows.Next() {
		var inmobiliaria dominio.Inmobiliaria
		err := rows.Scan(
			&inmobiliaria.ID,
			&inmobiliaria.Nombre,
			&inmobiliaria.RUC,
			&inmobiliaria.Direccion,
			&inmobiliaria.Telefono,
			&inmobiliaria.Email,
			&inmobiliaria.SitioWeb,
			&inmobiliaria.Descripcion,
			&inmobiliaria.LogoURL,
			&inmobiliaria.Activa,
			&inmobiliaria.FechaCreacion,
			&inmobiliaria.FechaActualizacion,
		)
		if err != nil {
			return nil, fmt.Errorf("error al leer inmobiliaria activa: %w", err)
		}
		inmobiliarias = append(inmobiliarias, inmobiliaria)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar inmobiliarias activas: %w", err)
	}

	return inmobiliarias, nil
}

// Actualizar modifica una inmobiliaria existente
func (r *InmobiliariaRepositoryPostgres) Actualizar(inmobiliaria *dominio.Inmobiliaria) error {
	// Actualizar fecha de modificación
	inmobiliaria.ActualizarFecha()

	query := `
		UPDATE inmobiliarias SET 
			nombre = $2, ruc = $3, direccion = $4, telefono = $5, 
			email = $6, sitio_web = $7, descripcion = $8, 
			logo_url = $9, activa = $10, fecha_actualizacion = $11
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		inmobiliaria.ID,
		inmobiliaria.Nombre,
		inmobiliaria.RUC,
		inmobiliaria.Direccion,
		inmobiliaria.Telefono,
		inmobiliaria.Email,
		inmobiliaria.SitioWeb,
		inmobiliaria.Descripcion,
		inmobiliaria.LogoURL,
		inmobiliaria.Activa,
		inmobiliaria.FechaActualizacion,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar inmobiliaria: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inmobiliaria no encontrada: %s", inmobiliaria.ID)
	}

	log.Printf("Inmobiliaria actualizada exitosamente: %s", inmobiliaria.ID)
	return nil
}

// Desactivar marca una inmobiliaria como inactiva (soft delete)
func (r *InmobiliariaRepositoryPostgres) Desactivar(id string) error {
	query := `UPDATE inmobiliarias SET activa = FALSE, fecha_actualizacion = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al desactivar inmobiliaria: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar desactivación: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inmobiliaria no encontrada: %s", id)
	}

	log.Printf("Inmobiliaria desactivada exitosamente: %s", id)
	return nil
}

// BuscarPorNombre busca inmobiliarias por nombre usando búsqueda de texto completo
func (r *InmobiliariaRepositoryPostgres) BuscarPorNombre(nombre string) ([]dominio.Inmobiliaria, error) {
	query := `
		SELECT id, nombre, ruc, direccion, telefono, email, 
			   sitio_web, descripcion, logo_url, activa, 
			   fecha_creacion, fecha_actualizacion
		FROM inmobiliarias 
		WHERE activa = TRUE
		  AND to_tsvector('spanish', nombre || ' ' || descripcion) @@ plainto_tsquery('spanish', $1)
		ORDER BY ts_rank(to_tsvector('spanish', nombre || ' ' || descripcion), plainto_tsquery('spanish', $1)) DESC
	`

	rows, err := r.db.Query(query, nombre)
	if err != nil {
		return nil, fmt.Errorf("error al buscar inmobiliarias: %w", err)
	}
	defer rows.Close()

	var inmobiliarias []dominio.Inmobiliaria

	for rows.Next() {
		var inmobiliaria dominio.Inmobiliaria
		err := rows.Scan(
			&inmobiliaria.ID,
			&inmobiliaria.Nombre,
			&inmobiliaria.RUC,
			&inmobiliaria.Direccion,
			&inmobiliaria.Telefono,
			&inmobiliaria.Email,
			&inmobiliaria.SitioWeb,
			&inmobiliaria.Descripcion,
			&inmobiliaria.LogoURL,
			&inmobiliaria.Activa,
			&inmobiliaria.FechaCreacion,
			&inmobiliaria.FechaActualizacion,
		)
		if err != nil {
			return nil, fmt.Errorf("error al leer resultado de búsqueda: %w", err)
		}
		inmobiliarias = append(inmobiliarias, inmobiliaria)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar resultados de búsqueda: %w", err)
	}

	return inmobiliarias, nil
}
