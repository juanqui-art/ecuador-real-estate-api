package servicio

import (
	"fmt"
	"log"
	"time"

	"realty-core/internal/dominio"
	"realty-core/internal/repositorio"
)

// UsuarioService maneja la lógica de negocio para usuarios
type UsuarioService struct {
	repo             repositorio.UsuarioRepository
	inmobiliariaRepo repositorio.InmobiliariaRepository
}

// NuevoUsuarioService crea una nueva instancia del servicio
func NuevoUsuarioService(repo repositorio.UsuarioRepository, inmobiliariaRepo repositorio.InmobiliariaRepository) *UsuarioService {
	return &UsuarioService{
		repo:             repo,
		inmobiliariaRepo: inmobiliariaRepo,
	}
}

// Crear crea un nuevo usuario con validaciones
func (s *UsuarioService) Crear(nombre, apellido, email, telefono, cedula string, tipoUsuario string) (*dominio.Usuario, error) {
	// Validar que no existe otro usuario con el mismo email
	existente, err := s.repo.ObtenerPorEmail(email)
	if err == nil && existente != nil {
		return nil, fmt.Errorf("ya existe un usuario con email: %s", email)
	}

	// Validar que no existe otro usuario con la misma cédula
	existente, err = s.repo.ObtenerPorCedula(cedula)
	if err == nil && existente != nil {
		return nil, fmt.Errorf("ya existe un usuario con cédula: %s", cedula)
	}

	// Crear nuevo usuario
	usuario := dominio.NuevoUsuario(nombre, apellido, email, telefono, cedula, tipoUsuario)

	// Validar datos
	if err := usuario.Validar(); err != nil {
		return nil, fmt.Errorf("datos de usuario inválidos: %w", err)
	}

	// Guardar en repositorio
	if err := s.repo.Crear(usuario); err != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", err)
	}

	log.Printf("Usuario creado exitosamente: %s - %s %s", usuario.ID, usuario.Nombre, usuario.Apellido)
	return usuario, nil
}

// CrearComprador crea un nuevo usuario comprador con presupuesto
func (s *UsuarioService) CrearComprador(nombre, apellido, email, telefono, cedula string, presupuestoMin, presupuestoMax *float64, provinciasInteres, tiposInteres []string) (*dominio.Usuario, error) {
	usuario, err := s.Crear(nombre, apellido, email, telefono, cedula, "comprador")
	if err != nil {
		return nil, err
	}

	// Configurar preferencias de comprador
	usuario.PresupuestoMin = presupuestoMin
	usuario.PresupuestoMax = presupuestoMax
	usuario.ProvinciasInteres = provinciasInteres
	usuario.TiposPropiedadInteres = tiposInteres

	// Validar presupuesto
	if err := usuario.ValidarPresupuesto(); err != nil {
		return nil, fmt.Errorf("presupuesto inválido: %w", err)
	}

	// Actualizar con preferencias
	if err := s.repo.Actualizar(usuario); err != nil {
		return nil, fmt.Errorf("error al configurar preferencias de comprador: %w", err)
	}

	return usuario, nil
}

// CrearAgente crea un nuevo agente asociado a una inmobiliaria
func (s *UsuarioService) CrearAgente(nombre, apellido, email, telefono, cedula, inmobiliariaID string) (*dominio.Usuario, error) {
	// Validar que la inmobiliaria existe y está activa
	inmobiliaria, err := s.inmobiliariaRepo.ObtenerPorID(inmobiliariaID)
	if err != nil {
		return nil, fmt.Errorf("inmobiliaria no encontrada: %w", err)
	}
	if !inmobiliaria.Activa {
		return nil, fmt.Errorf("no se puede asignar agentes a inmobiliarias inactivas")
	}

	usuario, err := s.Crear(nombre, apellido, email, telefono, cedula, "agente")
	if err != nil {
		return nil, err
	}

	// Asignar inmobiliaria
	usuario.InmobiliariaID = &inmobiliariaID

	// Actualizar con inmobiliaria
	if err := s.repo.Actualizar(usuario); err != nil {
		return nil, fmt.Errorf("error al asignar inmobiliaria al agente: %w", err)
	}

	log.Printf("Agente creado y asignado a inmobiliaria %s: %s", inmobiliariaID, usuario.ID)
	return usuario, nil
}

// ObtenerPorID obtiene un usuario por su ID
func (s *UsuarioService) ObtenerPorID(id string) (*dominio.Usuario, error) {
	if id == "" {
		return nil, fmt.Errorf("ID de usuario requerido")
	}

	usuario, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	return usuario, nil
}

// ObtenerPorEmail obtiene un usuario por su email
func (s *UsuarioService) ObtenerPorEmail(email string) (*dominio.Usuario, error) {
	if email == "" {
		return nil, fmt.Errorf("email requerido")
	}

	usuario, err := s.repo.ObtenerPorEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario por email: %w", err)
	}

	return usuario, nil
}

// ObtenerPorCedula obtiene un usuario por su cédula
func (s *UsuarioService) ObtenerPorCedula(cedula string) (*dominio.Usuario, error) {
	if cedula == "" {
		return nil, fmt.Errorf("cédula requerida")
	}

	usuario, err := s.repo.ObtenerPorCedula(cedula)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario por cédula: %w", err)
	}

	return usuario, nil
}

// ObtenerTodos obtiene todos los usuarios
func (s *UsuarioService) ObtenerTodos() ([]dominio.Usuario, error) {
	usuarios, err := s.repo.ObtenerTodos()
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}

	return usuarios, nil
}

// ObtenerCompradores obtiene todos los usuarios compradores
func (s *UsuarioService) ObtenerCompradores() ([]dominio.Usuario, error) {
	usuarios, err := s.repo.ObtenerPorTipo("comprador")
	if err != nil {
		return nil, fmt.Errorf("error al obtener compradores: %w", err)
	}

	return usuarios, nil
}

// ObtenerVendedores obtiene todos los usuarios vendedores
func (s *UsuarioService) ObtenerVendedores() ([]dominio.Usuario, error) {
	usuarios, err := s.repo.ObtenerPorTipo("vendedor")
	if err != nil {
		return nil, fmt.Errorf("error al obtener vendedores: %w", err)
	}

	return usuarios, nil
}

// ObtenerAgentes obtiene todos los usuarios agentes
func (s *UsuarioService) ObtenerAgentes() ([]dominio.Usuario, error) {
	usuarios, err := s.repo.ObtenerPorTipo("agente")
	if err != nil {
		return nil, fmt.Errorf("error al obtener agentes: %w", err)
	}

	return usuarios, nil
}

// ObtenerAgentesPorInmobiliaria obtiene agentes de una inmobiliaria específica
func (s *UsuarioService) ObtenerAgentesPorInmobiliaria(inmobiliariaID string) ([]dominio.Usuario, error) {
	if inmobiliariaID == "" {
		return nil, fmt.Errorf("ID de inmobiliaria requerido")
	}

	usuarios, err := s.repo.ObtenerPorInmobiliaria(inmobiliariaID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener agentes de inmobiliaria: %w", err)
	}

	return usuarios, nil
}

// Actualizar actualiza un usuario existente
func (s *UsuarioService) Actualizar(id string, nombre, apellido, email, telefono string, fechaNacimiento *time.Time, avatarURL, bio string) (*dominio.Usuario, error) {
	// Obtener usuario existente
	usuario, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Actualizar campos
	usuario.Nombre = nombre
	usuario.Apellido = apellido
	usuario.Email = email
	usuario.Telefono = telefono
	usuario.FechaNacimiento = fechaNacimiento
	usuario.AvatarURL = avatarURL
	usuario.Bio = bio

	// Validar datos actualizados
	if err := usuario.Validar(); err != nil {
		return nil, fmt.Errorf("datos actualizados inválidos: %w", err)
	}

	// Guardar cambios
	if err := s.repo.Actualizar(usuario); err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	log.Printf("Usuario actualizado exitosamente: %s", usuario.ID)
	return usuario, nil
}

// ActualizarPreferenciasComprador actualiza las preferencias de un comprador
func (s *UsuarioService) ActualizarPreferenciasComprador(id string, presupuestoMin, presupuestoMax *float64, provinciasInteres, tiposInteres []string) (*dominio.Usuario, error) {
	// Obtener usuario
	usuario, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Verificar que es comprador
	if usuario.TipoUsuario != "comprador" {
		return nil, fmt.Errorf("solo los compradores pueden tener preferencias de búsqueda")
	}

	// Actualizar preferencias
	usuario.PresupuestoMin = presupuestoMin
	usuario.PresupuestoMax = presupuestoMax
	usuario.ProvinciasInteres = provinciasInteres
	usuario.TiposPropiedadInteres = tiposInteres

	// Validar presupuesto
	if err := usuario.ValidarPresupuesto(); err != nil {
		return nil, fmt.Errorf("presupuesto inválido: %w", err)
	}

	// Guardar cambios
	if err := s.repo.Actualizar(usuario); err != nil {
		return nil, fmt.Errorf("error al actualizar preferencias: %w", err)
	}

	log.Printf("Preferencias de comprador actualizadas: %s", usuario.ID)
	return usuario, nil
}

// CambiarInmobiliaria cambia la inmobiliaria de un agente
func (s *UsuarioService) CambiarInmobiliaria(agenteID, nuevaInmobiliariaID string) (*dominio.Usuario, error) {
	// Obtener agente
	usuario, err := s.repo.ObtenerPorID(agenteID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Verificar que es agente
	if usuario.TipoUsuario != "agente" {
		return nil, fmt.Errorf("solo los agentes pueden cambiar de inmobiliaria")
	}

	// Validar nueva inmobiliaria
	inmobiliaria, err := s.inmobiliariaRepo.ObtenerPorID(nuevaInmobiliariaID)
	if err != nil {
		return nil, fmt.Errorf("inmobiliaria no encontrada: %w", err)
	}
	if !inmobiliaria.Activa {
		return nil, fmt.Errorf("no se puede asignar agentes a inmobiliarias inactivas")
	}

	// Cambiar inmobiliaria
	usuario.InmobiliariaID = &nuevaInmobiliariaID

	// Guardar cambios
	if err := s.repo.Actualizar(usuario); err != nil {
		return nil, fmt.Errorf("error al cambiar inmobiliaria: %w", err)
	}

	log.Printf("Agente %s cambiado a inmobiliaria %s", agenteID, nuevaInmobiliariaID)
	return usuario, nil
}

// Desactivar desactiva un usuario (soft delete)
func (s *UsuarioService) Desactivar(id string) error {
	if id == "" {
		return fmt.Errorf("ID de usuario requerido")
	}

	// Verificar que el usuario existe
	_, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Desactivar
	if err := s.repo.Desactivar(id); err != nil {
		return fmt.Errorf("error al desactivar usuario: %w", err)
	}

	log.Printf("Usuario desactivado exitosamente: %s", id)
	return nil
}

// Reactivar reactiva un usuario
func (s *UsuarioService) Reactivar(id string) (*dominio.Usuario, error) {
	// Obtener usuario
	usuario, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Cambiar estado a activo
	usuario.Activo = true

	// Actualizar
	if err := s.repo.Actualizar(usuario); err != nil {
		return nil, fmt.Errorf("error al reactivar usuario: %w", err)
	}

	log.Printf("Usuario reactivado exitosamente: %s", id)
	return usuario, nil
}

// BuscarPorNombre busca usuarios por nombre
func (s *UsuarioService) BuscarPorNombre(nombre string) ([]dominio.Usuario, error) {
	if nombre == "" {
		return nil, fmt.Errorf("término de búsqueda requerido")
	}

	usuarios, err := s.repo.BuscarPorNombre(nombre)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuarios: %w", err)
	}

	return usuarios, nil
}

// BuscarCompradoresParaPropiedad busca compradores que puedan pagar una propiedad
func (s *UsuarioService) BuscarCompradoresParaPropiedad(precioPropiedad float64) ([]dominio.Usuario, error) {
	if precioPropiedad <= 0 {
		return nil, fmt.Errorf("precio de propiedad debe ser mayor a 0")
	}

	usuarios, err := s.repo.BuscarCompradores(precioPropiedad)
	if err != nil {
		return nil, fmt.Errorf("error al buscar compradores: %w", err)
	}

	return usuarios, nil
}

// ValidarCedula valida una cédula ecuatoriana
func (s *UsuarioService) ValidarCedula(cedula string) error {
	// Crear usuario temporal solo para validación
	temp := &dominio.Usuario{Cedula: cedula}
	return temp.ValidarCedula()
}

// ValidarEmail valida un email
func (s *UsuarioService) ValidarEmail(email string) error {
	// Crear usuario temporal solo para validación
	temp := &dominio.Usuario{Email: email}
	return temp.ValidarEmail()
}

// ValidarTelefono valida un teléfono ecuatoriano
func (s *UsuarioService) ValidarTelefono(telefono string) error {
	// Crear usuario temporal solo para validación
	temp := &dominio.Usuario{Telefono: telefono}
	return temp.ValidarTelefono()
}

// ObtenerEstadisticas obtiene estadísticas generales de usuarios
func (s *UsuarioService) ObtenerEstadisticas() (map[string]interface{}, error) {
	todos, err := s.repo.ObtenerTodos()
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas: %w", err)
	}

	compradores, err := s.repo.ObtenerPorTipo("comprador")
	if err != nil {
		return nil, fmt.Errorf("error al obtener compradores: %w", err)
	}

	vendedores, err := s.repo.ObtenerPorTipo("vendedor")
	if err != nil {
		return nil, fmt.Errorf("error al obtener vendedores: %w", err)
	}

	agentes, err := s.repo.ObtenerPorTipo("agente")
	if err != nil {
		return nil, fmt.Errorf("error al obtener agentes: %w", err)
	}

	// Contar activos e inactivos
	activos := 0
	for _, usuario := range todos {
		if usuario.Activo {
			activos++
		}
	}

	estadisticas := map[string]interface{}{
		"total":              len(todos),
		"activos":            activos,
		"inactivos":          len(todos) - activos,
		"compradores":        len(compradores),
		"vendedores":         len(vendedores),
		"agentes":            len(agentes),
		"porcentaje_activos": float64(activos) / float64(len(todos)) * 100,
	}

	return estadisticas, nil
}
