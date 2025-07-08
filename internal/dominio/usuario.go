package dominio

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Usuario representa un usuario del sistema inmobiliario
type Usuario struct {
	// ID es el identificador único del usuario
	ID string `json:"id" db:"id"`

	// Información personal
	Nombre          string     `json:"nombre" db:"nombre"`
	Apellido        string     `json:"apellido" db:"apellido"`
	Email           string     `json:"email" db:"email"`
	Telefono        string     `json:"telefono" db:"telefono"`
	Cedula          string     `json:"cedula" db:"cedula"`                               // Cédula ecuatoriana (10 dígitos)
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty" db:"fecha_nacimiento"` // Fecha de nacimiento

	// Tipo y estado del usuario
	TipoUsuario string `json:"tipo_usuario" db:"tipo_usuario"` // comprador, vendedor, agente, admin
	Activo      bool   `json:"activo" db:"activo"`             // Usuario activo/inactivo

	// Preferencias de búsqueda (para compradores)
	PresupuestoMin        *float64 `json:"presupuesto_min,omitempty" db:"presupuesto_min"`                 // Presupuesto mínimo
	PresupuestoMax        *float64 `json:"presupuesto_max,omitempty" db:"presupuesto_max"`                 // Presupuesto máximo
	ProvinciasInteres     []string `json:"provincias_interes,omitempty" db:"provincias_interes"`           // Provincias de interés
	TiposPropiedadInteres []string `json:"tipos_propiedad_interes,omitempty" db:"tipos_propiedad_interes"` // Tipos de propiedades de interés

	// Perfil
	AvatarURL string `json:"avatar_url,omitempty" db:"avatar_url"` // URL del avatar
	Bio       string `json:"bio,omitempty" db:"bio"`               // Biografía del usuario

	// Relación con Inmobiliaria (para agentes)
	InmobiliariaID *string `json:"inmobiliaria_id,omitempty" db:"inmobiliaria_id"` // FK hacia inmobiliarias

	// Campos de auditoría
	FechaCreacion      time.Time `json:"fecha_creacion" db:"fecha_creacion"`
	FechaActualizacion time.Time `json:"fecha_actualizacion" db:"fecha_actualizacion"`
}

// Constantes para tipos de usuario
const (
	TipoUsuarioComprador = "comprador" // Usuario que busca comprar/alquilar
	TipoUsuarioVendedor  = "vendedor"  // Usuario que vende/alquila propiedades
	TipoUsuarioAgente    = "agente"    // Agente inmobiliario
	TipoUsuarioAdmin     = "admin"     // Administrador del sistema
)

// Constantes para estados de usuario
const (
	EstadoUsuarioActivo     = "activo"     // Usuario activo
	EstadoUsuarioInactivo   = "inactivo"   // Usuario inactivo (temporal)
	EstadoUsuarioSuspendido = "suspendido" // Usuario suspendido (sanción)
)

// NuevoUsuario crea un nuevo usuario básico
func NuevoUsuario(nombre, apellido, email, telefono, cedula, tipoUsuario string) *Usuario {
	id := uuid.New().String()

	return &Usuario{
		ID:                    id,
		Nombre:                strings.TrimSpace(nombre),
		Apellido:              strings.TrimSpace(apellido),
		Email:                 strings.ToLower(strings.TrimSpace(email)),
		Telefono:              telefono,
		Cedula:                cedula,
		TipoUsuario:           tipoUsuario,
		Activo:                true, // Activo por defecto
		ProvinciasInteres:     []string{},
		TiposPropiedadInteres: []string{},
		FechaCreacion:         time.Now(),
		FechaActualizacion:    time.Now(),
	}
}

// ActualizarFecha actualiza la fecha de modificación
func (u *Usuario) ActualizarFecha() {
	u.FechaActualizacion = time.Now()
}

// Validar valida los campos obligatorios del usuario
func (u *Usuario) Validar() error {
	if strings.TrimSpace(u.Nombre) == "" {
		return fmt.Errorf("nombre es requerido")
	}
	if strings.TrimSpace(u.Apellido) == "" {
		return fmt.Errorf("apellido es requerido")
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("email es requerido")
	}
	if strings.TrimSpace(u.Telefono) == "" {
		return fmt.Errorf("teléfono es requerido")
	}
	if strings.TrimSpace(u.Cedula) == "" {
		return fmt.Errorf("cédula es requerida")
	}

	// Validaciones de formato
	if err := u.ValidarEmail(); err != nil {
		return err
	}
	if err := u.ValidarTelefono(); err != nil {
		return err
	}
	if err := u.ValidarCedula(); err != nil {
		return err
	}
	if err := u.ValidarTipoUsuario(); err != nil {
		return err
	}

	return nil
}

// ValidarEmail valida el formato del email
func (u *Usuario) ValidarEmail() error {
	if !EsEmailValido(u.Email) {
		return fmt.Errorf("formato de email inválido: %s", u.Email)
	}
	return nil
}

// ValidarTelefono valida el formato del teléfono ecuatoriano
func (u *Usuario) ValidarTelefono() error {
	if !EsTelefonoValidoEcuador(u.Telefono) {
		return fmt.Errorf("formato de teléfono inválido para Ecuador: %s", u.Telefono)
	}
	return nil
}

// ValidarCedula valida la cédula ecuatoriana
func (u *Usuario) ValidarCedula() error {
	cedula := strings.TrimSpace(u.Cedula)
	if len(cedula) != 10 {
		return fmt.Errorf("cédula debe tener 10 dígitos")
	}

	// Verificar que solo contenga números
	if !regexp.MustCompile(`^\d{10}$`).MatchString(cedula) {
		return fmt.Errorf("cédula debe contener solo números")
	}

	return nil
}

// ValidarTipoUsuario valida el tipo de usuario
func (u *Usuario) ValidarTipoUsuario() error {
	tiposValidos := []string{"comprador", "vendedor", "agente", "admin"}
	for _, tipo := range tiposValidos {
		if u.TipoUsuario == tipo {
			return nil
		}
	}
	return fmt.Errorf("tipo de usuario inválido: %s", u.TipoUsuario)
}

// ValidarPresupuesto valida el presupuesto de un comprador
func (u *Usuario) ValidarPresupuesto() error {
	if u.PresupuestoMin != nil && *u.PresupuestoMin < 0 {
		return fmt.Errorf("presupuesto mínimo no puede ser negativo")
	}
	if u.PresupuestoMax != nil && *u.PresupuestoMax < 0 {
		return fmt.Errorf("presupuesto máximo no puede ser negativo")
	}
	if u.PresupuestoMin != nil && u.PresupuestoMax != nil && *u.PresupuestoMin > *u.PresupuestoMax {
		return fmt.Errorf("presupuesto mínimo no puede ser mayor al máximo")
	}
	return nil
}

// Activar activa el usuario
func (u *Usuario) Activar() {
	u.Activo = true
	u.ActualizarFecha()
}

// Desactivar desactiva el usuario
func (u *Usuario) Desactivar() {
	u.Activo = false
	u.ActualizarFecha()
}

// ConfigurarPerfil actualiza la información del perfil
func (u *Usuario) ConfigurarPerfil(avatarURL, bio string) {
	u.AvatarURL = avatarURL
	u.Bio = bio
	u.ActualizarFecha()
}

// AsignarInmobiliaria asigna una inmobiliaria a un agente
func (u *Usuario) AsignarInmobiliaria(inmobiliariaID string) {
	u.InmobiliariaID = &inmobiliariaID
	u.ActualizarFecha()
}

// ConfigurarPreferencias establece las preferencias de búsqueda
func (u *Usuario) ConfigurarPreferencias(provincias, tipos []string) {
	u.ProvinciasInteres = provincias
	u.TiposPropiedadInteres = tipos
	u.ActualizarFecha()
}

// AgregarTipoInteres añade un tipo de propiedad de interés
func (u *Usuario) AgregarTipoInteres(tipo string) error {
	// Validar que el tipo sea válido
	tiposValidos := []string{TipoCasa, TipoDepartamento, TipoTerreno, TipoComercial}
	tipoValido := false
	for _, t := range tiposValidos {
		if t == tipo {
			tipoValido = true
			break
		}
	}

	if !tipoValido {
		return fmt.Errorf("tipo de propiedad no válido: %s", tipo)
	}

	// Verificar que no esté ya agregado
	for _, t := range u.TiposPropiedadInteres {
		if t == tipo {
			return nil // Ya existe, no hacer nada
		}
	}

	u.TiposPropiedadInteres = append(u.TiposPropiedadInteres, tipo)
	u.ActualizarFecha()
	return nil
}

// AgregarProvinciaInteres añade una provincia de interés
func (u *Usuario) AgregarProvinciaInteres(provincia string) error {
	if !EsProvinciaValida(provincia) {
		return fmt.Errorf("provincia no válida: %s", provincia)
	}

	// Verificar que no esté ya agregada
	for _, p := range u.ProvinciasInteres {
		if p == provincia {
			return nil // Ya existe, no hacer nada
		}
	}

	u.ProvinciasInteres = append(u.ProvinciasInteres, provincia)
	u.ActualizarFecha()
	return nil
}

// ConfigurarNotificaciones establece las preferencias de notificaciones
func (u *Usuario) ConfigurarNotificaciones(notificaciones, newsletter bool) {
	// Note: These fields would need to be added to the struct if needed
	u.ActualizarFecha()
}

// EsTipoUsuarioValido verifica si el tipo de usuario es válido
func EsTipoUsuarioValido(tipo string) bool {
	tiposValidos := []string{TipoUsuarioComprador, TipoUsuarioVendedor, TipoUsuarioAgente, TipoUsuarioAdmin}

	for _, t := range tiposValidos {
		if t == tipo {
			return true
		}
	}
	return false
}

// EsEstadoUsuarioValido verifica si el estado del usuario es válido
func EsEstadoUsuarioValido(estado string) bool {
	estadosValidos := []string{EstadoUsuarioActivo, EstadoUsuarioInactivo, EstadoUsuarioSuspendido}

	for _, e := range estadosValidos {
		if e == estado {
			return true
		}
	}
	return false
}

// EsCedulaEcuatoriana verifica si una cédula tiene formato ecuatoriano válido
func EsCedulaEcuatoriana(cedula string) bool {
	// Remover espacios y guiones
	cedulaLimpia := strings.ReplaceAll(strings.ReplaceAll(cedula, " ", ""), "-", "")

	// Debe tener exactamente 10 dígitos
	if len(cedulaLimpia) != 10 {
		return false
	}

	// Verificar que todos sean dígitos
	for _, r := range cedulaLimpia {
		if r < '0' || r > '9' {
			return false
		}
	}

	// Los primeros dos dígitos deben corresponder a una provincia (01-24)
	provincia := cedulaLimpia[:2]
	if provincia < "01" || provincia > "24" {
		return false
	}

	// Algoritmo de verificación de cédula ecuatoriana
	digitos := make([]int, 10)
	for i, r := range cedulaLimpia {
		digitos[i] = int(r - '0')
	}

	// El tercer dígito debe ser menor a 6 (para personas naturales)
	if digitos[2] >= 6 {
		return false
	}

	// Verificar dígito verificador
	suma := 0
	for i := 0; i < 9; i++ {
		if i%2 == 0 {
			// Posiciones pares (0,2,4,6,8)
			producto := digitos[i] * 2
			if producto > 9 {
				producto -= 9
			}
			suma += producto
		} else {
			// Posiciones impares (1,3,5,7)
			suma += digitos[i]
		}
	}

	digito_verificador := (10 - (suma % 10)) % 10
	return digito_verificador == digitos[9]
}

// ConfigurarCedula establece y valida la cédula del usuario
func (u *Usuario) ConfigurarCedula(cedula string) error {
	if cedula != "" && !EsCedulaEcuatoriana(cedula) {
		return fmt.Errorf("cédula ecuatoriana no válida: %s", cedula)
	}

	u.Cedula = cedula
	u.ActualizarFecha()
	return nil
}

// ObtenerNombreCompleto retorna el nombre completo del usuario
func (u *Usuario) ObtenerNombreCompleto() string {
	return strings.TrimSpace(u.Nombre + " " + u.Apellido)
}

// ObtenerResumen retorna un resumen del usuario para listados
func (u *Usuario) ObtenerResumen() map[string]interface{} {
	return map[string]interface{}{
		"id":             u.ID,
		"nombre":         u.ObtenerNombreCompleto(),
		"email":          u.Email,
		"telefono":       u.Telefono,
		"tipo_usuario":   u.TipoUsuario,
		"activo":         u.Activo,
		"fecha_creacion": u.FechaCreacion,
	}
}

// PuedeRecibirNotificaciones verifica si el usuario puede recibir notificaciones
func (u *Usuario) PuedeRecibirNotificaciones() bool {
	return u.Activo
}

// EsComprador verifica si es un usuario comprador
func (u *Usuario) EsComprador() bool {
	return u.TipoUsuario == TipoUsuarioComprador
}

// EsVendedor verifica si es un usuario vendedor
func (u *Usuario) EsVendedor() bool {
	return u.TipoUsuario == TipoUsuarioVendedor
}

// EsAgente verifica si es un agente inmobiliario
func (u *Usuario) EsAgente() bool {
	return u.TipoUsuario == TipoUsuarioAgente
}

// EsAdmin verifica si es un administrador
func (u *Usuario) EsAdmin() bool {
	return u.TipoUsuario == TipoUsuarioAdmin
}

// TienePresupuestoConfigurado verifica si tiene presupuesto configurado
func (u *Usuario) TienePresupuestoConfigurado() bool {
	return u.PresupuestoMin != nil || u.PresupuestoMax != nil
}

// PropiedadEnRangoPresupuesto verifica si una propiedad está en el rango de presupuesto
func (u *Usuario) PropiedadEnRangoPresupuesto(precioPropiedad float64) bool {
	if !u.TienePresupuestoConfigurado() {
		return true // Si no tiene presupuesto, todas las propiedades son válidas
	}

	if u.PresupuestoMin != nil && precioPropiedad < *u.PresupuestoMin {
		return false
	}

	if u.PresupuestoMax != nil && precioPropiedad > *u.PresupuestoMax {
		return false
	}

	return true
}
