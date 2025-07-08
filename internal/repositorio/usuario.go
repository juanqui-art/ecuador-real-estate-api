package repositorio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"realty-core/internal/dominio"
)

// UsuarioRepository define las operaciones de acceso a datos para usuarios
type UsuarioRepository interface {
	Crear(usuario *dominio.Usuario) error
	ObtenerPorID(id string) (*dominio.Usuario, error)
	ObtenerPorEmail(email string) (*dominio.Usuario, error)
	ObtenerPorCedula(cedula string) (*dominio.Usuario, error)
	ObtenerTodos() ([]dominio.Usuario, error)
	ObtenerPorTipo(tipo string) ([]dominio.Usuario, error)
	ObtenerPorInmobiliaria(inmobiliariaID string) ([]dominio.Usuario, error)
	Actualizar(usuario *dominio.Usuario) error
	Desactivar(id string) error
	BuscarPorNombre(nombre string) ([]dominio.Usuario, error)
	BuscarCompradores(precioPropiedad float64) ([]dominio.Usuario, error)
}

// UsuarioRepositoryPostgres implementa UsuarioRepository usando PostgreSQL
type UsuarioRepositoryPostgres struct {
	db *sql.DB
}

// NuevoUsuarioRepositoryPostgres crea una nueva instancia del repositorio
func NuevoUsuarioRepositoryPostgres(db *sql.DB) *UsuarioRepositoryPostgres {
	return &UsuarioRepositoryPostgres{db: db}
}

// Crear inserta un nuevo usuario en la base de datos
func (r *UsuarioRepositoryPostgres) Crear(usuario *dominio.Usuario) error {
	// Convertir slices a JSON para almacenar en JSONB
	provinciasJSON, err := json.Marshal(usuario.ProvinciasInteres)
	if err != nil {
		return fmt.Errorf("error al convertir provincias a JSON: %w", err)
	}

	tiposJSON, err := json.Marshal(usuario.TiposPropiedadInteres)
	if err != nil {
		return fmt.Errorf("error al convertir tipos de propiedad a JSON: %w", err)
	}

	query := `
		INSERT INTO usuarios (
			id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			tipo_usuario, activo, presupuesto_min, presupuesto_max,
			provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			inmobiliaria_id, fecha_creacion, fecha_actualizacion
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)
	`

	_, err = r.db.Exec(
		query,
		usuario.ID,
		usuario.Nombre,
		usuario.Apellido,
		usuario.Email,
		usuario.Telefono,
		usuario.Cedula,
		usuario.FechaNacimiento,
		usuario.TipoUsuario,
		usuario.Activo,
		usuario.PresupuestoMin,
		usuario.PresupuestoMax,
		string(provinciasJSON),
		string(tiposJSON),
		usuario.AvatarURL,
		usuario.Bio,
		usuario.InmobiliariaID,
		usuario.FechaCreacion,
		usuario.FechaActualizacion,
	)

	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	log.Printf("Usuario creado exitosamente: %s", usuario.ID)
	return nil
}

// ObtenerPorID busca un usuario por su ID
func (r *UsuarioRepositoryPostgres) ObtenerPorID(id string) (*dominio.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			   tipo_usuario, activo, presupuesto_min, presupuesto_max,
			   provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			   inmobiliaria_id, fecha_creacion, fecha_actualizacion
		FROM usuarios 
		WHERE id = $1
	`

	var usuario dominio.Usuario
	var provinciasJSON, tiposJSON string

	err := r.db.QueryRow(query, id).Scan(
		&usuario.ID,
		&usuario.Nombre,
		&usuario.Apellido,
		&usuario.Email,
		&usuario.Telefono,
		&usuario.Cedula,
		&usuario.FechaNacimiento,
		&usuario.TipoUsuario,
		&usuario.Activo,
		&usuario.PresupuestoMin,
		&usuario.PresupuestoMax,
		&provinciasJSON,
		&tiposJSON,
		&usuario.AvatarURL,
		&usuario.Bio,
		&usuario.InmobiliariaID,
		&usuario.FechaCreacion,
		&usuario.FechaActualizacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado: %s", id)
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	// Convertir JSON de vuelta a slices
	if provinciasJSON != "" {
		err = json.Unmarshal([]byte(provinciasJSON), &usuario.ProvinciasInteres)
		if err != nil {
			return nil, fmt.Errorf("error al convertir provincias desde JSON: %w", err)
		}
	}

	if tiposJSON != "" {
		err = json.Unmarshal([]byte(tiposJSON), &usuario.TiposPropiedadInteres)
		if err != nil {
			return nil, fmt.Errorf("error al convertir tipos desde JSON: %w", err)
		}
	}

	return &usuario, nil
}

// ObtenerPorEmail busca un usuario por su email
func (r *UsuarioRepositoryPostgres) ObtenerPorEmail(email string) (*dominio.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			   tipo_usuario, activo, presupuesto_min, presupuesto_max,
			   provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			   inmobiliaria_id, fecha_creacion, fecha_actualizacion
		FROM usuarios 
		WHERE email = $1
	`

	var usuario dominio.Usuario
	var provinciasJSON, tiposJSON string

	err := r.db.QueryRow(query, email).Scan(
		&usuario.ID,
		&usuario.Nombre,
		&usuario.Apellido,
		&usuario.Email,
		&usuario.Telefono,
		&usuario.Cedula,
		&usuario.FechaNacimiento,
		&usuario.TipoUsuario,
		&usuario.Activo,
		&usuario.PresupuestoMin,
		&usuario.PresupuestoMax,
		&provinciasJSON,
		&tiposJSON,
		&usuario.AvatarURL,
		&usuario.Bio,
		&usuario.InmobiliariaID,
		&usuario.FechaCreacion,
		&usuario.FechaActualizacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado con email: %s", email)
		}
		return nil, fmt.Errorf("error al obtener usuario por email: %w", err)
	}

	// Convertir JSON de vuelta a slices
	if provinciasJSON != "" {
		err = json.Unmarshal([]byte(provinciasJSON), &usuario.ProvinciasInteres)
		if err != nil {
			return nil, fmt.Errorf("error al convertir provincias desde JSON: %w", err)
		}
	}

	if tiposJSON != "" {
		err = json.Unmarshal([]byte(tiposJSON), &usuario.TiposPropiedadInteres)
		if err != nil {
			return nil, fmt.Errorf("error al convertir tipos desde JSON: %w", err)
		}
	}

	return &usuario, nil
}

// ObtenerPorCedula busca un usuario por su cédula
func (r *UsuarioRepositoryPostgres) ObtenerPorCedula(cedula string) (*dominio.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			   tipo_usuario, activo, presupuesto_min, presupuesto_max,
			   provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			   inmobiliaria_id, fecha_creacion, fecha_actualizacion
		FROM usuarios 
		WHERE cedula = $1
	`

	var usuario dominio.Usuario
	var provinciasJSON, tiposJSON string

	err := r.db.QueryRow(query, cedula).Scan(
		&usuario.ID,
		&usuario.Nombre,
		&usuario.Apellido,
		&usuario.Email,
		&usuario.Telefono,
		&usuario.Cedula,
		&usuario.FechaNacimiento,
		&usuario.TipoUsuario,
		&usuario.Activo,
		&usuario.PresupuestoMin,
		&usuario.PresupuestoMax,
		&provinciasJSON,
		&tiposJSON,
		&usuario.AvatarURL,
		&usuario.Bio,
		&usuario.InmobiliariaID,
		&usuario.FechaCreacion,
		&usuario.FechaActualizacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado con cédula: %s", cedula)
		}
		return nil, fmt.Errorf("error al obtener usuario por cédula: %w", err)
	}

	// Convertir JSON de vuelta a slices
	if provinciasJSON != "" {
		err = json.Unmarshal([]byte(provinciasJSON), &usuario.ProvinciasInteres)
		if err != nil {
			return nil, fmt.Errorf("error al convertir provincias desde JSON: %w", err)
		}
	}

	if tiposJSON != "" {
		err = json.Unmarshal([]byte(tiposJSON), &usuario.TiposPropiedadInteres)
		if err != nil {
			return nil, fmt.Errorf("error al convertir tipos desde JSON: %w", err)
		}
	}

	return &usuario, nil
}

// ObtenerTodos devuelve todos los usuarios
func (r *UsuarioRepositoryPostgres) ObtenerTodos() ([]dominio.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			   tipo_usuario, activo, presupuesto_min, presupuesto_max,
			   provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			   inmobiliaria_id, fecha_creacion, fecha_actualizacion
		FROM usuarios 
		ORDER BY nombre, apellido ASC
	`

	return r.ejecutarConsultaMultiple(query)
}

// ObtenerPorTipo devuelve usuarios filtrados por tipo
func (r *UsuarioRepositoryPostgres) ObtenerPorTipo(tipo string) ([]dominio.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			   tipo_usuario, activo, presupuesto_min, presupuesto_max,
			   provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			   inmobiliaria_id, fecha_creacion, fecha_actualizacion
		FROM usuarios 
		WHERE tipo_usuario = $1 AND activo = TRUE
		ORDER BY nombre, apellido ASC
	`

	return r.ejecutarConsultaMultipleConParametro(query, tipo)
}

// ObtenerPorInmobiliaria devuelve usuarios (agentes) de una inmobiliaria específica
func (r *UsuarioRepositoryPostgres) ObtenerPorInmobiliaria(inmobiliariaID string) ([]dominio.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			   tipo_usuario, activo, presupuesto_min, presupuesto_max,
			   provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			   inmobiliaria_id, fecha_creacion, fecha_actualizacion
		FROM usuarios 
		WHERE inmobiliaria_id = $1 AND activo = TRUE
		ORDER BY nombre, apellido ASC
	`

	return r.ejecutarConsultaMultipleConParametro(query, inmobiliariaID)
}

// Actualizar modifica un usuario existente
func (r *UsuarioRepositoryPostgres) Actualizar(usuario *dominio.Usuario) error {
	// Actualizar fecha de modificación
	usuario.ActualizarFecha()

	// Convertir slices a JSON
	provinciasJSON, err := json.Marshal(usuario.ProvinciasInteres)
	if err != nil {
		return fmt.Errorf("error al convertir provincias a JSON: %w", err)
	}

	tiposJSON, err := json.Marshal(usuario.TiposPropiedadInteres)
	if err != nil {
		return fmt.Errorf("error al convertir tipos de propiedad a JSON: %w", err)
	}

	query := `
		UPDATE usuarios SET 
			nombre = $2, apellido = $3, email = $4, telefono = $5, 
			cedula = $6, fecha_nacimiento = $7, tipo_usuario = $8, 
			activo = $9, presupuesto_min = $10, presupuesto_max = $11,
			provincias_interes = $12, tipos_propiedad_interes = $13, 
			avatar_url = $14, bio = $15, inmobiliaria_id = $16, 
			fecha_actualizacion = $17
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		usuario.ID,
		usuario.Nombre,
		usuario.Apellido,
		usuario.Email,
		usuario.Telefono,
		usuario.Cedula,
		usuario.FechaNacimiento,
		usuario.TipoUsuario,
		usuario.Activo,
		usuario.PresupuestoMin,
		usuario.PresupuestoMax,
		string(provinciasJSON),
		string(tiposJSON),
		usuario.AvatarURL,
		usuario.Bio,
		usuario.InmobiliariaID,
		usuario.FechaActualizacion,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario no encontrado: %s", usuario.ID)
	}

	log.Printf("Usuario actualizado exitosamente: %s", usuario.ID)
	return nil
}

// Desactivar marca un usuario como inactivo (soft delete)
func (r *UsuarioRepositoryPostgres) Desactivar(id string) error {
	query := `UPDATE usuarios SET activo = FALSE, fecha_actualizacion = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al desactivar usuario: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar desactivación: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario no encontrado: %s", id)
	}

	log.Printf("Usuario desactivado exitosamente: %s", id)
	return nil
}

// BuscarPorNombre busca usuarios por nombre usando búsqueda de texto completo
func (r *UsuarioRepositoryPostgres) BuscarPorNombre(nombre string) ([]dominio.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			   tipo_usuario, activo, presupuesto_min, presupuesto_max,
			   provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			   inmobiliaria_id, fecha_creacion, fecha_actualizacion
		FROM usuarios 
		WHERE to_tsvector('spanish', nombre || ' ' || apellido) @@ plainto_tsquery('spanish', $1)
		ORDER BY ts_rank(to_tsvector('spanish', nombre || ' ' || apellido), plainto_tsquery('spanish', $1)) DESC
	`

	return r.ejecutarConsultaMultipleConParametro(query, nombre)
}

// BuscarCompradores busca compradores que puedan pagar cierto precio
func (r *UsuarioRepositoryPostgres) BuscarCompradores(precioPropiedad float64) ([]dominio.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email, telefono, cedula, fecha_nacimiento,
			   tipo_usuario, activo, presupuesto_min, presupuesto_max,
			   provincias_interes, tipos_propiedad_interes, avatar_url, bio,
			   inmobiliaria_id, fecha_creacion, fecha_actualizacion
		FROM usuarios 
		WHERE tipo_usuario = 'comprador' 
		  AND activo = TRUE
		  AND presupuesto_min IS NOT NULL 
		  AND presupuesto_max IS NOT NULL
		  AND $1 >= presupuesto_min 
		  AND $1 <= presupuesto_max
		ORDER BY presupuesto_max DESC
	`

	return r.ejecutarConsultaMultipleConParametro(query, precioPropiedad)
}

// ejecutarConsultaMultiple - helper para consultas que devuelven múltiples usuarios
func (r *UsuarioRepositoryPostgres) ejecutarConsultaMultiple(query string) ([]dominio.Usuario, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar consulta: %w", err)
	}
	defer rows.Close()

	return r.procesarFilasUsuarios(rows)
}

// ejecutarConsultaMultipleConParametro - helper para consultas con parámetro
func (r *UsuarioRepositoryPostgres) ejecutarConsultaMultipleConParametro(query string, param interface{}) ([]dominio.Usuario, error) {
	rows, err := r.db.Query(query, param)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar consulta con parámetro: %w", err)
	}
	defer rows.Close()

	return r.procesarFilasUsuarios(rows)
}

// procesarFilasUsuarios - helper para procesar filas de usuarios
func (r *UsuarioRepositoryPostgres) procesarFilasUsuarios(rows *sql.Rows) ([]dominio.Usuario, error) {
	var usuarios []dominio.Usuario

	for rows.Next() {
		var usuario dominio.Usuario
		var provinciasJSON, tiposJSON string

		err := rows.Scan(
			&usuario.ID,
			&usuario.Nombre,
			&usuario.Apellido,
			&usuario.Email,
			&usuario.Telefono,
			&usuario.Cedula,
			&usuario.FechaNacimiento,
			&usuario.TipoUsuario,
			&usuario.Activo,
			&usuario.PresupuestoMin,
			&usuario.PresupuestoMax,
			&provinciasJSON,
			&tiposJSON,
			&usuario.AvatarURL,
			&usuario.Bio,
			&usuario.InmobiliariaID,
			&usuario.FechaCreacion,
			&usuario.FechaActualizacion,
		)
		if err != nil {
			return nil, fmt.Errorf("error al leer usuario: %w", err)
		}

		// Convertir JSON de vuelta a slices
		if provinciasJSON != "" {
			err = json.Unmarshal([]byte(provinciasJSON), &usuario.ProvinciasInteres)
			if err != nil {
				// Si falla la conversión, continuar sin provincias
				usuario.ProvinciasInteres = []string{}
			}
		}

		if tiposJSON != "" {
			err = json.Unmarshal([]byte(tiposJSON), &usuario.TiposPropiedadInteres)
			if err != nil {
				// Si falla la conversión, continuar sin tipos
				usuario.TiposPropiedadInteres = []string{}
			}
		}

		usuarios = append(usuarios, usuario)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar usuarios: %w", err)
	}

	return usuarios, nil
}
