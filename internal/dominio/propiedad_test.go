package dominio

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNuevaPropiedad prueba la creación de una nueva propiedad
func TestNuevaPropiedad(t *testing.T) {
	// Datos de prueba
	titulo := "Casa moderna en Samborondón"
	descripcion := "Hermosa casa con piscina y jardín"
	provincia := "Guayas"
	ciudad := "Samborondón"
	tipo := "casa"
	precio := 250000.0

	// Crear propiedad
	propiedad := NuevaPropiedad(titulo, descripcion, provincia, ciudad, tipo, precio)

	// Verificaciones usando testify
	assert.NotNil(t, propiedad, "La propiedad no debe ser nil")
	assert.NotEmpty(t, propiedad.ID, "ID debe ser generado")
	assert.NotEmpty(t, propiedad.Slug, "Slug debe ser generado")
	assert.Equal(t, titulo, propiedad.Titulo)
	assert.Equal(t, descripcion, propiedad.Descripcion)
	assert.Equal(t, provincia, propiedad.Provincia)
	assert.Equal(t, ciudad, propiedad.Ciudad)
	assert.Equal(t, tipo, propiedad.Tipo)
	assert.Equal(t, precio, propiedad.Precio)
	assert.Equal(t, "disponible", propiedad.Estado, "Estado por defecto debe ser 'disponible'")
	assert.Equal(t, 0, propiedad.Dormitorios, "Dormitorios por defecto debe ser 0")
	assert.Equal(t, float32(0), propiedad.Banos, "Baños por defecto debe ser 0")
	assert.Equal(t, float64(0), propiedad.AreaM2, "Área por defecto debe ser 0")

	// Verificar que las fechas son recientes (dentro de los últimos 5 segundos)
	now := time.Now()
	assert.WithinDuration(t, now, propiedad.FechaCreacion, 5*time.Second)
	assert.WithinDuration(t, now, propiedad.FechaActualizacion, 5*time.Second)
}

// TestPropiedad_EsValida usa table-driven tests para probar la validación
func TestPropiedad_EsValida(t *testing.T) {
	// Table-driven tests - patrón estándar en Go
	tests := []struct {
		name      string
		propiedad *Propiedad
		esperado  bool
	}{
		{
			name: "propiedad válida",
			propiedad: &Propiedad{
				Titulo:    "Casa en Quito",
				Precio:    100000,
				Provincia: "Pichincha",
				Ciudad:    "Quito",
				Tipo:      "casa",
			},
			esperado: true,
		},
		{
			name: "título vacío",
			propiedad: &Propiedad{
				Titulo:    "",
				Precio:    100000,
				Provincia: "Pichincha",
				Ciudad:    "Quito",
				Tipo:      "casa",
			},
			esperado: false,
		},
		{
			name: "precio cero",
			propiedad: &Propiedad{
				Titulo:    "Casa en Quito",
				Precio:    0,
				Provincia: "Pichincha",
				Ciudad:    "Quito",
				Tipo:      "casa",
			},
			esperado: false,
		},
		{
			name: "precio negativo",
			propiedad: &Propiedad{
				Titulo:    "Casa en Quito",
				Precio:    -50000,
				Provincia: "Pichincha",
				Ciudad:    "Quito",
				Tipo:      "casa",
			},
			esperado: false,
		},
		{
			name: "provincia vacía",
			propiedad: &Propiedad{
				Titulo:    "Casa en Quito",
				Precio:    100000,
				Provincia: "",
				Ciudad:    "Quito",
				Tipo:      "casa",
			},
			esperado: false,
		},
		{
			name: "ciudad vacía",
			propiedad: &Propiedad{
				Titulo:    "Casa en Quito",
				Precio:    100000,
				Provincia: "Pichincha",
				Ciudad:    "",
				Tipo:      "casa",
			},
			esperado: false,
		},
		{
			name: "tipo vacío",
			propiedad: &Propiedad{
				Titulo:    "Casa en Quito",
				Precio:    100000,
				Provincia: "Pichincha",
				Ciudad:    "Quito",
				Tipo:      "",
			},
			esperado: false,
		},
	}

	// Ejecutar cada test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := tt.propiedad.EsValida()
			assert.Equal(t, tt.esperado, resultado, "Validación no coincide para: %s", tt.name)
		})
	}
}

// TestPropiedad_ActualizarFecha verifica que se actualice la fecha correctamente
func TestPropiedad_ActualizarFecha(t *testing.T) {
	propiedad := NuevaPropiedad("Test", "Descripción", "Pichincha", "Quito", "casa", 100000)

	// Guardar fecha original
	fechaOriginal := propiedad.FechaActualizacion

	// Esperar un poco para que la fecha sea diferente
	time.Sleep(10 * time.Millisecond)

	// Actualizar fecha
	propiedad.ActualizarFecha()

	// Verificar que la fecha cambió
	assert.True(t, propiedad.FechaActualizacion.After(fechaOriginal),
		"La fecha de actualización debe ser posterior a la original")
}

// TestEsProvinciaValida prueba la validación de provincias ecuatorianas
func TestEsProvinciaValida(t *testing.T) {
	tests := []struct {
		name      string
		provincia string
		esperado  bool
	}{
		{"Pichincha válida", "Pichincha", true},
		{"Guayas válida", "Guayas", true},
		{"Azuay válida", "Azuay", true},
		{"Provincia inexistente", "Madrid", false},
		{"Cadena vacía", "", false},
		{"Caso sensible", "pichincha", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := EsProvinciaValida(tt.provincia)
			assert.Equal(t, tt.esperado, resultado)
		})
	}
}

// TestGenerarSlug prueba la generación de slugs SEO
func TestGenerarSlug(t *testing.T) {
	tests := []struct {
		name     string
		titulo   string
		id       string
		esperado string // Verificaremos que contenga ciertos elementos
	}{
		{
			name:     "título simple",
			titulo:   "Casa en Quito",
			id:       "12345678-1234-1234-1234-123456789012",
			esperado: "casa-en-quito-12345678",
		},
		{
			name:     "título con caracteres especiales",
			titulo:   "¡Casa Súper Moderna!",
			id:       "12345678-1234-1234-1234-123456789012",
			esperado: "casa-súper-moderna-12345678",
		},
		{
			name:     "título con múltiples espacios",
			titulo:   "Casa    con    espacios",
			id:       "12345678-1234-1234-1234-123456789012",
			esperado: "casa-con-espacios-12345678",
		},
		{
			name:     "título muy largo",
			titulo:   "Esta es una casa extremadamente lujosa y moderna con múltiples características",
			id:       "12345678-1234-1234-1234-123456789012",
			esperado: "esta-es-una-casa-extremadamente-lujosa-y-moderna-c-12345678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := GenerarSlug(tt.titulo, tt.id)
			assert.Equal(t, tt.esperado, resultado)

			// Verificar que el slug es válido
			assert.True(t, EsSlugValido(resultado), "El slug generado debe ser válido")
		})
	}
}

// TestEsSlugValido prueba la validación de slugs
func TestEsSlugValido(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		esperado bool
	}{
		{"slug válido", "casa-en-quito-12345678", true},
		{"slug con números", "casa-123-abc", true},
		{"slug simple", "casa", true},
		{"slug vacío", "", false},
		{"slug con mayúsculas", "Casa-En-Quito", false},
		{"slug con espacios", "casa en quito", false},
		{"slug con caracteres especiales", "casa_en_quito", false},
		{"slug que empieza con guión", "-casa-en-quito", false},
		{"slug que termina con guión", "casa-en-quito-", false},
		{"slug solo guiones", "---", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := EsSlugValido(tt.slug)
			assert.Equal(t, tt.esperado, resultado)
		})
	}
}

// TestPropiedad_ActualizarSlug verifica que se actualice el slug correctamente
func TestPropiedad_ActualizarSlug(t *testing.T) {
	propiedad := NuevaPropiedad("Título Original", "Descripción", "Pichincha", "Quito", "casa", 100000)
	slugOriginal := propiedad.Slug

	// Cambiar el título
	propiedad.Titulo = "Nuevo Título"

	// Actualizar slug
	propiedad.ActualizarSlug()

	// Verificar que el slug cambió
	assert.NotEqual(t, slugOriginal, propiedad.Slug)
	assert.Contains(t, propiedad.Slug, "nuevo-título")
	assert.True(t, EsSlugValido(propiedad.Slug))
}

// TestConstantes verifica que las constantes estén definidas correctamente
func TestConstantes(t *testing.T) {
	// Verificar tipos de propiedad
	assert.Equal(t, "casa", TipoCasa)
	assert.Equal(t, "departamento", TipoDepartamento)
	assert.Equal(t, "terreno", TipoTerreno)
	assert.Equal(t, "comercial", TipoComercial)

	// Verificar estados
	assert.Equal(t, "disponible", EstadoDisponible)
	assert.Equal(t, "vendida", EstadoVendida)
	assert.Equal(t, "alquilada", EstadoAlquilada)
	assert.Equal(t, "reservada", EstadoReservada)
}

// TestProvinciasEcuador verifica que tengamos todas las provincias
func TestProvinciasEcuador(t *testing.T) {
	// Verificar que tengamos 24 provincias
	assert.Len(t, ProvinciasEcuador, 24, "Ecuador debe tener 24 provincias")

	// Verificar algunas provincias específicas
	provinciasEsperadas := []string{"Pichincha", "Guayas", "Azuay", "Manabí", "Loja"}

	for _, provincia := range provinciasEsperadas {
		assert.Contains(t, ProvinciasEcuador, provincia, "Debe contener la provincia: %s", provincia)
	}
}

// BenchmarkGenerarSlug - ejemplo de benchmark en Go
func BenchmarkGenerarSlug(b *testing.B) {
	titulo := "Casa moderna en Samborondón con piscina y jardín"
	id := "12345678-1234-1234-1234-123456789012"

	// b.N es el número de iteraciones que Go determina automáticamente
	for i := 0; i < b.N; i++ {
		_ = GenerarSlug(titulo, id)
	}
}

// BenchmarkEsProvinciaValida - benchmark para validación de provincias
func BenchmarkEsProvinciaValida(b *testing.B) {
	provincia := "Pichincha"

	for i := 0; i < b.N; i++ {
		_ = EsProvinciaValida(provincia)
	}
}

// TestPropiedadCompleta - test de integración que verifica todo el flujo
func TestPropiedadCompleta(t *testing.T) {
	// Crear propiedad
	propiedad := NuevaPropiedad(
		"Hermosa casa en Samborondón",
		"Casa moderna de 3 pisos con acabados de lujo",
		"Guayas",
		"Samborondón",
		"casa",
		285000,
	)

	// Verificar que se creó correctamente
	require.NotNil(t, propiedad)
	require.True(t, propiedad.EsValida())

	// Actualizar algunos campos
	propiedad.Dormitorios = 4
	propiedad.Banos = 3.5
	propiedad.AreaM2 = 320.5
	propiedad.Sector = "Vía a la Costa"

	// Cambiar título y verificar que se actualice el slug
	tituloAnterior := propiedad.Titulo
	slugAnterior := propiedad.Slug

	propiedad.Titulo = "Casa de Lujo en Samborondón"
	propiedad.ActualizarSlug()
	propiedad.ActualizarFecha()

	// Verificaciones finales
	assert.NotEqual(t, tituloAnterior, propiedad.Titulo)
	assert.NotEqual(t, slugAnterior, propiedad.Slug)
	assert.True(t, propiedad.EsValida())
	assert.Equal(t, 4, propiedad.Dormitorios)
	assert.Equal(t, float32(3.5), propiedad.Banos)
	assert.Equal(t, 320.5, propiedad.AreaM2)
	assert.Equal(t, "Vía a la Costa", propiedad.Sector)
}
