package dominio

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Propiedad representa una propiedad inmobiliaria
// En Go, los tipos se definen con 'type NombreTipo struct'
type Propiedad struct {
	// ID es el identificador único de la propiedad
	// uuid.UUID es más robusto que string para IDs
	ID string `json:"id" db:"id"`

	// Slug para SEO - URL amigable generada desde el título
	Slug string `json:"slug" db:"slug"`

	// Información básica de la propiedad
	Titulo      string  `json:"titulo" db:"titulo"`
	Descripcion string  `json:"descripcion" db:"descripcion"`
	Precio      float64 `json:"precio" db:"precio"`

	// Ubicación específica para Ecuador
	Provincia string `json:"provincia" db:"provincia"`
	Ciudad    string `json:"ciudad" db:"ciudad"`
	Sector    string `json:"sector" db:"sector"`     // Opcional
	Direccion string `json:"direccion" db:"direccion"` // Opcional

	// Características de la propiedad
	Tipo        string  `json:"tipo" db:"tipo"`               // casa, departamento, terreno, comercial
	Estado      string  `json:"estado" db:"estado"`           // disponible, vendida, alquilada, reservada
	Dormitorios int     `json:"dormitorios" db:"dormitorios"` // Número de dormitorios
	Banos       float32 `json:"banos" db:"banos"`             // Puede ser 2.5 baños por ejemplo
	AreaM2      float64 `json:"area_m2" db:"area_m2"`         // Área en metros cuadrados

	// Campos de auditoría
	// En Go, time.Time es el tipo estándar para fechas
	FechaCreacion      time.Time `json:"fecha_creacion" db:"fecha_creacion"`
	FechaActualizacion time.Time `json:"fecha_actualizacion" db:"fecha_actualizacion"`
}

// NuevaPropiedad crea una nueva propiedad con slug SEO generado automáticamente
func NuevaPropiedad(titulo, descripcion, provincia, ciudad, tipo string, precio float64) *Propiedad {
	// uuid.New() genera un UUID v4 aleatorio
	id := uuid.New().String()
	
	// Generar slug SEO desde el título
	slug := GenerarSlug(titulo, id)
	
	return &Propiedad{
		ID:                 id,
		Slug:               slug,
		Titulo:             titulo,
		Descripcion:        descripcion,
		Precio:             precio,
		Provincia:          provincia,
		Ciudad:             ciudad,
		Tipo:               tipo,
		Estado:             "disponible", // Estado por defecto
		Dormitorios:        0,
		Banos:              0,
		AreaM2:             0,
		FechaCreacion:      time.Now(),
		FechaActualizacion: time.Now(),
	}
}

// ActualizarFecha actualiza la fecha de modificación
func (p *Propiedad) ActualizarFecha() {
	p.FechaActualizacion = time.Now()
}

// EsValida valida los campos obligatorios de la propiedad
func (p *Propiedad) EsValida() bool {
	// En Go, usamos && para AND lógico
	return p.Titulo != "" &&
		p.Precio > 0 &&
		p.Provincia != "" &&
		p.Ciudad != "" &&
		p.Tipo != ""
}

// Constantes para tipos de propiedad
// En Go, las constantes se definen con 'const'
const (
	TipoCasa         = "casa"
	TipoDepartamento = "departamento"
	TipoTerreno      = "terreno"
	TipoComercial    = "comercial"
)

// Constantes para estados
const (
	EstadoDisponible = "disponible"
	EstadoVendida    = "vendida"
	EstadoAlquilada  = "alquilada"
	EstadoReservada  = "reservada"
)

// ProvinciasEcuador lista las provincias válidas del Ecuador
var ProvinciasEcuador = []string{
	"Azuay", "Bolívar", "Cañar", "Carchi", "Chimborazo",
	"Cotopaxi", "El Oro", "Esmeraldas", "Galápagos",
	"Guayas", "Imbabura", "Loja", "Los Ríos", "Manabí",
	"Morona Santiago", "Napo", "Orellana", "Pastaza",
	"Pichincha", "Santa Elena", "Santo Domingo",
	"Sucumbíos", "Tungurahua", "Zamora Chinchipe",
}

// EsProvinciaValida verifica si una provincia es válida en Ecuador
func EsProvinciaValida(provincia string) bool {
	for _, p := range ProvinciasEcuador {
		if p == provincia {
			return true
		}
	}
	return false
}

// GenerarSlug crea un slug SEO amigable desde el título de la propiedad
func GenerarSlug(titulo, id string) string {
	// Convertir a minúsculas
	slug := strings.ToLower(titulo)
	
	// Reemplazar caracteres especiales y espacios con guiones
	// Esta expresión regular mantiene solo letras, números y espacios
	slug = regexp.MustCompile(`[^a-záéíóúñ0-9\s]+`).ReplaceAllString(slug, "")
	
	// Reemplazar espacios múltiples con un solo espacio
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, " ")
	
	// Convertir espacios a guiones
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Remover guiones al inicio y final
	slug = strings.Trim(slug, "-")
	
	// Truncar si es muy largo (máximo 50 caracteres antes del ID)
	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.Trim(slug, "-")
	}
	
	// Agregar ID corto al final para evitar duplicados
	// Tomamos solo los primeros 8 caracteres del UUID
	idCorto := id
	if len(id) > 8 {
		idCorto = id[:8]
	}
	
	return slug + "-" + idCorto
}

// ActualizarSlug regenera el slug cuando el título cambia
func (p *Propiedad) ActualizarSlug() {
	p.Slug = GenerarSlug(p.Titulo, p.ID)
}

// EsSlugValido verifica si un string puede ser un slug válido
func EsSlugValido(slug string) bool {
	if slug == "" {
		return false
	}
	
	// Un slug válido solo contiene letras minúsculas, números y guiones
	// No puede empezar o terminar con guión
	matched, _ := regexp.MatchString(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`, slug)
	return matched
}