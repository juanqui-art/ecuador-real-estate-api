package dominio

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Inmobiliaria representa una empresa inmobiliaria
type Inmobiliaria struct {
	// ID es el identificador único de la inmobiliaria
	ID string `json:"id" db:"id"`

	// Información básica de la empresa
	Nombre      string `json:"nombre" db:"nombre"`
	RUC         string `json:"ruc" db:"ruc"`             // RUC ecuatoriano (13 dígitos)
	Direccion   string `json:"direccion" db:"direccion"` // Dirección física
	Descripcion string `json:"descripcion" db:"descripcion"`

	// Información de contacto
	Telefono string `json:"telefono" db:"telefono"`
	Email    string `json:"email" db:"email"`
	SitioWeb string `json:"sitio_web" db:"sitio_web"`
	LogoURL  string `json:"logo_url" db:"logo_url"` // URL del logo

	// Estado de la inmobiliaria
	Activa bool `json:"activa" db:"activa"`

	// Campos de auditoría
	FechaCreacion      time.Time `json:"fecha_creacion" db:"fecha_creacion"`
	FechaActualizacion time.Time `json:"fecha_actualizacion" db:"fecha_actualizacion"`
}

// NuevaInmobiliaria crea una nueva inmobiliaria
func NuevaInmobiliaria(nombre, ruc, direccion, telefono, email string) *Inmobiliaria {
	id := uuid.New().String()

	return &Inmobiliaria{
		ID:                 id,
		Nombre:             nombre,
		RUC:                ruc,
		Direccion:          direccion,
		Telefono:           telefono,
		Email:              email,
		Activa:             true, // Por defecto está activa
		FechaCreacion:      time.Now(),
		FechaActualizacion: time.Now(),
	}
}

// ActualizarFecha actualiza la fecha de modificación
func (i *Inmobiliaria) ActualizarFecha() {
	i.FechaActualizacion = time.Now()
}

// Validar valida los campos obligatorios de la inmobiliaria
func (i *Inmobiliaria) Validar() error {
	if strings.TrimSpace(i.Nombre) == "" {
		return fmt.Errorf("nombre es requerido")
	}
	if strings.TrimSpace(i.RUC) == "" {
		return fmt.Errorf("RUC es requerido")
	}
	if strings.TrimSpace(i.Direccion) == "" {
		return fmt.Errorf("dirección es requerida")
	}
	if strings.TrimSpace(i.Telefono) == "" {
		return fmt.Errorf("teléfono es requerido")
	}
	if strings.TrimSpace(i.Email) == "" {
		return fmt.Errorf("email es requerido")
	}

	// Validar formato de email
	if err := i.ValidarEmail(); err != nil {
		return err
	}

	// Validar formato de teléfono
	if err := i.ValidarTelefono(); err != nil {
		return err
	}

	// Validar formato de RUC
	if err := i.ValidarRUC(); err != nil {
		return err
	}

	return nil
}

// ValidarRUC valida el formato del RUC ecuatoriano
func (i *Inmobiliaria) ValidarRUC() error {
	ruc := strings.TrimSpace(i.RUC)
	if len(ruc) != 13 {
		return fmt.Errorf("RUC debe tener 13 dígitos")
	}

	// Verificar que solo contenga números
	if !regexp.MustCompile(`^\d{13}$`).MatchString(ruc) {
		return fmt.Errorf("RUC debe contener solo números")
	}

	// Verificar que termine en 001 (empresas)
	if !strings.HasSuffix(ruc, "001") {
		return fmt.Errorf("RUC de empresa debe terminar en 001")
	}

	return nil
}

// ValidarEmail valida el formato del email
func (i *Inmobiliaria) ValidarEmail() error {
	if !EsEmailValido(i.Email) {
		return fmt.Errorf("formato de email inválido: %s", i.Email)
	}
	return nil
}

// ValidarTelefono valida el formato del teléfono ecuatoriano
func (i *Inmobiliaria) ValidarTelefono() error {
	if !EsTelefonoValidoEcuador(i.Telefono) {
		return fmt.Errorf("formato de teléfono inválido para Ecuador: %s", i.Telefono)
	}
	return nil
}

// Activar marca la inmobiliaria como activa
func (i *Inmobiliaria) Activar() {
	i.Activa = true
	i.ActualizarFecha()
}

// Desactivar marca la inmobiliaria como inactiva
func (i *Inmobiliaria) Desactivar() {
	i.Activa = false
	i.ActualizarFecha()
}

// ConfigurarRUC establece el RUC de la inmobiliaria
func (i *Inmobiliaria) ConfigurarRUC(ruc string) error {
	i.RUC = strings.TrimSpace(ruc)
	if err := i.ValidarRUC(); err != nil {
		return err
	}
	i.ActualizarFecha()
	return nil
}

// ConfigurarSitioWeb establece la URL del sitio web
func (i *Inmobiliaria) ConfigurarSitioWeb(url string) error {
	if url != "" && !EsURLValida(url) {
		return fmt.Errorf("URL del sitio web no es válida: %s", url)
	}

	i.SitioWeb = url
	i.ActualizarFecha()
	return nil
}

// EsEmailValido verifica si el email tiene un formato válido
func EsEmailValido(email string) bool {
	// Expresión regular básica para validar emails
	patron := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(patron, email)
	return matched
}

// EsURLValida verifica si la URL tiene un formato válido
func EsURLValida(url string) bool {
	// Expresión regular básica para validar URLs
	patron := `^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`
	matched, _ := regexp.MatchString(patron, url)
	return matched
}

// EsTelefonoValidoEcuador verifica si el teléfono tiene formato ecuatoriano
func EsTelefonoValidoEcuador(telefono string) bool {
	// Remover espacios y caracteres especiales
	telefonoLimpio := regexp.MustCompile(`[^\d]`).ReplaceAllString(telefono, "")

	// Teléfonos ecuatorianos:
	// Fijos: 07-XXXXXXX (9 dígitos con código de área)
	// Móviles: 09-XXXXXXXX (10 dígitos)
	// Con código país: +593-X-XXXXXXX

	// Verificar longitud y patrones comunes
	if len(telefonoLimpio) == 9 {
		// Teléfono fijo sin código país
		return strings.HasPrefix(telefonoLimpio, "0")
	} else if len(telefonoLimpio) == 10 {
		// Teléfono móvil sin código país
		return strings.HasPrefix(telefonoLimpio, "09")
	} else if len(telefonoLimpio) == 12 {
		// Con código país +593
		return strings.HasPrefix(telefonoLimpio, "593")
	}

	return false
}

// FormatearTelefono normaliza el formato del teléfono ecuatoriano
func FormatearTelefono(telefono string) string {
	// Remover espacios y caracteres especiales excepto +
	telefonoLimpio := regexp.MustCompile(`[^\d+]`).ReplaceAllString(telefono, "")

	// Si empieza con +593, mantener el formato
	if strings.HasPrefix(telefonoLimpio, "+593") {
		return telefonoLimpio
	}

	// Si empieza con 593, agregar +
	if strings.HasPrefix(telefonoLimpio, "593") {
		return "+" + telefonoLimpio
	}

	// Si es número local, agregar código país
	if strings.HasPrefix(telefonoLimpio, "0") && len(telefonoLimpio) >= 9 {
		return "+593" + telefonoLimpio[1:] // Remover el 0 inicial
	}

	return telefonoLimpio
}

// ConfigurarContacto actualiza la información de contacto de la inmobiliaria
func (i *Inmobiliaria) ConfigurarContacto(telefono, email string) error {
	// Validar email
	if !EsEmailValido(email) {
		return fmt.Errorf("email no válido: %s", email)
	}

	// Formatear y validar teléfono
	telefonoFormateado := FormatearTelefono(telefono)
	if !EsTelefonoValidoEcuador(telefonoFormateado) {
		return fmt.Errorf("teléfono no válido para Ecuador: %s", telefono)
	}

	i.Telefono = telefonoFormateado
	i.Email = strings.ToLower(strings.TrimSpace(email))
	i.ActualizarFecha()

	return nil
}

// ObtenerNombreCompleto retorna el nombre completo para mostrar
func (i *Inmobiliaria) ObtenerNombreCompleto() string {
	if i.Activa {
		return i.Nombre
	}
	return i.Nombre + " (Inactiva)"
}

// ObtenerResumen retorna un resumen de la inmobiliaria para listados
func (i *Inmobiliaria) ObtenerResumen() map[string]interface{} {
	return map[string]interface{}{
		"id":        i.ID,
		"nombre":    i.ObtenerNombreCompleto(),
		"ruc":       i.RUC,
		"telefono":  i.Telefono,
		"email":     i.Email,
		"direccion": i.Direccion,
		"activa":    i.Activa,
	}
}
