package dominio

import (
	"fmt"
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
	Sector    string `json:"sector" db:"sector"`       // Opcional
	Direccion string `json:"direccion" db:"direccion"` // Opcional

	// Geolocalización para mapas
	Latitud            float64 `json:"latitud" db:"latitud"`                         // Coordenada GPS latitud (-4.0 a 2.0 para Ecuador)
	Longitud           float64 `json:"longitud" db:"longitud"`                       // Coordenada GPS longitud (-92.0 a -75.0 para Ecuador)
	PrecisionUbicacion string  `json:"precision_ubicacion" db:"precision_ubicacion"` // exacta, aproximada, sector

	// Características de la propiedad
	Tipo        string  `json:"tipo" db:"tipo"`               // casa, departamento, terreno, comercial
	Estado      string  `json:"estado" db:"estado"`           // disponible, vendida, alquilada, reservada
	Dormitorios int     `json:"dormitorios" db:"dormitorios"` // Número de dormitorios
	Banos       float32 `json:"banos" db:"banos"`             // Puede ser 2.5 baños por ejemplo
	AreaM2      float64 `json:"area_m2" db:"area_m2"`         // Área en metros cuadrados

	// Imágenes y Media
	ImagenPrincipal string   `json:"imagen_principal" db:"imagen_principal"` // URL de imagen principal
	Imagenes        []string `json:"imagenes,omitempty" db:"imagenes"`       // Array de URLs de imágenes (JSON en DB)
	VideoTour       string   `json:"video_tour,omitempty" db:"video_tour"`   // URL de video tour (opcional)
	Tour360         string   `json:"tour_360,omitempty" db:"tour_360"`       // URL de tour virtual 360° (opcional)

	// Precios Adicionales
	PrecioAlquiler *float64 `json:"precio_alquiler,omitempty" db:"precio_alquiler"` // Para propiedades en alquiler
	GastosComunes  *float64 `json:"gastos_comunes,omitempty" db:"gastos_comunes"`   // Gastos de condominio/administración
	PrecioM2       *float64 `json:"precio_m2,omitempty" db:"precio_m2"`             // Precio por metro cuadrado

	// Características Detalladas
	AnoConstruccion *int   `json:"ano_construccion,omitempty" db:"ano_construccion"` // Año de construcción
	Pisos           *int   `json:"pisos,omitempty" db:"pisos"`                       // Número de pisos de la propiedad
	EstadoPropiedad string `json:"estado_propiedad" db:"estado_propiedad"`           // nueva, usada, remodelada
	Amoblada        bool   `json:"amoblada" db:"amoblada"`                           // Si viene amoblada

	// Amenidades (para filtros del frontend)
	Garage            bool `json:"garage" db:"garage"`                         // Tiene garage/estacionamiento
	Piscina           bool `json:"piscina" db:"piscina"`                       // Tiene piscina
	Jardin            bool `json:"jardin" db:"jardin"`                         // Tiene jardín
	Terraza           bool `json:"terraza" db:"terraza"`                       // Tiene terraza
	Balcon            bool `json:"balcon" db:"balcon"`                         // Tiene balcón
	Seguridad         bool `json:"seguridad" db:"seguridad"`                   // Condominio con seguridad
	Ascensor          bool `json:"ascensor" db:"ascensor"`                     // Tiene ascensor
	AireAcondicionado bool `json:"aire_acondicionado" db:"aire_acondicionado"` // Tiene aire acondicionado

	// Marketing y SEO
	Tags            []string `json:"tags,omitempty" db:"tags"`               // Tags para búsqueda ["lujo", "vista-al-mar"]
	Destacada       bool     `json:"destacada" db:"destacada"`               // Propiedad destacada/premium
	VisitasContador int      `json:"visitas_contador" db:"visitas_contador"` // Contador de vistas

	// Relación con Inmobiliaria (Foreign Key)
	InmobiliariaID *string `json:"inmobiliaria_id,omitempty" db:"inmobiliaria_id"` // FK opcional

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
		ID:          id,
		Slug:        slug,
		Titulo:      titulo,
		Descripcion: descripcion,
		Precio:      precio,
		Provincia:   provincia,
		Ciudad:      ciudad,
		Tipo:        tipo,
		Estado:      EstadoDisponible, // Estado por defecto usando constante
		Dormitorios: 0,
		Banos:       0,
		AreaM2:      0,
		// Geolocalización por defecto (0,0 indica que no está configurada)
		Latitud:            0,
		Longitud:           0,
		PrecisionUbicacion: PrecisionSector, // Precisión por defecto
		// Imágenes y media (vacías por defecto)
		ImagenPrincipal: "",
		Imagenes:        []string{},
		VideoTour:       "",
		Tour360:         "",
		// Precios adicionales (nil = no definidos)
		PrecioAlquiler: nil,
		GastosComunes:  nil,
		PrecioM2:       nil,
		// Características detalladas
		AnoConstruccion: nil,
		Pisos:           nil,
		EstadoPropiedad: EstadoPropiedadUsada, // Por defecto usada
		Amoblada:        false,
		// Amenidades (false por defecto)
		Garage:            false,
		Piscina:           false,
		Jardin:            false,
		Terraza:           false,
		Balcon:            false,
		Seguridad:         false,
		Ascensor:          false,
		AireAcondicionado: false,
		// Marketing y SEO
		Tags:            []string{},
		Destacada:       false,
		VisitasContador: 0,
		// Relación
		InmobiliariaID: nil, // Sin inmobiliaria por defecto
		// Auditoría
		FechaCreacion:      time.Now(),
		FechaActualizacion: time.Now(),
	}
}

// NuevaPropiedadConInmobiliaria crea una nueva propiedad asociada a una inmobiliaria
func NuevaPropiedadConInmobiliaria(titulo, descripcion, provincia, ciudad, tipo string, precio float64, inmobiliariaID string) *Propiedad {
	propiedad := NuevaPropiedad(titulo, descripcion, provincia, ciudad, tipo, precio)
	propiedad.InmobiliariaID = &inmobiliariaID
	return propiedad
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

// Constantes para precisión de ubicación
const (
	PrecisionExacta     = "exacta"     // Coordenadas exactas de la propiedad
	PrecisionAproximada = "aproximada" // Coordenadas aproximadas (cuadra)
	PrecisionSector     = "sector"     // Solo a nivel de sector/barrio
)

// Constantes para estado de la propiedad
const (
	EstadoPropiedadNueva      = "nueva"      // Propiedad nueva/a estrenar
	EstadoPropiedadUsada      = "usada"      // Propiedad usada/habitada
	EstadoPropiedadRemodelada = "remodelada" // Propiedad remodelada/renovada
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

// EsCoordenadasValidasEcuador verifica si las coordenadas están dentro de Ecuador
func EsCoordenadasValidasEcuador(latitud, longitud float64) bool {
	// Límites geográficos aproximados de Ecuador
	// Latitud: -5.0 (sur) a 2.0 (norte)
	// Longitud: -92.0 (oeste) a -75.0 (este)

	latitudValida := latitud >= -5.0 && latitud <= 2.0
	longitudValida := longitud >= -92.0 && longitud <= -75.0

	return latitudValida && longitudValida
}

// EsPrecisionUbicacionValida verifica si la precisión de ubicación es válida
func EsPrecisionUbicacionValida(precision string) bool {
	precisionesValidas := []string{PrecisionExacta, PrecisionAproximada, PrecisionSector}

	for _, p := range precisionesValidas {
		if p == precision {
			return true
		}
	}
	return false
}

// TieneUbicacionConfigurada verifica si la propiedad tiene coordenadas GPS configuradas
func (p *Propiedad) TieneUbicacionConfigurada() bool {
	// (0,0) indica que no está configurada
	return p.Latitud != 0 || p.Longitud != 0
}

// ConfigurarUbicacion establece las coordenadas GPS de la propiedad
func (p *Propiedad) ConfigurarUbicacion(latitud, longitud float64, precision string) error {
	// Validar coordenadas
	if !EsCoordenadasValidasEcuador(latitud, longitud) {
		return fmt.Errorf("coordenadas fuera del territorio ecuatoriano: lat=%.6f, lng=%.6f", latitud, longitud)
	}

	// Validar precisión
	if !EsPrecisionUbicacionValida(precision) {
		return fmt.Errorf("precisión de ubicación no válida: %s", precision)
	}

	p.Latitud = latitud
	p.Longitud = longitud
	p.PrecisionUbicacion = precision
	p.ActualizarFecha()

	return nil
}

// AsociarInmobiliaria asocia la propiedad con una inmobiliaria
func (p *Propiedad) AsociarInmobiliaria(inmobiliariaID string) {
	if strings.TrimSpace(inmobiliariaID) != "" {
		p.InmobiliariaID = &inmobiliariaID
		p.ActualizarFecha()
	}
}

// DesasociarInmobiliaria remueve la asociación con la inmobiliaria
func (p *Propiedad) DesasociarInmobiliaria() {
	p.InmobiliariaID = nil
	p.ActualizarFecha()
}

// TieneInmobiliaria verifica si la propiedad está asociada a una inmobiliaria
func (p *Propiedad) TieneInmobiliaria() bool {
	return p.InmobiliariaID != nil && *p.InmobiliariaID != ""
}

// ObtenerInmobiliariaID retorna el ID de la inmobiliaria asociada
func (p *Propiedad) ObtenerInmobiliariaID() string {
	if p.TieneInmobiliaria() {
		return *p.InmobiliariaID
	}
	return ""
}

// ConfigurarImagenes establece las imágenes de la propiedad
func (p *Propiedad) ConfigurarImagenes(imagenPrincipal string, imagenes []string) {
	p.ImagenPrincipal = imagenPrincipal
	p.Imagenes = imagenes
	p.ActualizarFecha()
}

// AgregarImagen añade una imagen a la galería
func (p *Propiedad) AgregarImagen(url string) {
	if url != "" {
		p.Imagenes = append(p.Imagenes, url)
		p.ActualizarFecha()
	}
}

// TieneImagenes verifica si la propiedad tiene imágenes configuradas
func (p *Propiedad) TieneImagenes() bool {
	return p.ImagenPrincipal != "" || len(p.Imagenes) > 0
}

// ConfigurarPrecioAlquiler establece el precio de alquiler
func (p *Propiedad) ConfigurarPrecioAlquiler(precio float64) {
	if precio > 0 {
		p.PrecioAlquiler = &precio
		p.ActualizarFecha()
	}
}

// CalcularPrecioM2 calcula y establece el precio por metro cuadrado
func (p *Propiedad) CalcularPrecioM2() {
	if p.AreaM2 > 0 {
		precioM2 := p.Precio / p.AreaM2
		p.PrecioM2 = &precioM2
		p.ActualizarFecha()
	}
}

// ConfigurarAmenidades establece múltiples amenidades de una vez
func (p *Propiedad) ConfigurarAmenidades(garage, piscina, jardin, seguridad bool) {
	p.Garage = garage
	p.Piscina = piscina
	p.Jardin = jardin
	p.Seguridad = seguridad
	p.ActualizarFecha()
}

// AgregarTag añade un tag para búsquedas
func (p *Propiedad) AgregarTag(tag string) {
	if tag != "" && !p.TieneTag(tag) {
		p.Tags = append(p.Tags, strings.ToLower(strings.TrimSpace(tag)))
		p.ActualizarFecha()
	}
}

// TieneTag verifica si la propiedad tiene un tag específico
func (p *Propiedad) TieneTag(tag string) bool {
	tagLower := strings.ToLower(strings.TrimSpace(tag))
	for _, t := range p.Tags {
		if t == tagLower {
			return true
		}
	}
	return false
}

// MarcarComoDestacada marca la propiedad como destacada
func (p *Propiedad) MarcarComoDestacada() {
	p.Destacada = true
	p.ActualizarFecha()
}

// QuitarDestacado quita el estado destacado
func (p *Propiedad) QuitarDestacado() {
	p.Destacada = false
	p.ActualizarFecha()
}

// IncrementarVisitas incrementa el contador de visitas
func (p *Propiedad) IncrementarVisitas() {
	p.VisitasContador++
	// No actualizamos fecha para visitas (sería demasiado frecuente)
}

// EsEstadoPropiedadValido verifica si el estado de la propiedad es válido
func EsEstadoPropiedadValido(estado string) bool {
	estadosValidos := []string{EstadoPropiedadNueva, EstadoPropiedadUsada, EstadoPropiedadRemodelada}

	for _, e := range estadosValidos {
		if e == estado {
			return true
		}
	}
	return false
}

// ConfigurarCaracteristicas establece las características detalladas de la propiedad
func (p *Propiedad) ConfigurarCaracteristicas(anoConstruccion, pisos int, estado string, amoblada bool) error {
	// Validar año de construcción
	anoActual := time.Now().Year()
	if anoConstruccion > 0 && (anoConstruccion < 1800 || anoConstruccion > anoActual+2) {
		return fmt.Errorf("año de construcción inválido: %d", anoConstruccion)
	}

	// Validar estado de propiedad
	if !EsEstadoPropiedadValido(estado) {
		return fmt.Errorf("estado de propiedad inválido: %s", estado)
	}

	// Validar número de pisos
	if pisos > 0 && pisos > 50 {
		return fmt.Errorf("número de pisos inválido: %d", pisos)
	}

	if anoConstruccion > 0 {
		p.AnoConstruccion = &anoConstruccion
	}
	if pisos > 0 {
		p.Pisos = &pisos
	}
	p.EstadoPropiedad = estado
	p.Amoblada = amoblada
	p.ActualizarFecha()

	return nil
}

// ObtenerResumenAmenidades retorna un resumen de las amenidades para mostrar
func (p *Propiedad) ObtenerResumenAmenidades() []string {
	var amenidades []string

	if p.Garage {
		amenidades = append(amenidades, "Garage")
	}
	if p.Piscina {
		amenidades = append(amenidades, "Piscina")
	}
	if p.Jardin {
		amenidades = append(amenidades, "Jardín")
	}
	if p.Terraza {
		amenidades = append(amenidades, "Terraza")
	}
	if p.Seguridad {
		amenidades = append(amenidades, "Seguridad")
	}
	if p.Ascensor {
		amenidades = append(amenidades, "Ascensor")
	}
	if p.AireAcondicionado {
		amenidades = append(amenidades, "Aire Acondicionado")
	}

	return amenidades
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

	// Un slug válido solo contiene letras minúsculas, números, acentos y guiones
	// No puede empezar o terminar con guión
	matched, _ := regexp.MatchString(`^[a-záéíóúñ0-9]([a-záéíóúñ0-9-]*[a-záéíóúñ0-9])?$`, slug)
	return matched
}
