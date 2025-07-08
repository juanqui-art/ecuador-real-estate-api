package servicio

import (
	"fmt"
	"log"

	"realty-core/internal/dominio"
	"realty-core/internal/repositorio"
)

// InmobiliariaService maneja la lógica de negocio para inmobiliarias
type InmobiliariaService struct {
	repo repositorio.InmobiliariaRepository
}

// NuevoInmobiliariaService crea una nueva instancia del servicio
func NuevoInmobiliariaService(repo repositorio.InmobiliariaRepository) *InmobiliariaService {
	return &InmobiliariaService{repo: repo}
}

// Crear crea una nueva inmobiliaria con validaciones
func (s *InmobiliariaService) Crear(nombre, ruc, direccion, telefono, email string) (*dominio.Inmobiliaria, error) {
	// Validar que no existe otra inmobiliaria con el mismo RUC
	existente, err := s.repo.ObtenerPorRUC(ruc)
	if err == nil && existente != nil {
		return nil, fmt.Errorf("ya existe una inmobiliaria con RUC: %s", ruc)
	}

	// Validar que no existe otra inmobiliaria con el mismo email
	// Nota: Necesitaríamos agregar este método al repositorio
	// Por ahora asumimos que el constraint de DB manejará esto

	// Crear nueva inmobiliaria
	inmobiliaria := dominio.NuevaInmobiliaria(nombre, ruc, direccion, telefono, email)

	// Validar datos
	if err := inmobiliaria.Validar(); err != nil {
		return nil, fmt.Errorf("datos de inmobiliaria inválidos: %w", err)
	}

	// Guardar en repositorio
	if err := s.repo.Crear(inmobiliaria); err != nil {
		return nil, fmt.Errorf("error al crear inmobiliaria: %w", err)
	}

	log.Printf("Inmobiliaria creada exitosamente: %s - %s", inmobiliaria.ID, inmobiliaria.Nombre)
	return inmobiliaria, nil
}

// ObtenerPorID obtiene una inmobiliaria por su ID
func (s *InmobiliariaService) ObtenerPorID(id string) (*dominio.Inmobiliaria, error) {
	if id == "" {
		return nil, fmt.Errorf("ID de inmobiliaria requerido")
	}

	inmobiliaria, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener inmobiliaria: %w", err)
	}

	return inmobiliaria, nil
}

// ObtenerPorRUC obtiene una inmobiliaria por su RUC
func (s *InmobiliariaService) ObtenerPorRUC(ruc string) (*dominio.Inmobiliaria, error) {
	if ruc == "" {
		return nil, fmt.Errorf("RUC requerido")
	}

	inmobiliaria, err := s.repo.ObtenerPorRUC(ruc)
	if err != nil {
		return nil, fmt.Errorf("error al obtener inmobiliaria por RUC: %w", err)
	}

	return inmobiliaria, nil
}

// ObtenerTodas obtiene todas las inmobiliarias
func (s *InmobiliariaService) ObtenerTodas() ([]dominio.Inmobiliaria, error) {
	inmobiliarias, err := s.repo.ObtenerTodas()
	if err != nil {
		return nil, fmt.Errorf("error al obtener inmobiliarias: %w", err)
	}

	return inmobiliarias, nil
}

// ObtenerActivas obtiene solo las inmobiliarias activas
func (s *InmobiliariaService) ObtenerActivas() ([]dominio.Inmobiliaria, error) {
	inmobiliarias, err := s.repo.ObtenerActivas()
	if err != nil {
		return nil, fmt.Errorf("error al obtener inmobiliarias activas: %w", err)
	}

	return inmobiliarias, nil
}

// Actualizar actualiza una inmobiliaria existente
func (s *InmobiliariaService) Actualizar(id string, nombre, direccion, telefono, email, sitioWeb, descripcion, logoURL string, activa bool) (*dominio.Inmobiliaria, error) {
	// Obtener inmobiliaria existente
	inmobiliaria, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("inmobiliaria no encontrada: %w", err)
	}

	// Actualizar campos
	inmobiliaria.Nombre = nombre
	inmobiliaria.Direccion = direccion
	inmobiliaria.Telefono = telefono
	inmobiliaria.Email = email
	inmobiliaria.SitioWeb = sitioWeb
	inmobiliaria.Descripcion = descripcion
	inmobiliaria.LogoURL = logoURL
	inmobiliaria.Activa = activa

	// Validar datos actualizados
	if err := inmobiliaria.Validar(); err != nil {
		return nil, fmt.Errorf("datos actualizados inválidos: %w", err)
	}

	// Guardar cambios
	if err := s.repo.Actualizar(inmobiliaria); err != nil {
		return nil, fmt.Errorf("error al actualizar inmobiliaria: %w", err)
	}

	log.Printf("Inmobiliaria actualizada exitosamente: %s", inmobiliaria.ID)
	return inmobiliaria, nil
}

// ActualizarPerfil actualiza solo información de perfil (sin cambiar estado activo)
func (s *InmobiliariaService) ActualizarPerfil(id string, nombre, direccion, telefono, email, sitioWeb, descripcion, logoURL string) (*dominio.Inmobiliaria, error) {
	// Obtener inmobiliaria existente
	inmobiliaria, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("inmobiliaria no encontrada: %w", err)
	}

	// Mantener el estado actual y solo actualizar perfil
	return s.Actualizar(id, nombre, direccion, telefono, email, sitioWeb, descripcion, logoURL, inmobiliaria.Activa)
}

// Desactivar desactiva una inmobiliaria (soft delete)
func (s *InmobiliariaService) Desactivar(id string) error {
	if id == "" {
		return fmt.Errorf("ID de inmobiliaria requerido")
	}

	// Verificar que la inmobiliaria existe
	_, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return fmt.Errorf("inmobiliaria no encontrada: %w", err)
	}

	// Desactivar
	if err := s.repo.Desactivar(id); err != nil {
		return fmt.Errorf("error al desactivar inmobiliaria: %w", err)
	}

	log.Printf("Inmobiliaria desactivada exitosamente: %s", id)
	return nil
}

// Reactivar reactiva una inmobiliaria
func (s *InmobiliariaService) Reactivar(id string) (*dominio.Inmobiliaria, error) {
	// Obtener inmobiliaria
	inmobiliaria, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("inmobiliaria no encontrada: %w", err)
	}

	// Cambiar estado a activa
	inmobiliaria.Activa = true

	// Actualizar
	if err := s.repo.Actualizar(inmobiliaria); err != nil {
		return nil, fmt.Errorf("error al reactivar inmobiliaria: %w", err)
	}

	log.Printf("Inmobiliaria reactivada exitosamente: %s", id)
	return inmobiliaria, nil
}

// BuscarPorNombre busca inmobiliarias por nombre
func (s *InmobiliariaService) BuscarPorNombre(nombre string) ([]dominio.Inmobiliaria, error) {
	if nombre == "" {
		return nil, fmt.Errorf("término de búsqueda requerido")
	}

	inmobiliarias, err := s.repo.BuscarPorNombre(nombre)
	if err != nil {
		return nil, fmt.Errorf("error al buscar inmobiliarias: %w", err)
	}

	return inmobiliarias, nil
}

// ValidarRUC valida un RUC ecuatoriano usando las reglas del dominio
func (s *InmobiliariaService) ValidarRUC(ruc string) error {
	// Crear inmobiliaria temporal solo para validación
	temp := &dominio.Inmobiliaria{RUC: ruc}
	return temp.ValidarRUC()
}

// ValidarEmail valida un email usando las reglas del dominio
func (s *InmobiliariaService) ValidarEmail(email string) error {
	// Crear inmobiliaria temporal solo para validación
	temp := &dominio.Inmobiliaria{Email: email}
	return temp.ValidarEmail()
}

// ValidarTelefono valida un teléfono ecuatoriano usando las reglas del dominio
func (s *InmobiliariaService) ValidarTelefono(telefono string) error {
	// Crear inmobiliaria temporal solo para validación
	temp := &dominio.Inmobiliaria{Telefono: telefono}
	return temp.ValidarTelefono()
}

// ObtenerEstadisticas obtiene estadísticas generales de inmobiliarias
// Nota: Este método requeriría funciones adicionales en el repositorio
func (s *InmobiliariaService) ObtenerEstadisticas() (map[string]interface{}, error) {
	todas, err := s.repo.ObtenerTodas()
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas: %w", err)
	}

	activas, err := s.repo.ObtenerActivas()
	if err != nil {
		return nil, fmt.Errorf("error al obtener inmobiliarias activas: %w", err)
	}

	estadisticas := map[string]interface{}{
		"total":              len(todas),
		"activas":            len(activas),
		"inactivas":          len(todas) - len(activas),
		"porcentaje_activas": float64(len(activas)) / float64(len(todas)) * 100,
	}

	return estadisticas, nil
}
