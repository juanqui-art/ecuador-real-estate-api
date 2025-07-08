package servicio

import (
	"database/sql"
	"testing"

	"realty-core/internal/dominio"
	"realty-core/internal/repositorio"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockPropiedadRepository es un mock más avanzado usando testify/mock
type MockPropiedadRepository struct {
	mock.Mock
}

func (m *MockPropiedadRepository) Crear(propiedad *dominio.Propiedad) error {
	args := m.Called(propiedad)
	return args.Error(0)
}

func (m *MockPropiedadRepository) ObtenerPorID(id string) (*dominio.Propiedad, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dominio.Propiedad), args.Error(1)
}

func (m *MockPropiedadRepository) ObtenerPorSlug(slug string) (*dominio.Propiedad, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dominio.Propiedad), args.Error(1)
}

func (m *MockPropiedadRepository) ObtenerTodas() ([]dominio.Propiedad, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dominio.Propiedad), args.Error(1)
}

func (m *MockPropiedadRepository) Actualizar(propiedad *dominio.Propiedad) error {
	args := m.Called(propiedad)
	return args.Error(0)
}

func (m *MockPropiedadRepository) Eliminar(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test Suite para PropiedadService
type PropiedadServiceTestSuite struct {
	suite.Suite
	service  *PropiedadService
	mockRepo *MockPropiedadRepository
}

func (suite *PropiedadServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockPropiedadRepository)
	suite.service = NuevoPropiedadService(suite.mockRepo)
}

func (suite *PropiedadServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Tests para CrearPropiedad
func (suite *PropiedadServiceTestSuite) TestCrearPropiedad_Exitoso() {
	// Configurar mock
	suite.mockRepo.On("Crear", mock.AnythingOfType("*dominio.Propiedad")).Return(nil)

	// Ejecutar
	propiedad, err := suite.service.CrearPropiedad(
		"Casa moderna en Samborondón",
		"Hermosa casa con piscina",
		"Guayas",
		"Samborondón",
		"casa",
		250000,
	)

	// Verificar
	suite.NoError(err)
	suite.NotNil(propiedad)
	suite.Equal("Casa moderna en Samborondón", propiedad.Titulo)
	suite.Equal("Guayas", propiedad.Provincia)
	suite.Equal("casa", propiedad.Tipo)
	suite.Equal(250000.0, propiedad.Precio)
	suite.Equal("disponible", propiedad.Estado)
}

func (suite *PropiedadServiceTestSuite) TestCrearPropiedad_TituloVacio() {
	// No configuramos mock porque no debería llegar al repositorio

	// Ejecutar
	propiedad, err := suite.service.CrearPropiedad(
		"", // Título vacío
		"Descripción",
		"Guayas",
		"Samborondón",
		"casa",
		250000,
	)

	// Verificar
	suite.Error(err)
	suite.Nil(propiedad)
	suite.Contains(err.Error(), "título es requerido")
}

func (suite *PropiedadServiceTestSuite) TestCrearPropiedad_TituloMuyCorto() {
	propiedad, err := suite.service.CrearPropiedad(
		"Casa", // Título muy corto (menos de 10 caracteres)
		"Descripción",
		"Guayas",
		"Samborondón",
		"casa",
		250000,
	)

	suite.Error(err)
	suite.Nil(propiedad)
	suite.Contains(err.Error(), "título debe tener al menos 10 caracteres")
}

func (suite *PropiedadServiceTestSuite) TestCrearPropiedad_PrecioInvalido() {
	tests := []struct {
		name   string
		precio float64
	}{
		{"precio cero", 0},
		{"precio negativo", -100000},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			propiedad, err := suite.service.CrearPropiedad(
				"Casa moderna en Samborondón",
				"Descripción",
				"Guayas",
				"Samborondón",
				"casa",
				tt.precio,
			)

			suite.Error(err)
			suite.Nil(propiedad)
			suite.Contains(err.Error(), "precio debe ser mayor a 0")
		})
	}
}

func (suite *PropiedadServiceTestSuite) TestCrearPropiedad_ProvinciaInvalida() {
	propiedad, err := suite.service.CrearPropiedad(
		"Casa moderna en Madrid",
		"Descripción",
		"Madrid", // Provincia que no existe en Ecuador
		"Madrid",
		"casa",
		250000,
	)

	suite.Error(err)
	suite.Nil(propiedad)
	suite.Contains(err.Error(), "provincia no válida")
}

func (suite *PropiedadServiceTestSuite) TestCrearPropiedad_TipoInvalido() {
	propiedad, err := suite.service.CrearPropiedad(
		"Casa moderna en Samborondón",
		"Descripción",
		"Guayas",
		"Samborondón",
		"mansion", // Tipo no válido
		250000,
	)

	suite.Error(err)
	suite.Nil(propiedad)
	suite.Contains(err.Error(), "tipo de propiedad no válido")
}

func (suite *PropiedadServiceTestSuite) TestCrearPropiedad_ErrorRepositorio() {
	// Configurar mock para retornar error
	suite.mockRepo.On("Crear", mock.AnythingOfType("*dominio.Propiedad")).Return(assert.AnError)

	propiedad, err := suite.service.CrearPropiedad(
		"Casa moderna en Samborondón",
		"Descripción",
		"Guayas",
		"Samborondón",
		"casa",
		250000,
	)

	suite.Error(err)
	suite.Nil(propiedad)
	suite.Contains(err.Error(), "error al crear propiedad")
}

// Tests para ObtenerPropiedad
func (suite *PropiedadServiceTestSuite) TestObtenerPropiedad_Exitoso() {
	propiedadEsperada := dominio.NuevaPropiedad(
		"Casa Test",
		"Descripción test",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	suite.mockRepo.On("ObtenerPorID", propiedadEsperada.ID).Return(propiedadEsperada, nil)

	propiedad, err := suite.service.ObtenerPropiedad(propiedadEsperada.ID)

	suite.NoError(err)
	suite.NotNil(propiedad)
	suite.Equal(propiedadEsperada.ID, propiedad.ID)
}

func (suite *PropiedadServiceTestSuite) TestObtenerPropiedad_IDVacio() {
	propiedad, err := suite.service.ObtenerPropiedad("")

	suite.Error(err)
	suite.Nil(propiedad)
	suite.Contains(err.Error(), "ID de propiedad requerido")
}

func (suite *PropiedadServiceTestSuite) TestObtenerPropiedad_NoEncontrada() {
	suite.mockRepo.On("ObtenerPorID", "id-inexistente").Return(nil, sql.ErrNoRows)

	propiedad, err := suite.service.ObtenerPropiedad("id-inexistente")

	suite.Error(err)
	suite.Nil(propiedad)
}

// Tests para ObtenerPropiedadPorSlug
func (suite *PropiedadServiceTestSuite) TestObtenerPropiedadPorSlug_Exitoso() {
	propiedadEsperada := dominio.NuevaPropiedad(
		"Casa Test",
		"Descripción test",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	suite.mockRepo.On("ObtenerPorSlug", propiedadEsperada.Slug).Return(propiedadEsperada, nil)

	propiedad, err := suite.service.ObtenerPropiedadPorSlug(propiedadEsperada.Slug)

	suite.NoError(err)
	suite.NotNil(propiedad)
	suite.Equal(propiedadEsperada.Slug, propiedad.Slug)
}

func (suite *PropiedadServiceTestSuite) TestObtenerPropiedadPorSlug_SlugVacio() {
	propiedad, err := suite.service.ObtenerPropiedadPorSlug("")

	suite.Error(err)
	suite.Nil(propiedad)
	suite.Contains(err.Error(), "slug de propiedad requerido")
}

func (suite *PropiedadServiceTestSuite) TestObtenerPropiedadPorSlug_SlugInvalido() {
	propiedad, err := suite.service.ObtenerPropiedadPorSlug("Slug Con Espacios")

	suite.Error(err)
	suite.Nil(propiedad)
	suite.Contains(err.Error(), "formato de slug inválido")
}

// Tests para ListarPropiedades
func (suite *PropiedadServiceTestSuite) TestListarPropiedades_Exitoso() {
	propiedadesEsperadas := []dominio.Propiedad{
		*dominio.NuevaPropiedad("Casa 1", "Desc 1", "Pichincha", "Quito", "casa", 100000),
		*dominio.NuevaPropiedad("Casa 2", "Desc 2", "Guayas", "Guayaquil", "casa", 200000),
	}

	suite.mockRepo.On("ObtenerTodas").Return(propiedadesEsperadas, nil)

	propiedades, err := suite.service.ListarPropiedades()

	suite.NoError(err)
	suite.Len(propiedades, 2)
}

func (suite *PropiedadServiceTestSuite) TestListarPropiedades_ErrorRepositorio() {
	suite.mockRepo.On("ObtenerTodas").Return(nil, assert.AnError)

	propiedades, err := suite.service.ListarPropiedades()

	suite.Error(err)
	suite.Nil(propiedades)
}

// Tests para ActualizarPropiedad
func (suite *PropiedadServiceTestSuite) TestActualizarPropiedad_Exitoso() {
	propiedadExistente := dominio.NuevaPropiedad(
		"Casa Original",
		"Descripción original",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	suite.mockRepo.On("ObtenerPorID", propiedadExistente.ID).Return(propiedadExistente, nil)
	suite.mockRepo.On("Actualizar", mock.AnythingOfType("*dominio.Propiedad")).Return(nil)

	propiedad, err := suite.service.ActualizarPropiedad(
		propiedadExistente.ID,
		"Casa Actualizada",
		"Descripción actualizada",
		"Guayas",
		"Guayaquil",
		"casa",
		180000,
	)

	suite.NoError(err)
	suite.NotNil(propiedad)
	suite.Equal("Casa Actualizada", propiedad.Titulo)
	suite.Equal(180000.0, propiedad.Precio)
}

func (suite *PropiedadServiceTestSuite) TestActualizarPropiedad_NoEncontrada() {
	suite.mockRepo.On("ObtenerPorID", "id-inexistente").Return(nil, sql.ErrNoRows)

	propiedad, err := suite.service.ActualizarPropiedad(
		"id-inexistente",
		"Casa Actualizada",
		"Descripción",
		"Pichincha",
		"Quito",
		"casa",
		180000,
	)

	suite.Error(err)
	suite.Nil(propiedad)
}

// Tests para EliminarPropiedad
func (suite *PropiedadServiceTestSuite) TestEliminarPropiedad_Exitoso() {
	propiedadExistente := dominio.NuevaPropiedad(
		"Casa a Eliminar",
		"Descripción",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	suite.mockRepo.On("ObtenerPorID", propiedadExistente.ID).Return(propiedadExistente, nil)
	suite.mockRepo.On("Eliminar", propiedadExistente.ID).Return(nil)

	err := suite.service.EliminarPropiedad(propiedadExistente.ID)

	suite.NoError(err)
}

func (suite *PropiedadServiceTestSuite) TestEliminarPropiedad_IDVacio() {
	err := suite.service.EliminarPropiedad("")

	suite.Error(err)
	suite.Contains(err.Error(), "ID de propiedad requerido")
}

func (suite *PropiedadServiceTestSuite) TestEliminarPropiedad_NoEncontrada() {
	suite.mockRepo.On("ObtenerPorID", "id-inexistente").Return(nil, sql.ErrNoRows)

	err := suite.service.EliminarPropiedad("id-inexistente")

	suite.Error(err)
}

// Tests para FiltrarPorProvincia
func (suite *PropiedadServiceTestSuite) TestFiltrarPorProvincia_Exitoso() {
	todasLasPropiedades := []dominio.Propiedad{
		*dominio.NuevaPropiedad("Casa 1", "Desc 1", "Pichincha", "Quito", "casa", 100000),
		*dominio.NuevaPropiedad("Casa 2", "Desc 2", "Guayas", "Guayaquil", "casa", 200000),
		*dominio.NuevaPropiedad("Casa 3", "Desc 3", "Pichincha", "Quito", "casa", 150000),
	}

	suite.mockRepo.On("ObtenerTodas").Return(todasLasPropiedades, nil)

	propiedades, err := suite.service.FiltrarPorProvincia("Pichincha")

	suite.NoError(err)
	suite.Len(propiedades, 2) // Solo las de Pichincha
	for _, propiedad := range propiedades {
		suite.Equal("Pichincha", propiedad.Provincia)
	}
}

func (suite *PropiedadServiceTestSuite) TestFiltrarPorProvincia_ProvinciaVacia() {
	propiedades, err := suite.service.FiltrarPorProvincia("")

	suite.Error(err)
	suite.Nil(propiedades)
	suite.Contains(err.Error(), "provincia requerida")
}

func (suite *PropiedadServiceTestSuite) TestFiltrarPorProvincia_ProvinciaInvalida() {
	propiedades, err := suite.service.FiltrarPorProvincia("Madrid")

	suite.Error(err)
	suite.Nil(propiedades)
	suite.Contains(err.Error(), "provincia no válida")
}

// Tests para FiltrarPorRangoPrecio
func (suite *PropiedadServiceTestSuite) TestFiltrarPorRangoPrecio_Exitoso() {
	todasLasPropiedades := []dominio.Propiedad{
		*dominio.NuevaPropiedad("Casa 1", "Desc 1", "Pichincha", "Quito", "casa", 80000),
		*dominio.NuevaPropiedad("Casa 2", "Desc 2", "Guayas", "Guayaquil", "casa", 150000),
		*dominio.NuevaPropiedad("Casa 3", "Desc 3", "Pichincha", "Quito", "casa", 250000),
	}

	suite.mockRepo.On("ObtenerTodas").Return(todasLasPropiedades, nil)

	propiedades, err := suite.service.FiltrarPorRangoPrecio(100000, 200000)

	suite.NoError(err)
	suite.Len(propiedades, 1) // Solo la de 150000
	suite.Equal(150000.0, propiedades[0].Precio)
}

func (suite *PropiedadServiceTestSuite) TestFiltrarPorRangoPrecio_PreciosNegativos() {
	propiedades, err := suite.service.FiltrarPorRangoPrecio(-100000, 200000)

	suite.Error(err)
	suite.Nil(propiedades)
	suite.Contains(err.Error(), "los precios deben ser positivos")
}

func (suite *PropiedadServiceTestSuite) TestFiltrarPorRangoPrecio_MinMayorQueMax() {
	propiedades, err := suite.service.FiltrarPorRangoPrecio(200000, 100000)

	suite.Error(err)
	suite.Nil(propiedades)
	suite.Contains(err.Error(), "precio mínimo no puede ser mayor que precio máximo")
}

// Tests para ObtenerEstadisticas
func (suite *PropiedadServiceTestSuite) TestObtenerEstadisticas_Exitoso() {
	todasLasPropiedades := []dominio.Propiedad{
		*dominio.NuevaPropiedad("Casa 1", "Desc 1", "Pichincha", "Quito", "casa", 100000),
		*dominio.NuevaPropiedad("Depto 1", "Desc 2", "Guayas", "Guayaquil", "departamento", 80000),
		*dominio.NuevaPropiedad("Casa 2", "Desc 3", "Pichincha", "Quito", "casa", 120000),
	}

	// Cambiar el estado de una propiedad para tener variedad
	todasLasPropiedades[1].Estado = "vendida"

	suite.mockRepo.On("ObtenerTodas").Return(todasLasPropiedades, nil)

	stats, err := suite.service.ObtenerEstadisticas()

	suite.NoError(err)
	suite.NotNil(stats)
	suite.Equal(3, stats["total_propiedades"])
	suite.Equal(100000.0, stats["precio_promedio"]) // (100000 + 80000 + 120000) / 3

	porTipo := stats["por_tipo"].(map[string]int)
	suite.Equal(2, porTipo["casa"])
	suite.Equal(1, porTipo["departamento"])

	porEstado := stats["por_estado"].(map[string]int)
	suite.Equal(2, porEstado["disponible"])
	suite.Equal(1, porEstado["vendida"])
}

func (suite *PropiedadServiceTestSuite) TestObtenerEstadisticas_SinPropiedades() {
	suite.mockRepo.On("ObtenerTodas").Return([]dominio.Propiedad{}, nil)

	stats, err := suite.service.ObtenerEstadisticas()

	suite.NoError(err)
	suite.NotNil(stats)
	suite.Equal(0, stats["total_propiedades"])
	suite.Equal(float64(0), stats["precio_promedio"])
}

// Ejecutar la suite de tests
func TestPropiedadServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PropiedadServiceTestSuite))
}

// Tests unitarios adicionales

// TestNuevoPropiedadService verifica la creación del servicio
func TestNuevoPropiedadService(t *testing.T) {
	mockRepo := new(MockPropiedadRepository)
	service := NuevoPropiedadService(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

// Test de integración usando el mock simple del repositorio
func TestIntegracionServicioConMockSimple(t *testing.T) {
	// Usar el mock simple del archivo de repositorio para tests de integración
	// (Este requeriría importar el mock del package repositorio)
	// Por ahora, solo verificamos que el servicio funciona con cualquier implementación de la interfaz

	mockRepo := new(MockPropiedadRepository)
	service := NuevoPropiedadService(mockRepo)

	assert.NotNil(t, service)
	assert.Implements(t, (*repositorio.PropiedadRepository)(nil), mockRepo)
}

// Benchmark para CrearPropiedad
func BenchmarkPropiedadService_CrearPropiedad(b *testing.B) {
	mockRepo := new(MockPropiedadRepository)
	mockRepo.On("Crear", mock.AnythingOfType("*dominio.Propiedad")).Return(nil)

	service := NuevoPropiedadService(mockRepo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CrearPropiedad(
			"Casa de Benchmark",
			"Descripción para benchmark",
			"Pichincha",
			"Quito",
			"casa",
			150000,
		)
	}
}

// Test de validación con datos reales ecuatorianos
func TestValidacionDatosEcuatorianos(t *testing.T) {
	mockRepo := new(MockPropiedadRepository)
	mockRepo.On("Crear", mock.AnythingOfType("*dominio.Propiedad")).Return(nil)

	service := NuevoPropiedadService(mockRepo)

	tests := []struct {
		name      string
		provincia string
		esperado  bool
	}{
		{"Pichincha", "Pichincha", true},
		{"Guayas", "Guayas", true},
		{"Manabí", "Manabí", true},
		{"Galápagos", "Galápagos", true},
		{"Provincia inexistente", "Cataluña", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CrearPropiedad(
				"Casa de prueba en "+tt.provincia,
				"Descripción de prueba",
				tt.provincia,
				"Ciudad",
				"casa",
				100000,
			)

			if tt.esperado {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "provincia no válida")
			}
		})
	}

	mockRepo.AssertExpectations(t)
}
