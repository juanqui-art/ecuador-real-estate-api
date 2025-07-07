package servicio

import (
	"fmt"
	"strings"

	"realty-core/internal/dominio"
	"realty-core/internal/repositorio"
)

// PropiedadService maneja la lógica de negocio para propiedades
type PropiedadService struct {
	repo repositorio.PropiedadRepository
}

// NuevoPropiedadService crea una nueva instancia del servicio
func NuevoPropiedadService(repo repositorio.PropiedadRepository) *PropiedadService {
	return &PropiedadService{repo: repo}
}

// CrearPropiedad crea una nueva propiedad con validaciones
func (s *PropiedadService) CrearPropiedad(titulo, descripcion, provincia, ciudad, tipo string, precio float64) (*dominio.Propiedad, error) {
	// Validar datos de entrada
	if err := s.validarDatosCreacion(titulo, provincia, ciudad, tipo, precio); err != nil {
		return nil, err
	}

	// Limpiar y normalizar datos
	titulo = strings.TrimSpace(titulo)
	descripcion = strings.TrimSpace(descripcion)
	provincia = strings.TrimSpace(provincia)
	ciudad = strings.TrimSpace(ciudad)
	tipo = strings.ToLower(strings.TrimSpace(tipo))

	// Crear la propiedad
	propiedad := dominio.NuevaPropiedad(titulo, descripcion, provincia, ciudad, tipo, precio)

	// Validar la propiedad completa
	if !propiedad.EsValida() {
		return nil, fmt.Errorf("datos de propiedad inválidos")
	}

	// Guardar en la base de datos
	if err := s.repo.Crear(propiedad); err != nil {
		return nil, fmt.Errorf("error al crear propiedad: %w", err)
	}

	return propiedad, nil
}

// ObtenerPropiedad busca una propiedad por ID
func (s *PropiedadService) ObtenerPropiedad(id string) (*dominio.Propiedad, error) {
	if id == "" {
		return nil, fmt.Errorf("ID de propiedad requerido")
	}

	propiedad, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener propiedad: %w", err)
	}

	return propiedad, nil
}

// ObtenerPropiedadPorSlug busca una propiedad por slug SEO
func (s *PropiedadService) ObtenerPropiedadPorSlug(slug string) (*dominio.Propiedad, error) {
	if slug == "" {
		return nil, fmt.Errorf("slug de propiedad requerido")
	}

	// Validar que el slug tenga formato correcto
	if !dominio.EsSlugValido(slug) {
		return nil, fmt.Errorf("formato de slug inválido: %s", slug)
	}

	propiedad, err := s.repo.ObtenerPorSlug(slug)
	if err != nil {
		return nil, fmt.Errorf("error al obtener propiedad por slug: %w", err)
	}

	return propiedad, nil
}

// ListarPropiedades obtiene todas las propiedades
func (s *PropiedadService) ListarPropiedades() ([]dominio.Propiedad, error) {
	propiedades, err := s.repo.ObtenerTodas()
	if err != nil {
		return nil, fmt.Errorf("error al listar propiedades: %w", err)
	}

	return propiedades, nil
}

// ActualizarPropiedad modifica una propiedad existente
func (s *PropiedadService) ActualizarPropiedad(id string, titulo, descripcion, provincia, ciudad, tipo string, precio float64) (*dominio.Propiedad, error) {
	// Validar que la propiedad existe
	propiedad, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, fmt.Errorf("propiedad no encontrada: %w", err)
	}

	// Validar nuevos datos
	if err := s.validarDatosCreacion(titulo, provincia, ciudad, tipo, precio); err != nil {
		return nil, err
	}

	// Actualizar campos
	propiedad.Titulo = strings.TrimSpace(titulo)
	propiedad.Descripcion = strings.TrimSpace(descripcion)
	propiedad.Provincia = strings.TrimSpace(provincia)
	propiedad.Ciudad = strings.TrimSpace(ciudad)
	propiedad.Tipo = strings.ToLower(strings.TrimSpace(tipo))
	propiedad.Precio = precio

	// Validar la propiedad actualizada
	if !propiedad.EsValida() {
		return nil, fmt.Errorf("datos de propiedad actualizados son inválidos")
	}

	// Guardar cambios
	if err := s.repo.Actualizar(propiedad); err != nil {
		return nil, fmt.Errorf("error al actualizar propiedad: %w", err)
	}

	return propiedad, nil
}

// EliminarPropiedad elimina una propiedad por ID
func (s *PropiedadService) EliminarPropiedad(id string) error {
	if id == "" {
		return fmt.Errorf("ID de propiedad requerido")
	}

	// Verificar que la propiedad existe antes de eliminar
	_, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return fmt.Errorf("propiedad no encontrada: %w", err)
	}

	// Eliminar la propiedad
	if err := s.repo.Eliminar(id); err != nil {
		return fmt.Errorf("error al eliminar propiedad: %w", err)
	}

	return nil
}

// FiltrarPorProvincia filtra propiedades por provincia
func (s *PropiedadService) FiltrarPorProvincia(provincia string) ([]dominio.Propiedad, error) {
	if provincia == "" {
		return nil, fmt.Errorf("provincia requerida")
	}

	// Validar que la provincia es válida
	if !dominio.EsProvinciaValida(provincia) {
		return nil, fmt.Errorf("provincia no válida: %s", provincia)
	}

	// Obtener todas las propiedades y filtrar
	// Nota: En una implementación real, esto debería hacerse en el repositorio
	propiedades, err := s.repo.ObtenerTodas()
	if err != nil {
		return nil, fmt.Errorf("error al obtener propiedades: %w", err)
	}

	var propiedadesFiltradas []dominio.Propiedad
	for _, propiedad := range propiedades {
		if propiedad.Provincia == provincia {
			propiedadesFiltradas = append(propiedadesFiltradas, propiedad)
		}
	}

	return propiedadesFiltradas, nil
}

// FiltrarPorRangoPrecio filtra propiedades por rango de precio
func (s *PropiedadService) FiltrarPorRangoPrecio(precioMin, precioMax float64) ([]dominio.Propiedad, error) {
	if precioMin < 0 || precioMax < 0 {
		return nil, fmt.Errorf("los precios deben ser positivos")
	}

	if precioMin > precioMax {
		return nil, fmt.Errorf("precio mínimo no puede ser mayor que precio máximo")
	}

	// Obtener todas las propiedades y filtrar
	propiedades, err := s.repo.ObtenerTodas()
	if err != nil {
		return nil, fmt.Errorf("error al obtener propiedades: %w", err)
	}

	var propiedadesFiltradas []dominio.Propiedad
	for _, propiedad := range propiedades {
		if propiedad.Precio >= precioMin && propiedad.Precio <= precioMax {
			propiedadesFiltradas = append(propiedadesFiltradas, propiedad)
		}
	}

	return propiedadesFiltradas, nil
}

// validarDatosCreacion valida los datos básicos para crear/actualizar una propiedad
func (s *PropiedadService) validarDatosCreacion(titulo, provincia, ciudad, tipo string, precio float64) error {
	// Validar campos obligatorios
	if strings.TrimSpace(titulo) == "" {
		return fmt.Errorf("título es requerido")
	}

	if len(strings.TrimSpace(titulo)) < 10 {
		return fmt.Errorf("título debe tener al menos 10 caracteres")
	}

	if len(strings.TrimSpace(titulo)) > 255 {
		return fmt.Errorf("título no puede exceder 255 caracteres")
	}

	if strings.TrimSpace(provincia) == "" {
		return fmt.Errorf("provincia es requerida")
	}

	if strings.TrimSpace(ciudad) == "" {
		return fmt.Errorf("ciudad es requerida")
	}

	if strings.TrimSpace(tipo) == "" {
		return fmt.Errorf("tipo es requerido")
	}

	if precio <= 0 {
		return fmt.Errorf("precio debe ser mayor a 0")
	}

	// Validar provincia ecuatoriana
	if !dominio.EsProvinciaValida(provincia) {
		return fmt.Errorf("provincia no válida: %s", provincia)
	}

	// Validar tipo de propiedad
	tiposValidos := []string{dominio.TipoCasa, dominio.TipoDepartamento, dominio.TipoTerreno, dominio.TipoComercial}
	tipoLower := strings.ToLower(strings.TrimSpace(tipo))
	
	tipoValido := false
	for _, tipoPermitido := range tiposValidos {
		if tipoLower == tipoPermitido {
			tipoValido = true
			break
		}
	}

	if !tipoValido {
		return fmt.Errorf("tipo de propiedad no válido: %s. Tipos permitidos: %v", tipo, tiposValidos)
	}

	return nil
}

// ObtenerEstadisticas devuelve estadísticas básicas de las propiedades
func (s *PropiedadService) ObtenerEstadisticas() (map[string]interface{}, error) {
	propiedades, err := s.repo.ObtenerTodas()
	if err != nil {
		return nil, fmt.Errorf("error al obtener propiedades: %w", err)
	}

	stats := make(map[string]interface{})
	stats["total_propiedades"] = len(propiedades)

	// Contar por tipo
	tipoCount := make(map[string]int)
	// Contar por estado
	estadoCount := make(map[string]int)
	// Calcular precio promedio
	var precioTotal float64

	for _, propiedad := range propiedades {
		tipoCount[propiedad.Tipo]++
		estadoCount[propiedad.Estado]++
		precioTotal += propiedad.Precio
	}

	stats["por_tipo"] = tipoCount
	stats["por_estado"] = estadoCount
	
	if len(propiedades) > 0 {
		stats["precio_promedio"] = precioTotal / float64(len(propiedades))
	} else {
		stats["precio_promedio"] = 0
	}

	return stats, nil
}